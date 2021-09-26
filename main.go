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
	"io/ioutil"
	"log"
	"os"

	"github.com/kelseyhightower/envconfig"
)

var (
	version string = "0.0.0"
	commit  string = ""
)

var goenv struct {
	Debug     bool   `default:"false"`
	LogPrefix string `default:"[cocker]"`
}

var logInfo *log.Logger
var logWarn *log.Logger
var logErr *log.Logger
var logDebug *log.Logger

var dockerfile []byte
var dockerfilePath string

func handleEnv() {
	if err := envconfig.Process("CC", &goenv); err != nil {
		logErr.Fatalf("Failed to process env: %s", err)
		os.Exit(1)
	}

	// setup log outputs
	if goenv.Debug || flagDebug {
		logInfo.SetFlags(log.LstdFlags | log.Llongfile | log.Lmsgprefix)
		logWarn.SetFlags(log.LstdFlags | log.Llongfile | log.Lmsgprefix)
		logErr.SetFlags(log.LstdFlags | log.Llongfile | log.Lmsgprefix)
		logDebug.SetOutput(os.Stderr)
	}

	logInfo.SetPrefix(goenv.LogPrefix + " INFO ")
	logWarn.SetPrefix(goenv.LogPrefix + " WARN ")
	logErr.SetPrefix(goenv.LogPrefix + " ERROR ")
	logDebug.SetPrefix(goenv.LogPrefix + " DEBUG ")

}

func showHelp() {
	fmt.Println("")
	fmt.Println("Cocker is pre processor for Dockerfile.")
	fmt.Println("")
	fmt.Println("Usage:")
	fmt.Println("  $ cocker [options...] filename")
	fmt.Println("  $ cat Dockerfile | cocker [options...]")
	fmt.Println("")
	fmt.Println("Options:")
	fmt.Println("   -m --merge   : Merge RUN mode (cannot use -s option)")
	fmt.Println("   -s --split   : Split RUN mode (cannot use -m options")
	fmt.Println("   -i --include : Include another Dockerfile")
	fmt.Println("   -d --debug   : Print debug messages")
	fmt.Println("   --version    : Show version number")
	fmt.Println("   -h --help    : Show help")
	fmt.Println("")
	fmt.Println("Environment Variables:")
	fmt.Println("   CC_DEBUG=true : Print debug messages (=--debug option)")
	fmt.Println("")
}

func init() {
	logInfo = log.New(os.Stdout, "[cocker] INFO ", log.LstdFlags|log.Lmsgprefix)
	logErr = log.New(os.Stderr, "[cocker] ERROR ", log.LstdFlags|log.Lmsgprefix)
	logWarn = log.New(os.Stderr, "[cocker] WARN ", log.LstdFlags|log.Lmsgprefix)
	logDebug = log.New(ioutil.Discard, "[cocker] DEBUG ", log.LstdFlags|log.Llongfile|log.Lmsgprefix)
}

var (
	flagMerge   bool
	flagSplit   bool
	flagInclude bool
	flagDebug   bool
	flagHelp    bool
	flagVersion bool
)

func setupFlags() {
	flagDebug = true
	flag.BoolVar(&flagMerge, "m", false, "Merge RUN mode")
	flag.BoolVar(&flagSplit, "s", false, "Split RUN mode")
	flag.BoolVar(&flagInclude, "i", false, "Include Dockerfile using '#include filename' comment")
	flag.BoolVar(&flagDebug, "d", false, "Print debug messages")
	flag.BoolVar(&flagMerge, "merge", false, "Merge RUN mode (=-m)")
	flag.BoolVar(&flagSplit, "split", false, "Split RUN mode (=-s)")
	flag.BoolVar(&flagInclude, "include", false, "Include Dockerfile using '#include filename' comment (=-i)")
	flag.BoolVar(&flagDebug, "debug", false, "Print debug messages (=-d)")
	flag.BoolVar(&flagHelp, "h", false, "Show help")
	flag.BoolVar(&flagVersion, "version", false, "Show version")
	flag.Parse()
}

func main() {
	setupFlags()
	handleEnv()

	if flagVersion {
		fmt.Println("Cocker " + version)
		fmt.Println("Commit " + commit)
		fmt.Println("Source https://github.com/t-matsuo/cocker")
		os.Exit(0)
	}

	if flagHelp || (!flagMerge && !flagSplit && !flagInclude) || (flagMerge && flagSplit) {
		showHelp()
		os.Exit(0)
	}

	readDockerFile()
	if flagInclude {
		includeDockerfile()
	}
	if flagMerge {
		mergeRun()
	}
	if flagSplit {
		splitRun()
	}
}
