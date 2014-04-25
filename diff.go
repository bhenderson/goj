package goj

import (
	"bytes"
	"io/ioutil"
	// "log"
	"os"
	"os/exec"
)

func Diff(v1, v2 *Val) (b []byte, err error) {
	// TODO error handling
	b1, _ := v1.MarshalJSON()
	err = tempFile(b1, func(f1 *os.File) {
		b2, _ := v2.MarshalJSON()
		err = tempFile(b2, func(f2 *os.File) {
			// if tempFile returns an error, the callback won't be called.
			b, err = exec.Command("git", "diff", "--color=always", "--no-index", f1.Name(), f2.Name()).Output()
			// TODO err?
			b = cleanDiff(b, v1, v2, f1, f2)
		})
	})

	return
}

func cleanDiff(b []byte, v1, v2 *Val, f1, f2 File) []byte {
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

	l1 = cleanName(l1, f1, v1)
	buf.Write(l1)

	l2 = cleanName(l2, f2, v2)
	buf.Write(l2)

	buf.Write(b)

	return buf.Bytes()
}

func cleanName(l []byte, f File, v *Val) []byte {
	name1 := f.Name()
	name2 := v.FileName()
	if name2[0] != '/' {
		name1 = name1[1:]
	}
	return bytes.Replace(l, []byte(name1), []byte(name2), 1)
}

func tempFile(b []byte, cb func(*os.File)) (err error) {
	f, err := ioutil.TempFile("", "goj")

	if err != nil {
		return
	}

	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()

	f.Write(b)

	// rewind
	f.Seek(0, 0)
	// log.Println(f.Name()[1:])
	cb(f)

	return
}
