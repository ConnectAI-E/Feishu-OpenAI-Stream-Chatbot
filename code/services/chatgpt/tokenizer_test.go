package chatgpt

import "testing"

func TestCalcTokenLength(t *testing.T) {
	type args struct {
		text string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "eng",
			args: args{
				text: "hello world",
			},
			want: 2,
		},
		{
			name: "cn",
			args: args{
				text: "我和我的祖国",
			},
			want: 13,
		},
		{
			name: "empty",
			args: args{
				text: "",
			},
			want: 0,
		},
		{
			name: "empty",
			args: args{
				text: " ",
			},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CalcTokenLength(tt.args.text); got != tt.want {
				t.Errorf("CalcTokenLength() = %v, want %v", got, tt.want)
			}
		})
	}
}
