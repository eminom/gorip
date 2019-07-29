package main

import (
	"github.com/elazarl/goproxy"
	"github.com/elazarl/goproxy/ext/auth"

	"flag"
	"log"
	"net"
	"net/http"
	"net/url"

	"./blocker"
	"./conn"
	"./xau"
)

var (
	verbose      = flag.Bool("v", false, "should every proxy request be logged to stdout")
	addr         = flag.String("addr", ":8080", "proxy listen address")
	showRipAhora = flag.Bool("showrip", false, "test rip-list")
	nextHop      = flag.String("hop", "", "next hop for prxy")

	bLog bool
)

func init() {
	flag.BoolVar(&bLog, "log", false, "log for debugging")
	bLog = true // remove this line when debug is done
	flag.Parse()
}

func main() {

	*showRipAhora = true
	if *showRipAhora {
		blocker.DefAl.Dump()
	}
	proxy := goproxy.NewProxyHttpServer()
	if xau.HaySecret {
		log.Print("auth is on")
		proxy.OnRequest().Do(auth.Basic("despacito", xau.CheckUserPasswd))
	} else {
		log.Print("no auth at proxy")
	}

	// The old way
	if false {
		proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			addr, e := net.ResolveTCPAddr("", req.RemoteAddr)
			//if e != nil || !addr.IP.Equal(misIP) {
			if e != nil || !blocker.DefAl.IsAllowed(addr.IP) {
				log.Printf("Not allowed for <%v>", req.RemoteAddr)
				return req, goproxy.NewResponse(req, goproxy.ContentTypeText, http.StatusForbidden, "not today")
			}
			return req, nil
		})
	}

	proxy.OnRequest().HandleConnect(goproxy.FuncHttpsHandler(
		func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
			return goproxy.OkConnect, host
		},
	))

	proxy.Verbose = *verbose

	if len(*nextHop) > 0 {
		pUrl, err := url.Parse(*nextHop)
		if nil == err {
			log.Printf("using proxy hop: %v", pUrl)
			proxy.Tr.Proxy = http.ProxyURL(pUrl)
		} else {
			log.Printf("url.Pase error: %v", err)
		}
	} else {
		log.Printf("proxy hop is disabled")
	}

	var (
		ln  net.Listener
		err error
	)

	if !blocker.DefAl.IsEmpty() {
		log.Printf("ip blocking-is on")
		ln, err = conn.NewTcpListener(*addr, func(addr string) bool {
			tcpAddr, err := net.ResolveTCPAddr("", addr)
			if nil != err {
				return false // unknown address(not tcp4 addr)
			}
			//log.Printf("filter tcp-addr: <%v>:<%v>", tcpAddr.IP, tcpAddr.Port)
			return blocker.DefAl.IsAllowed(tcpAddr.IP)
		})
	} else {
		log.Printf("no ip-blocking")
		ln, err = net.Listen("tcp", *addr)
	}

	if err != nil {
		log.Fatal(err)
	}
	svr := http.Server{
		Handler: proxy,
	}
	log.Printf("entering serve")
	log.Fatal(svr.Serve(ln))
	// log.Fatal(http.ListenAndServe(*addr, proxy))
}
