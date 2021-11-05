// Copyright 2016 Evans. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gcv

import (
	"image"
	"image/color"

	"gocv.io/x/gocv"
)

// IMRead read the image file to gocv.Mat
func IMRead(file string, flag ...int) gocv.Mat {
	f1 := 4
	if len(flag) > 0 {
		f1 = flag[0]
	}
	return gocv.IMRead(file, gocv.IMReadFlag(f1))
}

// IMWrite write the gocv.Mat to file
func IMWrite(name string, img gocv.Mat) bool {
	return gocv.IMWrite(name, img)
}

// IMWrite write the image.Image to file
func ImgWrite(name string, img image.Image) bool {
	im, err := ImgToMat(img)
	if err != nil {
		return false
	}

	return IMWrite(name, im)
}

// ImgToMat trans image.Image to gocv.Mat
func ImgToMat(img image.Image) (gocv.Mat, error) {
	return gocv.ImageToMatRGB(img)
}

// ImgToMatA trans image.Image to gocv.Mat
func ImgToMatA(img image.Image) (gocv.Mat, error) {
	return gocv.ImageToMatRGBA(img)
}

// MatToImg trans gocv.Mat to image.Image
func MatToImg(m1 gocv.Mat) (image.Image, error) {
	return m1.ToImage()
}

// Show show the gocv.Mat image
func Show(img gocv.Mat, args ...interface{}) {
	wName := "show"
	if len(args) > 0 {
		wName = args[0].(string)
	}

	window := gocv.NewWindow(wName)
	defer window.Close()

	h, w := GetSize(img)
	if len(args) > 2 {
		w = args[1].(int)
		h = args[2].(int)
	}
	window.ResizeWindow(w, h)

	window.IMShow(img)
	window.WaitKey(0)
}

// GetSize get the cv.Mat width, hight
func GetSize(img gocv.Mat) (int, int) {
	return img.Rows(), img.Cols()
}

// Resize resize the image
func Resize(img gocv.Mat, sz image.Point, w, h float64) gocv.Mat {
	dst := gocv.NewMat()
	gocv.Resize(img, &dst, sz, w, h, gocv.InterpolationFlags(1))
	return dst
}

// Rotate rotate the image to 90, 180, 270 degrees
func Rotate(img gocv.Mat, args ...interface{}) gocv.Mat {
	angle := 0
	if len(args) > 0 {
		angle = args[0].(int)
	}
	code := gocv.RotateFlag(angle)
	mask := gocv.NewMat()

	gocv.Rotate(img, &mask, code)
	return mask
}

// Crop crop the cv.Mat by rect with crabcut
func Crop(img gocv.Mat, rect image.Rectangle) gocv.Mat {
	mask := gocv.NewMat()

	bgdModel := gocv.NewMat()
	defer bgdModel.Close()
	fgdModel := gocv.NewMat()
	defer fgdModel.Close()

	gocv.GrabCut(img, &mask, rect, &bgdModel, &fgdModel, 1, gocv.GCEval)
	return mask
}

// MarkPoint mark the point draw line to img
func MarkPoint(img gocv.Mat, point image.Point, args ...interface{}) gocv.Mat {
	circle := false
	if len(args) > 0 {
		circle = args[0].(bool)
	}
	colors := uint8(100)
	if len(args) > 1 {
		colors = args[1].(uint8)
	}
	radius := 20
	if len(args) > 2 {
		radius = args[2].(int)
	}

	if circle {
		gocv.Circle(&img, point, radius, color.RGBA{255, 0, 0, 0}, 2)
	}

	rg := color.RGBA{colors, 0, 0, 0}
	x, y := point.X, point.Y
	gocv.Line(&img, image.Point{x - radius, y}, image.Point{x + radius, y},
		rg, 2)
	gocv.Line(&img, image.Point{x, y - radius}, image.Point{x, y + radius},
		rg, 2)
	return img
}

// MaskImg rectangle the mark to img
func MaskImg(img *gocv.Mat, mark image.Rectangle, args ...interface{}) {
	lineW := -1
	if len(args) > 1 {
		lineW = args[1].(int)
	}
	offset := int(lineW / 2)

	colors := color.RGBA{255, 255, 255, 0}
	if len(args) > 0 {
		colors = args[0].(color.RGBA)
	}

	min := mark.Min
	max := mark.Max
	mark = image.Rectangle{
		Min: image.Point{min.X - offset, min.Y - offset},
		Max: image.Point{max.X + lineW, max.Y + lineW}}
	gocv.Rectangle(img, mark, colors, lineW)
}
