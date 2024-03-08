package meta

import "time"

type ObjectMeta struct {
	ID        uint64    `json:"id,omitempty" gorm:"primary_key;AUTO_INCREMENT;column:id"`
	CreatedAt time.Time `json:"created_at,omitempty" gorm:"column:created_at"`
	CreatedBy uint64    `json:"created_by,omitempty" gorm:"column:created_by"`
	UpdatedAt time.Time `json:"updated_at,omitempty" gorm:"column:updated_at"`
}
