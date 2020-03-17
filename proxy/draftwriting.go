package proxy

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type DraftWriting struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Pen         uint64             `json:"pen" bson:"pen"`
	Author      string             `json:"author" bson:"author"`
	Draft       string `json:"draft" bson:"draft"`
	Book        uint64             `json:"book" bson:"book"`
	Page        uint16             `json:"page" bson:"page"`
	Number      uint16             `json:"number" bson:"number"`
	Duration    uint16             `json:"duration" bson:"duration"`
	Dots        string             `json:"dots" bson:"dots"`
}

func CreateDraftWriting(info *DraftWriting) error {
	_, err := insertOne(TableDraftWriting, info)
	return err
}

func GetDraftWritingNextID() uint64 {
	num, _ := getSequenceNext(TableDraftWriting)
	return num
}

func GetDraftWriting(uid string) (*DraftWriting, error) {
	result, err := findOne(TableDraftWriting, uid)
	if err != nil {
		return nil, err
	}
	model := new(DraftWriting)
	err1 := result.Decode(&model)
	if err1 != nil {
		return nil, err1
	}
	return model, err
}

func GetDraftWritingsByDraft(draft string) ([]*DraftWriting, error) {
	filter := bson.M{"draft": draft, "deleteAt": new(time.Time)}
	cursor, err := findManyByOpts(TableDraftWriting, filter, options.Find().SetProjection(bson.M{"dots": 0}))
	if err != nil {
		return nil, err
	}
	var items = make([]*DraftWriting, 0, 200)
	for cursor.Next(context.Background()) {
		var node DraftWriting
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, &node)
		}
	}
	return items, nil
}

func GetDraftWritingDots(uid string) (*DraftWriting, error) {
	result, err := findOneOfField(TableDraftWriting, uid, bson.M{"dots": 1, "dotNum": 1})
	if err != nil {
		return nil, err
	}
	model := new(DraftWriting)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func UpdateDraftWritingDots(uid string, stamp uint64, num uint16, hex string, duration uint16) error {
	msg := bson.M{"dotStamp": stamp, "dotNum": num, "dots": hex, "duration": duration, "updatedAt": time.Now()}
	_, err := updateOne(TableDraftWriting, uid, msg)
	return err
}

func dropDraftWriting() error {
	return dropOne(TableDraftWriting)
}
