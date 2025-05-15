package api

import (
	"fmt"
	"path"
)

const (
	prefix string = "/api"
)

// ApiV1 constructs an API route for version 1 using the provided route parts.
//
// Parameters:
//   - routeParts: A variadic parameter of strings representing the parts of the route to be joined.
//
// Returns:
//
//	A string representing the complete API route for version 1, including the API prefix and version.
func ApiV1(routeParts ...string) string {
	return buildApi(1, routeParts...)
}

func buildApi(version uint8, routeParts ...string) string {
	ver := fmt.Sprintf("v%v", version)
	route := path.Join(prefix, ver)
	for i := range routeParts {
		route = path.Join(route, routeParts[i])
	}
	return route
}
