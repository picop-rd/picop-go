package header

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"

	"go.opentelemetry.io/otel/baggage"
)

var (
	SignatureV1 = []byte{'\x0D', '\x0A', '\x0D', '\x0A', '\x00', '\x0D', '\x0A', '\x51', '\x55', '\x49', '\x54', '\x0A', '\x01'}

	ErrNoBCoP              = errors.New("BCoP: signature not present")
	ErrCannotReadV1Header  = errors.New("BCoP: cannot read v1 header")
	ErrCannotParseBaggage  = errors.New("BCoP: cannot parse baggage")
	ErrLineMustEndWithCrlf = errors.New("BCoP: header must end with \\r\\n")
)

type Header struct {
	version byte
	baggage baggage.Baggage
}

func NewV1(bag baggage.Baggage) *Header {
	return &Header{
		version: 1,
		baggage: bag,
	}
}

func (h *Header) String() string {
	return fmt.Sprintf("BCoP Header{ version: %d, baggage: %s }", h.version, h.baggage.String())
}

func (h *Header) Format() []byte {
	ret := append(SignatureV1, []byte(h.baggage.String())...)
	return append(ret, '\r', '\n')
}

// WriteTo renders a proxy protocol header in a format and writes it to an io.Writer.
func (h *Header) WriteTo(w io.Writer) (int64, error) {
	buf := h.Format()
	return bytes.NewBuffer(buf).WriteTo(w)
}

func Parse(r *bufio.Reader) (*Header, error) {
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

func parseV1(r *bufio.Reader) (*Header, error) {
	buf := []byte{}
	for {
		b, err := r.ReadByte()
		if err != nil {
			return nil, fmt.Errorf(ErrCannotReadV1Header.Error()+": %v", err)
		}
		buf = append(buf, b)
		if b == '\n' {
			break
		}
	}

	if len(buf) < 2 || buf[len(buf)-2] != '\r' {
		return nil, ErrLineMustEndWithCrlf
	}
	header := &Header{
		version: 1,
	}
	bagStr := string(buf[:len(buf)-2])
	bag, err := baggage.Parse(bagStr)
	if err != nil {
		return nil, fmt.Errorf(ErrCannotParseBaggage.Error()+": %v", err)
	}
	header.baggage = bag
	return header, nil
}
