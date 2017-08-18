package storage

type IStorage interface {
	InsertEntity(entity Entity) error
	InsertEntityMultiple(entities []Entity) error
	SelectEntity(entity Entity, conditions []Condition) error
	SelectEntityMultiple(entities interface{}, conditions []Condition) error
}

type Entity interface {
	GetId() int32
	TableName() string
	GetFields() []string
	GetValues() []interface{}
	GetFieldPointers() []interface{}
}

type Condition struct {
	Param string
	Value string
	Operator string
	Join  string
}