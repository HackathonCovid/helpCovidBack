package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
)

// Message struct
type Message struct {
	ID        uint64    `gorm:"primary_key;auto_increment" json:"id"`
	TransmitterID uint64   `gorm:"not null" json:"user_id"`
	RecipientID uint64    `gorm:"not null" json:"mission_id"`
	Body      string    `gorm:"text;not null;" json:"body"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// Prepare : prepare statements
func (c *Message) Prepare() {
	c.ID = 0
	c.Body = html.EscapeString(strings.TrimSpace(c.Body))
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
}

// Validate : validation rules
func (c *Message) Validate(action string) map[string]string {
	var errorMessages = make(map[string]string)
	var err error

	switch strings.ToLower(action) {
	case "update":
		if c.Body == "" {
			err = errors.New("Required Body in a Message")
			errorMessages["Required_body"] = err.Error()
		}
	default:
		if c.Body == "" {
			err = errors.New("Required body to Message")
			errorMessages["Required_body"] = err.Error()
		}
	}
	return errorMessages
}

// SaveMessage : function to save a Message linked to a user
func (c *Message) SaveMessage(db *gorm.DB, mission *Mission) (*Message, error) {

//	err :=db.Model(&mission).Association("Messages").Append(&c).Error;
	//server.DB.Model(&mission).Association("Messages").Append(&Message);
	err := db.Debug().Create(&c).Error
	if err != nil {
		return &Message{}, err
	}
	return c, nil
}

// GetMessages : function to get all the Messages for a mission and a user
func (c *Message) GetMessages(db *gorm.DB, pid uint64) (*[]Message, error) {

	Messages := []Message{}
	err := db.Debug().Model(&Message{}).Where("mission_id = ?", pid).Order("created_at desc").Find(&Messages).Error
	if err != nil {
		return &[]Message{}, err
	}
	return &Messages, err
}

// UpdateAMessage : funtion to update a Message
func (c *Message) UpdateAMessage(db *gorm.DB) (*Message, error) {

	var err error
	err = db.Debug().Model(&Message{}).Where("id = ?", c.ID).Updates(Message{Body: c.Body, UpdatedAt: time.Now()}).Error
	if err != nil {
		return &Message{}, err
	}
	return c, nil
}

// DeleteAMessage : function to delete a Message
func (c *Message) DeleteAMessage(db *gorm.DB) (int64, error) {

	db = db.Debug().Model(&Message{}).Where("id = ?", c.ID).Take(&Message{}).Delete(&Message{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

// DeleteUserMessages : When a user is deleted, we also delete the Messages that the user had
func (c *Message) DeleteUserMessages(db *gorm.DB, uid uint64) (int64, error) {
	Messages := []Message{}
	db = db.Debug().Model(&Message{}).Where("user_id = ?", uid).Find(&Messages).Delete(&Messages)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

// DeleteMissionMessages : When a mission is deleted, we also delete the Messages that the mission had
func (c *Message) DeleteMissionMessages(db *gorm.DB, pid uint64) (int64, error) {
	Messages := []Message{}
	db = db.Debug().Model(&Message{}).Where("mission_id = ?", pid).Find(&Messages).Delete(&Messages)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
