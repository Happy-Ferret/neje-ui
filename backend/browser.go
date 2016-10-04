/*
 * Copyright (c) 2016, Shinya Yagyu
 * All rights reserved.
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice,
 *    this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 * 3. Neither the name of the copyright holder nor the names of its
 *    contributors may be used to endorse or promote products derived from this
 *    software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE
 * LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 */

package backend

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

const (
	//Default represents default browser.
	Default = iota
	//AppOptionChrome  chrome browser with options.
	AppOptionChrome
)

var optionAppStyleChrome = " --disable-extension  --new-window --app=%s"

func browsers(t int) []string {
	var paths []string
	switch t {
	case Default:
		paths = defaultPaths()
		if exe := os.Getenv("BROWSER"); exe != "" {
			paths = append([]string{exe}, paths...)
		}
		for i := range paths {
			paths[i] += " %s"
		}
	case AppOptionChrome:
		paths = chromePaths()
		for i := range paths {
			paths[i] += optionAppStyleChrome
		}
	}
	return paths
}

func tryBrowser(t int, p string) error {
	cmds := browsers(t)
	for _, v := range cmds {
		v = fmt.Sprintf(v, p)
		// Separate command and arguments for exec.Command.
		args := strings.Split(v, " ")
		if len(args) == 0 {
			continue
		}
		log.Println("executing", v)
		viewer := exec.Command(args[0], args[1:]...)
		//viewer.Stderr = os.Stderr
		if err := viewer.Start(); err != nil {
			log.Println(err)
			continue
		}
		return nil
	}
	return errors.New("no browsers are found")
}
