package httpclient

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// Middleware represents a middleware function
type Middleware func(next Handler) Handler

// Handler represents an HTTP handler
type Handler func(*http.Request) (*Response, error)

// HandlerStack represents a stack of middleware handlers
type HandlerStack struct {
	handler  Handler
	stack    []Middleware
	position int
}

// NewHandlerStack creates a new handler stack
func NewHandlerStack(handler Handler) *HandlerStack {
	return &HandlerStack{
		handler:  handler,
		stack:    make([]Middleware, 0),
		position: 0,
	}
}

// Push adds middleware to the stack
func (hs *HandlerStack) Push(middleware Middleware) {
	hs.stack = append(hs.stack, middleware)
}

// Next executes the next middleware in the stack
func (hs *HandlerStack) Next(req *http.Request) (*Response, error) {
	if hs.position >= len(hs.stack) {
		return hs.handler(req)
	}

	middleware := hs.stack[hs.position]
	hs.position++

	return middleware(hs.Next)(req)
}

// Reset resets the stack position
func (hs *HandlerStack) Reset() {
	hs.position = 0
}

// Common middleware functions

// LoggingMiddleware logs request and response information
func LoggingMiddleware(logger Logger) Middleware {
	return func(next Handler) Handler {
		return func(req *http.Request) (*Response, error) {
			logger.Logf("Request: %s %s", req.Method, req.URL.String())
			
			resp, err := next(req)
			if err != nil {
				logger.Logf("Error: %v", err)
				return nil, err
			}
			
			logger.Logf("Response: %d", resp.StatusCode)
			return resp, nil
		}
	}
}

// RetryMiddleware retries failed requests
func RetryMiddleware(maxRetries int, backoff BackoffStrategy) Middleware {
	return func(next Handler) Handler {
		return func(req *http.Request) (*Response, error) {
			var lastErr error
			
			for attempt := 0; attempt <= maxRetries; attempt++ {
				resp, err := next(req)
				if err == nil {
					return resp, nil
				}
				
				lastErr = err
				if attempt < maxRetries {
					delay := backoff.Delay(attempt)
					select {
					case <-req.Context().Done():
						return nil, req.Context().Err()
					case <-time.After(delay):
						continue
					}
				}
			}
			
			return nil, lastErr
		}
	}
}

// TimeoutMiddleware adds timeout to requests
func TimeoutMiddleware(timeout time.Duration) Middleware {
	return func(next Handler) Handler {
		return func(req *http.Request) (*Response, error) {
			ctx, cancel := context.WithTimeout(req.Context(), timeout)
			defer cancel()
			
			req = req.WithContext(ctx)
			return next(req)
		}
	}
}

// Logger interface for logging
type Logger interface {
	Logf(format string, args ...interface{})
}

// SimpleLogger implements Logger interface
type SimpleLogger struct{}

func (sl *SimpleLogger) Logf(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

// BackoffStrategy interface for retry backoff
type BackoffStrategy interface {
	Delay(attempt int) time.Duration
}

// ExponentialBackoff implements exponential backoff
type ExponentialBackoff struct {
	BaseDelay time.Duration
	MaxDelay  time.Duration
}

func (eb *ExponentialBackoff) Delay(attempt int) time.Duration {
	delay := eb.BaseDelay * time.Duration(1<<attempt)
	if delay > eb.MaxDelay {
		delay = eb.MaxDelay
	}
	return delay
} 