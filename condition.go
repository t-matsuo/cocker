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
	"os"
	"regexp"
)

func haveStartConditionComment(line []byte) (bool, int, string) {
	mode := noDef
	env := ""
	isIfDef, _ := regexp.Match(`^#ifdef `, line)
	if isIfDef {
		mode = ifDef
	}
	isIfNotDef, _ := regexp.Match(`^#ifndef `, line)
	if isIfNotDef {
		mode = ifNotDef
	}
	if mode == noDef {
		return false, noDef, env
	}

	regSpaceSplit := regexp.MustCompile(`\s* \s*`)
	regSpaceSplitResult := regSpaceSplit.Split(string(line), -1)
	if len(regSpaceSplitResult) == 1 ||
		len(regSpaceSplitResult) > 2 {
		return false, mode, env
	}
	env = regSpaceSplitResult[1]
	if env == "" {
		return false, mode, env
	}
	return true, mode, env
}

func haveEndConditionComment(line []byte) bool {
	isEndIf, _ := regexp.Match(`^#endif`, line)
	if isEndIf {
		return true
	}
	return false
}

func handleConditionRecursive(
	scanner *bufio.Scanner,
	newDockerfile *[]byte,
	currentMode int,
	isSkip bool,
	depth int) {

	logDebug.Printf("--------------------- %d --------------------\n", depth)
	defer logDebug.Printf("--------------------- defer %d ----------------\n", depth)

	for scanner.Scan() {
		currentLine := make([]byte, len(scanner.Bytes()))
		copy(currentLine, scanner.Bytes())
		logDebug.Println(("currentLine: " + string(currentLine)))

		if haveEndConditionComment(currentLine) {
			logDebug.Println("CONDITION OP: have #endif")
			return
		}

		if (currentMode != ifDef && currentMode != ifNotDef) && haveEndConditionComment(currentLine) {
			logErr.Fatalln("invalid #endif")
		}

		isCondition, mode, env := haveStartConditionComment(currentLine)
		if isCondition {
			if mode == ifDef {
				_, val := os.LookupEnv(env)
				if val == true {
					logDebug.Println("CONDITION OP: Start ifdef and match condition")
					handleConditionRecursive(scanner, newDockerfile, mode, isSkip, depth+1)
				} else {
					logDebug.Println("CONDITION OP: Start ifdef but not match condition")
					handleConditionRecursive(scanner, newDockerfile, mode, true, depth+1)
				}
				continue
			}
			if mode == ifNotDef {
				_, val := os.LookupEnv(env)
				if val == true {
					logDebug.Println("CONDITION OP: Start ifndef but not match condition")
					handleConditionRecursive(scanner, newDockerfile, mode, true, depth+1)
				} else {
					logDebug.Println("CONDITION OP: Start ifndef and match condition")
					handleConditionRecursive(scanner, newDockerfile, mode, isSkip, depth+1)
				}
				continue
			}
		}
		if isSkip {
			logDebug.Println("CONDITION OP: Skipping")
			handleConditionRecursive(scanner, newDockerfile, currentMode, isSkip, depth+1)
			return
		}
		logDebug.Println("CONDITION OP: adding currentLine: " + string(currentLine))
		appendLineln(newDockerfile, currentLine)
	}
	if currentMode != noDef {
		logErr.Fatalln("no #endif")
	}
	return
}

func handleCondition() {
	logDebug.Println("handling ifdef,ifndef")

	buf := bytes.NewBuffer(dockerfile)
	scanner := bufio.NewScanner(buf)
	newDockerfile := make([]byte, 0, 100000)
	var tmpDockerfile []byte
	clearTmpDockerfile(&tmpDockerfile)

	handleConditionRecursive(
		scanner,
		&newDockerfile,
		noDef,
		false,
		1,
	)
	dockerfile = newDockerfile
	if flagCondition {
		outputDockerFile()
	}
}
