package bellbox

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_"github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

func (c DbConfig) GetDb() *gorm.DB {
	if db != nil {
		return db
	}
	connectionString := fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s", c.Host, c.Port, c.User, c.DbName, c.Password)
	var err error
	db, err = gorm.Open("postgres", connectionString)
	if err != nil {
		panic("Could not make connection to postgres db with connection details: " + err.Error())
	}
	db.AutoMigrate(User{}, UserToken{}, Bell{}, Bellringer{}, Message{})
	return db
}
