package main

// gothello is a little http-server which just sits somewhere
// and answers ("hello") http-requests with quotes from shakespeare's
// play 'otello'.
//
// it's only purpose is to serve as little helper tool to check ssh or
// vpn-tunnels, general network issues etc.
//
// author: mathias gumz <mg@2hoch5.com>
// license: do whatever you like.

import (
	"flag"
	"fmt"
	"io"
	"math/rand"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

// http://en.wikiquote.org/wiki/Othello
var QUOTES = []string{
	"Keep up your bright swords, for the dew will rust them.\n",
	"Then put up your pipes in your bag, for I'll away. Go, vanish into air, away!\n",
	"O! now, for ever Farewell the tranquil mind; farewell content!\n",
	"Put out the light, and then put out the light.\n",
	"I would not put a thief in my mouth to steal my brains.\n",
	"Rude am I in my speech, And little blessed with the soft phrase of peace.\n",
}

func gothello(w http.ResponseWriter, _ *http.Request) {
	io.WriteString(w, QUOTES[rand.Intn(len(QUOTES))])
}

func main() {
	binds := []string{":8080"}
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "answer http-requests with quotes from othello")
		fmt.Fprintf(os.Stderr, "usage:\n    %s bind1 [bind2] [bind3] [...]\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
	flag.Parse()
	if flag.NArg() > 0 {
		binds = flag.Args()
	}

	http.HandleFunc("/", gothello)
	servers := sync.WaitGroup{}
	for _, addr := range binds {
		conn, err := net.Listen("tcp", addr)
		if err != nil {
			fmt.Fprint(os.Stderr, err)
			continue
		}
		servers.Add(1)
		go http.Serve(conn, nil)
	}

	servers.Wait()
}
