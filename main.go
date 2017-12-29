package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/html/charset"
)

func main() {
	files, err := ioutil.ReadDir(".")
	check(err)
	buf := make([]byte, 512)
	for _, file := range files {
		if file.Mode()&os.ModeType != 0 {
			// skip directories, and other special file types
			continue
		}
		name := file.Name()
		newName := name + ".to8"
		in, err := os.Open(name)
		check(err)
		_, err = in.Read(buf)
		check(err)
		in.Seek(0, 0)
		contentType := http.DetectContentType(buf)
		conv, err := charset.NewReader(in, contentType)
		out, err := os.Create(newName)
		check(err)
		fmt.Printf("Converting %q (%s) ...", name, contentType)
		_, err = io.Copy(out, conv)
		check(err)
		check(out.Close())
		check(in.Close())
		check(os.Rename(newName, name))
		fmt.Printf("  done\n")
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
