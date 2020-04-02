package models

import (
	"errors"
	"html"
	"strings"
	"time"

	"github.com/jinzhu/gorm"
	uuid "github.com/twinj/uuid"
)

// Mission Struct
type Mission struct {
	UUID             uuid.UUID `gorm:"type:uuid;unique_index;" json:"uuid"`
	ID               uint64    `gorm:"primary_key;auto_increment" json:"id"`
	Title            string    `gorm:"size:255;not null;" json:"title"`
	Description      string    `gorm:"text;not null;" json:"description"`
	StartDate        time.Time `gorm:"not null;" json:"start_date"`
	EndDate          time.Time `gorm:"null;" json:"end_date"`
	NbDays           int       `gorm:"null;" json:"nb_days"`
	NbPeopleRequired int       `gorm:"null;" json:"nb_people_required"`
	SkillsRequired   string    `gorm:"text;not null;" json:"skills_required"`
	NightOrDay       string    `gorm:"size:150;not null;" json:"night_or_day"`
	AddressHospital  string    `gorm:"size:255;not null;" json:"address_hospital"`
	Author           User      `json:"author"`
	AuthorID         uint64    `gorm:"not null" json:"author_id"`
	CreatedAt        time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt        time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Comments []*Comment `gorm:"many2many:mission_comments;association_foreignkey:id;foreignkey:id" json:"comments,omitempty"`
	Users []*User `gorm:"many2many:mission_users;association_foreignkey:id;foreignkey:id" json:"users,omitempty"`
}

// TableName : Gorm related
func (m *Mission) TableName() string {
	return "missions"
}

// BeforeSave : Method before Save
func (m *Mission) BeforeSave(scope *gorm.Scope) error {
	scope.SetColumn("UUID", uuid.NewV4())
	scope.SetColumn("CreatedAt", time.Now())
	return nil
}

// BeforeUpdate is gorm hook that is triggered on every updated on vote struct
func (m *Mission) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("UpdatedAt", time.Now())
	return nil
}

// Validate : function to check the data
func (m *Mission) Validate() map[string]string {
	var err error
	var errorMessages = make(map[string]string)

	if m.Title == "" {
		err = errors.New("Required Title")
		errorMessages["Required_title"] = err.Error()
	}
	if m.Description == "" {
		err = errors.New("Required Description")
		errorMessages["Required_description"] = err.Error()
	}
	if m.StartDate.IsZero() {
		err = errors.New("Required Start Date")
		errorMessages["Required_start_date"] = err.Error()
	}
	if m.SkillsRequired == "" {
		err = errors.New("Required Skills")
		errorMessages["Required_skills"] = err.Error()
	}
	if m.NightOrDay == "" {
		err = errors.New("Required Night or Day")
		errorMessages["Required_night_or_day"] = err.Error()
	}
	if m.AddressHospital == "" {
		err = errors.New("Required Address")
		errorMessages["Required_address"] = err.Error()
	}
	if m.AuthorID < 1 {
		err = errors.New("Required Author")
		errorMessages["Required_author"] = err.Error()
	}

	return errorMessages
}

//Prepare : prepare a mission
func (m *Mission) Prepare() {
	m.Title = html.EscapeString(strings.TrimSpace(m.Title))
	m.Description = html.EscapeString(strings.TrimSpace(m.Description))
	m.NightOrDay = html.EscapeString(strings.TrimSpace(m.NightOrDay))
	m.SkillsRequired = html.EscapeString(strings.TrimSpace(m.SkillsRequired))
	m.AddressHospital = html.EscapeString(strings.TrimSpace(m.AddressHospital))
	m.Author = User{}
	m.CreatedAt = time.Now()
	m.UpdatedAt = time.Now()
}

// SaveMission : Method Save Mission, triggered on every saved on trip struct
func (m *Mission) SaveMission(db *gorm.DB) (*Mission, error) {
	var err error
	err = db.Debug().Model(&Mission{}).Create(&m).Error
	if err != nil {
		return &Mission{}, err
	}
	if m.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", m.AuthorID).Take(&m.Author).Error
		if err != nil {
			return &Mission{}, err
		}
	}
	return m, nil
}

// FindAllMissions : function to find all missions
func (m *Mission) FindAllMissions(db *gorm.DB) (*[]Mission, error) {
	var err error
	missions := []Mission{}
	err = db.Debug().Model(&Mission{}).Limit(100).Preload("Users").Order("created_at desc").Find(&missions).Error
	if err != nil {  
		return &[]Mission{}, err
	}
	if len(missions) > 0 {
		for i := range missions {
			err := db.Debug().Model(&User{}).Where("id = ?", missions[i].AuthorID).Take(&missions[i].Author).Error
			if err != nil {
				return &[]Mission{}, err
			}
		}
	}
	return &missions, nil
}

// FindMissionByID : function to find a mission with an ID
func (m *Mission) FindMissionByID(db *gorm.DB, pid uint64) (*Mission, error) {
	var err error
	err = db.Debug().Model(&Mission{}).Where("id = ?", pid).Preload("Users").Preload("Comments").Take(&m).Error
	if err != nil {
		return &Mission{}, err
	}
	if m.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", m.AuthorID).Take(&m.Author).Error
		if err != nil {
			return &Mission{}, err
		}
	}
	return m, nil
}

// UpdateAMission : function to update a mission
func (m *Mission) UpdateAMission(db *gorm.DB, pid uint64) (*Mission, error) {
	var err error
	db = db.Debug().Model(&Mission{}).Where("id = ?", pid).Take(&Mission{}).UpdateColumns(
		map[string]interface{}{
			"title":              m.Title,
			"description":        m.Description,
			"start_date":         m.StartDate,
			"end_date":           m.EndDate,
			"nb_days":            m.NbDays,
			"nb_people_required": m.NbPeopleRequired,
			"skills_required":    m.SkillsRequired,
			"night_or_day":       m.NightOrDay,
			"address_hospital":   m.AddressHospital,
			"updated_at":         time.Now(),
		},
	)
	err = db.Debug().Model(&Mission{}).Where("id = ?", pid).Take(&m).Error
	if err != nil {
		return &Mission{}, err
	}
	if m.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", m.AuthorID).Take(&m.Author).Error
		if err != nil {
			return &Mission{}, err
		}
	}
	return m, nil
}

// DeleteAMission : function to delete a mission
func (m *Mission) DeleteAMission(db *gorm.DB, pid uint64, uid uint64) (int64, error) {
	db = db.Debug().Model(&Mission{}).Where("id = ? and author_id = ?", pid, uid).Take(&Mission{}).Delete(&Mission{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Mission not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

// FindUserMissions : function to get all missions for a user
func (m *Mission) FindUserMissions(db *gorm.DB, uid uint64) (*[]Mission, error) {

	var err error
	missions := []Mission{}
	err = db.Debug().Model(&Mission{}).Where("author_id = ?", uid).Limit(100).Order("created_at desc").Find(&missions).Error
	if err != nil {
		return &[]Mission{}, err
	}
	if len(missions) > 0 {
		for i := range missions {
			err := db.Debug().Model(&User{}).Where("id = ?", missions[i].AuthorID).Take(&missions[i].Author).Error
			if err != nil {
				return &[]Mission{}, err
			}
		}
	}
	return &missions, nil
}

// DeleteUserMissions : When a user is deleted, we also delete the missions that the user had
func (m *Mission) DeleteUserMissions(db *gorm.DB, uid uint64) (int64, error) {
	missions := []Mission{}
	db = db.Debug().Model(&Mission{}).Where("author_id = ?", uid).Find(&missions).Delete(&missions)
	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
