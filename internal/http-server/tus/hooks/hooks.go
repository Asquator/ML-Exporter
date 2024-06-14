package hooks

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/tus/tusd/v2/pkg/handler"
)

func PreuploadHook(hook handler.HookEvent) (handler.HTTPResponse, handler.FileInfoChanges, error) {
	res := handler.HTTPResponse{}
	meta := hook.HTTPRequest.Header.Get("Upload-Metadata")
	fmt.Println("metar", handler.ParseMetadataHeader(meta))

	fname, ok := handler.ParseMetadataHeader(meta)["filename"]

	fmt.Println("name", fname)
	if !ok {
		// set reject upload to true
		res.StatusCode = 400
		res.Body = "no filename provided"
		return res, handler.FileInfoChanges{}, nil
	}

	fileInfo := handler.FileInfoChanges{
		ID: fname,
	}

	return res, fileInfo, nil
}

func CompleteUploadHook(hook handler.HookEvent, path string) error {
	//TODO: ERORR HANDLING
	id := hook.Upload.ID

	infoPath := filepath.Join(path, id+".info")
	lockPath := filepath.Join(path, id+".lock")

	os.Remove(infoPath)
	os.ReadDir(lockPath)

	return nil
}
