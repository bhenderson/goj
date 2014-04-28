package goj

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDiff(t *testing.T) {
	v1 := testVal(t, `{"a":"b"}`)
	v2 := testVal(t, `{"a":"c"}`)
	v2.file.(*testFile).n = "/some/file"

	exp := `[1;33m--- a/test input[m
[1;33m+++ b/some/file[m
[1;35m@@ -1,3 +1,3 @@[m
 {[m
[1;31m-  "a": "b"[m
[1;32m+[m[1;32m  "a": "c"[m
 }[m
\ No newline at end of file[m
`

	b, err := Diff(v1, v2)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, exp, string(b))
}

func TestDiff_single(t *testing.T) {
	d := testDecoder(t, `{"a":"b"}`).Decode("")
	v1 := <-d

	exp := `[1;33m--- a/test input[m
[1;33m+++ b/test input[m
[1;35m@@ -1,3 +1 @@[m
[1;31m-{[m
[1;31m-  "a": "b"[m
[1;31m-}[m
\ No newline at end of file[m
[1;32m+[m[1;32mnull[m
\ No newline at end of file[m
`

	b, err := Diff(v1, nil)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, exp, string(b))

	exp = `[1;33m--- a/test input[m
[1;33m+++ b/test input[m
[1;35m@@ -1 +1,3 @@[m
[1;31m-null[m
\ No newline at end of file[m
[1;32m+[m[1;32m{[m
[1;32m+[m[1;32m  "a": "b"[m
[1;32m+[m[1;32m}[m
\ No newline at end of file[m
`

	b, err = Diff(nil, v1)

	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, exp, string(b))
}

func BenchmarkDiff_same(b *testing.B) {
	d := testDecoder(b, input+input).Decode("")
	v1, v2 := <-d, <-d
	if v1 == nil || v2 == nil {
		b.Fatal("v1 or v2 are nil")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Diff(v1, v2)
	}
}

func BenchmarkDiff(b *testing.B) {
	// change the last price
	last := input[:len(input)-10] + "0}}}"
	d := testDecoder(b, input+last).Decode("")
	v1, v2 := <-d, <-d
	if v1 == nil || v2 == nil {
		b.Fatal("v1 or v2 are nil")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Diff(v1, v2)
	}
}
