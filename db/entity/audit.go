package entity

import "time"

// Audit presents the required entities' column mapping for trace auditing.
type Audit struct {
	CreatedBy        uint      `gorm:"column:CreatedBy"`
	CreatedDatim     time.Time `gorm:"column:CreatedDatim;autoCreateTime;<-:create"`
	LastUpdatedBy    uint      `gorm:"column:LastUpdatedBy"`
	LastUpdatedDatim time.Time `gorm:"column:LastUpdatedDatim;autoUpdateTime"`
}

// NewAudit creates the audit on given base entity and current user identifier
func NewAudit(base *Base, userID uint) *Audit {
	auditor := &Audit{
		LastUpdatedBy: userID,
	}
	if base == nil {
		auditor.CreatedBy = userID
	}
	return auditor
}
