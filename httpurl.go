package httpurl

import (
	"fmt"
	"github.com/pkg/errors"
	"net/url"
	"path"
	"strings"
)

type Builder struct {
	url   *url.URL
	query url.Values
}

type literal string

// NewBuilder returns an empty Builder. You must set the scheme and the host or else the Build() function
// will error.
func NewBuilder() *Builder {
	return &Builder{
		url:   &url.URL{},
		query: url.Values{},
	}
}

// FromString returns a Builder.
func FromString(rawurl string) (*Builder, error) {
	u, err := url.ParseRequestURI(rawurl)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse string URL: %s", rawurl)
	}

	return &Builder{
		url:   u,
		query: u.Query(),
	}, nil
}

// FromLiteral returns a Builder. The rawurl parameter must be a literal string since the type is not exported. This
// function panics instead of returning an error.
func FromLiteral(rawurl literal) *Builder {
	u, err := url.ParseRequestURI(string(rawurl))
	if err != nil {
		panic(fmt.Sprintf("failed to parse literal URL: %s", rawurl))
	}

	return &Builder{
		url:   u,
		query: u.Query(),
	}
}

// FromURL returns a Builder.
func FromURL(u *url.URL) *Builder {
	// create a copy of url.URL
	t := *u
	return &Builder{
		url:   &t,
		query: u.Query(),
	}
}

// Scheme sets the scheme. Only "http" and "https" are considered valid schemes.
func (b *Builder) Scheme(scheme string) *Builder {
	b.url.Scheme = scheme
	return b
}

// Host sets the host. Use "host:port" to set the port.
func (b *Builder) Host(host string) *Builder {
	b.url.Host = host
	return b
}

// AddPathSegment adds a path segment. It is safe to use this function with externally controlled data; newpath
// is escaped. A malicious user cannot leverage tricks involving ".." or "/".
func (b *Builder) AddPathSegment(newpath string) *Builder {
	b.url.Path = path.Join(b.url.Path, url.PathEscape(newpath))
	return b
}

// AddQueryParam adds a query parameter. Existing query parameters with the same name are preserved.
func (b *Builder) AddQueryParam(key string, value interface{}) *Builder {
	b.query.Add(key, fmt.Sprintf("%v", value))
	return b
}

// SetQueryParam sets a query parameter. Existing query parameters with the same name are dropped.
func (b *Builder) SetQueryParam(key string, value interface{}) *Builder {
	b.query.Set(key, fmt.Sprintf("%v", value))
	return b
}

// RemoveAllQueryParams removes all query parameters.
func (b *Builder) RemoveAllQueryParams() *Builder {
	b.query = url.Values{}
	return b
}

// RemovePathSegment drops a path segment at a given segment.
func (b *Builder) RemovePathSegment(index int) *Builder {
	p := strings.Split(b.url.Path, "/")
	var r []string
	for i := 1; i < len(p); i++ {
		if i == index+1 {
			continue
		}
		r = append(r, p[i])
	}
	b.url.Path = path.Join(r...)
	return b
}

// Build builds the final *url.URL.
func (b *Builder) Build() (*url.URL, error) {
	if b.url.Scheme != "http" && b.url.Scheme != "https" {
		return nil, errors.New(fmt.Sprintf("invalid scheme: %s", b.url.Scheme))
	}
	if b.url.Hostname() == "" {
		return nil, errors.New("host not set")
	}

	// make a copy before returning
	t := *b.url

	// set the query
	t.RawQuery = b.query.Encode()

	// fix the path if empty
	if t.Path == "" {
		t.Path = "/"
	}

	return &t, nil
}

// IsDomain is a helper function to check if a domain matches a url.URL. E.g. given http://www.example.com/,
// IsDomain returns true only for "www.example.com".
func IsDomain(u *url.URL, domain string) bool {
	return u.Hostname() == domain
}

// IsSubdomainOf is a helper function to check if a url.URL is a subdomain of a given domain. E.g. given
// http://www.example.com/, IsSubdomainOf returns true for "example.com", as well as for "com".
func IsSubdomainOf(u *url.URL, domain string) bool {
	return strings.HasSuffix(u.Hostname(), "."+domain)
}

// IsDomainOrSubdomainOf is a helper function to if a url.URL is a domain or subdomain of a given domain. E.g. given
// http://www.example.com/, IsDomainOrSubdomainOf returns true for "www.example.com", "example.com", and "com".
func IsDomainOrSubdomainOf(u *url.URL, domain string) bool {
	return IsDomain(u, domain) || IsSubdomainOf(u, domain)
}
