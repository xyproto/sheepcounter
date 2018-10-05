package sheepcounter

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestInterface(t *testing.T) {
	sc := NewSheepCounter(httptest.NewRecorder())
	var _ http.ResponseWriter = sc
}

func TestCounting(t *testing.T) {
	counted := make(chan int64)
	done := make(chan bool)
	srv := &http.Server{
		Addr: ":8080",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			sc := NewSheepCounter(w)
			fmt.Fprintln(sc, "Hi") // 3 bytes
			counted <- sc.Counter()
			done <- true
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
	go http.Get("http://localhost:8080/")
	if <-counted != 3 {
		t.Error("Wrong byte count!")
	}
}