package cookiejar

import (
	"net/http"
	"net/url"
	"strings"
	"time"
)

type CookieJar struct {
	cookies map[string]*http.Cookie
	created map[string]time.Time
}

func NewCookieJar() (jar *CookieJar) {
	return &CookieJar{
		cookies: make(map[string]*http.Cookie),
		created: make(map[string]time.Time),
	}
}

func (jar *CookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	for _, cookie := range cookies {
		key := cookie.String()
		jar.cookies[key] = cookie
		jar.created[key] = time.Now()
	}
}

func (jar *CookieJar) Prune() {
	for key, cookie := range jar.cookies {
		if !cookie.Expires.IsZero() && cookie.Expires.Before(time.Now()) {
			delete(jar.cookies, key)
			delete(jar.created, key)
		}

		if cookie.MaxAge > 0 {
			created := jar.created[key]
			expires := created.Add(time.Duration(cookie.MaxAge) * time.Second)

			if expires.Before(time.Now()) {
				delete(jar.cookies, key)
				delete(jar.created, key)
			}
		}
	}
}

func (jar *CookieJar) Cookies(u *url.URL) (cookies []*http.Cookie) {
	for _, cookie := range jar.cookies {
		if cookie.Domain != "" && !strings.HasSuffix(u.Host, cookie.Domain) {
			continue
		}

		if cookie.Path != "" && !strings.HasPrefix(u.Path, cookie.Path) {
			continue
		}

		if cookie.Secure && u.Scheme != "https" {
			continue
		}

		cookies = append(cookies, cookie)
	}

	return cookies
}
