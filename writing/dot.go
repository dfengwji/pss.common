package writing

import (
	"fmt"
	"github.com/pkg/errors"
	"image/color"
	"strconv"
)

const (
	DotActionDown  = 0
	DotActionMove  = 1
	DotActionUp    = 2
	DotActionHover = 3
)

const DotHexLength = 31

type DotInfo struct {
	Action uint8
	Scale  uint8
	Force  uint16
	TX      uint32
	TY      uint32
	FX     uint32
	FY     uint32
	Color  uint32
	X float32
	Y float32
	Stamp  uint64
	hex    string
}

type PointSample struct {
	Color  color.Color
	Action uint8
	Force  uint16
	Scale  uint8
	X      float64
	Y      float64
}

func (mine *DotInfo) String() string {
	tmp := fmt.Sprintf("act = %d;force = %d; x = %f, y = %f",
		mine.Action, mine.Force, mine.X, mine.Y)
	return tmp
}

func (mine *DotInfo) Hex() string {
	return mine.hex
}


func (mine *DotInfo)coordinate(tx, ty, fx, fy uint32) {
	mine.TX = tx
	mine.TY = ty
	mine.FX = fx
	mine.FY = fy
	mine.X = float32(tx) + float32(fx)/100
	mine.Y = float32(ty) + float32(fy)/100
}

// Transform 根据纸张大小和画布大小转换XY坐标
// TODO: 动态适配纸张大小
func (mine *DotInfo) Transform(w float32, h float32, code Vector2, dist Vector2) (float32, float32) {
	w = w - 10 // TODO: 移植到video中
	h = h - 10

	w1 := float32(code.X)
	h1 := float32(code.Y)

	// 以画布宽度为准
	canvasW := w
	canvasH := canvasW * h1 / w1

	// 如果高度超出，则以画布高度为准
	if canvasH > h {
		canvasH = h
		canvasW = canvasH * w1 / h1
	}

	x2 := (mine.X * float32(dist.X) * canvasW) / w1
	y2 := (mine.Y * float32(dist.Y) * canvasH) / h1
	return x2, y2
}

func (mine *DotInfo) ParseHex(hex string) error {
	if len(hex) < DotHexLength {
		return errors.New("the dot hex length is less than 30")
	}
	mine.hex = hex
	act, _ := strconv.ParseUint(hex[0:2], 10, 8)
	mine.Action = uint8(act)
	f, _ := strconv.ParseUint(hex[2:5], 16, 32)
	mine.Force = uint16(f)
	s, _ := strconv.ParseUint(hex[5:6], 16, 32)
	mine.Scale = uint8(s)
	x, _ := strconv.ParseUint(hex[6:10], 16, 32)
	y, _ := strconv.ParseUint(hex[10:14], 16, 32)
	fx, _ := strconv.ParseUint(hex[14:16], 16, 32)
	fy, _ := strconv.ParseUint(hex[16:18], 16, 32)
	col, _ := strconv.ParseUint(hex[18:20], 16, 32)
	mine.Color = uint32(col)
	mine.Stamp, _ = strconv.ParseUint(hex[20:], 16, 64)
	mine.coordinate(uint32(x),  uint32(y), uint32(fx), uint32(fy))
	return nil
}

func (mine *DotInfo) ParseHexV1(hex string, hexLength int) error {
	if len(hex) < hexLength {
		return errors.New(fmt.Sprintf("the length of hex is less than %v", hexLength))
	}
	mine.hex = hex
	act, _ := strconv.ParseUint(hex[0:2], 10, 8)
	mine.Action = uint8(act)
	f, _ := strconv.ParseUint(hex[2:5], 16, 32)
	mine.Force = uint16(f)
	s, _ := strconv.ParseUint(hex[5:6], 16, 32)
	mine.Scale = uint8(s)
	x, _ := strconv.ParseUint(hex[6:10], 16, 32)
	y, _ := strconv.ParseUint(hex[10:14], 16, 32)
	fx, _ := strconv.ParseUint(hex[14:16], 16, 32)
	fy, _ := strconv.ParseUint(hex[16:18], 16, 32)
	col, _ := strconv.ParseUint(hex[18:20], 16, 32)
	mine.Color = uint32(col)
	mine.Stamp, _ = strconv.ParseUint(hex[20:], 16, 64)
	mine.coordinate(uint32(x),  uint32(y), uint32(fx), uint32(fy))
	return nil
}

func (mine *DotInfo) AdjustTimestamp(_offset uint64) {
	mine.Stamp = mine.Stamp + _offset

}
