package proxy

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type PublicStyle struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Book        string             `json:"book" bson:"book"`
	Page        uint16             `json:"page" bson:"page"`
	Background  string             `json:"background" bson:"background"`
	Example     string             `json:"examNo" bson:"examNo"`
	StartX      uint32             `json:"sx" bson:"sx"`
	StartY      uint32             `json:"sy" bson:"sy"`
	EndX        uint32             `json:"ex" bson:"ex"`
	EndY        uint32             `json:"ey" bson:"ey"`
	Exercise    string             `json:"exercise" bson:"exercise"`
	Score       uint16             `json:"score" bson:"score"`
}

func CreatePublicStyle(info *PublicStyle) error {
	_, err := insertOne(TablePublicStyle, info)
	if err != nil {
		return err
	}
	return nil
}

func GetPublicStyleNextID() uint64 {
	num, _ := getSequenceNext(TablePublicStyle)
	return num
}

func GetAllPublicStyles() ([]*PublicStyle, error) {
	cursor, err1 := findAll(TablePublicStyle, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*PublicStyle, 0, 10000)
	for cursor.Next(context.Background()) {
		var node = new(PublicStyle)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetPublicStylesByBook(book string) ([]*PublicStyle, error) {
	msg := bson.M{"book": book, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TablePublicStyle, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*PublicStyle, 0, 10000)
	for cursor.Next(context.Background()) {
		var node = new(PublicStyle)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetPublicStyleByExam(book string, exam string) (*PublicStyle, error) {
	msg := bson.M{"book": book, "exercise": exam}
	result, err := findOneBy(TablePublicStyle, msg)
	if err != nil {
		return nil, err
	}
	model := new(PublicStyle)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func UpdatePublicStyleBook(uid string, book string) error {
	msg := bson.M{"book": book, "updatedAt": time.Now()}
	_, err := updateOne(TablePublicStyle, uid, msg)
	return err
}

func UpdatePublicStyleBase(uid string, page uint16, x uint32, y uint32, ex uint32, ey uint32, bg string) error {
	msg := bson.M{"page": page, "sx": x, "sy": y, "ex": ex, "ey": ey, "background": bg, "updatedAt": time.Now()}
	_, err := updateOne(TablePublicStyle, uid, msg)
	return err
}

func UpdatePublicStyleScore(uid string, score uint16) error {
	msg := bson.M{"score": score, "updatedAt": time.Now()}
	_, err := updateOne(TablePublicStyle, uid, msg)
	return err
}

func GetPublicStyle(uid string) (*PublicStyle, error) {
	result, err := findOne(TablePublicStyle, uid)
	if err != nil {
		return nil, err
	}
	model := new(PublicStyle)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}
