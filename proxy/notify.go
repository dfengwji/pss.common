package proxy

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Notify struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`

	SendTime    time.Time `json:"sendAt" bson:"sendAt"`
	Author      string    `json:"author" bson:"author"`
	Description string    `json:"desc" bson:"desc"`
	Voice       string    `json:"voice" bson:"voice"`
	Targets     []string  `json:"targets" bson:"target"`
	Images      []string  `json:"images" bson:"images"`
}

func CreateNotify(info *Notify) error {
	_, err := insertOne(TableNotify, info)
	if err != nil {
		return err
	}
	return nil
}

func GetNotifyNextID() uint64 {
	num, _ := getSequenceNext(TableNotify)
	return num
}

func GetNotifyCount() int64 {
	num, _ := getCount(TableNotify)
	return num
}

func GetNotify(uid string) (*Notify, error) {
	result, err := findOne(TableNotify, uid)
	if err != nil {
		return nil, err
	}
	model := new(Notify)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetNotifyByID(id uint64) (*Notify, error) {
	msg := bson.M{"id": id}
	result, err := findOneBy(TableNotify, msg)
	if err != nil {
		return nil, err
	}
	model := new(Notify)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func RemoveNotify(uid string) error {
	_, err := removeOne(TableNotify, uid)
	return err
}

func GetAllNotifies() ([]*Notify, error) {
	cursor, err1 := findAll(TableNotify, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Notify, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Notify)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetNotifiesByAuthor(author uint64) ([]*Notify, error) {
	cursor, err1 := findMany(TableNotify, bson.M{"author": author}, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Notify, 0, 100)
	for cursor.Next(context.Background()) {
		var node = new(Notify)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateNotifyBase(uid, title, desc string, send time.Time) error {
	msg := bson.M{"name": title, "desc": desc, "sendAt": send, "updatedAt": time.Now()}
	_, err := updateOne(TableNotify, uid, msg)
	return err
}
