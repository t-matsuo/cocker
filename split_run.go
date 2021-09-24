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

func recursiveSplitRun(
	scanner *bufio.Scanner,
	newDockerfile *[]byte,
	tmpDockerfile *[]byte,
	insideRun bool,
	needRUN bool,
	depth int) {

	if !scanner.Scan() {
		return
	}

	currentLine := scanner.Bytes()
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
		",isCommnet:", isComment,
		",insideRun:", insideRun,
	)

	if !insideRun && (isEmptyLine || isStartWithoutRun || (!isStartWithRun && isEndWithBackslashOnly) || isComment) {
		logDebug.Println("SPLIT OP : not inside RUN")
		clearTmpDockerfile(tmpDockerfile)
		appendLineln(newDockerfile, currentLine)
		recursiveSplitRun(
			scanner,
			newDockerfile,
			tmpDockerfile,
			insideRun,
			false,
			depth+1,
		)
		return
	}

	if insideRun {
		if isStartWithoutRun {
			logDebug.Println("SPLIT OP : end of RUN")
			appendLineln(tmpDockerfile, currentLine)
			return
		}
		logDebug.Println("SPLIT OP : insdie RUN")
		if isComment {
			recursiveSplitRun(
				scanner,
				newDockerfile,
				tmpDockerfile,
				insideRun,
				needRUN,
				depth+1,
			)
		} else if isEndWithBackslashOnly {
			recursiveSplitRun(
				scanner,
				newDockerfile,
				tmpDockerfile,
				insideRun,
				false,
				depth+1,
			)
		} else {
			recursiveSplitRun(
				scanner,
				newDockerfile,
				tmpDockerfile,
				insideRun,
				true,
				depth+1,
			)
		}
		if isComment {
			insertLineln(tmpDockerfile, currentLine)
			return
		}
		if !isStartWithRun && !needRUN {
			insertLineln(tmpDockerfile, deleteAmpersandAndBackslash(currentLine))
			return
		}
		if !isStartWithRun && isEndWithBackslashOnly {
			logDebug.Println("have BackSlash Only")
			insertLineln(tmpDockerfile, currentLine)
			return
		}
		if isEmptyLine {
			insertLineln(tmpDockerfile, currentLine)
			return
		}
		insertLineln(tmpDockerfile, addRun(deleteAmpersandAndBackslash(currentLine)))
		return
	}

	if needRUN == true {
		logErr.Fatal("BUG: needRUN = true")
	}
	if isStartWithRun {
		logDebug.Println("SPLIT OP : beginning of RUN")
		insideRun = true
		if isEndWithBackslashOnly {
			recursiveSplitRun(
				scanner,
				newDockerfile,
				tmpDockerfile,
				insideRun,
				false,
				depth+1,
			)

		} else {
			recursiveSplitRun(
				scanner,
				newDockerfile,
				tmpDockerfile,
				insideRun,
				true,
				depth+1,
			)
		}

		if isEndWithBackslashOnly {
			insertLineln(tmpDockerfile, currentLine)
		} else {
			insertLineln(tmpDockerfile, addRun(deleteAmpersandAndBackslash(currentLine)))
		}

		logDebug.Print("---- tmpDockerfile  ---- \n" + string(*tmpDockerfile))
		// read tmpDockerfile and add into newDockerfile
		tmpBuf := bytes.NewBuffer(*tmpDockerfile)
		tmpScanner := bufio.NewScanner(tmpBuf)
		for tmpScanner.Scan() {
			tmpLine := tmpScanner.Bytes()
			appendLineln(newDockerfile, tmpLine)
			tmpLine = tmpScanner.Bytes()
		}
		clearTmpDockerfile(tmpDockerfile)
		logDebug.Print("---- newDockerfile -----\n" + string(*newDockerfile))

		insideRun = false
		recursiveSplitRun(
			scanner,
			newDockerfile,
			tmpDockerfile,
			insideRun,
			needRUN,
			depth+1,
		)
		return
	}

	logDebug.Println("SPLIT OP : call default recursive")
	appendLineln(newDockerfile, currentLine)
	recursiveSplitRun(
		scanner,
		newDockerfile,
		tmpDockerfile,
		insideRun,
		needRUN,
		depth+1,
	)
	return
}

func splitRun() {
	logDebug.Println("Splitting Run")
	buf := bytes.NewBuffer(dockerfile)
	scanner := bufio.NewScanner(buf)
	newDockerfile := make([]byte, 0, 100000)
	var tmpDockerfile []byte
	clearTmpDockerfile(&tmpDockerfile)
	recursiveSplitRun(
		scanner,
		&newDockerfile,
		&tmpDockerfile,
		false, true, 1,
	)
	fmt.Print(string(newDockerfile))
}
