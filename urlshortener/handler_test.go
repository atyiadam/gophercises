package urlshort

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

const (
	validYAML = `
- path: /urlshort
  url: https://github.com/gophercises/urlshort
- path: /urlshort-final
  url: https://github.com/gophercises/urlshort/tree/solution
`
	invalidYAML = `
  - path: /urlshort
  url: https://github.com/gophercises/urlshort
  `
)

var testPathsToUrls = map[string]string{
	"/test": "https://godoc.org/github.com/gophercises/urlshort",
}

var fallback = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Page not found")
})

func TestMapHandler_ReturnsNonNilHandler(t *testing.T) {
	handler := MapHandler(testPathsToUrls, fallback)

	if handler == nil {
		t.Error("MapHandler returned nil, expected a handler function")
	}
}

func TestMapHandler_FallbackUsed(t *testing.T) {
	handler := MapHandler(testPathsToUrls, fallback)

	req := httptest.NewRequest("GET", "/missing", nil)

	w := httptest.NewRecorder()

	handler(w, req)

	response := w.Body.String()
	expected := "Page not found"
	if response != expected {
		t.Errorf("Expected '%s', got '%s'", expected, response)
	}
}

func TestMapHandler_PathFound(t *testing.T) {
	redirectUrl := "https://godoc.org/github.com/gophercises/urlshort"

	handler := MapHandler(testPathsToUrls, fallback)

	req := httptest.NewRequest("GET", "/test", nil)

	w := httptest.NewRecorder()

	handler(w, req)

	httpResponseCode := w.Code
	httpHeaderLocation := w.Header().Get("Location")

	if httpResponseCode != http.StatusFound {
		t.Errorf("Wrong http status. Expected '302', got '%d'", httpResponseCode)
	}

	if httpHeaderLocation != redirectUrl {
		t.Errorf("Wrong http header location. Expected '%s', got '%s'", redirectUrl, httpHeaderLocation)
	}
}

func TestMapHandler_FallbackNotCalled(t *testing.T) {
	fallbackCalled := false

	testFallback := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fallbackCalled = true
	})

	handler := MapHandler(testPathsToUrls, testFallback)

	req := httptest.NewRequest("GET", "/test", nil)

	w := httptest.NewRecorder()

	handler(w, req)

	if fallbackCalled {
		t.Error("Fallback should not have been called")
	}
}

func TestYAMLHandler_RedirectsFromFallbackMapHandler(t *testing.T) {
	redirectUrl := "https://godoc.org/github.com/gophercises/urlshort"

	mapHandler := MapHandler(testPathsToUrls, fallback)

	yaml := ``

	yamlHandler, err := YAMLHandler([]byte(yaml), mapHandler)

	if err != nil {
		t.Error(err)
	}

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	yamlHandler(w, req)

	httpResponseCode := w.Code
	httpHeaderLocation := w.Header().Get("Location")

	if httpResponseCode != http.StatusFound {
		t.Errorf("Wrong http status. Expected '302', got '%d'", httpResponseCode)
	}

	if httpHeaderLocation != redirectUrl {
		t.Errorf("Wrong http header location. Expected '%s', got '%s'", redirectUrl, httpHeaderLocation)
	}
}

func TestYAMLHandler_ShowsFallbacksFallback(t *testing.T) {
	mapHandler := MapHandler(testPathsToUrls, fallback)

	yaml := ``

	yamlHandler, _ := YAMLHandler([]byte(yaml), mapHandler)

	req := httptest.NewRequest("GET", "/missing", nil)
	w := httptest.NewRecorder()

	yamlHandler(w, req)

	expected := "Page not found"
	response := w.Body.String()

	if response != expected {
		t.Errorf("Expected '%s', got '%s'", expected, response)
	}
}

func TestYAMLHandler_EntryFound(t *testing.T) {
	redirectUrl := "https://github.com/gophercises/urlshort/tree/solution"

	yamlHandler, _ := YAMLHandler([]byte(validYAML), nil)

	req := httptest.NewRequest("GET", "/urlshort-final", nil)
	w := httptest.NewRecorder()

	yamlHandler(w, req)

	httpResponseCode := w.Code
	httpHeaderLocation := w.Header().Get("Location")

	if httpResponseCode != http.StatusFound {
		t.Errorf("Wrong http status. Expected '302', got '%d'", httpResponseCode)
	}

	if httpHeaderLocation != redirectUrl {
		t.Errorf("Wrong http header location. Expected '%s', got '%s'", redirectUrl, httpHeaderLocation)
	}
}

func TestYAMLHandler_YAMLParseError(t *testing.T) {
	mapHandler := MapHandler(testPathsToUrls, fallback)

	yamlHandler, err := YAMLHandler([]byte(invalidYAML), mapHandler)

	if err == nil {
		t.Errorf("Expected YAML parsing error. Got nil.")
	}

	req := httptest.NewRequest("GET", "/missing", nil)
	w := httptest.NewRecorder()

	yamlHandler(w, req)

	expected := "Page not found"
	response := w.Body.String()

	if response != expected {
		t.Errorf("Expected '%s', got '%s'", expected, response)
	}
}
