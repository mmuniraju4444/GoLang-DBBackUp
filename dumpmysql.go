package main

import (
	"dumpdatabase"
	"env"
	"encoding/json"
	"helper"
	"sync"
	"fmt"
	"strings"
)
type BackupData struct {
	envData *env.Env
	data map[string]string
	helper *helper.Helper
}
var backupData *BackupData

func Init() {
	backupData = new(BackupData)
	backupData.envData = &env.Env{}
	backupData.envData.SetEnvFileName(".env")
	backupData.envData.Load()
	backupData.helper = &helper.Helper{}
}
var wg sync.WaitGroup
func main() {
	Init() // Initiate the BackupData Class
	table := backupData.envData.Getenv("TABLE_CONDITION")
	err := json.Unmarshal([]byte(table), &backupData.data)
	backupData.helper.LogError(err)
	// Loop all the table in goroutines
	for tableName, condition := range backupData.data {
		wg.Add(1)
		go DumpData(
			strings.TrimSpace(tableName), 
			strings.TrimSpace(condition))
	}
	wg.Wait()
}

func DumpData(tableName string, condition string) {
	defer wg.Done()
	backupData.helper.LogInfo(fmt.Sprintf("BackUp Started for %s ",tableName))
	dump :=&dumpdatabase.DumpDB{}
	dump.InitDump()
	dump.SetTableName(string(tableName))
	dump.SetWhereCondition(string(condition))
	dump.DumpDatabase()
	backupData.helper.LogInfo(fmt.Sprintf("BackUp End for %s ",tableName))
}
