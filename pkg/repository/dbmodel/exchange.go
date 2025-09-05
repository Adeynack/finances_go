package dbmodel

import (
	"time"

	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

type Exchange struct {
	bun.BaseModel `bun:"table:exchanges"`

	ID          uuid.UUID `bun:",pk,type:uuid,default:uuid_generate_v4()"`
	CreatedAt   time.Time `bun:",notnull,nullzero,default:now()"`
	UpdatedAt   time.Time `bun:",notnull,nullzero,default:now()"`
	Date        time.Time `bun:",type=date,notnull,nullzero"`
	RegisterID  uuid.UUID `bun:",type=uuid,notnull,nullzero"`
	Cheque      *string   `bun:",nullzero"`
	Description string    `bun:",notnull,nullzero"`
	Memo        *string   `bun:",nullzero"`
	Status      string    `bun:",notnull,nullzero,default:uncleared"`

	Splits []Split `bun:",rel:has-many,join:id=exchange_id"`
}
