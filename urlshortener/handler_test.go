package urlshort

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testCase struct {
	name             string
	path             string
	fallbackHandler  http.Handler
	expectedStatus   int
	expectedLocation string
	expectedBody     string
}

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

func TestMapHandler(t *testing.T) {
	tests := []testCase{
		{
			name:             "Path found - redirects correctly",
			path:             "/test",
			fallbackHandler:  fallback,
			expectedStatus:   http.StatusFound,
			expectedLocation: "https://godoc.org/github.com/gophercises/urlshort",
			expectedBody:     "",
		},
		{
			name:             "Path not found - fallback used",
			path:             "/missing",
			fallbackHandler:  fallback,
			expectedStatus:   http.StatusOK,
			expectedLocation: "",
			expectedBody:     "Page not found",
		},
		{
			name: "Path found - fallback not called",
			path: "/test",
			fallbackHandler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				t.Error("Fallback was unexpectedly called when path was fond")
			}),
			expectedStatus:   http.StatusFound,
			expectedLocation: "https://godoc.org/github.com/gophercises/urlshort",
			expectedBody:     "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			handler := MapHandler(testPathsToUrls, tc.fallbackHandler)

			req := httptest.NewRequest("GET", tc.path, nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			if w.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, w.Code)
			}

			if tc.expectedLocation != "" && w.Header().Get("Location") != tc.expectedLocation {
				t.Errorf("expected location %s, got %s", tc.expectedLocation, w.Header().Get("Location"))
			}

			if tc.expectedBody != "" && w.Body.String() != tc.expectedBody {
				t.Errorf("expected body %s, got %s", tc.expectedBody, w.Body.String())
			}
		})
	}
}

func TestYAMLHandler(t *testing.T) {
	tests := []testCase{
		{
			name:             "Path found - redirects from fallback map",
			path:             "/test",
			fallbackHandler:  MapHandler(testPathsToUrls, fallback),
			expectedStatus:   http.StatusFound,
			expectedLocation: "https://godoc.org/github.com/gophercises/urlshort",
			expectedBody:     "",
		},
		{
			name:             "Path not found - shows fallback's fallback",
			path:             "/missing",
			fallbackHandler:  MapHandler(testPathsToUrls, fallback),
			expectedStatus:   http.StatusOK,
			expectedLocation: "",
			expectedBody:     "Page not found",
		},
		{
			name:             "Path found - redirects from yaml",
			path:             "/urlshort-final",
			fallbackHandler:  MapHandler(testPathsToUrls, fallback),
			expectedStatus:   http.StatusFound,
			expectedLocation: "https://github.com/gophercises/urlshort/tree/solution",
			expectedBody:     "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mapHandler := MapHandler(testPathsToUrls, tc.fallbackHandler)
			yamlHandler, err := YAMLHandler([]byte(validYAML), mapHandler)

			if err != nil {
				t.Errorf("could not create handler. error %s", err)
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", tc.path, nil)

			yamlHandler.ServeHTTP(w, r)

			if w.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, w.Code)
			}

			if tc.expectedLocation != "" && w.Header().Get("Location") != tc.expectedLocation {
				t.Errorf("expected location %s, got %s", tc.expectedLocation, w.Header().Get("Location"))
			}

			if tc.expectedBody != "" && w.Body.String() != tc.expectedBody {
				t.Errorf("expected body %s, got %s", tc.expectedBody, w.Body.String())
			}
		})
	}
}

func TestYAMLHandler_ParseError(t *testing.T) {
	mapHandler := MapHandler(testPathsToUrls, fallback)

	_, err := YAMLHandler([]byte(invalidYAML), mapHandler)

	if err == nil {
		t.Errorf("Expected YAML parsing error. Got nil.")
	}
}

// func TestMapHandler_FallbackUsed(t *testing.T) {
// 	handler := MapHandler(testPathsToUrls, fallback)
//
// 	req := httptest.NewRequest("GET", "/missing", nil)
//
// 	w := httptest.NewRecorder()
//
// 	handler(w, req)
//
// 	response := w.Body.String()
// 	expected := "Page not found"
// 	if response != expected {
// 		t.Errorf("Expected '%s', got '%s'", expected, response)
// 	}
// }
//
// func TestMapHandler_PathFound(t *testing.T) {
// 	redirectUrl := "https://godoc.org/github.com/gophercises/urlshort"
//
// 	handler := MapHandler(testPathsToUrls, fallback)
//
// 	req := httptest.NewRequest("GET", "/test", nil)
//
// 	w := httptest.NewRecorder()
//
// 	handler(w, req)
//
// 	httpResponseCode := w.Code
// 	httpHeaderLocation := w.Header().Get("Location")
//
// 	if httpResponseCode != http.StatusFound {
// 		t.Errorf("Wrong http status. Expected '302', got '%d'", httpResponseCode)
// 	}
//
// 	if httpHeaderLocation != redirectUrl {
// 		t.Errorf("Wrong http header location. Expected '%s', got '%s'", redirectUrl, httpHeaderLocation)
// 	}
// }
//
// func TestMapHandler_FallbackNotCalled(t *testing.T) {
// 	fallbackCalled := false
//
// 	testFallback := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		fallbackCalled = true
// 	})
//
// 	handler := MapHandler(testPathsToUrls, testFallback)
//
// 	req := httptest.NewRequest("GET", "/test", nil)
//
// 	w := httptest.NewRecorder()
//
// 	handler(w, req)
//
// 	if fallbackCalled {
// 		t.Error("Fallback should not have been called")
// 	}
// }

// func TestYAMLHandler_RedirectsFromFallbackMapHandler(t *testing.T) {
// 	redirectUrl := "https://godoc.org/github.com/gophercises/urlshort"
//
// 	mapHandler := MapHandler(testPathsToUrls, fallback)
//
// 	yaml := ``
//
// 	yamlHandler, err := YAMLHandler([]byte(yaml), mapHandler)
//
// 	if err != nil {
// 		t.Error(err)
// 	}
//
// 	req := httptest.NewRequest("GET", "/test", nil)
// 	w := httptest.NewRecorder()
//
// 	yamlHandler(w, req)
//
// 	httpResponseCode := w.Code
// 	httpHeaderLocation := w.Header().Get("Location")
//
// 	if httpResponseCode != http.StatusFound {
// 		t.Errorf("Wrong http status. Expected '302', got '%d'", httpResponseCode)
// 	}
//
// 	if httpHeaderLocation != redirectUrl {
// 		t.Errorf("Wrong http header location. Expected '%s', got '%s'", redirectUrl, httpHeaderLocation)
// 	}
// }

// func TestYAMLHandler_ShowsFallbacksFallback(t *testing.T) {
// 	mapHandler := MapHandler(testPathsToUrls, fallback)
//
// 	yaml := ``
//
// 	yamlHandler, _ := YAMLHandler([]byte(yaml), mapHandler)
//
// 	req := httptest.NewRequest("GET", "/missing", nil)
// 	w := httptest.NewRecorder()
//
// 	yamlHandler(w, req)
//
// 	expected := "Page not found"
// 	response := w.Body.String()
//
// 	if response != expected {
// 		t.Errorf("Expected '%s', got '%s'", expected, response)
// 	}
// }

// func TestYAMLHandler_EntryFound(t *testing.T) {
// 	redirectUrl := "https://github.com/gophercises/urlshort/tree/solution"
//
// 	yamlHandler, _ := YAMLHandler([]byte(validYAML), nil)
//
// 	req := httptest.NewRequest("GET", "/urlshort-final", nil)
// 	w := httptest.NewRecorder()
//
// 	yamlHandler(w, req)
//
// 	httpResponseCode := w.Code
// 	httpHeaderLocation := w.Header().Get("Location")
//
// 	if httpResponseCode != http.StatusFound {
// 		t.Errorf("Wrong http status. Expected '302', got '%d'", httpResponseCode)
// 	}
//
// 	if httpHeaderLocation != redirectUrl {
// 		t.Errorf("Wrong http header location. Expected '%s', got '%s'", redirectUrl, httpHeaderLocation)
// 	}
// }
