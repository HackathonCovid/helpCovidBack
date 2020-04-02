package tests

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/HackathonCovid/helpCovidBack/api/controllers"
	"github.com/HackathonCovid/helpCovidBack/api/models"
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
)

var server = controllers.Server{}
var userInstance = models.User{}
var missionInstance = models.Mission{}
var commentInstance = models.Comment{}
var applyInstance = models.Apply{}

func TestMain(m *testing.M) {
	// We have this part because we don't use circle CI
	var err error
	err = godotenv.Load(os.ExpandEnv("./../.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}
	// UNTIL HERE
	Database()

	os.Exit(m.Run())
}

func Database() {
	var err error

	////////////////////////////////// UNCOMMENT THIS WHILE TESTING ON LOCAL(WITHOUT USING CIRCLE CI) ///////////////////////
	TestDbDriver := os.Getenv("TEST_DB_DRIVER")
	if TestDbDriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", os.Getenv("TEST_DB_HOST"), os.Getenv("TEST_DB_PORT"), os.Getenv("TEST_DB_USER"), os.Getenv("TEST_DB_NAME"), os.Getenv("TEST_DB_PASSWORD"))
		server.DB, err = gorm.Open(TestDbDriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database\n", TestDbDriver)
		}
	}
	/////////////////////////////////  END OF LOCAL TEST DATABASE SETUP ///////////////////////////////////////////////////

	////////////////////////////////// COMMENT THIS WHILE TESTING ON LOCAL (USING CIRCLE CI)  //////////////////////
	// WE HAVE TO INPUT TESTING DATA MANUALLY BECAUSE CIRCLECI, CANNOT READ THE ".env" FILE WHICH, WE WOULD HAVE ADDED THE TEST CONFIG THERE
	// SO MANUALLY ADD THE NAME OF THE DATABASE, THE USER AND THE PASSWORD, AS SEEN BELOW:
	//DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", "127.0.0.1", "5432", "steven", "forum_db_test", "password")
	//server.DB, err = gorm.Open("postgres", DBURL)
	//if err != nil {
	//	fmt.Printf("Cannot connect to %s database\n", "postgres")
	//	log.Fatal("This is the error:", err)
	//} else {
	//	fmt.Printf("We are connected to the %s database\n", "postgres")
	//}
	//////////////////////////////// END OF USING CIRCLE CI ////////////////////////////////////////////////////////////////
}

func refreshUserTable() error {
	err := server.DB.DropTableIfExists(&models.User{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.User{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed table")
	return nil
}

func seedOneUser() (models.User, error) {

	refreshUserTable()

	input1 := "1996-02-08"
	layout1 := "2006-01-02"
	t, _ := time.Parse(layout1, input1)

	user := models.User{
		Firstname:    "Pet",
		Lastname:     "Last",
		Email:        "pet@gmail.com",
		Password:     "password",
		Isvolunteer:  1,
		Dateofbirth:  t,
		Sexe:         "Homme",
		City:         "Paris",
		PhoneNumber:  "0625321458",
		Adress:       "3 boulevard de la République",
		Description:  "Developer Go",
		Degree:       "Test",
		Longitude:    "",
		Latitude:     "",
		HospitalName: "",
	}

	err := server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func seedUsers() ([]models.User, error) {

	var err error
	if err != nil {
		return nil, err
	}

	input1 := "1996-02-08"
	layout1 := "2006-01-02"
	t, _ := time.Parse(layout1, input1)

	users := []models.User{
		models.User{
			Firstname:    "Steven",
			Lastname:     "Victor",
			Email:        "steven@gmail.com",
			Password:     "password",
			Isvolunteer:  1,
			Dateofbirth:  t,
			Sexe:         "Homme",
			City:         "Paris",
			PhoneNumber:  "0625321458",
			Adress:       "3 boulevard de la République",
			Description:  "Developer Go",
			Degree:       "Test",
			Longitude:    "",
			Latitude:     "",
			HospitalName: "",
		},
		models.User{
			Firstname:    "Kenny",
			Lastname:     "Morris",
			Email:        "kenny@gmail.com",
			Password:     "password",
			Isvolunteer:  1,
			Dateofbirth:  t,
			Sexe:         "Homme",
			City:         "Paris",
			PhoneNumber:  "0625321458",
			Adress:       "3 boulevard de la République",
			Description:  "Developer Go",
			Degree:       "Test",
			Longitude:    "",
			Latitude:     "",
			HospitalName: "",
		},
	}

	for i := range users {
		err := server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			return []models.User{}, err
		}
	}
	return users, nil
}

func refreshUserAndMissionTable() error {

	err := server.DB.DropTableIfExists(&models.User{}, &models.Mission{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.User{}, &models.Mission{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed tables")
	return nil
}

func seedOneUserAndOneMission() (models.User, models.Mission, error) {

	err := refreshUserAndMissionTable()
	if err != nil {
		return models.User{}, models.Mission{}, err
	}

	input1 := "1996-02-08"
	layout1 := "2006-01-02"
	t, _ := time.Parse(layout1, input1)

	user := models.User{
		Firstname:    "Sam",
		Lastname:     "Phil",
		Email:        "sam@gmail.com",
		Password:     "password",
		Isvolunteer:  1,
		Dateofbirth:  t,
		Sexe:         "Homme",
		City:         "Paris",
		PhoneNumber:  "0625321458",
		Adress:       "3 boulevard de la République",
		Description:  "Developer Go",
		Degree:       "Test",
		Longitude:    "",
		Latitude:     "",
		HospitalName: "",
	}

	err = server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		return models.User{}, models.Mission{}, err
	}

	input := "2020-01-06"
	layout := "2006-01-02"
	start, _ := time.Parse(layout, input)

	input2 := "2020-10-11"
	layout2 := "2006-01-02"
	end, _ := time.Parse(layout2, input2)

	mission := models.Mission{
		Title:            "Aide clinique paris 11",
		Description:      "Aider pour les salles de réa",
		StartDate:        start,
		EndDate:          end,
		NbDays:           5,
		NbPeopleRequired: 2,
		SkillsRequired:   "Aucune",
		NightOrDay:       "Night",
		AddressHospital:  "30 rue Kilford",
		AuthorID:         user.ID,
	}

	err = server.DB.Model(&models.Mission{}).Create(&mission).Error
	if err != nil {
		return models.User{}, models.Mission{}, err
	}
	return user, mission, nil
}

func seedUsersAndMissions() ([]models.User, []models.Mission, error) {

	var err error

	if err != nil {
		return []models.User{}, []models.Mission{}, err
	}

	input1 := "1996-02-08"
	layout1 := "2006-01-02"
	t, _ := time.Parse(layout1, input1)

	var users = []models.User{
		models.User{
			Firstname:    "Steven",
			Lastname:     "Victor",
			Email:        "steven@gmail.com",
			Password:     "password",
			Isvolunteer:  1,
			Dateofbirth:  t,
			Sexe:         "Homme",
			City:         "Paris",
			PhoneNumber:  "0625321458",
			Adress:       "3 boulevard de la République",
			Description:  "Developer Go",
			Degree:       "Test",
			Longitude:    "",
			Latitude:     "",
			HospitalName: "",
		},
		models.User{
			Firstname:    "Kenny",
			Lastname:     "Morris",
			Email:        "kenny@gmail.com",
			Password:     "password",
			Isvolunteer:  1,
			Dateofbirth:  t,
			Sexe:         "Homme",
			City:         "Paris",
			PhoneNumber:  "0625321458",
			Adress:       "3 boulevard de la République",
			Description:  "Developer Go",
			Degree:       "Test",
			Longitude:    "",
			Latitude:     "",
			HospitalName: "",
		},
	}

	input := "2020-09-01"
	layout := "2006-01-02"
	start, _ := time.Parse(layout, input)

	input2 := "2020-10-01"
	layout2 := "2006-01-02"
	end, _ := time.Parse(layout2, input2)

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
		},
		models.Mission{
			Title:            "Aide hopital paris 11",
			Description:      "Aider pour les malades",
			StartDate:        start,
			EndDate:          end,
			NbDays:           5,
			NbPeopleRequired: 2,
			SkillsRequired:   "Aucune",
			NightOrDay:       "Night",
			AddressHospital:  "30 rue Kilford",
		},
	}

	for i := range users {
		err = server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		missions[i].AuthorID = users[i].ID

		err = server.DB.Model(&models.Mission{}).Create(&missions[i]).Error
		if err != nil {
			log.Fatalf("cannot seed missions table: %v", err)
		}
	}
	return users, missions, nil
}

func refreshUserMissionAndApplyTable() error {
	err := server.DB.DropTableIfExists(&models.User{}, &models.Mission{}, &models.Apply{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.User{}, &models.Mission{}, &models.Apply{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed user, mission and apply tables")
	return nil
}

func seedUsersMissionsAndApplies() (models.Mission, []models.User, []models.Apply, error) {
	// The idea here is: two users can apply to one mission
	var err error

	input1 := "1996-02-08"
	layout1 := "2006-01-02"
	t, _ := time.Parse(layout1, input1)

	var users = []models.User{
		models.User{
			Firstname:    "Dwayne",
			Lastname:     "Heu",
			Email:        "dwayne@gmail.com",
			Password:     "password",
			Isvolunteer:  1,
			Dateofbirth:  t,
			Sexe:         "Homme",
			City:         "Paris",
			PhoneNumber:  "0625321458",
			Adress:       "3 boulevard de la République",
			Description:  "Developer Go",
			Degree:       "Test",
			Longitude:    "",
			Latitude:     "",
			HospitalName: "",
		},
		models.User{
			Firstname:    "Jen",
			Lastname:     "Stella",
			Email:        "jen@gmail.com",
			Password:     "password",
			Isvolunteer:  1,
			Dateofbirth:  t,
			Sexe:         "Femme",
			City:         "Paris",
			PhoneNumber:  "0625321458",
			Adress:       "3 boulevard de la République",
			Description:  "Developer Go",
			Degree:       "Test",
			Longitude:    "",
			Latitude:     "",
			HospitalName: "",
		},
	}

	input := "2020-09-01"
	layout := "2006-01-02"
	start, _ := time.Parse(layout, input)

	input2 := "2020-10-01"
	layout2 := "2006-01-02"
	end, _ := time.Parse(layout2, input2)

	mission := models.Mission{
		Title:            "Hopital au sud de la France",
		Description:      "Aider les malades",
		StartDate:        start,
		EndDate:          end,
		NbDays:           5,
		NbPeopleRequired: 2,
		SkillsRequired:   "Aucune",
		NightOrDay:       "Night",
		AddressHospital:  "30 rue Kilford",
	}

	err = server.DB.Model(&models.Mission{}).Create(&mission).Error
	if err != nil {
		log.Fatalf("cannot seed mission table: %v", err)
	}

	var applies = []models.Apply{
		models.Apply{
			UserID:    1,
			MissionID: mission.ID,
		},
		models.Apply{
			UserID:    2,
			MissionID: mission.ID,
		},
	}

	for i := range users {
		err = server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		err = server.DB.Model(&models.Apply{}).Create(&applies[i]).Error
		if err != nil {
			log.Fatalf("cannot seed applies table: %v", err)
		}
	}
	return mission, users, applies, nil
}

func refreshUserMissionAndCommentTable() error {
	err := server.DB.DropTableIfExists(&models.Mission{}, &models.Comment{}, &models.User{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.Mission{}, &models.Comment{}, &models.User{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed user, mission and comment tables")
	return nil
}

func seedUsersMissionsAndComments() (models.Mission, []models.User, []models.Comment, error) {
	// The idea here is: two users can comment one mission
	var err error

	input1 := "1996-02-08"
	layout1 := "2006-01-02"
	t, _ := time.Parse(layout1, input1)

	var users = []models.User{
		models.User{
			Firstname:    "Jon",
			Lastname:     "Doe",
			Email:        "doe@gmail.com",
			Password:     "password",
			Isvolunteer:  1,
			Dateofbirth:  t,
			Sexe:         "Homme",
			City:         "Paris",
			PhoneNumber:  "0625321458",
			Adress:       "3 boulevard de la République",
			Description:  "Developer Go",
			Degree:       "Test",
			Longitude:    "",
			Latitude:     "",
			HospitalName: "",
		},
		models.User{
			Firstname:    "Jenny",
			Lastname:     "Ellen",
			Email:        "jenny@gmail.com",
			Password:     "password",
			Isvolunteer:  1,
			Dateofbirth:  t,
			Sexe:         "Femme",
			City:         "Paris",
			PhoneNumber:  "0625321458",
			Adress:       "3 boulevard de la République",
			Description:  "Developer Go",
			Degree:       "Test",
			Longitude:    "",
			Latitude:     "",
			HospitalName: "",
		},
	}

	input := "2020-09-01"
	layout := "2006-01-02"
	start, _ := time.Parse(layout, input)

	input2 := "2020-10-01"
	layout2 := "2006-01-02"
	end, _ := time.Parse(layout2, input2)

	mission := models.Mission{
		Title:            "Aide clinique",
		Description:      "Aider pour les salles de réa",
		StartDate:        start,
		EndDate:          end,
		NbDays:           5,
		NbPeopleRequired: 2,
		SkillsRequired:   "Aucune",
		NightOrDay:       "Night",
		AddressHospital:  "30 rue Kilford",
	}

	err = server.DB.Model(&models.Mission{}).Create(&mission).Error
	if err != nil {
		log.Fatalf("cannot seed mission table: %v", err)
	}

	var comments = []models.Comment{
		models.Comment{
			Body:      "user 1 made this comment",
			UserID:    1,
			MissionID: mission.ID,
		},
		models.Comment{
			Body:      "user 2 made this comment",
			UserID:    2,
			MissionID: mission.ID,
		},
	}

	for i := range users {
		err = server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		err = server.DB.Model(&models.Comment{}).Create(&comments[i]).Error
		if err != nil {
			log.Fatalf("cannot seed comments table: %v", err)
		}
	}
	return mission, users, comments, nil
}

func refreshUserAndResetPasswordTable() error {
	err := server.DB.DropTableIfExists(&models.User{}, &models.ResetPassword{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.User{}, &models.ResetPassword{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed user and resetpassword tables")
	return nil
}

// Seed the reset password table with the token
func seedResetPassword() (models.ResetPassword, error) {

	resetDetails := models.ResetPassword{
		Token: "awesometoken",
		Email: "pet@example.com",
	}
	err := server.DB.Model(&models.ResetPassword{}).Create(&resetDetails).Error
	if err != nil {
		return models.ResetPassword{}, err
	}
	return resetDetails, nil
}
