package auth

import (
	"github.com/bwmarrin/snowflake"
	routing "github.com/qiangxue/fasthttp-routing"
)

type Middleware interface {
	CreateHash(pass string) (string, error)
	CheckHash(hash, pass string) bool
	CreateSessionKey() (string, error)
	CreateSession(ctx *routing.Context, uid snowflake.ID, remember bool) (string, error)
	Login(ctx *routing.Context) bool
	CheckRequestAuth(ctx *routing.Context) error
	LogOut(ctx *routing.Context) error
}
