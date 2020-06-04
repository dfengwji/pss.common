package proxy

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Member struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`

	Name       string `json:"name" bson:"name"`
	Number     string `json:"number" bson:"number"`
	Sex        uint8  `json:"sex" bson:"sex"`
	Type       uint8  `json:"type" bson:"type"`
	Phone      string `json:"phone" bson:"phone"`
	Photo      string `json:"photo" bson:"photo"`
	Passwords  string `json:"psw" bson:"psw"`
	Specialty  string `json:"specialty" bson:"specialty"`
	Desc       string `json:"desc" bson:"desc"`
	Experience string `json:"experience" bson:"experience"`
}

func CreateMember(info *Member) error {
	_, err := insertOne(TableMember, info)
	if err != nil {
		return err
	}
	return nil
}

func GetMemberNextID() uint64 {
	num, _ := getSequenceNext(TableRoleID)
	return num
}

func GetMember(uid string) (*Member, error) {
	result, err := findOne(TableMember, uid)
	if err != nil {
		return nil, err
	}
	model := new(Member)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetAllMembers() ([]*Member, error) {
	cursor, err1 := findAll(TableMember, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Member, 0)
	for cursor.Next(context.Background()) {
		var node = new(Member)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetMemberByPhone(phone string) (*Member, error) {
	msg := bson.M{"phone": phone}
	result, err := findOneBy(TableMember, msg)
	if err != nil {
		return nil, err
	}
	model := new(Member)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func RemoveMember(uid string) error {
	_, err := removeOne(TableMember, uid)
	return err
}

func UpdateMemberBase(uid string, name string, sex uint8, kind uint8, desc string) error {
	msg := bson.M{"name": name, "sex": sex, "type":kind,
		"desc": desc, "updatedAt": time.Now()}
	_, err := updateOne(TableMember, uid, msg)
	return err
}

func UpdateMemberMore(uid string, spec string, experience string) error {
	msg := bson.M{"specialty": spec, "experience": experience,
		"updatedAt": time.Now()}
	_, err := updateOne(TableMember, uid, msg)
	return err
}

func UpdateMemberPsw(uid string, psw string) error {
	msg := bson.M{"psw": psw, "updatedAt": time.Now()}
	_, err := updateOne(TableMember, uid, msg)
	return err
}

func UpdateMemberPhoto(uid string, cover string) error {
	msg := bson.M{"photo": cover, "updatedAt": time.Now()}
	_, err := updateOne(TableMember, uid, msg)
	return err
}

func UpdateMemberPhone(uid string, cover string) error {
	msg := bson.M{"photo": cover, "updatedAt": time.Now()}
	_, err := updateOne(TableMember, uid, msg)
	return err
}

