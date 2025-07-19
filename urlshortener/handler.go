package urlshort

import (
	"gopkg.in/yaml.v2"
	"net/http"
)

// PathToUrl represents a single path-to-URL mapping from YAML
type PathToUrl struct {
	Path string `yaml:"path"`
	Url  string `yaml:"url"`
}

// MapHandler will return an http.HandlerFunc (which also
// implements http.Handler) that will attempt to map any
// paths (keys in the map) to their corresponding URL (values
// that each key in the map points to, in string format).
// If the path is not provided in the map, then the fallback
// http.Handler will be called instead.
func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if path, ok := pathsToUrls[r.URL.Path]; ok {
			http.Redirect(w, r, path, 302)
		} else {
			fallback.ServeHTTP(w, r)
		}
	}
}

// YAMLHandler will parse the provided YAML and then return
// an http.HandlerFunc (which also implements http.Handler)
// that will attempt to map any paths to their corresponding
// URL. If the path is not provided in the YAML, then the
// fallback http.Handler will be called instead.
//
// YAML is expected to be in the format:
//
//   - path: /some-path
//     url: https://www.some-url.com/demo
//
// The only errors that can be returned all related to having
// invalid YAML data.
//
// See MapHandler to create a similar http.HandlerFunc via
// a mapping of paths to urls.
func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	yamlPathToUrls := []PathToUrl{}
	pathToUrls := make(map[string]string)

	err := yaml.Unmarshal(yml, &yamlPathToUrls)
	if err != nil {
		return nil, err
	}

	for _, v := range yamlPathToUrls {
		pathToUrls[v.Path] = v.Url
	}

	return MapHandler(pathToUrls, fallback), nil
}
