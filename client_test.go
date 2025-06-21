package httpclient

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client := NewClient(
		WithBaseURL("https://api.example.com"),
		WithTimeout(5*time.Second),
		WithHeaders(map[string]string{"User-Agent": "Test"}),
	)

	if client.baseURL != "https://api.example.com" {
		t.Errorf("Expected baseURL to be 'https://api.example.com', got '%s'", client.baseURL)
	}

	if client.timeout != 5*time.Second {
		t.Errorf("Expected timeout to be 5s, got %v", client.timeout)
	}

	if client.headers["User-Agent"] != "Test" {
		t.Errorf("Expected User-Agent header to be 'Test', got '%s'", client.headers["User-Agent"])
	}
}

func TestClient_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "success"}`))
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))
	resp, err := client.Get("/test", nil)

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if resp.GetStatusCode() != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.GetStatusCode())
	}

	expected := `{"message": "success"}`
	if resp.GetBody() != expected {
		t.Errorf("Expected body '%s', got '%s'", expected, resp.GetBody())
	}
}

func TestClient_Post_JSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected Content-Type to be application/json, got %s", r.Header.Get("Content-Type"))
		}

		var data map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			t.Errorf("Failed to decode JSON: %v", err)
		}

		if data["name"] != "John" {
			t.Errorf("Expected name to be 'John', got %v", data["name"])
		}

		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"id": 123}`))
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))
	data := map[string]interface{}{"name": "John", "age": 30}

	resp, err := client.Post("/users", &RequestOptions{
		JSON: data,
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if resp.GetStatusCode() != http.StatusCreated {
		t.Errorf("Expected status 201, got %d", resp.GetStatusCode())
	}
}

func TestClient_Post_FormData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
			t.Errorf("Expected Content-Type to be application/x-www-form-urlencoded, got %s", r.Header.Get("Content-Type"))
		}

		if err := r.ParseForm(); err != nil {
			t.Errorf("Failed to parse form: %v", err)
		}

		if r.FormValue("username") != "john_doe" {
			t.Errorf("Expected username to be 'john_doe', got %s", r.FormValue("username"))
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "success"}`))
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))
	resp, err := client.Post("/login", &RequestOptions{
		FormData: map[string]string{
			"username": "john_doe",
			"password": "secret123",
		},
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if resp.GetStatusCode() != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.GetStatusCode())
	}
}

func TestClient_QueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("page") != "1" {
			t.Errorf("Expected page=1, got %s", r.URL.Query().Get("page"))
		}
		if r.URL.Query().Get("limit") != "10" {
			t.Errorf("Expected limit=10, got %s", r.URL.Query().Get("limit"))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"data": []}`))
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))
	resp, err := client.Get("/users", &RequestOptions{
		QueryParams: map[string]string{
			"page":  "1",
			"limit": "10",
		},
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if resp.GetStatusCode() != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.GetStatusCode())
	}
}

func TestClient_Authentication(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			t.Error("Expected basic auth to be present")
		}
		if username != "user" || password != "pass" {
			t.Errorf("Expected auth to be user:pass, got %s:%s", username, password)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"authenticated": true}`))
	}))
	defer server.Close()

	client := NewClient(
		WithBaseURL(server.URL),
		WithAuth("user", "pass"),
	)

	resp, err := client.Get("/protected", nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if resp.GetStatusCode() != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.GetStatusCode())
	}
}

func TestMultipartData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Content-Type"), "multipart/form-data") {
			t.Errorf("Expected multipart content type, got %s", r.Header.Get("Content-Type"))
		}

		if err := r.ParseMultipartForm(32 << 20); err != nil {
			t.Errorf("Failed to parse multipart form: %v", err)
		}

		if r.FormValue("description") != "Test file" {
			t.Errorf("Expected description to be 'Test file', got %s", r.FormValue("description"))
		}

		file, header, err := r.FormFile("file")
		if err != nil {
			t.Errorf("Failed to get file: %v", err)
		}
		defer file.Close()

		if header.Filename != "test.txt" {
			t.Errorf("Expected filename to be 'test.txt', got %s", header.Filename)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"uploaded": true}`))
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))
	multipartData := NewMultipartData()
	multipartData.AddField("description", "Test file")
	multipartData.AddFileFromBytes("file", "test.txt", []byte("Hello, World!"))

	resp, err := client.Post("/upload", &RequestOptions{
		Multipart: multipartData,
	})

	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if resp.GetStatusCode() != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.GetStatusCode())
	}
}

func TestResponse_UnmarshalJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"name": "John", "age": 30}`))
	}))
	defer server.Close()

	client := NewClient(WithBaseURL(server.URL))
	resp, err := client.Get("/user", nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	var data struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	if err := resp.UnmarshalJSON(&data); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if data.Name != "John" {
		t.Errorf("Expected name to be 'John', got %s", data.Name)
	}

	if data.Age != 30 {
		t.Errorf("Expected age to be 30, got %d", data.Age)
	}
} 