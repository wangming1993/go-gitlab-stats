package lib

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

func PanicIfError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func WriteFile(out interface{}, fileName string) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}

	buf, err := json.MarshalIndent(out, "", "\t")
	if err != nil {
		return err
	}
	_, err = file.Write(buf)
	return err
}

func WriteJSON(out interface{}, name string) error {
	return WriteFile(out, "output/"+name+".json")
}

func StatsStartTime() time.Time {
	return time.Date(SettinsService.StatsYear(), time.January, 1, 0, 0, 0, 0, time.UTC)
}

func StatsEndTime() time.Time {
	return time.Date(SettinsService.StatsYear(), time.December, 31, 23, 59, 59, 0, time.UTC)
}
