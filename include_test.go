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

import (
	"testing"
)

func Test_haveIncludeComment(t *testing.T) {
	type args struct {
		line []byte
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		want1 string
		want2 int
		want3 string
	}{
		{
			name:  "normal",
			args:  args{line: []byte(`#include Dockerfile.inc`)},
			want:  true,
			want1: "Dockerfile.inc",
			want2: noDef,
			want3: "",
		},
		{
			name:  "ifdef",
			args:  args{line: []byte(`#include Dockerfile.inc ifdef TEST_ENV`)},
			want:  true,
			want1: "Dockerfile.inc",
			want2: ifDef,
			want3: "TEST_ENV",
		},
		{
			name:  "ifndef",
			args:  args{line: []byte(`#include Dockerfile.inc ifndef TEST_ENV`)},
			want:  true,
			want1: "Dockerfile.inc",
			want2: ifNotDef,
			want3: "TEST_ENV",
		},
		{
			name:  "ifdef with some space",
			args:  args{line: []byte(`#include   Dockerfile.inc  ifdef  TEST_ENV`)},
			want:  true,
			want1: "Dockerfile.inc",
			want2: ifDef,
			want3: "TEST_ENV",
		},
		{
			name:  "ifndef",
			args:  args{line: []byte(`#include Dockerfile.inc ifndef TEST_ENV`)},
			want:  true,
			want1: "Dockerfile.inc",
			want2: ifNotDef,
			want3: "TEST_ENV",
		},
		{
			name:  "invalid args1",
			args:  args{line: []byte(`#include`)},
			want:  false,
			want1: "",
			want2: noDef,
			want3: "",
		},
		{
			name:  "invalid args3",
			args:  args{line: []byte(`#include Dockerfile.inc ifdef`)},
			want:  false,
			want1: "",
			want2: noDef,
			want3: "",
		},
		{
			name:  "invalid args5",
			args:  args{line: []byte(`#include Dockerfile.inc ifdef TEST_ENV foo`)},
			want:  false,
			want1: "",
			want2: noDef,
			want3: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2, got3 := haveIncludeComment(tt.args.line)
			if got != tt.want {
				t.Errorf("haveIncludeComment() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("haveIncludeComment() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("haveIncludeComment() got2 = %v, want %v", got2, tt.want2)
			}
			if got3 != tt.want3 {
				t.Errorf("haveIncludeComment() got3 = %v, want %v", got3, tt.want3)
			}
		})
	}
}

func Test_getIncludeFilename(t *testing.T) {
	type args struct {
		line        string
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
			args:  args{line: "Dockerfile.inc", currentPath: "subpath"},
			want:  "subpath/Dockerfile.inc",
			want1: "subpath",
		},
		{
			name:  "subsub",
			args:  args{line: "Dockerfile.inc", currentPath: "subpath/subsubpath"},
			want:  "subpath/subsubpath/Dockerfile.inc",
			want1: "subpath/subsubpath",
		},
		{
			name:  "fullpath",
			args:  args{line: "Dockerfile.inc", currentPath: "/root/test/subpath"},
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
