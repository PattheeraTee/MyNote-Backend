package main

import (
	"fmt"
	"log"
	"miw/adapters/gormRepository"
	"miw/adapters/httpHandler"
	"miw/database"
	"miw/entities"
	"miw/middleware"
	"miw/usecases/service"
	"github.com/gofiber/fiber/v2"
)

func main() {
	cfg := database.LoadConfig()
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort)

	database, err := database.NewDatabaseConnection(dsn)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	// สร้างตารางอัตโนมัติโดยใช้ AutoMigrate
	err = database.AutoMigrate(
		&entities.User{},
		&entities.Note{},
		&entities.Reminder{},
		&entities.Tag{},
		// &entities.NoteTag{},
		&entities.ShareNote{},
		&entities.Event{},
		&entities.ToDo{},
	)

	if err != nil {
		log.Fatal("Failed to migrate tables:", err)
	}

	// สร้าง Repository และ Service
	userRepo := gormRepository.NewGormUserRepository(database)
	noteRepo := gormRepository.NewGormNoteRepository(database)
	tagRepo := gormRepository.NewGormTagRepository(database)

	userService := service.NewUserService(userRepo)
	noteService := service.NewNoteService(noteRepo)
	tagService := service.NewTagService(tagRepo)

	// สร้าง Handlers สำหรับ HTTP
	userHandler := httpHandler.NewHttpUserHandler(userService)
	noteHandler := httpHandler.NewHttpNoteHandler(noteService)
	tagHandler := httpHandler.NewHttpTagHandler(tagService)

	// สร้าง Fiber App และเพิ่ม Middleware
	app := fiber.New()

	//********************************************
	// User
	//********************************************
	app.Post("/register", userHandler.Register)
	app.Post("/login", userHandler.Login)

	app.Post("/forgot-password", userHandler.ForgotPassword)
	app.Post("/reset-password", userHandler.ResetPassword)

	app.Get("/user/:userid", middleware.AuthMiddleware, userHandler.GetUser)        // ดูข้อมูล user
	app.Put("/user/:userid", middleware.AuthMiddleware, userHandler.ChangeUsername) // แก้ไข username

	//********************************************
	// Note
	//********************************************
	app.Post("/note",middleware.AuthMiddleware, noteHandler.CreateNoteHandler)    // สร้าง note	
	app.Get("/note/:userid",middleware.AuthMiddleware, noteHandler.GetAllNoteHandler) // ดู note
	app.Put("/note/:noteid", middleware.AuthMiddleware,noteHandler.UpdateNoteHandler) // แก้ไข note
	app.Post("/note/add-tag",middleware.AuthMiddleware, noteHandler.AddTagToNoteHandler)
	app.Post("/note/remove-tag",middleware.AuthMiddleware,  noteHandler.RemoveTagFromNoteHandler)
	app.Delete("/note/:noteid",middleware.AuthMiddleware, noteHandler.DeleteNoteHandler) // ลบ note
	app.Put("/note/restore/:noteid",middleware.AuthMiddleware, noteHandler.RestoreNoteHandler)

	//********************************************
	// Tag
	//********************************************
	app.Post("/tag", middleware.AuthMiddleware, tagHandler.CreateTagHandler) // สร้าง tag
	app.Get("/tag/:tagid",middleware.AuthMiddleware,  tagHandler.GetTagHandler) // ดู tag
	app.Put("/tag/:tagid", middleware.AuthMiddleware, tagHandler.UpdateTagNameHandler) // แก้ไขชื่อ tag
	app.Delete("/tag/:tagid", middleware.AuthMiddleware, tagHandler.DeleteTagHandler) // ลบ tag
	
	// เริ่มเซิร์ฟเวอร์
	if err := app.Listen(":8000"); err != nil {
		log.Fatal("Failed to start server:", err)
	}

}

