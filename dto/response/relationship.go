package response

import "time"

type Relationship struct {
	Key       string          `json:"key"`
	Status    string          `json:"status"`
	Actor     ProfileResponse `json:"actor"`
	Partner   ProfileResponse `json:"partner"`
	CreatedAt time.Time       `json:"created_at"`
	UpdatedAt time.Time       `json:"updated_at"`
}
