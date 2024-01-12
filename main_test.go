package main

import (
	"reflect"
	"testing"
)

func Test_readBytesToPluginMap(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want map[string]struct{}
	}{
		{
			name: "plain",
			in: `
plugin1
plugin2
`,
			want: map[string]struct{}{
				"plugin1": {},
				"plugin2": {},
			},
		},
		{
			name: "with_comments",
			in: `
# comment
plugin1 # comment after plugin
# disabled-plugin
	# empty line but with whitespace
`,
			want: map[string]struct{}{
				"plugin1": {},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := readBytesToPluginMap([]byte(tt.in)); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("readBytesToPluginMap() = %v, want %v", got, tt.want)
			}
		})
	}
}
