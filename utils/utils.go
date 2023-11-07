package utils

import (
	"fmt"
	"image"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"
)

// PrintCmdErr print error in console
func PrintCmdErr(err error) {
	_, err = fmt.Fprintf(os.Stderr, "Error: '%s' \n", err)
	if err != nil {
		panic(err)
	}
}

func GenExpireTime() time.Duration {
	min := 2
	max := 5
	return time.Duration(rand.Intn(max-min)+min) * time.Minute
}

func ParseBaseUri(url string) string {
	reg, err := regexp.Compile("(https?://[^/]*)")
	if err != nil {
		// print log
		return ""
	}
	subs := reg.FindStringSubmatch(url)

	if len(subs) > 1 {
		return subs[1]
	}
	return ""
}

func BuildUrl(baseUri string, path string) string {
	return strings.TrimSuffix(ParseBaseUri(baseUri), "/") + "/" + strings.TrimPrefix(path, "/")
}

func ValidJpgImage(jpgPath string) (valid bool, err error) {
	file, err := os.Open(jpgPath)
	if err != nil {
		return
	}
	defer file.Close()
	//image.Decode 函数用于解码图像文件，并返回一个 image.Image 接口类型的对象，代表解码后的图像。这个函数会自动识别图像的格式，并根据格式进行解码。使用 image.Decode 函数可以获取完整的图像数据，可以对图像进行处理、修改和保存等操作。
	//image.DecodeConfig 函数用于获取图像文件的基本信息，而不需要完全解码图像。它返回一个 image.Config 类型的对象，包含图像的宽度、高度、颜色模式等信息，但不包含图像的像素数据。使用 image.DecodeConfig 函数可以快速获取图像的基本信息，而无需完全解码图像。
	_, format, err := image.DecodeConfig(file)
	if err != nil {
		return
	}

	if format != "jpeg" {
		return false, err
	}

	return true, err
}
