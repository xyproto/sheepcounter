package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/xyproto/sheepcounter"
)

const (
	title = "SheepCounter"
	style = `body { margin: 4em; background: wheat; color: black; font-family: terminus, "courier new", courier; font-size: 1.1em; } a:link { color: #403020; } a:visited { color: #403020; } a:hover { color: #605040; } a:active { color: #605040; } #counter { color: red; }`
	page  = "<!doctype html><html><head><style>%s</style><title>%s</title><body>%s</body></html>"
)

var totalBytesWritten uint64

func helloHandler(w http.ResponseWriter, r *http.Request) {
	sc := sheepcounter.New(w)
	body := `<p>Here are the <a href="/counter">counted bytes</a>.</p>`
	fmt.Fprintf(sc, page, style, title, body)
	written, err := sc.UCounter2()
	if err != nil {
		// Log an error and return
		log.Printf("error: %s\n", err)
		return
	}
	atomic.AddUint64(&totalBytesWritten, written)
	log.Printf("counted %d bytes\n", written)
}

func counterHandler(w http.ResponseWriter, r *http.Request) {
	sc := sheepcounter.New(w)
	body := fmt.Sprintf(`<p>Total bytes sent from the server (without counting this response): <span id="counter">%d</span></p><p><a href="/">Back</a></p>`, atomic.LoadUint64(&totalBytesWritten))
	fmt.Fprintf(sc, page, style, title, body)
	written, err := sc.UCounter2()
	if err != nil {
		// Log an error and return
		log.Printf("error: %s\n", err)
		return
	}
	atomic.AddUint64(&totalBytesWritten, written)
	log.Printf("counted %d bytes\n", written)
}

func main() {
	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/counter", counterHandler)

	httpAddr := os.Getenv("HTTP_ADDR")
	if httpAddr == "" {
		httpAddr = ":4000"
	}

	log.Println("Serving on " + httpAddr)
	log.Fatal(http.ListenAndServe(httpAddr, nil))
}
