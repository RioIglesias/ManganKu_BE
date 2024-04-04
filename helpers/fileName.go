package helpers

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
	"time"
)

func FileName() string {
	randomNum := rand.Intn(10000)
	re := regexp.MustCompile(`[\s-:]`)
	fileName := re.ReplaceAllString(time.Now().Format("2006-01-02 15:04:05"), "") + strconv.Itoa(int(randomNum))

	newImageName := fmt.Sprintf("%s_%d", fileName, time.Now().UnixNano())
	return newImageName
}
