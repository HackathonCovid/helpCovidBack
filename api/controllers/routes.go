package controllers

import "github.com/HackathonCovid/helpCovidBack/api/middlewares"

func (s *Server) initializeRoutes() {

	v1 := s.Router.Group("/api")
	{

		// Login Route
		v1.POST("/login", s.Login)

		// Reset password:
		v1.POST("/password/forgot", s.ForgotPassword)
		v1.POST("/password/reset", s.ResetPassword)
		v1.POST("/password/userreset", middlewares.TokenAuthMiddleware(), s.ResetUserPassword)

		//Users routes
		v1.POST("/users", s.CreateUser)
		v1.GET("/users", s.GetUsers)
		v1.GET("/users/:id", s.GetUser)
		v1.PUT("/users/:id", middlewares.TokenAuthMiddleware(), s.UpdateUser)
		v1.DELETE("/users/:id", middlewares.TokenAuthMiddleware(), s.DeleteUser)

		// Missions routes
		v1.POST("/missions", middlewares.TokenAuthMiddleware(), s.CreateMission)
		v1.GET("/missions", s.GetMissions)
		v1.GET("/missions/:id", s.GetMission)
		v1.PUT("/missions/:id", middlewares.TokenAuthMiddleware(), s.UpdateMission)
		v1.DELETE("/missions/:id", middlewares.TokenAuthMiddleware(), s.DeleteMission)
		v1.GET("/user_missions/:id", s.GetUserMissions)
		//v1.POST("/user_missions/:id",  middlewares.TokenAuthMiddleware(), s.AddUserToMission)
		//v1.DELETE("/user_missions/:id", middlewares.TokenAuthMiddleware(), s.DeleteUserFromMission)

		// Applies routes
		v1.GET("/applies/:id", s.GetApplies)
		v1.POST("/applies/:id", middlewares.TokenAuthMiddleware(), s.ApplyMission)
		v1.DELETE("/applies/:id", middlewares.TokenAuthMiddleware(), s.WithdrawApply)
		v1.GET("/userapplies/:id", s.GetAppliesById)
		v1.POST("/validate/:id", s.ValidateApply)

		// Comments routes
		v1.POST("/comments/:id", middlewares.TokenAuthMiddleware(), s.CreateComment)
		v1.GET("/comments/:id", s.GetComments)
		v1.PUT("/comments/:id", middlewares.TokenAuthMiddleware(), s.UpdateComment)
		v1.DELETE("/comments/:id", middlewares.TokenAuthMiddleware(), s.DeleteComment)
	}
}
