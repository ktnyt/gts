#compdef gts

function _gts_annotate {
    _arguments \
        "-h[show help]" \
        "--help[show help]" \
        "--version[print the version number]" \
        "-F[output file format (defaults to same as input)]" \
        "--format[output file format (defaults to same as input)]" \
        "--no-cache[do not use or create cache]" \
        "-o[output sequence file (specifying `-` will force standard output)]" \
        "--output[output sequence file (specifying `-` will force standard output)]" \
        "*::files:_files"
}

function _gts_cache_list {
    _arguments \
        "-h[show help]" \
        "--help[show help]" \
        "--version[print the version number]" \
         \
        "*::files:_files"
}

function _gts_cache_path {
    _arguments \
        "-h[show help]" \
        "--help[show help]" \
        "--version[print the version number]" \
         \
        "*::files:_files"
}

function _gts_cache_purge {
    _arguments \
        "-h[show help]" \
        "--help[show help]" \
        "--version[print the version number]" \
         \
        "*::files:_files"
}

function _gts_cache {
    local line

    function _commands {
        local -a commands
        commands=(
            'list:list the cache files'
            'path:print the cache directory path'
            'purge:delete all cache files'
        )
        _describe 'command' commands
    }

    _arguments -C \
        "-h[show help]" \
        "--help[show help]" \
        "--version[print the version number]" \
        "1: :_commands" \
        "*::arg:->args"

    case $line[1] in
        list)  _gts_cache_list ;;
        path)  _gts_cache_path ;;
        purge) _gts_cache_purge ;;
        *) ;;
    esac
}

function _gts_clear {
    _arguments \
        "-h[show help]" \
        "--help[show help]" \
        "--version[print the version number]" \
        "-F[output file format (defaults to same as input)]" \
        "--format[output file format (defaults to same as input)]" \
        "--no-cache[do not use or create cache]" \
        "-o[output sequence file (specifying `-` will force standard output)]" \
        "--output[output sequence file (specifying `-` will force standard output)]" \
        "*::files:_files"
}

function _gts_complement {
    _arguments \
        "-h[show help]" \
        "--help[show help]" \
        "--version[print the version number]" \
        "-F[output file format (defaults to same as input)]" \
        "--format[output file format (defaults to same as input)]" \
        "--no-cache[do not use or create cache]" \
        "-o[output sequence file (specifying `-` will force standard output)]" \
        "--output[output sequence file (specifying `-` will force standard output)]" \
        "*::files:_files"
}

function _gts_define {
    _arguments \
        "-h[show help]" \
        "--help[show help]" \
        "--version[print the version number]" \
        "-F[output file format (defaults to same as input)]" \
        "--format[output file format (defaults to same as input)]" \
        "--no-cache[do not use or create cache]" \
        "-o[output sequence file (specifying `-` will force standard output)]" \
        "--output[output sequence file (specifying `-` will force standard output)]" \
        "-q[qualifier key-value pairs (syntax: key=value))]" \
        "--qualifier[qualifier key-value pairs (syntax: key=value))]" \
        "*::files:_files"
}

function _gts_delete {
    _arguments \
        "-h[show help]" \
        "--help[show help]" \
        "--version[print the version number]" \
        "-e[remove features contained in the deleted regions]" \
        "--erase[remove features contained in the deleted regions]" \
        "-F[output file format (defaults to same as input)]" \
        "--format[output file format (defaults to same as input)]" \
        "--no-cache[do not use or create cache]" \
        "-o[output sequence file (specifying `-` will force standard output)]" \
        "--output[output sequence file (specifying `-` will force standard output)]" \
        "*::files:_files"
}

function _gts_extract {
    _arguments \
        "-h[show help]" \
        "--help[show help]" \
        "--version[print the version number]" \
        "-F[output file format (defaults to same as input)]" \
        "--format[output file format (defaults to same as input)]" \
        "--no-cache[do not use or create cache]" \
        "-o[output sequence file (specifying `-` will force standard output)]" \
        "--output[output sequence file (specifying `-` will force standard output)]" \
        "*::files:_files"
}

function _gts_infix {
    _arguments \
        "-h[show help]" \
        "--help[show help]" \
        "--version[print the version number]" \
        "-e[extend existing feature locations when inserting instead of splitting them]" \
        "--embed[extend existing feature locations when inserting instead of splitting them]" \
        "-F[output file format (defaults to same as input)]" \
        "--format[output file format (defaults to same as input)]" \
        "--no-cache[do not use or create cache]" \
        "-o[output sequence file (specifying `-` will force standard output)]" \
        "--output[output sequence file (specifying `-` will force standard output)]" \
        "*::files:_files"
}

function _gts_insert {
    _arguments \
        "-h[show help]" \
        "--help[show help]" \
        "--version[print the version number]" \
        "-e[extend existing feature locations when inserting instead of splitting them]" \
        "--embed[extend existing feature locations when inserting instead of splitting them]" \
        "-F[output file format (defaults to same as input)]" \
        "--format[output file format (defaults to same as input)]" \
        "--no-cache[do not use or create cache]" \
        "-o[output sequence file (specifying `-` will force standard output)]" \
        "--output[output sequence file (specifying `-` will force standard output)]" \
        "*::files:_files"
}

function _gts_join {
    _arguments \
        "-h[show help]" \
        "--help[show help]" \
        "--version[print the version number]" \
        "-c[output the sequence as circular if possible]" \
        "--circular[output the sequence as circular if possible]" \
        "-F[output file format (defaults to same as input)]" \
        "--format[output file format (defaults to same as input)]" \
        "--no-cache[do not use or create cache]" \
        "-o[output sequence file (specifying `-` will force standard output)]" \
        "--output[output sequence file (specifying `-` will force standard output)]" \
        "*::files:_files"
}

function _gts_length {
    _arguments \
        "-h[show help]" \
        "--help[show help]" \
        "--version[print the version number]" \
        "-o[output file (specifying `-` will force standard output)]" \
        "--output[output file (specifying `-` will force standard output)]" \
        "*::files:_files"
}

function _gts_pick {
    _arguments \
        "-h[show help]" \
        "--help[show help]" \
        "--version[print the version number]" \
        "-f[pick features instead of sequences]" \
        "--feature[pick features instead of sequences]" \
        "-F[output file format (defaults to same as input)]" \
        "--format[output file format (defaults to same as input)]" \
        "--no-cache[do not use or create cache]" \
        "-o[output sequence file (specifying `-` will force standard output)]" \
        "--output[output sequence file (specifying `-` will force standard output)]" \
        "*::files:_files"
}

function _gts_query {
    _arguments \
        "-h[show help]" \
        "--help[show help]" \
        "--version[print the version number]" \
        "-d[string to insert between columns]" \
        "--delimiter[string to insert between columns]" \
        "--empty[allow missing qualifiers to be reported]" \
        "-H[do not print the header line]" \
        "--no-header[do not print the header line]" \
        "-I[do not report the sequence identifier]" \
        "--no-seqid[do not report the sequence identifier]" \
        "-K[do not report the feature key]" \
        "--no-key[do not report the feature key]" \
        "-L[do not report the feature location]" \
        "--no-location[do not report the feature location]" \
        "--no-cache[do not use or create cache]" \
        "-n[qualifier name(s) to select]" \
        "--name[qualifier name(s) to select]" \
        "-o[output table file (specifying `-` will force standard output)]" \
        "--output[output table file (specifying `-` will force standard output)]" \
        "--source[include the source feature(s)]" \
        "-t[string to insert between qualifier values]" \
        "--separator[string to insert between qualifier values]" \
        "*::files:_files"
}

function _gts_repair {
    _arguments \
        "-h[show help]" \
        "--help[show help]" \
        "--version[print the version number]" \
        "-F[output file format (defaults to same as input)]" \
        "--format[output file format (defaults to same as input)]" \
        "--no-cache[do not use or create cache]" \
        "-o[output sequence file (specifying `-` will force standard output)]" \
        "--output[output sequence file (specifying `-` will force standard output)]" \
        "*::files:_files"
}

function _gts_reverse {
    _arguments \
        "-h[show help]" \
        "--help[show help]" \
        "--version[print the version number]" \
        "-F[output file format (defaults to same as input)]" \
        "--format[output file format (defaults to same as input)]" \
        "--no-cache[do not use or create cache]" \
        "-o[output sequence file (specifying `-` will force standard output)]" \
        "--output[output sequence file (specifying `-` will force standard output)]" \
        "*::files:_files"
}

function _gts_rotate {
    _arguments \
        "-h[show help]" \
        "--help[show help]" \
        "--version[print the version number]" \
        "-F[output file format (defaults to same as input)]" \
        "--format[output file format (defaults to same as input)]" \
        "--no-cache[do not use or create cache]" \
        "-o[output sequence file (specifying `-` will force standard output)]" \
        "--output[output sequence file (specifying `-` will force standard output)]" \
        "*::files:_files"
}

function _gts_search {
    _arguments \
        "-h[show help]" \
        "--help[show help]" \
        "--version[print the version number]" \
        "-e[match the exact pattern even for ambiguous letters]" \
        "--exact[match the exact pattern even for ambiguous letters]" \
        "-F[output file format (defaults to same as input)]" \
        "--format[output file format (defaults to same as input)]" \
        "-k[key for the reported oligomer region features]" \
        "--key[key for the reported oligomer region features]" \
        "--no-cache[do not use or create cache]" \
        "--no-complement[do not match the complement strand]" \
        "-o[output sequence file (specifying `-` will force standard output)]" \
        "--output[output sequence file (specifying `-` will force standard output)]" \
        "-q[qualifier key-value pairs (syntax: key=value))]" \
        "--qualifier[qualifier key-value pairs (syntax: key=value))]" \
        "*::files:_files"
}

function _gts_select {
    _arguments \
        "-h[show help]" \
        "--help[show help]" \
        "--version[print the version number]" \
        "-F[output file format (defaults to same as input)]" \
        "--format[output file format (defaults to same as input)]" \
        "--no-cache[do not use or create cache]" \
        "-o[output sequence file (specifying `-` will force standard output)]" \
        "--output[output sequence file (specifying `-` will force standard output)]" \
        "-s[strand to select features from (`both`, `forward`, or `reverse`)]" \
        "--strand[strand to select features from (`both`, `forward`, or `reverse`)]" \
        "-v[select features that do not match the given criteria]" \
        "--invert-match[select features that do not match the given criteria]" \
        "*::files:_files"
}

function _gts_sort {
    _arguments \
        "-h[show help]" \
        "--help[show help]" \
        "--version[print the version number]" \
        "-F[output file format (defaults to same as input)]" \
        "--format[output file format (defaults to same as input)]" \
        "--no-cache[do not use or create cache]" \
        "-o[output sequence file (specifying `-` will force standard output)]" \
        "--output[output sequence file (specifying `-` will force standard output)]" \
        "-r[reverse the sort order]" \
        "--reverse[reverse the sort order]" \
        "*::files:_files"
}

function _gts_split {
    _arguments \
        "-h[show help]" \
        "--help[show help]" \
        "--version[print the version number]" \
        "-F[output file format (defaults to same as input)]" \
        "--format[output file format (defaults to same as input)]" \
        "--no-cache[do not use or create cache]" \
        "-o[output sequence file (specifying `-` will force standard output)]" \
        "--output[output sequence file (specifying `-` will force standard output)]" \
        "*::files:_files"
}

function _gts_summary {
    _arguments \
        "-h[show help]" \
        "--help[show help]" \
        "--version[print the version number]" \
        "-F[suppress feature summary]" \
        "--no-feature[suppress feature summary]" \
        "--no-cache[do not use or create cache]" \
        "-o[output file (specifying `-` will force standard output)]" \
        "--output[output file (specifying `-` will force standard output)]" \
        "-Q[suppress qualifier summary]" \
        "--no-qualifier[suppress qualifier summary]" \
        "*::files:_files"
}

function _gts {
    local line

    function _commands {
        local -a commands
        commands=(
            'annotate:merge features from a feature list file into a sequence'
            'cache:manage gts cache files'
            'clear:remove all features from the sequence (excluding source features)'
            'complement:compute the complement of the given sequence'
            'define:define a new feature'
            'delete:delete a region of the given sequence(s)'
            'extract:extract the sequences referenced by the features'
            'infix:infix input sequence(s) into the host sequence(s)'
            'insert:insert guest sequence(s) into the input sequence(s)'
            'join:join the sequences contained in the files'
            'length:report the length of the sequence(s)'
            'pick:pick sequence(s) from multiple sequences'
            'query:query information from the given sequence'
            'repair:repair fragmented features'
            'reverse:reverse order of the given sequence(s)'
            'rotate:shift the coordinates of a circular sequence'
            'search:search for a subsequence and annotate its results'
            'select:select features using the given feature selector(s)'
            'sort:sort the list of sequences'
            'split:split the sequence at the provided locations'
            'summary:report a brief summary of the sequence(s)'
        )
        _describe 'command' commands
    }

    _arguments -C \
        "-h[show help]" \
        "--help[show help]" \
        "--version[print the version number]" \
        "1: :_commands" \
        "*::arg:->args"

    case $line[1] in
        annotate)   _gts_annotate ;;
        cache)      _gts_cache ;;
        clear)      _gts_clear ;;
        complement) _gts_complement ;;
        define)     _gts_define ;;
        delete)     _gts_delete ;;
        extract)    _gts_extract ;;
        infix)      _gts_infix ;;
        insert)     _gts_insert ;;
        join)       _gts_join ;;
        length)     _gts_length ;;
        pick)       _gts_pick ;;
        query)      _gts_query ;;
        repair)     _gts_repair ;;
        reverse)    _gts_reverse ;;
        rotate)     _gts_rotate ;;
        search)     _gts_search ;;
        select)     _gts_select ;;
        sort)       _gts_sort ;;
        split)      _gts_split ;;
        summary)    _gts_summary ;;
        *) ;;
    esac
}

