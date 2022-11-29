package header

import (
	"bufio"
	"bytes"
	"errors"
	"io"
)

var (
	TestHeaderString = "BCoP TEST\n"
	ErrNoBCoP        = errors.New("BCoP: signature not present")
)

type Header struct{}

func (h *Header) String() string {
	return TestHeaderString
}

func (h *Header) Format() []byte {
	return []byte(h.String())
}

// WriteTo renders a proxy protocol header in a format and writes it to an io.Writer.
func (h *Header) WriteTo(w io.Writer) (int64, error) {
	buf := h.Format()
	return bytes.NewBuffer(buf).WriteTo(w)
}

func Read(reader *bufio.Reader) (*Header, error) {
	sign, err := reader.Peek(len([]byte(TestHeaderString)))
	if err != nil {
		if err == io.EOF {
			return nil, ErrNoBCoP
		}
		return nil, err
	}

	if bytes.Equal(sign, []byte(TestHeaderString)) {
		return &Header{}, nil
	}
	return nil, ErrNoBCoP
}
