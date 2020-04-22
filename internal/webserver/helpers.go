package webserver

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/myrunes/backend/internal/static"

	routing "github.com/qiangxue/fasthttp-routing"
	"github.com/valyala/fasthttp"
)

var emptyResponseBody = []byte("{}")

var (
	headerUserAgent    = []byte("User-Agent")
	headerCacheControl = []byte("Cache-Control")
	headerETag         = []byte("ETag")

	headerCacheControlValue = []byte("max-age=2592000; must-revalidate; proxy-revalidate;  public")

	bcryptPrefix = []byte("$2a")
)

var defStatusBoddies = map[int][]byte{
	http.StatusOK:           []byte("{\n  \"code\": 200,\n  \"message\": \"ok\"\n}"),
	http.StatusCreated:      []byte("{\n  \"code\": 201,\n  \"message\": \"created\"\n}"),
	http.StatusNotFound:     []byte("{\n  \"code\": 404,\n  \"message\": \"not found\"\n}"),
	http.StatusUnauthorized: []byte("{\n  \"code\": 401,\n  \"message\": \"unauthorized\"\n}"),
}

// jsonError writes the error message of err and the
// passed status to response context and aborts the
// execution of following registered handlers ONLY IF
// err != nil.
// This function always returns a nil error that the
// default error handler can be bypassed.
func jsonError(ctx *routing.Context, err error, status int) error {
	if err != nil {
		ctx.Response.Header.SetContentType("application/json")
		ctx.SetStatusCode(status)
		ctx.SetBodyString(fmt.Sprintf("{\n  \"code\": %d,\n  \"message\": \"%s\"\n}",
			status, err.Error()))
		ctx.Abort()
	}
	return nil
}

// jsonResponse tries to parse the passed interface v
// to JSON and writes it to the response context body
// as same as the passed status code.
// If the parsing fails, this will result in a jsonError
// output of the error with status 500.
// This function always returns a nil error.
func jsonResponse(ctx *routing.Context, v interface{}, status int) error {
	var err error
	data := emptyResponseBody

	if v == nil {
		if d, ok := defStatusBoddies[status]; ok {
			data = d
		}
	} else {
		if static.Release != "TRUE" {
			data, err = json.MarshalIndent(v, "", "  ")
		} else {
			data, err = json.Marshal(v)
		}
		if err != nil {
			return jsonError(ctx, err, fasthttp.StatusInternalServerError)
		}
	}

	ctx.Response.Header.SetContentType("application/json")
	ctx.SetStatusCode(status)
	_, err = ctx.Write(data)

	return jsonError(ctx, err, fasthttp.StatusInternalServerError)
}

func jsonCachableResponse(ctx *routing.Context, v interface{}, status int) error {
	var err error
	data := emptyResponseBody

	if v == nil {
		if d, ok := defStatusBoddies[status]; ok {
			data = d
		}
	} else {
		if static.Release != "TRUE" {
			data, err = json.MarshalIndent(v, "", "  ")
		} else {
			data, err = json.Marshal(v)
		}
		if err != nil {
			return jsonError(ctx, err, fasthttp.StatusInternalServerError)
		}
	}

	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.Header.SetBytesKV(headerCacheControl, headerCacheControlValue)
	ctx.Response.Header.SetBytesK(headerETag, getETag(data, true))
	ctx.SetStatusCode(status)
	_, err = ctx.Write(data)

	return jsonError(ctx, err, fasthttp.StatusInternalServerError)
}

// parseJSONBody tries to parse a requests JSON
// body to the passed object pointer. If the
// parsing fails, this will result in a jsonError
// output with status 400.
// This function always returns a nil error.
func parseJSONBody(ctx *routing.Context, v interface{}) error {
	data := ctx.PostBody()
	err := json.Unmarshal(data, v)
	if err != nil {
		jsonError(ctx, err, fasthttp.StatusBadRequest)
	}
	return err
}

func (ws *WebServer) addHeaders(ctx *routing.Context) error {
	ctx.Response.Header.SetServer("MYRUNES v." + static.AppVersion)

	if ws.config.PublicAddr != "" && ws.config.EnableCors {
		ctx.Response.Header.Set("Access-Control-Allow-Origin", ws.config.PublicAddr)
		ctx.Response.Header.Set("Access-Control-Allow-Headers", "authorization, content-type, set-cookie, cookie, server")
		ctx.Response.Header.Set("Access-Control-Allow-Methods", "POST, GET, DELETE, OPTIONS")
		ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
	}

	return nil
}

func checkPageName(pageName, guess string, tollerance float64) bool {
	if pageName == "" || guess == "" {
		return false
	}

	lenPageName := float64(len(strings.Replace(pageName, " ", "", -1)))
	lenGuesses := float64(len(strings.Replace(guess, " ", "", -1)))

	pageNameSplit := strings.Split(strings.ToLower(pageName), " ")
	guessSplit := strings.Split(strings.ToLower(guess), " ")

	var matchedChars int
	for _, wordName := range pageNameSplit {
		for _, guessName := range guessSplit {
			if wordName == guessName {
				matchedChars += len(wordName)
			}
		}
	}

	return float64(matchedChars)/lenPageName >= (1-tollerance) &&
		float64(matchedChars)/lenGuesses >= (1-tollerance)
}

func getETag(body []byte, weak bool) string {
	hash := sha1.Sum(body)

	weakTag := ""
	if weak {
		weakTag = "W/"
	}

	tag := fmt.Sprintf("%s\"%x\"", weakTag, hash)

	return tag
}

func isOldPasswordHash(hash []byte) bool {
	return bytes.HasPrefix(hash, bcryptPrefix)
}
