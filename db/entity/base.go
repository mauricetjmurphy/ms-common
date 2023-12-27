package entity

// Base presents the basic entities' column mapping
type Base struct {
	// ID describes the primary key column
	ID uint `gorm:"column:Id;primaryKey;autoIncrement:true"`
}

// NewBase creates the base entity
func NewBase(id uint) *Base {
	return &Base{ID: id}
}

func NewBaseInt64(value int64) *Base {
	if value > 0 {
		return NewBase(uint(value))
	}
	return nil
}

func BaseNotNil(e *Base) bool {
	return e != nil && e.ID > 0
}

func BaseNil(e *Base) bool {
	return !BaseNotNil(e)
}
