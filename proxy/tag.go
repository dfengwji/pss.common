package proxy

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Tag struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
}

func GetTag(uid string) (*Tag, error) {
	result, err := findOne(TableExercisesTag, uid)
	if err != nil {
		return nil, err
	}
	model := new(Tag)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetTagNextID() uint64 {
	num, _ := getSequenceNext(TableExercisesTag)
	return num
}

func GetAllTags() ([]*Tag, error) {
	cursor, err1 := findAll(TableExercisesTag, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Tag, 0, 200)
	for cursor.Next(context.Background()) {
		var node = new(Tag)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}
