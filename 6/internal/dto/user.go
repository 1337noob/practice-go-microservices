package dto

type CreateUserDto struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}

type UpdateUserDto struct {
	Name  string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	Age   int    `json:"age,omitempty"`
}

type UserDto struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int    `json:"age"`
}
