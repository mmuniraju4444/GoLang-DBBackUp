package env

import (
	"bufio"
	"helper"
	"os"
	"strings"
	"errors"
)

type Env struct {
	fileName    string ".env"
	fileHandler *os.File
	envData     map[string]string
	helper      *helper.Helper
}

func (env *Env) SetEnvFileName(filename string) *Env {
	env.fileName = filename
	return env
}

func (env *Env) init() *Env {
	env.envData = make(map[string]string)
	env.helper = &helper.Helper{}
	return env
}

func (env *Env) loadFile() *Env {
	envFile, err := os.Open(env.fileName)
	env.helper.LogError(err)
	env.fileHandler = envFile
	return env
}

func (env *Env) closeFile() *Env {
	defer func() {
		env.fileHandler = nil
	}()
	env.fileHandler.Close()
	return env
}

func (env *Env) Load() *Env {
	env.init()
	env.loadFile()
	env.setEnvData()
	env.closeFile()
	return env
}

func (env *Env) Getenv(key string) string{
	return os.Getenv(key)
}

func (env *Env) getEnvData() *Env{
	scanner := bufio.NewScanner(env.fileHandler)
	var fileLines []string
	for scanner.Scan() {
		fileLines = append(fileLines, scanner.Text())
	}
	for _, line := range fileLines {
		key, value, err := env.getKeyAndValue(line)
		env.helper.LogError(err)
		env.envData[key] = strings.TrimSpace(value)
	}
	return env
}
func (env *Env) setEnvData() {
	env.getEnvData()
	for key, value := range env.envData {
		os.Setenv(key, value)
	}
}

func (env *Env) getKeyAndValue(line string) (key, value string, err error) {

	lineSplit := strings.SplitN(line, "=", 2)
	if len(lineSplit) != 2 {
		err = errors.New("Should be in format `KEY`=`VALUE` but found : " + strings.Join(lineSplit, " "))
		return
	}
	key = lineSplit[0]
	value = lineSplit[1]
	return
}
