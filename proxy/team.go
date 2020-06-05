package proxy


import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Team struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`

	MaxNum      uint16    `json:"maxNum" bson:"maxNum"`
	Name        string    `json:"name" bson:"name"`
	Creator  string `json:"creator" bson:"creator"`
	Remark	 string `json:"remark" bson:"remark"`
	Master    string   `json:"master" bson:"master"`
	Members  []string `json:"members" bson:"members"`
	Department  string `json:"department" bson:"department"`
}

func CreateTeam(info *Team) error {
	_, err := insertOne(TableTeam, info)
	if err != nil {
		return err
	}
	return nil
}

func GetTeamNextID() uint64 {
	num, _ := getSequenceNext(TableTeam)
	return num
}

func GetTeamesByDepartment(department string) ([]*Team, error) {
	msg := bson.M{"department": department, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableTeam, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Team, 0, 200)
	for cursor.Next(context.Background()) {
		var node = new(Team)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetTeam(uid string) (*Team, error) {
	result, err := findOne(TableTeam, uid)
	if err != nil {
		return nil, err
	}
	model := new(Team)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func UpdateTeamBase(uid string, name string, max uint16, remark string) error {
	msg := bson.M{"name": name, "maxNum": max, "remark": remark, "updatedAt": time.Now()}
	_, err := updateOne(TableTeam, uid, msg)
	return err
}

func UpdateTeamMaster(uid string, teacher string) error {
	msg := bson.M{"master": teacher, "updatedAt": time.Now()}
	_, err := updateOne(TableTeam, uid, msg)
	return err
}

func UpdateTeamMembers(uid string, members []string) error {
	msg := bson.M{"members": members, "updatedAt": time.Now()}
	_, err := updateOne(TableTeam, uid, msg)
	return err
}

func UpdateTeamBook(uid string, book string) error {
	msg := bson.M{"book": book, "updatedAt": time.Now()}
	_, err := updateOne(TableTeam, uid, msg)
	return err
}

func RemoveTeam(uid string) bool {
	_, err := removeOne(TableTeam, uid)
	if err == nil {
		return true
	}
	return false
}

func AppendTeamMember(uid string, member string) error {
	if len(member) < 1 {
		return errors.New("the member uid is empty")
	}
	msg := bson.M{"members": member}
	_, err := appendElement(TableTeam, uid, msg)
	return err
}

func UnbindTeamMember(uid string, member string) error {
	if len(member) < 1 {
		return errors.New("the member uid is empty")
	}
	msg := bson.M{"members": member}
	_, err := removeElement(TableTeam, uid, msg)
	return err
}

