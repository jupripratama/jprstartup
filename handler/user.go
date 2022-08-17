package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jupripratama/jprstartup/auth"
	"github.com/jupripratama/jprstartup/helper"
	"github.com/jupripratama/jprstartup/user"
)

type userHandler struct {
	userService user.Service
	authService auth.Service
}

func NewUserHandler(userService user.Service, authService auth.Service) *userHandler {
	return &userHandler{userService, authService}
}

func (h *userHandler) RegisterUser(c *gin.Context) {
	var input user.RegisterUserInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		// pindah ke kelper
		// var errors []string

		// for _, e := range err.(validator.ValidationErrors) {
		// 	errors = append(errors, e.Error())
		// }
		// berubah menjadi
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"errors": errors}

		response := helper.APIResponse("Account has been failed", http.StatusUnprocessableEntity, "not successfully", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": false,
			"user":    response,
		})
		return
	}

	newUser, err := h.userService.RegisterUser(input)
	if err != nil {
		response := helper.APIResponse("Account has been failed", http.StatusBadRequest, "not successfully", nil)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": false,
			"user":    response,
		})
		return
	}

	// token, e
	token, err := h.authService.GenerateToken(newUser.ID)
	if err != nil {
		response := helper.APIResponse("Token generation has been failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": false,
			"user":    response,
		})
		return
	}

	formatter := user.FormatUser(newUser, token)

	response := helper.APIResponse("Account has been registered successfully", http.StatusOK, "successfully", formatter)

	c.JSON(http.StatusOK, gin.H{
		"message": true,
		"user":    response,
	})

	// tangkap input user
	// map input dari user ke struct RegisterUser
	// struct di atas kita passing sepabagai parameter service
}

func (h *userHandler) Login(c *gin.Context) {
	// user masukan input (email dan password)
	// input dari handler

	var input user.LoginInput

	err := c.ShouldBind(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"error": errors}

		response := helper.APIResponse("Login has been failed", http.StatusUnprocessableEntity, "not successfully", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": false,
			"user":    response,
		})
		return
	}

	loggedinUser, err := h.userService.Login(input)

	if err != nil {
		errorMessage := gin.H{"error": err.Error()}
		response := helper.APIResponse("Login has been failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": false,
			"user":    response,
		})
		return
	}

	token, err := h.authService.GenerateToken(loggedinUser.ID)
	if err != nil {
		response := helper.APIResponse("Login Token generation has been failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": false,
			"user":    response,
		})
		return
	}

	formatter := user.FormatUser(loggedinUser, token)
	response := helper.APIResponse("Successfully loggedin input", http.StatusOK, "successfully", formatter)
	c.JSON(http.StatusOK, response)

	// map input dari user ke input struct
	// input struct passing service
	// di service mencari dg bantuan repository user dengan email  x

	// mencocokan password
}

func (h *userHandler) ChekEmailAvailability(c *gin.Context) {
	// ada input email dari user
	// input email di-maping ke struct input
	// struct input di-passing ke service
	// service akan memangil  repository - email sudah ada atau belum
	// repository - db
	var input user.CekEmailInput

	err := c.ShouldBind(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMessage := gin.H{"error": errors}

		response := helper.APIResponse("Chek Email Failed", http.StatusUnprocessableEntity, "errorMessage", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": false,
			"user":    response,
		})
		return
	}
	isEmailAvailable, err := h.userService.IsEmailAvailable(input)
	if err != nil {
		errorMessage := gin.H{"error": "server error"}
		response := helper.APIResponse("Email cheking request failed", http.StatusUnprocessableEntity, "error", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": false,
			"user":    response,
		})
		return
	}
	data := gin.H{
		"is_available": isEmailAvailable,
	}

	metaMesage := "Email has been registered"

	if isEmailAvailable {
		metaMesage = "Email is available"
	}

	response := helper.APIResponse(metaMesage, http.StatusOK, "successfully", data)
	c.JSON(http.StatusOK, gin.H{
		"message": true,
		"user":    response,
	})

}

func (h *userHandler) UploadAvatar(c *gin.Context) {
	// input gambar user
	//simpan gambar difolder image
	//di  service kita pangil repo
	//jwt
	// repo ambil data user yang id = 1
	// repo update data simpan lokasi file
	file, err := c.FormFile("avatar")
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("failed to create avatar file", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": false,
			"data":    response,
		})
		return
	}

	currentUser := c.MustGet("currentUser").(user.User)
	userID := currentUser.ID
	//path := "images/avatars/" + file.Filename
	path := fmt.Sprintf("images/avatars/%d-%s", userID, file.Filename)

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("failed to upload file image avatar", http.StatusBadRequest, "error Save Upload  File", data)

		c.JSON(http.StatusBadRequest, gin.H{
			"mesage": false,
			"data":   response,
		})
		return
	}
	_, err = h.userService.SaveAvatar(userID, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIResponse("failed to upload file image avatar", http.StatusBadRequest, "error Save Avatar", data)

		c.JSON(http.StatusBadRequest, gin.H{
			"mesage": false,
			"data":   response,
		})
		return
	}
	data := gin.H{"is_uploaded": true}
	response := helper.APIResponse("succeed to upload file image avatar", http.StatusOK, "succeed", data)
	c.JSON(http.StatusOK, gin.H{
		"mesage": true,
		"data":   response,
	})
}
