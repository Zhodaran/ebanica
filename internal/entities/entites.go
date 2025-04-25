package entities

type UserAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type User struct {
	ID        int          `json:"id"`
	Name      string       `json:"name"`
	Email     string       `json:"email"`
	DeletedAt *string      `json:"deleted_at"` // Для логического удаления
	Books     map[int]Book `json:"books"`
}

type Book struct {
	Index     int    `json:"index"`
	Book      string `json:"book"`
	Author    string `json:"author"`
	Block     *bool  `json:"block"`
	TakeCount int    `json:"take_count"`
}
