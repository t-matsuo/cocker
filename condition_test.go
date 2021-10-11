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

func Test_haveStartConditionComment(t *testing.T) {
	type args struct {
		line []byte
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		want1 int
		want2 string
	}{
		{
			name:  "normal ifdef",
			args:  args{line: []byte("#ifdef TEST_ENV")},
			want:  true,
			want1: ifDef,
			want2: "TEST_ENV",
		},
		{
			name:  "normal ifndef",
			args:  args{line: []byte("#ifndef TEST_ENV")},
			want:  true,
			want1: ifNotDef,
			want2: "TEST_ENV",
		},
		{
			name:  "invalid args",
			args:  args{line: []byte("#ifdef")},
			want:  false,
			want1: noDef,
			want2: "",
		},
		{
			name:  "invalid args with space",
			args:  args{line: []byte("#ifdef ")},
			want:  false,
			want1: ifDef,
			want2: "",
		},
		{
			name:  "many space",
			args:  args{line: []byte("#ifdef   TEST_ENV")},
			want:  true,
			want1: ifDef,
			want2: "TEST_ENV",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := haveStartConditionComment(tt.args.line)
			if got != tt.want {
				t.Errorf("haveStartConditionComment() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("haveStartConditionComment() got1 = %v, want %v", got1, tt.want1)
			}
			if got2 != tt.want2 {
				t.Errorf("haveStartConditionComment() got2 = %v, want %v", got2, tt.want2)
			}
		})
	}
}
