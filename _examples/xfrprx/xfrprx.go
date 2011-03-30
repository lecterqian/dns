package main

// Xfrprx is a proxy that intercepts notify messages
// and then performs a ixfr/axfr to get the new 
// zone contents. 
// This zone is then checked cryptographically is
// everything is correct.
// When the message is deemed correct a remote 
// server is sent a notify to retrieve the ixfr/axfr.
// If a new DNSKEY record is seen for the apex and
// it validates it writes this record to disk and
// this new key will be used in future validations.

import (
	"os"
	"os/signal"
	"fmt"
	"dns"
)

// Static amount of RRs...
type zone struct {
	name string
	rrs  [10000]dns.RR
	size int
}

var Zone zone

func handle(d *dns.Conn, i *dns.Msg) {
/* send response here, how ??? */
	if i.MsgHdr.Response == true {
		return
	}
	handleNotify(d, i)
	handleXfr(d, i)
}

func qhandle(d *dns.Conn, i *dns.Msg) {
        // We should send i to d.RemoteAddr
        // simpleQuery here

        // what do we do with the reply
        ///handle HERE!!?? Need globals or stuff in d...

        //        in/out channel must be accessible
}

func listen(addr string, e chan os.Error, tcp string) {
	switch tcp {
	case "tcp":
		err := dns.ListenAndServeTCP(addr, handle)
		e <- err
	case "udp":
		err := dns.ListenAndServeUDP(addr, handle)
		e <- err
	}
	return
}

func query(e chan os.Error, tcp string) {
        switch tcp {
        case "tcp":
                err := dns.QueryAndServeTCP(qhandle)
                e <- err
        case "udp":
                err := dns.QueryAndServeUDP(qhandle)
                e <- err
        }
        return
}

func main() {
	err := make(chan os.Error)

	// Outgoing queries
        dns.QueryInitChannels()
	go query(err, "tcp")
        go query(err, "udp")

	// Incoming queries
	go listen("127.0.0.1:8053", err, "tcp")
	go listen("[::1]:8053", err, "tcp")
	go listen("127.0.0.1:8053", err, "udp")
	go listen("[::1]:8053", err, "udp")

forever:
	for {
		select {
		case e := <-err:
			fmt.Printf("Error received, stopping: %s\n", e.String())
			break forever
		case <-signal.Incoming:
			fmt.Printf("Signal received, stopping")
			break forever
                case q := <-dns.QueryReply:
                        var _ = q
                        /* ... */
		}
	}
	close(err)

}
