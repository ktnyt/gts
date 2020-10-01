# gts-complement(1) -- compute the complement of the given sequence(s)

## SYNOPSIS

gts-complement [--version] [-h | --help] [<args>] <input>

## DESCRIPTION

**gts-complement** takes a single sequence input and return the complemented
sequence as output. Any features present in the sequence will be relocated to
the complement strand. This command _will not_ reverse the sequence. To obtain
the reversed sequence, use **gts-reverse(1)**.

## OPTIONS

  * `<input>`:
    Input sequence file (may be omitted if standard input is provided). See
    gts-seqin(7) for a list of currently supported list of sequence formats.

  * `-F <format>`, `--format=<format>`:
    Output file format (defaults to same as input). See gts-seqout(7) for a
    list of currently supported list of sequence formats. The format specified
    with this option will override the file type detection from the output
    filename.

  * `-o <output>`, `--output=<output>`:
    Output sequence file (specifying `-` will force standard output). The
    output file format will be automatically detected from the filename if none
    is specified with the `-F` or `--format` option.

## BUGS

**gts-complement** currently has no known bugs.

## AUTHORS

**gts-complement** is written and maintained by Kotone Itaya.

## SEE ALSO

gts(1), gts-reverse(1), gts-seqin(7), gts-seqout(7)