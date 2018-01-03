package main

import (
	"errors"
	"io"
)

func NewSkipBOMWriter(w io.Writer) io.Writer { return &skipBomWriter{w: w} }

type skipBomWriter struct {
	w       io.Writer
	skipped bool
}

func (w *skipBomWriter) Write(p []byte) (n int, err error) {
	if !w.skipped {
		switch len(p) {
		case 0, 1, 2, 3:
			return -1, errors.New("need at least 4 bytes on first write")
		default:
			w.skipped = true
			switch {
			case isUTF32(p[:4]) || isUTF32LE(p[:4]):
				x, err := w.w.Write(p[4:])
				return x + 4, err
			case isUTF16(p[:2]) || isUTF16LE(p[:2]):
				x, err := w.w.Write(p[2:])
				return x + 2, err
			case isUTF8(p[:3]):
				x, err := w.w.Write(p[3:])
				return x + 3, err
			default:
				return w.w.Write(p)
			}
		}

	}
	return w.w.Write(p)
}

func isUTF32(b []byte) bool   { return b[0] == 0x00 && b[1] == 0x00 && b[2] == 0xFE && b[3] == 0xFF }
func isUTF32LE(b []byte) bool { return b[0] == 0x00 && b[1] == 0x00 && b[2] == 0xFF && b[3] == 0xFE }
func isUTF16(b []byte) bool   { return b[0] == 0xFE && b[1] == 0xFF }
func isUTF16LE(b []byte) bool { return b[0] == 0xFF && b[1] == 0xFE }
func isUTF8(b []byte) bool    { return b[0] == 0xEF && b[1] == 0xBB && b[2] == 0xBF }

func skipBOM(r io.ReadSeeker) {
	var bom [4]byte
	switch n, err := r.Read(bom[:]); {
	case err != nil:
		panic(err)
	case n < 4:
		// skip processing
	case isUTF32(bom[:4]) || isUTF32LE(bom[:4]):
		if _, err := r.Seek(4, 0); err != nil {
			panic(err)
		}
		return
	case isUTF16(bom[:2]) || isUTF16LE(bom[:2]):
		if _, err := r.Seek(2, 0); err != nil {
			panic(err)
		}
		return
	case isUTF8(bom[:3]):
		if _, err := r.Seek(3, 0); err != nil {
			panic(err)
		}
		return
	default:
	}
	// anything other than successful read of 4 bytes; bail.
	if _, err := r.Seek(0, 0); err != nil {
		panic(err)
	}
}
