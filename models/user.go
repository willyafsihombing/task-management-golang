package models

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type User struct {
	Id        int       `gorm:"type:int;primaryKey;autoIncrement" json:"id"`
	Role      string    `gorm:"type:varchar(10);" json:"role"`
	Name      string    `gorm:"type:varchar(255);" json:"name"`
	Email     string    `gorm:"type:varchar(50);" json:"email"`
	Password  string    `gorm:"type:varchar(255);" json:"-"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Token     string    `gorm:"-" json:"token,omitempty"`
	// Task      []Task    `gorm:"constraint:OnDelete:CASCADE" json:"tasks, omitempty"` //has many
}

func (u *User) AfterDelete(tx *gorm.DB) (err error) {
	tx.Clauses(clause.Returning{}).Where("userId = ?", u.Id).Delete(&Task{})
	return
}
