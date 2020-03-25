package proxy

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type EduClass struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`

	Name        string    `json:"name" bson:"name"`
	Book        string    `json:"book" bson:"book"`
	MaxNum      uint16    `json:"maxNum" bson:"maxNum"`
	FromDate    time.Time `json:"fromDate" bson:"fromDate"`
	ToDate      time.Time `json:"toDate" bson:"toDate"`
	DayTime     string    `json:"dayTime" bson:"dayTime"`
	HoldSeconds uint16    `json:"hold" bson:"hold"`
	Remark      string    `json:"remark" bson:"remark"`

	NickName   string `json:"nickName" bson:"nickName"`

	Scene       string   `json:"scene" bson:"scene"`
	Master 		string   `json:"master" bson:"master"`
	Teachers    []string `json:"teachers" bson:"teachers"`
	Students    []string `json:"students" bson:"students"`
	Workbooks   []string `json:"workbooks" bson:"workbooks"`
}

func CreateClass(info *EduClass) error {
	_, err := insertOne(TableClasses, info)
	if err != nil {
		return err
	}
	return nil
}

func GetClassNextID() uint64 {
	num, _ := getSequenceNext(TableClasses)
	return num
}

func GetClassesByScene(scene string) ([]*EduClass, error) {
	msg := bson.M{"scene": scene, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableClasses, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*EduClass, 0, 200)
	for cursor.Next(context.Background()) {
		var node = new(EduClass)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetClass(uid string) (*EduClass, error) {
	result, err := findOne(TableClasses, uid)
	if err != nil {
		return nil, err
	}
	model := new(EduClass)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func UpdateClassBase(uid string, name string, max uint16, remark string) error {
	msg := bson.M{"name": name, "maxNum": max, "remark": remark, "updatedAt": time.Now()}
	_, err := updateOne(TableClasses, uid, msg)
	return err
}

func UpdateClassTime(uid string, from time.Time, to time.Time, day string, hold uint16) error {
	msg := bson.M{"fromDate": from, "toDate": to, "dayTime": day,
		"hold": hold, "updatedAt": time.Now()}
	_, err := updateOne(TableClasses, uid, msg)
	return err
}

func UpdateClassMaster(uid string, teacher string) error {
	msg := bson.M{"master": teacher, "updatedAt": time.Now()}
	_, err := updateOne(TableClasses, uid, msg)
	return err
}

func UpdateClassTeachers(uid string, teachers []string) error {
	msg := bson.M{"teachers": teachers, "updatedAt": time.Now()}
	_, err := updateOne(TableClasses, uid, msg)
	return err
}

func UpdateClassBook(uid string, book string) error {
	msg := bson.M{"book": book, "updatedAt": time.Now()}
	_, err := updateOne(TableClasses, uid, msg)
	return err
}

func RemoveClass(uid string) bool {
	_, err := removeOne(TableClasses, uid)
	if err == nil {
		return true
	}
	return false
}

func AppendClassStudent(uid string, student string) error {
	msg := bson.M{"students": student}
	_, err := appendElement(TableClasses, uid, msg)
	return err
}

func AppendClassTeacher(uid string, teacher string) error {
	msg := bson.M{"teachers": teacher}
	_, err := appendElement(TableClasses, uid, msg)
	return err
}

func UnbindClassTeacher(uid string, teacher string) error {
	msg := bson.M{"teachers": teacher}
	_, err := removeElement(TableClasses, uid, msg)
	return err
}

func UnbindClassStudent(uid string, student string) error {
	msg := bson.M{"students": student}
	_, err := removeElement(TableClasses, uid, msg)
	return err
}
