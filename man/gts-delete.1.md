# gts-delete(1) -- delete a region of the given sequence(s)

## SYNOPSIS

gts-delete [--version] [-h | --help] [<args>] <locator> <input>

## DESCRIPTION

**gts-delete** takes a single sequence input and deletes the specified region.
The region to be deleted is specified using a `locator`. A locator is a
combination of one of `point location`, `range location`, or `selector`, and a
`modifier` in the form `[selector|point|range][@modifier]`. See gts-locator(7)
for a more in-depth explanation of a locator. Refer to the EXAMPLES for some
examples to get started.

Features that were present in the region being deleted will be shifted as being
in between the bases at the deletion point. Such features can be completely
erased from the sequence if the `-e` or `--erase` option is provided.

## OPTIONS

  * `<locator>`:
    A locator string ([selector|point|range][@modifier]). See gts-locator(7)
    for more details.

  * `<input>`:
    Input sequence file (may be omitted if standard input is provided). See
    gts-seqin(7) for a list of currently supported list of sequence formats.

  * `-F <format>`, `--format=<format>`:
    Output file format (defaults to same as input). See gts-seqout(7) for a
    list of currently supported list of sequence formats. The format specified
    with this option will override the file type detection from the output
    filename.

  * `-e`, `--erase`:
    Remove features contained in the deleted regions.

  * `-o <output>`, `--output=<output>`:
    Output sequence file (specifying `-` will force standard output). The
    output file format will be automatically detected from the filename if none
    is specified with the `-F` or `--format` option.

## EXAMPLES

Delete bases 100 to 200:

    $ gts delete 100..200 <input>

Delete all regions of `misc_feature` and its features:

    $ gts delete --erase misc_feature <input>

Delete 20 bases upstream of every `CDS`:

    $ gts delete CDS^-20..^ <input>

## BUGS

**gts-delete** currently has no known bugs.

## AUTHORS

**gts-delete** is written and maintained by Kotone Itaya.

## SEE ALSO

gts(1), gts-insert(1), gts-locator(7), gts-selector(7), gts-seqin(7),
gts-seqout(7)