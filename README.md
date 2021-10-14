# gcv

[![Build Status](https://github.com/vcaesar/gcv/workflows/Go/badge.svg)](https://github.com/vcaesar/gcv/commits/master)
[![Build Status](https://travis-ci.org/vcaesar/gcv.svg)](https://travis-ci.org/vcaesar/gcv)
[![CircleCI Status](https://circleci.com/gh/vcaesar/gcv.svg?style=shield)](https://circleci.com/gh/vcaesar/gcv)
[![codecov](https://codecov.io/gh/vcaesar/gcv/branch/master/graph/badge.svg)](https://codecov.io/gh/vcaesar/gcv)
[![Go Report Card](https://goreportcard.com/badge/github.com/vcaesar/gcv)](https://goreportcard.com/report/github.com/vcaesar/gcv)
[![GoDoc](https://godoc.org/github.com/vcaesar/gcv?status.svg)](https://godoc.org/github.com/vcaesar/gcv)
[![Release](https://github-release-version.herokuapp.com/github/vcaesar/gcv/release.svg?style=flat)](https://github.com/vcaesar/gcv/releases/latest)


## Use

```go
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
```