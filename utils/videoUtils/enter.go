package videoUtils

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

// SaveUploadedFile 保存上传的文件到指定目录
func SaveUploadedFile(file *multipart.FileHeader, dir string) (string, error) {
	// 创建目录
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", errors.New("创建目录失败")
	}

	// 打开上传文件
	src, err := file.Open()
	if err != nil {
		return "", errors.New("打开上传文件失败")
	}
	defer src.Close()

	// 创建目标文件
	filename := generateFilename(file.Filename)
	dstPath := filepath.Join(dir, filename)
	dst, err := os.Create(dstPath)
	if err != nil {
		return "", errors.New("创建目标文件失败")
	}
	defer dst.Close()

	// 复制文件内容
	if _, err = io.Copy(dst, src); err != nil {
		return "", errors.New("保存文件内容失败")
	}

	return dstPath, nil
}

// generateFilename 生成唯一的文件名
func generateFilename(original string) string {
	ext := filepath.Ext(original)
	return time.Now().Format("20060102150405") + ext
}
