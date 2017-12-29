package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"golang.org/x/net/html/charset"
)

var (
	dryRun      = false
	recurse     = false
	excludeDirs = ".git,.hg,.svn"
	excludes    []string
	force       = false
	verbose     = false
	dir         = "."
)

func init() {
	flag.BoolVar(&dryRun, "dry", dryRun, "perform dry run only")
	flag.BoolVar(&recurse, "recurse", recurse, "find files recursively")
	flag.StringVar(&excludeDirs, "exclude", excludeDirs, "comma separated directories to exclude (in recurse only)")
	flag.BoolVar(&force, "force", force, "convert all files regardless of content-type")
	flag.BoolVar(&verbose, "verbose", verbose, "verbose progress logging")
	flag.StringVar(&dir, "dir", dir, "directory to convert")
}

func main() {
	flag.Parse()

	buf := make([]byte, 1024)
	for _, name := range fileNames() {
		newName := name + ".to8"

		in, err := os.Open(name)
		if err != nil {
			log.Printf("error: %v", err)
			continue
		}
		_, err = in.Read(buf)
		if err != nil {
			log.Printf("error read %q: %v", name, err)
			if err = in.Close(); err != nil {
				log.Printf("error close: %q: %v", name, err)
			}
			continue
		}
		if _, err = in.Seek(0, 0); err != nil {
			log.Printf("error seek beginning: %q: %v", name, err)
			if err = in.Close(); err != nil {
				log.Printf("error close: %q: %v", name, err)
			}
			continue
		}

		switch ct := http.DetectContentType(buf); {
		case force || strings.HasPrefix(ct, "text/"):
			conv, err := charset.NewReader(in, ct)
			if err != nil {
				log.Printf("error create charset reader for %q: %v", name, err)
				continue
			}
			out, err := os.Create(newName)
			if err != nil {
				log.Printf("error create working file %q: %v", newName, err)
				continue
			}
			if verbose {
				log.Printf("Converting %q (%s)", name, ct)
			}
			if _, err = io.Copy(out, conv); err != nil {
				log.Printf("error convert %q to working file %q: %v", name, newName, err)
				continue
			}
			if err = out.Close(); err != nil {
				log.Printf("error close %q: %v", newName, err)
			}
			if err = in.Close(); err != nil {
				log.Printf("error close %q: %v", name, err)
			}
			if !dryRun {
				if err = os.Rename(newName, name); err != nil {
					log.Printf("error move %q over %q: %v", newName, name, err)
					continue
				}
			}
		default:
			if verbose {
				log.Printf("skipping %q: won't convert Content-Type %q", name, ct)
			}
		}
	}
}

func fileNames() []string {
	var names []string
	for _, name := range flag.Args() {
		switch fi, err := os.Stat(name); {
		case err != nil:
			log.Printf("error stat %q: %v", name, err)
		case fi.Mode()&os.ModeType != 0:
			log.Printf("error convert %q: not a regular file", name)
		default:
			names = append(names, name)
		}
	}

	if len(names) == 0 {
		names = append(names, recurseFileNames(dir)...)
	}

	return names
}

func recurseFileNames(root string) []string {
	var names []string
	if files, err := ioutil.ReadDir(root); err != nil {
		log.Printf("error read %q: %v", root, err)
	} else {
		for _, fi := range files {
			if recurse && fi.IsDir() {
				if !shouldExclude(fi.Name()) {
					names = append(names, recurseFileNames(path.Join(root, fi.Name()))...)
				}
				continue
			}
			if fi.Mode()&os.ModeType != 0 {
				if verbose {
					log.Printf("skipping non-regular file %q", fi.Name())
				}
				continue
			}
			names = append(names, path.Join(root, fi.Name()))
		}
	}
	return names
}

func shouldExclude(name string) bool {
	if len(excludeDirs) == 0 {
		return false
	}
	if len(excludes) == 0 {
		excludes = strings.Split(excludeDirs, ",")
		for ix, ex := range excludes {
			excludes[ix] = strings.TrimSpace(ex)
		}
	}
	for _, ex := range excludes {
		if ex == name {
			return true
		}
	}
	return false
}
