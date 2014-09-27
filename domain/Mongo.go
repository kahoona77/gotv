package domain

import (
  "labix.org/v2/mgo"
  "labix.org/v2/mgo/bson"
  )

type GoTvRepository  struct {
  Collection *mgo.Collection
}

func NewRepository(session  *mgo.Session, collectionName string) *GoTvRepository {
    repo := new(GoTvRepository)
    repo.Collection = session.DB("xtv").C(collectionName)
    return repo
}

func (this GoTvRepository) All(results interface{}) error{
  return this.Collection.Find(nil).All(results)
}

func (this GoTvRepository) FindById(docId string, result MongoDomain) error{
  return this.Collection.FindId(docId).One(result)
}

func (this GoTvRepository) Remove(docId string) error{
  return this.Collection.RemoveId(docId)
}

func (this GoTvRepository) Save(docId string,doc MongoDomain) (info *mgo.ChangeInfo, err error) {
  if (docId == "") {
    docId = bson.NewObjectId().Hex()
    doc.SetId (docId)
  }

  return this.Collection.UpsertId(docId, doc)
}
