package logger

import (
	"bufio"
	"log"
	"os"

	"github.com/perdokcat/TermoTune/config"
)

var (
	INFO_Logger *log.Logger
	WARN_Logger *log.Logger
	ERROR_Logger *log.Logger
)

func initLogger() {
	log_file_path, err := os.OpenFile(
		config.GetConfig().TermoTunePath,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 
		0666)
	// if the log file path is not valid, using stdout for logging	
	if err != nil {
		log_file_path = os.Stdout
	}		
	INFO_Logger = log.New(
		log_file_path,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile,
	)
	WARN_Logger = log.New(
		log_file_path,
		"WARN: ",
		log.Ldate|log.Ltime|log.Lshortfile,
	)
	ERROR_Logger = log.New(
		log_file_path,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile,
	)
}

func  LogError(msgErr TermoTuneError, extra ...any){ 
	ERROR_Logger.Println(
		msgErr,
		extra,
	)
}

func LogInfor(msg string, extra ...any) {
	INFO_Logger.Println(
		msg,
		extra,
	)
}

func LogWarn(warn string, extra ...any) {
	WARN_Logger.Println(warn, extra)
}


func GetLog() ([]string, error) {
	log_file_path, err := os.Open(config.GetConfig().LogFile)

	if err != nil  {
		return nil, err
	}

	var LastLines []string 
	scanner := bufio.NewScanner(log_file_path)
	for scanner.Scan() {
		LastLines = append(LastLines, scanner.Text())
		if(len(LastLines) > 200) {
			LastLines = LastLines[1:]
		}
	}

	return LastLines, nil
}

