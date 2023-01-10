package header

import (
	"errors"
	"net/textproto"
	"strings"
)

var headerNewlineToSpace = strings.NewReplacer("\n", " ", "\r", " ")

type MIMEHeader struct {
	textproto.MIMEHeader
}

func NewMIMEHeader() MIMEHeader {
	return MIMEHeader{textproto.MIMEHeader{}}
}

func (h MIMEHeader) String() string {
	// ref: net/http.Header.writeSubset
	ret := make([]string, 0, len(h.MIMEHeader)*4)
	for k, vs := range h.MIMEHeader {
		for _, v := range vs {
			v = headerNewlineToSpace.Replace(v)
			v = textproto.TrimString(v)
			ret = append(ret, k, ":", v, "\r\n")
		}
	}
	if len(ret) <= 1 {
		return ""
	}
	return strings.Join(ret[:len(ret)-1], "") // 最後の\r\nを削除
}

func parseMIMEHeader(data string) (MIMEHeader, error) {
	kvs := strings.Split(data, "\r\n")
	h := NewMIMEHeader()
	for _, kv := range kvs {
		p := strings.Split(kv, ":")
		if len(p) != 2 {
			return MIMEHeader{}, errors.New("invalid header")
		}
		k, v := p[0], p[1]
		h.Add(k, v)
	}
	return h, nil
}
