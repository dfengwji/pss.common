package proxy

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type PrivateBook struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	Author      string             `json:"author" bson:"author"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Type        uint8              `json:"type" bson:"type"`
	Cover       string             `json:"cover" bson:"cover"`
	Remark      string             `json:"remark" bson:"remark"`
	Count       uint16             `json:"count" bson:"count"`
	Exams       []string           `json:"exams" bson:"exams"`
}

func CreateMisBook(info *PrivateBook) error {
	_, err := insertOne(TablePrivateBook, info)
	if err != nil {
		return err
	}
	return nil
}

func GetPrivateBookNextID() uint64 {
	num, _ := getSequenceNext(TablePrivateBook)
	return num
}

func GetPrivateBookCount() int64 {
	num, _ := getCount(TablePrivateBook)
	return num
}

func GetPrivateBook(uid string) (*PrivateBook, error) {
	result, err := findOne(TablePrivateBook, uid)
	if err != nil {
		return nil, err
	}
	model := new(PrivateBook)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func RemovePrivateBook(uid string) error {
	_, err := removeOne(TablePrivateBook, uid)
	return err
}

func GetPrivateBooks() ([]*PrivateBook, error) {
	cursor, err1 := findAll(TablePrivateBook, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*PrivateBook, 0, 200)
	for cursor.Next(context.Background()) {
		var node = new(PrivateBook)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetPrivateBooksByAuthor(author string) ([]*PrivateBook, error) {
	msg := bson.M{"author": author, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TablePrivateBook, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*PrivateBook, 0, 200)
	for cursor.Next(context.Background()) {
		var node = new(PrivateBook)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdatePrivateBookCover(uid string, cover string) error {
	msg := bson.M{"cover": cover, "updatedAt": time.Now()}
	_, err := updateOne(TablePrivateBook, uid, msg)
	return err
}

func UpdatePrivateBookBase(uid string, name string, remark string) error {
	msg := bson.M{"name": name, "remark": remark, "updatedAt": time.Now()}
	_, err := updateOne(TablePrivateBook, uid, msg)
	return err
}

func UpdatePrivateBookExams(uid string, exams []string) error {
	msg := bson.M{"exams": exams, "updatedAt": time.Now()}
	_, err := updateOne(TablePrivateBook, uid, msg)
	return err
}

func AppendPrivateBookExam(uid string, exam string) error {
	msg := bson.M{"exams": exam}
	_, err := appendElement(TablePrivateBook, uid, msg)
	return err
}

func UnbindPrivateBookExam(uid string, exam string) error {
	msg := bson.M{"exams": exam}
	_, err := removeElement(TablePrivateBook, uid, msg)
	return err
}
