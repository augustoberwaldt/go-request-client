package httpclient

import (
	"context"
	"sync"
)

// Promise represents an asynchronous request result
type Promise struct {
	response *Response
	err      error
	done     chan struct{}
	once     sync.Once
}

// NewPromise creates a new promise
func NewPromise() *Promise {
	return &Promise{
		done: make(chan struct{}),
	}
}

// Resolve resolves the promise with a response
func (p *Promise) Resolve(response *Response) {
	p.once.Do(func() {
		p.response = response
		close(p.done)
	})
}

// Reject rejects the promise with an error
func (p *Promise) Reject(err error) {
	p.once.Do(func() {
		p.err = err
		close(p.done)
	})
}

// Wait waits for the promise to be resolved or rejected
func (p *Promise) Wait() (*Response, error) {
	<-p.done
	return p.response, p.err
}

// Then executes a function when the promise is resolved
func (p *Promise) Then(fn func(*Response) (*Response, error)) *Promise {
	newPromise := NewPromise()
	
	go func() {
		resp, err := p.Wait()
		if err != nil {
			newPromise.Reject(err)
			return
		}
		
		newResp, newErr := fn(resp)
		if newErr != nil {
			newPromise.Reject(newErr)
			return
		}
		
		newPromise.Resolve(newResp)
	}()
	
	return newPromise
}

// Catch executes a function when the promise is rejected
func (p *Promise) Catch(fn func(error) error) *Promise {
	newPromise := NewPromise()
	
	go func() {
		resp, err := p.Wait()
		if err != nil {
			newErr := fn(err)
			newPromise.Reject(newErr)
			return
		}
		
		newPromise.Resolve(resp)
	}()
	
	return newPromise
}

// AsyncClient extends Client with async functionality
type AsyncClient struct {
	*Client
}

// NewAsyncClient creates a new async client
func NewAsyncClient(options ...ClientOption) *AsyncClient {
	return &AsyncClient{
		Client: NewClient(options...),
	}
}

// SendAsync sends an asynchronous request
func (ac *AsyncClient) SendAsync(method, path string, options *RequestOptions) *Promise {
	promise := NewPromise()
	
	go func() {
		resp, err := ac.Request(method, path, options)
		if err != nil {
			promise.Reject(err)
			return
		}
		
		promise.Resolve(resp)
	}()
	
	return promise
}

// GetAsync sends an asynchronous GET request
func (ac *AsyncClient) GetAsync(path string, options *RequestOptions) *Promise {
	return ac.SendAsync("GET", path, options)
}

// PostAsync sends an asynchronous POST request
func (ac *AsyncClient) PostAsync(path string, options *RequestOptions) *Promise {
	return ac.SendAsync("POST", path, options)
}

// PutAsync sends an asynchronous PUT request
func (ac *AsyncClient) PutAsync(path string, options *RequestOptions) *Promise {
	return ac.SendAsync("PUT", path, options)
}

// DeleteAsync sends an asynchronous DELETE request
func (ac *AsyncClient) DeleteAsync(path string, options *RequestOptions) *Promise {
	return ac.SendAsync("DELETE", path, options)
}

// PatchAsync sends an asynchronous PATCH request
func (ac *AsyncClient) PatchAsync(path string, options *RequestOptions) *Promise {
	return ac.SendAsync("PATCH", path, options)
}

// ConcurrentRequest represents a concurrent request
type ConcurrentRequest struct {
	Method  string
	Path    string
	Options *RequestOptions
}

// ConcurrentResponse represents a concurrent response
type ConcurrentResponse struct {
	Index    int
	Response *Response
	Error    error
}

// SendConcurrent sends multiple requests concurrently
func (ac *AsyncClient) SendConcurrent(requests []ConcurrentRequest) []ConcurrentResponse {
	results := make([]ConcurrentResponse, len(requests))
	var wg sync.WaitGroup
	
	for i, req := range requests {
		wg.Add(1)
		go func(index int, method, path string, options *RequestOptions) {
			defer wg.Done()
			
			resp, err := ac.Request(method, path, options)
			results[index] = ConcurrentResponse{
				Index:    index,
				Response: resp,
				Error:    err,
			}
		}(i, req.Method, req.Path, req.Options)
	}
	
	wg.Wait()
	return results
}

// SendConcurrentWithLimit sends multiple requests concurrently with a limit
func (ac *AsyncClient) SendConcurrentWithLimit(requests []ConcurrentRequest, limit int) []ConcurrentResponse {
	results := make([]ConcurrentResponse, len(requests))
	semaphore := make(chan struct{}, limit)
	var wg sync.WaitGroup
	
	for i, req := range requests {
		wg.Add(1)
		go func(index int, method, path string, options *RequestOptions) {
			defer wg.Done()
			
			semaphore <- struct{}{} // Acquire
			defer func() { <-semaphore }() // Release
			
			resp, err := ac.Request(method, path, options)
			results[index] = ConcurrentResponse{
				Index:    index,
				Response: resp,
				Error:    err,
			}
		}(i, req.Method, req.Path, req.Options)
	}
	
	wg.Wait()
	return results
}

// WaitAll waits for all promises to complete
func WaitAll(promises ...*Promise) []ConcurrentResponse {
	results := make([]ConcurrentResponse, len(promises))
	var wg sync.WaitGroup
	
	for i, promise := range promises {
		wg.Add(1)
		go func(index int, p *Promise) {
			defer wg.Done()
			
			resp, err := p.Wait()
			results[index] = ConcurrentResponse{
				Index:    index,
				Response: resp,
				Error:    err,
			}
		}(i, promise)
	}
	
	wg.Wait()
	return results
} 