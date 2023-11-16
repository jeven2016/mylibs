package utils

import (
	"fmt"
	"log"
	"net/url"
	"path/filepath"
	"testing"
)

func TestPicValidation(t *testing.T) {
	path := "/mnt/files/comic/cartoon18/hello/cover.jpg"
	valid, err := ValidJpgImage(path)
	if err != nil {
		t.Fatal(err)
	}
	log.Println(valid)

}

func TestFileExt(t *testing.T) {
	testurl := "http://test.example/abc.png?a=123&b=123"
	parsedURL, err := url.Parse(testurl)
	if err != nil {
		fmt.Println("无法解析 URL:", err)
		return
	}
	filename := filepath.Base(parsedURL.Path)
	ext := filepath.Ext(filename)
	if ext != ".png" {
		t.Fatal("not a valid file extension")
	}
}
