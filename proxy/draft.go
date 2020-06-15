package proxy

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type VideoDraft struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Name        string             `json:"name" bson:"name"`
	Remark      string             `json:"remark" bson:"remark"`
	Author      string             `json:"author" bson:"author"`
	Cover       string             `json:"cover" bson:"cover"`
	Video       string             `json:"video" bson:"video"`
	Duration    uint32             `json:"duration" bson:"duration"`
	Menu        string             `json:"menu" bson:"menu"`
	Start       uint64             `json:"start" bson:"start"`
	TaskUID     string             `json:"task" bson:"task"`
	Open        uint8 				`json:"open" bson:"open"`
	AutoFit     uint8 				`json:"fit" bson:"fit"`
	WaterMark   uint8 				`json:"watermark" bson:"watermark"`
	Tags 		[]string 			`json:"tags" bson:"tags"`
	Audios      []string           `json:"audios" bson:"audios"`
	OpenTargets	[]string 			`json:"targets" bson:"targets"`
	Events      []DraftEvent 		`json:"events" bson:"events"`
}

type DraftEvent struct {
	Type uint8 `json:"type" bson:"type"`
	Stamp uint64 `json:"stamp" bson:"stamp"`
	URL string `json:"url" bson:"url"`
	X float32 `json:"x" bson:"x"`
	Y float32 `json:"y" bson:"y"`
	Width float32 `json:"width" bson:"width"`
	Height float32 `json:"height" bson:"height"`
}

func CreateCourseDraft(info *VideoDraft) error {
	_, err := insertOne(TableVideoDraft, info)
	if err != nil {
		return err
	}
	return nil
}

func GetAllCourseDrafts() ([]*VideoDraft, error) {
	cursor, err1 := findAll(TableVideoDraft, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*VideoDraft, 0, 10000)
	for cursor.Next(context.Background()) {
		var node = new(VideoDraft)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetCourseDraftNextID() uint64 {
	num, _ := getSequenceNext(TableVideoDraft)
	return num
}

func RemoveCourseDraft(uid string) error {
	_, err := removeOne(TableVideoDraft, uid)
	return err
}

func GetCourseDraft(uid string) (*VideoDraft, error) {
	result, err := findOne(TableVideoDraft, uid)
	if err != nil {
		return nil, err
	}
	model := new(VideoDraft)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetCourseDraftsByAuthor(author string) ([]*VideoDraft, error) {
	msg := bson.M{"author": author, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableVideoDraft, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*VideoDraft, 0, 200)
	for cursor.Next(context.Background()) {
		var node = new(VideoDraft)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateCourseDraftBase(uid string, name string, remark string) error {
	msg := bson.M{"name": name, "remark": remark, "updatedAt": time.Now()}
	_, err := updateOne(TableVideoDraft, uid, msg)
	return err
}

func UpdateCourseDraftVideo(uid string, video string) error {
	msg := bson.M{"video": video, "updatedAt": time.Now()}
	_, err := updateOne(TableVideoDraft, uid, msg)
	return err
}

func UpdateCourseDraftMenu(uid string, menu string) error {
	msg := bson.M{"menu": menu, "updatedAt": time.Now()}
	_, err := updateOne(TableVideoDraft, uid, msg)
	return err
}

func UpdateCourseDraftTags(uid string, tags []string) error {
	msg := bson.M{"tags": tags, "updatedAt": time.Now()}
	_, err := updateOne(TableVideoDraft, uid, msg)
	return err
}

func UpdateCourseDraftCover(uid string, cover string) error {
	msg := bson.M{"cover": cover, "updatedAt": time.Now()}
	_, err := updateOne(TableVideoDraft, uid, msg)
	return err
}

func UpdateCourseDraftTask(uid string, task string) error {
	msg := bson.M{"task": task, "updatedAt": time.Now()}
	_, err := updateOne(TableVideoDraft, uid, msg)
	return err
}

func UpdateCourseDraftTargets(uid string,open uint8, targets []string) error {
	msg := bson.M{"open":open, "targets": targets, "updatedAt": time.Now()}
	_, err := updateOne(TableVideoDraft, uid, msg)
	return err
}


