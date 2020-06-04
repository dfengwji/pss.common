package proxy

import (
	"context"
	"errors"
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

	RealName   string `json:"realName" bson:"realName"`
	Grade      uint8  `json:"grade" bson:"grade"`
	Counter    uint8  `json:"counter" bson:"counter"`
	EnrolYear  uint16 `json:"enrolYear" bson:"enrolYear"`
	EnrolMonth uint8  `json:"enrolMonth" bson:"enrolMonth"`

	Scene     string   `json:"scene" bson:"scene"`
	Master    string   `json:"master" bson:"master"`
	Members  []string `json:"members" bson:"members"`
	Students  []string `json:"students" bson:"students"`
	Workbooks []string `json:"workbooks" bson:"workbooks"`
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

func UpdateClassMembers(uid string, members []string) error {
	msg := bson.M{"members": members, "updatedAt": time.Now()}
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
	if len(student) < 1 {
		return errors.New("the student uid is empty")
	}
	msg := bson.M{"students": student}
	_, err := appendElement(TableClasses, uid, msg)
	return err
}

func AppendClassMember(uid string, member string) error {
	if len(member) < 1 {
		return errors.New("the member uid is empty")
	}
	msg := bson.M{"members": member}
	_, err := appendElement(TableClasses, uid, msg)
	return err
}

func UnbindClassMember(uid string, member string) error {
	if len(member) < 1 {
		return errors.New("the member uid is empty")
	}
	msg := bson.M{"members": member}
	_, err := removeElement(TableClasses, uid, msg)
	return err
}

func UnbindClassStudent(uid string, student string) error {
	if len(student) < 1 {
		return errors.New("the student uid is empty")
	}
	msg := bson.M{"students": student}
	_, err := removeElement(TableClasses, uid, msg)
	return err
}
