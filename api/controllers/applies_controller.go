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

// ApplyMission : function to apply to mission
func (server *Server) ApplyMission(c *gin.Context) {

	//clear previous error if any
	errList = map[string]string{}

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
	uid, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}
	// check if the user exist
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
	// check if the mission exist
	mission := models.Mission{}
	err = server.DB.Debug().Model(models.Mission{}).Where("id = ?", pid).Take(&mission).Error
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	apply := models.Apply{}
	apply.UserID = user.ID
	apply.MissionID = mission.ID

	applyCreated, err := apply.SaveApply(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		errList = formattedError
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":   http.StatusCreated,
		"response": applyCreated,
	})
}

// GetApplies : funtion to get the applies
func (server *Server) GetApplies(c *gin.Context) {
	//clear previous error if any
	errList = map[string]string{}

	missionID := c.Param("id")

	// Is a valid mission id given to us?
	pid, err := strconv.ParseUint(missionID, 10, 64)
	if err != nil {
		fmt.Println("this is the error: ", err)
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
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

	apply := models.Apply{}

	applies, err := apply.GetAppliesInfo(server.DB, pid)
	if err != nil {
		errList["No_applies"] = "No Applies found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": applies,
	})
}

// ValidateApply : function to validate an apply of a mission
func (server *Server) ValidateApply(c *gin.Context) {
	//clear previous error if any
	errList = map[string]string{}

	applyID := c.Param("id")
	// Check if the apply id is valid
	pid, err := strconv.ParseUint(applyID, 10, 64)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}

	//Check if the apply exist
	origApply := models.Apply{}
	err = server.DB.Debug().Model(models.Apply{}).Where("id = ?", pid).Take(&origApply).Error
	if err != nil {
		errList["No_apply"] = "No Apply Found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
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
	apply := models.Apply{}
	err = json.Unmarshal(body, &apply)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	apply.ID = origApply.ID
	apply.UserID = origApply.UserID
	apply.MissionID = origApply.MissionID

	applyUpdated, err := apply.FindApplyAndUpdateByID(server.DB)
	if err != nil {
		formattedError := formaterror.FormatError(err.Error())
		errList = formattedError
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": applyUpdated,
	})

}

// WithdrawApply : funtion to withdraw an apply of a mission
func (server *Server) WithdrawApply(c *gin.Context) {

	applyID := c.Param("id")
	// Is a valid apply id given to us?
	lid, err := strconv.ParseUint(applyID, 10, 64)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}
	// Is this user authenticated?
	uid, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}
	// Check if the apply exist
	apply := models.Apply{}
	err = server.DB.Debug().Model(models.Apply{}).Where("id = ?", lid).Take(&apply).Error
	if err != nil {
		errList["No_apply"] = "No Apply Found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	// Is the authenticated user, the owner of this apply?
	if uid != apply.UserID {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	// If all the conditions are met, delete the apply
	_, err = apply.DeleteApply(server.DB)
	if err != nil {
		errList["Other_error"] = "Please try again later"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": "Apply deleted",
	})
}

func (server *Server) GetAppliesById(c *gin.Context) {
	//clear previous error if any
	errList = map[string]string{}

	missionID := c.Param("id")

	// Is a valid mission id given to us?
	pid, err := strconv.ParseUint(missionID, 10, 64)
	if err != nil {
		fmt.Println("this is the error: ", err)
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}
	// check if the user exist
	user := models.User{}
	err = server.DB.Debug().Model(models.User{}).Where("id = ?", pid).Take(&user).Error
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}

	apply := models.Apply{}

	applies, err := apply.GetAppliesByUserId(server.DB, pid)
	if err != nil {
		errList["No_applies"] = "No Applies found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": applies,
	})
}
