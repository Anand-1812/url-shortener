package middleware

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestRateLimitEnforcesBurstPerClient(t *testing.T) {
	var handled atomic.Int32
	limited := RateLimit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handled.Add(1)
		w.WriteHeader(http.StatusNoContent)
	}), 100, 2)

	for i := 0; i < 2; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "198.51.100.10:1234"
		rec := httptest.NewRecorder()
		limited.ServeHTTP(rec, req)
		if rec.Code != http.StatusNoContent {
			t.Fatalf("request %d status = %d, want %d", i+1, rec.Code, http.StatusNoContent)
		}
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "198.51.100.10:1234"
	rec := httptest.NewRecorder()
	limited.ServeHTTP(rec, req)
	if rec.Code != http.StatusTooManyRequests {
		t.Fatalf("third request status = %d, want %d", rec.Code, http.StatusTooManyRequests)
	}
	if handled.Load() != 2 {
		t.Fatalf("handler ran %d times, want 2", handled.Load())
	}
}

func TestRateLimitDoesNotHoldLockWhileHandlerRuns(t *testing.T) {
	firstStarted := make(chan struct{})
	releaseFirst := make(chan struct{})
	secondStarted := make(chan struct{})
	var calls atomic.Int32

	limited := RateLimit(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if calls.Add(1) == 1 {
			close(firstStarted)
			<-releaseFirst
			return
		}
		close(secondStarted)
	}), 100, 10)

	firstDone := make(chan struct{})
	go func() {
		defer close(firstDone)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "198.51.100.11:1234"
		limited.ServeHTTP(httptest.NewRecorder(), req)
	}()

	<-firstStarted
	go func() {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "198.51.100.12:1234"
		limited.ServeHTTP(httptest.NewRecorder(), req)
	}()

	select {
	case <-secondStarted:
	case <-time.After(time.Second):
		close(releaseFirst)
		t.Fatal("second request was blocked while the first handler was running")
	}
	close(releaseFirst)
	<-firstDone
}
