/*
Copyright 2021 MATSUO Takatoshi and Cocker Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import "testing"

func Test_getIncludeFilename(t *testing.T) {
	type args struct {
		line        []byte
		currentPath string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		{
			name:  "normal",
			args:  args{line: []byte(`Dockerfile.inc`), currentPath: "subpath"},
			want:  "subpath/Dockerfile.inc",
			want1: "subpath",
		},
		{
			name:  "subsub",
			args:  args{line: []byte(`Dockerfile.inc`), currentPath: "subpath/subsubpath"},
			want:  "subpath/subsubpath/Dockerfile.inc",
			want1: "subpath/subsubpath",
		},
		{
			name:  "fullpath",
			args:  args{line: []byte(`Dockerfile.inc`), currentPath: "/root/test/subpath"},
			want:  "/root/test/subpath/Dockerfile.inc",
			want1: "/root/test/subpath",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := getIncludeFilename(tt.args.line, tt.args.currentPath)
			if got != tt.want {
				t.Errorf("getIncludeFilename() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("getIncludeFilename() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
