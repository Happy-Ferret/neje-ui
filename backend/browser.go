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

var optionAppStyleChrome = []string{"--disable-extension", "--new-window", "--app=%s"}

//replaceURL replaces %s in orig to url.
func replaceURL(orig []string, url string) {
	for i, o := range orig {
		if strings.Contains(o, "%s") {
			orig[i] = fmt.Sprintf(o, url)
		}
	}
}

//browsers returns the browser paths.
//if t==Default this returns one of default browser.
//if t=AppOptionChrome this returns one of chrome browser
func browsers(t int) ([]string, []string) {
	var cmds, opt []string
	switch t {
	case Default:
		cmds, opt = defaultPaths()
		if exe := os.Getenv("BROWSER"); exe != "" {
			cmds = append([]string{exe}, cmds...)
		}
		opt = append(opt, "%s")
	case AppOptionChrome:
		cmds, opt = chromePaths()
		opt = append(opt, optionAppStyleChrome...)
	}
	return cmds, opt
}

//ErrNoBrowser is an error that no browsers are found.
var ErrNoBrowser = errors.New("no browsers are found")

//tryBrowser tries to start browsers and opens URL p.
//if t==Default this tries to start default browser.
//if t=AppOptionChrome this tries to chrome browser
func tryBrowser(t int, p string) error {
	cmds, opt := browsers(t)
	replaceURL(opt, p)
	for _, c := range cmds {
		log.Println("executing", c, opt)
		viewer := exec.Command(c, opt...)
		//viewer.Stderr = os.Stderr
		if err := viewer.Start(); err != nil {
			log.Println(err)
			continue
		}
		return nil
	}
	return ErrNoBrowser
}
