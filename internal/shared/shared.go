package shared

import routing "github.com/qiangxue/fasthttp-routing"

var (
	headerXForwardedFor = []byte("X-Forwarded-For")
)

// GetIPAddr returns the IP address as string
// from the given request context. When the
// request contains a 'X-Forwarded-For' header,
// the value of this will be returned as address.
// Else, the conext remote address will be returned.
func GetIPAddr(ctx *routing.Context) string {
	forwardedfor := ctx.Request.Header.PeekBytes(headerXForwardedFor)
	if forwardedfor != nil && len(forwardedfor) > 0 {
		return string(forwardedfor)
	}

	return ctx.RemoteIP().String()
}
