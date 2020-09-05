package assets

import (
	"fmt"
	"io"
	"net/http"

	"github.com/myrunes/backend/internal/storage"
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
	for champ := range cChampNames {
		url := fmt.Sprintf(avatarCDNURL, champ)
		resp, err := http.Get(url)
		if err != nil {
			cError <- err
			continue
		}

		if resp.StatusCode >= 400 {
			cError <- fmt.Errorf("resuest failed with code %d", resp.StatusCode)
			continue
		}

		err = ah.put(champ, resp.Body, resp.ContentLength)
		if err != nil {
			cError <- err
			continue
		}
	}

	close(cError)
}

func (ah *AvatarHandler) put(champ string, reader io.Reader, size int64) error {
	return ah.storage.PutObject(avatarBucketName, getObjectName(champ), reader, size, avatarMimeType)
}

func getObjectName(champ string) string {
	return fmt.Sprintf("%s.png", champ)
}
