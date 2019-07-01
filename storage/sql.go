package storage

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"
	"github.com/teamnsrg/mida/jstrace"
	"github.com/teamnsrg/mida/log"
)

type SqlCall struct {
	gorm.Model

	CrawlID   int `gorm:"not null"`
	IsolateID int `gorm:"not null"`
	ScriptID  int `gorm:"not null"`
}

// CreatePostgresConnection connects to the postgres server and creates the specified database, if it does not
// already exist.
func CreatePostgresConnection(host string, port string, dbName string) (*gorm.DB, error) {

	log.Log.Info("Attempting connection:")
	log.Log.Infof("Host: [ %s ]", host)
	log.Log.Infof("Port: [ %s ]", port)
	log.Log.Infof("DB: [ %s ]", dbName)
	log.Log.Infof("Username: [ %s ]", viper.GetString("postgresuser"))
	log.Log.Infof("Password: [ %s ]", viper.GetString("postgrespass"))

	db, err := gorm.Open("postgres",
		"host="+host+
			" port="+port+
			" user="+viper.GetString("postgresuser")+
			" dbname="+"postgres"+
			" password="+viper.GetString("postgrespass"))
	if err != nil {
		return nil, err
	}

	// This will error if the database already exists. That's okay - we are going to connect to it anyway
	db = db.Exec("CREATE DATABASE " + dbName + " WITH TEMPLATE mida_template OWNER mida;")

	db, err = gorm.Open("postgres",
		"host=" + host +
			" port=" + port +
			" user="+ viper.GetString("postgresuser") +
			" dbname=" + dbName +
			" password=" + viper.GetString("postgrespass") )
	if err != nil {
		return nil, err
	}
	return db, nil

}

func ClosePostgresConnection(db *gorm.DB) error {
	return db.Close()
}

func StoreJSTraceToDB(db *gorm.DB, trace jstrace.JSTrace) error {
	return nil
}