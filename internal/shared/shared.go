package shared

import routing "github.com/qiangxue/fasthttp-routing"

var (
	headerXForwardedFor = []byte("X-Forwarded-For")
)

func GetIPAddr(ctx *routing.Context) string {
	forwardedfor := ctx.Request.Header.PeekBytes(headerXForwardedFor)
	if forwardedfor != nil && len(forwardedfor) > 0 {
		return string(forwardedfor)
	}

	return ctx.RemoteIP().String()
}
