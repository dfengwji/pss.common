package proxy

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type CourseTheme struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Name        string             `json:"name" bson:"name"`
	Remark      string             `json:"remark" bson:"remark"`
	Courses     []string 			`json:"courses" bson:"courses"`
}

func CreateCourseTheme(info *CourseTheme) error {
	_, err := insertOne(TableCourseTheme, info)
	if err != nil {
		return err
	}
	return nil
}

func GetAllCourseThemes() ([]*CourseTheme, error) {
	cursor, err1 := findAll(TableCourseTheme, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*CourseTheme, 0, 100)
	for cursor.Next(context.Background()) {
		var node = new(CourseTheme)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetCourseThemeNextID() uint64 {
	num, _ := getSequenceNext(TableCourseTheme)
	return num
}

func RemoveCourseTheme(uid string) error {
	_, err := removeOne(TableCourseTheme, uid)
	return err
}

func GetCourseTheme(uid string) (*CourseTheme, error) {
	result, err := findOne(TableCourseTheme, uid)
	if err != nil {
		return nil, err
	}
	model := new(CourseTheme)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func UpdateCourseThemeBase(uid string, name string, remark string) error {
	msg := bson.M{"name": name, "remark": remark, "updatedAt": time.Now()}
	_, err := updateOne(TableCourseTheme, uid, msg)
	return err
}

func UpdateThemeCourses(uid string, courses []string) error {
	msg := bson.M{"courses": courses, "updatedAt": time.Now()}
	_, err := updateOne(TableCourseTheme, uid, msg)
	return err
}

func AppendCourseInTheme(uid string, course string) error {
	if len(course) < 1 {
		return errors.New("the course uid is empty")
	}
	msg := bson.M{"courses": course}
	_, err := appendElement(TableCourseTheme, uid, msg)
	return err
}

func UnbindCourseInTheme(uid string, course string) error {
	if len(course) < 1 {
		return errors.New("the course uid is empty")
	}
	msg := bson.M{"courses": course}
	_, err := removeElement(TableCourseTheme, uid, msg)
	return err
}