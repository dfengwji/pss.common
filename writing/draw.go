package writing

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/fogleman/gg"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/cmplx"
)

const (
	ColorBlue = 1
	ColorGreen = 2
	ColorCyan = 3
	ColorRed = 4
	ColorPurple = 5
	ColorYellow = 6
	ColorBlack = 7
)

const (
	ColorHexBlue = "#4082E3"
	ColorHexGreen = "#65F44D"
	ColorHexCyan = "#40E1D8"
	ColorHexRed = "#EF3737"
	ColorHexPurple = "#CF0070"
	ColorHexYellow = "#FFDA57"
	ColorHexBlack= "#0A0A0A"
)

type Vector2 struct {
	X float64
	Y float64
}

var (
	red   color.Color = color.RGBA{239, 55, 55, 255}
	blue  color.Color = color.RGBA{64, 130, 227, 255}
	green color.Color = color.RGBA{101, 244, 77, 255}
	black color.Color = color.RGBA{10, 10, 10, 255}
	cyan  color.Color = color.RGBA{64, 225, 216, 255}
	purple color.Color = color.RGBA{207, 0, 112, 255}
	yellow color.Color = color.RGBA{255, 218, 87, 255}
)

var canvasSize = Vector2{X: 600, Y: 844}
var paperSize = Vector2{X: 290, Y: 210}
var paperOffset = Vector2{X: 0, Y: 0}
var codepointX = 1.524
var codepointY = 1.524
var lastPoint = new(Vector2)

var canvas *gg.Context
var imageNum = 0
var pointIndex = -1
var gWidth = 1.0 //笔迹粗细

var gX0, gX1, gX2, gX3 float64
var gY0, gY1, gY2, gY3 float64
var gP0, gP1, gP2, gP3 float64
var gVx01, gVy01, gNX0, gNY0 float64
var gVx21, gVy21 float64
var gNorm float64
var gNX2, gNY2 float64

func DrawPoints(points []*DotInfo, uid string) error {
	if points == nil {
		return errors.New("the points is nil")
	}
	length := len(points)
	if length < 1 {
		return errors.New("the points is empty")
	}
	//t := fmt.Sprintf("try draw points num = %d that uid = %s", length, uid)
	path := fmt.Sprintf("files/images/img_%s.png", uid)
	if canvas == nil {
		canvas = gg.NewContext(int(canvasSize.X), int(canvasSize.Y))
		canvas.SetRGB(0, 0, 0)
		canvas.SetLineCapRound()
	}else{
		canvas = gg.NewContext(int(canvasSize.X), int(canvasSize.Y))
		canvas.SetRGB(0, 0, 0)
		canvas.SetLineCapRound()
	}
	imageNum += length
	isUp := false
	for i := 0; i < length; i++ {
		up := drawGraph(points[i])
		if up {
			isUp = true
		}
	}

	if isUp {
		err2 := canvas.SavePNG(path)
		if err2 != nil {
			return err2
		}
		buf := new(bytes.Buffer)
		err := png.Encode(buf, canvas.Image())
		if err != nil {
			return err
		}
	}
	return nil
}

func DrawPointSamples(points []*PointSample) {
	path := "files/images/test-3.png"
	if canvas == nil {
		canvas = gg.NewContext(int(canvasSize.X), int(canvasSize.Y))
		canvas.SetRGB(0, 0, 0)
		canvas.SetLineCapRound()
	}
	//imageNum += 1
	for i := 0; i < len(points); i++ {
		drawGraph2(points[i])
	}

	err2 := canvas.SavePNG(path)
	if err2 != nil {
		fmt.Println(err2.Error())
	}
}

func drawGraph(point *DotInfo) bool {
	if point == nil {
		return false
	}

	coordinateX := float64(point.FX)/100.0 + float64(point.X)
	coordinateY := float64(point.FY)/100.0 + float64(point.Y)
	xx := coordinateX * canvasSize.X
	ax := paperSize.X / codepointX
	px := xx / ax
	x := px + paperOffset.X
	//x := roundNum(px, 13)

	yy := coordinateY * canvasSize.Y
	ay := paperSize.Y / codepointY
	py := yy / ay
	y := py + paperOffset.Y
	color := getColor(point.Color)
	//y := roundNum(py, 13)
	if point.Action == DotActionDown {
		touchDown(x, y, color, point.Scale, point.Force)
	} else if point.Action == DotActionMove {
		touchMove(x, y, color, point.Scale, point.Force)
	} else if point.Action == DotActionUp {
		touchUp(x, y, color)
		return true
	}
	return false
}

func drawGraph2(point *PointSample) {
	if point == nil {
		return
	}

	xx := point.X * canvasSize.X
	ax := paperSize.X / codepointX
	px := xx / ax
	x := px + paperOffset.X
	//x := roundNum(px, 13)

	yy := point.Y * canvasSize.Y
	ay := paperSize.Y / codepointY
	py := yy / ay
	y := py + paperOffset.Y
	//y := roundNum(py, 13)

	if point.Action == DotActionDown {
		touchDown(x, y, point.Color, point.Scale, point.Force)
	} else if point.Action == DotActionMove {
		touchMove(x, y, point.Color, point.Scale, point.Force)
	} else if point.Action == DotActionUp {
		touchUp(x, y, point.Color)
	}
}

func getColor(color uint32) color.Color {
	if color == ColorBlue {
		return blue
	}else if color == ColorBlack {
		return black
	}else if color == ColorCyan {
		return cyan
	}else if color == ColorGreen {
		return green
	}else  if color == ColorPurple {
		return purple
	}else if color == ColorYellow {
		return yellow
	}else if color == ColorRed {
		return red
	}else{
		return black
	}
}

func touchDown(x float64, y float64, color color.Color,scale uint8, force uint16) {
	pointIndex = 0

	canvas.SetColor(color)
	drawPen(float64(scale), 0.0, 0.0, gWidth, x, y, force, DotActionDown)
	//canvas.DrawPoint(x, y, 1)
	canvas.Stroke()
	lastPoint.X = x
	lastPoint.Y = y
}

func touchMove(x float64, y float64,color color.Color,scale uint8, force uint16) {
	pointIndex += 1
	//fmt.Println()
	//fmt.Printf("x = %f; y = %f", x, y)
	canvas.SetColor(color)
	drawPen(float64(scale), 0.0, 0.0, gWidth, x, y, force, DotActionMove)
	//canvas.SetLineWidth(lineWidth)
	//canvas.DrawLine(lastPoint.X, lastPoint.Y, x, y)
	canvas.Stroke()
	lastPoint.X = x
	lastPoint.Y = y
}

func touchUp(x float64, y float64, color color.Color) {
	pointIndex += 1
	canvas.SetColor(color)
	//canvas.SetLineWidth(lineWidth)
	//canvas.DrawLine(lastPoint.X, lastPoint.Y, x, y)
	drawPen(1.0, 0.0, 0.0, gWidth, x, y, 0, DotActionUp)
	canvas.Stroke()
	lastPoint.X = x
	lastPoint.Y = y
	pointIndex = -1
}

func drawPen(scale float64, offsetX float64, offsetY float64, penWidth float64, x float64, y float64, force uint16, ntype uint8) {
	//DV.paint.setStrokeCap(Paint.Cap.ROUND)
	//DV.paint.setStyle(Paint.Style.FILL)
	if pointIndex == 0 { //down
		gX0 = x*scale + offsetX + 0.1
		gY0 = y*scale + offsetY
		//g_p0 = Math.max(1, penWidth * 3 * force / 1023) * scale;
		gP0 = getPenWidth(penWidth, force) * scale
		canvas.DrawPoint(gX0, gY0, 0.5)
		return
	}

	if pointIndex == 1 {
		gX1 = x*scale + offsetX + 0.1
		gY1 = y*scale + offsetY
		//g_p1 = Math.max(1, penWidth * 3 * force / 1023) * scale;
		gP1 = getPenWidth(penWidth, force) * scale

		gVx01 = gX1 - gX0
		gVy01 = gY1 - gY0
		// instead of dividing tangent/norm by two, we multiply norm by 2
		gNorm = math.Sqrt(gVx01*gVx01+gVy01*gVy01+0.0001) * 2
		gVx01 = gVx01 / gNorm * gP0
		gVy01 = gVy01 / gNorm * gP0
		gNX0 = gVy01
		gNY0 = -gVx01
		return
	}

	if pointIndex > 1 && pointIndex < 10000 {
		// (x0,y0) and (x2,y2) are midpoints, (x1,y1) and (x3,y3) are actual
		gX3 = x*scale + offsetX + 0.1
		gY3 = y*scale + offsetY
		//g_p3 = Math.max(1, penWidth * 3 * force / 1023) * scale;
		gP3 = getPenWidth(penWidth, force) * scale

		gX2 = (gX1 + gX3) / 2
		gY2 = (gY1 + gY3) / 2
		gP2 = (gP1 + gP3) / 2
		gVx21 = gX1 - gX2
		gVy21 = gY1 - gY2
		gNorm = math.Sqrt(gVx21*gVx21+gVy21*gVy21+0.0001) * 2

		gVx21 = gVx21 / gNorm * gP2
		gVy21 = gVy21 / gNorm * gP2
		gNX2 = -gVy21
		gNY2 = gVx21

		canvas.MoveTo(gX0+gNX0, gY0+gNY0)
		// The + boundary of the stroke
		canvas.CubicTo(gX1+gNX0, gY1+gNY0, gX1+gNX2, gY1+gNY2, gX2+gNX2, gY2+gNY2)
		// round out the cap
		canvas.CubicTo(gX2+gNX2-gVx21, gY2+gNY2-gVy21, gX2-gNX2-gVx21, gY2-gNY2-gVy21, gX2-gNX2, gY2-gNY2)
		// THe - boundary of the stroke
		canvas.CubicTo(gX1-gNX2, gY1-gNY2, gX1-gNX0, gY1-gNY0, gX0-gNX0, gY0-gNY0)
		// round out the other cap
		canvas.CubicTo(gX0-gNX0-gVx01, gY0-gNY0-gVy01, gX0+gNX0-gVx01, gY0+gNY0-gVy01, gX0+gNX0, gY0+gNY0)

		if ntype == 2 {
			canvas.SetLineWidth(gP3)
			canvas.DrawLine(gX1, gY1, gX3, gY3)
		}
		gX0 = gX2
		gY0 = gY2
		gP0 = gP2
		gX1 = gX3
		gY1 = gY3
		gP1 = gP3
		gVx01 = -gVx21
		gVy01 = -gVy21
		gNX0 = gNX2
		gNY0 = gNY2
		return
	}
	if pointIndex >= 10000 { //Last Point
		gX2 = x*scale + offsetX + 0.1
		gY2 = y*scale + offsetY
		//g_p2 = Math.max(1, penWidth * 3 * force / 1023) * scale;
		gP2 = getPenWidth(penWidth, force) * scale
		gVx21 = gX1 - gX2
		gVy21 = gY1 - gY2
		gNorm = math.Sqrt(gVx21*gVx21+gVy21*gVy21+0.0001) * 2
		gVx21 = gVx21 / gNorm * gP2
		gVy21 = gVy21 / gNorm * gP2
		gNX2 = -gVy21
		gNY2 = gVx21

		canvas.MoveTo(gX0+gNX0, gY0+gNY0)
		canvas.CubicTo(gX1+gNX0, gY1+gNY0, gX1+gNX2, gY1+gNY2, gX2+gNX2, gY2+gNY2)
		canvas.CubicTo(gX2+gNX2-gVx21, gY2+gNY2-gVy21, gX2-gNX2-gVx21, gY2-gNY2-gVy21, gX2-gNX2, gY2-gNY2)
		canvas.CubicTo(gX1-gNX2, gY1-gNY2, gX1-gNX0, gY1-gNY0, gX0-gNX0, gY0-gNY0)
		canvas.CubicTo(gX0-gNX0-gVx01, gY0-gNY0-gVy01, gX0+gNX0-gVx01, gY0+gNY0-gVy01, gX0+gNX0, gY0+gNY0)
		return
	}
}

func getPenWidth(width float64, force uint16) float64 {
	var m = 1.0
	if width == 1 {
		if force <= 50 {
			m = 0.8
		} else if force > 50 && force <= 90 {
			m = 1.0
		} else if force > 90 && force <= 120 {
			m = 1.2
		} else if force > 120 && force <= 150 {
			m = 1.4
		} else if force > 150 && force <= 190 {
			m = 1.6
		} else if force > 190 && force <= 210 {
			m = 1.8
		} else if force > 210 && force <= 330 {
			m = 1.9
		} else if force > 330 && force <= 500 {
			m = 2.0
		} else if force > 500 && force <= 650 {
			m = 2.1
		} else if force > 650 && force <= 800 {
			m = 2.2
		} else if force > 800 {
			m = 2.4
		}
	} else if width == 2 {
		if force <= 50 {
			m = 1.6
		} else if force > 50 && force <= 90 {
			m = 2.0
		} else if force > 90 && force <= 120 {
			m = 2.4
		} else if force > 120 && force <= 150 {
			m = 2.8
		} else if force > 150 && force <= 190 {
			m = 3.2
		} else if force > 190 && force <= 210 {
			m = 3.6
		} else if force > 210 && force <= 330 {
			m = 3.8
		} else if force > 330 && force <= 500 {
			m = 4.0
		} else if force > 500 && force <= 650 {
			m = 4.2
		} else if force > 650 && force <= 800 {
			m = 4.4
		} else if force > 800 {
			m = 4.8
		}
	} else if width == 3 {
		if force <= 50 {
			m = 2.4
		} else if force > 50 && force <= 90 {
			m = 3.0
		} else if force > 90 && force <= 120 {
			m = 3.6
		} else if force > 120 && force <= 150 {
			m = 4.2
		} else if force > 150 && force <= 190 {
			m = 4.8
		} else if force > 190 && force <= 210 {
			m = 5.4
		} else if force > 210 && force <= 330 {
			m = 5.7
		} else if force > 330 && force <= 500 {
			m = 6.0
		} else if force > 500 && force <= 650 {
			m = 6.3
		} else if force > 650 && force <= 800 {
			m = 6.6
		} else if force > 800 {
			m = 7.2
		}
	} else if width == 4 {
		if force <= 50 {
			m = 3.2
		} else if force > 50 && force <= 90 {
			m = 4.0
		} else if force > 90 && force <= 120 {
			m = 4.8
		} else if force > 120 && force <= 150 {
			m = 5.6
		} else if force > 150 && force <= 190 {
			m = 6.4
		} else if force > 190 && force <= 210 {
			m = 7.2
		} else if force > 210 && force <= 330 {
			m = 7.6
		} else if force > 330 && force <= 500 {
			m = 8.0
		} else if force > 500 && force <= 650 {
			m = 8.4
		} else if force > 650 && force <= 800 {
			m = 8.8
		} else if force > 800 {
			m = 9.6
		}
	} else if width == 5 {
		if force <= 50 {
			m = 4.0
		} else if force > 50 && force <= 90 {
			m = 5.0
		} else if force > 90 && force <= 120 {
			m = 6.0
		} else if force > 120 && force <= 150 {
			m = 7.0
		} else if force > 150 && force <= 190 {
			m = 8.0
		} else if force > 190 && force <= 210 {
			m = 9.0
		} else if force > 210 && force <= 330 {
			m = 9.5
		} else if force > 330 && force <= 500 {
			m = 10.0
		} else if force > 500 && force <= 650 {
			m = 10.5
		} else if force > 650 && force <= 800 {
			m = 11.0
		} else if force > 800 {
			m = 12.0
		}
	} else if width == 6 {
		if force <= 50 {
			m = 4.8
		} else if force > 50 && force <= 90 {
			m = 6.0
		} else if force > 90 && force <= 120 {
			m = 7.2
		} else if force > 120 && force <= 150 {
			m = 8.4
		} else if force > 150 && force <= 190 {
			m = 9.6
		} else if force > 190 && force <= 210 {
			m = 10.8
		} else if force > 210 && force <= 330 {
			m = 11.4
		} else if force > 330 && force <= 500 {
			m = 12.0
		} else if force > 500 && force <= 650 {
			m = 12.6
		} else if force > 650 && force <= 800 {
			m = 13.2
		} else if force > 800 {
			m = 14.4
		}
	} else {
		if force <= 50 {
			m = 3.2
		} else if force > 50 && force <= 90 {
			m = 4.0
		} else if force > 90 && force <= 120 {
			m = 4.8
		} else if force > 120 && force <= 150 {
			m = 5.6
		} else if force > 150 && force <= 190 {
			m = 6.4
		} else if force > 190 && force <= 210 {
			m = 7.2
		} else if force > 210 && force <= 330 {
			m = 7.6
		} else if force > 330 && force <= 500 {
			m = 8.0
		} else if force > 500 && force <= 650 {
			m = 8.4
		} else if force > 650 && force <= 800 {
			m = 8.8
		} else if force > 800 {
			m = 9.6
		}
	}
	return m
}

func roundNum(number float64, fractionDigits float64) float64 {
	return math.Round(number*math.Pow(10, fractionDigits)) / math.Pow(10, fractionDigits)
}

func abs(x int) int {
	if x >= 0 {
		return x
	}
	return -x
}

func fractal(img *image.RGBA) {
	dx := img.Bounds().Max.X
	dy := img.Bounds().Max.Y
	for x := 0; x < dx; x++ {
		for y := 0; y < dy; y++ {
			z := mandelbrot(complex(
				(float64(y)/float64(dy)*3 - 2.15),
				(float64(x)/float64(dx)*3 - 1.5),
			))
			img.Set(x, y, color.RGBA{-z ^ 2, -z ^ 2, -z ^ 2, 255})
		}
	}
}

func fractal2(img *image.RGBA) {
	dx := canvasSize.X
	dy := canvasSize.Y
	userX := -0.794591379577363
	userY := 0.16093921135504
	zoom := int64(9990999900)
	iterations := 255
	xShift := float64(dx / 2)
	yShift := float64(dy / 2)

	for v := 0; v < int(dy); v++ {
		for u := 0; u < int(dx); u++ {
			x := float64(u) - xShift
			y := (float64(v) * -1) + float64(dy) - yShift
			x = x + userX*float64(zoom)
			y = y + userY*float64(zoom)
			x = x / float64(zoom)
			y = y / float64(zoom)

			level := mandelb(x, y, iterations)
			if level == iterations {
				img.Set(u, v, black)
			} else {
				clr := 255 - uint8(level*255/iterations)
				img.Set(u, v, color.RGBA{clr, clr, clr, 255})
			}

		}
	}
}

func mandelb(x0, y0 float64, iter int) int {
	x := x0
	y := y0
	for i := 0; i < iter; i++ {
		real2 := x * x
		imag2 := y * y
		if (real2 + imag2) > 4.0 {
			return i
		}
		y = 2*x*y + y0
		x = real2 - imag2 + x0
	}
	return iter
}

func mandelbrot(in complex128) uint8 {
	n := in
	for i := uint8(0) + 255; i > 0; i-- {
		if cmplx.Abs(n) > 2 {
			return i
		}
		n = cmplx.Pow(n, complex(2, 0)) + in
	}
	return 255
}
