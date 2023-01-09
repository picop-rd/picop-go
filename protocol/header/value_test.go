package header

import (
	"testing"
)

func TestMIMEHeader_String(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		header MIMEHeader
		want   string
	}{
		{
			name:   "正しくフォーマットできる",
			header: makeMIMEHeader("key1", "value1"),
			want:   "Key1:value1",
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.header.String(); got != tt.want {
				t.Errorf("MIMEHeader.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
