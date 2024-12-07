package filetools

import (
	"bytes"
	"net/http"
)

func DownloadFile(localFilePath string, downloadFileURL string) error {

	// download
	res, err := http.Get(downloadFileURL)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	var buf bytes.Buffer
	if err := res.Write(&buf); err != nil {
		panic(err)
	}

	// write to file
	return Write(localFilePath, buf.Bytes())

}
