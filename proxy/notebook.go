package proxy

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type NoteBook struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	Name        string             `json:"name" bson:"name"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`

	Used uint8 			`json:"used" bson:"used"`
	Owner string 		`json:"owner" bson:"owner"`
	Style  string 		`json:"style" bson:"style"`
	Remark string 		`json:"remark" bson:"remark"`
}

func CreateNoteBook(info *NoteBook) error {
	_, err := insertOne(TableNoteBook, info)
	if err != nil {
		return err
	}
	return nil
}

func GetNoteBookNextID() uint64 {
	num, _ := getSequenceNext(TableBookID)
	return num
}

func GetNoteBookCount() int64 {
	num, _ := getCount(TableNoteBook)
	return num
}

func GetNoteBook(uid string) (*NoteBook, error) {
	result, err := findOne(TableNoteBook, uid)
	if err != nil {
		return nil, err
	}
	model := new(NoteBook)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetNoteBookByID(id uint64) (*NoteBook, error) {
	msg := bson.M{"id": id}
	result, err := findOneBy(TableNoteBook, msg)
	if err != nil {
		return nil, err
	}
	model := new(NoteBook)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func RemoveNoteBook(uid string) error {
	_, err := removeOne(TableNoteBook, uid)
	return err
}

func GetAllNoteBooks() ([]*NoteBook, error) {
	cursor, err1 := findAll(TableNoteBook, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*NoteBook, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(NoteBook)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetNoteBooksByOwner(owner string) ([]*NoteBook, error) {
	cursor, err1 := findMany(TableNoteBook, bson.M{"owner": owner}, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*NoteBook, 0, 200)
	for cursor.Next(context.Background()) {
		var node = new(NoteBook)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateNoteBookUsed(uid string, used uint8) error {
	msg := bson.M{"used": used, "updatedAt": time.Now()}
	_, err := updateOne(TableNoteBook, uid, msg)
	return err
}

func UpdateNoteBookStyle(uid string, style string) error {
	msg := bson.M{"style": style, "updatedAt": time.Now()}
	_, err := updateOne(TableNoteBook, uid, msg)
	return err
}

func UpdateNoteBookBase(uid, name, remark string) error {
	msg := bson.M{"name": name, "remark": remark, "updatedAt": time.Now()}
	_, err := updateOne(TableNoteBook, uid, msg)
	return err
}

func UpdateNoteBookOwner(uid string, owner string) error {
	msg := bson.M{"owner": owner, "updatedAt": time.Now()}
	_, err := updateOne(TableNoteBook, uid, msg)
	return err
}
