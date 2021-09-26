/*
Copyright 2021 MATSUO Takatoshi and cocker Authors

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
	"fmt"
)

// return true if current line need '&& \'
func recursiveMergeRun(
	scanner *bufio.Scanner,
	newDockerfile *[]byte,
	tmpDockerfile *[]byte,
	insideRun bool,
	depth int) bool {

	if !scanner.Scan() {
		return false
	}

	currentLine := make([]byte, len(scanner.Bytes()))
	copy(currentLine, scanner.Bytes())
	isStartWithRun := haveRun(currentLine)
	isEndWithAmpBackslash := haveAmpersandAndBackslash(currentLine)
	isEndWithBackslashOnly := haveBackslashOnly(currentLine)
	isEmptyLine := haveEmptyLine(currentLine)
	isComment := haveComment(currentLine)
	isStartWithoutRun := haveWithoutRun(currentLine)

	logDebug.Printf("------------------------ %d -------------------------\n", depth)
	defer logDebug.Printf("------------------------ defer %d -------------------\n", depth)
	logDebug.Println(
		"currentLine: ", string(currentLine),
	)
	logDebug.Println(
		"isStartWithRun:", isStartWithRun,
		",isEndWithAmpBackslash:", isEndWithAmpBackslash,
		",isEndWithBackslashOnly:", isEndWithBackslashOnly,
		",isEmptyLine:", isEmptyLine,
		",isComment:", isComment,
		",insideRun:", insideRun,
	)

	if !insideRun && (isEmptyLine || isStartWithoutRun || (!isStartWithRun && isEndWithBackslashOnly) || isComment) {
		logDebug.Println("MERGE OP: not inside RUN")
		clearTmpDockerfile(tmpDockerfile)
		appendLineln(newDockerfile, currentLine)
		recursiveMergeRun(
			scanner,
			newDockerfile,
			tmpDockerfile,
			insideRun,
			depth+1,
		)
		return false
	}

	if insideRun {
		if isStartWithoutRun {
			logDebug.Println("MERGE OP : end of RUN")
			appendLineln(tmpDockerfile, currentLine)
			return false
		}
		logDebug.Println("MERGE OP : insdie RUN")
		needAmpBackSlash := recursiveMergeRun(
			scanner,
			newDockerfile,
			tmpDockerfile,
			insideRun,
			depth+1,
		)
		if isComment {
			insertLineln(tmpDockerfile, currentLine)
		} else if needAmpBackSlash {
			if !isEmptyLine {
				insertLineln(tmpDockerfile, addAmpersandAndBackslash(currentLine))
			}
		} else {
			insertLineln(tmpDockerfile, currentLine)
			if haveRun(currentLine) {
				needAmpBackSlash = true
			}
		}
		return needAmpBackSlash
	}

	if isStartWithRun {
		logDebug.Println("MERGE OP : beginning of RUN")
		insideRun = true
		needAmpBack := recursiveMergeRun(
			scanner,
			newDockerfile,
			tmpDockerfile,
			insideRun,
			depth+1,
		)
		logDebug.Println("need ampasand and backslash :", needAmpBack)

		if needAmpBack {
			appendLineln(newDockerfile, addAmpersandAndBackslash(currentLine))
		} else {
			appendLineln(newDockerfile, currentLine)
		}

		logDebug.Print("---- tmpDockerfile  ---- \n" + string(*tmpDockerfile))
		// read tmpDockerfile and add into newDockerfile
		tmpBuf := bytes.NewBuffer(*tmpDockerfile)
		tmpScanner := bufio.NewScanner(tmpBuf)
		for tmpScanner.Scan() {
			tmpLine := tmpScanner.Bytes()
			appendLineln(newDockerfile, deleteRun(tmpLine))
		}
		clearTmpDockerfile(tmpDockerfile)
		logDebug.Print("---- newDockerfile -----\n" + string(*newDockerfile))

		insideRun = false
		recursiveMergeRun(
			scanner,
			newDockerfile,
			tmpDockerfile,
			insideRun,
			depth+1,
		)
		return false
	}

	logDebug.Println("MERGE OP : call default recursive")
	appendLineln(newDockerfile, currentLine)
	recursiveMergeRun(
		scanner,
		newDockerfile,
		tmpDockerfile,
		insideRun,
		depth+1,
	)
	return false
}

func mergeRun() {
	logDebug.Println("Merging Run")
	buf := bytes.NewBuffer(dockerfile)
	scanner := bufio.NewScanner(buf)
	newDockerfile := make([]byte, 0, 100000)
	var tmpDockerfile []byte
	clearTmpDockerfile(&tmpDockerfile)
	recursiveMergeRun(
		scanner,
		&newDockerfile,
		&tmpDockerfile,
		false, 1,
	)
	fmt.Print(string(newDockerfile))
}
