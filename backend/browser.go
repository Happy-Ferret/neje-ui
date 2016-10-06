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

var optionAppStyleChrome = "--disable-extension  --new-window --app=%s"

func browsers(t int) ([]string, string) {
	var cmds []string
	opt := ""
	switch t {
	case Default:
		cmds, opt = defaultPaths()
		if exe := os.Getenv("BROWSER"); exe != "" {
			cmds = append([]string{exe}, cmds...)
		}
		if opt == "" {
			opt = "%s"
		} else {
			opt += " %s"
		}
	case AppOptionChrome:
		cmds, opt = chromePaths()
		if opt == "" {
			opt = optionAppStyleChrome
		} else {
			opt += " " + optionAppStyleChrome
		}
	}
	return cmds, opt
}

func tryBrowser(t int, p string) error {
	cmds, opt := browsers(t)
	o := fmt.Sprintf(opt, p)
	args := strings.Split(o, " ")
	filteredArgs := args[:0]
	for _, arg := range args {
		filteredArgs = append(filteredArgs, strings.Replace(arg, "^", " ", -1))
	}
	for _, c := range cmds {
		// Separate command and arguments for exec.Command.
		if len(args) == 0 {
			continue
		}
		log.Println("executing", c, o)
		viewer := exec.Command(c, filteredArgs...)
		viewer.Stderr = os.Stderr
		if err := viewer.Start(); err != nil {
			log.Println(err)
			continue
		}
		return nil
	}
	return errors.New("no browsers are found")
}
