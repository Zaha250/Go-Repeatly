package user

type Id int64
type User struct {
	Id        Id
	TgId      int64
	Timezone  string
	Name      string
	UserName  string
	CreatedAt string
}
