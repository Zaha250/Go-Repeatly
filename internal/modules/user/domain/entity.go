package domain

import "time"

type ID int64

type User struct {
	ID         ID
	TelegramID int64
	FirstName  string
	Username   string
	Status     string
	Timezone   string
	CreatedAt  time.Time
}
