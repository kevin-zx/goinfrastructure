package stringer

import (
	"reflect"
	"testing"
)

func TestRemoveDuplicateAndBlankString(t *testing.T) {
	type args struct {
		strs []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "1",
			args: args{
				strs: []string{"a", "b", "c", "", "c", "a"},
			},
			want: []string{"a", "b", "c"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := RemoveDuplicateAndBlankString(tt.args.strs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoveDuplicateAndBlankString() = %v, want %v", got, tt.want)
			}
		})
	}
}
