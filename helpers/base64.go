// Di dalam package helpers
package helpers

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/png"
	"log"
)

func Base64toPng(imageBase64 string) ([]byte, error) {
	reader := base64.NewDecoder(base64.StdEncoding, bytes.NewReader([]byte(imageBase64)))
	m, _, err := image.Decode(reader)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	var imgByte []byte
	buffer := &bytes.Buffer{}
	err = png.Encode(buffer, m)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	imgByte = buffer.Bytes()

	return imgByte, nil
}
