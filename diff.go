package goj

import (
	"bytes"
	"io/ioutil"
	// "log"
	"os"
	"os/exec"
)

func Diff(d1, d2 *Decoder) (b []byte, err error) {
	err = tempFile(d1.StringColorless(), func(f1 *os.File) {
		err = tempFile(d2.StringColorless(), func(f2 *os.File) {
			// if tempFile returns an error, the callback won't be called.
			b, err = exec.Command("git", "diff", "--color=always", "--no-index", f1.Name(), f2.Name()).Output()
			// TODO err?
			b = cleanDiff(b, d1, d2, f1, f2)
		})
	})

	return
}

func cleanDiff(b []byte, d1, d2 *Decoder, f1, f2 *os.File) []byte {
	buf := bytes.NewBuffer(b)

	// skip first two lines
	// git diff specific
	buf.ReadString('\n')
	buf.ReadString('\n')

	// --- filename
	l1, _ := buf.ReadBytes('\n')
	// +++ filename
	l2, _ := buf.ReadBytes('\n')

	b = buf.Bytes()
	buf.Truncate(0)

	l1 = bytes.Replace(l1, []byte(f1.Name()), []byte(d1.FileName()), 1)
	buf.Write(l1)

	l2 = bytes.Replace(l2, []byte(f2.Name()), []byte(d2.FileName()), 1)
	buf.Write(l2)

	buf.Write(b)

	return buf.Bytes()
}

func tempFile(s string, cb func(*os.File)) (err error) {
	f, err := ioutil.TempFile("", "goj")

	if err != nil {
		return
	}

	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()

	f.WriteString(s)

	// rewind
	f.Seek(0, 0)
	// log.Println(f.Name()[1:])
	cb(f)

	return
}
