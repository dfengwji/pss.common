package proxy

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type MicroCourse struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Name        string             `json:"name" bson:"name"`
	Remark      string             `json:"remark" bson:"remark"`
	Duration    uint32             `json:"duration" bson:"duration"`
	Status      uint8              `json:"status" bson:"status"`
	OpenMode    uint8              `json:"open" bson:"open"`
	Author      string             `json:"author" bson:"author"`
	Owner       string             `json:"owner" bson:"owner"`
	Draft       string             `json:"draft" bson:"draft"`
	Cover       string             `json:"cover" bson:"cover"`
	Menu        string             `json:"menu" bson:"menu"`
	Video       string             `json:"video" bson:"video"`
	Tags        []string           `json:"tags" bson:"tags"`
	OpenTargets	[]string 			`json:"targets" bson:"targets"`
}

func CreateMicroCourse(info *MicroCourse) error {
	_, err := insertOne(TableMicroCourse, info)
	if err != nil {
		return err
	}
	return nil
}

func GetAllMicroCourses() ([]*MicroCourse, error) {
	cursor, err1 := findAll(TableMicroCourse, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*MicroCourse, 0, 10000)
	for cursor.Next(context.Background()) {
		var node = new(MicroCourse)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetMicroCourseNextID() uint64 {
	num, _ := getSequenceNext(TableMicroCourse)
	return num
}

func RemoveMicroCourse(uid string) error {
	_, err := removeOne(TableMicroCourse, uid)
	return err
}

func GetMicroCourse(uid string) (*MicroCourse, error) {
	result, err := findOne(TableMicroCourse, uid)
	if err != nil {
		return nil, err
	}
	model := new(MicroCourse)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetMicroCoursesByAuthor(author string) ([]*MicroCourse, error) {
	msg := bson.M{"author": author, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableMicroCourse, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*MicroCourse, 0, 200)
	for cursor.Next(context.Background()) {
		var node = new(MicroCourse)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateMicroCourseBase(uid string, name string, remark string) error {
	msg := bson.M{"name": name, "remark": remark, "updatedAt": time.Now()}
	_, err := updateOne(TableMicroCourse, uid, msg)
	return err
}

func UpdateMicroCourseOpen(uid string, open uint8, targets []string) error {
	msg := bson.M{"open": open, "targets":targets, "updatedAt": time.Now()}
	_, err := updateOne(TableMicroCourse, uid, msg)
	return err
}

func UpdateMicroCourseCover(uid string, cover string) error {
	msg := bson.M{"cover": cover, "updatedAt": time.Now()}
	_, err := updateOne(TableMicroCourse, uid, msg)
	return err
}

func UpdateMicroCourseStatus(uid string, state uint8) error {
	msg := bson.M{"status": state, "updatedAt": time.Now()}
	_, err := updateOne(TableMicroCourse, uid, msg)
	return err
}

func UpdateMicroCourseTags(uid string, tags []string) error {
	msg := bson.M{"tags": tags, "updatedAt": time.Now()}
	_, err := updateOne(TableMicroCourse, uid, msg)
	return err
}

func UpdateMicroCourseMenu(uid string, menu string) error {
	msg := bson.M{"menu": menu, "updatedAt": time.Now()}
	_, err := updateOne(TableMicroCourse, uid, msg)
	return err
}
