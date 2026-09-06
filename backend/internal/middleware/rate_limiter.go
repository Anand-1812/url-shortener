// package middleware with the function of rate limiting
package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/tomasen/realip"
	"golang.org/x/time/rate"
)

type client struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

func RateLimit(next http.Handler, rps, bursts int) http.Handler {
	var mu sync.Mutex
	clients := make(map[string]*client)

	go func() {
		for {
			time.Sleep(time.Minute)
			mu.Lock()

			for ip, c := range clients {
				if time.Since(c.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := realip.FromRequest(r)

		mu.Lock()
		c, exist := clients[ip]

		if !exist {
			c = &client{
				limiter: rate.NewLimiter(rate.Limit(rps), bursts),
			}
			clients[ip] = c
		}

		c.lastSeen = time.Now()

		allowed := c.limiter.Allow()
		mu.Unlock()

		if !allowed {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte(`{"error": "too many requests, please slow down"}`))
			return
		}

		next.ServeHTTP(w, r)
	})
}
