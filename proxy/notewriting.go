package proxy

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type NoteWriting struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`

	Pen      uint64 `json:"pen" bson:"pen"`
	Writer     uint64 `json:"writer" bson:"writer"`
	Book     string `json:"book" bson:"book"`
	Page     uint16 `json:"page" bson:"page"`
	DotBook  uint64 `json:"dotBook" bson:"dotBook"`
	DotPage  uint16 `json:"dotPage" bson:"dotPage"`
	DotStamp uint64 `json:"dotStamp" bson:"dotStamp"`
	DotNum   uint32 `json:"dotNum" bson:"dotNum"`
	Duration uint16 `json:"duration" bson:"duration"`
	Snapshot string `json:"snapshot" bson:"snapshot"`
	Color    string `json:"color" bson:"color"`
	Dots     string `json:"dots" bson:"dots"`
	Bytes 	[]byte `json:"bytes" bson:"bytes"`
}

func CreateNoteWriting(info *NoteWriting) error {
	_, err := insertOne(TableNoteWriting, info)
	if err != nil {
		return err
	}
	return nil
}

func GetNoteWritingNextID() uint64 {
	num, _ := getSequenceNext(TableNoteWriting)
	return num
}

func GetNoteWriting(uid string) (*NoteWriting, error) {
	result, err := findOne(TableNoteWriting, uid)
	if err != nil {
		return nil, err
	}
	model := new(NoteWriting)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetNoteWritingByPage(writer uint64, book string, page uint16) (*NoteWriting, error) {
	filter := bson.M{"writer": writer, "book": book, "page": page, "deleteAt": new(time.Time)}
	result, err := findOneBy(TableNoteWriting, filter)
	if err != nil {
		return nil, err
	}
	model := new(NoteWriting)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetNoteWritingsByTime(writer uint64, book string, from, to time.Time) ([]*NoteWriting, error) {
	filter := bson.M{"writer": writer, "book": book, "createdAt": bson.M{"$gte": from , "$lte": to}}
	cursor, err := findMany(TableNoteWriting, filter, 0)
	if err != nil {
		return nil, err
	}
	var items = make([]*NoteWriting, 0, 20)
	for cursor.Next(context.Background()) {
		var node NoteWriting
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, &node)
		}
	}
	return items, nil
}

func GetNoteWritingsInDuration(writer uint64, from, to time.Time) ([]*NoteWriting, error) {
	filter := bson.M{"writer": writer, "createdAt": bson.M{"$gte": from , "$lte": to}}
	cursor, err := findMany(TableNoteWriting, filter, 0)
	if err != nil {
		return nil, err
	}
	var items = make([]*NoteWriting, 0, 20)
	for cursor.Next(context.Background()) {
		var node NoteWriting
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, &node)
		}
	}
	return items, nil
}

func GetNoteWritingsByWriter(writer uint64, book string) ([]*NoteWriting, error) {
	filter := bson.M{"writer": writer, "book": book, "deleteAt": new(time.Time)}
	cursor, err := findMany(TableNoteWriting, filter, 20)
	if err != nil {
		return nil, err
	}
	var items = make([]*NoteWriting, 0, 20)
	for cursor.Next(context.Background()) {
		var node NoteWriting
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, &node)
		}
	}
	return items, nil
}

func GetAllNoteWritings() ([]*NoteWriting, error) {
	cursor, err1 := findAll(TableNoteWriting, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*NoteWriting, 0)
	for cursor.Next(context.Background()) {
		var node = new(NoteWriting)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func RemoveNoteWriting(uid string) error {
	_, err := removeOne(TableNoteWriting, uid)
	return err
}

func UpdateNoteWritingSnapshot(uid string, snapshot string) error {
	msg := bson.M{"snapshot": snapshot, "updatedAt": time.Now()}
	_, err := updateOne(TableNoteWriting, uid, msg)
	return err
}

func UpdateNoteWritingDots(uid string, stamp uint64, num uint32, dots string, duration uint16) error {
	msg := bson.M{"dotStamp": stamp, "dotNum": num, "dots": dots, "duration": duration, "updatedAt": time.Now()}
	_, err := updateOne(TableNoteWriting, uid, msg)
	return err
}

func UpdateNoteWritingBytes(uid string, stamp uint64, num uint32, dots []byte, duration uint16) error {
	msg := bson.M{"dotStamp": stamp, "dotNum": num, "bytes": dots, "duration": duration, "updatedAt": time.Now()}
	_, err := updateOne(TableNoteWriting, uid, msg)
	return err
}
