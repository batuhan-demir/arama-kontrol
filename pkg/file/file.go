package file

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func Download(url string, output string) {

	dir := "./files"

	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	out, err := os.Create(fmt.Sprintf("%s/%s", dir, output))
	if err != nil {
		panic(err)
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		panic(err)
	}
}
