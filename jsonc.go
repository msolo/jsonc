// jsonc provides an analogous API to the standard json package, but
// allows json data to contain comments // ...  or /* ... */. These
// functions strip comments and allow JSON parsing to proceed as
// expected using the standard json package.
package jsonc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

// See json.Unmarsal.
func Unmarshal(data []byte, v interface{}) error {
	buf, err := StripComments(data)
	if err != nil {
		return err
	}
	return json.Unmarshal(buf, v)
}

type jsoncReader struct {
	r   io.Reader
	buf *bytes.Buffer
}

func (jr *jsoncReader) Read(b []byte) (n int, err error) {
	if jr.buf == nil {
		in, err := ioutil.ReadAll(jr.r)
		if err != nil {
			return 0, err
		}
		stripped, err := StripComments(in)
		if err != nil {
			return 0, err
		}
		jr.buf = bytes.NewBuffer(stripped)
	}
	return jr.buf.Read(b)
}

// FIXME(msolo) This strips a whole buffer at a time rather than reading incrementally from the underlying reader. No one should confuse JSONC for something high performance, but we needed waste too many resources.
//
// See json.NewDecoder.
func NewDecoder(r io.Reader) *json.Decoder {
	jr := &jsoncReader{r: r}
	return json.NewDecoder(jr)
}

// Return a copy of the input with all comments stripped.
func StripComments(data []byte) ([]byte, error) {
	l := lex("jsonc-strip", string(data))
	buf := make([]byte, 0, len(data))
	for {
		i := l.yield()
		if i.typ == itemEOF {
			break
		} else if i.typ == itemError {
			return nil, fmt.Errorf(i.val)
		} else if i.typ != itemComment && i.typ != itemEOF {
			buf = append(buf, i.val...)
		}
	}
	return buf, nil
}
