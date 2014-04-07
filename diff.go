package goj

import (
	"io/ioutil"
	"os/exec"
)

func diff(d1, d2 *Decoder) ([]byte, error) {
	f1, err := ioutil.TempFile("", "goj")
	f2, err := ioutil.TempFile("", "goj")

	f1Name := f1.Name()
	f2Name := f2.Name()

	defer func() {
		f1.Close()
		f2.Close()
		os.Remove(f1Name)
		os.Remove(f2Name)
	}()

	f1.WriteString(d1.String())
	f2.WriteString(d2.String())

	// rewind
	f1.Seek(0, 0)
	f2.Seek(0, 0)

	// diff command configurable
	// diff args configurable
	return exec.Command("diff", f1Name, f2Name).Output()
}

func diff1(d1, d2 *Decoder) ([]byte, error) {
	tempFile(d1.String(), func(f1 *os.File) {
		tempFile(d2.String(), func(f2 *os.File) {
			return exec.Command("diff", f1.Name(), f2.Name()).Output()
		})
	})
}

func tempFile(s string, cb func(*os.File)) {
	// TODO do something with the error.
	f, _ := ioutil.TempFile("", "goj")

	defer func() {
		f.Close()
		os.Remove(f.Name())
	}()

	f.WriteString(s)

	// rewind
	f.Seek(0, 0)
	cb(f)
}
