// Copyright 2016 Evans. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package gcv

import (
	"image"

	"github.com/vcaesar/imgo"
	"gocv.io/x/gocv"
)

// FindImgFile find image file in subfile
func FindImgFile(file, subFile string, flag ...int) (float32, float32, image.Point, image.Point) {
	f1 := 4
	if len((flag)) > 0 {
		f1 = flag[0]
	}
	flags := gocv.IMReadFlag(f1)
	imrgb := gocv.IMRead(file, flags)
	temp := gocv.IMRead(subFile, flags)

	return FindImgMat(imrgb, temp)
}

// FindImg find image in the subImg
func FindImg(img, subImg image.Image) (float32, float32, image.Point, image.Point) {
	m1, _ := ImgToMat(img)
	m2, _ := ImgToMat(subImg)
	return FindImgMat(m1, m2)
}

// FindImgByte find image in the subImg by []byte
func FindImgByte(img, subImg []byte) (float32, float32, image.Point, image.Point) {
	m1, _ := imgo.ByteToImg(img)
	m2, _ := imgo.ByteToImg(subImg)
	return FindImg(m1, m2)
}

// FindImgXY find image in the subImg return x, y
func FindImgXY(img, subImg image.Image) (int, int) {
	_, _, _, maxLoc := FindImg(img, subImg)
	return maxLoc.X, maxLoc.Y
}

// FindImgMat find the image Mat in the temp Mat
func FindImgMat(imRgb, temp gocv.Mat) (float32, float32, image.Point, image.Point) {
	res := gocv.NewMat()
	defer res.Close()
	msk := gocv.NewMat()
	defer msk.Close()

	gocv.MatchTemplate(imRgb, temp, &res, gocv.TmCcoeffNormed, msk)
	minVal, maxVal, minLoc, maxLoc := gocv.MinMaxLoc(res)

	return minVal, maxVal, minLoc, maxLoc
}

// ImgToMat trans image.Image to gocv.Mat
func ImgToMat(img image.Image) (gocv.Mat, error) {
	return gocv.ImageToMatRGB(img)
}

// MatToImg trans gocv.Mat to image.Image
func MatToImg(m1 gocv.Mat) (image.Image, error) {
	return m1.ToImage()
}

// Show show the gocv.Mat image
func Show(img gocv.Mat) {
	window := gocv.NewWindow("show")
	defer window.Close()
	window.ResizeWindow(1200, 800)
	window.IMShow(img)
	window.WaitKey(0)
}
