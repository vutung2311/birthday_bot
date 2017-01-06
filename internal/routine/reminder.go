package routine

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"birthday-bot/internal/model"
	"birthday-bot/internal/pkg/dingtalk"
	"birthday-bot/internal/repository"

	"github.com/beevik/ntp"
)

const peopleWhoHaveBirthday = "People who have birthday today"

func GetDingTalkReminderRoutine(db *sql.DB, accessToken string) func() {
	client := dingtalk.NewClient(accessToken)
	vnTimeLocation := time.FixedZone("Saigon Time", int((7 * time.Hour).Seconds()))
	notifyChecker := birthdayNotifyChecker{
		period:  24 * time.Hour,
		storage: make(map[model.PersonBirthday]NotifiedBirthday, 10),
	}
	return func() {
		for {
			err := checkAndSendNotification(db, client, vnTimeLocation, &notifyChecker)
			if err != nil {
				log.Printf("ERROR: %v", err.Error())
			}
			time.Sleep(5 * time.Second)
		}
	}
}

func checkAndSendNotification(
	db *sql.DB,
	client *dingtalk.Client,
	location *time.Location,
	notifyChecker *birthdayNotifyChecker,
) error {
	now, err := ntp.Time("asia.pool.ntp.org")
	if err != nil {
		return err
	}

	birthdays, err := repository.GetBirthdayAt(db, now.In(location))
	if err != nil {
		return err
	}
	if len(birthdays) == 0 {
		return nil
	}

	people := make([]string, 0, 10)
	for _, birthday := range birthdays {
		if notifyChecker.Check(birthday) {
			continue
		}
		people = append(people, birthday.PersonName)
	}
	if len(people) == 0 {
		return nil
	}

	text := fmt.Sprintf("# %s\n", peopleWhoHaveBirthday)
	for _, person := range people {
		text += fmt.Sprintf("* **%s**\n", person)
	}
	err = client.SendMarkdown(
		peopleWhoHaveBirthday,
		text,
		nil,
		false,
	)
	if err != nil {
		return err
	}
	for _, birthday := range birthdays {
		notifyChecker.Notified(birthday)
	}
	return nil
}

type NotifiedBirthday struct {
	Notified   bool
	NotifiedAt time.Time
}

type birthdayNotifyChecker struct {
	storage map[model.PersonBirthday]NotifiedBirthday
	period  time.Duration
}

func (c *birthdayNotifyChecker) Notified(b model.PersonBirthday) {
	c.storage[b] = NotifiedBirthday{
		Notified:   true,
		NotifiedAt: time.Now(),
	}
	for k, v := range c.storage {
		if time.Now().Sub(v.NotifiedAt) > c.period {
			delete(c.storage, k)
		}
	}
}

func (c *birthdayNotifyChecker) Check(b model.PersonBirthday) bool {
	notified := c.storage[b]
	if notified.Notified && time.Now().Sub(notified.NotifiedAt) < c.period {
		return true
	}
	return false
}
