package main

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/dfengwji/pss.common/writing"
)

func TestDot(t *testing.T) {
	// <setup code>
	t.Run("A=1", func(t *testing.T) {
		dot := writing.DotInfo{}
		fmt.Println(dot.String())
	})
	// <tear-down code>
}

func Test_SavePNG(t *testing.T) {
	bytes, err := ioutil.ReadFile("/tmp/dots.txt")
	if nil != err {
		t.Error(err)
		t.Fail()
		return
	}
	hex := string(bytes)
	num := len(hex) / writing.DotHexLength
	c := num * writing.DotHexLength
	t.Logf("origin -> length of hex is %v, number is %v", len(hex), float32(len(hex))/float32(writing.DotHexLength))
	t.Logf("wanted -> length of hex is %v, number is %v", c, num)
	if len(hex) != c {
		t.Error("the length not equal")
		t.Fail()
		return
	}

	dots := make([]*writing.DotInfo, num)
	for i := 0; i < len(dots); i++ {
		hexVal := hex[i*writing.DotHexLength : (i+1)*writing.DotHexLength]
		//t.Log(hexVal)
		dots[i] = new(writing.DotInfo)
		dots[i].ParseHexV1(hexVal, writing.DotHexLength)
		//t.Log(dots[i])
	}

	options := writing.RenderOptions{
		PaperWidth:  210,
		PaperHeight: 290,
		ImageWidth:  1080,
		ImageHeight: 1920,
		PenSize:     1.0,
		ClipSize:    320,
		PaddingSize: 120,
	}
	if err := writing.SavePNG(dots, options, "/tmp/dots.png"); nil != err {
		t.Error(err)
		t.Fail()
	}
}
