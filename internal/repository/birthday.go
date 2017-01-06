package repository

import (
	"database/sql"
	"time"

	"birthday-bot/internal/model"
)

func GetAllBirthdays(db *sql.DB) ([]model.PersonBirthday, error) {
	rows, err := db.Query(`SELECT * FROM birthdays`)
	if err != nil {
		return nil, err
	}

	var (
		birthdays []model.PersonBirthday
		birthday  model.PersonBirthday
	)

	for rows.Next() {
		err = rows.Scan(&birthday.Id, &birthday.PersonName, &birthday.Birthday)
		if err != nil {
			return nil, err
		}
		birthdays = append(birthdays, birthday)
	}

	return birthdays, nil
}

func GetBirthdayAt(db *sql.DB, now time.Time) ([]model.PersonBirthday, error) {
	rows, err := db.Query(`
		SELECT * FROM birthdays
		WHERE strftime('%d', birthday) = strftime('%d', ?)
		AND strftime('%m', birthday) = strftime('%m', ?)`, now, now)
	if err != nil {
		return nil, err
	}

	var (
		birthdays []model.PersonBirthday
		birthday  model.PersonBirthday
	)

	for rows.Next() {
		err = rows.Scan(&birthday.Id, &birthday.PersonName, &birthday.Birthday)
		if err != nil {
			return nil, err
		}
		birthdays = append(birthdays, birthday)
	}

	return birthdays, nil
}

func SaveBirthday(db *sql.DB, birthday *model.PersonBirthday) error {
	stmt, err := db.Prepare(`
		INSERT INTO birthdays(person_name, birthday)
		VALUES(?, ?)
	`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(birthday.PersonName, birthday.Birthday.ToTime())
	return err
}

func UpdateBirthday(db *sql.DB, bd *model.PersonBirthday) error {
	stmt, err := db.Prepare(`
		UPDATE birthdays SET person_name=?, birthday=?
		WHERE id=?
	`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(bd.PersonName, bd.Birthday.ToTime(), bd.Id)
	return err
}

func DeleteBirthday(db *sql.DB, id int64) error {
	stmt, err := db.Prepare("DELETE FROM birthdays WHERE id=?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id)
	return err
}

func DeleteAllBirthdays(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM birthdays")
	return err
}
