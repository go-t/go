package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/mgutz/ansi"
)

var (
	flagPretty = flag.Bool("pretty", false, "Pretty print")
	flagDebug  = flag.Bool("debug", false, "Print message type")

	NEWLINE = []byte{'\n'}
)

func writeln(buf []byte) {
	n := len(buf)
	for n > 0 && (buf[n-1] == '\n' || buf[n-1] == '\r') {
		n -= 1
	}
	os.Stdout.Write(buf[:n])
	os.Stdout.Write(NEWLINE)
}

func main() {
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Println("Usage: wscat [options] <url>\n")
		fmt.Println("options:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	u, err := url.Parse(flag.Arg(0))
	if err != nil {
		panic(err)
	}
	h := http.Header{"Origin": {"http://" + u.Host}}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), h)

	if err != nil {
		panic(err)
	}

	go func() {
		for {
			t, buf, err := conn.ReadMessage()
			if *flagDebug {
				log.Println("type:", t, len(buf), err)
				continue
			}
			if *flagPretty {
				color := ansi.ColorCode("green")
				reset := ansi.ColorCode("reset")
				fmt.Fprintf(os.Stdout, "%s%s%s\n", color, strings.Repeat("<", 10), reset)

				var v interface{}
				if err := json.Unmarshal(buf, &v); err != nil {
					writeln(buf)
				} else {
					buf, err = json.MarshalIndent(v, "", "  ")
					if err != nil {
						panic(err)
					}
					writeln(buf)
				}
			} else {
				writeln(buf)
			}
		}
	}()

	b := bufio.NewReader(os.Stdin)
	for {
		line, err := b.ReadString('\n')
		if err != nil {
			log.Println(err)
			break
		}
		if err := conn.WriteMessage(websocket.TextMessage, []byte(line)); err != nil {
			break
		}
	}

	conn.Close()
}
