package reader

import (
	"io"
	"strings"
	"testing"
)

func TestCsv(t *testing.T) {
	type args struct {
		reader io.Reader
		f      func(rcs []string)
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nomral",
			args: args{
				reader: strings.NewReader(`name,age,height,weight
kevin,18,175,130`),
				f: func(rcs []string) {},
			},
			wantErr: false,
		}, {
			name: "wrong fmt csv",
			args: args{
				reader: strings.NewReader(`name,age,height,weight
kevin,18,175,130
somebody,1
`),
				f: func(rcs []string) {},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Csv(tt.args.reader, tt.args.f); (err != nil) != tt.wantErr {
				t.Errorf("Csv() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
