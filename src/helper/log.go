package helper

import (
	"log"
)

func (helper *Helper) LogError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (helper *Helper) LogInfo(text string) {
	log.Println(text)
}