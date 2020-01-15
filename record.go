package gts

import (
	"fmt"
	"io"
	"os"

	pars "gopkg.in/pars.v2"
)

// Record is the interface for sequence records with metadata and features.
type Record interface {
	Metadata() interface{}
	FeatureTable
	MutableSequence
}

// NewRecord creates a new record.
func NewRecord(meta interface{}, ff []Feature, p []byte) Record {
	seq := Seq(p)
	switch v := meta.(type) {
	case GenBankFields:
		return GenBank{v, FeatureList(ff), NewSequenceServer(seq)}
	default:
		err := fmt.Errorf("gts does not know how to create a record using metadata of type `%T`", v)
		panic(err)
	}
}

// DefaultFormatter returns the default formatter for the given record.
func DefaultFormatter(rec Record) Formatter {
	switch rec.Metadata().(type) {
	case GenBankFields:
		return GenBankFormatter{rec}
	default:
		return GenBankFormatter{rec}
	}
}

// RecordParser attempts to parse a single sequence record.
var RecordParser = pars.Any(GenBankParser)

// RecordScanner scans one sequence record at a time.
type RecordScanner struct {
	s Scanner
	r Record
}

func newRecordScanner(r io.Reader) *MultiParserScanner {
	return NewMultiParserScanner(r,
		GenBankParser,
		DecoderParser(gbDecCtor(NewMsgpackDecoder)),
	)
}

// NewRecordScanner creates a new RecordScanner.
func NewRecordScanner(r io.Reader) *RecordScanner {
	return &RecordScanner{newRecordScanner(r), nil}
}

// NewRecordFileScanner creates a new specialized RecordScanner based on the
// given filename.
func NewRecordFileScanner(f *os.File) *RecordScanner {
	var s Scanner
	switch GetFileType(f.Name()) {
	case GenBankFlat:
		s = GenBankFlatScanner(f)
	case GenBankPack:
		s = GenBankPackScanner(f)
	default:
		s = newRecordScanner(f)
	}
	return &RecordScanner{s, nil}
}

// Scan advances the Scanner to the next Record.
func (s *RecordScanner) Scan() bool {
	ok := s.s.Scan()
	if !ok {
		return false
	}
	s.r, ok = s.s.Value().(Record)
	return ok
}

// Record returns the most recent Record generated by a call to Scan.
func (s *RecordScanner) Record() Record {
	return s.r
}

// Err returns the first non-EOF error that was encountered by the Scanner.
func (s RecordScanner) Err() error {
	return s.s.Err()
}
