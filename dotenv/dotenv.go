package dotenv

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"strings"
)

func ParseDotenv() {
	re := regexp.MustCompile(`^.+=.+$`)
	f, err := os.Open(".env")
	if err != nil {
		log.Fatal("err opening .env", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if re.Match([]byte(line)) {
			split := strings.Split(line, "=")
			name := strings.TrimSpace(split[0])
			value := strings.TrimSpace(split[1])
			os.Setenv(name, value)
		}
	}
}
