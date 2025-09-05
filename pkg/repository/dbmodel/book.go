package dbmodel

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Book struct {
	bun.BaseModel `bun:"table:books"`

	ID                     uuid.UUID `bun:",pk,type:uuid,default:uuid_generate_v4()"`
	CreatedAt              time.Time `bun:",notnull,nullzero,default:now()"`
	UpdatedAt              time.Time `bun:",notnull,nullzero,default:now()"`
	Name                   string    `bun:",notnull,nullzero"`
	OwnerID                uuid.UUID `bun:",type:uuid,notnull,nullzero"`
	DefaultCurrencyIsoCode string    `bun:",notnull,nullzero"`

	Owner *User `bun:"rel:belongs-to,join:owner_id=id"`
}
