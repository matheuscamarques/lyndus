package utils

import (
	"bufio"
	"encoding/base64"
	"errors"
	"image"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// FileLoad function for to sload file
func FileLoad(path string) (file *os.File, err error) {

	////log.Printf("-->Carregando Arquivo Path: %v", path)
	_, err = os.Stat(path)
	if err == nil {
		file, err = os.Open(path)
	}
	return file, err
}

// TrimToNum function remove non numeric
func TrimToNum(r rune) bool {
	return r < '0' || '9' < r
}

// GetPathFile function for to generate path location of path
func GetPathFile(ID, IDFile int, configPath, modulePath, prefix, extension string) string {
	pathFile := "/"
	strID := strconv.Itoa(ID)
	for _, c := range strID {
		pathFile += string(c) + "/"
	}
	path := configPath + "/" + modulePath + pathFile + strID + "/"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}
	return filepath.FromSlash(path + prefix + strconv.Itoa(IDFile) + "." + extension)
}

// GetPathFile function for to generate path location of path
func GetBMPathFile(ID int, IDFile string, configPath, modulePath, prefix, extension string) string {
	pathFile := "/"
	strID := strconv.Itoa(ID)
	for _, c := range strID {
		pathFile += string(c) + "/"
	}
	path := configPath + "/" + modulePath + pathFile
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}
	return filepath.FromSlash(path + prefix + IDFile + "." + extension)
}

// PathImageToBase64Str converte imagem para string base64
func PathImageToBase64Str(path string) (imgBase64Str string, err error) {
	file, err := FileLoad(path)
	if err != nil {
		return imgBase64Str, err
	}
	fi, err := file.Stat()
	if err != nil {
		return imgBase64Str, err
	}

	var size int64 = fi.Size()
	buf := make([]byte, size)

	fReader := bufio.NewReader(file)
	fReader.Read(buf)
	imgBase64Str = base64.StdEncoding.EncodeToString(buf)

	return "data:image/png;base64," + imgBase64Str, err
}

// Base64ImageSave converte base64 string de imagem para byte.
func Base64ImageSave(strImage, filePath string) error {
	var err error
	if len(strImage) < 10 {
		err = errors.New("Imagem inválida . ")
		return err
	}

	dataImage := strings.SplitAfterN(strImage, ",", 2)
	if len(dataImage) == 2 {
		strImage = dataImage[1]
	} else {
		strImage = dataImage[0]
	}
	//typeContent := dataImage[0][:len(dataImage[0])-8]

	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(strImage))

	var img image.Image

	img, _, err = image.Decode(reader)

	if err != nil {
		return err
	}

	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	err = png.Encode(f, img)
	return err
}
