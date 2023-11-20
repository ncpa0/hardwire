package utils

import (
	neturl "net/url"
	"strings"
)

// ex. schema: /product/:id/:revision, url: /product/123/1
// result: {id: 123, revision: 1}
func ParseUrlParams(schema string, url string) map[string]string {
	result := make(map[string]string)
	URL, err := neturl.Parse(url)
	urlPath := URL.Path

	if err != nil {
		return result
	}

	if schema[0] == '/' {
		schema = schema[1:]
	}
	if urlPath[0] == '/' {
		urlPath = urlPath[1:]
	}

	schemaSegments := strings.Split(schema, "/")
	urlSegments := strings.Split(urlPath, "/")

	for i, segment := range schemaSegments {
		if len(urlSegments) <= i {
			break
		}
		if strings.HasPrefix(segment, ":") {
			result[segment[1:]] = urlSegments[i]
		} else if segment != urlSegments[i] {
			break
		}
	}

	return result
}
