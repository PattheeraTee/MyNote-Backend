package httpHandler

import (
	"miw/entities"
	"miw/usecases/service"
	"strconv"
	"github.com/gofiber/fiber/v2"
	"fmt"
)

type NoteResponse struct {
	NoteID     uint          `json:"note_id"`
	UserID     uint          `json:"user_id"`
	Title      string        `json:"title"`
	Content    string        `json:"content"`
	Color      string        `json:"color"`
	Priority   int           `json:"priority"`
	IsTodo     bool          `json:"is_todo"`
	TodoStatus bool          `json:"todo_status"`
	CreatedAt  string        `json:"created_at"`
	UpdatedAt  string        `json:"updated_at"`
	DeletedAt  string       `json:"deleted_at"`
	Tags       []string      `json:"Tags"`
	Reminders []entities.Reminder `json:"Reminders"`
	Event      interface{}   `json:"Event"`
}
type ReminderResponse struct {
	ReminderID uint   `json:"reminder_id"`
	NoteID     uint   `json:"note_id"`
	Content    string `json:"content"`
	DateTime   string `json:"datetime"`
}


type HttpNoteHandler struct {
	noteUseCase service.NoteUseCase
}

func NewHttpNoteHandler(useCase service.NoteUseCase) *HttpNoteHandler {
	return &HttpNoteHandler{noteUseCase: useCase}
}

func (h *HttpNoteHandler) CreateNoteHandler(c *fiber.Ctx) error {
	note := new(entities.Note)

	// รับข้อมูลโน้ตจาก Body ของ request
	if err := c.BodyParser(note); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	// เรียกใช้ฟังก์ชันสร้างโน้ต
	if err := h.noteUseCase.CreateNote(note); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Could not create note")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Note created successfully",
		"note":    note,
	})
}

func (h *HttpNoteHandler) GetAllNoteHandler(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid user id")
	}

	notes, err := h.noteUseCase.GetAllNote(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Notes not found for this user")
	}

	// สร้าง response list โดยใช้ struct NoteResponse
	var response []NoteResponse
	for _, note := range notes {
		tags := []string{}
		for _, tag := range note.Tags {
			tags = append(tags, tag.TagName)
		}

		response = append(response, NoteResponse{
			NoteID:     note.NoteID,
			UserID:     note.UserID,
			Title:      note.Title,
			Content:    note.Content,
			Color:      note.Color,
			Priority:   note.Priority,
			IsTodo:     note.IsTodo,
			TodoStatus: note.TodoStatus,
			CreatedAt:  note.CreatedAt,
			UpdatedAt:  note.UpdatedAt,
			DeletedAt:  note.DeletedAt,
			Tags:       tags,
			Reminders:  note.Reminders,
			Event:      note.Event,
		})
	}

	// ส่ง JSON response พร้อมการจัดเรียงลำดับที่ต้องการ
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"notes": response,
	})
}

func (h *HttpNoteHandler) UpdateNoteHandler(c *fiber.Ctx) error {
	// ดึง Note ID จาก URL พารามิเตอร์
	noteID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid note ID")
	}

	// รับข้อมูลโน้ตจาก Body ของ request
	note := new(entities.Note)
	if err := c.BodyParser(note); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	// เรียกใช้ฟังก์ชันแก้ไขโน้ต
	updatedNote, err := h.noteUseCase.UpdateNote(uint(noteID), note)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Could not update note")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Note updated successfully",
		"note":    updatedNote,
	})
}

func (h *HttpNoteHandler) AddTagToNoteHandler(c *fiber.Ctx) error {
    var request struct {
        NoteID uint `json:"note_id"`
        TagID  uint `json:"tag_id"`
    }

    // Parse JSON body into the request struct
    if err := c.BodyParser(&request); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
    }

    // Call the use case to add a single tag to the note
    err := h.noteUseCase.AddTagToNote(request.NoteID, request.TagID)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": fmt.Sprintf("Failed to add tag %d to note %d: %v", request.TagID, request.NoteID, err),
        })
    }

    // Return success message
    return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Tag added successfully"})
}

