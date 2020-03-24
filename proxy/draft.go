package proxy

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type CourseDraft struct {
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
	Tags 		[]string 			`json:"tags" bson:"tags"`
	Audios      []string           `json:"audios" bson:"audios"`
	OpenTargets	[]string 			`json:"targets" bson:"targets"`
}

func CreateCourseDraft(info *CourseDraft) error {
	_, err := insertOne(TableCourseDraft, info)
	if err != nil {
		return err
	}
	return nil
}

func GetAllCourseDrafts() ([]*CourseDraft, error) {
	cursor, err1 := findAll(TableCourseDraft, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*CourseDraft, 0, 10000)
	for cursor.Next(context.Background()) {
		var node = new(CourseDraft)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetCourseDraftNextID() uint64 {
	num, _ := getSequenceNext(TableCourseDraft)
	return num
}

func RemoveCourseDraft(uid string) error {
	_, err := removeOne(TableCourseDraft, uid)
	return err
}

func GetCourseDraft(uid string) (*CourseDraft, error) {
	result, err := findOne(TableCourseDraft, uid)
	if err != nil {
		return nil, err
	}
	model := new(CourseDraft)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetCourseDraftsByAuthor(author string) ([]*CourseDraft, error) {
	msg := bson.M{"author": author, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableCourseDraft, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*CourseDraft, 0, 200)
	for cursor.Next(context.Background()) {
		var node = new(CourseDraft)
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
	_, err := updateOne(TableCourseDraft, uid, msg)
	return err
}

func UpdateCourseDraftVideo(uid string, video string) error {
	msg := bson.M{"video": video, "updatedAt": time.Now()}
	_, err := updateOne(TableCourseDraft, uid, msg)
	return err
}

func UpdateCourseDraftTags(uid string, tags []string) error {
	msg := bson.M{"tags": tags, "updatedAt": time.Now()}
	_, err := updateOne(TableCourseDraft, uid, msg)
	return err
}

func UpdateCourseDraftCover(uid string, cover string) error {
	msg := bson.M{"cover": cover, "updatedAt": time.Now()}
	_, err := updateOne(TableCourseDraft, uid, msg)
	return err
}

func UpdateCourseDraftTask(uid string, task string) error {
	msg := bson.M{"task": task, "updatedAt": time.Now()}
	_, err := updateOne(TableCourseDraft, uid, msg)
	return err
}

func UpdateCourseDraftTargets(uid string,open uint8, targets []string) error {
	msg := bson.M{"open":open, "targets": targets, "updatedAt": time.Now()}
	_, err := updateOne(TableCourseDraft, uid, msg)
	return err
}


