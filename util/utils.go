package util

import (
	"io"
	"net/http"
	"os"
	"strings"
)

func Fetch(url string, path string) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	if response.StatusCode != 200 {
		return Errorf("%s: %s", response.Status, url)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if strings.HasPrefix(url, MASTER_URL) {
		return os.WriteFile(path, body, 0644)
	} else {
		return os.WriteFile(path, body, 0777)
	}
}
