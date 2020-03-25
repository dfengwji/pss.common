package writing

import (
	"strconv"
	"strings"
)

const (
	//从客户端收到的消息
	MessageReqConnect = "0001"
	//连接成功发送给客户端的消息，客户端开始发送数据
	MessageRespConnect = "0002"
	//接收坐标数据
	MessageReqData  = "0003"
	MessageRespData = "0004"
	//心跳
	MessageReqHeart  = "0005"
	MessageRespHeart = "0006"
	//发送坐标数据完成
	MessageReqComplete = "0007"
)

type HeaderInfo struct {
	Seq string `json:"-"`
	//用户ID
	UID string `json:"-"`
	//笔的ID, 笔盒里面预留字段
	SSRC   string `json:"-"`
	Type   string `json:"-"`
	userID uint64 `json:"-"`
	penID  uint64 `json:"-"`
}

func (mine *HeaderInfo) GetHex() string {
	return mine.Seq + mine.UID + mine.SSRC + mine.Type
}

func (mine *HeaderInfo) String() string {
	return strings.Join([]string{"seq", mine.Seq, "uid", mine.UID, "ssrc", mine.SSRC, "msg id", mine.Type}, ",")
}

func (mine *HeaderInfo) ParseHex(hex string) {
	if len(hex) < 28 {
		return
	}
	mine.Seq = hex[0:8]
	mine.UID = hex[8:16]
	mine.SSRC = hex[16:24]
	mine.Type = hex[24:28]
	mine.userID, _ = strconv.ParseUint(mine.UID, 16, 64)
	mine.penID, _ = strconv.ParseUint(mine.SSRC, 16, 64)
}

func (mine *HeaderInfo) User() uint64 {
	return mine.userID
}

func (mine *HeaderInfo) Pen() uint64 {
	return mine.penID
}
