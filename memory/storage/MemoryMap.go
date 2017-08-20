package storage

type MemoryMap map[string]EntityMap

type EntityMap map[int32]interface{}

var Db MemoryMap

// todo later
func InitMemoryMap() {
	Db = MemoryMap{
		"user" : EntityMap{},
		"location" : EntityMap{},
		"visit" : EntityMap{},
	}
}

/*
func (memoryMap MemoryMap) InsertEntity(entity Entity) (err error) {
	// добавляю
	memoryMap[entity.TableName()][entity.GetId()] = entity

	// если связная сущность, обновляю связанные сущности
	if entity.HasForeignRelations() {
		entityValue := reflect.ValueOf(entity)

		for i := 0; i < entityValue.NumField(); i++ {
			*/
/*field := entityValue.Type().Field(i)

			relation := field.Tag.Get("relation")
			value := entityValue.Field(i).Interface().(Entity)

			if relation != "" {
				memoryMap[relation][value.GetId()]
			}*//*

		}
	}

	return
}

func (memoryMap MemoryMap) InsertEntityMultiple(entities []Entity) error {
	panic("implement me")
}

func (memoryMap MemoryMap) SelectEntity(entity Entity, conditions []Condition) error {
	panic("implement me")
}

func (memoryMap MemoryMap) SelectEntityMultiple(entities interface{}, fields []string, joins []Join, conditions []Condition, sort Sort) error {
	panic("implement me")
}

func (memoryMap MemoryMap) GetAverage(entity Entity, avgColumn string, tableJoins []Join, conditions []Condition) (average float32, err error) {
	panic("implement me")
}

func (memoryMap MemoryMap) UpdateEntity(entity Entity, params map[string]interface{}, conditions []Condition) (rowsAffected int64, err error) {
	panic("implement me")
}
*/
