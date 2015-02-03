package service

import (
	"github.com/kahoona77/gotv/domain"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

const database string = "xtv"

type MongoService struct {
	Session *mgo.Session
}

func CreateMongoService() *MongoService {
	ms := new(MongoService)

	//creating db
	session, err := mgo.Dial("localhost") //mgo.Dial("192.168.56.101") //mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}

	ms.Session = session
	return ms
}

func (ms *MongoService) Close() {
	ms.Session.Close()
}

func (ms *MongoService) GetRepo(collection string) *GoTvRepository {
	return NewRepository(ms.Session, collection)
}

// ++++++ GoTvRepository +++++++++++

type GoTvRepository struct {
	Collection *mgo.Collection
}

func NewRepository(session *mgo.Session, collectionName string) *GoTvRepository {
	repo := new(GoTvRepository)
	repo.Collection = session.DB("xtv").C(collectionName)
	return repo
}

func (this GoTvRepository) All(results interface{}) error {
	return this.Collection.Find(nil).All(results)
}

func (this GoTvRepository) CountAll() (int, error) {
	return this.Collection.Find(nil).Count()
}

func (this GoTvRepository) FindWithQuery(query *bson.M, results interface{}) error {
	return this.Collection.Find(query).All(results)
}

func (this GoTvRepository) FindById(docId string, result domain.MongoDomain) error {
	return this.Collection.FindId(docId).One(result)
}

func (this GoTvRepository) FindFirst(result domain.MongoDomain) error {
	return this.Collection.Find(nil).One(result)
}

func (this GoTvRepository) Remove(docId string) error {
	return this.Collection.RemoveId(docId)
}

func (this GoTvRepository) RemoveAll(query *bson.M) (info *mgo.ChangeInfo, err error) {
	return this.Collection.RemoveAll(query)
}

func (this GoTvRepository) Save(docId string, doc domain.MongoDomain) (info *mgo.ChangeInfo, err error) {
	if docId == "" {
		docId = bson.NewObjectId().Hex()
		doc.SetId(docId)
	}
	return this.Collection.UpsertId(docId, doc)
}
