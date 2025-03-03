package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/panjf2000/gnet"
	"github.com/panjf2000/gnet/ringbuffer"
)

func main() {
	var port int
	var loops int
	var udp bool
	var trace bool
	var reuseport bool

	flag.IntVar(&port, "port", 5000, "server port")
	flag.BoolVar(&udp, "udp", false, "listen on udp")
	flag.BoolVar(&reuseport, "reuseport", false, "reuseport (SO_REUSEPORT)")
	flag.BoolVar(&trace, "trace", false, "print packets to console")
	flag.IntVar(&loops, "loops", -1, "num loops")
	flag.Parse()

	var events gnet.Events
	events.NumLoops = loops
	events.OnInitComplete = func(srv gnet.Server) (action gnet.Action) {
		log.Printf("echo server started on port %d (loops: %d)", port, srv.NumLoops)
		if reuseport {
			log.Printf("reuseport")
		}
		return
	}
	events.React = func(c gnet.Conn, inBuf *ringbuffer.RingBuffer) (out []byte, action gnet.Action) {
		top, tail := inBuf.PreReadAll()
		out = append(top, tail...)
		inBuf.Reset()

		if trace {
			log.Printf("%s", strings.TrimSpace(string(top)+string(tail)))
		}
		return
	}
	scheme := "tcp"
	if udp {
		scheme = "udp"
	}
	log.Fatal(gnet.Serve(events, fmt.Sprintf("%s://:%d", scheme, port)))
}
