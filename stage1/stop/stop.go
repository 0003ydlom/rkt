// Copyright 2016 The rkt Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//+build linux

package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

var (
	force bool
)

func init() {
	flag.BoolVar(&force, "force", false, "Forced stopping")
}

func readIntFromFile(path string) (i int, err error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return
	}
	_, err = fmt.Sscanf(string(b), "%d", &i)
	return
}

// Send signal to process and all its child processes
func recursiveKill(signal syscall.Signal, pid int) int {

	var outbuf bytes.Buffer
	var rc int

	cpidCmd := exec.Command("pgrep", "-P", strconv.Itoa(pid))
	cpidCmd.Stdout = &outbuf
	cpidCmd.Run()
	cpidsStr := strings.Split(outbuf.String(), "\n")

	for i := 0; i < len(cpidsStr)-1; i++ {
		cpid, _ := strconv.Atoi(cpidsStr[i])
		rc = recursiveKill(signal, cpid)
		if rc != 0 {
			return rc
		}
	}

	rc = stop(signal, pid)
	if rc != 0 {
		return rc
	}

	return 0

}

func stop(signal syscall.Signal, pid int) int {

	if err := syscall.Kill(pid, signal); err != nil {
		fmt.Fprintf(os.Stderr, "error sending %v to %d: %v\n", signal, pid, err)
		return 1
	}

	return 0
}

func main() {
	flag.Parse()

	pid, err := readIntFromFile("ppid")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading pid: %v\n", err)
		return
	}

	if force {
		recursiveKill(syscall.SIGKILL, pid)
	} else {
		os.Exit(stop(syscall.SIGTERM, pid))
	}

	return

}
