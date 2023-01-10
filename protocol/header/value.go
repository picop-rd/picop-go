package header

import (
	"errors"
	"net/textproto"
	"sort"
	"strings"
	"sync"
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
	kvs, sorter, l := h.sortedKeyValues()
	defer headerSorterPool.Put(sorter)
	ret := make([]string, 0, l*4)
	for _, kv := range kvs {
		for _, v := range kv.values {
			v = headerNewlineToSpace.Replace(v)
			v = textproto.TrimString(v)
			ret = append(ret, kv.key, ":", v, "\r\n")
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

type keyValues struct {
	key    string
	values []string
}

// A headerSorter implements sort.Interface by sorting a []keyValues
// by key. It's used as a pointer, so it can fit in a sort.Interface
// interface value without allocation.
type headerSorter struct {
	kvs []keyValues
}

func (s *headerSorter) Len() int           { return len(s.kvs) }
func (s *headerSorter) Swap(i, j int)      { s.kvs[i], s.kvs[j] = s.kvs[j], s.kvs[i] }
func (s *headerSorter) Less(i, j int) bool { return s.kvs[i].key < s.kvs[j].key }

var headerSorterPool = sync.Pool{
	New: func() any { return new(headerSorter) },
}

func (h MIMEHeader) sortedKeyValues() (kvs []keyValues, hs *headerSorter, length int) {
	hs = headerSorterPool.Get().(*headerSorter)
	if cap(hs.kvs) < len(h.MIMEHeader) {
		hs.kvs = make([]keyValues, 0, len(h.MIMEHeader))
	}
	kvs = hs.kvs[:0]
	for k, vv := range h.MIMEHeader {
		kvs = append(kvs, keyValues{k, vv})
		length += len(vv)
	}
	hs.kvs = kvs
	sort.Sort(hs)
	return kvs, hs, length
}
