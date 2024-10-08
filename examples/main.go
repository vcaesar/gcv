package main

import (
	"fmt"

	"github.com/go-vgo/robotgo"
	"github.com/vcaesar/gcv"
)

func main() {
	// save images
	img, _ := robotgo.CaptureImg()
	img1, _ := robotgo.CaptureImg(18, 4, 20, 20)

	rs := gcv.FindAllImg(img1, img)
	if len(rs) > 0 {
		fmt.Println("find: ", rs[0].TopLeft.Y, rs[0].Rects.TopLeft.X, rs[0].ImgSize.H)
	}
	fmt.Println("find: ", rs)

	m1, _ := gcv.ImgToMat(img)
	m2, _ := gcv.ImgToMat(img1)

	rs = gcv.FindAllTemplate(m1, m2, 0.8)
	fmt.Println("find: ", rs)

	rs = gcv.FindAllSift(m1, m2)
	fmt.Println("find sift: ", rs)
}
