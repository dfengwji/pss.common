package writing

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"io"
)

const (
	OCRTypeText  = "Text"
	OCRTypeMath  = "Math"
	OCRTypeGraph = "Graph"
)

type RecognizerInfo struct {
	Type    string         `json:"contentType"`
	DPIx    uint32         `json:"xDPI"`
	DPIy    uint32         `json:"yDPI"`
	Config  *ConfigInfo     `json:"configuration"`
	Strokes []*StrokeGroup `json:"strokeGroups"`
	dots    []*DotInfo     `json:"-"`
}

type ConfigInfo struct {
	Math MathInfo `json:"math"`
}

type MathInfo struct {
	Solver SolverInfo `json:"solver"`
}

type SolverInfo struct {
	Enable bool `json:"enable"`
}

type StrokeGroup struct {
	PenStyle *PenStyleInfo `json:"penStyle"`
	Strokes  []*Stroke     `json:"strokes"`
}

type PenStyleInfo struct {
}

type Stroke struct {
	Xs []float64 `json:"x"`
	Ys []float64 `json:"y"`
	Ts []uint64  `json:"t"`
}

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

func (mine *RecognizerInfo) NewStroke() *Stroke {
	stroke := new(Stroke)
	stroke.Xs = make([]float64, 0, 100)
	stroke.Ys = make([]float64, 0, 100)
	stroke.Ts = make([]uint64, 0, 100)
	return stroke
}

func (mine *RecognizerInfo) SetDots(kind string, array []*DotInfo) {
	mine.dots = array
	mine.Type = kind
}

func (mine *RecognizerInfo) Encrypt(msg string, appKey string, hmacKey string) (string,error) {
	h := hmac.New(sha512.New, []byte(appKey+hmacKey))
	_, err := io.WriteString(h, msg)
	if err != nil {
		return "",err
	}
	return fmt.Sprintf("%x", h.Sum(nil)),nil
}

func (mine *RecognizerInfo) Bytes() []byte {
	mine.DPIx = 96
	mine.DPIy = 96
	if mine.Type == OCRTypeMath {
		var config = new(ConfigInfo)
		config.Math.Solver.Enable = false
		mine.Config = config
	}else{
		mine.Config = nil
	}
	mine.Strokes = make([]*StrokeGroup, 0, 10)
	group := new(StrokeGroup)
	group.PenStyle = nil
	mine.Strokes = append(mine.Strokes, group)
	length := len(mine.dots)
	group.Strokes = make([]*Stroke, 0, 10)
	point := mine.basePoint()
	stroke := mine.NewStroke()
	group.Strokes = append(group.Strokes, stroke)
	for i := 0; i < length; i += 1 {
		dot := mine.dots[i]
		s := group.Strokes[len(group.Strokes) - 1]
		s.Xs = append(s.Xs, float64(dot.X*100+dot.FX)-point.X)
		s.Ys = append(s.Ys, float64(dot.Y*100+dot.FY)-point.Y)
		s.Ts = append(s.Ts, dot.stamp)
		if dot.Action == DotActionUp && i < length - 1 {
			tmp := mine.NewStroke()
			group.Strokes = append(group.Strokes, tmp)
		}
	}
	data, _ := json.Marshal(mine)
	return data
}

func (mine *RecognizerInfo) basePoint() *Point {
	var length = len(mine.dots)
	if length < 1 {
		return nil
	}
	var point = new(Point)
	point.X = float64(mine.dots[0].X*100 + mine.dots[0].FX)
	point.Y = float64(mine.dots[0].Y*100 + mine.dots[0].FY)
	for i := 0; i < length; i += 1 {
		dot := mine.dots[i]
		xx := float64(dot.X*100 + dot.FX)
		yy := float64(dot.Y*100 + dot.FY)
		if point.X > xx {
			point.X = xx
		}
		if point.Y > yy {
			point.Y = yy
		}
	}
	return point
}
