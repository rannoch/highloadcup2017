package storage

var Db IStorage

type IStorage interface {
	InsertEntity(entity Entity) error
	InsertEntityMultiple(entities []Entity) error
	SelectEntity(entity Entity, conditions []Condition) error
	SelectEntityMultiple(entities interface{}, fields []string, joins []Join, conditions []Condition, sort Sort) error
	GetAverage(entity Entity, avgColumn string, tableJoins []Join, conditions []Condition) (average float32, err error)
	UpdateEntity(entity Entity, params map[string]interface{}, conditions []Condition) (rowsAffected int64, err error)
}

type Entity interface {
	GetId() int32
	TableName() string
	GetFields(alias string) []string
	GetValues() []interface{}
	GetFieldPointers(with []string) []interface{}
	ValidateParams(params map[string]interface{}, scenario string) (result bool)
	SetParams(params map[string]interface{})
	HasForeignRelations() bool
}

type Condition struct {
	Param         string
	Value         string
	Operator      string
	JoinCondition string
}

type Join struct {
	Name      string
	Type      string
	Condition Condition
}

type Sort struct {
	Fields []string
	Direction string
}