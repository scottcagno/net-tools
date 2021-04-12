package user

type User struct {
	ID       int    `json:"id" html:"id"`
	Name     string `json:"name" html:"name"`
	Username string `json:"username" html:"username"`
	Password string `json:"password" html:"password"`
	Role     int    `json:"role" html:"role"`
	Active   bool   `json:"active" html:"active"`
}

func (u *User) GetID() int {
	return u.ID
}
