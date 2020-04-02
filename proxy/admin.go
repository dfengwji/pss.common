package proxy

import (
	"context"
	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Admin struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Role        uint8              `json:"role" bson:"role"`
	Password    string             `json:"psw" bson:"psw"`
	Phone       string             `json:"phone" bson:"phone"`
	Email       string             `json:"email" bson:"email"`
	Scene       string             `json:"scene" bson:"scene"`
}

func CreateAdmin(info *Admin) (error, string) {
	_, err := insertOne(TableAdmin, info)
	if err != nil {
		return err, ""
	}
	return nil, info.UID.Hex()
}

func GetAdminNextID() uint64 {
	num, _ := getSequenceNext(TableAdmin)
	return num
}

func UpdateAdminPassword(uid string, psw string) error {
	msg := bson.M{"psw": psw, "updatedAt": time.Now()}
	_, err := updateOne(TableAdmin, uid, msg)
	return err
}

func UpdateAdminScene(uid string, scene string) error {
	msg := bson.M{"scene": scene, "updatedAt": time.Now()}
	_, err := updateOne(TableAdmin, uid, msg)
	return err
}

func UpdateAdminInfo(uid string, name string, phone string, email string) error {
	msg := bson.M{"phone": phone, "email": email, "name": name, "updatedAt": time.Now()}
	_, err := updateOne(TableAdmin, uid, msg)
	return err
}

func GetAdmin(uid string) (*Admin, error) {
	admin := new(Admin)
	result, err := findOne(TableAdmin, uid)
	if err != nil {
		return nil, err
	}
	err1 := result.Decode(admin)
	if err1 != nil {
		return nil, err1
	}
	return admin, nil
}

func GetAllAdmins() ([]*Admin, error) {
	var items = make([]*Admin, 0, 100)
	cursor, err1 := findAll(TableAdmin, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	/*err := cursor.All(context.Background(), items)
	if err != nil {
		warn(err.Error())
		return nil
	}*/
	for cursor.Next(context.Background()) {
		var node = new(Admin)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func RemoveAdmin(uid string) error {
	_, err := removeOne(TableAdmin, uid)
	return err
}

func dropAdmin() error {
	err := dropOne(TableAdmin)
	return err
}

func recoverAdmins(data []gjson.Result) bool {
	admins := make([]*Admin, 0)
	for _, value := range data {
		var admin = new(Admin)
		for key, value := range value.Map() {
			switch key {
			case "uid":
				admin.UID, _ = primitive.ObjectIDFromHex(value.String())
			case "base.id":
				admin.ID = value.Uint()
			case "base.created_at":
				admin.CreatedTime = value.Time()
			case "updatedAt":
				admin.UpdatedTime = value.Time()
			case "base.delete_at":
				admin.DeleteTime = value.Time()
			case "name":
				admin.Name = value.String()
			case "role":
				admin.Role = uint8(value.Uint())
			case "psw":
				admin.Password = value.String()
			case "phone":
				admin.Phone = value.String()
			case "email":
				admin.Email = value.String()
			case "organize":
				admin.Scene = value.String()
			}
		}
		admins = append(admins, admin)
	}
	err := dropAdmin()
	if err != nil {
		return false
	}
	for i := 0; i < len(admins); i++ {
		CreateAdmin(admins[i])
	}
	return true
}
