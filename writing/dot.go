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
	Force  uint16
	Scale  uint8
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
	tmp := fmt.Sprintf("act = %d;force = %d; tx = %d, ty = %d, fx = %d, fy = %d",
		mine.Action, mine.Force, mine.TX, mine.TY, mine.FX, mine.FY)
	return tmp
}

func (mine *DotInfo) Hex() string {
	return mine.hex
}


func (mine *DotInfo)Coordinate() (x float32, y float32) {
	x = float32(mine.TX) + float32(mine.FX)/100
	y = float32(mine.TY) + float32(mine.FY)/100
	return x,y
}

func (mine *DotInfo)DistX() float32 {
	return float32(mine.TX) + float32(mine.FX)/100
}

func (mine *DotInfo)DistY() float32 {
	return float32(mine.TY) + float32(mine.FY)/100
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

	x2 := (mine.DistX() * float32(dist.X) * canvasW) / w1
	y2 := (mine.DistY() * float32(dist.Y) * canvasH) / h1
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
	mine.TX = uint32(x)
	y, _ := strconv.ParseUint(hex[10:14], 16, 32)
	mine.TY = uint32(y)
	fx, _ := strconv.ParseUint(hex[14:16], 16, 32)
	mine.FX = uint32(fx)
	fy, _ := strconv.ParseUint(hex[16:18], 16, 32)
	mine.FY = uint32(fy)
	col, _ := strconv.ParseUint(hex[18:20], 16, 32)
	mine.Color = uint32(col)
	mine.Stamp, _ = strconv.ParseUint(hex[20:], 16, 64)
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
	mine.TX = uint32(x)
	y, _ := strconv.ParseUint(hex[10:14], 16, 32)
	mine.TY = uint32(y)
	fx, _ := strconv.ParseUint(hex[14:16], 16, 32)
	mine.FX = uint32(fx)
	fy, _ := strconv.ParseUint(hex[16:18], 16, 32)
	mine.FY = uint32(fy)
	col, _ := strconv.ParseUint(hex[18:20], 16, 32)
	mine.Color = uint32(col)
	mine.Stamp, _ = strconv.ParseUint(hex[20:], 16, 64)
	return nil
}

func (mine *DotInfo) AdjustTimestamp(_offset uint64) {
	mine.Stamp = mine.Stamp + _offset

}
