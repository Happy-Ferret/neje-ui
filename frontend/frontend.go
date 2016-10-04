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

package frontend

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"time"

	"github.com/gopherjs/gopherjs/js"
	"github.com/gopherjs/websocket"
)

//Ping  is struct to be pinged.
type Ping struct{}

//Ping responses to a ping..
func (p *Ping) Ping(m *string, response *string) error {
	*response = "OK"
	return nil
}

//Frontend represents RPC of browser side
type Frontend struct {
	client *rpc.Client
	s      net.Conn
	c      net.Conn
}

//ping pings the server.
func (b *Frontend) ping() {
	go func() {
		count := 0
		m := ""
		r := ""
		for range time.Tick(time.Second) {
			if errr := b.Call("Ping.Ping", &m, &r); errr != nil {
				log.Println(errr)
				count++
			} else {
				count = 0
			}
			if count > 2 {
				log.Println("disconnected from server, reconnecting...")
				b.Close()
				if err := b.connect(); err != nil {
					log.Println(err)
					js.Global.Get("window").Call("alert", "cannot connect to server.")
				}
				count = 0
			}
		}
	}()
}

//connect connects two connections to the server.
func (b *Frontend) connect() error {
	var err error
	port := js.Global.Get("window").Get("location").Get("port").String()
	log.Println(port)
	b.s, err = websocket.Dial("ws://localhost:" + port + "/ws-client") // Blocks until connection is established
	if err != nil {
		log.Println("ws://localhost:"+port+"/ws-client", err)
		return err
	}
	log.Println("connected to ws-client")
	go jsonrpc.ServeConn(b.s)

	b.c, err = websocket.Dial("ws://localhost:" + port + "/ws-server") // Blocks until connection is established
	if err != nil {
		log.Println("ws://localhost:"+port+"/ws-server", err)
		return err
	}
	log.Println("connected to ws-server")
	b.client = jsonrpc.NewClient(b.c)
	return nil
}

//New registers strs to RPC as funcs, connects websocket and returns Frontend obj.
func New(strs ...interface{}) (*Frontend, error) {
	b := &Frontend{}
	if err := rpc.Register(new(Ping)); err != nil {
		return nil, err
	}
	for _, str := range strs {
		if err := rpc.Register(str); err != nil {
			return nil, err
		}
	}
	b.connect()
	b.ping()
	return b, nil
}

//Call Calls calls RPC.
func (b *Frontend) Call(m string, args interface{}, reply interface{}) error {
	return b.client.Call(m, args, reply)
}

//Close closes RPC client and server connections.
func (b *Frontend) Close() {
	if b.c != nil {
		if err := b.c.Close(); err != nil {
			log.Println(err)
		}
	}
	if b.s != nil {
		if err := b.s.Close(); err != nil {
			log.Println(err)
		}
	}
}
