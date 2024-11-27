package httpHandler

import (
	"miw/entities"
	"miw/usecases/service"
	"strconv"
	"github.com/gofiber/fiber/v2"
	"fmt"
	// "strings"
	// "time"
)

type NoteResponse struct {
	NoteID     uint                `json:"note_id"`
	UserID     uint                `json:"user_id"`
	Title      string              `json:"title"`
	Content    string              `json:"content,omitempty"` // ซ่อนถ้าไม่มีค่า
	Color      string              `json:"color"`
	Priority   int                 `json:"priority"`
	IsTodo     bool                `json:"is_todo"`
	IsAllDone  bool                `json:"is_all_done"`       // เพิ่มฟิลด์นี้
	TodoItems  []ToDoResponse      `json:"todo_items"`        // เพิ่มรายการ ToDo
	CreatedAt  string              `json:"created_at"`
	UpdatedAt  string              `json:"updated_at"`
	DeletedAt  string              `json:"deleted_at,omitempty"` // ซ่อนถ้าไม่มีค่า
	Tags       []string            `json:"tags"`
	Reminder   []entities.Reminder `json:"reminder"`
	Event      interface{}         `json:"event"`
}

type ReminderResponse struct {
	ReminderID uint   `json:"reminder_id"`
	NoteID     uint   `json:"note_id"`
	Content    string `json:"content"`
	DateTime   string `json:"datetime"`
}

type ToDoResponse struct {
	ID      uint   `json:"id"`
	Content string `json:"content"`
	IsDone  bool   `json:"is_done"`
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

	// ดึง UserID จาก Context (Middleware)
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}
	note.UserID = userID

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
	// ดึง UserID จาก Context (Middleware)
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	// ดึงข้อมูลโน้ตทั้งหมดของ User
	notes, err := h.noteUseCase.GetAllNote(userID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Notes not found for this user")
	}

	// แปลงผลลัพธ์เป็น JSON Response
	var response []NoteResponse
	for _, note := range notes {
		tags := []string{}
		for _, tag := range note.Tags {
			tags = append(tags, tag.TagName)
		}

		// แปลง TodoItems จาก entities.ToDo เป็น ToDoResponse
		var todoResponses []ToDoResponse
		for _, todo := range note.TodoItems {
			todoResponses = append(todoResponses, ToDoResponse{
				ID:      todo.ID,
				Content: todo.Content,
				IsDone:  todo.IsDone,
			})
		}

		response = append(response, NoteResponse{
			NoteID:     note.NoteID,
			UserID:     note.UserID,
			Title:      note.Title,
			Content:    note.Content,
			Color:      note.Color,
			Priority:   note.Priority,
			IsTodo:     note.IsTodo,
			IsAllDone:  note.IsAllDone,
			TodoItems:  todoResponses,
			CreatedAt:  note.CreatedAt,
			UpdatedAt:  note.UpdatedAt,
			DeletedAt:  note.DeletedAt,
			Tags:       tags,
			Reminder:   note.Reminder,
			Event:      note.Event,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"notes": response,
	})
}


func (h *HttpNoteHandler) UpdateNoteHandler(c *fiber.Ctx) error {
	noteID, err := strconv.Atoi(c.Params("noteid"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid note ID"})
	}

	note := new(entities.Note)
	if err := c.BodyParser(note); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	updatedNote, err := h.noteUseCase.UpdateNote(uint(noteID), note, userID)
	if err != nil {
		if err.Error() == "note not found or does not belong to the user" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "You are not authorized to update this note"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not update note"})
	}

	return c.JSON(updatedNote)
}


func (h *HttpNoteHandler) AddTagToNoteHandler(c *fiber.Ctx) error {
    var request struct {
        NoteID uint `json:"note_id"`
        TagID  uint `json:"tag_id"`
    }

    // รับข้อมูลจาก Body
    if err := c.BodyParser(&request); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
    }

    // ดึง UserID จาก Context
    userID, ok := c.Locals("user_id").(uint)
    if !ok {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
    }

    // เรียกใช้ Service Layer
    if err := h.noteUseCase.AddTagToNote(request.NoteID, request.TagID, userID); err != nil {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": err.Error()})
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Tag added successfully"})
}

func (h *HttpNoteHandler) RemoveTagFromNoteHandler(c *fiber.Ctx) error {
    var request struct {
        NoteID uint `json:"note_id"`
        TagID  uint `json:"tag_id"`
    }

    // Parse the request body
    if err := c.BodyParser(&request); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
    }

    // ดึง UserID จาก Context
    userID, ok := c.Locals("user_id").(uint)
    if !ok {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
    }

    // Call the use case to remove the tag from the note
    if err := h.noteUseCase.RemoveTagFromNote(request.NoteID, request.TagID, userID); err != nil {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "error": fmt.Sprintf("Failed to remove tag %d from note %d: %v", request.TagID, request.NoteID, err),
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Tag removed successfully"})
}


func (h *HttpNoteHandler) DeleteNoteHandler(c *fiber.Ctx) error {
    // ดึง Note ID จากพารามิเตอร์
    noteID, err := strconv.Atoi(c.Params("noteid"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid note ID"})
    }

    // ดึง UserID จาก Context
    userID, ok := c.Locals("user_id").(uint)
    if !ok {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
    }

    // เรียกใช้ Use Case เพื่อลบโน้ต
    if err := h.noteUseCase.DeleteNoteById(uint(noteID), userID); err != nil {
        if err.Error() == "note not found or does not belong to the user" {
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "You are not authorized to delete this note"})
        }
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": fmt.Sprintf("Failed to delete note with ID %d: %v", noteID, err),
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Note deleted successfully"})
}

func (h *HttpNoteHandler) RestoreNoteHandler(c *fiber.Ctx) error {
    // ดึง Note ID จากพารามิเตอร์
    noteID, err := strconv.Atoi(c.Params("noteid"))
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid note ID"})
    }

    // ดึง UserID จาก Context
    userID, ok := c.Locals("user_id").(uint)
    if !ok {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
    }

    // เรียกใช้ Use Case เพื่อกู้คืนโน้ต
    if err := h.noteUseCase.RestoreNoteById(uint(noteID), userID); err != nil {
        if err.Error() == "note not found or does not belong to the user" {
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "You are not authorized to restore this note"})
        }
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": fmt.Sprintf("Failed to restore note with ID %d: %v", noteID, err),
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Note restored successfully"})
}

