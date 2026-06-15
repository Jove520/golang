package vo

type User struct {
	BaseVo
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Power    int    `json:"power"`
	Status   bool   `json:"status"`
}
