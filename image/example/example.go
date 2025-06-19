// Package main 提供图片处理模块的使用示例
package main

import (
	"fmt"
	"image/color"
	"log"
	"os"
	"path/filepath"

	imageutil "github.com/gophertool/tool/image"
)

func main() {
	// 创建图片加载器
	loader := imageutil.NewLoader()

	// 示例1：从文件加载图片
	exampleLoadFromFile(loader)

	// 示例2：从URL加载图片
	exampleLoadFromURL(loader)

	// 示例3：从Base64加载图片
	exampleLoadFromBase64(loader)

	// 示例4：从字节数组加载图片
	exampleLoadFromBytes(loader)

	// 示例5：保存图片
	exampleSaveImage(loader)
}

// 示例1：从文件加载图片
func exampleLoadFromFile(loader imageutil.Loader) {
	fmt.Println("示例1：从文件加载图片")

	// 注意：需要提供一个有效的图片文件路径
	filePath := "path/to/image.jpg"
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		fmt.Printf("文件 %s 不存在，跳过此示例\n\n", filePath)
		return
	}

	img, err := loader.LoadFromFile(filePath)
	if err != nil {
		log.Printf("从文件加载图片失败: %v\n\n", err)
		return
	}

	// 打印图片信息
	bounds := img.Bounds()
	fmt.Printf("成功加载图片，尺寸: %dx%d\n\n", bounds.Dx(), bounds.Dy())
}

// 示例2：从URL加载图片
func exampleLoadFromURL(loader imageutil.Loader) {
	fmt.Println("示例2：从URL加载图片")

	// 使用一个公共的测试图片URL
	url := "https://via.placeholder.com/150"

	img, err := loader.LoadFromURL(url)
	if err != nil {
		log.Printf("从URL加载图片失败: %v\n\n", err)
		return
	}

	// 打印图片信息
	bounds := img.Bounds()
	fmt.Printf("成功从URL加载图片，尺寸: %dx%d\n\n", bounds.Dx(), bounds.Dy())
}

// 示例3：从Base64加载图片
func exampleLoadFromBase64(loader imageutil.Loader) {
	fmt.Println("示例3：从Base64加载图片")

	// 这是一个1x1像素的PNG图片的base64编码
	base64Str := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+P+/HgAFeAJ5jfZixgAAAABJRU5ErkJggg=="

	img, err := loader.LoadFromBase64(base64Str)
	if err != nil {
		log.Printf("从Base64加载图片失败: %v\n\n", err)
		return
	}

	// 打印图片信息
	bounds := img.Bounds()
	fmt.Printf("成功从Base64加载图片，尺寸: %dx%d\n\n", bounds.Dx(), bounds.Dy())

	// 获取像素颜色
	pixel := img.At(0, 0)
	r, g, b, a := color.RGBAModel.Convert(pixel).(color.RGBA).RGBA()
	fmt.Printf("像素颜色 (0,0): R=%d, G=%d, B=%d, A=%d\n\n", r>>8, g>>8, b>>8, a>>8)
}

// 示例4：从字节数组加载图片
func exampleLoadFromBytes(loader imageutil.Loader) {
	fmt.Println("示例4：从字节数组加载图片")

	// 这里我们从文件读取字节，作为示例
	// 注意：需要提供一个有效的图片文件路径
	filePath := "path/to/image.jpg"
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		fmt.Printf("文件 %s 不存在，使用测试数据代替\n", filePath)

		// 使用一个简单的1x1像素图片数据作为替代
		// 这是PNG文件头的一部分，足够让程序识别为PNG格式
		data := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}

		// 获取图片格式
		format, err := imageutil.GetImageFormat(data)
		if err != nil {
			log.Printf("获取图片格式失败: %v\n\n", err)
			return
		}

		fmt.Printf("测试数据的图片格式: %s\n\n", format)
		return
	}

	// 读取文件内容
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("读取文件失败: %v\n\n", err)
		return
	}

	// 获取图片格式
	format, err := imageutil.GetImageFormat(data)
	if err != nil {
		log.Printf("获取图片格式失败: %v\n\n", err)
		return
	}

	fmt.Printf("图片格式: %s\n", format)

	// 从字节数组加载图片
	img, err := loader.LoadFromBytes(data)
	if err != nil {
		log.Printf("从字节数组加载图片失败: %v\n\n", err)
		return
	}

	// 打印图片信息
	bounds := img.Bounds()
	fmt.Printf("成功从字节数组加载图片，尺寸: %dx%d\n\n", bounds.Dx(), bounds.Dy())
}

// 示例5：保存图片
func exampleSaveImage(loader imageutil.Loader) {
	fmt.Println("示例5：保存图片")

	// 从Base64加载一个简单的图片
	base64Str := "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+P+/HgAFeAJ5jfZixgAAAABJRU5ErkJggg=="

	img, err := loader.LoadFromBase64(base64Str)
	if err != nil {
		log.Printf("加载图片失败: %v\n\n", err)
		return
	}

	// 创建临时目录用于保存图片
	tmpDir, err := os.MkdirTemp("", "imageutil-example")
	if err != nil {
		log.Printf("创建临时目录失败: %v\n\n", err)
		return
	}
	defer os.RemoveAll(tmpDir) // 清理

	// 保存为PNG
	pngPath := filepath.Join(tmpDir, "example.png")
	err = imageutil.SaveImage(img, pngPath, "png")
	if err != nil {
		log.Printf("保存PNG图片失败: %v\n\n", err)
		return
	}

	// 保存为JPEG
	jpegPath := filepath.Join(tmpDir, "example.jpg")
	err = imageutil.SaveImage(img, jpegPath, "jpeg")
	if err != nil {
		log.Printf("保存JPEG图片失败: %v\n\n", err)
		return
	}

	// 获取文件信息
	pngInfo, err := os.Stat(pngPath)
	if err != nil {
		log.Printf("获取PNG文件信息失败: %v\n\n", err)
		return
	}

	jpegInfo, err := os.Stat(jpegPath)
	if err != nil {
		log.Printf("获取JPEG文件信息失败: %v\n\n", err)
		return
	}

	fmt.Printf("成功保存图片:\n")
	fmt.Printf("PNG文件: %s (大小: %d字节)\n", pngPath, pngInfo.Size())
	fmt.Printf("JPEG文件: %s (大小: %d字节)\n\n", jpegPath, jpegInfo.Size())
}
