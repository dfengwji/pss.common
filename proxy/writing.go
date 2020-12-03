package proxy

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Writing struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Pen         uint64             `json:"pen" bson:"pen"`
	Writer        uint64             `json:"writer" bson:"writer"`
	Book        string             `json:"book" bson:"book"`
	Exercise    string             `json:"exercise" bson:"exercise"`
	Style       string             `json:"style" bson:"style"`
	Color       string             `json:"color" bson:"color"`
	Duration    uint16             `json:"duration" bson:"duration"`
	DotPage     uint16             `json:"dotPage" bson:"dotPage"`
	DotBook     uint64             `json:"dotBook" bson:"dotBook"`
	DotStamp    uint64             `json:"dotStamp" bson:"dotStamp"`
	DotNum      uint32             `json:"dotNum" bson:"dotNum"`
	Dots        string             `json:"dots" bson:"dots"`
}

func CreateWriting(info *Writing) error {
	_, err := insertOne(TableWriting, info)
	return err
}

func GetWritingNextID() uint64 {
	num, _ := getSequenceNext(TableWriting)
	return num
}

func GetWriting(uid string) (*Writing, error) {
	result, err := findOne(TableWriting, uid)
	if err != nil {
		return nil, err
	}
	model := new(Writing)
	err1 := result.Decode(&model)
	if err1 != nil {
		return nil, err1
	}
	return model, err
}

/**
查询不包含dots字段的数据
 */
func GetWritingsByWriter(writer uint64, book string) ([]*Writing, error) {
	filter := bson.M{"writer": writer, "book": book, "deleteAt": new(time.Time)}
	cursor, err := findManyByOpts(TableWriting, filter, options.Find().SetProjection(bson.M{"dots": 0}))
	if err != nil {
		return nil, err
	}
	var items = make([]*Writing, 0, 200)
	for cursor.Next(context.Background()) {
		var node Writing
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, &node)
		}
	}
	return items, nil
}

func GetWritingByExercise(writer uint64, book string, exercise string) (*Writing, error) {
	msg := bson.M{"writer": writer, "book": book, "exercise": exercise}
	result, err := findOneBy(TableWriting, msg)
	if err != nil {
		return nil, err
	}
	model := new(Writing)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetWritingByStyle(writer uint64, style string) (*Writing, error) {
	msg := bson.M{"writer": writer, "style": style}
	result, err := findOneBy(TableWriting, msg)
	if err != nil {
		return nil, err
	}
	model := new(Writing)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetWritingDots(uid string) (*Writing, error) {
	result, err := findOneOfField(TableWriting, uid, bson.M{"dots": 1, "dotNum": 1})
	if err != nil {
		return nil, err
	}
	model := new(Writing)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func UpdateWritingDots(uid string, stamp uint64, num uint32, hex string, duration uint16) error {
	msg := bson.M{"dotStamp": stamp, "dotNum": num, "dots": hex, "duration": duration, "updatedAt": time.Now()}
	_, err := updateOne(TableWriting, uid, msg)
	return err
}

func GetAllWritings() ([]*Writing, error) {
	cursor, err1 := findAllByOpts(TableWriting, options.Find().SetProjection(bson.M{"dots": 0}))
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Writing, 0, 200)
	for cursor.Next(context.Background()) {
		var node = new(Writing)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func dropWriting() error {
	return dropOne(TableWriting)
}
