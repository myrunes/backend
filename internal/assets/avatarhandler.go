package assets

import (
	"fmt"
	"io"
	"net/http"

	"github.com/myrunes/backend/internal/logger"
	"github.com/myrunes/backend/internal/storage"
	"github.com/myrunes/backend/pkg/workerpool"
)

const (
	avatarCDNURL     = "https://www.mobafire.com/images/avatars/%s-classic.png"
	avatarBucketName = "championavatars"
	avatarMimeType   = "image/png"
)

type AvatarHandler struct {
	storage storage.Middleware
}

func NewAvatarHandler(st storage.Middleware) *AvatarHandler {
	return &AvatarHandler{st}
}

func (ah *AvatarHandler) Get(champ string) (io.ReadCloser, int64, error) {
	return ah.storage.GetObject(avatarBucketName, getObjectName(champ))
}

func (ah *AvatarHandler) FetchAll(cChampNames chan string, cError chan error) {
	wp := workerpool.New(5)

	go func() {
		for res := range wp.Results() {
			if err, _ := res.(error); err != nil {
				cError <- err
			}
		}
	}()

	for champ := range cChampNames {
		wp.Push(ah.jobFetchSingle, champ)
	}
	wp.Close()

	wp.WaitBlocking()
	close(cError)
}

func (ah *AvatarHandler) jobFetchSingle(workerId int, params ...interface{}) interface{} {
	champ := params[0].(string)

	logger.Info("ASSETSHANDLER :: [%d] fetch champion avatar asset of '%s'...", workerId, champ)

	url := fmt.Sprintf(avatarCDNURL, champ)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		return fmt.Errorf("resuest failed with code %d", resp.StatusCode)
	}

	return ah.put(champ, resp.Body, resp.ContentLength)
}

func (ah *AvatarHandler) put(champ string, reader io.Reader, size int64) error {
	return ah.storage.PutObject(avatarBucketName, getObjectName(champ), reader, size, avatarMimeType)
}

func getObjectName(champ string) string {
	return fmt.Sprintf("%s.png", champ)
}
