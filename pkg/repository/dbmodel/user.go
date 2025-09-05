package dbmodel

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users"`

	ID                uuid.UUID `bun:",pk,type:uuid,default:uuid_generate_v4()"`
	CreatedAt         time.Time `bun:",notnull,nullzero,default:current_timestamp"`
	UpdatedAt         time.Time `bun:",notnull,nullzero,default:current_timestamp"`
	Email             string    `bun:",notnull,nullzero,unique"`
	EncryptedPassword string    `bun:",notnull,nullzero"`
	Admin             bool      `bun:",notnull,default=false"`
	DisplayName       string    `bun:",notnull,nullzero"`
}
