package header

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name      string
		input     []byte
		wantValue string
		wantErr   bool
	}{
		{
			name:      "headerのみを正常に受理できる",
			input:     append(SignatureV1, []byte("key1=value1\r\n")...),
			wantValue: "key1=value1",
		},
		{
			name:      "header+追加データを正常に受理できる",
			input:     append(SignatureV1, []byte("key1=value1\r\ntesttest")...),
			wantValue: "key1=value1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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
