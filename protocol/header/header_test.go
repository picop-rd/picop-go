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
		wantValue *Header
		wantErr   bool
	}{
		{
			name:      "It can accept only a header.",
			input:     makeV1HeaderByte(11, "key1:value1"),
			wantValue: makeV1HeaderStruct("key1", "value1"),
		},
		{
			name:      "It can accept a header and additional data.",
			input:     makeV1HeaderByte(11, "key1:value1testtest"),
			wantValue: makeV1HeaderStruct("key1", "value1"),
		},
		{
			name:      "It can accept headers.",
			input:     makeV1HeaderByte(39, "key1:value1\r\nkey2:value21\r\nkey2:value22"),
			wantValue: makeV1HeaderStruct("key1", "value1", "key2", "value21", "key2", "value22"),
		},
		{
			name:      "It can accept an empty value",
			input:     makeV1HeaderByte(5, "key1:"),
			wantValue: makeV1HeaderStruct("key1", ""),
		},
		{
			name:      "It can raise an error if the value has :",
			input:     makeV1HeaderByte(11, "key1:va:ue1"),
			wantValue: nil,
			wantErr:   true,
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
			opts := cmp.AllowUnexported(Header{}, MIMEHeader{})
			if diff := cmp.Diff(tt.wantValue, got, opts); diff != "" {
				t.Errorf("Parse() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestHeader_Format(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		value *Header
		want  []byte
	}{
		{
			name:  "It can format a header.",
			value: makeV1HeaderStruct("key1", "value1"),
			want:  makeV1HeaderByte(11, "Key1:value1"),
		},
		{
			name:  "It can format headers.",
			value: makeV1HeaderStruct("key1", "value1", "key2", "value21", "key2", "value22"),
			want:  makeV1HeaderByte(39, "Key1:value1\r\nKey2:value21\r\nKey2:value22"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if diff := cmp.Diff(tt.want, tt.value.Format()); diff != "" {
				t.Errorf("Header.Format() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func makeV1HeaderByte(l int, d string) []byte {
	ret := append(SignatureV1, byte(l))
	return append(ret, []byte(d)...)
}

func makeV1HeaderStruct(kv ...string) *Header {
	h := NewV1()
	for i := 0; i < len(kv); i = i + 2 {
		h.Add(kv[i], kv[i+1])
	}
	return h
}
