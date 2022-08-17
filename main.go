package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jupripratama/jprstartup/auth"
	"github.com/jupripratama/jprstartup/campaign"
	"github.com/jupripratama/jprstartup/handler"
	"github.com/jupripratama/jprstartup/helper"
	"github.com/jupripratama/jprstartup/user"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := "root:password@tcp(127.0.0.1:3306)/bwastartup?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal(err.Error())
	}

	userRepository := user.NewRepository(db)
	campaignRepository := campaign.NewRepository(db)

	userService := user.NewService(userRepository)
	//buatcek kambar manual  gambar
	// userService.SaveAvatar(1, "image/1-profile.png")
	/* tes  service
	input := user.LoginInput{
		Email:    "dsa@gmail.com",
		Password: "passwoerd",
	}

	user, err := userService.Login(input)
	if err != nil {
		fmt.Println("salah  broo")
		fmt.Println(err.Error())
	}
	fmt.Println(user.Email)
	fmt.Println(user.Name)
	*/
	/* mencoba tes cari user
	userByEmai, err := userRepository.FindByEmail("dsa@gmail.com")
	if err != nil {
		fmt.Println(err.Error())
	}
	// fmt.Println(userByEmai.Name)
	if userByEmai.ID == 0 {
		fmt.Println("User Tidak ketemu")
	} else {
		fmt.Println(userByEmai.Name)
	}
	*/

	campaignService := campaign.NewService(campaignRepository)

	// tes campaigns
	// campaigns, err := campaignRepository.FindByUserID(1)
	// fmt.Println("debug")
	// fmt.Println("debug")
	// fmt.Println("debug")
	// fmt.Println(len(campaigns))
	// for _, campaign := range campaigns {
	// 	fmt.Println(campaign.Name)
	// 	if len(campaign.CampaignImages) > 0 {
	// 		fmt.Println(campaign.CampaignImages[0].FileName)
	// 	}
	// }

	campaigns, _ := campaignService.FindCampaign(2)
	fmt.Println(len(campaigns))
	authService := auth.NewService()
	/* tes GenerateToken Jwt
	fmt.Println(authService.GenerateToken(1001))
	*/
	/* tes ValidateToken  Jwt
	token, err := authService.ValidateToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoyfQ.xFy841awdutKeIZEM4j-DBWcHAz-6dxUPMPBLkpSbbI")
	if err != nil {
		fmt.Println("EROR")

	}
	if token.Valid {
		fmt.Println("VALID")

	} else {
		fmt.Println("Invalid token")

	}
	*/

	userHandler := handler.NewUserHandler(userService, authService)

	router := gin.Default()
	api := router.Group("/api/v1")

	api.POST("/users", userHandler.RegisterUser)
	api.POST("/sessions", userHandler.Login)
	api.POST("/email_checkers", userHandler.ChekEmailAvailability)
	api.POST("/avatars", authMiddleware(authService, userService), userHandler.UploadAvatar)

	router.Run(":8080")

	//instput dari user
	//handler, mapping input dari user -> struct input
	//service : mellakukan mapping dari struct input ke struct user
	//repository '
	//db

}

func authMiddleware(authService auth.Service, userService user.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}
		//Bearer tokentoken
		// 0           1
		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		token, err := authService.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claim, ok := token.Claims.(jwt.MapClaims)

		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		userID := int(claim["user_id"].(float64))

		user, err := userService.GetUserByID(userID)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			c.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		c.Set("currentUser", user)
		// ambil nilai header Authorization: Bearer tokentoken
		// dari Authorization, kita ambil nilai tokennya saja
		// kita validasi tokennya
		// kita ambil user_id
		// ambil user dari db berdasarkan user_id lewat service
		// kita set context isinya user

	}
}
