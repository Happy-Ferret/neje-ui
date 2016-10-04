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

package main

import (
	"flag"
	"log"
	"time"

	"github.com/utamaro/neje-ui/backend"
)

//Msg  is struct to be called from remote by rpc.
type Msg struct{}

//Message writes a message to the browser.
func (t *Msg) Message(m *string, response *string) error {
	*response = "OK, I heard that you said\"" + *m + "\""
	return nil
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	browser := 0
	flag.IntVar(&browser, "type", 0, "0:deault browser 1:app style chrome")
	flag.Parse()
	ws, err := backend.New(browser, "ex.html", new(Msg))
	if err != nil {
		log.Fatal(err)
	}
	i := 0
	for {
		select {
		case <-ws.Finished:
			log.Println("browser was closed. Exiting...")
			return
		case <-time.After(10 * time.Second):
			i++
			log.Println("writing", i, "to browser")
			msg := "Now " + time.Now().String() + " at server!"
			reply := ""
			if err := ws.Call("GUI.Write", &msg, &reply); err != nil {
				log.Fatal(err)
			}
		}
	}
}
