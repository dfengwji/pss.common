package proxy

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Student struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`

	Number    string `json:"number" bson:"number"`
	WXCreator string `json:"creator" bson:"creator"`
	Name      string `json:"name" bson:"name"`
	Birthday  string `json:"birthday" bson:"birthday"`
	Sex       uint8  `json:"sex" bson:"sex"`
	Photo     string `json:"photo" bson:"photo"`
	Phone     string `json:"phone" bson:"phone"`
	IDCard    string `json:"card" bson:"card"`

	Books []string `json:"books" bson:"books"`
}

func CreateStudent(info *Student) error {
	_, err := insertOne(TableStudent, info)
	if err != nil {
		return err
	}
	return nil
}

func GetStudentNextID() uint64 {
	num, _ := getSequenceNext(TableRoleID)
	return num
}

func GetStudent(uid string) (*Student, error) {
	result, err := findOne(TableStudent, uid)
	if err != nil {
		return nil, err
	}
	model := new(Student)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetStudentByIDCard(card string) (*Student, error) {
	msg := bson.M{"card": card}
	result, err := findOneBy(TableStudent, msg)
	if err != nil {
		return nil, err
	}
	model := new(Student)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetAllStudents() ([]*Student, error) {
	cursor, err1 := findAll(TableStudent, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Student, 0)
	for cursor.Next(context.Background()) {
		var node = new(Student)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func RemoveStudent(uid string) error {
	_, err := removeOne(TableStudent, uid)
	return err
}

func UpdateStudentBase(uid string, name string, sex uint8, birth string) error {
	msg := bson.M{"name": name, "sex": sex,
		"birthday": birth, "updatedAt": time.Now()}
	_, err := updateOne(TableStudent, uid, msg)
	return err
}

func UpdateStudentPhoto(uid string, photo string) error {
	msg := bson.M{"photo": photo, "updatedAt": time.Now()}
	_, err := updateOne(TableStudent, uid, msg)
	return err
}

func UpdateStudentPhone(uid string, phone string) error {
	msg := bson.M{"phone": phone, "updatedAt": time.Now()}
	_, err := updateOne(TableStudent, uid, msg)
	return err
}

func AppendStudentBook(uid string, book string) error {
	if len(book) < 1 {
		return errors.New("the book uid is empty")
	}
	msg := bson.M{"books": book}
	_, err := appendElement(TableStudent, uid, msg)
	return err
}
