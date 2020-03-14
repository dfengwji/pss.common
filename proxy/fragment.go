package proxy

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

const (
	FragmentStatusIdle = 0
	FragmentStatusUsed = 1
)

type Fragment struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64    `json:"id" bson:"id"`
	Name        string    `json:"name" bson:"name"`
	CreatedTime time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time `json:"deleteAt" bson:"deleteAt"`
	Status      uint8     `json:"status" bson:"status"`
	Pen         uint64    `json:"pen" bson:"pen"`
	User        uint64    `json:"user" bson:"user"`
	Color       string    `json:"color" bson:"color"`
	ExamStyle   string    `json:"style" bson:"style"`
	Stamp       uint64    `json:"stamp" bson:"stamp"`
	DotBook     uint64    `json:"dotBook" bson:"dotBook"`
	DotPage     uint16    `json:"dotPage" bson:"dotPage"`
	DotNum      uint16    `json:"dotNum" bson:"dotNum"`
	Dots        string    `json:"dots" bson:"dots"`
}

func CreateFragment(info *Fragment) error {
	_, err := insertOne(TableFragment, info)
	return err
}

func GetFragmentNextID() uint64 {
	num,_ := getSequenceNext(TableFragment)
	return num
}

func GetFragment(uid string) (*Fragment,error) {
	result, err := findOne(TableFragment, uid)
	if err != nil {
		return nil,err
	}
	model := new(Fragment)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil,err1
	}
	return model,nil
}

func GetIdleFragments() ([]*Fragment,error) {
	msg := bson.M{"status": 0, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableFragment, msg, 0)
	if err1 != nil {
		return nil,err1
	}
	var items = make([]*Fragment, 0, 200)
	for cursor.Next(context.Background()) {
		var node = new(Fragment)
		if err := cursor.Decode(node); err != nil {
			return nil,err
		} else {
			items = append(items, node)
		}
	}
	return items,nil
}

func UpdateFragmentStatus(uid string, status uint8) error {
	msg := bson.M{"status": status, "updatedAt": time.Now()}
	_, err := updateOne(TableFragment, uid, msg)
	return err
}

func GetAllFragments() ([]*Fragment,error) {
	cursor, err1 := findAll(TableFragment, 0)
	if err1 != nil {
		return nil,err1
	}
	var items = make([]*Fragment, 0, 200)
	for cursor.Next(context.Background()) {
		var node = new(Fragment)
		if err := cursor.Decode(node); err != nil {
			return nil,err
		} else {
			items = append(items, node)
		}
	}
	return items,nil
}

func RemoveFragment(uid string) bool {
	_, err := removeOne(TableFragment, uid)
	if err == nil {
		return true
	}
	return false
}

func dropFragment() error {
	err := dropOne(TableFragment)
	return err
}
