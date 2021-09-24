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
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
)

func readFile(filename string) []byte {
	var r io.Reader
	f, err := os.Open(filename)
	defer f.Close()
	if err != nil {
		logErr.Fatalf("%s not found\n", filename)
		return nil
	}
	r = f
	file, _ := ioutil.ReadAll(r)
	return file
}

func haveIncludeCommnet(line []byte) bool {
	isInclude, _ := regexp.Match(`^#include `, line)
	return isInclude
}

func getIncludeFilename(line []byte) string {
	regGetIncludeFilnename := regexp.MustCompile(`^#include *`)
	filename := string(regGetIncludeFilnename.ReplaceAll(line, []byte(``)))
	logDebug.Println("dockerfilePath=", dockerfilePath)
	logDebug.Println("filename=", filename)
	if filepath.IsAbs(filename) {
		return filename
	}
	return dockerfilePath + "/" + filename
}

func includeDockerfileRecursiveFile(filename string, depth int) []byte {
	file := readFile(filename)
	return includeDockerfileRecursive(file, depth)
}

func includeDockerfileRecursive(file []byte, depth int) []byte {
	logDebug.Printf("--------------------- %d --------------------\n", depth)
	defer logDebug.Printf("--------------------- defer %d ----------------\n", depth)
	newDockerfile := make([]byte, 0, 100000)

	buf := bytes.NewBuffer(file)
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		line := scanner.Bytes()
		logDebug.Println(("line: " + string(line)))
		if haveIncludeCommnet(line) {
			subFile := includeDockerfileRecursiveFile(getIncludeFilename(line), depth+1)
			logDebug.Println("adding file" + string(line))
			appendLine(&newDockerfile, subFile)
			logDebug.Print("newDockerfile\n" + string(*&newDockerfile))
			continue
		}
		logDebug.Println("adding", string(line))
		appendLineln(&newDockerfile, line)
		logDebug.Print("newDockerfile\n" + string(*&newDockerfile))
	}
	return newDockerfile
}

func includeDockerfile() {
	logDebug.Println("Including Dockerfile")
	newDockerfile := includeDockerfileRecursive(dockerfile, 1)
	dockerfile = newDockerfile
	if !flagMerge && !flagSplit {
		outputDockerFile()
	}
}
