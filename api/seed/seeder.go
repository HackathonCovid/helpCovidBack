package seed

import (
	"log"
	"time"

	"github.com/HackathonCovid/helpCovidBack/api/models"
	"github.com/jinzhu/gorm"
)

// Load : Validation and Join
func Load(db *gorm.DB) {

	input := "2020-10-01"
	layout := "2006-01-02"
	start, _ := time.Parse(layout, input)

	input2 := "2020-10-12"
	layout2 := "2006-01-02"
	end, _ := time.Parse(layout2, input2)

	var users = []models.User{
		models.User{
			Firstname:   "Steve",
			Lastname:    "Victor",
			Email:       "steve@gmail.com",
			Isvolunteer: 1,
			Password:    "password",
			TypeOrga:    "Clinique",
			OrgaName:    "Paris Bichat",
			City:        "Paris",
			PhoneNumber: "0625321458",
			Adress:      "3 boulevard de la République",
			Description: "Developer Go",
			Degree:      "Test",
			Longitude:   "",
			Latitude:    "",
		},
		models.User{
			Firstname:   "Kevin",
			Lastname:    "Feige",
			Email:       "feige@gmail.com",
			Isvolunteer: 1,
			Password:    "password",
			TypeOrga:    "Hopital",
			OrgaName:    "Paris Salpatriere",
			City:        "Paris",
			PhoneNumber: "0625321458",
			Adress:      "3 boulevard de la République",
			Description: "Developer Go",
			Degree:      "Test",
			Longitude:   "",
			Latitude:    "",
		},
		models.User{
			Firstname:   "Test",
			Lastname:    "Victor",
			Email:       "steve123@gmail.com",
			Isvolunteer: 1,
			Password:    "password",
			TypeOrga:    "Ehpad",
			OrgaName:    "Bichat 2",
			City:        "Paris",
			PhoneNumber: "0625321458",
			Adress:      "3 boulevard de la République",
			Description: "Developer Go",
			Degree:      "Test",
			Longitude:   "",
			Latitude:    "",
		},
	}

	var missions = []models.Mission{
		models.Mission{
			Title:            "Aide clinique paris 11",
			Description:      "Aider pour les salles de réa",
			StartDate:        start,
			EndDate:          end,
			NbDays:           5,
			NbPeopleRequired: 2,
			SkillsRequired:   "Aucune",
			NightOrDay:       "Night",
			AddressHospital:  "30 rue Kilford",
			AuthorID:         1,
		},
		models.Mission{
			Title:            "Aide clinique 92",
			Description:      "Aider pour le ménage",
			StartDate:        start,
			EndDate:          end,
			NbDays:           3,
			NbPeopleRequired: 2,
			SkillsRequired:   "Aucune",
			NightOrDay:       "Day",
			AddressHospital:  "36 rue du Général Leclerc",
			AuthorID:         1,
		},
		models.Mission{
			Title:            "Aide hopital",
			Description:      "Aider pour ranger matériel",
			StartDate:        start,
			EndDate:          end,
			NbDays:           3,
			NbPeopleRequired: 2,
			SkillsRequired:   "Aucune",
			NightOrDay:       "Day",
			AddressHospital:  "403 avenue de la République",
			AuthorID:         2,
		},
	}

	var comments = []models.Comment{
		models.Comment{
			Body:      "Super, merci",
			UserID:    1,
			MissionID: 1,
		},
		models.Comment{
			Body:      "Ne fonctionne pas",
			UserID:    2,
			MissionID: 1,
		},
		models.Comment{
			Body:      "Test comment",
			UserID:    1,
			MissionID: 2,
		},
	}

	err := db.Debug().DropTableIfExists(&models.Comment{}, &models.Mission{}, &models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.User{}, &models.Mission{}, &models.Comment{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	err = db.Debug().Model(&models.Mission{}).AddForeignKey("author_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	err = db.Debug().Model(&models.Comment{}).AddForeignKey("user_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	err = db.Debug().Model(&models.Comment{}).AddForeignKey("mission_id", "missions(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	for i := range users {
		err = db.Debug().Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		missions[i].AuthorID = users[i].ID
		//comments[i].UserID = users[i].ID

		err = db.Debug().Model(&models.Mission{}).Create(&missions[i]).Error
		if err != nil {
			log.Fatalf("cannot seed missions table: %v", err)
		}
		//comments[i].MissionID = missions[i].ID

		err = db.Debug().Model(&models.Comment{}).Create(&comments[i]).Error
		if err != nil {
			log.Fatalf("cannot seed comments table: %v", err)
		}
	}
}
