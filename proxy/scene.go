package proxy

import (
	"context"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Scene struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`

	Status   uint8    `json:"status" bson:"status"`
	Type     uint8    `json:"type" bson:"type"`
	Name     string   `json:"name" bson:"name"`
	Place    string `json:"place" bson:"place"`
	Address  string   `json:"address" bson:"address"`
	Mater    string   `json:"mater" bson:"master"`
	Phone    string   `json:"phone" bson:"phone"`
	Desc     string   `json:"desc" bson:"desc"`
	Icon     string   `json:"icon" bson:"icon"`
	Teachers []string `json:"teachers" bson:"teachers"`
	Students []string `json:"students" bson:"students"`
	Books    []string `json:"books" bson:"books"`
}

func CreateScene(info *Scene) error {
	_, err := insertOne(TableScene, info)
	if err != nil {
		return err
	}
	return nil
}

func GetSceneNextID() uint64 {
	num, _ := getSequenceNext(TableScene)
	return num
}

func GetScene(uid string) (*Scene, error) {
	result, err := findOne(TableScene, uid)
	if err != nil {
		return nil, err
	}
	model := new(Scene)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetAllScenes() ([]*Scene, error) {
	var items = make([]*Scene, 0, 100)
	cursor, err1 := findAll(TableScene, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var node = new(Scene)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func HadSceneByName(name string) (bool, error) {
	msg := bson.M{"name": name}
	return hadOne(TableScene, msg)
}

func UpdateSceneBase(uid string, name string, master string, phone string, desc string) error {
	msg := bson.M{"name": name, "master": master,
		"phone": phone, "desc": desc, "updatedAt": time.Now()}
	_, err := updateOne(TableScene, uid, msg)
	return err
}

func UpdateSceneIcon(uid string, icon string) error {
	msg := bson.M{"icon": icon, "updatedAt": time.Now()}
	_, err := updateOne(TableScene, uid, msg)
	return err
}

func UpdateSceneStatus(uid string, status uint8) error {
	msg := bson.M{"status": status, "updatedAt": time.Now()}
	_, err := updateOne(TableScene, uid, msg)
	return err
}

func UpdateScenePlace(uid string, place string) error {
	msg := bson.M{"place": place, "updatedAt": time.Now()}
	_, err := updateOne(TableScene, uid, msg)
	return err
}

func RemoveScene(uid string) bool {
	_, err := removeOne(TableScene, uid)
	if err == nil {
		return true
	}
	return false
}

func AppendSceneClass(uid string, class string) error {
	if len(class) < 1 {
		return errors.New("the class uid is empty")
	}
	msg := bson.M{"classes": class}
	_, err := appendElement(TableScene, uid, msg)
	return err
}

func UnbindSceneClass(uid string, class string) error {
	if len(class) < 1 {
		return errors.New("the class uid is empty")
	}
	msg := bson.M{"classes": class}
	_, err := removeElement(TableScene, uid, msg)
	return err
}

func AppendSceneStudent(uid string, student string) error {
	if len(student) < 1 {
		return errors.New("the student uid is empty")
	}
	msg := bson.M{"students": student}
	_, err := appendElement(TableScene, uid, msg)
	return err
}

func UnbindSceneStudent(uid string, student string) error {
	if len(student) < 1 {
		return errors.New("the student uid is empty")
	}
	msg := bson.M{"students": student}
	_, err := removeElement(TableScene, uid, msg)
	return err
}

func AppendSceneTeacher(uid string, teacher string) error {
	if len(teacher) < 1 {
		return errors.New("the teacher uid is empty")
	}
	msg := bson.M{"teachers": teacher}
	_, err := appendElement(TableScene, uid, msg)
	return err
}

func UnbindSceneTeacher(uid string, teacher string) error {
	if len(teacher) < 1 {
		return errors.New("the teacher uid is empty")
	}
	msg := bson.M{"teachers": teacher}
	_, err := removeElement(TableScene, uid, msg)
	return err
}

func AppendSceneBook(uid string, book string) error {
	if len(book) < 1 {
		return errors.New("the book uid is empty")
	}
	msg := bson.M{"books": book}
	_, err := appendElement(TableScene, uid, msg)
	return err
}

func AppendSceneAdmin(uid string, admin string) error {
	if len(admin) < 1 {
		return errors.New("the admin uid is empty")
	}
	msg := bson.M{"admins": admin}
	_, err := appendElement(TableScene, uid, msg)
	return err
}

func UnbindSceneAdmin(uid string, admin string) error {
	if len(admin) < 1 {
		return errors.New("the admin uid is empty")
	}
	msg := bson.M{"admins": admin}
	_, err := removeElement(TableScene, uid, msg)
	return err
}
