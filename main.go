package main

import (
	"github.com/elazarl/goproxy"
	"github.com/elazarl/goproxy/ext/auth"

	"flag"
	"log"
	"net"
	"net/http"
	"net/url"

	"gorip/blocker"
	"gorip/conn"
	"gorip/xau"
)

func main() {

	var (
		verbose      = flag.Bool("v", false, "should every proxy request be logged to stdout")
		addr         = flag.String("addr", ":8080", "proxy listen address")
		showRipAhora = flag.Bool("showrip", false, "test rip-list")
		nextHop      = flag.String("hop", "", "next hop for prxy")
	)
	flag.Parse()

	if *showRipAhora {
		blocker.DefAl.Dump()
		return
	}
	proxy := goproxy.NewProxyHttpServer()
	proxy.OnRequest().Do(auth.Basic("despacito", xau.CheckUserPasswd))

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
			if *verbose {
				log.Printf("using proxy hop: %v", pUrl)
			}
			proxy.Tr.Proxy = http.ProxyURL(pUrl)
		}
	}

	ln, err := conn.NewTcpListener(*addr, func(addr string) bool {
		tcpAddr, err := net.ResolveTCPAddr("", addr)
		if nil != err {
			return false // unknown address(not tcp4 addr)
		}
		//log.Printf("filter tcp-addr: <%v>:<%v>", tcpAddr.IP, tcpAddr.Port)
		return blocker.DefAl.IsAllowed(tcpAddr.IP)
	})
	if err != nil {
		log.Fatal(err)
	}
	svr := http.Server{
		Handler: proxy,
	}
	log.Fatal(svr.Serve(ln))
	// log.Fatal(http.ListenAndServe(*addr, proxy))
}
