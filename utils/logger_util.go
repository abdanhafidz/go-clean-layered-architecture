package utils

import (
	"fmt"
	"log"
	"os"
)

func InternalErrorLog(err_log error) {
	fmt.Println("There is an error!")

	file, err := os.OpenFile("logs/error_log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	if err != nil {
		log.Fatal(err)
	}
	log.Println("Error Log :", err_log)
	log.SetOutput(file)
}

func SecurityLog(security_log string) {
	fmt.Println("There is an error!")

	file, err := os.OpenFile("logs/security_log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)

	if err != nil {
		log.Fatal(err)
	}
	log.Println("Security Log :", security_log)
	log.SetOutput(file)
}
