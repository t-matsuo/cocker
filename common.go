/*
Copyright 2021 MATSUO Takatoshi and cocker Author

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
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

const (
	noDef    = 0
	ifDef    = 1
	ifNotDef = 2
)

func readDockerFile() {
	var filename string
	var r io.Reader

	if args := flag.Args(); len(args) > 0 {
		filename = args[0]
	}
	switch filename {
	case "-":
		r = os.Stdin
		dockerfilePath = "."
	case "":
		if terminal.IsTerminal(int(syscall.Stdin)) {
			logErr.Fatalln("File is not specified")
		}
		r = os.Stdin
		dockerfilePath = "."
	default:
		f, err := os.Open(filename)
		defer f.Close()
		if err != nil {
			logErr.Fatalf("%s not found\n", filename)
		}
		r = f
		dockerfilePath = filepath.Dir(filename)
	}
	dockerfile, _ = ioutil.ReadAll(r)
}

func outputDockerFile() {
	fmt.Print(string(dockerfile))
}

func haveRun(line []byte) bool {
	isStartWithRun, _ := regexp.Match(`^RUN`, line)
	return isStartWithRun
}

func addRun(line []byte) []byte {
	if !haveRun(line) {
		regAddRun := regexp.MustCompile(`^ *`)
		return regAddRun.ReplaceAll(line, []byte(`RUN `))
	}
	return line
}

func deleteRun(line []byte) []byte {
	if len(line) == 0 {
		return line
	}
	regDeleteBackslash := regexp.MustCompile(`^RUN *`)
	return regDeleteBackslash.ReplaceAll(line, []byte(`    `))
}

func haveAmpersandAndBackslash(line []byte) bool {
	isEndWithAmpBackslash, _ := regexp.Match(`&&[ ,	]*\\[ ,	]*$`, line)
	return isEndWithAmpBackslash
}

func haveBackslashOnly(line []byte) bool {
	if haveAmpersandAndBackslash(line) {
		return false
	}
	isEndWithAmpBackslashOnly, _ := regexp.Match(`\\$`, line)
	return isEndWithAmpBackslashOnly
}

func addAmpersandAndBackslash(line []byte) []byte {
	if haveBackslashOnly(line) {
		return line
	}
	if !haveAmpersandAndBackslash(line) {
		regAddBackslash := regexp.MustCompile(`$`)
		return regAddBackslash.ReplaceAll(line, []byte(` && \`))
	}
	return line
}

func deleteAmpersandAndBackslash(line []byte) []byte {
	if len(line) == 0 {
		return line
	}
	regDeleteBackslash := regexp.MustCompile(`[ ,	]*&&[ ,	]*\\[ ,	]*$`)
	return regDeleteBackslash.ReplaceAll(line, []byte(``))
}

func haveEmptyLine(line []byte) bool {
	isEmptyLine, _ := regexp.Match(`^[ ,	]*$`, line)
	return isEmptyLine
}

func haveWithoutRun(line []byte) bool {
	isStartWithoutRun, _ := regexp.Match(`^FROM|^CMD|^LABEL|^MAINTAINER|^EXPOSE|^ENV|^ADD|^COPY|^ENTRYPOINT|^VOLUME|^USER|^WORKDIR|^ARG|^ONBUILD|^STOPSIGNAL|^HEALTHCHECK|^SHELL`, line)
	return isStartWithoutRun
}

func haveComment(line []byte) bool {
	isComment, _ := regexp.Match(`^[ ,	]*#`, line)
	return isComment
}

func appendLine(target *[]byte, line []byte) {
	if line == nil {
		return
	}
	*target = append(*target, line...)
}

func appendLineln(target *[]byte, line []byte) {
	if line == nil {
		return
	}
	*target = append(*target, line...)
	*target = append(*target, []byte("\n")...)
}

func insertLineln(target *[]byte, line []byte) {
	newLine := make([]byte, len(line))
	copy(newLine, line)
	newLine = append(newLine, []byte("\n")...)
	*target = append(newLine, *target...)
}

func clearTmpDockerfile(target *[]byte) {
	*target = nil
}
