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
	"log"
	"net"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"

	"golang.org/x/net/websocket"
)

//Ping  is struct to be pinged.
type Ping struct{}

//Ping responses to a ping..
func (p *Ping) Ping(m *string, response *string) error {
	*response = "OK"
	return nil
}

//Backend represents web server side.
type Backend struct {
	client      *rpc.Client
	closeClient chan struct{}
	//Finished is a chan which signals browsser was closed.
	Finished chan struct{}
}

//ping pings the server.
func (w *Backend) ping() {
	go func() {
		count := 0
		m := ""
		r := ""
		for range time.Tick(time.Second) {
			if w.client == nil {
				continue
			}
			if errr := w.Call("Ping.Ping", &m, &r); errr != nil {
				count++
			} else {
				count = 0
			}
			if count > 2 {
				log.Println("browser seems to be closed.")
				w.Finished <- struct{}{}
				return
			}
		}
	}()
}

//New registers strs to RPC as funcs, starts web server, open firstpage by browser,
// and  returns Backend obj.
//t selects chrome or the default browser.
func New(t int, firstPage string, strs ...interface{}) (*Backend, error) {
	if err := rpc.Register(new(Ping)); err != nil {
		return nil, err
	}
	for _, str := range strs {
		if err := rpc.Register(str); err != nil {
			return nil, err
		}
	}
	w := &Backend{
		closeClient: make(chan struct{}),
	}
	addr := w.start()
	err := tryBrowser(t, "http://"+addr.String()+"/"+firstPage)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	w.ping()
	return w, nil
}

//Close closes client RPC connection.
func (w *Backend) Close() {
	w.closeClient <- struct{}{}
}

//Call calls calls RPC.
func (w *Backend) Call(m string, args interface{}, reply interface{}) error {
	if w.client == nil {
		return errors.New("not connected to server yet")
	}
	return w.client.Call(m, args, reply)
}

//regHandlers registers handlers to http return servemux.
func (w *Backend) regHandlers() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/ws-server",
		func(rw http.ResponseWriter, req *http.Request) {
			log.Println("connected to ws-server")
			s := websocket.Server{
				Handler: websocket.Handler(func(ws *websocket.Conn) {
					jsonrpc.ServeConn(ws)
					log.Println("ws-server was disconnected")
				}),
			}
			s.ServeHTTP(rw, req)
		})
	mux.HandleFunc("/ws-client",
		func(rw http.ResponseWriter, req *http.Request) {
			log.Println("connected to ws-client")
			s := websocket.Server{
				Handler: websocket.Handler(func(ws *websocket.Conn) {
					w.client = jsonrpc.NewClient(ws)
					<-w.closeClient //wait for Close()
				}),
			}
			s.ServeHTTP(rw, req)
		})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, r.URL.Path[1:])
	})
	return mux
}

//start starts webserver at localhost:0(i.e. free port) and returns listen address.
func (w *Backend) start() net.Addr {
	mux := w.regHandlers()
	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		if err := http.Serve(l, mux); err != nil {
			log.Fatal(err)
		}
	}()
	return l.Addr()
}
