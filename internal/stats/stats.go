package stats

import (
	"elisBot/internal/config"
	"github.com/robfig/cron/v3"
	"log"
	"time"
)

var messagesCount [24][60]int

func Inc(hour, minute int) {
	messagesCount[hour][minute]++
}

func GetCount(period config.Period, hour, minute int) int {
	switch period {
	case config.MINUTE:
		return getMinuteCount(hour, minute)
	case config.HOUR:
		return getHourCount(hour)
	case config.DAY:
		return getDayCount()
	default:
		return 0
	}
}

func getMinuteCount(hour int, minute int) int {
	return messagesCount[hour][minute]
}

func getHourCount(hour int) int {
	sum := 0
	for i := 0; i < len(messagesCount[hour]); i++ {
		sum += messagesCount[hour][i]
	}
	return sum
}

func getDayCount() int {
	sum := 0
	for i := 0; i < len(messagesCount); i++ {
		sum += getHourCount(i)
	}
	return sum
}

func clear() {
	for i := 0; i < len(messagesCount); i++ {
		for j := 0; j < len(messagesCount[i]); j++ {
			messagesCount[i][j] = 0
		}
	}
}

func ScheduleClear() {
	mow, err := time.LoadLocation(config.StatsClearTimezone)
	if err != nil {
		log.Panic(err)
	}
	c := cron.New(cron.WithLocation(mow))
	_, err = c.AddFunc(config.StatsClearCron, clear)
	if err != nil {
		log.Panic(err)
	}

	c.Run()
}
