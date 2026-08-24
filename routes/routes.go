package routes

import (
	"github.com/gin-gonic/gin"
	"myutilityx.com/http"
	"myutilityx.com/middlewares"
)

func RegisterRoutes() *gin.Engine {
	server := gin.Default()
	server.MaxMultipartMemory = 32 << 20 // 32 MiB
	server.Use(middlewares.CORSMiddleware())
	authenticated := server.Group("/")
	authenticated.Use(middlewares.Authenticate)
	//Home
	server.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{"API": "WORKING"})
	})
	//links
	authenticated.POST("/url/shrink", addLink)
	authenticated.GET("/url", getAllLinks)
	authenticated.DELETE("/url/:id", deleteUrl)
	server.GET("/:shorturl", getSingleUrl)
	server.POST("/graphql", http.GraphQLHandler())
	server.GET("/check-db-connection", checkMongoDBConnection)
	//users
	server.POST("/user/register", register)
	server.POST("/contact-form", contactForm)
	server.GET("/user/verify", verifyEmail)
	server.POST("/user/login", login)
	server.POST("/user/reset-password", resetPassword)
	server.POST("/user/reset-password/confirm", resetPasswordVerify)
	authenticated.POST("/user/update-role", setAdminRole)
	//SMS
	server.POST("/save-sms", saveSms)
	server.GET("/get-sms", getSms)
	//expenses statistics
	authenticated.POST("/expense/add", addExpense)
	authenticated.PUT("/expense/update/:id", updateExpense)
	authenticated.GET("/expense/all", GetAllExpenses)

	server.GET("/get-file/:fileId", getFile)
	authenticated.POST("/upload-file", uploadFile)
	authenticated.GET("/files", getFiles)
	authenticated.DELETE("/file/:fileId", deleteFile)
	return server
}
