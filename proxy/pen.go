package proxy

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Pen struct {
	UID         primitive.ObjectID `bson:"_id"`
	Status      uint8              `json:"status" bson:"status"`
	ID          uint64             `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	UploadTime  time.Time          `json:"uploadAt" bson:"uploadAt"`
	UseTime     time.Time          `json:"useAt" bson:"useAt"`
	Mac         string             `json:"mac" bson:"mac"`
	User        string             `json:"user" bson:"user"`
	Child       string             `json:"child" bson:"child"`
}

func CreatePen(info *Pen) error {
	_, err := insertOne(TablePen, info)
	return err
}

func GetPenNextID() uint64 {
	num, _ := getSequenceNext(TablePen)
	return num
}

func GetPen(uid string) (*Pen, error) {
	result, err := findOne(TablePen, uid)
	if err != nil {
		return nil, err
	}
	model := new(Pen)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetPenByMac(mac string) (*Pen, error) {
	msg := bson.M{"mac": mac}
	result, err := findOneBy(TablePen, msg)
	if err != nil {
		return nil, err
	}
	model := new(Pen)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func HadPen(mac string) error {
	msg := bson.M{"mac": mac}
	_, err := hadOne(TablePen, msg)
	return err
}

func GetAllPens() ([]*Pen, error) {
	cursor, err1 := findAll(TablePen, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Pen, 0, 200)
	for cursor.Next(context.Background()) {
		var node = new(Pen)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetPensByUser(user string) ([]*Pen, error) {
	msg := bson.M{"user": user, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TablePen, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Pen, 0, 10)
	for cursor.Next(context.Background()) {
		var node = new(Pen)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func RemovePen(uid string) bool {
	_, err := removeOne(TablePen, uid)
	if err == nil {
		return true
	}
	return false
}

func UpdatePenUser(uid string, user string, child string) error {
	msg := bson.M{"user": user, "child": child, "updatedAt": time.Now()}
	_, err := updateOne(TablePen, uid, msg)
	return err
}

func UpdatePenName(uid string, name string) error {
	msg := bson.M{"name": name, "updatedAt": time.Now()}
	_, err := updateOne(TablePen, uid, msg)
	return err
}

func dropPen() error {
	return dropOne(TablePen)
}
