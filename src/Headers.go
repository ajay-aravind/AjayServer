package main

const (
	HeaderAccept             = "Accept"
	HeaderAcceptCharset      = "Accept-Charset"
	HeaderAcceptEncoding     = "Accept-Encoding"
	HeaderAcceptLanguage     = "Accept-Language"
	HeaderAuthorization      = "Authorization"
	HeaderCacheControl       = "Cache-Control"
	HeaderConnection         = "Connection"
	HeaderContentEncoding    = "Content-Encoding"
	HeaderContentLanguage    = "Content-Language"
	HeaderContentLength      = "Content-Length"
	HeaderContentLocation    = "Content-Location"
	HeaderContentType        = "Content-Type"
	HeaderCookie             = "Cookie"
	HeaderDate               = "Date"
	HeaderETag               = "ETag"
	HeaderExpect             = "Expect"
	HeaderExpires            = "Expires"
	HeaderForwarded          = "Forwarded"
	HeaderFrom               = "From"
	HeaderHost               = "Host"
	HeaderIfMatch            = "If-Match"
	HeaderIfModifiedSince    = "If-Modified-Since"
	HeaderIfNoneMatch        = "If-None-Match"
	HeaderIfRange            = "If-Range"
	HeaderIfUnmodifiedSince  = "If-Unmodified-Since"
	HeaderLastModified       = "Last-Modified"
	HeaderLocation           = "Location"
	HeaderOrigin             = "Origin"
	HeaderPragma             = "Pragma"
	HeaderProxyAuthenticate  = "Proxy-Authenticate"
	HeaderProxyAuthorization = "Proxy-Authorization"
	HeaderRange              = "Range"
	HeaderReferer            = "Referer"
	HeaderRetryAfter         = "Retry-After"
	HeaderServer             = "Server"
	HeaderSetCookie          = "Set-Cookie"
	HeaderTE                 = "TE"
	HeaderTrailer            = "Trailer"
	HeaderTransferEncoding   = "Transfer-Encoding"
	HeaderUserAgent          = "User-Agent"
	HeaderVary               = "Vary"
	HeaderWWWAuthenticate    = "WWW-Authenticate"
	HeaderXForwardedFor      = "X-Forwarded-For"
)

/*+Add(key string, value string) //Add adds the given 'key: value' header.
+Set(key string, value string) // adds/updates the give header
+SetConnectionClose() //sets 'Connection: close' header
+GetConnectionClose() bool // gets if 'Connection: close' header is present
+GetConnectionUpgrade() bool //ConnectionUpgrade returns true if 'Connection: Upgrade' header is set.
+SetConnectionUpgrade()
+SetContentEncoding()
+GetContentEncoding()
+SetContentLength() //Negative content-length sets 'Transfer-Encoding: chunked' header.
+setContentType()
+SetCookie()
+CopyTo()
+del(key string) //delets header with given key
+DelAllCookies()
+DelCookie(key string)
+Reset() //resets headers, deletes all header keys
+Peek()
+PeekAll()
+PeekCookie()
+Get/Set Protocol()
+GetServerHeader()
+SetCanonical(key, value) // making key in canonical form, but do i need it?
+SetLastModified()
+SetNoDefaultContentType() // if set to true, what do i set? Do i expose a const and set it??
+SetProtocol()
+SetStatusCode()
+SetStatusMessage()
+AddTrailer() // Add Trailed RequestHeaders, this is only supported in  chunked transfer encoding
+IsTls()
+GetLocalAddr
+GetLocalIp()
+GetRequestPath()
+GetContentType()
+GetUserAgent()*/

type Headers struct {
	headersList map[string][]string
}
