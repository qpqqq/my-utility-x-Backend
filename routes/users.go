package routes

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"myutilityx.com/db"
	"myutilityx.com/errors"
	"myutilityx.com/mailS"
	"myutilityx.com/models"
	"myutilityx.com/utils"
)

func register(ctx *gin.Context) {
	var user models.User
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.ErrBindingUserData)
		return
	}

	user.Role = utils.RoleUser

	err = user.FindByEmail()
	if err == nil {
		ctx.JSON(http.StatusConflict, errors.ErrUserAlreadyExists)
		return
	}

	err = user.Save()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	token, err := utils.GenerateToken(user.Email, user.Username, user.ID, time.Hour*2, utils.Role(user.Role))

	if err != nil {
		ctx.JSON(http.StatusOK, errors.ErrSomethingWentWrong)
		return
	}
	_, err = mailS.SendSimpleMessage(os.Getenv("API_URL")+"user/verify?token="+token.AccessToken, user.Email, user.Username, "d-958c75cdb588424fb80e49688fb2c3da")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.ErrSendingVerificationEmail)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"user": user})
}

func login(ctx *gin.Context) {

	var user *models.User
	err := ctx.ShouldBindJSON(&user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.ErrBindingUserData)
		return
	}

	err = user.ValidateCredintials(user.Password)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, err.Error())
		return
	}
	td, err := utils.GenerateToken(user.Email, user.Username, user.ID, time.Hour*2, utils.Role(user.Role))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.ErrGeneratingToken)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"user": user, "token": td.AccessToken, "refreshToken": td.RefreshToken})
}

func verifyEmail(ctx *gin.Context) {
	token, ok := ctx.GetQuery("token")

	if !ok {
		ctx.JSON(http.StatusInternalServerError, errors.ErrIncompleteRequest)
	}

	userId, _, err := utils.VerifyToken(token)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.ErrIncorrectOrExpiredToken)
	}

	var user models.User

	user.ID = userId

	err = user.Update(bson.M{"isVerified": true})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "oops!" + err.Error()})

	}
	// ctx.JSON(http.StatusOK, gin.H{"message": "user is verified."})
	ctx.Redirect(http.StatusTemporaryRedirect, "https://app.mux04.com/auth/login")
}

func resetPassword(ctx *gin.Context) {
	var user models.User
	err := ctx.ShouldBindJSON(&user)
	if user.Email == "" {
		ctx.JSON(http.StatusInternalServerError, errors.ErrEmptyEmail)
		return
	}

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong (643)! " + err.Error()})
		return
	}

	err = user.FindByEmail()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.ErrFindingByEmail)
		return
	}

	token, err := utils.GenerateToken(user.Email, user.Username, user.ID, time.Minute*15, utils.Role(user.Role))

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong (732)!"})
		return
	}
	_, err = mailS.SendSimpleMessage("https://mux04.com/auth/reset-password/confirm/"+token.AccessToken, user.Email, user.Username, "d-325e3a95b2fb497d9c293519596f6a45")

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.ErrSendingResetPasswordEmail)
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Reset password email sent successfully! please check your inbox."})
}

func resetPasswordVerify(ctx *gin.Context) {
	token := ctx.Query("token")
	var user models.User
	if token == "" {
		ctx.JSON(http.StatusInternalServerError, errors.ErrSomethingWentWrong)
		return
	}

	var passwords struct {
		OldPassword string `binding:"required"`
		NewPassword string `binding:"required"`
	}

	err := ctx.ShouldBindJSON(&passwords)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.ErrParsing)
		return
	}

	userid, _, err := utils.VerifyToken(token)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errors.ErrIncorrectOrExpiredToken)
		return
	}

	user.ID = userid

	err = user.VerifyAndUpdatePassword(passwords.OldPassword)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errors.ErrVerifyingAndUpdatePassword)
		return
	}

	hashedPassword, err := utils.HashPassword(passwords.NewPassword)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error!" + err.Error()})
		return
	}

	err = user.Update(bson.M{"password": hashedPassword})

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error!" + err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{"message": ""})
}

func checkMongoDBConnection(ctx *gin.Context) {
	_, _, err := db.Init()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to connect to MongoDB//" + os.Getenv("MONGO_URL")})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Successfully connected to MongoDB"})
}

func setAdminRole(ctx *gin.Context) {
	// currentRole, exists := ctx.Get("role")
	// if !exists || currentRole != utils.RoleAdmin {
	// 	ctx.JSON(http.StatusForbidden, gin.H{"error": "Are you serious? (:"})
	// 	return
	// }

	var requestBody struct {
		UserID string `json:"userId" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	userID, err := primitive.ObjectIDFromHex(requestBody.UserID)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var user models.User
	user.ID = userID

	err = user.Update(bson.M{"role": string(utils.RoleAdmin)})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user role"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User role updated to admin successfully"})
}
