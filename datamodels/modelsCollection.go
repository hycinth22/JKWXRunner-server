package datamodels

var ModelsCollection []interface{} = nil

func init() {
	ModelsCollection = append(ModelsCollection,
		&Account{},
		&AccountLog{},
		&CacheUserInfo{},
		&CacheUserSportResult{},
		&Device{},
		&Token{},
		&UserIDRelation{},
	)
}
