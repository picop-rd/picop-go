package header

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

var (
	SignatureV1 = []byte{'\x0D', '\x0A', '\x0D', '\x0A', '\x00', '\x0D', '\x0A', '\x51', '\x55', '\x49', '\x54', '\x0A', '\x01'}

	ErrNoBCoP             = errors.New("BCoP: signature not present")
	ErrCannotReadV1Header = errors.New("BCoP: cannot read v1 header")
)

type Header struct {
	version byte
	value   MIMEHeader
}

// NewV1 return BCoP V1 Header.
func NewV1(value MIMEHeader) *Header {
	return &Header{
		version: 1,
		value:   value,
	}
}

func (h Header) Get() MIMEHeader {
	return h.value
}

func (h *Header) Set(value MIMEHeader) {
	h.value = value
}

func (h Header) String() string {
	return fmt.Sprintf("BCoP Header{ version: %d, value: %s }", h.version, h.value.String())
}

func (h Header) Format() []byte {
	b := []byte(h.value.String())
	ret := append(SignatureV1, byte(len(b)))
	return append(ret, b...)
}

// WriteTo renders a proxy protocol header in a format and writes it to an io.Writer.
func (h Header) WriteTo(w io.Writer) (int64, error) {
	buf := h.Format()
	return bytes.NewBuffer(buf).WriteTo(w)
}

func Parse(r io.Reader) (*Header, error) {
	sign := make([]byte, 13)
	_, err := r.Read(sign)
	if err != nil {
		if err == io.EOF {
			return nil, ErrNoBCoP
		}
		return nil, err
	}

	if bytes.Equal(sign, SignatureV1) {
		return parseV1(r)
	}
	return nil, ErrNoBCoP
}

func parseV1(r io.Reader) (*Header, error) {
	length := make([]byte, 1)
	_, err := r.Read(length)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", ErrCannotReadV1Header.Error(), err)
	}

	data := make([]byte, int(length[0]))
	_, err = r.Read(data)
	if err != nil {
		return nil, fmt.Errorf("%s: %v", ErrCannotReadV1Header.Error(), err)
	}

	mimeHeader, err := parseMIMEHeader(string(data))
	if err != nil {
		return nil, fmt.Errorf("%s: %v", ErrCannotReadV1Header.Error(), err)
	}

	return &Header{
		version: 1,
		value:   mimeHeader,
	}, nil
}
