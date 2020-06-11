package proxy

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type NotebookStyle struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`

	Type        uint8              `json:"type" bson:"type"`
	Page        uint16             `json:"page" bson:"page"`
	Cover       string             `json:"cover" bson:"cover"`
	Background  string             `json:"background" bson:"background"`
	OriginBook  string             `json:"origin" bson:"origin"`
}

func CreateNotebookStyle(info *NotebookStyle) error {
	_, err := insertOne(TableNotebookStyle, info)
	if err != nil {
		return err
	}
	return nil
}

func GetNotebookStyleNextID() uint64 {
	num, _ := getSequenceNext(TableNotebookStyle)
	return num
}

func GetAllNotebookStyles() ([]*NotebookStyle, error) {
	cursor, err1 := findAll(TableNotebookStyle, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*NotebookStyle, 0, 10)
	for cursor.Next(context.Background()) {
		var node = new(NotebookStyle)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetNotebookStylesByOrigin(book string) ([]*NotebookStyle, error) {
	msg := bson.M{"origin": book, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableNotebookStyle, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*NotebookStyle, 0, 30)
	for cursor.Next(context.Background()) {
		var node = new(NotebookStyle)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateNotebookStyleBase(uid, name, origin, cover, bg string) error {
	msg := bson.M{"name": name, "origin": origin, "cover": cover, "background": bg, "updatedAt": time.Now()}
	_, err := updateOne(TableNotebookStyle, uid, msg)
	return err
}

func GetNotebookStyle(uid string) (*NotebookStyle, error) {
	result, err := findOne(TableNotebookStyle, uid)
	if err != nil {
		return nil, err
	}
	model := new(NotebookStyle)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}
