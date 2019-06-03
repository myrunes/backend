package auth

import (
	"github.com/bwmarrin/snowflake"
	"github.com/qiangxue/fasthttp-routing"
)

type Middleware interface {
	CreateHash(pass []byte) ([]byte, error)
	CheckHash(hash, pass []byte) bool
	CreateSessionKey() (string, error)
	CreateSession(ctx *routing.Context, uid snowflake.ID, remember bool) error
	Login(ctx *routing.Context) bool
	CheckRequestAuth(ctx *routing.Context) error
	LogOut(ctx *routing.Context) error
}
