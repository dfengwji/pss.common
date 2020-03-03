package proxy

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type OriginBook struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	Author      string             `json:"author" bson:"author"`
	BN          string             `json:"isbn" bson:"isbn"`
	CreatedTime time.Time          `json:"createAt" bson:"createAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Type        string             `json:"type" bson:"type"`
	Subject     uint8              `json:"subject" bson:"subject"`
}

type OriginStyle struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Book        string             `json:"book" bson:"book"`
	Page        uint32             `json:"pageNo" bson:"pageNo"`
	Exam        string             `json:"examNo" bson:"examNo"`
	StartX      uint32             `json:"sx" bson:"sx"`
	StartY      uint32             `json:"sy" bson:"sy"`
	EndX        uint32             `json:"ex" bson:"ex"`
	EndY        uint32             `json:"ey" bson:"ey"`
	Exercise    string             `json:"exercise" bson:"exercise"`
}

func CreateOriginBook(info *OriginBook) error {
	_, err := insertOne(TableOriginal, info)
	if err != nil {
		return err
	}
	return nil
}

func GetOriginBookNextID() uint64 {
	num, _ := getSequenceNext(TableOriginal)
	return num
}

func GetOriginBookCount() int64 {
	num, _ := getCount(TableOriginal)
	return num
}

func GetAllOriginBooks() ([]*OriginBook, error) {
	cursor, err1 := findAll(TableOriginal, 100)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	var items = make([]*OriginBook, 0, 200)
	for cursor.Next(context.Background()) {
		var node = new(OriginBook)
		if err := cursor.Decode(node); err != nil {
			fmt.Println(err.Error())
			return nil, err
		} else {
			items = append(items, node)
		}
	}

	return items, nil
}

func CreateOriginStyle(info *OriginStyle) error {
	_, err := insertOne(TableOriginStyle, info)
	if err != nil {
		return err
	}
	return nil
}

func GetOriginStyleNextID() uint64 {
	num, _ := getSequenceNext(TableOriginStyle)
	return num
}

func GetAllOriginStyles() ([]*OriginStyle, error) {
	cursor, err1 := findAll(TableOriginStyle, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*OriginStyle, 0, 200)
	for cursor.Next(context.Background()) {
		var node = new(OriginStyle)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}
