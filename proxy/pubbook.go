package proxy

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type PublicBook struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	Author      string             `json:"author" bson:"author"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Subject     uint16             `json:"type" bson:"type"`
	Cover       string             `json:"cover" bson:"cover"`
	Score       uint16             `json:"score" bson:"score"`
	PDF         string				`json:"pdf" bson:"pdf"`
	Styles      []string           `json:"styles" bson:"styles"`
}

func CreatePublicBook(info *PublicBook) error {
	_, err := insertOne(TablePublicBook, info)
	if err != nil {
		return err
	}
	return nil
}

func GetPublicBookNextID() uint64 {
	num, _ := getSequenceNext(TablePublicBook)
	return num
}

func GetPublicBookCount() int64 {
	num, _ := getCount(TablePublicBook)
	return num
}

func GetPublicBook(uid string) (*PublicBook, error) {
	result, err := findOne(TablePublicBook, uid)
	if err != nil {
		return nil, err
	}
	model := new(PublicBook)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func RemovePublicBook(uid string) error {
	_, err := removeOne(TablePublicBook, uid)
	return err
}

func GetAllPublicBooks() ([]*PublicBook, error) {
	cursor, err1 := findAll(TablePublicBook, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*PublicBook, 0, 200)
	for cursor.Next(context.Background()) {
		var node = new(PublicBook)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetPublicBooksByAuthor(author string) ([]*PublicBook, error) {
	msg := bson.M{"author": author, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TablePublicBook, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*PublicBook, 0, 200)
	for cursor.Next(context.Background()) {
		var node = new(PublicBook)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdatePublicBookCover(uid string, cover string) error {
	msg := bson.M{"cover": cover, "updatedAt": time.Now()}
	_, err := updateOne(TablePublicBook, uid, msg)
	return err
}

func UpdatePublicBookPDF(uid string, pdf string) error {
	msg := bson.M{"pdf": pdf, "updatedAt": time.Now()}
	_, err := updateOne(TablePublicBook, uid, msg)
	return err
}

func UpdatePublicBookBase(uid string, name string, subject uint16, score uint16) error {
	msg := bson.M{"name": name, "subject": subject, "score": score, "updatedAt": time.Now()}
	_, err := updateOne(TablePublicBook, uid, msg)
	return err
}

func UpdatePublicBookStyles(uid string, styles []string) error {
	msg := bson.M{"styles": styles, "updatedAt": time.Now()}
	_, err := updateOne(TablePublicBook, uid, msg)
	return err
}

func AppendPublicBookStyle(uid string, style string) error {
	msg := bson.M{"styles": style}
	_, err := appendElement(TablePublicBook, uid, msg)
	return err
}

func UnbindPublicBookStyle(uid string, style string) error {
	msg := bson.M{"styles": style}
	_, err := removeElement(TablePublicBook, uid, msg)
	return err
}
