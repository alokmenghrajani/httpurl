package httpurl

import (
	"github.com/stretchr/testify/require"
	"net/url"
	"testing"
)

func TestNewBuilder(t *testing.T) {
	builder := NewBuilder()
	builder.Scheme("foobar")
	_, err := builder.Build()
	require.Error(t, err)

	builder.Scheme("http")
	_, err = builder.Build()
	require.Error(t, err)

	builder.Host("example.com")
	u, err := builder.Build()
	require.NoError(t, err)
	require.Equal(t, "http://example.com/", u.String())
}

func TestFromLiteral(t *testing.T) {
	panicF := func() {
		_ = FromLiteral("not a url\n")
	}
	require.Panics(t, panicF)

	builder := FromLiteral("http://example.com/foo/bar?a=1&b=2")
	builder.AddQueryParam("b", 3)
	u, err := builder.Build()
	require.NoError(t, err)
	require.Equal(t, "http://example.com/foo/bar?a=1&b=2&b=3", u.String())
}

func TestFromString(t *testing.T) {
	s := "not a url\n"
	_, err := FromString(s)
	require.Error(t, err)

	s = "http://example.com/foo/bar"
	s += "?a=1&b=2"
	builder, err := FromString(s)
	require.NoError(t, err)

	builder.SetQueryParam("b", 3)
	u, err := builder.Build()
	require.NoError(t, err)
	require.Equal(t, "http://example.com/foo/bar?a=1&b=3", u.String())
}

func TestFromURL(t *testing.T) {
	u := &url.URL{
		Scheme: "http",
		Host:   "example.com",
		Path:   "/foo/bar",
	}
	builder := FromURL(u)
	builder.SetQueryParam("b", 3)
	u, err := builder.Build()
	require.NoError(t, err)
	require.Equal(t, "http://example.com/foo/bar?b=3", u.String())

	builder.AddPathSegment("../baz/meh")
	u, err = builder.Build()
	require.NoError(t, err)
	require.Equal(t, "http://example.com/foo/bar/..%252Fbaz%252Fmeh?b=3", u.String())

	builder.RemovePathSegment(2)
	u, err = builder.Build()
	require.NoError(t, err)
	require.Equal(t, "http://example.com/foo/bar?b=3", u.String())

	builder.RemoveAllQueryParams()
	builder.Host("www.example.com")
	u, err = builder.Build()
	require.NoError(t, err)
	require.Equal(t, "http://www.example.com/foo/bar", u.String())

	require.False(t, IsDomain(u, "example"))
	require.False(t, IsDomain(u, "example.com"))
	require.True(t, IsDomain(u, "www.example.com"))

	require.False(t, IsSubdomainOf(u, "example"))
	require.True(t, IsSubdomainOf(u, "example.com"))
	require.False(t, IsSubdomainOf(u, "www.example.com"))

	require.False(t, IsDomainOrSubdomainOf(u, "example"))
	require.True(t, IsDomainOrSubdomainOf(u, "example.com"))
	require.True(t, IsDomainOrSubdomainOf(u, "www.example.com"))
}
