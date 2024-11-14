package main

import (
	"fmt"
	"miw/entities"
	"miw/handler"
	"os"
	"strconv"
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"gopkg.in/gomail.v2"
	"gorm.io/gorm"
)

func registerHandler(db *gorm.DB, c *fiber.Ctx) error {
	user := new(entities.User)

	if err := c.BodyParser(user); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	err := handler.CreateUser(db, user)

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "user created successfully",
	})
}

func loginHandler(db *gorm.DB, c *fiber.Ctx) error {
	user := new(entities.User)

	if err := c.BodyParser(user); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	token, err := handler.Login(db, user)

	if err != nil {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 72),
		HTTPOnly: true,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "login successful",
	})
}

func getUserHandler(db *gorm.DB, c *fiber.Ctx) error {
	// ดึง user_id จาก Context
	authUserID := c.Locals("user_id").(uint)

	// ดึง user_id จาก URL และแปลงเป็น uint
	id, err := strconv.Atoi(c.Params("id"))

	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	requestedUserID := uint(id)

	// ตรวจสอบว่า authUserID ตรงกับ requestedUserID หรือไม่
	if authUserID != requestedUserID {
		// ห้ามเข้าถึงหากไม่ใช่ข้อมูลของตัวเอง
		return c.Status(fiber.StatusForbidden).SendString("You are not authorized to access this user's information.")
	}

	user := handler.GetUser(db, uint(id))

	if user == nil {
		return c.Status(fiber.StatusNotFound).SendString("user not found")
	}

	return c.JSON(user)
}

func changeUsernameHandler(db *gorm.DB, c *fiber.Ctx) error {
	// ดึง user_id จาก Context
	authUserID := c.Locals("user_id").(uint)

	// ดึง user_id จาก URL และแปลงเป็น uint
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid user ID")
	}
	requestedUserID := uint(id)

	// ตรวจสอบว่า authUserID ตรงกับ requestedUserID หรือไม่
	if authUserID != requestedUserID {
		// ห้ามเข้าถึงหากไม่ใช่ข้อมูลของตัวเอง
		return c.Status(fiber.StatusForbidden).SendString("You are not authorized to access this user's information.")
	}

	// อ่านข้อมูล username ใหม่จาก body
	var requestBody map[string]interface{}
	if err := c.BodyParser(&requestBody); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	// ตรวจสอบว่ามีฟิลด์อื่น ๆ นอกเหนือจาก username หรือไม่
	if username, exists := requestBody["username"]; exists && len(requestBody) == 1 {
		user := new(entities.User)
		user.UserID = requestedUserID
		user.Username = username.(string)

		// เรียกใช้งานฟังก์ชันเพื่ออัปเดต user
		result := handler.UpdateUsername(db, user)
		if result != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("could not update username")
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "username updated successfully",
		})
	}

	return c.Status(fiber.StatusBadRequest).SendString("Only 'username' field is allowed")
}


func forgotPasswordHandler(db *gorm.DB, c *fiber.Ctx) error {
	// อ่าน email จาก body โดยใช้ map[string]interface{} เพื่อเช็คจำนวนฟิลด์
	var requestBody map[string]interface{}
	if err := c.BodyParser(&requestBody); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	// ตรวจสอบว่ามีฟิลด์อื่นนอกเหนือจาก email หรือไม่
	if email, exists := requestBody["email"]; exists && len(requestBody) == 1 {
		emailStr := email.(string)
		// ค้นหา user ในระบบตามอีเมลที่ส่งมา
		user := new(entities.User)
		if err := db.Where("email = ?", emailStr).First(user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusNotFound).SendString("email not found")
			}
			return c.Status(fiber.StatusInternalServerError).SendString("could not find user")
		}

		// สร้าง JWT token สำหรับรีเซ็ตรหัสผ่าน
		jwtSecret := os.Getenv("JWT_SECRET")
		token := jwt.New(jwt.SigningMethodHS256)
		claims := token.Claims.(jwt.MapClaims)
		claims["user_id"] = user.UserID
		claims["email"] = user.Email
		claims["exp"] = time.Now().Add(time.Hour * 1).Unix() // token หมดอายุใน 1 ชั่วโมง

		resetToken, err := token.SignedString([]byte(jwtSecret))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("could not create token")
		}

		// ส่งอีเมลที่มีลิงก์สำหรับรีเซ็ตรหัสผ่าน
		resetURL := fmt.Sprintf("http://localhost:8000/reset-password?token=%s", resetToken)
		if err := sendResetEmail(user.Email, resetURL); err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("could not send email")
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "password reset link sent",
		})
	} else {
		return c.Status(fiber.StatusBadRequest).SendString("Only 'email' field is allowed")
	}
}

// ฟังก์ชันส่งอีเมลรีเซ็ตรหัสผ่าน
func sendResetEmail(email, resetURL string) error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", "pattheera.t@kkumail.com")
	mailer.SetHeader("To", email)
	mailer.SetHeader("Subject", "Password Reset Request")
	mailer.SetBody("text/plain", fmt.Sprintf("Click here to reset your password: %s", resetURL))

	dialer := gomail.NewDialer("smtp.gmail.com", 587, os.Getenv("MAIL_EMAIL"), os.Getenv("MAIL_PASSWORD"))

	return dialer.DialAndSend(mailer)
}

func resetPasswordHandler(db *gorm.DB, c *fiber.Ctx) error {
	// รับ token จาก query parameter
	tokenString := c.Query("token")
	jwtSecret := os.Getenv("JWT_SECRET")

	// ตรวจสอบ token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	claims := token.Claims.(jwt.MapClaims)
	email := claims["email"].(string)

	// รับรหัสผ่านใหม่จาก body และตรวจสอบว่ามีเพียงฟิลด์ password เท่านั้น
	var requestBody map[string]interface{}
	if err := c.BodyParser(&requestBody); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	// ตรวจสอบว่ามีฟิลด์อื่นนอกเหนือจาก password หรือไม่
	if password, exists := requestBody["password"]; exists && len(requestBody) == 1 {
		passwordStr := password.(string)

		// ค้นหาผู้ใช้ด้วยอีเมลที่ตรงกับโทเค็น
		user := new(entities.User)
		if err := db.Where("email = ?", email).First(user).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return c.Status(fiber.StatusNotFound).SendString("user not found")
			}
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		// อัปเดตรหัสผ่านใหม่
		user.Password = passwordStr
		if err := handler.UpdatePassword(db, user); err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "password updated successfully",
		})
	} else {
		return c.Status(fiber.StatusBadRequest).SendString("Only 'password' field is allowed")
	}
}
