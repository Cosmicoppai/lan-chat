package admin

import (
	"bufio"
	"errors"
	"log"
	"os"
	"strings"
)

func LoadEnv(envFile string) {
	if envFile != "" {
		file, err := os.Open(envFile)
		if err != nil {
			log.Fatalln(err)
		}
		defer func(file *os.File) {
			_ = file.Close()
		}(file)
		scanner := bufio.NewScanner(file)
		index := -1
		var currLine string
		for scanner.Scan() {
			currLine = scanner.Text()
			if index = strings.Index(currLine, "="); index != -1 {
				err := os.Setenv(currLine[0:index], currLine[index+1:])
				if err != nil {
					log.Fatalln(err)
				}
			}
		}
	} else {
		_ = errors.New("unable to locate the env File")
	}
}
