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
		wantValue string
		wantErr   bool
	}{
		{
			name:      "headerのみを正常に受理できる",
			input:     makeV1Header(11, "key1=value1"),
			wantValue: "key1=value1",
		},
		{
			name:      "header+追加データを正常に受理できる",
			input:     makeV1Header(11, "key1=value1testtest"),
			wantValue: "key1=value1",
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
			opts := cmp.AllowUnexported(Header{})
			if diff := cmp.Diff(NewV1(tt.wantValue), got, opts); diff != "" {
				t.Errorf("Parse() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestHeader_Format(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		value string
		want  []byte
	}{
		{
			name:  "正しくヘッダをフォーマットできる",
			value: "key1=value1",
			want:  makeV1Header(11, "key1=value1"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			h := NewV1(tt.value)
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
