package proxy

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Teacher struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`

	Name       string `json:"name" bson:"name"`
	Number     string `json:"number" bson:"number"`
	Sex        uint8  `json:"sex" bson:"sex"`
	Phone      string `json:"phone" bson:"phone"`
	Photo      string `json:"photo" bson:"photo"`
	Passwords  string `json:"psw" bson:"psw"`
	Specialty  string `json:"specialty" bson:"specialty"`
	Desc       string `json:"desc" bson:"desc"`
	Experience string `json:"experience" bson:"experience"`
}

func CreateTeacher(info *Teacher) error {
	_, err := insertOne(TableTeacher, info)
	if err != nil {
		return err
	}
	return nil
}

func GetTeacherNextID() uint64 {
	num, _ := getSequenceNext(TableTeacher)
	return num
}

func GetTeacher(uid string) (*Teacher, error) {
	result, err := findOne(TableTeacher, uid)
	if err != nil {
		return nil, err
	}
	model := new(Teacher)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetAllTeachers() ([]*Teacher, error) {
	cursor, err1 := findAll(TableTeacher, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Teacher, 0)
	for cursor.Next(context.Background()) {
		var node = new(Teacher)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetTeacherByPhone(phone string) (*Teacher, error) {
	msg := bson.M{"phone": phone}
	result, err := findOneBy(TableTeacher, msg)
	if err != nil {
		return nil, err
	}
	model := new(Teacher)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func RemoveTeacher(uid string) bool {
	_, err := removeOne(TableTeacher, uid)
	if err == nil {
		return true
	}
	return false
}

func UpdateTeacherBase(uid string, name string, sex uint8, desc string) error {
	msg := bson.M{"name": name, "sex": sex,
		"desc": desc, "updatedAt": time.Now()}
	_, err := updateOne(TableTeacher, uid, msg)
	return err
}

func UpdateTeacherMore(uid string, spec string, experience string) error {
	msg := bson.M{"specialty": spec, "experience": experience,
		"updatedAt": time.Now()}
	_, err := updateOne(TableTeacher, uid, msg)
	return err
}

func UpdateTeacherPsw(uid string, psw string) error {
	msg := bson.M{"psw": psw, "updatedAt": time.Now()}
	_, err := updateOne(TableTeacher, uid, msg)
	return err
}

func UpdateTeacherPhoto(uid string, cover string) error {
	msg := bson.M{"photo": cover, "updatedAt": time.Now()}
	_, err := updateOne(TableTeacher, uid, msg)
	return err
}

func UpdateTeacherPhone(uid string, cover string) error {
	msg := bson.M{"photo": cover, "updatedAt": time.Now()}
	_, err := updateOne(TableTeacher, uid, msg)
	return err
}
