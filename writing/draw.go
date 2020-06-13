package writing

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/cmplx"

	"github.com/fogleman/gg"
)

const (
	ColorBlue   = 1
	ColorGreen  = 2
	ColorCyan   = 3
	ColorRed    = 4
	ColorPurple = 5
	ColorYellow = 6
	ColorBlack  = 7
)

const (
	ColorHexBlue   = "#4082E3"
	ColorHexGreen  = "#65F44D"
	ColorHexCyan   = "#40E1D8"
	ColorHexRed    = "#EF3737"
	ColorHexPurple = "#CF0070"
	ColorHexYellow = "#FFDA57"
	ColorHexBlack  = "#0A0A0A"
)

type Vector2 struct {
	X float64
	Y float64
}

var (
	red    color.Color = color.RGBA{239, 55, 55, 255}
	blue   color.Color = color.RGBA{64, 130, 227, 255}
	green  color.Color = color.RGBA{101, 244, 77, 255}
	black  color.Color = color.RGBA{10, 10, 10, 255}
	cyan   color.Color = color.RGBA{64, 225, 216, 255}
	purple color.Color = color.RGBA{207, 0, 112, 255}
	yellow color.Color = color.RGBA{255, 218, 87, 255}
)

var paperOffset = Vector2{X: 0, Y: 0}
var codepointX = 1.524
var codepointY = 1.524
var lastPoint = new(Vector2)

var pointIndex = -1

var gX0, gX1, gX2, gX3 float64
var gY0, gY1, gY2, gY3 float64
var gP0, gP1, gP2, gP3 float64
var gVx01, gVy01, gNX0, gNY0 float64
var gVx21, gVy21 float64
var gNorm float64
var gNX2, gNY2 float64

// 渲染选项
type RenderOptions struct {
	PaperWidth  int     // 纸张宽度, 单位为毫米
	PaperHeight int     // 纸张高度, 单位为毫米
	ImageWidth  int     // 图片宽度, 单位为像素
	ImageHeight int     // 图片高度, 单位为像素
	PenSize     float64 // 笔迹粗细， 单位为像素
	ClipSize    int     // 裁剪框，单位为像素，如果笔迹未超出裁剪线框，则按笔迹的包围盒开始裁剪，也就是放大了
	PaddingSize int     // 填充尺寸，裁剪后，填充空白边框的尺寸
}

// 渲染的点
type RenderPoint struct {
	Dot *DotInfo // 原始点
	X   float64  // 实际渲染坐标X
	Y   float64  // 实际渲染坐标Y
}

func DrawPoints(points []*DotInfo, canvasSize Vector2, paperSize Vector2, path string) error {
	if points == nil {
		return errors.New("the points is nil")
	}
	length := len(points)
	if length < 1 {
		return errors.New("the points is empty")
	}
	//t := fmt.Sprintf("try draw points num = %d that uid = %s", length, uid)
	var canvas *gg.Context
	canvas = gg.NewContext(int(canvasSize.X), int(canvasSize.Y))
	canvas.SetRGB(0, 0, 0)
	canvas.SetLineCapRound()
	isUp := false
	for i := 0; i < length; i++ {
		up := drawGraph(canvas, points[i], canvasSize, paperSize)
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

func DrawImage(points []*DotInfo, canvasSize Vector2, paperSize Vector2, path string) (*bytes.Buffer, error) {
	if points == nil {
		return nil, errors.New("the points is nil")
	}
	length := len(points)
	if length < 1 {
		return nil, errors.New("the points is empty")
	}
	//t := fmt.Sprintf("try draw points num = %d that uid = %s", length, uid)
	var canvas *gg.Context
	canvas = gg.NewContext(int(canvasSize.X), int(canvasSize.Y))
	canvas.SetRGB(0, 0, 0)
	canvas.SetLineCapRound()
	isUp := false
	for i := 0; i < length; i++ {
		up := drawGraph(canvas, points[i], canvasSize, paperSize)
		if up {
			isUp = true
		}
	}

	if isUp {
		buf := bytes.NewBuffer(nil)
		err := png.Encode(buf, canvas.Image())
		if err != nil {
			return nil, err
		}
		if len(path) > 0 {
			err2 := canvas.SavePNG(path)
			if err2 != nil {
				return nil, err2
			}
		}
		return buf, nil
	}
	return nil,errors.New("the up action not found")
}

func DrawPointSamples(points []*PointSample, size Vector2, paper Vector2) {
	path := "files/images/test-3.png"
	var canvas *gg.Context
	canvas = gg.NewContext(int(size.X), int(size.Y))
	canvas.SetRGB(0, 0, 0)
	canvas.SetLineCapRound()
	//imageNum += 1
	for i := 0; i < len(points); i++ {
		drawGraph2(canvas, points[i], size, paper)
	}

	err2 := canvas.SavePNG(path)
	if err2 != nil {
		fmt.Println(err2.Error())
	}
}

func SavePNG(_dots []*DotInfo, _options RenderOptions, _filepath string) error {
	length := len(_dots)
	if length < 1 {
		return errors.New("none points")
	}

	// 将纸张转换为等比适配图片大小的尺寸时,使用FitCenter模式,保留完整纸张，不切割和局部展示
	// 以纸张的比例和画布的宽度计算画布的等比高度
	canvasWidth := int(float64(_options.ImageHeight) * float64(_options.PaperWidth) / float64(_options.PaperHeight))
	// 以纸张的比例和画布的高度计算画布的等比宽度
	canvasHeight := int(float64(_options.ImageWidth) * float64(_options.PaperHeight) / float64(_options.PaperWidth))
	// 适配后，X轴或Y轴会存在空白，设定X轴和Y轴的偏移量，将纸张移动图片中央
	offsetX := float64(0)
	offsetY := float64(0)

	// 图片宽度比大于纸张宽度比
	// ....._____________.....
	// .    |           |    .
	// .    |           |    .
	// .    |           |    .
	// .    |           |    .
	// .    |  Paper    |    .  Image
	// .    |           |    .
	// .    |           |    .
	// .    |           |    .
	// ....._____________.....

	//fmt.Println(fmt.Sprintf("canvas : {width:%v, height:%v}", canvasWidth, canvasHeight))
	if (float64(_options.PaperHeight)/float64(_options.PaperWidth))-(float64(_options.ImageHeight)/float64(_options.ImageWidth)) > 0 {
		// 画布高度为图片高度
		canvasHeight = _options.ImageHeight
		// 计算X轴偏移量
		offsetX = float64(_options.ImageWidth-canvasWidth) / 2
	}

	// 图片宽度比大于纸张宽度比
	// .................
	// .               .
	// .               .
	// _________________
	// |               |
	// |               |
	// |               |
	// |               |
	// |     Paper     |
	// |               |
	// |               |
	// |               |
	// _________________
	// .               .
	// .     Image     .
	// .................

	if (float64(_options.ImageHeight)/float64(_options.ImageWidth))-(float64(_options.PaperHeight)/float64(_options.PaperWidth)) > 0 {
		// 画布宽度为图片宽度
		canvasWidth = _options.ImageWidth
		// 计算X轴偏移量
		offsetY = float64(_options.ImageHeight-canvasHeight) / 2
	}

	//fmt.Println(fmt.Sprintf("fit canvas : {width:%v, height:%v}", canvasWidth, canvasHeight))

	// 点的坐标系的左上角为原点，右下角为无限大
	// 画布上的点的包围框，左上角为(minX,minY)，右下角为(maxX, maxY)
	// !!! 初始化时，将小值设置为最大，大值设置为最小, 不是BUG，不要改我
	minX := float64(9223372036854775807)
	maxX := float64(0)
	minY := float64(9223372036854775807)
	maxY := float64(0)

	points := make([]*RenderPoint, len(_dots))
	for i := 0; i < length; i++ {
		points[i] = &RenderPoint{
			Dot: _dots[i],
		}
		dot := _dots[i]
		if dot == nil {
			continue
		}
		//fmt.Println(dot)

		// FX代表小数部分，X代表整数部分
		// 计算出点的浮点型的坐标
		dotX := float64(dot.FX)/100.0 + float64(dot.X)
		dotY := float64(dot.FY)/100.0 + float64(dot.Y)

		//fmt.Println(fmt.Sprintf("dotX: %v, dotY: %v", dotX, dotY))

		// codepint 表示笔触在纸上的点的大小
		// 将纸张大小按笔触单位为1进行标准化
		paperNormalizeWidth := float64(_options.PaperWidth) / codepointX
		paperNormalizeHeight := float64(_options.PaperHeight) / codepointY

		// 计算像素和物理单位的比例
		scaleX := float64(canvasWidth) / paperNormalizeWidth
		scaleY := float64(canvasWidth) / paperNormalizeHeight

		// 将笔的物理单位转化为像素
		x := dotX * scaleX
		y := dotY * scaleY

		// 添加偏移量，使其居中
		x = x + offsetX
		y = y + offsetY

		points[i].X = x
		points[i].Y = y

		// 矫正包围框
		if x > maxX {
			maxX = x
		}
		if x < minX {
			minX = x
		}
		if y > maxY {
			maxY = y
		}
		if y < minY {
			minY = y
		}
	}

	// 判断点的包围框是否未超出裁剪框
	if int(minX) > _options.ClipSize && int(maxY) < _options.ClipSize && int(minX) > _options.ClipSize && int(maxY) < _options.ClipSize {
		// 裁剪
		// 点的包围盒的高度和宽度
		//boundWidth := maxX - minX
		//boundHeight := maxY - minY
	}

	// 新建画布
	newCanvas := gg.NewContext(_options.ImageWidth, _options.ImageHeight)
	newCanvas.SetLineCapRound()
	// 设置白色
	newCanvas.SetRGB(1, 1, 1)
	// 使用设置的颜色清空画布
	newCanvas.Clear()
	// 渲染点
	for i := 0; i < len(points); i++ {
		point := points[i]
		if point.Dot == nil {
			continue
		}
		color := getColor(point.Dot.Color)
		penSize := float64(point.Dot.Scale) * _options.PenSize
		if point.Dot.Action == DotActionDown {
			touchDown(newCanvas, point.X, point.Y, color, penSize, 1.0, point.Dot.Force)
		} else if point.Dot.Action == DotActionMove {
			touchMove(newCanvas, point.X, point.Y, color, penSize, 1.0, point.Dot.Force)
		} else if point.Dot.Action == DotActionUp {
			touchUp(newCanvas, point.X, point.Y, color, penSize)
		}
	}
	return newCanvas.SavePNG(_filepath)
}

func drawGraph(canvas *gg.Context, point *DotInfo, size Vector2, paper Vector2) bool {
	if point == nil {
		return false
	}

	coordinateX := float64(point.FX)/100.0 + float64(point.X)
	coordinateY := float64(point.FY)/100.0 + float64(point.Y)
	xx := coordinateX * size.X
	ax := paper.X / codepointX
	px := xx / ax
	x := px + paperOffset.X
	//x := roundNum(px, 13)

	yy := coordinateY * size.Y
	ay := paper.Y / codepointY
	py := yy / ay
	y := py + paperOffset.Y
	color := getColor(point.Color)
	//y := roundNum(py, 13)
	//笔迹粗细
	gWidth := float64(point.Scale * 2.0)
	if point.Action == DotActionDown {
		touchDown(canvas, x, y, color, gWidth, 1.0, point.Force)
	} else if point.Action == DotActionMove {
		touchMove(canvas, x, y, color, gWidth, 1.0, point.Force)
	} else if point.Action == DotActionUp {
		touchUp(canvas, x, y, color, gWidth)
		return true
	}
	return false
}

func drawGraph2(canvas *gg.Context, point *PointSample, size Vector2, paper Vector2) {
	if point == nil {
		return
	}

	xx := point.X * size.X
	ax := paper.X / codepointX
	px := xx / ax
	x := px + paperOffset.X
	//x := roundNum(px, 13)

	yy := point.Y * size.Y
	ay := paper.Y / codepointY
	py := yy / ay
	y := py + paperOffset.Y
	//y := roundNum(py, 13)

	if point.Action == DotActionDown {
		touchDown(canvas, x, y, point.Color, 1.0, point.Scale, point.Force)
	} else if point.Action == DotActionMove {
		touchMove(canvas, x, y, point.Color, 1.0, point.Scale, point.Force)
	} else if point.Action == DotActionUp {
		touchUp(canvas, x, y, point.Color, 1.0)
	}
}

func getColor(color uint32) color.Color {
	if color == ColorBlue {
		return blue
	} else if color == ColorBlack {
		return black
	} else if color == ColorCyan {
		return cyan
	} else if color == ColorGreen {
		return green
	} else if color == ColorPurple {
		return purple
	} else if color == ColorYellow {
		return yellow
	} else if color == ColorRed {
		return red
	} else {
		return black
	}
}

func touchDown(_canvas *gg.Context, x float64, y float64, color color.Color, _penSize float64, scale uint8, force uint16) {
	pointIndex = 0

	_canvas.SetColor(color)
	drawPen(_canvas, float64(scale), 0.0, 0.0, _penSize, x, y, force, DotActionDown)
	//canvas.DrawPoint(x, y, 1)
	_canvas.Stroke()
	lastPoint.X = x
	lastPoint.Y = y
}

func touchMove(_canvas *gg.Context, x float64, y float64, color color.Color, _penSize float64, scale uint8, force uint16) {
	pointIndex += 1
	//fmt.Println()
	//fmt.Printf("x = %f; y = %f", x, y)
	_canvas.SetColor(color)
	drawPen(_canvas, float64(scale), 0.0, 0.0, _penSize, x, y, force, DotActionMove)
	//canvas.SetLineWidth(lineWidth)
	//canvas.DrawLine(lastPoint.X, lastPoint.Y, x, y)
	_canvas.Stroke()
	lastPoint.X = x
	lastPoint.Y = y
}

func touchUp(_canvas *gg.Context, x float64, y float64, color color.Color, _penSize float64) {
	pointIndex += 1
	_canvas.SetColor(color)
	//canvas.SetLineWidth(lineWidth)
	//canvas.DrawLine(lastPoint.X, lastPoint.Y, x, y)
	drawPen(_canvas, 1.0, 0.0, 0.0, _penSize, x, y, 0, DotActionUp)
	_canvas.Stroke()
	lastPoint.X = x
	lastPoint.Y = y
	pointIndex = -1
}

func drawPen(_canvas *gg.Context, scale float64, offsetX float64, offsetY float64, penWidth float64, x float64, y float64, force uint16, ntype uint8) {
	//DV.paint.setStrokeCap(Paint.Cap.ROUND)
	//DV.paint.setStyle(Paint.Style.FILL)
	ws := 1.0
	wh := ws * 20
	if pointIndex == 0 { //down
		gX0 = x*scale + offsetX + 0.1
		gY0 = y*scale + offsetY
		//g_p0 = Math.max(1, penWidth * 3 * force / 1023) * scale;
		gP0 = getPenWidth(penWidth, force) * ws
		_canvas.DrawPoint(gX0, gY0, 0.5)
		return
	}

	if pointIndex == 1 {
		gX1 = x*scale + offsetX + 0.1
		gY1 = y*scale + offsetY
		//g_p1 = Math.max(1, penWidth * 3 * force / 1023) * scale;
		gP1 = getPenWidth(penWidth, force) * ws

		gVx01 = gX1 - gX0
		gVy01 = gY1 - gY0
		// instead of dividing tangent/norm by two, we multiply norm by 2
		gNorm = math.Sqrt(gVx01*gVx01+gVy01*gVy01+0.0001) * wh
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
		gP3 = getPenWidth(penWidth, force) * ws

		gX2 = (gX1 + gX3) / 2
		gY2 = (gY1 + gY3) / 2
		gP2 = (gP1 + gP3) / 2
		gVx21 = gX1 - gX2
		gVy21 = gY1 - gY2
		gNorm = math.Sqrt(gVx21*gVx21+gVy21*gVy21+0.0001) * wh

		gVx21 = gVx21 / gNorm * gP2
		gVy21 = gVy21 / gNorm * gP2
		gNX2 = -gVy21
		gNY2 = gVx21

		_canvas.MoveTo(gX0+gNX0, gY0+gNY0)
		// The + boundary of the stroke
		_canvas.CubicTo(gX1+gNX0, gY1+gNY0, gX1+gNX2, gY1+gNY2, gX2+gNX2, gY2+gNY2)
		// round out the cap
		_canvas.CubicTo(gX2+gNX2-gVx21, gY2+gNY2-gVy21, gX2-gNX2-gVx21, gY2-gNY2-gVy21, gX2-gNX2, gY2-gNY2)
		// THe - boundary of the stroke
		_canvas.CubicTo(gX1-gNX2, gY1-gNY2, gX1-gNX0, gY1-gNY0, gX0-gNX0, gY0-gNY0)
		// round out the other cap
		_canvas.CubicTo(gX0-gNX0-gVx01, gY0-gNY0-gVy01, gX0+gNX0-gVx01, gY0+gNY0-gVy01, gX0+gNX0, gY0+gNY0)

		if ntype == 2 {
			_canvas.SetLineWidth(gP3)
			_canvas.DrawLine(gX1, gY1, gX3, gY3)
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
		gP2 = getPenWidth(penWidth, force) * ws
		gVx21 = gX1 - gX2
		gVy21 = gY1 - gY2
		gNorm = math.Sqrt(gVx21*gVx21+gVy21*gVy21+0.0001) * wh
		gVx21 = gVx21 / gNorm * gP2
		gVy21 = gVy21 / gNorm * gP2
		gNX2 = -gVy21
		gNY2 = gVx21

		_canvas.MoveTo(gX0+gNX0, gY0+gNY0)
		_canvas.CubicTo(gX1+gNX0, gY1+gNY0, gX1+gNX2, gY1+gNY2, gX2+gNX2, gY2+gNY2)
		_canvas.CubicTo(gX2+gNX2-gVx21, gY2+gNY2-gVy21, gX2-gNX2-gVx21, gY2-gNY2-gVy21, gX2-gNX2, gY2-gNY2)
		_canvas.CubicTo(gX1-gNX2, gY1-gNY2, gX1-gNX0, gY1-gNY0, gX0-gNX0, gY0-gNY0)
		_canvas.CubicTo(gX0-gNX0-gVx01, gY0-gNY0-gVy01, gX0+gNX0-gVx01, gY0+gNY0-gVy01, gX0+gNX0, gY0+gNY0)
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

func fractal2(img *image.RGBA, size Vector2) {
	dx := size.X
	dy := size.Y
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
