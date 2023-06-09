package header

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMIMEHeader_String(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		header MIMEHeader
		want   string
	}{
		{
			name:   "It can format a header.",
			header: makeMIMEHeader("key1", "value1"),
			want:   "Key1:value1",
		},
		{
			name:   "It can format headers.",
			header: makeMIMEHeader("key1", "value1", "key2", "value21", "key2", "value22"),
			want:   "Key1:value1\r\nKey2:value21\r\nKey2:value22",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.header.String()
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("MIMEHeader.String() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func makeMIMEHeader(kv ...string) MIMEHeader {
	h := NewMIMEHeader()
	for i := 0; i < len(kv); i = i + 2 {
		h.Add(kv[i], kv[i+1])
	}
	return h
}
