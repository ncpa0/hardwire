package configuration

import "fmt"

func GenerateCacheHeader(policy *CachingPolicy) string {
	header := ""

	if policy.NoStore {
		header += "no-store"
		return header
	}

	if policy.Private {
		header += "private"
	} else {
		header += "public"
	}

	if policy.MaxAge != 0 {
		header += fmt.Sprintf(", max-age=%d, must-revalidate", policy.MaxAge)
	} else if policy.NoCache {
		header += ", no-cache"
	}

	return header
}

func GenerateCacheHeaderForStaticRoute() string {
	return GenerateCacheHeader(Current.Caching.StaticRoutes)
}

func GenerateCacheHeaderForDynamicRoute() string {
	return GenerateCacheHeader(Current.Caching.DynamicRoutes)
}

func GenerateCacheHeaderForFragments() string {
	return GenerateCacheHeader(Current.Caching.Fragments)
}
