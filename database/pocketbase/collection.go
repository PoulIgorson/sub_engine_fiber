package pocketbase

import (
	"encoding/json"
	"log"
	"reflect"

	. "github.com/PoulIgorson/sub_engine_fiber/database/errors"
	. "github.com/PoulIgorson/sub_engine_fiber/database/interfaces"
	. "github.com/PoulIgorson/sub_engine_fiber/define"
)

var _ Table = &Collection{}

type Collection struct {
	db   *DataBase
	name string

	model   Model
	Objects ManagerI
}

func (collection *Collection) DB() DB {
	return collection.db
}

func (collection *Collection) Name() string {
	return collection.name
}

func (collection *Collection) Model() Model {
	return collection.model
}

func (collection *Collection) Get(idI any) (Model, error) {
	id, ok := idI.(string)
	if !ok {
		return nil, NewErrorf("pb: id must be string")
	}
	records, err := collection.db.pb.Filter(collection.name, map[string]any{"id": id})
	if err != nil {
		return nil, ToError(err)
	}
	if len(records) == 0 {
		return nil, NewErrorf("pb: record not found")
	}
	dataByte, _ := json.Marshal(records[0].data)
	model := collection.model.Create(collection.db, string(dataByte))
	collection.Objects.Store(model.Id().(string), model)
	return model, nil
}

func (collection *Collection) Save(model Model) error {
	dataByte, _ := json.Marshal(model)
	data := map[string]any{}
	json.Unmarshal(dataByte, &data)
	form := NewForm(collection.db.pb, NewRecord(collection.name, collection.db.pb))
	form.LoadData(data)
	id, err := form.Submit()
	if err != nil {
		return ToError(err)
	}

	field_id, err := Check(model, "ID")
	if err != nil {
		return NewErrorf("pb: " + err.Error())
	}
	field_id.Set(reflect.ValueOf(id))
	if collection.Objects.IsInstance() {
		collection.Objects.Store(model.Id(), model)
	}
	return nil
}

func (collection *Collection) Delete(idI any) error {
	id, ok := idI.(string)
	if !ok {
		return NewErrorf("pb: id must be string")
	}
	return ToError(collection.db.pb.Delete(collection.name, id))
}

// Pocketbase does not support DeleteAll
// All models will be deletting of one
func (collection *Collection) DeleteAll() error {
	log.Println("pocketbase does not support DeleteAll\nAll models will be deletting of one")
	for _, model := range collection.Objects.All() {
		model.Delete(collection.db)
	}
	return nil
}

func (collection *Collection) Count() uint {
	records, err := collection.db.pb.Filter(collection.name, map[string]any{})
	if err != nil {
		return 0
	}
	return uint(len(records))
}

func (collection *Collection) Manager() ManagerI {
	return collection.Objects
}

func (collection *Collection) SetManager(newManager ManagerI) {
	collection.Objects = newManager
}
