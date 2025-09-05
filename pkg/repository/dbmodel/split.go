package dbmodel

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Split struct {
	bun.BaseModel `bun:"table:splits"`

	ID                    uuid.UUID `bun:",pk,type:uuid,default:uuid_generate_v4()"`
	CreatedAt             time.Time `bun:",notnull,nullzero,default:now()"`
	UpdatedAt             time.Time `bun:",notnull,nullzero,default:now()"`
	ExchangeID            uuid.UUID `bun:",type:uuid,notnull,nullzero"`
	DestinationRegisterId uuid.UUID `bun:",type:uuid,notnull,nullzero"`
	Amount                int64     `bun:",notnull"`
	CounterpartAmount     *int64    `bun:","`
	Memo                  *string   `bun:",nullzero"`
	Status                string    `bun:",notnull,nullzero,default:uncleared"`

	Exchange *Exchange `bun:",rel:belongs-to,join=exchange_id=id"`
}
