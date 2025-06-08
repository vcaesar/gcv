// Copyright 2016 Evans. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 <LICENSE-APACHE or
// http://www.apache.org/licenses/LICENSE-2.0>
// This file may not be copied, modified, or distributed
// except according to those terms.

package gcv

import (
	"image"
	"image/color"

	"github.com/go-vgo/gt/hset"
	"github.com/vcaesar/imgo"
	"gocv.io/x/gocv"
)

var (
	// Sift use the Sift mode default when find template result <=0
	Sift = true
)

// FindImgFile find image file in subfile
func FindImgFile(tempFile, file string, flag ...int) (float32, float32, image.Point, image.Point) {
	return FindImgMatC(IMRead(file, flag...), IMRead(tempFile, flag...))
}

// FindImg find image in the subImg
func FindImg(subImg, imgSource image.Image) (float32, float32, image.Point, image.Point) {
	m1, _ := ImgToMatA(imgSource)
	m2, _ := ImgToMatA(subImg)
	return FindImgMatC(m1, m2)
}

// FindImgByte find image in the subImg by []byte
func FindImgByte(subImg, imgSource []byte) (float32, float32, image.Point, image.Point) {
	m1, _ := imgo.ByteToImg(imgSource)
	m2, _ := imgo.ByteToImg(subImg)
	return FindImg(m2, m1)
}

// FindImgX find image in the subImg return x, y
func FindImgX(subImg, imgSource image.Image) (int, int) {
	_, _, _, maxLoc := FindImg(subImg, imgSource)
	return maxLoc.X, maxLoc.Y
}

// FindImgMatC find the image Mat in the temp Mat and close gocv.Mat
func FindImgMatC(imgSource, temp gocv.Mat) (float32, float32, image.Point, image.Point) {
	defer imgSource.Close()
	defer temp.Close()
	return FindImgMat(imgSource, temp)
}

// FindImgMat find the image Mat in the temp Mat
func FindImgMat(imgSource, temp gocv.Mat) (float32, float32, image.Point, image.Point) {
	res := gocv.NewMat()
	defer res.Close()
	msk := gocv.NewMat()
	defer msk.Close()

	gocv.MatchTemplate(imgSource, temp, &res, gocv.TmCcoeffNormed, msk)
	minVal, maxVal, minLoc, maxLoc := gocv.MinMaxLoc(res)

	return minVal, maxVal, minLoc, maxLoc
}

// FlannbasedMatch new flann based match
func FlannbasedMatch(query, train gocv.Mat, k int) [][]gocv.DMatch {
	fb := gocv.NewFlannBasedMatcher()
	return fb.KnnMatch(query, train, k)
}

// Point the image X, Y point structure
type Point struct {
	X, Y int
}

// Size the image size structure
type Size struct {
	W, H int
}

// Rect the image rectangles structure
type Rect struct {
	TopLeft, TopRight       Point
	BottomLeft, BottomRight Point
}

// Result find template result structure
type Result struct {
	Middle, TopLeft Point // maxLoc
	Rects           Rect

	MaxVal  []float32
	ImgSize Size
}

// Find find all the img search in the img source by
// find all template and sift and return Result
func Find(imgSearch, imgSource image.Image, args ...interface{}) (r Result) {
	res := FindAll(imgSearch, imgSource, args...)
	if len(res) > 0 {
		r = res[0]
	}
	return
}

// FindX find all the img search in the img source by
// find all template and sift and return x, y
func FindX(imgSearch, imgSource image.Image, args ...interface{}) (x, y int) {
	res := Find(imgSearch, imgSource, args...)
	x, y = res.Middle.X, res.Middle.Y
	return
}

// FindAllX find all the img search in the img source by
// find all template and sift and return []x, []y
func FindAllX(imgSearch, imgSource image.Image, args ...interface{}) (x, y []int) {
	res := FindAll(imgSearch, imgSource, args...)
	for i := 0; i < len(res); i++ {
		x = append(x, res[i].Middle.X)
		y = append(y, res[i].Middle.Y)
	}
	return
}

// FindAllImg find the search image all template in the source image return []Result
func FindAllImg(imgSearch, imgSource image.Image, args ...interface{}) []Result {
	matSource, _ := ImgToMatA(imgSource)
	matSearch, _ := ImgToMatA(imgSearch)

	return FindAllTemplateC(matSource, matSearch, args...)
}

// FindAll find all the img search in the img source by
// find all template and sift and return []Result
func FindAll(imgSearch, imgSource image.Image, args ...interface{}) []Result {
	matSource, _ := ImgToMatA(imgSource)
	matSearch, _ := ImgToMatA(imgSearch)

	res := FindAllTemplateC(matSource, matSearch, args...)
	if len(res) <= 0 && Sift {
		res = FindAllSiftC(matSource, matSearch, args...)
	}
	return res
}

// FindAllImgFlie find the search image all template in the source image file
// return []Result
func FindAllImgFile(fileSearh, file string, args ...interface{}) []Result {
	return FindAllTemplateC(IMRead(file), IMRead(fileSearh), args...)
}

// FindMultiAllImgFile find the multi file search image all template
// in the file source image return [][]Result
func FindMultiAllImgFile(fileSearh []string, file string, args ...interface{}) [][]Result {
	matSource := IMRead(file)
	var matSearch []gocv.Mat
	for i := 0; i < len(fileSearh); i++ {
		search := IMRead(fileSearh[i])
		matSearch = append(matSearch, search)
	}
	return FindMultiAllTemplateC(matSource, matSearch, args...)
}

// FindMultiAllImg find the multi search image all template in the source image
// return [][]Result
func FindMultiAllImg(imgSearch []image.Image, imgSource image.Image, args ...interface{}) [][]Result {
	matSource, _ := ImgToMatA(imgSource)
	var matSearch []gocv.Mat
	for i := 0; i < len(imgSearch); i++ {
		search, _ := ImgToMatA(imgSearch[i])
		matSearch = append(matSearch, search)
	}

	return FindMultiAllTemplateC(matSource, matSearch, args...)
}

// FindMultiAllTemplate find the multi imgSearch all template in the imgSource return [][]Result
// and close gocv.Mat
func FindMultiAllTemplateC(imgSource gocv.Mat, imgSearch []gocv.Mat, args ...interface{}) (r [][]Result) {
	defer imgSource.Close()
	for i := 0; i < len(imgSearch); i++ {
		r = append(r, FindAllTemplateCS(imgSource, imgSearch[i], args...))
	}

	return
}

// FindMultiAllTemplate find the multi imgSearch all template in the imgSource return [][]Result
func FindMultiAllTemplate(imgSource gocv.Mat, imgSearch []gocv.Mat, args ...interface{}) (r [][]Result) {
	for i := 0; i < len(imgSearch); i++ {
		r = append(r, FindAllTemplate(imgSource, imgSearch[i], args...))
	}

	return
}

// FindAllTemplateCS find the imgSearch all template in the imgSource return []Result
// and close gocv.Mat
func FindAllTemplateCS(imgSource, imgSearch gocv.Mat, args ...interface{}) []Result {
	// defer imgSource.Close()
	defer imgSearch.Close()
	return FindAllTemplate(imgSource, imgSearch, args...)
}

// FindAllTemplateC find the imgSearch all template in the imgSource return []Result
// and close gocv.Mat
func FindAllTemplateC(imgSource, imgSearch gocv.Mat, args ...interface{}) []Result {
	defer imgSource.Close()
	defer imgSearch.Close()
	return FindAllTemplate(imgSource, imgSearch, args...)
}

// FindAllSiftC find the imgSearch all sift in the imgSource return []Result
// and close gocv.Mat
func FindAllSiftC(matSource, matSearch gocv.Mat, args ...interface{}) []Result {
	// defer matSource.Close()
	// defer matSearch.Close()
	return FindAllSift(matSource, matSearch, args...)
}

// FindAllTemplate find the imgSearch all template in the imgSource return []Result
func FindAllTemplate(matSource, matSearch gocv.Mat, args ...interface{}) []Result {
	threshold := float32(0.8)
	if len(args) > 0 {
		threshold = float32(args[0].(float64))
	}

	maxCount := 10
	if len(args) > 1 {
		maxCount = args[1].(int)
	}
	// rgb := false
	// if len(args) > 2 {
	// 	rgb = args[2].(bool)
	// }

	method := false
	if len(args) > 3 {
		method = args[3].(bool)
	}

	iGray := gocv.NewMat()
	defer iGray.Close()
	sGray := gocv.NewMat()
	defer sGray.Close()

	// if !rgb {
	gocv.CvtColor(matSource, &iGray, gocv.ColorRGBToGray)
	gocv.CvtColor(matSearch, &sGray, gocv.ColorRGBToGray)
	// }

	results := make([]Result, 0)
	for {
		_, maxVal, minLoc, maxLoc := FindImgMat(iGray, sGray)
		h, w := matSearch.Rows(), matSearch.Cols()
		if maxVal < threshold || len(results) > maxCount {
			break
		}

		if method {
			maxLoc = minLoc
		}

		rs := getVal(maxLoc, maxVal, w, h)
		results = append(results, rs)
		if len(results) <= 0 {
			return nil
		}

		Fill(iGray, rs.Rects)
		// Rectangle(iGray, maxLoc, w, h)
		// Show(iGray)
	}

	return results
}

func getVal(maxLoc image.Point, maxVal float32, w, h int) Result {
	rect := Rect{
		TopLeft:     Point{maxLoc.X, maxLoc.Y},
		BottomLeft:  Point{maxLoc.X, maxLoc.Y + h},
		BottomRight: Point{maxLoc.X + w, maxLoc.Y + h},
		TopRight:    Point{maxLoc.X + w, maxLoc.Y},
	}

	middle := image.Pt(maxLoc.X+w/2, maxLoc.Y+h/2)
	middlePoint := Point{middle.X, middle.Y}

	topLeft := Point{maxLoc.X, maxLoc.Y}
	maxVals := []float32{maxVal}
	size := Size{w, h}

	return Result{
		Middle:  middlePoint,
		TopLeft: topLeft, // MaxLoc
		Rects:   rect,
		MaxVal:  maxVals,
		ImgSize: size,
	}
}

// Rectangle rectangle the iGray image
func Rectangle(iGray gocv.Mat, maxLoc image.Point, w, h int) {
	r := image.Rect(
		int(maxLoc.X-w/2),
		int(maxLoc.Y-h/2),
		int(maxLoc.X+w/2),
		int(maxLoc.Y+h/2),
	)

	// white := color.RGBA{255, 255, 255, 255}
	black := color.RGBA{0, 0, 0, 0}
	gocv.Rectangle(&iGray, r, black, -1)
}

// Fill fillpoly the iGray image
func Fill(iGray gocv.Mat, rect Rect) {
	pts := [][]image.Point{{
		image.Pt(rect.TopLeft.X, rect.TopLeft.Y),
		image.Pt(rect.BottomLeft.X, rect.BottomLeft.Y),
		image.Pt(rect.BottomRight.X, rect.BottomLeft.Y),
		image.Pt(rect.TopRight.X, rect.TopRight.Y),
	}}

	pts1 := gocv.NewPointsVectorFromPoints(pts)
	defer pts1.Close()

	blue := color.RGBA{0, 0, 255, 0}
	gocv.FillPoly(&iGray, pts1, blue)
}

func findH(kpS, kpSrc []gocv.KeyPoint, goodDiff []gocv.DMatch) (gocv.Mat, gocv.Mat) {
	src := gocv.NewMatWithSize(len(goodDiff), 1, gocv.MatTypeCV64FC2)
	defer src.Close()
	dst := gocv.NewMatWithSize(len(goodDiff), 1, gocv.MatTypeCV64FC2)
	// defer dst.Close()
	mask := gocv.NewMat()
	// defer mask.Close()

	// Get the keypoints from the good matches
	for i := 0; i < len(goodDiff); i++ {
		src.SetDoubleAt(i, 0, kpS[goodDiff[i].QueryIdx].X)
		src.SetDoubleAt(i, 1, kpS[goodDiff[i].QueryIdx].Y)

		dst.SetDoubleAt(i, 0, kpSrc[goodDiff[i].TrainIdx].X)
		dst.SetDoubleAt(i, 1, kpSrc[goodDiff[i].TrainIdx].Y)
	}

	// find estimate H
	hm := gocv.FindHomography(src, dst, gocv.HomographyMethodRANSAC, 5.0,
		&mask, 2000, 0.95)

	return hm, mask
}

func transform(h, w int, hm gocv.Mat) gocv.Mat {
	dst := gocv.NewMatWithSize(4, 2, gocv.MatTypeCV32FC2)
	// defer dst.Close()

	// new mat from the img search
	src := gocv.NewMatWithSize(4, 2, gocv.MatTypeCV32FC2)
	defer src.Close()

	src.SetFloatAt(0, 0, 0)
	src.SetFloatAt(0, 1, 0)

	src.SetFloatAt(1, 0, 0)
	src.SetFloatAt(1, 1, float32(h-1))

	src.SetFloatAt(2, 0, float32(w-1))
	src.SetFloatAt(2, 1, float32(h-1))

	src.SetFloatAt(3, 0, float32(w-1))
	src.SetFloatAt(3, 1, 0)

	gocv.PerspectiveTransform(src, &dst, hm)
	return dst
}

func getPoint(ms [][]gocv.DMatch, kpS, kpSrc []gocv.KeyPoint, ratio float64) []gocv.DMatch {
	var good, goodDiff []gocv.DMatch

	// Filter matches low distance by ratio
	for i := 0; i < len(ms); i++ {
		if ms[i][0].Distance < ratio*ms[i][1].Distance {
			good = append(good, ms[i][0])
		}
	}

	// Remove the duplicates point
	h1 := hset.New()
	for i := 0; i < len(good); i++ {
		p1 := Point{
			X: int(kpSrc[good[i].TrainIdx].X),
			Y: int(kpSrc[good[i].TrainIdx].Y),
		}

		if !h1.Exists(p1) {
			h1.Add(p1)
			goodDiff = append(goodDiff, good[i])
		}
	}

	return goodDiff
}

// FindAllSift find the matSearch all sift in matSource return result
func FindAllSift(matSource, matSearch gocv.Mat, args ...interface{}) (res []Result) {
	sift := gocv.NewSIFT()
	defer sift.Close()
	//
	mask1 := gocv.NewMat()
	defer mask1.Close()
	mask2 := gocv.NewMat()
	defer mask2.Close()

	minMatch := 4
	ratio := 0.75 // 0.9

	// detect the feature and compute descriptor
	kpS, des := sift.DetectAndCompute(matSearch, mask1)
	kpSrc, deSrc := sift.DetectAndCompute(matSource, mask2)
	if len(kpS) < 2 || len(kpSrc) < 2 || len(kpS) < minMatch {
		return
	}

	ms := FlannbasedMatch(des, deSrc, 2)
	goodDiff := getPoint(ms, kpS, kpSrc, ratio)

	if len(goodDiff) == 0 {
		return
	}
	defer matSearch.Close()
	defer matSource.Close()

	h, w := GetSize(matSearch)
	// get the result value
	if len(goodDiff) == 1 {
		kpt := kpSrc[goodDiff[0].TrainIdx]
		middlePoint := Point{int(kpt.X), int(kpt.Y)}

		res = append(res, Result{
			Middle:  middlePoint,
			MaxVal:  []float32{0.5},
			Rects:   Rect{},
			ImgSize: Size{w, h},
		})

		return
	}

	// if len(goodDiff) > 3 {
	hm, mask := findH(kpS, kpSrc, goodDiff)
	dst := transform(h, w, hm)
	hm.Close()
	defer mask.Close()
	defer dst.Close()

	p1 := Point{int(dst.GetFloatAt(0, 0)), int(dst.GetFloatAt(0, 1))}
	p2 := Point{int(dst.GetFloatAt(1, 0)), int(dst.GetFloatAt(1, 1))}
	p3 := Point{int(dst.GetFloatAt(2, 0)), int(dst.GetFloatAt(2, 1))}
	p4 := Point{int(dst.GetFloatAt(3, 0)), int(dst.GetFloatAt(3, 1))}

	at := h / w
	res = append(res, Result{
		Middle:  Point{p1.X + (p3.X-p1.X)/2, p1.Y + (p3.Y-p1.Y)/2},
		TopLeft: p1,
		Rects: Rect{
			TopLeft:     p1,
			BottomLeft:  p2,
			BottomRight: p3,
			TopRight:    p4,
		},
		MaxVal:  []float32{float32(at), float32(len(goodDiff))},
		ImgSize: Size{w, h},
	})
	// }

	return
}
