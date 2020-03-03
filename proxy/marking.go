package proxy

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Marking struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Writing     string             `json:"writing" bson:"writing"`
	Status      uint8              `json:"result" bson:"result"`
	OCRText     string             `json:"ocr" bson:"ocr"`
	Score       uint16             `json:"score" bson:"score"`
	Postil      string             `json:"postil" bson:"postil"`
}

func CreateMarking(info *Marking) error {
	_, err := insertOne(TableMarking, info)
	return err
}

func GetMarkingNextID() uint64 {
	num, _ := getSequenceNext(TableMarking)
	return num
}

func GetAllMarkings() ([]*Marking, error) {
	cursor, err1 := findAll(TableMarking, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Marking, 0, 200)
	for cursor.Next(context.Background()) {
		var node = new(Marking)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetMarking(uid string) (*Marking, error) {
	result, err := findOne(TableMarking, uid)
	if err != nil {
		return nil, err
	}
	model := new(Marking)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetMarkingByWriting(uid string) (*Marking, error) {
	msg := bson.M{"writing": uid}
	result, err := findOneBy(TableMarking, msg)
	if err != nil {
		return nil, err
	}
	model := new(Marking)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func UpdateMarkingBase(uid string, st uint8, score uint16, postil string) error {
	msg := bson.M{"result": st, "score": score, "postil": postil, "updatedAt": time.Now()}
	_, err := updateOne(TableWriting, uid, msg)
	return err
}

func UpdateMarkingOCR(uid string, text string) error {
	msg := bson.M{"ocr": text, "updatedAt": time.Now()}
	_, err := updateOne(TableMarking, uid, msg)
	return err
}
