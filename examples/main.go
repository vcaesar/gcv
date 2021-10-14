package main

import (
	"fmt"

	"github.com/go-vgo/robotgo"
	"github.com/vcaesar/gcv"
)

func main() {
	img := robotgo.CaptureImg()
	img1 := robotgo.CaptureImg(18, 4, 20, 20)
	rs := gcv.FindAllImg(img1, img)
	fmt.Println("find: ", rs[0].TopLeft.Y, rs[0].Rects.TopLeft.X, rs[0].ImgSize.H)
	fmt.Println("find: ", rs)

	m1, _ := gcv.ImgToMat(img)
	m2, _ := gcv.ImgToMat(img1)
	rs = gcv.FindAllTemplate(m1, m2, 0.8)
	fmt.Println("find: ", rs)
}
