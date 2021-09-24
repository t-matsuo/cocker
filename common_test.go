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
	"reflect"
	"testing"
)

func Test_addRun(t *testing.T) {
	type args struct {
		line []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "normal",
			args: args{line: []byte(`ls -la`)},
			want: []byte(`RUN ls -la`),
		},
		{
			name: "nothing",
			args: args{line: []byte(`RUN ls -la`)},
			want: []byte(`RUN ls -la`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := addRun(tt.args.line); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("addRun() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_haveAmpersandAndBackslash(t *testing.T) {
	type args struct {
		line []byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "normal",
			args: args{line: []byte(`aaaa   &&   \`)},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := haveAmpersandAndBackslash(tt.args.line); got != tt.want {
				t.Errorf("haveAmpersandAndBackslash() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_haveBackslashOnly(t *testing.T) {
	type args struct {
		line []byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "normal",
			args: args{line: []byte(`aaaa  \`)},
			want: true,
		},
		{
			name: "have Ampersand",
			args: args{line: []byte(`aaaa   &&   \`)},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := haveBackslashOnly(tt.args.line); got != tt.want {
				t.Errorf("haveBackslashOnly() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_deleteRun(t *testing.T) {
	type args struct {
		line []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "normal",
			args: args{line: []byte(`ls -la`)},
			want: []byte(`ls -la`),
		},
		{
			name: "nothing",
			args: args{line: []byte(`RUN ls -la`)},
			want: []byte(`    ls -la`),
		},
		{
			name: "empty line",
			args: args{line: []byte(``)},
			want: []byte(``),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := deleteRun(tt.args.line); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("deleteRun() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_addAmpasandAndBackslash(t *testing.T) {
	type args struct {
		line []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "normal",
			args: args{line: []byte(`RUN ls -la`)},
			want: []byte(`RUN ls -la && \`),
		},
		{
			name: "nothing",
			args: args{line: []byte(`RUN ls -la && \`)},
			want: []byte(`RUN ls -la && \`),
		},
		{
			name: "some spaces and tabs",
			args: args{line: []byte(`RUN ls -la    	&&   	\`)},
			want: []byte(`RUN ls -la    	&&   	\`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := addAmpersandAndBackslash(tt.args.line); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("addBackslash() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func Test_deleteAmpersandAndBackslash(t *testing.T) {
	type args struct {
		line []byte
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "normal",
			args: args{line: []byte(`RUN ls -la && \`)},
			want: []byte(`RUN ls -la`),
		},
		{
			name: "nothing",
			args: args{line: []byte(`RUN ls -la`)},
			want: []byte(`RUN ls -la`),
		},
		{
			name: "some_spaces1",
			args: args{line: []byte(`RUN ls -la   && \`)},
			want: []byte(`RUN ls -la`),
		},
		{
			name: "some_spaces2",
			args: args{line: []byte(`RUN ls -la   &&   \  `)},
			want: []byte(`RUN ls -la`),
		},
		{
			name: "empty line",
			args: args{line: []byte(``)},
			want: []byte(``),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := deleteAmpersandAndBackslash(tt.args.line); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("deleteAmpersandAndBackslash() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}

func Test_appendtLineln(t *testing.T) {
	type args struct {
		target *[]byte
		line   []byte
	}
	tgt := []byte("FROM centos:8\n")
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "normal",
			args: args{target: &tgt, line: []byte("RUN ls -la")},
			want: []byte("FROM centos:8\nRUN ls -la\n"),
		},
		{
			name: "append empty []byte",
			args: args{target: &tgt, line: []byte("")},
			want: []byte("FROM centos:8\n\n"),
		},
		{
			name: "append nil",
			args: args{target: &tgt, line: nil},
			want: []byte("FROM centos:8\n"),
		},
	}
	for _, tt := range tests {
		tgt = []byte("FROM centos:8\n")
		t.Run(tt.name, func(t *testing.T) {
			appendLineln(tt.args.target, tt.args.line)
			if !reflect.DeepEqual(*tt.args.target, tt.want) {
				t.Errorf("appendLineln() = %v, want %v", string(*tt.args.target), string(tt.want))
			}
		})
	}
}

func Test_insertLineln(t *testing.T) {
	type args struct {
		target *[]byte
		line   []byte
	}
	tgt := []byte("FROM centos:8\n")
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "normal",
			args: args{target: &tgt, line: []byte("RUN ls -la")},
			want: []byte("RUN ls -la\nFROM centos:8\n"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			insertLineln(tt.args.target, tt.args.line)
			if !reflect.DeepEqual(*tt.args.target, tt.want) {
				t.Errorf("insertLineln() = %v, want %v", string(*tt.args.target), string(tt.want))
			}
		})
	}
}

func Test_haveWithoutRun(t *testing.T) {
	type args struct {
		line []byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "RUN",
			args: args{line: []byte("RUN ls -la")},
			want: false,
		},
		{
			name: "LABEL",
			args: args{line: []byte("LABEL aaa=aaa")},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := haveWithoutRun(tt.args.line); got != tt.want {
				t.Errorf("haveWithoutRun() = %v, want %v", got, tt.want)
			}
		})
	}
}
