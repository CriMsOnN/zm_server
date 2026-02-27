package dto

type CreateOrUpdateUserDTO struct {
	Name    string `json:"name"`
	Fivem   string `json:"fivem"`
	License string `json:"license"`
}

type UserJoinedDTO struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Identifier string `json:"identifier"`
}
