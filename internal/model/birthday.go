package model

import "birthday-bot/internal/pkg/birthday"

type PersonBirthday struct {
	Id         int64             `json:"-"`
	PersonName string            `json:"personName"`
	Birthday   birthday.Birthday `json:"birthday"`
}
