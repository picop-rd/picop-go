package header

import (
	"bytes"
	"io"
)

type Header struct{}

func (h *Header) String() string {
	return "BOPP TEST"
}

func (h *Header) Format() []byte {
	return []byte(h.String())
}

// WriteTo renders a proxy protocol header in a format and writes it to an io.Writer.
func (h *Header) WriteTo(w io.Writer) (int64, error) {
	buf := h.Format()
	return bytes.NewBuffer(buf).WriteTo(w)
}
