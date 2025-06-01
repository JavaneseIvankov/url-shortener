package model

import (
	"time"

	"github.com/google/uuid"
)

type ShortLink struct {
	 Id uuid.UUID
    ShortName  string
    OriginalUrl string
    CreatedAt  time.Time
    UserId     uuid.UUID
}
