// Copyright 2016 HenryLee. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/henrylee2cn/think/model"
	"github.com/henrylee2cn/thinkgo"
)

var appname string
var crupath, _ = os.Getwd()
var apptpl = "simple"

func main() {
	thinkgo.RemoveUseless()
	if len(os.Args) < 2 {
		help()
		return
	}
	switch os.Args[1] {
	case "new":
		newapp(os.Args[2:])
	case "run":
		runapp(os.Args[2:])
	}
}

func newapp(args []string) {
	switch len(args) {
	case 1:
		initVar(args)
	case 2:
		initVar(args)
		apptpl = args[1]
	default:
		newappHelp()
		return
	}
	thinkgo.Printf("[think] Create a thinkgo project named `%s` in the `%s` path.", appname, crupath)
	if isExist(crupath) {
		thinkgo.Printf("[think] The project path has conflic, do you want to build in: %s\n", crupath)
		thinkgo.Printf("[think] Do you want to overwrite it? [yes|no]]  ")
		if !askForConfirmation() {
			thinkgo.Fatalf("[think] Cancel...")
			return
		}
	}

	thinkgo.Printf("[think] Start create project...")

	switch apptpl {
	case "simple":
		model.SimplePro(crupath, appname)
	default:
		thinkgo.Fatalf("[think] `%s` template does not exist, reference:\n[simple]\n", apptpl)
	}

	thinkgo.Printf("[think] Create was successful")

	if err := os.Chdir(crupath); err != nil {
		thinkgo.Fatalf("[think] Create project fail: %v", err)
	}
	exit := make(chan bool)
	autobuild()
	newWatcher()
	for {
		select {
		case <-exit:
			runtime.Goexit()
		}
	}
}

func runapp(args []string) {
	switch len(args) {
	case 0, 1:
		initVar(args)
	default:
		runappHelp()
		return
	}
	if err := os.Chdir(crupath); err != nil {
		thinkgo.Fatalf("[think] Create project fail: %v", err)
	}
	exit := make(chan bool)
	autobuild()
	newWatcher()
	for {
		select {
		case <-exit:
			runtime.Goexit()
		}
	}
}

const helpInfo = `Think Usage:
        think command [arguments]

The commands are:
        new        create, compile and run (monitor changes) a new thinkgo project
        run        compile and run (monitor changes) an any existing go project

think new appname [apptpl]
        appname    specifies the path of the new thinkgo project
        apptpl     optionally, specifies the thinkgo project template type

think run [appname]
        appname    optionally, specifies the path of the new project
`

func help() {
	fmt.Println(helpInfo)
}

func newappHelp() {
	fmt.Println(helpInfo)
}

func runappHelp() {
	fmt.Println(helpInfo)
}

func initVar(args []string) {
	var dir string
	if len(args) > 0 {
		dir, appname = filepath.Split(args[0])
		if dir != "" {
			crupath = filepath.Join(dir, appname)
		} else {
			crupath = filepath.Join(crupath, appname)
		}
	} else {
		_, appname = filepath.Split(crupath)
	}
	var err error
	crupath = strings.TrimSpace(crupath)
	crupath, err = filepath.Abs(crupath)
	if err != nil {
		thinkgo.Fatalf("[think] Create project fail: %s", err)
	}
	crupath = strings.Replace(crupath, `\`, `/`, -1)
	crupath = strings.TrimRight(crupath, "/") + "/"
}
