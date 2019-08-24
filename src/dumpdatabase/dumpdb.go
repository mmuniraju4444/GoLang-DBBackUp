package dumpdatabase

import (
	"helper"
	"time"
	"fmt"
	"os"
	"os/exec"
)

type DumpDB struct {
	connection *Connection
	dumpCommand string
	dumpCommandPath string
	helper *helper.Helper
	fullPath string
	fileHandler *os.File
	filePermission os.FileMode
}

func (dump *DumpDB) InitDump() *DumpDB {
	dump.init()
	return dump
}

func (dump *DumpDB) SetWhereCondition(condition string) *DumpDB {
	dump.connection.WhereCondition = condition
	return dump
}

func (dump *DumpDB) SetTableName(tableName string) *DumpDB {
	dump.connection.TableName = tableName
	return dump
}

func (dump *DumpDB) init() {
	dump.connection = &Connection{}
	dump.helper = &helper.Helper{}
	dump.connection.PopulateFromEnv()
	dump.dumpCommand = "mysqldump"
	dump.filePermission = 0777
}

func (dump *DumpDB) setdumpCommandPath() {
	dumpCommandPath, err := exec.LookPath(dump.dumpCommand)
	dump.helper.LogError(err)
	dump.dumpCommandPath = dumpCommandPath
}

func (dump *DumpDB) generateFileName(file string) {
	today := time.Now()
	date := fmt.Sprintf("_%d_%02d_%02d", today.Year(), today.Month(), today.Day())
	fileName := string(file) + date + ".sql"
	dump.fullPath = string(dump.connection.BackUpPath) +  "/" + fileName
}

func (dump *DumpDB) openFile()  {
	fileHandler, err := os.OpenFile(dump.fullPath, os.O_RDWR|os.O_CREATE, dump.filePermission)
	dump.helper.LogError(err)
	dump.fileHandler = fileHandler
}

func (dump *DumpDB) closeFile() *DumpDB {
	defer func() {
		dump.fileHandler = nil
	}()
	dump.fileHandler.Close()
	return dump
}

func (dump *DumpDB) dumpData()  {
	command := exec.Command(dump.dumpCommandPath, "-h", dump.connection.Hostname, "-u", dump.connection.DatabaseUser, "-p"+string(dump.connection.DatabasePassword),
		dump.connection.DatabaseName, dump.connection.TableName, "--no-create-info",
	)
	if len(string(dump.connection.WhereCondition)) != 0 {
		command = exec.Command(dump.dumpCommandPath, "-h", dump.connection.Hostname, "-u", dump.connection.DatabaseUser, "-p"+string(dump.connection.DatabasePassword),
			dump.connection.DatabaseName, dump.connection.TableName, "--where="+string(dump.connection.WhereCondition), "--no-create-info",
		)
	}
	command.Stdout = dump.fileHandler
	dump.helper.LogError(command.Start())
	command.Wait()
}

func (dump *DumpDB) DumpDatabase()  {
	dump.setdumpCommandPath()
	dump.generateFileName(dump.connection.TableName)
	dump.openFile()
	defer dump.closeFile()
	dump.dumpData()
}
