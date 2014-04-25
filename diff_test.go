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
