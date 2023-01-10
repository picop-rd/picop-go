package header

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParse(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		input     []byte
		wantValue MIMEHeader
		wantErr   bool
	}{
		{
			name:      "headerのみを正常に受理できる",
			input:     makeV1Header(11, "key1:value1"),
			wantValue: makeMIMEHeader("key1", "value1"),
		},
		{
			name:      "header+追加データを正常に受理できる",
			input:     makeV1Header(11, "key1:value1testtest"),
			wantValue: makeMIMEHeader("key1", "value1"),
		},
		{
			name:      "複数headerを正常に受理できる",
			input:     makeV1Header(39, "key1:value1\r\nkey2:value21\r\nkey2:value22"),
			wantValue: makeMIMEHeader("key1", "value1", "key2", "value21", "key2", "value22"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			r := bufio.NewReader(bytes.NewReader(tt.input))
			got, err := Parse(r)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			wantHeader := NewV1()
			wantHeader.Set(tt.wantValue)
			opts := cmp.AllowUnexported(Header{}, MIMEHeader{})
			if diff := cmp.Diff(wantHeader, got, opts); diff != "" {
				t.Errorf("Parse() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestHeader_Format(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		value MIMEHeader
		want  []byte
	}{
		{
			name:  "正しくヘッダをフォーマットできる",
			value: makeMIMEHeader("key1", "value1"),
			want:  makeV1Header(11, "Key1:value1"),
		},
		{
			name:  "正しく複数ヘッダをフォーマットできる",
			value: makeMIMEHeader("key1", "value1", "key2", "value21", "key2", "value22"),
			want:  makeV1Header(39, "Key1:value1\r\nKey2:value21\r\nKey2:value22"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h := NewV1()
			h.Set(tt.value)
			if diff := cmp.Diff(tt.want, h.Format()); diff != "" {
				t.Errorf("Header.Format() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func makeV1Header(l int, d string) []byte {
	ret := append(SignatureV1, byte(l))
	return append(ret, []byte(d)...)
}

func makeMIMEHeader(kv ...string) MIMEHeader {
	h := NewMIMEHeader()
	for i := 0; i < len(kv); i = i + 2 {
		h.Add(kv[i], kv[i+1])
	}
	return h
}
