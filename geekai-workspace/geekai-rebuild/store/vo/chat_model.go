package vo

type ChatModel struct {
	BaseVo
	Name       string `json:"name"`
	Value      string `json:"value"`
	Power      int    `json:"power"`
	MaxTokens  int    `json:"max_tokens"`
	MaxContext int    `json:"max_context"`
	Enabled    bool   `json:"enabled"`
}
