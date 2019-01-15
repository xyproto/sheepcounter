package sheepcounter

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const addr = ":1911"

func TestInterface(t *testing.T) {
	sc := New(httptest.NewRecorder())
	var _ http.ResponseWriter = sc
	var _ http.Hijacker = sc
}

func TestCounting(t *testing.T) {
	counted := make(chan int64)
	done := make(chan bool)
	srv := &http.Server{
		Addr: addr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sc := NewSheepCounter(w)
			fmt.Fprintln(sc, "Hi") // 3 bytes
			counted <- sc.Counter()
		}),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		<-done
		srv.Shutdown(context.Background())
	}()
	go srv.ListenAndServe()
	go http.Get("http://localhost" + addr + "/")
	if <-counted != 3 {
		t.Error("Wrong byte count!")
	}
	done <- true
}
