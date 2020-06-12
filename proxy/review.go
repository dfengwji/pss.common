package proxy

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Review struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`

	Status      uint8              `json:"status" bson:"status"`
	Score       uint16             `json:"score" bson:"score"`
	Postil      string             `json:"postil" bson:"postil"`
	Author      string 				`json:"author" bson:"author"`
	Book        string 				`json:"book" bson:"book"`
	// 笔记本页面快照
	Pages     []PageSnap             `json:"pages" bson:"pages"`
}

type PageSnap struct {
	Index uint16 `json:"index" bson:"index"`
	Snapshot string `json:"snapshot" bson:"snapshot"`
}

func CreateReview(info *Review) error {
	_, err := insertOne(TableReview, info)
	return err
}

func GetReviewNextID() uint64 {
	num, _ := getSequenceNext(TableReview)
	return num
}

func GetAllReviews() ([]*Review, error) {
	cursor, err1 := findAll(TableReview, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Review, 0, 200)
	for cursor.Next(context.Background()) {
		var node = new(Review)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetReview(uid string) (*Review, error) {
	result, err := findOne(TableReview, uid)
	if err != nil {
		return nil, err
	}
	model := new(Review)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetReviewByAuthor(author string) (*Review, error) {
	msg := bson.M{"author": author}
	result, err := findOneBy(TableReview, msg)
	if err != nil {
		return nil, err
	}
	model := new(Review)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetReviewsByAuthor(owner string) ([]*Review, error) {
	cursor, err1 := findMany(TableNoteBook, bson.M{"author": owner}, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Review, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Review)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func RemoveReview(uid string) error {
	_, err := removeOne(TableReview, uid)
	return err
}

func UpdateReviewBase(uid string, st uint8, score uint16, postil string) error {
	msg := bson.M{"status": st, "score": score, "postil": postil, "updatedAt": time.Now()}
	_, err := updateOne(TableWriting, uid, msg)
	return err
}

