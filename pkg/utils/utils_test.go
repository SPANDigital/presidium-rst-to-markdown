package utils

import (
	"testing"
)

func TestGenerateSlug(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "one-word",
			args: args{
				input: "One",
			},
			want: "one",
		},
		{
			name: "two-word",
			args: args{
				input: "aLpHA BETA",
			},
			want: "alpha_beta",
		},
		{
			name: "three-word-with-tab",
			args: args{
				input: "a 	B c	e",
			},
			want: "a_b_c_e",
		},
		{
			name: `various-types-of-whitespace-tests`,
			args: args{
				input: `various types		of     whitespace
is used   here 		together`,
			},
			want: "various_types_of_whitespace_is_used_here_together",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateSlug(tt.args.input); got != tt.want {
				t.Errorf("GenerateSlug() = %v, want %v", got, tt.want)
			}
		})
	}
}
