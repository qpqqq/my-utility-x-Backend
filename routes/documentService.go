package routes

import (
	"encoding/base64"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"myutilityx.com/errors"
	grpc "myutilityx.com/gRPC"
	"myutilityx.com/models"
	"myutilityx.com/repository"
)

func getFile(ctx *gin.Context) {
	fileId := ctx.Param("fileId")

	req := &grpc.GetFileRequest{
		Filename: "testFileFromGo.pdf",
		FileId:   fileId,
	}

	client := grpc.Connect()
	response, err := client.GetFile(ctx, req)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve file" + err.Error()})
		return
	}

	decodedFile, err := base64.StdEncoding.DecodeString(response.File)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode file content"})
		return
	}
	ctx.Data(http.StatusOK, "application/octet-stream", decodedFile)
}

// Add this function to close the gRPC connection when your server shuts down
func CloseGRPCConnection() {
	grpc.Close()
}

func uploadFile(ctx *gin.Context) {
	userId, exist := ctx.Get("userId")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, errors.ErrUnAuthorized)
		return
	}

	var requestBody struct {
		Filename string `json:"filename"`
		Filedata string `json:"filedata"`
	}
	if err := ctx.ShouldBindJSON(&requestBody); err != nil {
		log.Printf("Error binding JSON: %v", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	// Extract the base64 data from the filedata string
	parts := strings.SplitN(requestBody.Filedata, ",", 2)
	if len(parts) != 2 {
		log.Printf("Invalid file data format")
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file data format"})
		return
	}
	base64Data := parts[1]

	req := &grpc.UploadFileRequest{
		Filename: requestBody.Filename,
		FileData: base64Data,
	}

	client := grpc.Connect()
	response, err := client.UploadFile(ctx, req)
	if err != nil {
		log.Printf("gRPC UploadFile error: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file: " + err.Error()})
		return
	}

	fileRepo := repository.NewFileRepository()
	file := models.File{
		FileName: req.Filename,
		FileId:   response.FileId,
		UserId:   userId.(primitive.ObjectID),
	}

	if err := fileRepo.AddFile(ctx, file); err != nil {
		log.Printf("Error adding file to database: %v", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file information"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully"})
}

func getFiles(ctx *gin.Context) {
	userId, exist := ctx.Get("userId")
	fileRepo := repository.NewFileRepository()

	if !exist {
		ctx.JSON(401, errors.ErrUnAuthorized)
	}

	if userId, ok := userId.(primitive.ObjectID); ok {
		result, err := fileRepo.GetAll(ctx, userId)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errors.ErrSomethingWentWrong)
		}
		ctx.JSON(http.StatusOK, result)

	}
}

func deleteFile(ctx *gin.Context) {
	userId, exist := ctx.Get("userId")
	fileId := ctx.Param("fileId")
	fileRepo := repository.NewFileRepository()

	if !exist {
		ctx.JSON(401, errors.ErrUnAuthorized)
		return
	}

	if userId, ok := userId.(primitive.ObjectID); ok {
		fileRepo.DeleteFile(ctx, fileId, userId)
	}

	req := &grpc.DeleteFileRequest{
		FileId: fileId,
	}

	client := grpc.Connect()
	res, err := client.DeleteFile(ctx, req)
	if err != nil || !res.Ok {
		ctx.JSON(400, errors.ErrSomethingWentWrong)
		return
	}

}
