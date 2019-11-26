package dumpdatabase

import (
	"helper"
	"time"
	"fmt"
	"os"
	"os/exec"
	"compress/gzip"
)

type DumpDB struct {
	connection *Connection
	dumpCommand string
	dumpCommandPath string
	helper *helper.Helper
	fileName string
	fullPath string
	fileHandler *gzip.Writer
	filePermission os.FileMode
}

func (dump *DumpDB) init() {
	dump.connection = &Connection{}
	dump.helper = &helper.Helper{}
	dump.connection.PopulateFromEnv()
	dump.dumpCommand = "mysqldump"
	dump.filePermission = 0744
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

func (dump *DumpDB) setDumpCommandPath() {
	dumpCommandPath, err := exec.LookPath(dump.dumpCommand)
	dump.helper.LogError(err)
	dump.dumpCommandPath = dumpCommandPath
}

func (dump *DumpDB) generateFileName(file string) {
	today := time.Now()
	date := fmt.Sprintf("_%d_%02d_%02d", today.Year(), today.Month(), today.Day())
	dump.fileName = string(file) + date + ".sql.gzip"
	dump.fullPath = string(dump.connection.BackUpPath) + string(os.PathSeparator) + dump.fileName
}

func (dump *DumpDB) openFile()  {
	file, err := os.OpenFile(dump.fullPath, os.O_RDWR|os.O_CREATE, dump.filePermission)
	dump.helper.LogError(err)
	gzipFile := gzip.NewWriter(file)
	gzipFile.Name = dump.fileName
	gzipFile.ModTime = time.Date(1977, time.May, 25, 0, 0, 0, 0, time.UTC)
	dump.fileHandler = gzipFile
}

func (dump *DumpDB) closeFile() *DumpDB {
	defer func() {
		dump.fileHandler = nil
	}()
	dump.fileHandler.Close()
	return dump
}

func (dump *DumpDB) dumpData()  {
	dump.openFile()
	defer dump.closeFile()
	// Create the Command with no Where Clause
	command := exec.Command(dump.dumpCommandPath, "-h", dump.connection.Hostname, "-u", dump.connection.DatabaseUser, "-p"+string(dump.connection.DatabasePassword),
		dump.connection.DatabaseName, dump.connection.TableName, "--no-create-info",
	)

	// Create the Command with Where Clause if Where Condition Present
	if len(string(dump.connection.WhereCondition)) != 0 {
		command = exec.Command(dump.dumpCommandPath, "-h", dump.connection.Hostname, "-u", dump.connection.DatabaseUser, "-p"+string(dump.connection.DatabasePassword),
			dump.connection.DatabaseName, dump.connection.TableName, "--where="+string(dump.connection.WhereCondition), "--no-create-info",
		)
	}
	// Attach Filehandler to command's STDIO
	command.Stdout = dump.fileHandler
	dump.helper.LogError(command.Start())
	command.Wait()
}

func (dump *DumpDB) DumpDatabase()  {
	// Dump the Query to .sql file
	dump.setDumpCommandPath()
	dump.generateFileName(dump.connection.TableName)
	dump.dumpData()
}
