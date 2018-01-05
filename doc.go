// to8 provides a command line tool to convert files from most other
// encodings to UTF8 with no BOM. This is useful because many Go stdlib
// utilities assume UTF8, with no byte-order mark (BOM)
//
//      Usage of to8:
//        -dir string
//            directory to convert (default ".")
//        -dry
//            perform dry run only
//        -exclude string
//            comma separated directories to exclude (in recurse only) (default ".git,.hg,.svn")
//        -force
//            convert all files regardless of content-type
//        -recurse
//            find files recursively
//        -verbose
//            verbose progress logging
package main
