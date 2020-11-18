package reader

import (
	"io"
	"strings"
	"testing"
)

func TestReadLines(t *testing.T) {
	type args struct {
		r io.Reader
		f func(line []byte)
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "normal line",
			args: args{
				r: strings.NewReader(`1
2
3`),
				f: func([]byte) {},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ReadLines(tt.args.r, tt.args.f); (err != nil) != tt.wantErr {
				t.Errorf("ReadLines() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
