[![Build Status](https://travis-ci.org/utamaro/neje-ui.svg?branch=master)](https://travis-ci.org/utamaro/neje-ui)
[![GoDoc](https://godoc.org/github.com/utamaro/neje-ui?status.svg)](https://godoc.org/github.com/utamaro/neje-ui)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/utamaro/neje-ui/master/LICENSE)


# neje-ui

**Don't Embed, Just Execute** Chrome browser for a UI in Go.

For now it's just a PoC (proof of concept).  Don't believe me so much. :)

I believe this works on Linux, Win and OS X, but I have only tested on Linux. 

![](http://imgur.com/2TSlOIp.gif)


## Overview

This library is a UI alternative for Go which uses the Chrome browser (or an alternate browser) that is already installed. The application communicates with the browser via a JSON-RPC websocket using [gopherjs](https://github.com/gopherjs/gopherjs). That means you can call funcs in the browser from the server (and vice versa) in [golang RPC-style](https://golang.org/pkg/net/rpc/) 
without worrying about the websocket and JavaScript.

You can write the server-side *and* client side program in Go.

## Requirements

This requires

* git
* go 1.7 (for gopherjs)
* gopherjs
```
go get -u github.com/gopherjs/gopherjs
```

## Installation

    $ go get -u github.com/utamaro/neje-ui


## Example
(This example omits error handlings for simplicity.)

## browser side

[ex.go](https://github.com/utamaro/wsrpc/blob/master/example/browser/ex.go)

```go

//GUI is struct to be called from remote by rpc.
type GUI struct{}

//Write writes a response from the server.
func (g *GUI) Write(msg *string, reply *string) error {
	//show welcome message:
	jquery.NewJQuery("#from_server").SetText(msg)
	return nil
}

func main() {
	b,_ := browser.New(new(GUI))
	jquery.NewJQuery("button").On(jquery.CLICK, func(e jquery.Event) {
		go func() {
			m := jquery.NewJQuery("#to_server").Val()
			response := ""
			b.Call("Msg.Message", &m, &response)
			//show welcome message:
			jquery.NewJQuery("#response").SetText(response)
		}()
	})

}


```

Then compile it with GopherJS to create ex.js:

```
go get  
gopherjs build ex.go
```

## webserver side

[ex.go](https://github.com/utamaro/wsrpc/blob/master/example/webserver/ex.go)

```go

//Msg  is struct to bel called from remote by rpc.
type Msg struct{}

//Message writes a message to the browser.
func (t *Msg) Message(m *string, response *string) error {
	*response = "OK, I heard that you said\"" + *m + "\""
	return nil
}

func main() {
	ws,_ := webserver.New("", "ex.html", new(Msg))

	for {
		select {
		case <-ws.Finished:
			log.Println("Browser was closed. Exiting...")
			return
		case <-time.After(10 * time.Second):
			msg := "Now " + time.Now().String() + " at server!"
			reply := ""
			ws.Call("GUI.Write", &msg, &reply)
		}
	}
}

```

Then copy ex.html and ex.js to the webserver directory,
```
go run ex.go
```

Then your Chrome browser (or something else if Chrome is not installed) will open automatically and
display the demo.

## What is It Doing?

1. Start a webserver including websocket on a free port at localhost.
1. Search Chrome browser. 
	1. If windows, read registry to get the path to Chrome. 
	2. If OS X, just run "open -a google chrome ".
    3. if Linux, run "google-chrome", "chrome", or something else.
1. If Chrome is found, run chrome with the options "--disable-extension --app=<url>"
for me1.. If Chrome is not found, 
	1. If Windows, just run "start <url>". 
	2. If OS X, just run "open <url>  ".
    3. if Linux, just run "xdg-open <url>"
1. Communicate between the webserver and browser using websockets.

## Why not embed Chrome lib?

1. Chrome lib is very big (about 100MB?) for a single app.
2. Chrome lib APIs are always changing.
3. One doesn't want to loose the eco system (easy to cross compile etc) of Go.
4. Chrome lib is difficult to understand. :(
5. Chrome browser has convenient options for an application (--app etc).

### Pros
 
 * Can make Go apps that can be cross compiled easily with a small size.

### Cons

* Can't control browser precisely--must control them via JavaScript, manually. (window size, menu etc.)
* Behaviour may be different for each platform if Chrome is not found.


# Contribution
Improvements to the codebase and pull requests are encouraged.


