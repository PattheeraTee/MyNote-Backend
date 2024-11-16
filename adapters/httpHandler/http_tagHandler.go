package httpHandler

import (
	"miw/entities"
	"miw/usecases/service"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type TagResponse struct {
	TagID   uint   `json:"tag_id"`
	TagName string `json:"tag_name"`
	Notes   []uint `json:"notes"`
}

type HttpTagHandler struct {
	tagUseCase service.TagUseCase
}

func NewHttpTagHandler(useCase service.TagUseCase) *HttpTagHandler {
	return &HttpTagHandler{tagUseCase: useCase}
}

func (h *HttpTagHandler) CreateTagHandler(c *fiber.Ctx) error {
	tag := new(entities.Tag)

	// รับข้อมูลแท็กจาก Body ของ request
	if err := c.BodyParser(tag); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid request body")
	}

	// เรียกใช้ฟังก์ชันสร้างแท็ก
	if err := h.tagUseCase.CreateTag(tag); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Could not create tag")
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Tag created successfully",
		"tag":     tag,
	})
}

func (h *HttpTagHandler) GetTagHandler(c *fiber.Ctx) error {
	tagID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid tag ID"})
	}

	tag, err := h.tagUseCase.GetTagById(uint(tagID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	// สร้างโครงสร้าง JSON ใหม่ที่แสดงเฉพาะ note_id
	var noteIDs []uint
	for _, note := range tag.Notes {
		noteIDs = append(noteIDs, note.NoteID)
	}

	response := TagResponse{
		TagID:   tag.TagID,
		TagName: tag.TagName,
		Notes:   noteIDs,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
