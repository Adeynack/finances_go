package model

// All lists all models that are to be managed by Gorm
func All() []any {
	return []any{
		User{},
	}
}
