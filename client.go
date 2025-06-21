package httpclient

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client represents an HTTP client similar to Guzzle
type Client struct {
	httpClient *http.Client
	baseURL    string
	headers    map[string]string
	timeout    time.Duration
	auth       *Auth
}

// Auth represents authentication credentials
type Auth struct {
	Username string
	Password string
}

// RequestOptions represents options for HTTP requests
type RequestOptions struct {
	Headers     map[string]string
	QueryParams map[string]string
	FormData    map[string]string
	JSON        interface{}
	Body        io.Reader
	Timeout     time.Duration
	Auth        *Auth
	Cookies     []*http.Cookie
	AllowRedirects bool
	Multipart   *MultipartData
}

// Response represents an HTTP response
type Response struct {
	*http.Response
	Body []byte
}

// NewClient creates a new HTTP client
func NewClient(options ...ClientOption) *Client {
	client := &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		headers: make(map[string]string),
	}

	for _, option := range options {
		option(client)
	}

	return client
}

// ClientOption is a function that configures a client
type ClientOption func(*Client)

// WithBaseURL sets the base URL for the client
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// WithTimeout sets the timeout for requests
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.timeout = timeout
		c.httpClient.Timeout = timeout
	}
}

// WithHeaders sets default headers for all requests
func WithHeaders(headers map[string]string) ClientOption {
	return func(c *Client) {
		for k, v := range headers {
			c.headers[k] = v
		}
	}
}

// WithAuth sets authentication credentials
func WithAuth(username, password string) ClientOption {
	return func(c *Client) {
		c.auth = &Auth{Username: username, Password: password}
	}
}

// Request sends an HTTP request
func (c *Client) Request(method, path string, options *RequestOptions) (*Response, error) {
	if options == nil {
		options = &RequestOptions{}
	}

	// Build URL
	requestURL := c.buildURL(path)
	if len(options.QueryParams) > 0 {
		requestURL = c.addQueryParams(requestURL, options.QueryParams)
	}

	// Prepare body
	body, contentType, err := c.prepareBody(options)
	if err != nil {
		return nil, err
	}

	// Create request
	req, err := http.NewRequest(method, requestURL, body)
	if err != nil {
		return nil, err
	}

	// Set headers
	c.setHeaders(req, options.Headers, contentType)

	// Set authentication
	if options.Auth != nil {
		req.SetBasicAuth(options.Auth.Username, options.Auth.Password)
	} else if c.auth != nil {
		req.SetBasicAuth(c.auth.Username, c.auth.Password)
	}

	// Set cookies
	for _, cookie := range options.Cookies {
		req.AddCookie(cookie)
	}

	// Send request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}

	return &Response{
		Response: resp,
		Body:     respBody,
	}, nil
}

// Get sends a GET request
func (c *Client) Get(path string, options *RequestOptions) (*Response, error) {
	return c.Request("GET", path, options)
}

// Post sends a POST request
func (c *Client) Post(path string, options *RequestOptions) (*Response, error) {
	return c.Request("POST", path, options)
}

// Put sends a PUT request
func (c *Client) Put(path string, options *RequestOptions) (*Response, error) {
	return c.Request("PUT", path, options)
}

// Delete sends a DELETE request
func (c *Client) Delete(path string, options *RequestOptions) (*Response, error) {
	return c.Request("DELETE", path, options)
}

// Patch sends a PATCH request
func (c *Client) Patch(path string, options *RequestOptions) (*Response, error) {
	return c.Request("PATCH", path, options)
}

// buildURL builds the complete URL
func (c *Client) buildURL(path string) string {
	if c.baseURL == "" {
		return path
	}
	return strings.TrimRight(c.baseURL, "/") + "/" + strings.TrimLeft(path, "/")
}

// addQueryParams adds query parameters to URL
func (c *Client) addQueryParams(requestURL string, params map[string]string) string {
	u, err := url.Parse(requestURL)
	if err != nil {
		return requestURL
	}

	q := u.Query()
	for k, v := range params {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	return u.String()
}

// prepareBody prepares the request body based on options
func (c *Client) prepareBody(options *RequestOptions) (io.Reader, string, error) {
	if options.Body != nil {
		return options.Body, "", nil
	}

	if options.JSON != nil {
		jsonData, err := json.Marshal(options.JSON)
		if err != nil {
			return nil, "", err
		}
		return bytes.NewBuffer(jsonData), "application/json", nil
	}

	if len(options.FormData) > 0 {
		formData := url.Values{}
		for k, v := range options.FormData {
			formData.Set(k, v)
		}
		return strings.NewReader(formData.Encode()), "application/x-www-form-urlencoded", nil
	}

	// Check for multipart data
	if options.Multipart != nil {
		return options.Multipart.ToReader()
	}

	return nil, "", nil
}

// setHeaders sets request headers
func (c *Client) setHeaders(req *http.Request, customHeaders map[string]string, contentType string) {
	// Set default headers
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	// Set custom headers
	for k, v := range customHeaders {
		req.Header.Set(k, v)
	}

	// Set content type if provided
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
}

// GetStatusCode returns the HTTP status code
func (r *Response) GetStatusCode() int {
	return r.StatusCode
}

// GetHeader returns a header value
func (r *Response) GetHeader(name string) string {
	return r.Header.Get(name)
}

// GetBody returns the response body as string
func (r *Response) GetBody() string {
	return string(r.Body)
}

// GetBodyBytes returns the response body as bytes
func (r *Response) GetBodyBytes() []byte {
	return r.Body
}

// UnmarshalJSON unmarshals the response body as JSON
func (r *Response) UnmarshalJSON(v interface{}) error {
	return json.Unmarshal(r.Body, v)
} 