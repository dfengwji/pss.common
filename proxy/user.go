package proxy

import (
	"context"
	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	UID          primitive.ObjectID `bson:"_id"`
	ID           uint64             `json:"id" bson:"id"`
	Name         string             `json:"name" bson:"name"`
	CreatedTime  time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime  time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime   time.Time          `json:"deleteAt" bson:"deleteAt"`
	Status       uint8              `json:"status" bson:"status"`
	Sex          uint8              `json:"sex" bson:"sex"`
	Phone        string             `json:"phone" bson:"phone"`
	Email        string             `json:"email" bson:"email"`
	Password     string             `json:"psw" bson:"psw"`
	Portrait     string             `json:"port" bson:"port"`
	Wechat       string             `json:"wechat" bson:"wechat"`
	Address      string             `json:"address" bson:"address"`
	AppointChild string             `json:"appoint" bson:"appoint"`
	Children     []string           `json:"children" bson:"children"`
	Role         string             `json:"role" bson:"role"`
}

func CreateUser(info *User) error {
	_, err := insertOne(TableUser, info)
	if err != nil {
		return err
	}
	return nil
}

func GetUserNextID() uint64 {
	num, _ := getSequenceNext(TableUser)
	return num
}

func GetUser(uid string) (*User, error) {
	result, err := findOne(TableUser, uid)
	if err != nil {
		return nil, err
	}
	model := new(User)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetUserByWX(wechat string) (*User, error) {
	msg := bson.M{"wechat": wechat}
	result, err := findOneBy(TableUser, msg)
	if err != nil {
		return nil, err
	}
	model := new(User)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetUserByPhone(phone string) (*User, error) {
	msg := bson.M{"phone": phone}
	result, err := findOneBy(TableUser, msg)
	if err != nil {
		return nil, err
	}
	model := new(User)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetAllUsers() ([]*User, error) {
	cursor, err1 := findAll(TableUser, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*User, 0, 200)
	for cursor.Next(context.Background()) {
		var node = new(User)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func RemoveUser(uid string) bool {
	_, err := removeOne(TableUser, uid)
	if err == nil {
		return true
	}
	return false
}

func UpdateUserBase(uid string, name string, sex uint8, port string) error {
	msg := bson.M{"name": name, "sex": sex, "port": port, "updatedAt": time.Now()}
	_, err := updateOne(TableUser, uid, msg)
	return err
}

func UpdateUserAddress(uid string, address string) error {
	msg := bson.M{"address": address, "updatedAt": time.Now()}
	_, err := updateOne(TableUser, uid, msg)
	return err
}

func UpdateUserPortrait(uid string, port string) error {
	msg := bson.M{"port": port, "updatedAt": time.Now()}
	_, err := updateOne(TableUser, uid, msg)
	return err
}

func UpdateUserAppoint(uid string, child string) error {
	msg := bson.M{"appoint": child, "updatedAt": time.Now()}
	_, err := updateOne(TableUser, uid, msg)
	return err
}

func UpdateUserPassword(uid string, psw string) error {
	msg := bson.M{"psw": psw, "updatedAt": time.Now()}
	_, err := updateOne(TableUser, uid, msg)
	return err
}

func UpdateUserPhone(uid string, phone string) error {
	msg := bson.M{"phone": phone, "updatedAt": time.Now()}
	_, err := updateOne(TableUser, uid, msg)
	return err
}

func UpdateUserRole(uid string, role string, phone string) error {
	msg := bson.M{"role": role, "phone": phone, "updatedAt": time.Now()}
	_, err := updateOne(TableUser, uid, msg)
	return err
}

func AppendUserChild(uid string, child string) error {
	msg := bson.M{"children": child}
	_, err := appendElement(TableUser, uid, msg)
	return err
}

func UnbindUserChild(uid string, child string) error {
	msg := bson.M{"children": child}
	_, err := removeElement(TableUser, uid, msg)
	return err
}

func dropUser() error {
	err := dropOne(TableUser)
	return err
}

func recoverUsers(data []gjson.Result) bool {
	array := make([]*User, 0)
	for _, value := range data {
		var info = new(User)
		for key, value := range value.Map() {
			switch key {
			case "_id":
				info.UID, _ = primitive.ObjectIDFromHex(value.String())
			case "id":
				info.ID = uint64(value.Int())
			case "base":
				for key, value := range value.Map() {
					switch key {
					case "created_at":
						info.CreatedTime = value.Time()
					case "updatedAt":
						info.UpdatedTime = value.Time()
					case "delete_at":
						info.DeleteTime = value.Time()
					case "name":
						info.Name = value.String()
					}
				}
			case "phone":
				info.Phone = value.String()
			case "sex":
				info.Sex = uint8(value.Uint())
			case "email":
				info.Email = value.String()
			case "status":
				info.Status = uint8(value.Uint())
			case "psw":
				info.Password = value.String()
			case "port":
				info.Portrait = value.String()
			case "wechat":
				info.Wechat = value.String()
			}
		}
		array = append(array, info)
	}
	err := dropUser()
	if err != nil {
		return false
	}
	for i := 0; i < len(array); i++ {
		_ = CreateUser(array[i])
	}
	return true
}
