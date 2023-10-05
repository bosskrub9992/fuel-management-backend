package models

import "time"

type GetHealthResponse struct {
	ServerStartTime time.Time `json:"serverStartTime"`
}
