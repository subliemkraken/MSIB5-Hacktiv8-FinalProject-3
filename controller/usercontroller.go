package controller

import (
	"FinalProject3/config"
	"FinalProject3/dto"
	"FinalProject3/middleware"
	"FinalProject3/model"
	"FinalProject3/model/validation"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(c *gin.Context) {
	var user model.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	user.Role = "member"

	if err := validation.ValidateStruct(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if config.DB.Where("email = ?", user.Email).Find(&user).RowsAffected > 0 {
		c.JSON(http.StatusConflict, gin.H{"message": "The email is already in use"})
		return
	}

	if err := user.HashPassword(); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to hash password"})
		return
	}

	if err := config.DB.Create(&user).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to create user"})
		return
	}

	userdto := dto.User{
		ID:        user.ID,
		Full_name: user.Full_name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	c.JSON(http.StatusCreated, userdto)
}

func LoginUser(c *gin.Context) {

	var body struct {
		Email    string
		Password string
	}

	if err := c.Bind(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	var user model.User

	config.DB.Find(&user, "email = ?", body.Email)

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Incorrect email or password"})
		return
	}

	signedToken, err := middleware.GenerateToken(user.ID, user.Email, user.Role)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Incorrect email or password"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": signedToken,
	})
}

func UpdateUser(c *gin.Context) {
	userData, _ := c.Get("userData")
	userClaim, ok := userData.(jwt.MapClaims)

	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
		return
	}

	userID := uint(userClaim["id"].(float64))

	var user model.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if err := validation.ValidateStruct(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	if config.DB.Model(&user).Where("id = ?", userID).Updates(&user).RowsAffected == 0 {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": "Failed to update user"})
		return
	}

	userdto := dto.User{
		ID:        userID,
		Full_name: user.Full_name,
		Email:     user.Email,
		CreatedAt: user.UpdatedAt,
	}

	c.JSON(http.StatusOK, userdto)
}

func DeleteUser(c *gin.Context) {
	userData, _ := c.Get("userData")
	userClaim, ok := userData.(jwt.MapClaims)

	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
		return
	}

	userID := uint(userClaim["id"].(float64))

	var user model.User

	if err := config.DB.Delete(&user, userID).Error; err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Your account has been succesfully deleted",
	})
}
