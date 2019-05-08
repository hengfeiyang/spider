package util

import (
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"

	"golang.org/x/image/draw"
)

// Scale 缩放图像
// 参数：src 原始图像路径，dst 目标图像路径，如果为空或相同则覆盖原始文件，w,h 为缩放后的宽高
// equalRate 参数决定是否等比缩放，开启等比缩放时，将大小控制在指定的区域内，图像尺寸小于等于给定尺寸
// cut 参数觉决定是否使用裁剪方案，当使用裁剪时，使用较小的缩放比，满足最小边在指定区域内，同时将多余（在区域外）的部分减掉
// 当使用 cut 启用裁剪时，忽略 equalRate的设置，使用等比缩放并裁剪
// 缩放函数有四个 draw.NearestNeighbor, draw.ApproxBiLinear, draw.BiLinear, draw.CatmullRom,
// 具体区别参见：https://godoc.org/github.com/golang/image/draw#pkg-variables//
// 使用方法，非常简单：
/*  err = Scale("2.jpg", "a1.png", 100, 100, false, false)
err = Scale("2.jpg", "a2.png", 100, 100, true, false)
err = Scale("2.jpg", "a3.png", 100, 100, false, true)
err = Scale("2.jpg", "a4.png", 100, 100, true, true)
*/
func Scale(src, dst string, w, h int, equalRate bool, cut bool) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}

	ext := filepath.Ext(src)
	var imgSrc image.Image
	// 判断源图片格式是否支持
	switch ext {
	case ".jpg", ".jpeg":
		imgSrc, err = jpeg.Decode(f)
	case ".png":
		imgSrc, err = png.Decode(f)
	case ".gif":
		imgSrc, err = gif.Decode(f)
	default:
		err = fmt.Errorf("unknown input image type")
	}
	// 读取完数据源图片就可以关闭了
	f.Close()

	if err != nil {
		return err
	}
	if dst == "" || dst == src {
		dst = src
	}
	f, err = os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	ext = filepath.Ext(dst)
	// 判断导出图片格式是否支持
	switch ext {
	case ".jpg", ".jpeg":
		err = nil
	case ".png":
		err = nil
	case ".gif":
		err = nil
	default:
		err = fmt.Errorf("unknown output image type")
	}
	if err != nil {
		return err
	}

	// 计算缩放比例，当启用裁剪时按小比率，否则为大比率
	if cut {
		// 当启用裁剪时，一定使用等比
		equalRate = true
	}
	var rate float64
	srcBound := imgSrc.Bounds()
	rateX := float64(srcBound.Dx()) / float64(w)
	rateY := float64(srcBound.Dy()) / float64(h)
	if rateX == rateY {
		rate = rateX // 正好等比
	} else {
		if rateX > rateY {
			if cut {
				rate = rateY
			} else {
				rate = rateX
			}
		} else {
			if cut {
				rate = rateX
			} else {
				rate = rateY
			}
		}
	}

	// 计算目标尺寸，启用等比时，目标尺寸不等于设定尺寸，所以要先计算比率
	dstW := int(float64(srcBound.Dx()) / rate)
	dstH := int(float64(srcBound.Dy()) / rate)
	// 如果没启用等比，一定没启用裁剪，直接赋值，不考虑比率问题
	if equalRate == false {
		dstW = w
		dstH = h
	}

	// 首次处理，先缩放
	imgDst := image.NewRGBA(image.Rect(0, 0, dstW, dstH))
	draw.CatmullRom.Scale(imgDst, imgDst.Bounds(), imgSrc, imgSrc.Bounds(), draw.Src, nil)

	// 不等比，或者正好比率相同也无需裁剪，直接保存
	if equalRate && rateX != rateY {
		// 裁剪图片，输出目标尺寸，截取中间部分
		imgDst = imgDst.SubImage(image.Rectangle{
			Min: image.Point{(dstW - w) / 2, (dstH - h) / 2},
			Max: image.Point{w + (dstW-w)/2, h + (dstH-h)/2},
		}).(*image.RGBA)
	}

	// 导出
	switch ext {
	case ".jpg", ".jpeg":
		err = jpeg.Encode(f, imgDst, &jpeg.Options{Quality: 75})
	case ".png":
		encoder := png.Encoder{
			CompressionLevel: png.BestCompression,
		}
		err = encoder.Encode(f, imgDst)
	case ".gif":
		err = gif.Encode(f, imgDst, nil)
	default:
		err = fmt.Errorf("unknown output image type")
	}
	return err
}
