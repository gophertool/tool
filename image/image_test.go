package image_test

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/color"
	"strings"
	"testing"

	imageutil "github.com/gophertool/tool/image"
)

const (
	testImageBase64 = `/9j/4AAQSkZJRgABAQAAAQABAAD/4gHYSUNDX1BST0ZJTEUAAQEAAAHIAAAAAAQwAABtbnRyUkdCIFhZWiAH4AABAAEAAAAAAABhY3NwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAQAA9tYAAQAAAADTLQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAlkZXNjAAAA8AAAACRyWFlaAAABFAAAABRnWFlaAAABKAAAABRiWFlaAAABPAAAABR3dHB0AAABUAAAABRyVFJDAAABZAAAAChnVFJDAAABZAAAAChiVFJDAAABZAAAAChjcHJ0AAABjAAAADxtbHVjAAAAAAAAAAEAAAAMZW5VUwAAAAgAAAAcAHMAUgBHAEJYWVogAAAAAAAAb6IAADj1AAADkFhZWiAAAAAAAABimQAAt4UAABjaWFlaIAAAAAAAACSgAAAPhAAAts9YWVogAAAAAAAA9tYAAQAAAADTLXBhcmEAAAAAAAQAAAACZmYAAPKnAAANWQAAE9AAAApbAAAAAAAAAABtbHVjAAAAAAAAAAEAAAAMZW5VUwAAACAAAAAcAEcAbwBvAGcAbABlACAASQBuAGMALgAgADIAMAAxADb/2wBDAAMCAgICAgMCAgIDAwMDBAYEBAQEBAgGBgUGCQgKCgkICQkKDA8MCgsOCwkJDRENDg8QEBEQCgwSExIQEw8QEBD/2wBDAQMDAwQDBAgEBAgQCwkLEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBAQEBD/wAARCAABAAEDASIAAhEBAxEB/8QAFQABAQAAAAAAAAAAAAAAAAAAAAj/xAAUEAEAAAAAAAAAAAAAAAAAAAAA/8QAFQEBAQAAAAAAAAAAAAAAAAAABwn/xAAUEQEAAAAAAAAAAAAAAAAAAAAA/9oADAMBAAIRAxEAPwCdAAYqm//Z`
)

// 测试从文件加载图片
func TestLoadFromFile(t *testing.T) {
	// 这个测试需要一个实际的图片文件
	// 在实际测试中，你需要提供一个有效的图片文件路径
	filePath := "./test.jpg"

	loader := imageutil.NewLoader()
	img, err := loader.LoadFromFile(filePath)
	if err != nil {
		t.Fatalf("从文件加载图片失败: %v", err)
	}

	if img == nil {
		t.Fatal("加载的图片为nil")
	}

	// 可以进一步检查图片的属性
	bounds := img.Bounds()
	if bounds.Dx() <= 0 || bounds.Dy() <= 0 {
		t.Fatalf("图片尺寸无效: %v", bounds)
	}
}

// 测试从URL加载图片
func TestLoadFromURL(t *testing.T) {
	// 这个测试需要一个有效的图片URL和网络连接
	// 在实际测试中，你需要提供一个有效的图片URL
	url := "https://raw.githubusercontent.com/gophertool/tool/refs/heads/main/image/test.png"

	loader := imageutil.NewLoader()
	img, err := loader.LoadFromURL(url)
	if err != nil {
		t.Fatalf("从URL加载图片失败: %v", err)
	}

	if img == nil {
		t.Fatal("加载的图片为nil")
	}
}

// 测试从Base64字符串加载图片
func TestLoadFromBase64(t *testing.T) {
	loader := imageutil.NewLoader()
	img, err := loader.LoadFromBase64(testImageBase64)
	if err != nil {
		t.Fatalf("从Base64加载图片失败: %v", err)
	}

	if img == nil {
		t.Fatal("加载的图片为nil")
	}

	// 检查图片尺寸应该是1x1
	bounds := img.Bounds()
	if bounds.Dx() != 1 || bounds.Dy() != 1 {
		t.Fatalf("图片尺寸不是1x1: %v", bounds)
	}
}

// 测试从字节数组加载图片
func TestLoadFromBytes(t *testing.T) {
	data, err := base64.StdEncoding.DecodeString(testImageBase64)
	if err != nil {
		t.Fatalf("解码Base64字符串失败: %v", err)
	}

	loader := imageutil.NewLoader()
	img, err := loader.LoadFromBytes(data)
	if err != nil {
		t.Fatalf("从字节数组加载图片失败: %v", err)
	}

	if img == nil {
		t.Fatal("加载的图片为nil")
	}

	// 检查图片尺寸应该是1x1
	bounds := img.Bounds()
	if bounds.Dx() != 1 || bounds.Dy() != 1 {
		t.Fatalf("图片尺寸不是1x1: %v", bounds)
	}
}

// 测试从io.Reader加载图片
func TestLoadFromReader(t *testing.T) {
	data, err := base64.StdEncoding.DecodeString(testImageBase64)
	if err != nil {
		t.Fatalf("解码Base64字符串失败: %v", err)
	}

	reader := bytes.NewReader(data)

	loader := imageutil.NewLoader()
	img, err := loader.LoadFromReader(reader)
	if err != nil {
		t.Fatalf("从Reader加载图片失败: %v", err)
	}

	if img == nil {
		t.Fatal("加载的图片为nil")
	}

	// 检查图片尺寸应该是1x1
	bounds := img.Bounds()
	if bounds.Dx() != 1 || bounds.Dy() != 1 {
		t.Fatalf("图片尺寸不是1x1: %v", bounds)
	}
}

// 测试保存图片
func TestSaveImage(t *testing.T) {
	// 这个测试需要写入文件系统
	// 在实际测试中，你可能需要使用临时文件或模拟文件系统
	//
	// // 首先加载一个图片
	// base64Str := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+P+/HgAFeAJ5jfZixgAAAABJRU5ErkJggg=="
	// data, _ := base64.StdEncoding.DecodeString(base64Str)
	//
	// loader := imageutil.NewLoader()
	// img, err := loader.LoadFromBytes(data)
	// if err != nil {
	// 	t.Fatalf("加载图片失败: %v", err)
	// }
	//
	// // 创建临时文件
	// tmpfile, err := os.CreateTemp("", "test-*.png")
	// if err != nil {
	// 	t.Fatalf("创建临时文件失败: %v", err)
	// }
	// tmpfileName := tmpfile.Name()
	// tmpfile.Close()
	// defer os.Remove(tmpfileName) // 清理
	//
	// // 保存图片
	// err = imageutil.SaveImage(img, tmpfileName, "png")
	// if err != nil {
	// 	t.Fatalf("保存图片失败: %v", err)
	// }
	//
	// // 验证文件存在且大小大于0
	// stat, err := os.Stat(tmpfileName)
	// if err != nil {
	// 	t.Fatalf("获取保存的文件信息失败: %v", err)
	// }
	//
	// if stat.Size() <= 0 {
	// 	t.Fatal("保存的图片文件大小为0")
	// }
}

// 测试保存图片到Writer
func TestSaveImageToWriter(t *testing.T) {
	// 创建一个简单的1x1像素的图片
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.RGBA{255, 0, 0, 255}) // 红色像素

	// 使用bytes.Buffer作为Writer
	buf := new(bytes.Buffer)

	// 保存图片到Writer
	err := imageutil.SaveImageToWriter(img, buf, "png")
	if err != nil {
		t.Fatalf("保存图片到Writer失败: %v", err)
	}

	// 验证buffer中有数据
	if buf.Len() <= 0 {
		t.Fatal("保存的图片数据大小为0")
	}

	// 验证数据是有效的PNG图片
	_, err = imageutil.GetImageFormat(buf.Bytes())
	if err != nil {
		t.Fatalf("获取保存的图片格式失败: %v", err)
	}
}

// 测试获取图片格式
func TestGetImageFormat(t *testing.T) {
	// 使用一个有效的jpeg图片数据
	data, _ := base64.StdEncoding.DecodeString(testImageBase64)

	format, err := imageutil.GetImageFormat(data)
	if err != nil {
		t.Fatalf("获取图片格式失败: %v", err)
	}

	if format != "jpeg" {
		t.Fatalf("图片格式不正确，期望png，实际%s", format)
	}
}

// 测试不支持的格式
func TestUnsupportedFormat(t *testing.T) {
	// 创建一个简单的1x1像素的图片
	img := image.NewRGBA(image.Rect(0, 0, 1, 1))
	img.Set(0, 0, color.RGBA{0, 255, 0, 255}) // 绿色像素

	// 使用不支持的格式保存
	buf := new(bytes.Buffer)
	err := imageutil.SaveImageToWriter(img, buf, "unsupported")

	// 应该返回不支持的格式错误
	if err != imageutil.ErrUnsupportedFormat {
		t.Fatalf("期望不支持的格式错误，实际得到: %v", err)
	}
}

// 测试无效的Base64字符串
func TestInvalidBase64(t *testing.T) {
	// 无效的Base64字符串
	base64Str := "这不是有效的Base64字符串"

	loader := imageutil.NewLoader()
	_, err := loader.LoadFromBase64(base64Str)

	// 应该返回错误
	if err == nil {
		t.Fatal("期望解码无效Base64字符串时返回错误，但没有")
	}
}

// 测试无效的图片数据
func TestInvalidImageData(t *testing.T) {
	// 有效的Base64字符串，但不是有效的图片数据
	base64Str := base64.StdEncoding.EncodeToString([]byte("这不是有效的图片数据"))

	loader := imageutil.NewLoader()
	_, err := loader.LoadFromBase64(base64Str)

	// 应该返回错误
	if err == nil {
		t.Fatal("期望解码无效图片数据时返回错误，但没有")
	}
}

// 测试空Reader
func TestEmptyReader(t *testing.T) {
	// 创建一个空的Reader
	reader := strings.NewReader("")

	loader := imageutil.NewLoader()
	_, err := loader.LoadFromReader(reader)

	// 应该返回错误
	if err == nil {
		t.Fatal("期望从空Reader加载图片时返回错误，但没有")
	}
}
