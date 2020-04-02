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

// CreateComment : function to create a comment
func (server *Server) CreateComment(c *gin.Context) {
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
	// check the token
	uid, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}
	// check if the user exists
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
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		errList["Invalid_body"] = "Unable to get request"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}
	comment := models.Comment{}
	err = json.Unmarshal(body, &comment)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	// enter the userid and the missionid. The comment body is automatically passed
	comment.UserID = uid
	comment.MissionID = pid

	comment.Prepare()
	errorMessages := comment.Validate("")
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	//	fmt.Printf("User id: %v", mission.ID)
	//	fmt.Printf(" Vote id: %v", comment.MissionID)
	fmt.Printf("User id: %v", uid)
	fmt.Printf("mission id: %v", pid)

	//	server.DB.Model(&vote).Association("Users").Append(&user);

	/*	commentCreated, err := comment.SaveComment(server.DB, &mission)
		if err != nil {
			formattedError := formaterror.FormatError(err.Error())
			errList = formattedError
			c.JSON(http.StatusNotFound, gin.H{
				"status": http.StatusNotFound,
				"error":  errList,
			})
			return
		}
		c.JSON(http.StatusCreated, gin.H{
			"status":   http.StatusCreated,
			"response": commentCreated,
		})*/

	err_2 := server.DB.Model(&mission).Association("Comments").Append(&comment).Error
	if err_2 != nil {
		formattedError := formaterror.FormatError(err_2.Error())
		errList = formattedError
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"status":   http.StatusCreated,
		"response": comment,
	})

}

// GetComments : funtion to get the comments
func (server *Server) GetComments(c *gin.Context) {
	//clear previous error if any
	errList = map[string]string{}

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
	// check if the mission exist:
	mission := models.Mission{}
	err = server.DB.Debug().Model(models.Mission{}).Where("id = ?", pid).Take(&mission).Error
	if err != nil {
		errList["No_mission"] = "No mission found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}

	comment := models.Comment{}

	comments, err := comment.GetComments(server.DB, pid)
	if err != nil {
		errList["No_comments"] = "No comments found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":   http.StatusOK,
		"response": comments,
	})
}

// UpdateComment : function to update a comment
func (server *Server) UpdateComment(c *gin.Context) {
	//clear previous error if any
	errList = map[string]string{}

	commentID := c.Param("id")
	// Check if the comment id is valid
	pid, err := strconv.ParseUint(commentID, 10, 64)
	if err != nil {
		errList["Invalid_request"] = "Invalid Request"
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"error":  errList,
		})
		return
	}
	//Check if the auth token is valid and  get the user id from it
	uid, err := auth.ExtractTokenID(c.Request)
	if err != nil {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}
	//Check if the comment exist
	origComment := models.Comment{}
	err = server.DB.Debug().Model(models.Comment{}).Where("id = ?", pid).Take(&origComment).Error
	if err != nil {
		errList["No_comment"] = "No Comment Found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	if uid != origComment.UserID {
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
	comment := models.Comment{}
	err = json.Unmarshal(body, &comment)
	if err != nil {
		errList["Unmarshal_error"] = "Cannot unmarshal body"
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	comment.Prepare()
	errorMessages := comment.Validate("")
	if len(errorMessages) > 0 {
		errList = errorMessages
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"status": http.StatusUnprocessableEntity,
			"error":  errList,
		})
		return
	}

	comment.ID = origComment.ID //this is important to tell the model the comment id to update, the other update field are set above
	comment.UserID = origComment.UserID
	comment.MissionID = origComment.MissionID

	commentUpdated, err := comment.UpdateAComment(server.DB)
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
		"response": commentUpdated,
	})
}

// DeleteComment : function to delete a comment
func (server *Server) DeleteComment(c *gin.Context) {
	commentID := c.Param("id")
	// Is a valid post id given to us?
	cid, err := strconv.ParseUint(commentID, 10, 64)
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

	// Check if the comment exist
	comment := models.Comment{}
	err = server.DB.Debug().Model(models.Comment{}).Where("id = ?", cid).Take(&comment).Error
	if err != nil {
		errList["No_comment"] = "No Comment Found"
		c.JSON(http.StatusNotFound, gin.H{
			"status": http.StatusNotFound,
			"error":  errList,
		})
		return
	}
	// Is the authenticated user, the owner of this comment?
	if uid != comment.UserID {
		errList["Unauthorized"] = "Unauthorized"
		c.JSON(http.StatusUnauthorized, gin.H{
			"status": http.StatusUnauthorized,
			"error":  errList,
		})
		return
	}
	// If all the conditions are met, delete the comment
	_, err = comment.DeleteAComment(server.DB)
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
		"response": "Comment deleted",
	})
}
