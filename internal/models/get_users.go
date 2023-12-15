package models

type GetUserDatum struct {
	ID              int64  `json:"id"`
	DefaultCarID    int64  `json:"defaultCarId"`
	Nickname        string `json:"nickname"`
	ProfileImageURL string `json:"profileImageUrl"`
}

type GetUserData struct {
	Data []GetUserDatum `json:"data"`
}
