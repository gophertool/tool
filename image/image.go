// Package image 提供图片处理相关的工具函数
package image

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"strings"
)

// Loader 是图片加载器接口，提供从不同来源加载图片的方法
type Loader interface {
	// LoadFromFile 从文件加载图片
	LoadFromFile(filePath string) (image.Image, error)

	// LoadFromURL 从URL加载图片
	LoadFromURL(url string) (image.Image, error)

	// LoadFromBase64 从Base64字符串加载图片
	LoadFromBase64(base64Str string) (image.Image, error)

	// LoadFromBytes 从字节数组加载图片
	LoadFromBytes(data []byte) (image.Image, error)

	// LoadFromReader 从io.Reader加载图片
	LoadFromReader(reader io.Reader) (image.Image, error)
}

// DefaultLoader 是默认的图片加载器实现
type DefaultLoader struct{}

// NewLoader 创建一个新的默认图片加载器
func NewLoader() Loader {
	return &DefaultLoader{}
}

// LoadFromFile 从文件加载图片
func (l *DefaultLoader) LoadFromFile(filePath string) (image.Image, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开图片文件失败: %w", err)
	}
	defer file.Close()

	return l.LoadFromReader(file)
}

// LoadFromURL 从URL加载图片
func (l *DefaultLoader) LoadFromURL(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("获取URL图片失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("获取URL图片失败，状态码: %d", resp.StatusCode)
	}

	return l.LoadFromReader(resp.Body)
}

// LoadFromBase64 从Base64字符串加载图片
func (l *DefaultLoader) LoadFromBase64(base64Str string) (image.Image, error) {
	// 移除可能的前缀，如 "data:image/jpeg;base64,"
	commaIndex := strings.Index(base64Str, ",")
	if commaIndex != -1 {
		base64Str = base64Str[commaIndex+1:]
	}

	// 尝试解码Base64字符串
	data, err := base64.StdEncoding.DecodeString(base64Str)
	if err != nil {
		// 如果标准编码失败，尝试URL安全编码
		data, err = base64.URLEncoding.DecodeString(base64Str)
		if err != nil {
			return nil, fmt.Errorf("解码Base64字符串失败: %w", err)
		}
	}

	return l.LoadFromBytes(data)
}

// LoadFromBytes 从字节数组加载图片
func (l *DefaultLoader) LoadFromBytes(data []byte) (image.Image, error) {
	reader := bytes.NewReader(data)
	return l.LoadFromReader(reader)
}

// LoadFromReader 从io.Reader加载图片
func (l *DefaultLoader) LoadFromReader(reader io.Reader) (image.Image, error) {
	img, format, err := image.Decode(reader)
	if err != nil {
		return nil, fmt.Errorf("解码图片失败: %w", err)
	}

	// 可以根据format做一些额外处理，这里只是简单返回解码后的图片
	_ = format
	return img, nil
}

// 支持的图片格式
var (
	ErrUnsupportedFormat = errors.New("不支持的图片格式")
)

// SaveImage 保存图片到文件
func SaveImage(img image.Image, filePath string, format string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("创建图片文件失败: %w", err)
	}
	defer file.Close()

	return SaveImageToWriter(img, file, format)
}

// SaveImageToWriter 保存图片到io.Writer
func SaveImageToWriter(img image.Image, writer io.Writer, format string) error {
	switch strings.ToLower(format) {
	case "jpeg", "jpg":
		return jpeg.Encode(writer, img, &jpeg.Options{Quality: 90})
	case "png":
		return png.Encode(writer, img)
	default:
		return ErrUnsupportedFormat
	}
}

// GetImageFormat 获取图片格式
func GetImageFormat(data []byte) (string, error) {
	_, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return "", fmt.Errorf("获取图片格式失败: %w", err)
	}
	return format, nil
}
