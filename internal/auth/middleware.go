package auth

import (
	"github.com/qiangxue/fasthttp-routing"
)

type Middleware interface {
	CreateHash(pass []byte) ([]byte, error)
	CheckHash(hash, pass []byte) bool
	CreateSessionKey() (string, error)
	Login(ctx *routing.Context) bool
	CheckRequestAuth(ctx *routing.Context) error
}
