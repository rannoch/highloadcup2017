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
