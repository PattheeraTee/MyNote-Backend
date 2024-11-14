package httpHandler

import (
	"time"
	"miw/entities"
	"miw/usecases/service"
	"strconv"
	"github.com/gofiber/fiber/v2"
)

type HttpUserHandler struct {
	userUseCase service.UserUseCase
}

func NewHttpUserHandler(useCase service.UserUseCase) *HttpUserHandler {
	return &HttpUserHandler{userUseCase: useCase}
}

func (h *HttpUserHandler) Register(c *fiber.Ctx) error {
	user := new(entities.User)
	if err := c.BodyParser(user); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if err := h.userUseCase.Register(user); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Could not register user")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "User registered successfully"})
}

func (h *HttpUserHandler) Login(c *fiber.Ctx) error {
	data := new(struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	})
	if err := c.BodyParser(data); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	token, err := h.userUseCase.Login(data.Email, data.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("Email or password is incorrect")
	}

	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 72),
		HTTPOnly: true,
	})

	return c.JSON(fiber.Map{"message": "Login successful"})
}

func (h *HttpUserHandler) ForgotPassword(c *fiber.Ctx) error {
	data := new(struct {
		Email string `json:"email"`
	})
	if err := c.BodyParser(data); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if err := h.userUseCase.SendResetPasswordEmail(data.Email); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Could not send reset email")
	}

	return c.JSON(fiber.Map{"message": "Reset password email sent"})
}

func (h *HttpUserHandler) ResetPassword(c *fiber.Ctx) error {
	// รับโทเค็นจาก query parameter
	tokenString := c.Query("token")

	// อ่านรหัสผ่านใหม่จาก body และตรวจสอบว่ามีเพียงฟิลด์ password เท่านั้น
	var requestBody map[string]interface{}
	if err := c.BodyParser(&requestBody); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	// ตรวจสอบว่ามีฟิลด์อื่นนอกเหนือจาก password หรือไม่
	if password, exists := requestBody["password"]; exists && len(requestBody) == 1 {
		passwordStr := password.(string)

		// ใช้ userUseCase.ResetPassword เพื่อเปลี่ยนรหัสผ่าน
		if err := h.userUseCase.ResetPassword(tokenString, passwordStr); err != nil {
			if err.Error() == "user not found" {
				return c.Status(fiber.StatusNotFound).SendString("user not found")
			}
			return c.Status(fiber.StatusUnauthorized).SendString("could not reset password")
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "password updated successfully",
		})
	}

	return c.Status(fiber.StatusBadRequest).SendString("Only 'password' field is allowed")
}

func (h *HttpUserHandler) GetUser(c *fiber.Ctx) error {
	id,err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid ID")
	}
	user, err := h.userUseCase.GetUser(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("User not found")
	}

	return c.JSON(user)
}

func (h *HttpUserHandler) ChangeUsername(c *fiber.Ctx) error {
	id,err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid ID")
	}

	var requestBody map[string]interface{}
	if err := c.BodyParser(&requestBody); err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
	}

	if newUsername, exists := requestBody["username"]; exists && len(requestBody) == 1 {
		newUsernameStr := newUsername.(string)

		if err := h.userUseCase.ChangeUsername(uint(id), newUsernameStr); err != nil {
			if err.Error() == "user not found" {
				return c.Status(fiber.StatusNotFound).SendString("user not found")
			}
			return c.Status(fiber.StatusInternalServerError).SendString("could not change username")
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "username updated successfully",
		})
	}

	return c.Status(fiber.StatusBadRequest).SendString("Only 'username' field is allowed")
}
