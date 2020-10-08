package seqio

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/go-ascii/ascii"
	"github.com/go-gts/gts"
	"github.com/go-pars/pars"
	"github.com/go-wrap/wrap"
)

// GenBankFields represents the fields of a GenBank record other than the
// features and sequence.
type GenBankFields struct {
	LocusName string
	Molecule  gts.Molecule
	Topology  gts.Topology
	Division  string
	Date      Date

	Definition string
	Accession  string
	Version    string
	DBLink     Dictionary
	Keywords   []string
	Source     Organism
	References []Reference
	Comment    string
	Contig     Contig

	Region gts.Region
}

// Slice returns a metadata sliced with the given region.
func (gbf GenBankFields) Slice(start, end int) interface{} {
	gbf.Region = gts.Segment{start, end}

	prefix := gbf.Molecule.Counter()
	parser := parseReferenceInfo(prefix)
	tryParse := func(info string) ([]gts.Segment, bool) {
		result, err := parser.Parse(pars.FromString(info))
		if err != nil {
			return nil, false
		}
		return result.Value.([]gts.Segment), true
	}

	refs := []Reference{}
	for _, ref := range gbf.References {
		info := ref.Info

		segs, ok := tryParse(info)
		switch {
		case ok:
			olap := []gts.Segment{}
			for _, seg := range segs {
				if seg.Overlap(start, end) {
					olap = append(olap, seg)
				}
			}
			if len(olap) > 0 {
				ss := make([]string, len(olap))
				for i, seg := range olap {
					head, tail := gts.Unpack(seg)
					head = gts.Max(0, head-start)
					tail = gts.Min(end-start, tail-start)
					ss[i] = fmt.Sprintf("%d to %d", head+1, tail)
				}
				ref.Info = fmt.Sprintf("(%s %s)", prefix, strings.Join(ss, "; "))
				refs = append(refs, ref)
			}
		default:
			refs = append(refs, ref)
		}
	}

	for i := range refs {
		refs[i].Number = i + 1
	}

	gbf.References = refs

	return gbf
}

// ID returns the ID of the sequence.
func (gbf GenBankFields) ID() string {
	if gbf.Version != "" {
		return gbf.Version
	}
	if gbf.Accession != "" {
		return gbf.Accession
	}
	return gbf.LocusName
}

// String satisifes the fmt.Stringer interface.
func (gbf GenBankFields) String() string {
	return fmt.Sprintf("%s %s", gbf.Version, gbf.Definition)
}

func toOriginLength(length int) int {
	lines := length / 60
	lastLine := length % 60
	blocks := lastLine / 10
	lastBlock := lastLine % 10
	ret := lines * 76
	if lastLine != 0 {
		ret += 10 + blocks*11
		if lastBlock != 0 {
			ret += lastBlock + 1
		}
	}
	return gts.Max(0, ret-1)
}

func fromOriginLength(length int) int {
	lines := length / 76
	lastLine := length%76 - 10
	blocks := lastLine / 11
	lastBlock := lastLine % 11
	ret := lines * 60
	if lastLine != 0 {
		ret += blocks*10 + lastBlock
	}
	return ret
}

type originIterator struct {
	line, block int
}

func (iter *originIterator) Next() (int, int) {
	i := iter.line*76 + iter.block*11 + 10
	j := iter.line*60 + iter.block*10
	iter.block++
	if iter.block == 6 {
		iter.block = 0
		iter.line++
	}
	return i, j
}

// Origin represents a GenBank sequence origin value.
type Origin struct {
	Buffer []byte
	Parsed bool
}

// NewOrigin formats a byte slice into GenBank sequence origin format.
func NewOrigin(p []byte) Origin {
	q := make([]byte, toOriginLength(len(p)))
	offset := 0
	for i := 0; i < len(p); i += 60 {
		if i != 0 {
			q[offset] = '\n'
			offset++
		}

		prefix := []byte(fmt.Sprintf("%9d", i+1))
		copy(q[offset:], prefix)
		offset += len(prefix)

		for j := 0; j < 60 && i+j < len(p); j += 10 {
			start := i + j
			end := gts.Min(start+10, len(p))

			q[offset] = ' '
			offset++

			copy(q[offset:], p[start:end])
			offset += end - start
		}
	}
	return Origin{q, false}
}

// Bytes converts the GenBank sequence origin into a byte slice.
func (o *Origin) Bytes() []byte {
	if !o.Parsed {
		if len(o.Buffer) < 10 {
			return nil
		}
		p := make([]byte, fromOriginLength(len(o.Buffer)))
		iter := originIterator{}
		i, j := iter.Next()

		for i < len(o.Buffer) {
			n := gts.Min(10, len(o.Buffer)-i)
			copy(p[j:j+n], o.Buffer[i:i+n])
			i, j = iter.Next()
		}
		o.Buffer = p
		o.Parsed = true
	}
	return o.Buffer
}

// String satisfies the fmt.Stringer interface.
func (o Origin) String() string {
	if !o.Parsed {
		return string(o.Buffer)
	}
	return string(NewOrigin(o.Buffer).Buffer)
}

// Len returns the actual sequence length.
func (o Origin) Len() int {
	if len(o.Buffer) == 0 {
		return 0
	}
	if o.Parsed {
		return len(o.Buffer)
	}
	return fromOriginLength(len(o.Buffer))
}

// GenBank represents a GenBank sequence record.
type GenBank struct {
	Fields GenBankFields
	Table  gts.FeatureTable
	Origin Origin
}

// NewGenBank creates a new GenBank object.
func NewGenBank(info GenBankFields, ff []gts.Feature, p []byte) GenBank {
	return GenBank{info, ff, NewOrigin(p)}
}

// Info returns the metadata of the sequence.
func (gb GenBank) Info() interface{} {
	return gb.Fields
}

// Features returns the feature table of the sequence.
func (gb GenBank) Features() gts.FeatureTable {
	return gb.Table
}

// Len returns the length of the sequence.
func (gb GenBank) Len() int {
	return gb.Origin.Len()
}

// Bytes returns the byte representation of the sequence.
func (gb GenBank) Bytes() []byte {
	return gb.Origin.Bytes()
}

// WithInfo creates a shallow copy of the given Sequence object and swaps the
// metadata with the given value.
func (gb GenBank) WithInfo(info interface{}) gts.Sequence {
	switch v := info.(type) {
	case GenBankFields:
		return GenBank{v, gb.Table, gb.Origin}
	default:
		return gts.New(v, gb.Features(), gb.Bytes())
	}
}

// WithFeatures creates a shallow copy of the given Sequence object and swaps
// the feature table with the given features.
func (gb GenBank) WithFeatures(ff gts.FeatureTable) gts.Sequence {
	return GenBank{gb.Fields, ff, gb.Origin}
}

// WithBytes creates a shallow copy of the given Sequence object and swaps the
// byte representation with the given byte slice.
func (gb GenBank) WithBytes(p []byte) gts.Sequence {
	return GenBank{gb.Fields, gb.Table, NewOrigin(p)}
}

// WithTopology creates a shallow copy of the given Sequence object and swaps
// the topology value with the given value.
func (gb GenBank) WithTopology(t gts.Topology) gts.Sequence {
	info := gb.Fields
	info.Topology = t
	return gb.WithInfo(info)
}

// String satisifes the fmt.Stringer interface.
func (gb GenBank) String() string {
	b := strings.Builder{}
	indent := "            "

	length := gb.Origin.Len()
	if length == 0 {
		length = gb.Fields.Contig.Region.Len()
	}

	date := strings.ToUpper(gb.Fields.Date.ToTime().Format("02-Jan-2006"))
	locus := fmt.Sprintf(
		"%-12s%-17s %10d bp %6s     %-9s%s %s", "LOCUS", gb.Fields.LocusName,
		length, gb.Fields.Molecule, gb.Fields.Topology, gb.Fields.Division, date,
	)

	b.WriteString(locus)

	definition := AddPrefix(gb.Fields.Definition, indent)
	b.WriteString("\nDEFINITION  " + definition + ".")
	b.WriteString("\nACCESSION   " + gb.Fields.Accession)
	if seg, ok := gb.Fields.Region.(gts.Segment); ok {
		head, tail := gts.Unpack(seg)
		b.WriteString(fmt.Sprintf(" REGION: %s", gts.Range(head, tail)))
	}
	b.WriteString("\nVERSION     " + gb.Fields.Version)

	for i, pair := range gb.Fields.DBLink {
		switch i {
		case 0:
			b.WriteString("\nDBLINK      ")
		default:
			b.WriteString("\n" + indent)
		}
		b.WriteString(fmt.Sprintf("%s: %s", pair.Key, pair.Value))
	}

	keywords := wrap.Space(strings.Join(gb.Fields.Keywords, "; ")+".", 67)
	keywords =
		AddPrefix(keywords, indent)
	b.WriteString("\nKEYWORDS    " + keywords)

	source := wrap.Space(gb.Fields.Source.Species, 67)
	source =
		AddPrefix(source, indent)
	b.WriteString("\nSOURCE      " + source)

	organism := wrap.Space(gb.Fields.Source.Name, 67)
	organism =
		AddPrefix(organism, indent)
	b.WriteString("\n  ORGANISM  " + organism)

	taxon := wrap.Space(strings.Join(gb.Fields.Source.Taxon, "; ")+".", 67)
	taxon = AddPrefix(taxon, indent)
	b.WriteString("\n" + indent + taxon)

	for _, ref := range gb.Fields.References {
		b.WriteString(fmt.Sprintf("\nREFERENCE   %d", ref.Number))
		if ref.Info != "" {
			pad := strings.Repeat(" ", 3-len(strconv.Itoa(ref.Number)))
			b.WriteString(pad + ref.Info)
		}
		if ref.Authors != "" {
			b.WriteString("\n  AUTHORS   " +
				AddPrefix(ref.Authors, indent))
		}
		if ref.Group != "" {
			b.WriteString("\n  CONSRTM   " +
				AddPrefix(ref.Group, indent))
		}
		if ref.Title != "" {
			b.WriteString("\n  TITLE     " +
				AddPrefix(ref.Title, indent))
		}
		if ref.Journal != "" {
			b.WriteString("\n  JOURNAL   " +
				AddPrefix(ref.Journal, indent))
		}
		if ref.Xref != nil {
			if v, ok := ref.Xref["PUBMED"]; ok {
				b.WriteString("\n   PUBMED   " + v)
			}
		}
		if ref.Comment != "" {
			b.WriteString("\n  REMARK    " +
				AddPrefix(ref.Comment, indent))
		}
	}

	if gb.Fields.Comment != "" {
		b.WriteString("\nCOMMENT     " +
			AddPrefix(gb.Fields.Comment, indent))
	}

	b.WriteString("\nFEATURES             Location/Qualifiers\n")

	gb.Table.Format("     ", 21).WriteTo(&b)

	if gb.Fields.Contig.String() != "" {
		b.WriteString(fmt.Sprintf("\nCONTIG      %s", gb.Fields.Contig))
	}

	if gb.Origin.Len() > 0 {
		b.WriteString("\nORIGIN      \n")
		b.WriteString(gb.Origin.String())
	}

	b.WriteString("\n//\n")

	return b.String()
}

// WriteTo satisfies the io.WriterTo interface.
func (gb GenBank) WriteTo(w io.Writer) (int64, error) {
	n, err := io.WriteString(w, gb.String())
	return int64(n), err
}

// GenBankFormatter implements the Formatter interface for GenBank files.
type GenBankFormatter struct {
	seq gts.Sequence
}

// WriteTo satisfies the io.WriterTo interface.
func (gf GenBankFormatter) WriteTo(w io.Writer) (int64, error) {
	switch seq := gf.seq.(type) {
	case GenBank:
		return seq.WriteTo(w)
	case *GenBank:
		return GenBankFormatter{*seq}.WriteTo(w)
	default:
		switch info := seq.Info().(type) {
		case GenBankFields:
			gb := NewGenBank(info, seq.Features(), seq.Bytes())
			return GenBankFormatter{gb}.WriteTo(w)
		default:
			return 0, fmt.Errorf("gts does not know how to format a sequence with metadata of type `%T` as GenBank", info)

		}
	}
}

var genbankLocusParser = pars.Seq(
	"LOCUS", pars.Spaces,
	pars.Word(ascii.Not(ascii.IsSpace)), pars.Spaces,
	pars.Int, " bp", pars.Spaces,
	pars.Word(ascii.Not(ascii.IsSpace)), pars.Spaces,
	pars.Word(ascii.Not(ascii.IsSpace)), pars.Spaces,
	pars.Count(pars.Byte(), 3).Map(pars.Cat), pars.Spaces,
	pars.AsParser(pars.Line).Map(func(result *pars.Result) (err error) {
		s := string(result.Token)
		date, err := AsDate(s)
		result.SetValue(date)
		return err
	}),
).Children(1, 2, 4, 7, 9, 11, 13)

func genbankFieldBodyParser(depth int) pars.Parser {
	indent := pars.String(strings.Repeat(" ", depth))
	return func(state *pars.State, result *pars.Result) error {
		pars.Line(state, result)
		tmp := *pars.NewTokenResult(result.Token)
		parser := pars.Many(pars.Seq(indent, pars.Line).Child(1))
		parser(state, result)
		children := append([]pars.Result{tmp}, result.Children...)
		result.SetChildren(children)
		return nil
	}
}

// GenBankParser attempts to parse a single GenBank record.
func GenBankParser(state *pars.State, result *pars.Result) error {
	if err := genbankLocusParser(state, result); err != nil {
		return err
	}

	pars.Cut(state, result)

	depth := len(result.Children[0].Token) + 5
	indent := pars.String(strings.Repeat(" ", depth))

	locus := string(result.Children[1].Token)
	length := result.Children[2].Value.(int)
	molecule, err := gts.AsMolecule(string(result.Children[3].Token))
	if err != nil {
		return pars.NewError(err.Error(), state.Position())
	}
	topology, err := gts.AsTopology(string(result.Children[4].Token))
	if err != nil {
		return pars.NewError(err.Error(), state.Position())
	}
	division := string(result.Children[5].Token)
	date := result.Children[6].Value.(Date)

	fields := GenBankFields{
		LocusName: locus,
		Molecule:  molecule,
		Topology:  topology,
		Division:  division,
		Date:      date,
		Region:    nil,
	}

	gb := GenBank{Fields: fields}

	fieldNameParser := pars.Word(ascii.IsUpper).Error(errors.New("expected field name"))
	fieldBodyParser := genbankFieldBodyParser(depth)
	end := pars.Seq("//", pars.EOL).Error(errors.New("expected end of record"))

	for {
		if end(state, result) == nil {
			result.SetValue(gb)
			return nil
		}
		if err := fieldNameParser(state, result); err != nil {
			return err
		}
		name := string(result.Token)
		paddingParser := pars.Count(' ', depth-len(name))

		if err := paddingParser(state, result); name != "ORIGIN" && err != nil {
			return pars.NewError("uneven indent", state.Position())
		}

		switch name {
		case "DEFINITION":
			parser := fieldBodyParser.Map(pars.Join([]byte("\n")))
			parser(state, result)
			token := bytes.TrimRight(result.Token, ".")
			gb.Fields.Definition = string(token)

		case "ACCESSION":
			pars.Line(state, result)
			gb.Fields.Accession = string(result.Token)

		case "VERSION":
			pars.Line(state, result)
			gb.Fields.Version = string(result.Token)

		case "DBLINK":
			headParser := pars.Seq(pars.Until(':'), ':', pars.Line).Children(0, 2)
			if err := headParser(state, result); err != nil {
				return err
			}
			db := string(result.Children[0].Token)
			id := string(result.Children[1].Token[1:])
			gb.Fields.DBLink.Set(db, id)

			tailParser := pars.Many(pars.Seq(indent, headParser).Child(1))
			tailParser(state, result)
			for _, child := range result.Children {
				db = string(child.Children[0].Token)
				id = string(child.Children[1].Token[1:])
				gb.Fields.DBLink.Set(db, id)
			}

		case "KEYWORDS":
			parser := fieldBodyParser.Map(pars.Join([]byte(" ")))
			parser(state, result)
			gb.Fields.Keywords =
				FlatFileSplit(string(result.Token))

		case "SOURCE":
			sourceParser := fieldBodyParser.Map(pars.Join([]byte("\n")))
			sourceParser(state, result)

			organism := Organism{}
			organism.Species = string(result.Token)

			organismLineParser := pars.Seq(
				pars.Spaces, []byte("ORGANISM"), pars.Spaces,
			).Map(pars.Cat)

			if organismLineParser(state, result) == nil {
				if len(result.Token) != depth {
					return pars.NewError("uneven indent", state.Position())
				}
				pars.Line(state, result)
				organism.Name = string(result.Token)

				taxonParser := pars.Many(
					pars.Seq(indent, pars.Line).Child(1),
				).Map(pars.Join([]byte(" ")))
				taxonParser(state, result)
				organism.Taxon =
					FlatFileSplit(string(result.Token))
			}

			gb.Fields.Source = organism

		case "REFERENCE":
			pars.Spaces(state, result)

			if err := pars.Int(state, result); err != nil {
				return pars.NewError("expected a reference number", state.Position())
			}

			number := result.Value.(int)

			reference := Reference{
				Number: number,
			}

			pad := strings.Repeat(" ", 3-len(strconv.Itoa(number)))
			infoParser := pars.Seq(pad, pars.Line).Child(1)
			if infoParser(state, result) == nil {
				reference.Info = string(result.Token)
			} else {
				pars.Line(state, result)
			}

			subfieldParser := pars.Seq(
				pars.Spaces,
				pars.Any(
					"AUTHORS",
					"CONSRTM",
					"TITLE",
					"JOURNAL",
					"PUBMED",
					"REMARK",
				),
				pars.Spaces,
			).Map(func(result *pars.Result) error {
				children := result.Children
				name := children[1].Value.(string)
				depth := len(name) + len(children[0].Token) + len(children[2].Token)
				*result = *pars.AsResults(name, depth)
				return nil
			})

			for subfieldParser(state, result) == nil {
				name := result.Children[0].Value.(string)
				if result.Children[1].Value.(int) != depth {
					return pars.NewError("uneven indent", state.Position())
				}
				parser := fieldBodyParser.Map(pars.Join([]byte("\n")))
				parser(state, result)
				switch name {
				case "AUTHORS":
					reference.Authors = string(result.Token)
				case "CONSRTM":
					reference.Group = string(result.Token)
				case "TITLE":
					reference.Title = string(result.Token)
				case "JOURNAL":
					reference.Journal = string(result.Token)
				case "PUBMED":
					reference.Xref = map[string]string{"PUBMED": string(result.Token)}
				case "REMARK":
					reference.Comment = string(result.Token)
				}
			}

			gb.Fields.References = append(gb.Fields.References, reference)

		case "COMMENT":
			parser := fieldBodyParser.Map(pars.Join([]byte("\n")))
			parser(state, result)
			gb.Fields.Comment = string(result.Token)

		case "FEATURES":
			pars.Line(state, result)
			parser := gts.FeatureTableParser("")
			if err := parser(state, result); err != nil {
				return err
			}
			gb.Table = gts.FeatureTable(result.Value.(gts.FeatureTable))

		case "CONTIG":
			contigParser := pars.Seq(
				"join(",
				pars.Until(':'), ':',
				pars.Int, "..", pars.Int,
				')', pars.Line,
			).Map(func(result *pars.Result) error {
				id := string(result.Children[1].Token)
				head := result.Children[3].Value.(int)
				tail := result.Children[5].Value.(int)
				contig := Contig{id, gts.Segment{head - 1, tail}}
				result.SetValue(contig)
				return nil
			})
			if err := contigParser(state, result); err != nil {
				return err
			}
			gb.Fields.Contig = result.Value.(Contig)

		case "ORIGIN":
			// Trim off excess whitespace.
			pars.Line(state, result)

			state.Push()
			if err := state.Request(toOriginLength(length) + 1); err != nil {
				return pars.NewError("not enough bytes in state", state.Position())
			}
			buffer := state.Buffer()
			state.Advance()
			gb.Origin = Origin{buffer[:len(buffer)-1], false}

		default:
			what := fmt.Sprintf("unexpected field name `%s`", name)
			return pars.NewError(what, state.Position())
		}
	}
}
