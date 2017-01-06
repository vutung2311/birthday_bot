package birthday

import (
	"time"
)

const (
	TimeFormat = "2006-01-02"
)

func ParseBirthday(value string) (Birthday, error) {
	t, err := time.Parse(TimeFormat, value)
	return Birthday(t), err
}

type Birthday time.Time

func (b Birthday) MarshalJSON() ([]byte, error) {
	s := b.ToTime().Format(`"` + TimeFormat + `"`)
	return []byte(s), nil
}

func (b *Birthday) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}

	var err error
	*b, err = parseJson(string(data))
	return err
}

func parseJson(value string) (Birthday, error) {
	t, err := time.Parse(`"`+TimeFormat+`"`, value)
	return Birthday(t), err
}

func (b Birthday) ToTime() time.Time {
	return time.Time(b)
}

func (b Birthday) Format(layout string) string {
	return b.ToTime().Format(layout)
}
