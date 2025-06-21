package main

import (
	"fmt"
	"log"
	"time"

	httpclient "github.com/go-request-client"
)

func main() {
	fmt.Println("=== Go HTTP Client Library Examples ===")

	// Example 1: Basic GET request
	basicGetExample()

	// Example 2: POST with JSON
	postJSONExample()

	// Example 3: Form data
	formDataExample()

	// Example 4: Multipart upload
	multipartExample()

	// Example 5: Authentication
	authExample()

	// Example 6: Async requests
	asyncExample()

	// Example 7: Concurrent requests
	concurrentExample()

	// Example 8: Middleware
	middlewareExample()
}

func basicGetExample() {
	fmt.Println("1. Basic GET Request:")
	
	client := httpclient.NewClient(
		httpclient.WithBaseURL("https://httpbin.org"),
		httpclient.WithTimeout(10*time.Second),
	)

	resp, err := client.Get("/get", &httpclient.RequestOptions{
		QueryParams: map[string]string{
			"param1": "value1",
			"param2": "value2",
		},
		Headers: map[string]string{
			"User-Agent": "Go-HTTP-Client/1.0",
		},
	})

	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Status: %d\n", resp.GetStatusCode())
	fmt.Printf("Body: %s\n\n", resp.GetBody())
}

func postJSONExample() {
	fmt.Println("2. POST with JSON:")

	client := httpclient.NewClient(
		httpclient.WithBaseURL("https://httpbin.org"),
	)

	data := map[string]interface{}{
		"name":    "John Doe",
		"email":   "john@example.com",
		"age":     30,
		"active":  true,
	}

	resp, err := client.Post("/post", &httpclient.RequestOptions{
		JSON: data,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	})

	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Status: %d\n", resp.GetStatusCode())
	fmt.Printf("Response: %s\n\n", resp.GetBody())
}

func formDataExample() {
	fmt.Println("3. Form Data POST:")

	client := httpclient.NewClient(
		httpclient.WithBaseURL("https://httpbin.org"),
	)

	resp, err := client.Post("/post", &httpclient.RequestOptions{
		FormData: map[string]string{
			"username": "john_doe",
			"password": "secret123",
			"action":   "login",
		},
	})

	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Status: %d\n", resp.GetStatusCode())
	fmt.Printf("Response: %s\n\n", resp.GetBody())
}

func multipartExample() {
	fmt.Println("4. Multipart Upload:")

	client := httpclient.NewClient(
		httpclient.WithBaseURL("https://httpbin.org"),
	)

	multipartData := httpclient.NewMultipartData()
	multipartData.AddField("description", "Test file upload")
	multipartData.AddFileFromBytes("file", "test.txt", []byte("Hello, World!"))

	resp, err := client.Post("/post", &httpclient.RequestOptions{
		Multipart: multipartData,
	})

	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Status: %d\n", resp.GetStatusCode())
	fmt.Printf("Response: %s\n\n", resp.GetBody())
}

func authExample() {
	fmt.Println("5. Authentication:")

	client := httpclient.NewClient(
		httpclient.WithBaseURL("https://httpbin.org"),
		httpclient.WithAuth("user", "pass"),
	)

	resp, err := client.Get("/basic-auth/user/pass", nil)

	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Status: %d\n", resp.GetStatusCode())
	fmt.Printf("Response: %s\n\n", resp.GetBody())
}

func asyncExample() {
	fmt.Println("6. Async Requests:")

	asyncClient := httpclient.NewAsyncClient(
		httpclient.WithBaseURL("https://httpbin.org"),
	)

	promise := asyncClient.GetAsync("/delay/1", nil)
	
	fmt.Println("Request sent asynchronously...")
	
	resp, err := promise.Wait()
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}

	fmt.Printf("Async response status: %d\n", resp.GetStatusCode())
	
	// Using Then/Catch
	promise2 := asyncClient.GetAsync("/get", nil)
	promise2.Then(func(resp *httpclient.Response) (*httpclient.Response, error) {
		fmt.Printf("Then callback - Status: %d\n", resp.GetStatusCode())
		return resp, nil
	}).Catch(func(err error) error {
		fmt.Printf("Catch callback - Error: %v\n", err)
		return err
	})
	
	promise2.Wait()
	fmt.Println()
}

func concurrentExample() {
	fmt.Println("7. Concurrent Requests:")

	asyncClient := httpclient.NewAsyncClient(
		httpclient.WithBaseURL("https://httpbin.org"),
	)

	requests := []httpclient.ConcurrentRequest{
		{Method: "GET", Path: "/get", Options: &httpclient.RequestOptions{}},
		{Method: "GET", Path: "/delay/1", Options: &httpclient.RequestOptions{}},
		{Method: "GET", Path: "/delay/2", Options: &httpclient.RequestOptions{}},
	}

	fmt.Println("Sending concurrent requests...")
	results := asyncClient.SendConcurrent(requests)

	for i, result := range results {
		if result.Error != nil {
			fmt.Printf("Request %d failed: %v\n", i, result.Error)
		} else {
			fmt.Printf("Request %d succeeded: %d\n", i, result.Response.GetStatusCode())
		}
	}
	fmt.Println()
}

func middlewareExample() {
	fmt.Println("8. Middleware Example:")

	// Note: In a real implementation, you would integrate middleware
	// into the client's request pipeline
	fmt.Println("Middleware would be integrated into the client pipeline")
	fmt.Println("This would provide logging, retry, timeout, and other features")
	fmt.Println()
} 