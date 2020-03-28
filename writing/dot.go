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

const DotHexLength = 30

type DotInfo struct {
	Action uint8
	Force  uint16
	Scale  uint8
	X      uint32
	Y      uint32
	FX     uint32
	FY     uint32
	Color  uint32
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
	tmp := fmt.Sprintf("act = %d;force = %d; x = %d, y = %d, fx = %d, fy = %d",
		mine.Action, mine.Force, mine.X, mine.Y, mine.FX, mine.FY)
	return tmp
}

func (mine *DotInfo) Hex() string {
	return mine.hex
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
	mine.X = uint32(x)
	y, _ := strconv.ParseUint(hex[10:14], 16, 32)
	mine.Y = uint32(y)
	fx, _ := strconv.ParseUint(hex[14:16], 16, 32)
	mine.FX = uint32(fx)
	fy, _ := strconv.ParseUint(hex[16:18], 16, 32)
	mine.FY = uint32(fy)
	col, _ := strconv.ParseUint(hex[18:20], 16, 32)
	mine.Color = uint32(col)
	mine.Stamp, _ = strconv.ParseUint(hex[20:30], 16, 64)
	return nil
}

func (this *DotInfo) AdjustTimestamp(_offset uint64) {
	this.Stamp = this.Stamp + _offset

}
