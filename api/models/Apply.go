package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

// Apply struct
type Apply struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	UserID    uint64    `gorm:"not null" json:"user_id"`
	MissionID uint64    `gorm:"not null" json:"mission_id"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Validate  int       `gorm:"default:0" json:"validate"`
	Mission   Mission   `json:"mission"`
	User      User      `json:"user"`
}

// SaveApply : function to apply to a mission
func (a *Apply) SaveApply(db *gorm.DB) (*Apply, error) {

	// Check if the auth user has applied to this mission before
	err := db.Debug().Model(&Apply{}).Where("mission_id = ? AND user_id = ?", a.MissionID, a.UserID).Take(&a).Error
	if err != nil {
		if err.Error() == "record not found" {
			// The user has not applied to this mission before, so lets save incomming apply
			err = db.Debug().Model(&Apply{}).Create(&a).Error
			if err != nil {
				return &Apply{}, err
			}
		}
	} else {
		// The user has applied it before, so create a custom error message
		err = errors.New("You already applied to this mission")
		return &Apply{}, err
	}
	return a, nil
}

// DeleteApply : function to delete an apply of a mission
func (a *Apply) DeleteApply(db *gorm.DB) (*Apply, error) {
	var err error
	var deletedApply *Apply

	err = db.Debug().Model(Apply{}).Where("id = ?", a.ID).Take(&a).Error
	if err != nil {
		return &Apply{}, err
	} else {
		//If the apply exist, save it in deleted apply and delete it
		deletedApply = a
		db = db.Debug().Model(&Apply{}).Where("id = ?", a.ID).Take(&Apply{}).Delete(&Apply{})
		if db.Error != nil {
			fmt.Println("cant delete apply: ", db.Error)
			return &Apply{}, db.Error
		}
	}
	return deletedApply, nil
}

// GetAppliesInfo : get the infos
func (a *Apply) GetAppliesInfo(db *gorm.DB, pid uint64) (*[]Apply, error) {

	applies := []Apply{}
	err := db.Debug().Model(&Apply{}).Where("mission_id = ?", pid).Find(&applies).Error
	if err != nil {
		return &[]Apply{}, err
	}
	if len(applies) > 0 {
		for i := range applies {
			err := db.Debug().Model(&User{}).Where("id = ?", applies[i].UserID).Take(&applies[i].User).Error
			if err != nil {
				return &[]Apply{}, err
			}
		}
	}
	return &applies, err
}

// GetAppliesByUserId : get the infos
func (a *Apply) GetAppliesByUserId(db *gorm.DB, pid uint64) (*[]Apply, error) {

	applies := []Apply{}
	err := db.Debug().Model(&Apply{}).Where("user_id = ?", pid).Find(&applies).Error
	if err != nil {
		return &[]Apply{}, err
	}
	if len(applies) > 0 {
		for i := range applies {
			err := db.Debug().Model(&Mission{}).Where("id = ?", applies[i].MissionID).Take(&applies[i].Mission).Error
			if err != nil {
				return &[]Apply{}, err
			}
		}
	}
	return &applies, err
}

// DeleteUserApplies : When a user is deleted, we also delete the applies that the user had
func (a *Apply) DeleteUserApplies(db *gorm.DB, uid uint64) (int64, error) {
	applies := []Apply{}
	db = db.Debug().Model(&Apply{}).Where("user_id = ?", uid).Find(&applies).Delete(&applies)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

// DeleteMissionApplies : When a mission is deleted, we also delete the applies that the mission had
func (a *Apply) DeleteMissionApplies(db *gorm.DB, pid uint64) (int64, error) {
	applies := []Apply{}
	db = db.Debug().Model(&Apply{}).Where("mission_id = ?", pid).Find(&applies).Delete(&applies)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

// FindApplyAndUpdateByID : When an appply is updated, we also delete the applies that the mission had
func (a *Apply) FindApplyAndUpdateByID(db *gorm.DB) (*Apply, error) {
	var err error

	err = db.Debug().Model(Apply{}).Where("id = ?", a.ID).Updates(Apply{Validate: a.Validate, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Apply{}, err
	}

	return a, nil
}
