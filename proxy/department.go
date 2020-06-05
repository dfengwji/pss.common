package proxy

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Department struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`

	Scene string 		`json:"scene" bson:"scene"`
	Remark string 		`json:"remark" bson:"remark"`
}

func CreateDepartment(info *Department) error {
	_, err := insertOne(TableDepartment, info)
	if err != nil {
		return err
	}
	return nil
}

func GetDepartmentNextID() uint64 {
	num, _ := getSequenceNext(TableDepartment)
	return num
}

func GetDepartmentCount() int64 {
	num, _ := getCount(TableDepartment)
	return num
}

func GetDepartment(uid string) (*Department, error) {
	result, err := findOne(TableDepartment, uid)
	if err != nil {
		return nil, err
	}
	model := new(Department)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetDepartmentByID(id uint64) (*Department, error) {
	msg := bson.M{"id": id}
	result, err := findOneBy(TableDepartment, msg)
	if err != nil {
		return nil, err
	}
	model := new(Department)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func RemoveDepartment(uid string) error {
	_, err := removeOne(TableDepartment, uid)
	return err
}

func GetAllDepartments() ([]*Department, error) {
	cursor, err1 := findAll(TableDepartment, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Department, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Department)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetDepartmentsByScene(scene string) ([]*Department, error) {
	cursor, err1 := findMany(TableDepartment, bson.M{"scene": scene}, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Department, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Department)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateDepartmentBase(uid, name, remark string) error {
	msg := bson.M{"name": name, "remark": remark, "updatedAt": time.Now()}
	_, err := updateOne(TableDepartment, uid, msg)
	return err
}