package dumpdatabase

import (
	"env"
)

type Connection struct {
	Hostname         string
	DatabaseUser     string
	DatabasePassword string
	DatabaseName     string
	TableName        string
	WhereCondition   string
	BackUpPath       string
}

func (db *Connection) PopulateFromEnv() {
	envData := &env.Env{}
	envData.SetEnvFileName(".env").Load()
	defer func() {
		envData = nil
	}()
	db.Hostname = envData.Getenv("HOSTNAME")
	db.DatabaseUser = envData.Getenv("DATABASE_USER")
	db.DatabasePassword = envData.Getenv("DATABASE_PASSWORD")
	db.DatabaseName = envData.Getenv("DATABASE_NAME")
	db.BackUpPath = envData.Getenv("BACKUP_PATH")
}
