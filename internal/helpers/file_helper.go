package helpers

import (
	"autotm-admin/internal/configs"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func CreateFolder(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

func DeleteFile(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}
	return os.RemoveAll(path)
}

func ConvertToWebP(srcPath, dstPath string, quality float32) error {
	img, err := imaging.Open(srcPath, imaging.AutoOrientation(true))
	if err != nil {
		return fmt.Errorf("fail to open image: %v", err)
	}

	outFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("failed to create webp file: %v", err)
	}
	defer outFile.Close()

	options, err := encoder.NewLossyEncoderOptions(encoder.PresetPhoto, quality)
	if err != nil {
		return fmt.Errorf("fail to create lossy encoder: %v", err)
	}

	if err := webp.Encode(outFile, img, options); err != nil {
		return fmt.Errorf("fail to convert to webp: %v", err)
	}
	return nil
}

func UploadImage(file multipart.File, header *multipart.FileHeader) (string, error) {
	cfg := configs.GetConfig()
	publicPath := cfg.FilePath

	ext := strings.ToLower(filepath.Ext(header.Filename))
	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	dirPath := filepath.Join(publicPath, "images")

	if err := CreateFolder(dirPath); err != nil {
		return "", fmt.Errorf("fail to create folder: %v", err)
	}

	fullPath := filepath.Join(dirPath, fileName)
	dst, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("fail to create file: %v", err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return "", fmt.Errorf("fail to copy file: %v", err)
	}

	if ext != ".webp" {
		webpFilename := strings.TrimSuffix(fileName, ext) + ".webp"
		webpFullPath := filepath.Join(dirPath, webpFilename)

		if err := ConvertToWebP(fullPath, webpFullPath, 90); err != nil {
			return "", fmt.Errorf("fail to convert to webp: %v", err)
		}

		if err := DeleteFile(fullPath); err != nil {
			fmt.Printf("Warning: failed to delete original file: %v\n", err)
		}

		return "/images/" + webpFilename, nil
	}

	return "/images/" + fileName, nil
}

func DeleteImage(relativePath string) error {
	cfg := configs.GetConfig()
	publicPath := cfg.FilePath

	cleanPath := strings.TrimPrefix(relativePath, "/")

	fullPath := filepath.Join(publicPath, strings.TrimPrefix(cleanPath, "uploads/"))

	if err := DeleteFile(fullPath); err != nil {
		return fmt.Errorf("failed to delete image: %v", err)
	}

	return nil
}
