package githubheaderauth

import (
	"net/http"
)

// GithuHeaderTransport is an http.RoundTripper that authenticates all requests
// using HTTP Header Authentication with the provided token. In Github this is
// usually generated in Github.com under Settings>Developer settings>Personal
// access tokens.
type GithuHeaderTransport struct {
	Token string // Github personal access token

	// Transport is the underlying HTTP transport to use when making requests.
	// It will default to http.DefaultTransport if nil
	Transport http.RoundTripper
}

// RoundTrip implements the RoundTripper interface.
// https://github.com/google/go-github/blob/master/github/github.go#L949
func (g *GithuHeaderTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// To set extra headers, we must make a copy of the Request so
	// that we don't modify the Request we were given. This is required by the
	// specification of http.RoundTripper.
	//
	// Since we are going to modify only req.Header here, we only need a deep copy
	// of req.Header.
	req2 := new(http.Request)
	*req2 = *req
	req2.Header = make(http.Header, len(req.Header))
	for k, s := range req.Header {
		req2.Header[k] = append([]string(nil), s...)
	}
	req2.Header["Authorization"] = append([]string(nil), "token "+g.Token)

	return g.transport().RoundTrip(req2)
}

// Client returns an *http.Client that makes requests that are authenticated
// using HTTP Header Authentication.
func (g *GithuHeaderTransport) Client() *http.Client {
	return &http.Client{Transport: g}
}

func (g *GithuHeaderTransport) transport() http.RoundTripper {
	if g.Transport != nil {
		return g.Transport
	}
	return http.DefaultTransport
}
