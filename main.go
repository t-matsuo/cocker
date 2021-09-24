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

var goenv struct {
	Debug      bool   `default:"false"`
	Log_prefix string `default:"[cocker]"`
}

var log_info *log.Logger
var log_warn *log.Logger
var log_err *log.Logger
var log_debug *log.Logger

var Dockerfile []byte
var DockerfilePath string

func handleEnv() {
	if err := envconfig.Process("CC", &goenv); err != nil {
		log_err.Fatalf("Failed to process env: %s", err)
		os.Exit(1)
	}

	// setup log outputs
	if goenv.Debug || flagDebug {
		log_info.SetFlags(log.LstdFlags | log.Llongfile | log.Lmsgprefix)
		log_warn.SetFlags(log.LstdFlags | log.Llongfile | log.Lmsgprefix)
		log_err.SetFlags(log.LstdFlags | log.Llongfile | log.Lmsgprefix)
		log_debug.SetOutput(os.Stderr)
	}

	log_info.SetPrefix(goenv.Log_prefix + " INFO ")
	log_warn.SetPrefix(goenv.Log_prefix + " WARN ")
	log_err.SetPrefix(goenv.Log_prefix + " ERROR ")
	log_debug.SetPrefix(goenv.Log_prefix + " DEBUG ")

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
	fmt.Println("   -m --merge : Merge RUN mode (cannot use -s option)")
	fmt.Println("   -s --split : Split RUN mode (cannot use -m options")
	fmt.Println("   -i --include : Include another Dockerfile")
	fmt.Println("   -d --debug : Print debug messages")
	fmt.Println("   -h --help  : Show help")
	fmt.Println("")
	fmt.Println("Environment Variables:")
	fmt.Println("   CC_DEBUG=true : Print debug messages (=--debug option)")
	fmt.Println("")
}

func init() {
	log_info = log.New(os.Stdout, "[cocker] INFO ", log.LstdFlags|log.Lmsgprefix)
	log_err = log.New(os.Stderr, "[cocker] ERROR ", log.LstdFlags|log.Lmsgprefix)
	log_warn = log.New(os.Stderr, "[cocker] WARN ", log.LstdFlags|log.Lmsgprefix)
	log_debug = log.New(ioutil.Discard, "[cocker] DEBUG ", log.LstdFlags|log.Llongfile|log.Lmsgprefix)
}

var (
	flagMerge   bool
	flagSplit   bool
	flagInclude bool
	flagDebug   bool
	flagHelp    bool
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
	flag.Parse()
}

func main() {
	setupFlags()
	handleEnv()

	if flagHelp || (!flagMerge && !flagSplit && !flagInclude) || (flagMerge && flagSplit) {
		showHelp()
		os.Exit(0)
	}

	ReadDockerFile()
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
