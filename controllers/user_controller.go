package controllers

import (
	"net/http"
	"time"
	"tusk/config"
	"tusk/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserController struct {
	DB *gorm.DB
}

type Claims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func (u *UserController) Login(c *gin.Context) {
	input := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user := models.User{}
	if err := u.DB.Where("email = ?", input.Email).Take(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{})
		return
	}
	errHash := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(input.Password),
	)
	if errHash != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email or Password is Wrong"})
		return
	}

	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		UserID: uint(user.Id),
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(config.JWTKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed generate token"})
		return
	}
	user.Token = tokenString

	c.JSON(http.StatusOK, user)

}

func (u *UserController) CreateAccount(c *gin.Context) {
	input := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
	}{}
	err := c.ShouldBindJSON(&input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	existingUser := models.User{}
	if err := u.DB.Where("email = ?", input.Email).
		First(&existingUser).Error; err == nil {

		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Email Already Exists",
		})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(input.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to hash password",
		})
		return
	}
	user := models.User{
		Email:    input.Email,
		Name:     input.Name,
		Password: string(hashedPassword),
		Role:     "Employee",
	}

	errDBCreate := u.DB.Create(&user).Error
	if errDBCreate != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errDBCreate.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (u *UserController) Delete(c *gin.Context) {
	id := c.Param("id")

	errDBDelete := u.DB.Delete(&models.User{}, id)
	if errDBDelete.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errDBDelete.Error.Error()})
		return
	}

	if errDBDelete.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Data not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message": "Deleted Successfully"})
}

func (u *UserController) GetEmployee(c *gin.Context) {

	users := []models.User{}

	result := u.DB.Select("id,name,email").Where("role = ?", "Employee").Find(&users)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"message": "Data Employee not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": users,
	})

}
