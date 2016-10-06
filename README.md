[![Build Status](https://travis-ci.org/utamaro/neje-ui.svg?branch=master)](https://travis-ci.org/utamaro/neje-ui)
[![GitHub license](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/utamaro/neje-ui/master/LICENSE)


# neje-ui

Do **Not Embed, Just Execute** the browser for an UI in Go.

### App option of chrome mode 
![](http://imgur.com/2TSlOIp.gif)

* You need to install Chrome browser.
* It looks like a real application.

### Default browser mode (ex. firefox) 
![](http://i.imgur.com/CbDrwWr.gif)

* You can install other browsers.
* Its appearance is "browser" itself. 
* You cannot control window size and window location by Javascript due to the  browser restrictions.

## Overview

This library is a UI alternative for Go which uses the Chrome browser (or a default browser) that is already installed. 
The application communicates with the browser via a JSON-RPC websocket using [GopherJS](https://github.com/gopherjs/gopherjs).
 That means you can call funcs in the browser from the server (and vice versa) in [golang RPC-style](https://golang.org/pkg/net/rpc/) 
without worrying about the websocket and JavaScript.

You can write the server-side *and* client side program in Go.

## Requirements

This requires

* git
* go 1.7 (for GopherJS)
* web browser
	* Newest Chrome browser is recommended.
	* firefox
	* This doens't work on Internet Explorer and Microsoft Edge (because GopherJS uses ECMAScript6?) 
* GopherJS

```
go get -u github.com/gopherjs/gopherjs
```

## Platforms

* Linux
* OSX
* Windows

## Installation

    $ go get -u github.com/utamaro/neje-ui

## Document

[frontend](https://godoc.org/github.com/utamaro/neje-ui/frontend)
[backend](https://godoc.org/github.com/utamaro/neje-ui/backend)


## Example
(This example omits error handlings for simplicity.)

## browser side

[ex.html](https://github.com/utamaro/neje-ui/blob/master/example/browser/ex.html)
[ex.go](https://github.com/utamaro/neje-ui/blob/master/example/browser/ex.go)

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
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	b,_ := frontend.New(new(GUI))
	//defer b.Close()
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

[ex.go](https://github.com/utamaro/neje-ui/blob/master/example/webserver/ex.go)

```go

//Msg  is struct to be called from remote by rpc.
type Msg struct{}

//Message writes a message to the browser.
func (t *Msg) Message(m *string, response *string) error {
	*response = "OK, I heard that you said\"" + *m + "\""
	return nil
}

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Llongfile)
	//for default browser mode
	ws, _ := backend.New(backend.Default, "ex.html", new(Msg))
	//for app option mode
//	ws, _ := backend.New(backend.AppOptionChrome, "ex.html", new(Msg))
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
			ws.Call("GUI.Write", &msg, &reply)
		}
	}
}

```

Then copy ex.html and ex.js to the webserver directory,
```
go run ex.go
```

Then your default browser will open automatically and display the demo.

## What is It Doing?

1. Start a webserver including websocket on a free port at localhost.
1. Search Chrome browser if app option mode.  
	1. If windows, read registry to get the path to Chrome. 
	2. If OS X, run "open -a google chrome ".
    3. if Linux, run "google-chrome", "chrome", or something else.
    1. If Chrome is found, run chrome with the options "--disable-extension --app=<url>"
1. If default browser mode, 
	1. If Windows, run "start <url>". 
	2. If OS X, run "open <url>  ".
    3. if Linux, run "xdg-open <url>"
1. The werbserver and browser starts communication browser using websockets.

## Why not embed Chrome lib?

1. Chrome lib is too big (about 100MB?) for a single app. 
Go is used to make a single app, not used in framework like electron. 
2. Chrome lib APIs are always changing.
3. Nobody wants to loose the eco system (easy to cross compile etc) of Go.
4. Chrome lib is too difficult to understand :(
5. Chrome browser has convenient options for an application (--app etc).

### Pros
 
 * You can make Go apps that can be cross compiled easily with a small size.

### Cons

* You can't control browser precisely--must control them via JavaScript, manually. (window size, menu etc.)
* Behaviour may be different for each platforms when you use default browser.


# Contribution
Improvements to the codebase and pull requests are encouraged.
Also any ideas including attractive demos are welcomed.

