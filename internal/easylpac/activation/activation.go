package activation

import (
	"bytes"
	"errors"
	"image"
	_ "image/jpeg"
	"os"
	"strings"

	"EasyLPAC/internal/easylpac/model"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"golang.design/x/clipboard"
)

var ErrInvalidCodeFormat = errors.New("invalid activation code format")

func DecodeCode(code string) (info model.PullInfo, confirmCodeNeeded bool, err error) {
	code = strings.TrimSpace(code)
	var ok bool
	if code, ok = strings.CutPrefix(code, "LPA:"); !ok {
		return info, false, ErrInvalidCodeFormat
	}

	// ref: https://www.gsma.com/esim/wp-content/uploads/2020/06/SGP.22-v2.2.2.pdf#page=111
	parts := strings.Split(code, "$")
	if len(parts) == 0 || parts[0] != "1" {
		return info, false, ErrInvalidCodeFormat
	}

	var codeNeeded string
	bindings := []*string{&info.SMDP, &info.MatchID, &info.ObjectID, &codeNeeded}
	for index, value := range parts[1:] {
		if index >= len(bindings) {
			break
		}
		*bindings[index] = strings.TrimSpace(value)
	}
	if info.SMDP == "" {
		return info, false, ErrInvalidCodeFormat
	}
	confirmCodeNeeded = codeNeeded == "1"
	return info, confirmCodeNeeded, nil
}

func CompleteCode(input string) string {
	if strings.HasPrefix(input, "LPA:1$") {
		return input
	}
	if strings.HasPrefix(input, "1$") {
		return "LPA:" + input
	}
	if strings.HasPrefix(input, "$") {
		return "LPA:1" + input
	}
	return input
}

func DecodeQRCodeFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer func() {
		_ = file.Close()
	}()

	img, _, err := image.Decode(file)
	if err != nil {
		return "", err
	}
	return decodeQRCode(img)
}

func DecodeQRCodeBytes(imageBytes []byte) (string, error) {
	img, _, err := image.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return "", err
	}
	return decodeQRCode(img)
}

func ReadClipboard() (clipboard.Format, []byte, error) {
	if err := clipboard.Init(); err != nil {
		panic(err)
	}

	result := clipboard.Read(clipboard.FmtText)
	if len(result) != 0 {
		return clipboard.FmtText, result, nil
	}

	result = clipboard.Read(clipboard.FmtImage)
	if len(result) != 0 {
		return clipboard.FmtImage, result, nil
	}

	return clipboard.FmtText, nil, errors.New("failed to read clipboard: not text or image")
}

func decodeQRCode(img image.Image) (string, error) {
	bitmap, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		return "", err
	}

	result, err := qrcode.NewQRCodeReader().Decode(bitmap, nil)
	if err != nil {
		return "", err
	}
	return result.String(), nil
}
