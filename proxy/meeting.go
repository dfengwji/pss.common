package proxy

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Meeting struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`

	Creator  string   `json:"creator" bson:"creator"`
	Operator string   `json:"operator" bson:"operator"`

	/**
	所属组织或者部门
	 */
	Group    string   `json:"group" bson:"group"`
	Remark   string   `json:"remark" bson:"remark"`
	Date     string `json:"date" bson:"date"`
	Members  []string `json:"members" bson:"members"`
}


func CreateMeeting(info *Meeting) error {
	_, err := insertOne(TableMeeting, info)
	if err != nil {
		return err
	}
	return nil
}

func GetMeetingNextID() uint64 {
	num, _ := getSequenceNext(TableMeeting)
	return num
}

func GetAllMeetings() ([]*Meeting, error) {
	cursor, err1 := findAll(TableMeeting, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Meeting, 0, 10)
	for cursor.Next(context.Background()) {
		var node = new(Meeting)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetMeetingsByGroup(group string) ([]*Meeting, error) {
	msg := bson.M{"group": group, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableMeeting, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Meeting, 0, 30)
	for cursor.Next(context.Background()) {
		var node = new(Meeting)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateMeetingBase(uid, name, remark string) error {
	msg := bson.M{"name": name, "remark":remark, "updatedAt": time.Now()}
	_, err := updateOne(TableMeeting, uid, msg)
	return err
}

func GetMeeting(uid string) (*Meeting, error) {
	result, err := findOne(TableMeeting, uid)
	if err != nil {
		return nil, err
	}
	model := new(Meeting)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func RemoveMeeting(uid string) error {
	_, err := removeOne(TableMeeting, uid)
	return err
}