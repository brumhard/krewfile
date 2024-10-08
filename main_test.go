package main

import (
	"reflect"
	"testing"
)

func Test_readPluginsFromKrew(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		wantPluginMap PluginMap
	}{
		{
			name: "krew index output",
			input: `
PLUGIN             VERSION
node-shell         v1.10.1
resource-capacity  v0.8.0
`,
			wantPluginMap: PluginMap{
				"node-shell":        struct{}{},
				"resource-capacity": struct{}{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := readPluginsFromKrew([]byte(tt.input)); !reflect.DeepEqual(got, tt.wantPluginMap) {
				t.Errorf("readIndexesFromKrew() = %v, want %v", got, tt.wantPluginMap)
			}
		})
	}
}

func Test_readIndexesFromKrew(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		wantIndexMap IndexMap
	}{
		{
			name: "krew index output",
			input: `
INDEX     URL
default   https://github.com/kubernetes-sigs/krew-index.git
netshoot  https://github.com/nilic/kubectl-netshoot.git
`,
			wantIndexMap: IndexMap{
				"default":  "https://github.com/kubernetes-sigs/krew-index.git",
				"netshoot": "https://github.com/nilic/kubectl-netshoot.git",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := readIndexesFromKrew([]byte(tt.input)); !reflect.DeepEqual(got, tt.wantIndexMap) {
				t.Errorf("readIndexesFromKrew() = %v, want %v", got, tt.wantIndexMap)
			}
		})
	}
}

func Test_readKrewfile(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		wantPluginMap PluginMap
		wantIndexMap  IndexMap
		wantErr       error
	}{
		{
			name: "plain",
			input: `
index index1 https://myindex1.url.home.arpa/
index index2 https://myindex2.url.home.arpa/

index1/plugin1
plugin2
      `,
			wantPluginMap: PluginMap{"index1/plugin1": {}, "plugin2": {}},
			wantIndexMap: IndexMap{
				"index1": "https://myindex1.url.home.arpa/",
				"index2": "https://myindex2.url.home.arpa/",
			},
			wantErr: nil,
		}, {
			name: "with comments",
			input: `
# index comment
index index1 https://myindex1.url.home.arpa/ # comment after index
# index disabled-index https://myindex2.url.home.arpa/

plugin1
# disabled-plugin
  # empty line but with whitespace
      `,
			wantPluginMap: PluginMap{"plugin1": {}},
			wantIndexMap: IndexMap{
				"index1": "https://myindex1.url.home.arpa/",
			},
			wantErr: nil,
		}, {
			name: "with error",
			input: `
invalid line
      `,
			wantPluginMap: nil,
			wantIndexMap:  nil,
			wantErr:       InvalidKrewfileLineError{line: "invalid line"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPluginMap, gotIndexMap, gotErr := readKrewfile([]byte(tt.input))

			if !reflect.DeepEqual(gotPluginMap, tt.wantPluginMap) {
				t.Errorf("readKrewfile()/PluginMap = %v, want %v", gotPluginMap, tt.wantPluginMap)
			}

			if !reflect.DeepEqual(gotIndexMap, tt.wantIndexMap) {
				t.Errorf("readKrewfile()/IndexMap = %v, want %v", gotIndexMap, tt.wantIndexMap)
			}

			if gotErr == nil {
				if tt.wantErr != nil {
					t.Errorf("readKrewfile()/err = %v, want %v", gotErr, tt.wantErr)
				}

				return
			}

			if gotErr.Error() != tt.wantErr.Error() {
				t.Errorf("readKrewfile()/err = %q, want %q", gotErr, tt.wantErr)
			}
		})
	}
}
