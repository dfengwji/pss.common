package proxy

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Apply struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`

	Applicant  string    `json:"applicant" bson:"applicant"`
	Student    string    `json:"student" bson:"student"`
	Status     uint8     `json:"status" bson:"status"`
	Scene      string    `json:"scene" bson:"scene"`
	Class      string    `json:"class" bson:"class"`
	SubmitTime time.Time `json:"submit" bson:"submit"`
}

func CreateApply(info *Apply) error {
	_, err := insertOne(TableApply, info)
	if err != nil {
		return err
	}
	return nil
}

func GetApplyNextID() uint64 {
	num, _ := getSequenceNext(TableApply)
	return num
}

func GetApply(uid string) (*Apply, error) {
	result, err := findOne(TableApply, uid)
	if err != nil {
		return nil, err
	}
	model := new(Apply)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetAppliesByClass(class string) ([]*Apply, error) {
	msg := bson.M{"class": class, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableStudent, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Apply, 0, 5)
	for cursor.Next(context.Background()) {
		var node = new(Apply)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetAppliesByApplicant(user string) ([]*Apply, error) {
	msg := bson.M{"applicant": user, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableStudent, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Apply, 0, 5)
	for cursor.Next(context.Background()) {
		var node = new(Apply)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateApply(uid string, status uint8) error {
	msg := bson.M{"status": status, "updatedAt": time.Now()}
	_, err := updateOne(TableApply, uid, msg)
	return err
}

func RemoveApply(uid string) error {
	_, err := removeOne(TableApply, uid)
	return err
}
