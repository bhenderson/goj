package goj

import (
	"bytes"
	"io/ioutil"
	// "log"
	"os"
	"os/exec"
)

func Diff(d1, d2 *Decoder) (b []byte, err error) {
	// reset color
	c1, c2 := d1.color, d2.color
	d1.color, d2.color = ColorNever, ColorNever
	defer func() {
		d1.color, d2.color = c1, c2
	}()

	err = tempFile(d1.String(), func(f1 *os.File) {
		err = tempFile(d2.String(), func(f2 *os.File) {
			// if tempFile returns an error, the callback won't be called.
			b, err = exec.Command("git", "diff", "--color=always", "--no-index", f1.Name(), f2.Name()).Output()
		})
	})

	b = cleanDiff(b, d1, d2)

	return
}

func cleanDiff(b []byte, d1, d2 *Decoder) []byte {
	buf := bytes.NewBuffer(b)

	// skip first four lines
	// git diff specific
	buf.ReadString('\n')
	buf.ReadString('\n')
	buf.ReadString('\n')
	buf.ReadString('\n')

	b = buf.Bytes()
	buf.Truncate(0)

	buf.Write(yellowb)
	buf.WriteString("--- ")
	buf.WriteString(d1.FileName())
	buf.Write(reset)
	buf.WriteRune('\n')

	buf.Write(yellowb)
	buf.WriteString("+++ ")
	buf.WriteString(d2.FileName())
	buf.WriteRune('\n')
	buf.Write(reset)

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
