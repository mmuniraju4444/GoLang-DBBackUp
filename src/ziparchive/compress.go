package ziparchive

import (
	"compress/gzip"
	"os"
	"helper"
	"io"
	"time"
)

type Compress struct {
	FilePath string
	FileName string
	CopyFilelocation string
	SaveLocaion string
	DeleteOriginalFile bool
	fileExtention string

	helper *helper.Helper
}

func (compress *Compress) init() {
	compress.helper = &helper.Helper{}
	compress.fileExtention = ".gzip"
}

func (compress *Compress) CompressFile() {
	compress.init() // Initiate Compress Class Data
	compress.copyData()
	// Remove Original File if DeleteOriginalFile flag set
	if compress.DeleteOriginalFile == true {
		os.Remove(compress.CopyFilelocation)
	}
}

func (compress *Compress) getFilePath() (string) {
	return compress.FilePath + string(os.PathSeparator) + compress.FileName + compress.fileExtention
}

func (compress *Compress) getZipFile(file io.Writer) (gzipFile *gzip.Writer){
	gzipFile = gzip.NewWriter(file)
	gzipFile.Name = compress.FileName
	gzipFile.ModTime = time.Date(1977, time.May, 25, 0, 0, 0, 0, time.UTC)
	return 
}

func (compress *Compress) copyData() {

	file, err := os.OpenFile(
		compress.getFilePath(), 
		os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 
		0644)
	compress.helper.LogError(err)

	defer file.Close()
	gzipFile := compress.getZipFile(file)
	defer gzipFile.Close()

	copyFile, err := os.Open(compress.CopyFilelocation)
	compress.helper.LogError(err)

	// Copy content of CopyFilelocation to Gzip File
	_, err = io.Copy(gzipFile, copyFile)
	compress.helper.LogError(err)
}

// Gzip the file to the Given BACKUP_PATH
	// dump.helper.LogInfo(fmt.Sprintf("Starting GZIP for %s ",dump.connection.TableName))
	// dump.compress.FilePath = string(dump.connection.BackUpPath)
	// dump.compress.FileName = dump.fileName
	// dump.compress.CopyFilelocation = dump.fullPath
	// dump.compress.DeleteOriginalFile = true
	// dump.compress.CompressFile()
	// dump.helper.LogInfo(fmt.Sprintf("Ending GZIP %s ",dump.connection.TableName))
