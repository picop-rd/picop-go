package header

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/google/go-cmp/cmp"
	"go.opentelemetry.io/otel/baggage"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name            string
		input           []byte
		wantBaggageFunc func() baggage.Baggage
		wantErr         bool
	}{
		{
			name:  "headerのみを正常に受理できる",
			input: append(SignatureV1, []byte("key1=value1\r\n")...),
			wantBaggageFunc: func() baggage.Baggage {
				m1, _ := baggage.NewMember("key1", "value1")
				bag, _ := baggage.New(m1)
				return bag
			},
		},
		{
			name:  "header+追加データを正常に受理できる",
			input: append(SignatureV1, []byte("key1=value1\r\ntesttest")...),
			wantBaggageFunc: func() baggage.Baggage {
				m1, _ := baggage.NewMember("key1", "value1")
				bag, _ := baggage.New(m1)
				return bag
			},
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
			opts := cmp.AllowUnexported(Header{}, baggage.Baggage{})
			if diff := cmp.Diff(NewV1(tt.wantBaggageFunc()), got, opts); diff != "" {
				t.Errorf("Parse() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
