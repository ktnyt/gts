package gts

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	ascii "gopkg.in/ascii.v1"
	pars "gopkg.in/pars.v2"
)

// Location represents a location in a sequence as defined by the INSDC feature
// table definition.
type Location interface {
	fmt.Stringer
	Len() int
	Shift(i, n int) Location
	Less(loc Location) bool
}

// Locatable represents a location that can locate a region within a sequence.
type Locatable interface {
	Location
	Complement() Locatable
	Locate(seq Sequence) Sequence
}

// LocatableParser attempts to parse some Locatable location.
var LocatableParser pars.Parser

// LocationParser attempts to parse some location.
var LocationParser pars.Parser

func init() {
	LocatableParser = pars.Any(
		RangeParser,
		BetweenParser,
		ComplementParser,
		JoinParser,
		PointParser,
	)

	LocationParser = pars.Any(
		RangeParser,
		BetweenParser,
		AmbiguousParser,
		ComplementParser,
		JoinParser,
		OrderParser,
		PointParser,
	)
}

func shift(pos, i, n int, closed bool) int {
	flag := i < pos
	if closed {
		flag = i <= pos
	}
	if flag {
		pos += n
		if pos < 0 {
			return 0
		}
	}
	return pos
}

// Between represents a position between two bases. This will only make logical
// sense if the start and end positions are directly adjacent.
type Between int

// String satisfies the fmt.Stringer interface.
func (between Between) String() string {
	return fmt.Sprintf("%d^%d", between, between+1)
}

// Len returns the total length spanned by the location.
func (between Between) Len() int {
	return 0
}

// Shift the location beyond the given position i by n.
func (between Between) Shift(i, n int) Location {
	return Between(shift(int(between), i, n, false))
}

// Less returns true if the location is less than the given location.
func (between Between) Less(loc Location) bool {
	switch v := loc.(type) {
	case Between:
		return int(between) < int(v)
	case Point:
		return int(between) <= int(v)
	case Ranged:
		return int(between) <= v.Start
	case Complemented:
		return between.Less(v[0])
	case Joined:
		for _, u := range v {
			if !between.Less(u) {
				return false
			}
		}
		return true
	case Ambiguous:
		return int(between) <= v[0]
	case Ordered:
		for _, u := range v {
			if !between.Less(u) {
				return false
			}
		}
		return true
	default:
		return true
	}
}

// Complement returns the complement location.
func (between Between) Complement() Locatable {
	return Complemented{between}
}

// Locate the sequence region represented by the location.
func (between Between) Locate(seq Sequence) Sequence {
	return New(seq.Info(), []byte{})
}

// BetweenParser will attempt to parse a Between loctation.
func BetweenParser(state *pars.State, result *pars.Result) error {
	state.Push()
	if err := pars.Int(state, result); err != nil {
		state.Pop()
		return err
	}
	start := result.Value.(int)
	c, err := pars.Next(state)
	if err != nil {
		state.Pop()
		return err
	}
	if c != '^' {
		err := pars.NewError("expected `^`", state.Position())
		state.Pop()
		return err
	}
	state.Advance()
	if err := pars.Int(state, result); err != nil {
		state.Pop()
		return err
	}
	end := result.Value.(int)
	if start+1 != end {
		return fmt.Errorf("%d^%d is not a valid location: coordinates should be adjacent", start, end)
	}
	result.SetValue(Between(start))
	state.Drop()
	return nil
}

// Point represents a single base position in a sequence.
type Point int

// String satisfies the fmt.Stringer interface.
func (point Point) String() string {
	return strconv.Itoa(int(point + 1))
}

// Len returns the total length spanned by the location.
func (point Point) Len() int {
	return 1
}

// Shift the location beyond the given position i by n.
func (point Point) Shift(i, n int) Location {
	pos := int(point)
	if n < 0 && i <= pos && pos+1 <= i-n {
		return Between(i)
	}
	return Point(shift(pos, i, n, true))
}

// Less returns true if the location is less than the given location.
func (point Point) Less(loc Location) bool {
	switch v := loc.(type) {
	case Between:
		return int(point) < int(v)
	case Point:
		return int(point) < int(v)
	case Ranged:
		switch {
		case int(point) < v.Start:
			return true
		default:
			return false
		}
	case Complemented:
		return point.Less(v[0])
	case Joined:
		for _, u := range v {
			if !point.Less(u) {
				return false
			}
		}
		return true
	case Ambiguous:
		return int(point) < v[0]
	case Ordered:
		for _, u := range v {
			if !point.Less(u) {
				return false
			}
		}
		return true
	default:
		return true
	}
}

// Complement returns the complement location.
func (point Point) Complement() Locatable {
	return Complemented{point}
}

// Locate the sequence region represented by the location.
func (point Point) Locate(seq Sequence) Sequence {
	return Slice(seq, int(point), int(point)+1)
}

// PointParser will attempt to parse a Point loctation.
var PointParser = pars.Parser(pars.Int).Map(func(result *pars.Result) error {
	point := result.Value.(int)
	result.SetValue(Point(point - 1))
	return nil
})

// Partial represents the partiality of a location range.
type Partial [2]bool

// Partiality values.
var (
	Complete    Partial = [2]bool{false, false}
	Partial5    Partial = [2]bool{true, false}
	Partial3    Partial = [2]bool{false, true}
	PartialBoth Partial = [2]bool{true, true}
)

// Ranged represents a contiguous region of bases in a sequence. The starting
// and ending positions of a Ranged may be partial.
type Ranged struct {
	Start   int
	End     int
	Partial Partial
}

// PartialRange returns the range between the start and end positions where the
// specified ends are partial. They can be Complete, Partial5, Partial3, or
// PartialBoth.
func PartialRange(start, end int, partial Partial) Locatable {
	if end <= start {
		panic(fmt.Errorf("Ranged bounds out of range [%d:%d]", start, end))
	}
	/* DISCUSS: should a complete, one base range be reduced to a Point?
	if partial == Complete && start+1 == end {
		return Point(start)
	}
	*/
	return Ranged{start, end, partial}
}

// Range returns the complete range between the start and end positions.
func Range(start, end int) Locatable {
	return PartialRange(start, end, Complete)
}

// String satisfies the fmt.Stringer interface.
func (ranged Ranged) String() string {
	b := strings.Builder{}
	if ranged.Partial[0] {
		b.WriteByte('<')
	}
	b.WriteString(strconv.Itoa(ranged.Start + 1))
	b.WriteString("..")
	if ranged.Partial[1] {
		b.WriteByte('>')
	}
	b.WriteString(strconv.Itoa(ranged.End))
	return b.String()
}

// Len returns the total length spanned by the location.
func (ranged Ranged) Len() int {
	return ranged.End - ranged.Start
}

// Shift the location beyond the given position i by n.
func (ranged Ranged) Shift(i, n int) Location {
	if n < 0 && i <= ranged.Start && ranged.End <= i-n {
		return Between(i)
	}
	start, end := shift(ranged.Start, i, n, true), shift(ranged.End, i, n, false)
	return PartialRange(start, end, ranged.Partial)
}

// Less returns true if the location is less than the given location.
func (ranged Ranged) Less(loc Location) bool {
	switch v := loc.(type) {
	case Between:
		return ranged.Start < int(v)
	case Point:
		return ranged.Start <= int(v)
	case Ranged:
		switch {
		case ranged.Start < v.Start:
			return true
		case v.Start < ranged.Start:
			return false
		case ranged.Partial[0] && !v.Partial[0]:
			return true
		case !ranged.Partial[0] && v.Partial[0]:
			return false
		case ranged.End < v.End:
			return true
		case v.End < ranged.End:
			return false
		case !ranged.Partial[1] && v.Partial[1]:
			return true
		case ranged.Partial[1] && !v.Partial[1]:
			return false
		default:
			return false
		}
	case Complemented:
		return ranged.Less(v[0])
	case Joined:
		for _, u := range v {
			if !ranged.Less(u) {
				return false
			}
		}
		return true
	case Ambiguous:
		switch {
		case ranged.Start < v[0]:
			return true
		case v[0] < ranged.Start:
			return false
		case ranged.Partial[0]:
			return true
		case ranged.End < v[1]:
			return true
		case v[1] < ranged.End:
			return false
		case ranged.Partial[1]:
			return false
		default:
			return false
		}
	case Ordered:
		for _, u := range v {
			if !ranged.Less(u) {
				return false
			}
		}
		return true
	default:
		return true
	}
}

// Complement returns the complement location.
func (ranged Ranged) Complement() Locatable {
	return Complemented{ranged}
}

// Locate the sequence region represented by the location.
func (ranged Ranged) Locate(seq Sequence) Sequence {
	return Slice(seq, ranged.Start, ranged.End)
}

// RangeParser attempts to parse a Ranged location.
func RangeParser(state *pars.State, result *pars.Result) error {
	state.Push()
	c, err := pars.Next(state)
	if err != nil {
		state.Pop()
		return err
	}
	partial5 := false
	if c == '<' {
		partial5 = true
		state.Advance()
	}
	if err := pars.Int(state, result); err != nil {
		state.Pop()
		return err
	}
	start := result.Value.(int) - 1
	if err := state.Request(2); err != nil {
		state.Pop()
		return err
	}
	if !bytes.Equal(state.Buffer(), []byte("..")) {
		err := pars.NewError("expected `..`", state.Position())
		state.Pop()
		return err
	}
	state.Advance()
	c, err = pars.Next(state)
	partial3 := false
	if c == '>' {
		partial3 = true
		state.Advance()
	}
	if err := pars.Int(state, result); err != nil {
		state.Pop()
		return err
	}
	end := result.Value.(int)

	// Some legacy entries have the partial marker in the end.
	c, err = pars.Next(state)
	if err == nil && c == '>' {
		partial3 = true
		state.Advance()
	}
	result.SetValue(Ranged{start, end, [2]bool{partial5, partial3}})
	state.Drop()
	return nil
}

// Complemented represents a location complemented for the given molecule type.
type Complemented [1]Locatable

// String satisfies the fmt.Stringer interface.
func (complement Complemented) String() string {
	return fmt.Sprintf("complement(%s)", complement[0])
}

// Len returns the total length spanned by the location.
func (complement Complemented) Len() int {
	return complement[0].Len()
}

// Shift the location beyond the given position i by n.
func (complement Complemented) Shift(i, n int) Location {
	return Complemented{complement[0].Shift(i, n).(Locatable)}
}

// Less returns true if the location is less than the given location.
func (complement Complemented) Less(loc Location) bool {
	return complement[0].Less(loc)
}

// Complement returns the complement location.
func (complement Complemented) Complement() Locatable {
	return complement[0]
}

// Locate the sequence region represented by the location.
func (complement Complemented) Locate(seq Sequence) Sequence {
	return Reverse(Complement(complement[0].Locate(seq)))
}

// ComplementParser attempts to parse a Complement location.
func ComplementParser(state *pars.State, result *pars.Result) error {
	state.Push()
	if err := state.Request(11); err != nil {
		state.Pop()
		return err
	}
	if !bytes.Equal(state.Buffer(), []byte("complement(")) {
		err := pars.NewError("expected `complement(`", state.Position())
		state.Pop()
		return err
	}
	state.Advance()
	if err := LocatableParser(state, result); err != nil {
		state.Pop()
		return err
	}
	c, err := pars.Next(state)
	if err != nil {
		state.Pop()
		return err
	}
	if c != ')' {
		err := pars.NewError("expected `)`", state.Position())
		state.Pop()
		return err
	}
	state.Advance()
	result.SetValue(result.Value.(Locatable).Complement())
	state.Drop()
	return nil
}

// LocatableList represents a singly linked list of Locatable objects.
type LocatableList struct {
	Data Locatable
	Next *LocatableList
}

// Len returns the length of the list.
func (ll *LocatableList) Len() int {
	if ll.Next == nil {
		if ll.Data == nil {
			return 0
		}
		return 1
	}
	return ll.Next.Len() + 1
}

// Slice returns the slice representation of the list.
func (ll *LocatableList) Slice() []Locatable {
	list := []Locatable{ll.Data}
	if ll.Next == nil {
		return list
	}
	return append(list, ll.Next.Slice()...)
}

// Push a Locatable object to the end of the list. If the Locatable object is
// equivalent to the last element, nothing happens. If the Locatable object can
// be joined with the last element to form a contiguous Locatable location, the
// last element will be replaced with the joined Locatable object.
func (ll *LocatableList) Push(loc Locatable) {
	if ll.Next != nil {
		ll.Next.Push(loc)
		return
	}

	if joined, ok := loc.(Joined); ok {
		for i := range joined {
			ll.Push(joined[i])
		}
		return
	}

	if ll.Data == nil {
		ll.Data = loc
		return
	}

	switch v := ll.Data.(type) {
	case Point:
		switch u := loc.(type) {
		case Point:
			if v == u {
				return
			}
		case Ranged:
			if int(v) == u.Start {
				ll.Data = u
				return
			}
		}

	case Ranged:
		switch u := loc.(type) {
		case Point:
			if v.End == int(u) {
				return
			}
		case Ranged:
			if v.End == u.Start {
				ll.Data = Ranged{v.Start, u.End, v.Partial}
				return
			}
		}

	case Complemented:
		if u, ok := loc.(Complemented); ok {
			tmp := LocatableList{u[0], nil}
			tmp.Push(v[0])
			ll.Data = Complemented{Join(tmp.Slice()...)}
		}
		return
	}

	ll.Next = &LocatableList{loc, nil}
}

// Joined represents a list of Locatable locations. It is strongly recommended
// this be constructed using the Join helper function to reduce the list of
// Locatable locations to the simplest representation.
type Joined []Locatable

// Join the given Locatable locations. Will panic if no argument is given. The
// locations will first be reduced to the simplest representation by merging
// adjacent identical locations and contiguous locations. If the resulting list
// of locations have only one element, the elemnt will be returuned. Otherwise,
// a Joined object will be returned.
func Join(locs ...Locatable) Locatable {
	list := LocatableList{}
	for _, loc := range locs {
		list.Push(loc)
	}

	switch list.Len() {
	case 0:
		panic("Join without arguments is not allowed")
	case 1:
		return list.Data
	default:
		return Joined(list.Slice())
	}
}

// String satisfies the fmt.Stringer interface.
func (joined Joined) String() string {
	tmp := make([]string, len(joined))
	for i, loc := range joined {
		tmp[i] = loc.String()
	}
	return fmt.Sprintf("join(%s)", strings.Join(tmp, ","))
}

// Len returns the total length spanned by the location.
func (joined Joined) Len() int {
	n := 0
	for _, loc := range joined {
		n += loc.Len()
	}
	return n
}

// Shift the location beyond the given position i by n.
func (joined Joined) Shift(i, n int) Location {
	locs := make([]Locatable, len(joined))
	for j, loc := range joined {
		locs[j] = loc.Shift(i, n).(Locatable)
	}
	return Join(locs...)
}

// Less returns true if the location is less than the given location.
func (joined Joined) Less(loc Location) bool {
	for _, v := range joined {
		if v.Less(loc) {
			return true
		}
	}
	return false
}

// Complement returns the complement location.
func (joined Joined) Complement() Locatable {
	return Complemented{joined}
}

// Locate the sequence region represented by the location.
func (joined Joined) Locate(seq Sequence) Sequence {
	p := make([]byte, joined.Len())
	n := 0
	for _, loc := range joined {
		n += copy(p[n:], loc.Locate(seq).Bytes())
	}
	return New(seq.Info(), p)
}

func locationDelimiter(state *pars.State, result *pars.Result) bool {
	state.Push()
	c, err := pars.Next(state)
	if err != nil {
		state.Pop()
		return false
	}
	if c != ',' {
		state.Pop()
		return false
	}
	state.Advance()
	c, err = pars.Next(state)
	for ascii.IsSpace(c) && err == nil {
		state.Advance()
		c, err = pars.Next(state)
	}
	state.Drop()
	return true
}

func multipleLocatableParser(state *pars.State, result *pars.Result) error {
	state.Push()
	if err := LocatableParser(state, result); err != nil {
		state.Pop()
		return err
	}
	locs := []Locatable{result.Value.(Locatable)}
	for locationDelimiter(state, result) {
		if err := LocatableParser(state, result); err != nil {
			state.Pop()
			return err
		}
		locs = append(locs, result.Value.(Locatable))
	}
	result.SetValue(locs)
	state.Drop()
	return nil
}

// JoinParser attempts to parse a Joined location.
func JoinParser(state *pars.State, result *pars.Result) error {
	state.Push()
	if err := state.Request(5); err != nil {
		state.Pop()
		return err
	}
	if !bytes.Equal(state.Buffer(), []byte("join(")) {
		err := pars.NewError("expected `join(`", state.Position())
		state.Pop()
		return err
	}
	state.Advance()
	if err := multipleLocatableParser(state, result); err != nil {
		return err
	}
	c, err := pars.Next(state)
	if err != nil {
		state.Pop()
		return err
	}
	if c != ')' {
		err := pars.NewError("expected `)`", state.Position())
		state.Pop()
		return err
	}
	state.Advance()
	result.SetValue(Join(result.Value.([]Locatable)...))
	state.Drop()
	return nil
}

// Ambiguous represents a single base within a given range.
type Ambiguous [2]int

// String satisfies the fmt.Stringer interface.
func (ambiguous Ambiguous) String() string {
	return fmt.Sprintf("%d.%d", ambiguous[0]+1, ambiguous[1])
}

// Len returns the total length spanned by the location.
func (ambiguous Ambiguous) Len() int {
	return 1
}

// Shift the location beyond the given position i by n.
func (ambiguous Ambiguous) Shift(i, n int) Location {
	if n < 0 && i <= ambiguous[0] && ambiguous[1] <= i-n {
		return Between(i)
	}
	return Ambiguous{
		shift(ambiguous[0], i, n, true),
		shift(ambiguous[1], i, n, false),
	}
}

// Less returns true if the location is less than the given location.
func (ambiguous Ambiguous) Less(loc Location) bool {
	switch v := loc.(type) {
	case Between:
		return ambiguous[0] < int(v)
	case Point:
		return ambiguous[0] <= int(v)
	case Ranged:
		switch {
		case ambiguous[0] < v.Start:
			return true
		case v.Start < ambiguous[0]:
			return false
		case v.Partial[0]:
			return false
		case ambiguous[1] < v.End:
			return true
		case v.End < ambiguous[1]:
			return false
		case v.Partial[1]:
			return true
		default:
			return false
		}
	case Complemented:
		return ambiguous.Less(v[0])
	case Joined:
		for _, u := range v {
			if !ambiguous.Less(u) {
				return false
			}
		}
		return true
	case Ambiguous:
		switch {
		case ambiguous[0] < v[0]:
			return true
		case v[0] < ambiguous[0]:
			return false
		case ambiguous[1] < v[1]:
			return true
		default:
			return false
		}
	case Ordered:
		for _, u := range v {
			if !ambiguous.Less(u) {
				return false
			}
		}
		return true
	default:
		return true
	}
}

// AmbiguousParser attempts to parse a Ambiguous location.
func AmbiguousParser(state *pars.State, result *pars.Result) error {
	state.Push()
	if err := pars.Int(state, result); err != nil {
		state.Pop()
		return err
	}
	start := result.Value.(int) - 1
	c, err := pars.Next(state)
	if err != nil {
		state.Pop()
		return err
	}
	if c != '.' {
		err := pars.NewError("expected `.`", state.Position())
		state.Pop()
		return err
	}
	state.Advance()
	if err := pars.Int(state, result); err != nil {
		state.Pop()
		return err
	}
	end := result.Value.(int)
	result.SetValue(Ambiguous{start, end})
	state.Drop()
	return nil
}

// Ordered represents multiple locations.
type Ordered []Location

func flattenLocations(locs []Location) []Location {
	list := []Location{}
	for i := range locs {
		switch loc := locs[i].(type) {
		case Ordered:
			list = append(list, flattenLocations([]Location(loc))...)
		default:
			list = append(list, loc)
		}
	}
	return list
}

// Order takes the given Locations and returns an Ordered containing the
// simplest form.
func Order(locs ...Location) Location {
	list := flattenLocations(locs)
	switch len(list) {
	case 0:
		panic("Order without arguments is not allowed")
	case 1:
		return list[0]
	default:
		return Ordered(list)
	}
}

// String satisfies the fmt.Stringer interface.
func (ordered Ordered) String() string {
	tmp := make([]string, len(ordered))
	for i, loc := range ordered {
		tmp[i] = loc.String()
	}
	return fmt.Sprintf("order(%s)", strings.Join(tmp, ","))
}

// Len returns the total length spanned by the location.
func (ordered Ordered) Len() int {
	n := 0
	for _, loc := range ordered {
		n += loc.Len()
	}
	return n
}

// Shift the location beyond the given position i by n.
func (ordered Ordered) Shift(i, n int) Location {
	locs := make([]Location, len(ordered))
	for j, loc := range ordered {
		locs[j] = loc.Shift(i, n)
	}
	return Order(locs...)
}

// Less returns true if the location is less than the given location.
func (ordered Ordered) Less(loc Location) bool {
	for _, v := range ordered {
		if v.Less(loc) {
			return true
		}
	}
	return false
}

func multipleLocationParser(state *pars.State, result *pars.Result) error {
	state.Push()
	if err := LocationParser(state, result); err != nil {
		state.Pop()
		return err
	}
	locs := []Location{result.Value.(Location)}
	for locationDelimiter(state, result) {
		if err := LocationParser(state, result); err != nil {
			state.Pop()
			return err
		}
		locs = append(locs, result.Value.(Location))
	}
	result.SetValue(locs)
	state.Drop()
	return nil
}

// OrderParser attempts to parse a Ordered location.
func OrderParser(state *pars.State, result *pars.Result) error {
	state.Push()
	if err := state.Request(6); err != nil {
		state.Pop()
		return err
	}
	if !bytes.Equal(state.Buffer(), []byte("order(")) {
		err := pars.NewError("expected `order(`", state.Position())
		state.Pop()
		return err
	}
	state.Advance()
	if err := multipleLocationParser(state, result); err != nil {
		return err
	}
	c, err := pars.Next(state)
	if err != nil {
		state.Pop()
		return err
	}
	if c != ')' {
		err := pars.NewError("expected `)`", state.Position())
		state.Pop()
		return err
	}
	state.Advance()
	result.SetValue(Order(result.Value.([]Location)...))
	state.Drop()
	return nil
}
