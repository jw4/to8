# to8
convert files to utf-8

## Install

`go get jw4.us/to8`

## Usage

By default this utility will process all named file arguments.

If no names are given, it will find all the regular files in the current directory.

If the Content-Type is determined to be `text/*` (or if `-force` is specified), a new file will be created with a suffix of `.to8` and the utility will copy the contents of the file to the new file after converting the charset to UTF-8.  After the copy is done (and if `-dry` is not specified), the new file will be copied over the original file, preserving the file mode bits.


```
  Usage of to8:
    -dir string
        directory to convert (default ".")
    -dry
        perform dry run only
    -exclude string
        comma separated directories to exclude (in recurse only) (default ".git,.hg,.svn")
    -force
        convert all files regardless of content-type
    -recurse
        find files recursively
    -verbose
        verbose progress logging
```
