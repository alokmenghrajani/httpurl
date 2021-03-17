package httpurl

import (
	"fmt"
	"net/url"
	"path"
	"strings"
)

// Parses a rawurl (must be an absolute URI) and panics in case of error.
func MustParse(rawurl string) *url.URL {
	u, err := url.ParseRequestURI(rawurl)
	if err != nil {
		panic(err)
	}
	return u
}

// AddQueryParam adds the value to key in the query parameters. It appends to any existing values associated with key.
func AddQueryParam(u *url.URL, key string, value interface{}) {
	q := u.Query()
	q.Add(key, fmt.Sprintf("%v", value))
	u.RawQuery = q.Encode()
}

// SetQueryParam sets the key to value in the query parameters. It replaces any existing values.
func SetQueryParam(u *url.URL, key string, value interface{}) {
	q := u.Query()
	q.Set(key, fmt.Sprintf("%v", value))
	u.RawQuery = q.Encode()
}

// AddPathSegment adds a path segment. It is safe to use this function with externally controlled data; newpath
// is escaped. A malicious user cannot leverage tricks involving ".." or "/". An empty newpath is ignored.
func AddPathSegment(u *url.URL, newpath string) {
	u.Path = path.Join(u.Path, url.PathEscape(newpath))
}

// RemovePathSegment drops a path segment at a given segment, counting from 0.
// E.g. RemovePathSegment(http://example.com/foo/bar/xyz, 1) would result in http://example.com/foo/xyz
func RemovePathSegment(u *url.URL, index int) {
	p := strings.Split(u.Path, "/")
	var r []string
	for i := 1; i < len(p); i++ {
		if i == index+1 {
			continue
		}
		r = append(r, p[i])
	}
	u.Path = path.Join(r...)
}

// IsDomain checks if a domain matches a url.URL. E.g. given http://www.example.com/, IsDomain returns true only for
// "www.example.com".
func IsDomain(u *url.URL, domain string) bool {
	return u.Hostname() == domain
}

// IsSubdomainOf checks if a url.URL is a subdomain of a given domain. E.g. given http://www.example.com/, IsSubdomainOf
// returns true for "example.com", as well as for "com".
func IsSubdomainOf(u *url.URL, domain string) bool {
	return strings.HasSuffix(u.Hostname(), "."+domain)
}

// IsDomainOrSubdomainOf checks if a url.URL is a domain or subdomain of a given domain. E.g. given
// http://www.example.com/, IsDomainOrSubdomainOf returns true for "www.example.com", "example.com", and "com".
func IsDomainOrSubdomainOf(u *url.URL, domain string) bool {
	return IsDomain(u, domain) || IsSubdomainOf(u, domain)
}
