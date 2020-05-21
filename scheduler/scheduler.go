package scheduler

import (
	"github.com/doubtnut/handler"
	"github.com/doubtnut/logging"
	"github.com/doubtnut/redis"
	redis1 "github.com/go-redis/redis"
	"github.com/jasonlvhit/gocron"
	"strings"
	"time"
)

var Logger = logging.NewLogger()

func Init() {
	schedule := gocron.NewScheduler()
	schedule.Every(1).Seconds().Do(Task)
	<-schedule.Start()

}

func Task() {
	//get keys
	Logger.Infof("scheduler triggered")
	var keys *redis1.StringSliceCmd
	keys = redis.Keys("user*")
	for i := range keys.Val() {
		key := keys.Val()[i]
		value, err := redis.GetValue(key)
		if err != nil {
			Logger.Errorf("Error getting key " + err.Error())
		}
		dateFormat := "2006-01-02 15:04:05.999999999 -0700 MST"
		t1, err := time.Parse(dateFormat, strings.Split(value, " m=")[0])
		if err != nil {
			Logger.Errorf("error splitting value " + err.Error())
		}
		//parse current time
		t2, err := time.Parse(dateFormat, strings.Split(time.Now().String(), " m=")[0])
		if err != nil {
			Logger.Errorf("error parsing timestamp " + err.Error())

		}
		duration := t2.Sub(t1)
		if (duration > 15) {
			//send email
			userId := strings.Split(key, "-")[1]
			err = handler.SendEmail("frommail@gmail.com", "tomail@gmail.com", userId)
			if err != nil {
				Logger.Errorf("error sending mail " + err.Error())
			}
			redis.Del(key)
		}

	}

}
