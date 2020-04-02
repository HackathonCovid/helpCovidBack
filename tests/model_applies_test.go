package tests

import (
	"log"
	"testing"

	"github.com/HackathonCovid/helpCovidBack/api/models"
	"github.com/stretchr/testify/assert"
)

func TestSaveAnApply(t *testing.T) {
	err := refreshUserMissionAndApplyTable()
	if err != nil {
		log.Fatalf("Error refreshing user, mission and apply table %v\n", err)
	}
	user, mission, err := seedOneUserAndOneMission()
	if err != nil {
		log.Fatalf("Cannot seed user and mission %v\n", err)
	}
	newApply := models.Apply{
		ID:        1,
		UserID:    user.ID,
		MissionID: mission.ID,
	}
	savedApply, err := newApply.SaveApply(server.DB)
	if err != nil {
		t.Errorf("this is the error getting the apply: %v\n", err)
		return
	}
	assert.Equal(t, newApply.ID, savedApply.ID)
	assert.Equal(t, newApply.UserID, savedApply.UserID)
	assert.Equal(t, newApply.MissionID, savedApply.MissionID)
}

func TestGetApplyInfoForAMission(t *testing.T) {

	err := refreshUserMissionAndApplyTable()
	if err != nil {
		log.Fatalf("Error refreshing user, mission and apply table %v\n", err)
	}
	mission, users, applies, err := seedUsersMissionsAndApplies()
	if err != nil {
		log.Fatalf("Error seeding user, mission and apply table %v\n", err)
	}
	//Where applyInstance is an instance of the post initialize in setup_test.go
	_, err = applyInstance.GetAppliesInfo(server.DB, mission.ID)
	if err != nil {
		t.Errorf("this is the error getting the applies: %v\n", err)
		return
	}
	assert.Equal(t, len(applies), 2)
	assert.Equal(t, len(users), 2) //two users apply to the mission
}

func TestDeleteAnApply(t *testing.T) {

	err := refreshUserMissionAndApplyTable()
	if err != nil {
		log.Fatalf("Error refreshing user, mission and apply table %v\n", err)
	}
	_, _, applies, err := seedUsersMissionsAndApplies()
	if err != nil {
		log.Fatalf("Error seeding user, mission and apply table %v\n", err)
	}
	// Delete the first apply
	for _, v := range applies {
		if v.ID == 2 {
			continue
		}
		applyInstance.ID = v.ID //applyInstance is defined in setup_test.go
	}
	deletedApply, err := applyInstance.DeleteApply(server.DB)
	if err != nil {
		t.Errorf("this is the error deleting the apply: %v\n", err)
		return
	}
	assert.Equal(t, deletedApply.ID, applyInstance.ID)
}

// When a mission is deleted, delete its applies
func TestDeleteAppliesForAMission(t *testing.T) {

	err := refreshUserMissionAndApplyTable()
	if err != nil {
		log.Fatalf("Error refreshing user, mission and apply table %v\n", err)
	}
	mission, _, _, err := seedUsersMissionsAndApplies()
	if err != nil {
		log.Fatalf("Error seeding user, mission and apply table %v\n", err)
	}
	numberDeleted, err := applyInstance.DeleteMissionApplies(server.DB, mission.ID)
	if err != nil {
		t.Errorf("this is the error deleting the apply: %v\n", err)
		return
	}
	assert.Equal(t, numberDeleted, int64(2))
}

// When a user is deleted, delete its applies
func TestDeleteAppliesForAUser(t *testing.T) {
	var userID uint64
	err := refreshUserMissionAndApplyTable()
	if err != nil {
		log.Fatalf("Error refreshing user, mission and apply table %v\n", err)
	}
	_, users, applies, err := seedUsersMissionsAndApplies()
	if err != nil {
		log.Fatalf("Error seeding user, mission and apply table %v\n", err)
	}
	for _, v := range applies {
		if v.ID == 2 {
			continue
		}
		applyInstance.ID = v.ID //applyInstance is defined in setup_test.go
	}
	// get the first user, this user has one apply
	for _, v := range users {
		if v.ID == 2 {
			continue
		}
		userID = v.ID
	}
	numberDeleted, err := applyInstance.DeleteUserApplies(server.DB, userID)
	if err != nil {
		t.Errorf("this is the error deleting the apply: %v\n", err)
		return
	}
	assert.Equal(t, numberDeleted, int64(1))
}
