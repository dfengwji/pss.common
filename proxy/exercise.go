package proxy

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Exercise struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Title       string             `json:"title" bson:"title"`
	Desc        string             `json:"desc" bson:"desc"`
	Grade       uint8              `json:"grade" bson:"grade"`
	Type        uint8              `json:"type" bson:"type"`
	Answer      string             `json:"answer" bson:"answer"`
	Owner		string 			   `json:"owner" bson:"owner"`
	Author      string             `json:"author" bson:"author"`
	Tags        []string           `json:"tags" bson:"tags"`
}

func CreateExercise(info *Exercise) error {
	_, err := insertOne(TableExercises, info)
	if err != nil {
		return err
	}
	return nil
}

func GetAllExercises() ([]*Exercise, error) {
	cursor, err1 := findAll(TableExercises, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Exercise, 0, 10000)
	for cursor.Next(context.Background()) {
		var node = new(Exercise)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetExerciseNextID() uint64 {
	num, _ := getSequenceNext(TableExercises)
	return num
}

func RemoveExercise(uid string) error {
	_, err := removeOne(TableExercises, uid)
	return err
}

func GetExercise(uid string) (*Exercise, error) {
	result, err := findOne(TableExercises, uid)
	if err != nil {
		return nil, err
	}
	model := new(Exercise)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetExercisesByAuthor(author string) ([]*Exercise, error) {
	msg := bson.M{"author": author, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableExercises, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Exercise, 0, 200)
	for cursor.Next(context.Background()) {
		var node = new(Exercise)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateExerciseBase(uid string, title string, desc string, answer string, grade uint8) error {
	msg := bson.M{"title": title, "desc": desc, "answer": answer,
		"grade": grade, "updatedAt": time.Now()}
	_, err := updateOne(TableExercises, uid, msg)
	return err
}

func UpdateExerciseTags(uid string, tags []string) error {
	msg := bson.M{"tags": tags, "updatedAt": time.Now()}
	_, err := updateOne(TableExercises, uid, msg)
	return err
}
