package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/HackathonCovid/helpCovidBack/api/auth"
	"github.com/HackathonCovid/helpCovidBack/api/models"
	"github.com/HackathonCovid/helpCovidBack/api/utils/formaterror"
	"github.com/gin-gonic/gin"
)

//CreateMission : function to create a mission
func (server *Server) CreateMission(c *gin.Context) {

	//clear previous error if any
	errList = map[string]string{}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["Invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	mission := models.Mission{}

	err = json.Unmarshal(body, &mission)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	uid, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	// check if the user exist:
	user := models.User{}
	err = server.DB.Debug().Model(models.User{}).Where("id = ?", uid).Take(&user).Error
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	mission.AuthorID = uid //the authenticated user is the one creating the mission

	mission.Prepare()
	errorMessages := mission.Validate()
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
/*
	reader := bufio.NewReader(mission.FileBase64)
    content, _ := ioutil.ReadAll(reader)
	mission.FileBase64 = base64.StdEncoding.EncodeToString(content)*/

	missionCreated, err := mission.SaveMission(server.DB)
	if err != nil {
		errList := formaterror.FormatError(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":   http.StatusCreated,
		"response": missionCreated,
	})
}

//GetMissions : function to get all missions
func (server *Server) GetMissions(c *gin.Context) {

	mission := models.Mission{}

	missions, err := mission.FindAllMissions(server.DB)
	if err != nil {
		errList["No_mission"] = "No Mission Found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": missions,
	})
}

//GetMission : function to get a mission
func (server *Server) GetMission(c *gin.Context) {

	missionID := c.Param("id")
	pid, err := strconv.ParseUint(missionID, 10, 64)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}
	mission := models.Mission{}

	missionReceived, err := mission.FindMissionByID(server.DB, pid)
	if err != nil {
		errList["No_mission"] = "No Mission Found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": missionReceived,
	})
}

//UpdateMission : function to update a mission
func (server *Server) UpdateMission(c *gin.Context) {

	//clear previous error if any
	errList = map[string]string{}

	missionID := c.Param("id")
	// Check if the mission id is valid
	pid, err := strconv.ParseUint(missionID, 10, 64)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}
	//Check if the auth token is valid and get the user id from it
	uid, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}
	//Check if the mission exist
	origMission := models.Mission{}
	err = server.DB.Debug().Model(models.Mission{}).Where("id = ?", pid).Take(&origMission).Error
	if err != nil {
		errList["No_mission"] = "No Mission Found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	if uid != origMission.AuthorID {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}
	// Read the data posted
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["Invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	// Start processing the request data
	mission := models.Mission{}
	err = json.Unmarshal(body, &mission)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	mission.ID = origMission.ID //this is important to tell the model, the mission id to update.
	mission.AuthorID = origMission.AuthorID

	mission.Prepare()
	errorMessages := mission.Validate()
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	missionUpdated, err := mission.UpdateAMission(server.DB, pid)
	if err != nil {
		errList := formaterror.FormatError(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": missionUpdated,
	})
}

//DeleteMission : function to delete a mission
func (server *Server) DeleteMission(c *gin.Context) {

	missionID := c.Param("id")
	// Is a valid mission id given to us?
	pid, err := strconv.ParseUint(missionID, 10, 64)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}

	fmt.Println("delete a mission")

	// Is this user authenticated ?
	uid, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}
	// Check if the mission exist
	mission := models.Mission{}
	err = server.DB.Debug().Model(models.Mission{}).Where("id = ?", pid).Take(&mission).Error
	if err != nil {
		errList["No_mission"] = "No Mission Found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	// Is the authenticated user, the owner of this mission ?
	if uid != mission.AuthorID {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}
	// If all the conditions are met, delete the mission
	_, err = mission.DeleteAMission(server.DB, pid, uid)
	if err != nil {
		errList["Other_error"] = "Please try again later"
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  errList,
		})
		return
	}
	comment := models.Comment{}
	apply := models.Apply{}

	// Also delete the comments and applies that this mission have.
	_, err = comment.DeleteMissionComments(server.DB, pid)

	_, err = apply.DeleteMissionApplies(server.DB, pid)

	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": "Mission deleted",
	})
}

//GetUserMissions : function to get the missions of a user
func (server *Server) GetUserMissions(c *gin.Context) {

	userID := c.Param("id")
	// Is a valid user id given to us?
	uid, err := strconv.ParseUint(userID, 10, 64)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}
	mission := models.Mission{}
	missions, err := mission.FindUserMissions(server.DB, uint64(uid))
	if err != nil {
		errList["No_mission"] = "No Mission Found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": missions,
	})
}
