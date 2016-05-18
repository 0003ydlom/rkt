// Copyright 2016 The rkt Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" IS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//+build linux

package stage0

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/coreos/rkt/common"

	"github.com/appc/spec/schema/types"
)

func StopPod(dir string, force bool, uuid *types.UUID) error {
	s1v, err := getStage1InterfaceVersion(dir)
	if err != nil {
		return fmt.Errorf("8error determining stage1 interface version: %v", err)
	}

	if !interfaceVersionSupportsStop(s1v) {
		return fmt.Errorf("9stop entrypoint not supported by stage1")
	}

	s1rootfs := common.Stage1RootfsPath(dir)

	if err := os.Chdir(dir); err != nil {
		return fmt.Errorf("Afailed changing to dir: %v", err)
	}

	ep, err := getStage1Entrypoint(dir, stopEntrypoint)
	if err != nil {
		return fmt.Errorf("Brkt stop not implemented for pod's stage1: %v", err)
	}
	args := []string{filepath.Join(s1rootfs, ep)}
	debug("Execing %s", ep)

	if force {
		args = append(args, "--force")
	}

	args = append(args, uuid.String())

	// STDOUT I STDERR ARE REMOVED
	c := exec.Cmd{
		Path: args[0],
		Args: args,
	}

	return c.Run()
}
