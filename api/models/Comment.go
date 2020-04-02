package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

// Comment struct
type Comment struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	UserID    uint64    `gorm:"not null" json:"user_id"`
	MissionID uint64    `gorm:"not null" json:"mission_id"`
	Body      string    `gorm:"text;not null;" json:"body"`
	User      User      `json:"user"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Prepare : prepare statements
func (c *Comment) Prepare() {
	c.ID = 0
	c.Body = html.EscapeString(strings.TrimSpace(c.Body))
	c.User = User{}
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
}

// Validate : validation rules
func (c *Comment) Validate(action string) map[string]string {
	var errorMessages = make(map[string]string)
	var err error

	switch strings.ToLower(action) {
	case "update":
		if c.Body == "" {
			err = errors.New("Required Body in a Comment")
			errorMessages["Required_body"] = err.Error()
		}
	default:
		if c.Body == "" {
			err = errors.New("Required body to Comment")
			errorMessages["Required_body"] = err.Error()
		}
	}
	return errorMessages
}

// SaveComment : function to save a comment linked to a user
func (c *Comment) SaveComment(db *gorm.DB, mission *Mission) (*Comment, error) {

	//	err :=db.Model(&mission).Association("Comments").Append(&c).Error;
	//server.DB.Model(&mission).Association("Comments").Append(&comment);
	err := db.Debug().Create(&c).Error
	if err != nil {
		return &Comment{}, err
	}
	if c.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", c.UserID).Take(&c.User).Error
		if err != nil {
			return &Comment{}, err
		}
	}
	return c, nil
}

// GetComments : function to get all the comments for a mission and a user
func (c *Comment) GetComments(db *gorm.DB, pid uint64) (*[]Comment, error) {

	comments := []Comment{}
	err := db.Debug().Model(&Comment{}).Where("mission_id = ?", pid).Order("created_at desc").Find(&comments).Error
	if err != nil {
		return &[]Comment{}, err
	}
	if len(comments) > 0 {
		for i := range comments {
			err := db.Debug().Model(&User{}).Where("id = ?", comments[i].UserID).Take(&comments[i].User).Error
			if err != nil {
				return &[]Comment{}, err
			}
		}
	}
	return &comments, err
}

// UpdateAComment : funtion to update a comment
func (c *Comment) UpdateAComment(db *gorm.DB) (*Comment, error) {

	var err error
	err = db.Debug().Model(&Comment{}).Where("id = ?", c.ID).Updates(Comment{Body: c.Body, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Comment{}, err
	}

	if c.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", c.UserID).Take(&c.User).Error
		if err != nil {
			return &Comment{}, err
		}
	}
	return c, nil
}

// DeleteAComment : function to delete a comment
func (c *Comment) DeleteAComment(db *gorm.DB) (int64, error) {

	db = db.Debug().Model(&Comment{}).Where("id = ?", c.ID).Take(&Comment{}).Delete(&Comment{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

// DeleteUserComments : When a user is deleted, we also delete the comments that the user had
func (c *Comment) DeleteUserComments(db *gorm.DB, uid uint64) (int64, error) {
	comments := []Comment{}
	db = db.Debug().Model(&Comment{}).Where("user_id = ?", uid).Find(&comments).Delete(&comments)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

// DeleteMissionComments : When a mission is deleted, we also delete the comments that the mission had
func (c *Comment) DeleteMissionComments(db *gorm.DB, pid uint64) (int64, error) {
	comments := []Comment{}
	db = db.Debug().Model(&Comment{}).Where("mission_id = ?", pid).Find(&comments).Delete(&comments)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
