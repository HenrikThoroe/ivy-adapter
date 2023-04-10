package mgmt

import (
	"errors"
	"io"
	"net/http"
	"os"

	"github.com/HenrikThoroe/ivy-adapter/internal/pkg/conf"
)

// DownloadEngine downloads the given engine instance and saves it to the engine store.
func DownloadEngine(inst *EngineInstance) error {
	path := inst.Path()
	tempPath := path + ".tmp"
	url := inst.URL()
	err := os.MkdirAll(conf.GetEngineStore(), os.ModePerm)

	if err != nil {
		return err
	}

	defer os.Remove(tempPath)

	if _, err := os.Stat(tempPath); err == nil {
		os.Remove(tempPath)
	}

	if _, err := os.Stat(path); err == nil {
		os.Remove(path)
	}

	if err := downloadFile(url, tempPath); err != nil {
		return err
	}

	os.Rename(tempPath, path)

	return nil
}

// downloadFile downloads a file from a given url and saves it to a given path.
func downloadFile(url string, path string) error {
	file, err := os.Create(path)
	resp, netErr := http.Get(url)

	defer file.Close()
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return errors.New("Could not download file: " + resp.Status)
	}

	if err != nil {
		return errors.New("Could not create file: " + err.Error())
	}

	if netErr != nil {
		return errors.New("Could not download file: " + netErr.Error())
	}

	_, ioErr := io.Copy(file, resp.Body)

	if ioErr != nil {
		return errors.New("Could not write to file: " + ioErr.Error())
	}

	return nil
}
