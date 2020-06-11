package proxy

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type PrivateStyle struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	DefBook     string             `json:"book" bson:"book"`
	Page        uint16             `json:"page" bson:"page"`
	StartX      uint32             `json:"sx" bson:"sx"`
	StartY      uint32             `json:"sy" bson:"sy"`
	EndX        uint32             `json:"ex" bson:"ex"`
	EndY        uint32             `json:"ey" bson:"ey"`
	Background  string             `json:"background" bson:"background"`
	Exercise    string             `json:"exercise" bson:"exercise"`
}

func CreatePrivateStyle(info *PrivateStyle) error {
	_, err := insertOne(TablePrivateStyle, info)
	if err != nil {
		return err
	}
	return nil
}

func GetPrivateStyleNextID() uint64 {
	num, _ := getSequenceNext(TablePrivateStyle)
	return num
}

func GetPrivateStyles() ([]*PrivateStyle, error) {
	cursor, err1 := findAll(TablePrivateStyle, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*PrivateStyle, 0, 1000)
	for cursor.Next(context.Background()) {
		var node = new(PrivateStyle)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetPrivateStylesByBook(book string) ([]*PrivateStyle, error) {
	msg := bson.M{"book": book, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TablePrivateStyle, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*PrivateStyle, 0, 1000)
	for cursor.Next(context.Background()) {
		var node = new(PrivateStyle)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetPrivateStyleByExam(book string, exam string) (*PrivateStyle, error) {
	msg := bson.M{"book": book, "exercise": exam}
	result, err := findOneBy(TablePrivateStyle, msg)
	if err != nil {
		return nil, err
	}
	model := new(PrivateStyle)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func UpdatePrivateStyle(uid string, book string) error {
	msg := bson.M{"book": book, "updatedAt": time.Now()}
	_, err := updateOne(TablePrivateStyle, uid, msg)
	return err
}

func UpdatePrivateStyleBase(uid string, page uint16, x uint32, y uint32, ex uint32, ey uint32, bg string) error {
	msg := bson.M{"page": page, "sx": x, "sy": y, "ex": ex, "ey": ey, "background": bg, "updatedAt": time.Now()}
	_, err := updateOne(TablePrivateStyle, uid, msg)
	return err
}

func GetPrivateStyle(uid string) (*PrivateStyle, error) {
	result, err := findOne(TablePrivateStyle, uid)
	if err != nil {
		return nil, err
	}
	model := new(PrivateStyle)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}
