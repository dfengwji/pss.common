package writing

import (
	"github.com/pkg/errors"
	"strconv"
)

const KindChirography = "01"

type FragmentInfo struct {
	Seq     string
	Counter string
	Size    string
	Color	string
	kind    string
	/**
	 * 如果是命令,则这里是命令数据,就没有下面的数据点了
	 */
	Command string
	/**
	 * 存放解析好了的命令
	 */
	CommandData string

	/*
	书的实例ID
	 */
	BookInstance uint64
	/*
	铺码书本的ID
	 */
	BookTemplate uint64
	Section      string
	Owner        string
	Page         uint16
	Stamp        uint64
	Mac          string
	/**
	 * 笔迹点的个数,每个点是13(B)
	 */
	PointNum int
	Points   []*DotInfo
	hex      string
	dotHex	string
}

func formatString(msg string) (string,error) {
	cli, err := strconv.ParseInt(msg, 16, 64)
	if err != nil {
		return "",err
	}
	return strconv.FormatInt(cli, 16),nil
}

func hexToUnit(msg string) (uint64,error) {
	num, err := strconv.ParseUint(msg, 16, 64)
	if err != nil {
		return 0,err
	}
	return num,nil
}

func hexToUnit16(msg string) (uint16,error) {
	num, err := strconv.ParseUint(msg, 16, 16)
	if err != nil {
		return 0,err
	}
	return uint16(num),nil
}

func (mine *FragmentInfo) GetDots() []*DotInfo {
	list := make([]*DotInfo, 0, mine.PointNum)
	for i := 0; i < mine.PointNum; i += 1 {
		list = append(list, mine.Points[i])
	}
	return list
}

func (mine *FragmentInfo) Hex() string {
	return mine.hex
}

func (mine *FragmentInfo) DotsHex() string {
	return mine.dotHex
}

func (mine *FragmentInfo) SetDotsHex(str string) {
	mine.dotHex = str
}

func (mine *FragmentInfo) HasPenUp() bool {
	for i := 0; i < len(mine.Points); i += 1 {
		if mine.Points[i].Action == DotActionUp {
			return true
		}
	}
	return false
}

func (mine *FragmentInfo) ParseHex(msgType string, hex string) error {
	max := 46
	if msgType == MessageReqConnect {
		if len(hex) < max {
			return nil
		}
		mine.Counter,_ = formatString(hex[0:8])
		mine.Seq = hex[8:16]
		length := hex[34:38]
		size, _ := strconv.ParseInt(length, 16, 32)
		//wifi笔盒这里的token就是mac
		mine.Mac = hex[38 : 38+size*2]
	} else if msgType == MessageReqData {
		if len(hex) < max {
			return errors.New("the message format length is must more than 38")
		}
		mine.hex = hex
		mine.BookInstance, _ = hexToUnit(hex[0:8])
		mine.Counter,_ = formatString(hex[8:16])
		mine.Color = hex[16:20]
		mine.kind = hex[20:22]
		if mine.kind == KindChirography {
			mine.BookTemplate,_ = hexToUnit(hex[22:26])
			mine.Section,_ = formatString(hex[26:30])
			mine.Owner,_ = formatString(hex[30:34])
			mine.Page,_ = hexToUnit16(hex[34:36])
			mine.Stamp,_ = hexToUnit(hex[36:44])
			num, _ := strconv.ParseInt(hex[44:max], 16, 32)
			mine.PointNum = int(num)
			mine.Points = make([]*DotInfo, 0, mine.PointNum)
			sub := hex[max:]
			mine.dotHex = sub
			if len(sub) < mine.PointNum * DotHexLength {
				return errors.New("the points hex length is less than max length that num = " + strconv.FormatInt(num,10))
			}else{
				for i := 0; i < int(mine.PointNum); i += 1 {
					var point = new(DotInfo)
					err := point.ParseHex(sub[i*DotHexLength:(i+1)*DotHexLength])
					if err != nil {
						return err
					} else {
						mine.Points = append(mine.Points, point)
					}
				}
			}
		} else {
			//不是01,就代表命令
			mine.Command = hex[14 : 14+1]
		}
	} else if msgType == MessageReqHeart {

	}
	return nil
}
