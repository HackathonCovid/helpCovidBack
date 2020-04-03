package models

import (
	"errors"
	"html"
	"log"
	"strings"
	"time"

	"github.com/HackathonCovid/helpCovidBack/api/security"

	"github.com/badoux/checkmail"
	"github.com/jinzhu/gorm"
	uuid "github.com/twinj/uuid"
	"golang.org/x/crypto/bcrypt"
)

// User Struct
type User struct {
	UUID        uuid.UUID  `gorm:"type:uuid;unique_index;" json:"uuid"`
	ID          uint64     `gorm:"primary_key;auto_increment" json:"id"`
	Firstname   string     `valid:"required,alpha,length(2|255)" json:"firstname"`
	Lastname    string     `valid:"required,alpha,length(2|255)" json:"lastname"`
	Email       string     `gorm:"size:100;not null;unique" valid:"email" json:"email"`
	Password    string     `gorm:"size:100;not null;" json:"password"`
	Isvolunteer int        `valid:"range(0|1),numeric" json:"is_volunteer"`
	TypeOrga    string     `gorm:"size:150;null;" json:"type_orga"`
	City        string     `gorm:"size:150;not null;" json:"city"`
	Adress      string     `gorm:"size:150;not null;" json:"adress"`
	PhoneNumber string     `gorm:"size:15;null" json:"phone_number"`
	OrgaName    string     `gorm:"size:150;null;" json:"orga_name"`
	Description string     `gorm:"text;null;" json:"description"`
	Degree      string     `gorm:"text;null;" json:"degree"`
	Longitude   string     `gorm:"text;null;" json:"longitude"`
	Latitude    string     `gorm:"text;null;" json:"latitude"`
	CreatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Messages    []*Message `gorm:"many2many:message_users;association_foreignkey:id;foreignkey:id" json:"messages,omitempty"`
}

// TableName : Gorm related
func (u *User) TableName() string {
	return "users"
}

// BeforeSave : call package security to hash the password
func (u *User) BeforeSave() error {
	hashedPassword, err := security.Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	u.UUID = uuid.NewV4()
	u.CreatedAt = time.Now()
	return nil
}

// BeforeUpdate is gorm hook that is triggered on every updated on user struct
func (u *User) BeforeUpdate(scope *gorm.Scope) error {
	scope.SetColumn("UpdatedAt", time.Now())
	return nil
}

// Prepare : prepare  statements
func (u *User) Prepare() {
	u.UUID = uuid.NewV4()
	u.ID = 0
	u.Firstname = html.EscapeString(strings.TrimSpace(u.Firstname))
	u.Lastname = html.EscapeString(strings.TrimSpace(u.Lastname))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.City = html.EscapeString(strings.TrimSpace(u.City))
	u.Adress = html.EscapeString(strings.TrimSpace(u.Adress))
	u.PhoneNumber = html.EscapeString(strings.TrimSpace(u.PhoneNumber))
	u.OrgaName = html.EscapeString(strings.TrimSpace(u.OrgaName))
	u.TypeOrga = html.EscapeString(strings.TrimSpace(u.TypeOrga))
	u.Description = html.EscapeString(strings.TrimSpace(u.Description))
	u.Degree = html.EscapeString(strings.TrimSpace(u.Degree))
	u.Longitude = html.EscapeString(strings.TrimSpace(u.Longitude))
	u.Latitude = html.EscapeString(strings.TrimSpace(u.Latitude))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

// VerifyPassword : This method compare the password with the hash
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// Validate : function to check the data
func (u *User) Validate(action string) map[string]string {
	var errorMessages = make(map[string]string)
	var err error

	switch strings.ToLower(action) {
	case "update":
		if u.Email == "" {
			err = errors.New("Required Email")
			errorMessages["Required_email"] = err.Error()
		}
		if u.Email != "" {
			if err = checkmail.ValidateFormat(u.Email); err != nil {
				err = errors.New("Invalid Email")
				errorMessages["Invalid_email"] = err.Error()
			}
		}

	case "login":
		if u.Password == "" {
			err = errors.New("Required Password")
			errorMessages["Required_password"] = err.Error()
		}
		if u.Email == "" {
			err = errors.New("Required Email")
			errorMessages["Required_email"] = err.Error()
		}
		if u.Email != "" {
			if err = checkmail.ValidateFormat(u.Email); err != nil {
				err = errors.New("Invalid Email")
				errorMessages["Invalid_email"] = err.Error()
			}
		}
	case "forgotpassword":
		if u.Email == "" {
			err = errors.New("Required Email")
			errorMessages["Required_email"] = err.Error()
		}
		if u.Email != "" {
			if err = checkmail.ValidateFormat(u.Email); err != nil {
				err = errors.New("Invalid Email")
				errorMessages["Invalid_email"] = err.Error()
			}
		}
	default:
		if u.Firstname == "" {
			err = errors.New("Required Firstname")
			errorMessages["Required_firstname"] = err.Error()
		}
		if u.Lastname == "" {
			err = errors.New("Required Lastname")
			errorMessages["Required_lastname"] = err.Error()
		}
		if u.City == "" {
			err = errors.New("Required City")
			errorMessages["Required_city"] = err.Error()
		}
		if u.Adress == "" {
			err = errors.New("Required Adress")
			errorMessages["Required_adress"] = err.Error()
		}
		if u.Password == "" {
			err = errors.New("Required Password")
			errorMessages["Required_password"] = err.Error()
		}
		if u.Password != "" && len(u.Password) < 6 {
			err = errors.New("Password should be atleast 6 characters")
			errorMessages["Invalid_password"] = err.Error()
		}
		if u.Email == "" {
			err = errors.New("Required Email")
			errorMessages["Required_email"] = err.Error()

		}
		if u.Email != "" {
			if err = checkmail.ValidateFormat(u.Email); err != nil {
				err = errors.New("Invalid Email")
				errorMessages["Invalid_email"] = err.Error()
			}
		}
	}
	return errorMessages

}

// SaveUser : Method Save User, triggered on every saved on user struct
func (u *User) SaveUser(db *gorm.DB) (*User, error) {
	var err error
	err = db.Debug().Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

// FindAllUsers : function to find all users
func (u *User) FindAllUsers(db *gorm.DB) (*[]User, error) {
	var err error
	users := []User{}
	err = db.Debug().Model(&User{}).Limit(100).Find(&users).Error
	if err != nil {
		return &[]User{}, err
	}
	return &users, err
}

// FindUserByID : function to find a user with an ID
func (u *User) FindUserByID(db *gorm.DB, uid uint64) (*User, error) {
	var err error
	err = db.Debug().Model(User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &User{}, errors.New("User Not Found")
	}
	return u, err
}

// UpdateAUser : update an user
func (u *User) UpdateAUser(db *gorm.DB, uid uint64) (*User, error) {
	// To hash the password
	err := u.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}
	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"firstname":    u.Firstname,
			"lastname":     u.Lastname,
			"email":        u.Email,
			"type_orga":    u.TypeOrga,
			"orga_name":    u.OrgaName,
			"city":         u.City,
			"adress":       u.Adress,
			"phone_number": u.PhoneNumber,
			"description":  u.Description,
			"degree":       u.Degree,
			"longitude":    u.Longitude,
			"latitude":     u.Latitude,
			"update_at":    time.Now(),
		},
	)
	if db.Error != nil {
		return &User{}, db.Error
	}
	// This is the display the updated user
	err = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

// DeleteAUser : function to delete the user
func (u *User) DeleteAUser(db *gorm.DB, uid uint64) (int64, error) {
	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).Delete(&User{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}

// UpdatePassword : funtion to update the password
func (u *User) UpdatePassword(db *gorm.DB) error {

	// To hash the password
	err := u.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}
	

	db = db.Debug().Model(&User{}).Where("email = ?", u.Email).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"password":  u.Password,
			"update_at": time.Now(),
		},
	)
	if db.Error != nil {
		return db.Error
	}
	return nil
}
