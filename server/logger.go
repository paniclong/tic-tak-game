package server

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

type Logger struct {
	file   *os.File
	prefix string
}

type ExtendData struct {
	Cell map[string]string
}

func CreateLogger() (error, *Logger) {
	var dirPath = "/app/log/"
	if d := os.Getenv("DIR_NAME_LOG"); d != "" {
		dirPath = d
	}

	_, err := os.Stat(dirPath)
	if err != nil {
		err := os.Mkdir(dirPath, 0777)
		if err != nil {
			return err, nil
		}
	}

	file, err := os.OpenFile(dirPath+"info.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		return err, nil
	}

	if file == nil {
		return err, nil
	}

	logger := new(Logger)
	logger.file = file

	return nil, logger
}

func (log *Logger) SetPrefix(message string) {
	log.prefix = message
}

func (log *Logger) Write(message string) {
	log.writeCustom(message, &ExtendData{})
}

func (log *Logger) WriteWithData(message string, data *ExtendData) {
	log.writeCustom(message, data)
}

func (log *Logger) writeCustom(message string, data *ExtendData) {
	var writingMessage string
	writingMessage = "[" + time.Now().Format("02.01.2006-15:04:05") + "]"

	if log.prefix != "" {
		writingMessage += fmt.Sprintf("[%s] ", log.prefix)
	}

	writingMessage += message

	if len(data.Cell) != 0 {
		cell, _ := json.Marshal(data)
		writingMessage += string(cell)
	}

	writingMessage += "\n"

	_, err := log.file.WriteString(writingMessage)
	if err != nil {
		fmt.Println(err)
	}
}
