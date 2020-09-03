package jsonc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
)

// These shortcut functions strip comments and allow JSON parsing to proceed as expected afterwards.

// See json.Unmarsal.
func Unmarshal(data []byte, v interface{}) error {
	buf, err := stripComments(data)
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
		stripped, err := stripComments(in)
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

func stripComments(data []byte) ([]byte, error) {
	l, items := lex("jsonc-strip", string(data))
	defer l.drain()
	buf := make([]byte, 0, len(data))
	for i := range items {
		if i.typ == itemError {
			return nil, fmt.Errorf(i.val)
		} else if i.typ != itemComment && i.typ != itemEOF {
			buf = append(buf, i.val...)
		}

	}
	return buf, nil
}
