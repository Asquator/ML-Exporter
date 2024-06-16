package hooks

import (
	"goml/internal/config"
	"os"
	"path/filepath"

	"github.com/tus/tusd/v2/pkg/handler"
)

func PreuploadHook(hook handler.HookEvent) (handler.HTTPResponse, handler.FileInfoChanges, error) {
	res := handler.HTTPResponse{}
	meta := hook.HTTPRequest.Header.Get("Upload-Metadata")

	fname, ok := handler.ParseMetadataHeader(meta)["filename"]

	if !ok {
		// set reject upload to true
		res.StatusCode = 400
		res.Body = "no filename provided"
		return res, handler.FileInfoChanges{}, nil
	}

	fileInfo := handler.FileInfoChanges{
		ID:       fname,
		MetaData: handler.MetaData{},
	}

	return res, fileInfo, nil
}

func CompleteUploadHook(hook handler.HookEvent, cfg *config.Config) error {
	//TODO: ERORR HANDLING
	id := hook.Upload.ID

	infoPath := filepath.Join(cfg.Tus.UploadPath, id+".info")
	lockPath := filepath.Join(cfg.Tus.UploadPath, id+".lock")

	os.Remove(infoPath)
	os.ReadDir(lockPath)

	//TODO: VALIDATE FILE
	os.Rename(filepath.Join(cfg.Tus.UploadPath, id), filepath.Join(cfg.StoragePath, id))

	return nil

}
