package proxy

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type CourseMenu struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Name        string             `json:"name" bson:"name"`
	Beta        uint8 				`json:"beta" bson:"beta"`
	Parent      string             `json:"parent" bson:"parent"`
	Cover       string             `json:"cover" bson:"cover"`
	Remark      string             `json:"remark" bson:"remark"`
}

func CreateCourseMenu(info *CourseMenu) error {
	_, err := insertOne(TableCourseMenu, info)
	if err != nil {
		return err
	}
	return nil
}

func GetAllCourseMenus() ([]*CourseMenu, error) {
	cursor, err1 := findAll(TableCourseMenu, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*CourseMenu, 0, 100)
	for cursor.Next(context.Background()) {
		var node = new(CourseMenu)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetCourseMenuNextID() uint64 {
	num, _ := getSequenceNext(TableCourseMenu)
	return num
}

func RemoveCourseMenu(uid string) error {
	_, err := removeOne(TableCourseMenu, uid)
	return err
}

func GetCourseMenu(uid string) (*CourseMenu, error) {
	result, err := findOne(TableCourseMenu, uid)
	if err != nil {
		return nil, err
	}
	model := new(CourseMenu)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetAllTopCourseMenus() ([]*CourseMenu, error) {
	msg := bson.M{"parent": "", "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableCourseMenu, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*CourseMenu, 0, 200)
	for cursor.Next(context.Background()) {
		var node = new(CourseMenu)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetCourseMenusByTop(parent string) ([]*CourseMenu, error) {
	msg := bson.M{"parent": parent, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableCourseMenu, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*CourseMenu, 0, 200)
	for cursor.Next(context.Background()) {
		var node = new(CourseMenu)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateCourseMenuBase(uid string, name string, remark string) error {
	msg := bson.M{"name": name, "remark": remark, "updatedAt": time.Now()}
	_, err := updateOne(TableCourseMenu, uid, msg)
	return err
}

func UpdateCourseMenuBeta(uid string, beta uint8) error {
	msg := bson.M{"beta": beta, "updatedAt": time.Now()}
	_, err := updateOne(TableCourseMenu, uid, msg)
	return err
}
