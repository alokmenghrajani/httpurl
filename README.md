# httpurl
[![license](http://img.shields.io/badge/license-apache_2.0-blue.svg?style=flat)](https://raw.githubusercontent.com/alokmenghrajani/httpurl/master/LICENSE)
[![travis](https://travis-ci.com/alokmenghrajani/httpurl.svg?branch=main)](https://travis-ci.com/github/alokmenghrajani/httpurl)
[![coverage](https://coveralls.io/repos/github/alokmenghrajani/httpurl/badge.svg?branch=main&service=github)](https://coveralls.io/github/alokmenghrajani/httpurl)
[![report](https://goreportcard.com/badge/github.com/alokmenghrajani/httpurl)](https://goreportcard.com/report/github.com/alokmenghrajani/httpurl)

Minimalistic Go library to make handling http URLs easier.

`httpurl` complements Golang's `url.URL` from the `net/url` standard library. The desire to implement `httpurl` stems
from:
- noticing developers are using string manipulation functions to manipulate URLs/URIs. Inevitably, we end up
  with bits and pieces of incorrect code (especially when mixing regular expressions and URLs for domain validation).
- the lack of `url.MustParse`.
- query manipulation requiring going back and forth from `url.Query` to `url.RawQuery`.
- somewhat inspired by Java's [OkHttpUrl](https://github.com/square/okhttp/blob/okhttp_4.9.x/okhttp/src/main/kotlin/okhttp3/HttpUrl.kt)
  library. An early attempt to use a builder pattern didn't make production code any shorter or easier to write.  

# Documentation

[https://pkg.go.dev/github.com/alokmenghrajani/httpurl](https://pkg.go.dev/github.com/alokmenghrajani/httpurl)

# Example

Instead of:

```
u, err := url.Parse("http://example.com/")
if err != nil {
    // Need to handle this error
}
q := u.Query()
q.Add("code", "1234")
u.RawQuery = q.Encode()
```

You can do:

```
u := httpurl.MustParse("http://example.com")
AddQueryParam(u, "code", 1234)
```

<p align="center">* * *</p>

And instead of :

```
basePath := "https://example.com/users/%d/products/%s/"

baseURL, err := url.Parse(fmt.Sprintf(basePath, userId, productId))
if err != nil {
    // Need to handle this error
}
```

You can do:

```
basePath := "https://example.com/users/{userId}/products/{productId}/"

baseURL := httpurl.MustParse(basePath)
httpurl.ExpandPath(baseURL, httpurl.ExpandMap{
  "userId": userId,
  "productId": productId,
})
```
