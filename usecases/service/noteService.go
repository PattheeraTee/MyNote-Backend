package service

import (
	"fmt"
	"miw/entities"
	"miw/usecases/repository"
	"time"
)

type NoteUseCase interface {
	CreateNote(note *entities.Note) error
	GetAllNote(userid uint) ([]entities.Note, error)
	UpdateColor(noteID uint, userID uint, color string) error
	UpdatePriority(noteID uint, userID uint, priority int) error
	UpdateTitleAndContent(noteID uint, userID uint, title string, content string, todoItems []entities.ToDo) error 
	UpdateStatus(noteID uint, userID uint, isTodo *bool, isAllDone *bool) error
	DeleteNoteById(noteID uint, userID uint) error
	RestoreNoteById(noteID uint, userID uint) error
	AddTagToNote(noteID uint, tagID uint, userID uint) error
	RemoveTagFromNote(noteID uint, tagID uint, userID uint) error
}

type NoteService struct {
	noteRepo repository.NoteRepository
}

func NewNoteService(noteRepo repository.NoteRepository) *NoteService {
	return &NoteService{
		noteRepo: noteRepo,
	}
}

func (s *NoteService) CreateNote(note *entities.Note) error {
	timeCreate := time.Now().Format("2006-01-02 15:04:05")
	note.CreatedAt = timeCreate

	// คำนวณ IsAllDone จาก TodoItems
	note.IsAllDone = true
	for _, todo := range note.TodoItems {
		if !todo.IsDone {
			note.IsAllDone = false
			break
		}
	}

	return s.noteRepo.CreateNote(note)
}

func (s *NoteService) GetAllNote(userid uint) ([]entities.Note, error) {
	return s.noteRepo.GetAllNoteByUserId(userid)
}

func (s *NoteService) UpdateColor(noteID uint, userID uint, color string) error {
	// ตรวจสอบว่า Note เป็นของ User หรือไม่
	_, err := s.noteRepo.GetNoteByIdAndUser(noteID, userID)
	if err != nil {
		return fmt.Errorf("note not found or does not belong to the user")
	}

	return s.noteRepo.UpdateNoteColor(noteID, userID, color)
}

func (s *NoteService) UpdatePriority(noteID uint, userID uint, priority int) error {
	// ตรวจสอบว่า Note เป็นของ User หรือไม่
	_, err := s.noteRepo.GetNoteByIdAndUser(noteID, userID)
	if err != nil {
		return fmt.Errorf("note not found or does not belong to the user")
	}

	return s.noteRepo.UpdateNotePriority(noteID, userID, priority)
}

func (s *NoteService) UpdateTitleAndContent(noteID uint, userID uint, title string, content string, todoItems []entities.ToDo) error {
	// ตรวจสอบว่า Note เป็นของ User หรือไม่
	note, err := s.noteRepo.GetNoteByIdAndUser(noteID, userID)
	if err != nil {
		return fmt.Errorf("note not found or does not belong to the user")
	}

	// Validation: ห้ามส่ง content และ todo_items พร้อมกัน
	if len(todoItems) > 0 && content != "" {
		return fmt.Errorf("note cannot have both content and todo_items")
	}

	// อัปเดต Title หากมีการส่งค่า
	if title != "" {
		note.Title = title
	}

	// ถ้ามี Content ให้ลบ TodoItems และอัปเดต Content
	if content != "" {
		note.Content = content
		note.TodoItems = nil // ลบ TodoItems
	}

	// ถ้ามี TodoItems ให้ลบ Content และอัปเดต TodoItems
	if len(todoItems) > 0 {
		note.TodoItems = todoItems
		note.Content = "" // ลบ Content
	}

	// อัปเดต UpdatedAt
	note.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")

	// บันทึกการอัปเดต
	return s.noteRepo.UpdateNoteTitleAndContent(note)
}


func (s *NoteService) UpdateStatus(noteID uint, userID uint, isTodo *bool, isAllDone *bool) error {
	// ตรวจสอบว่า Note เป็นของ User หรือไม่
	_, err := s.noteRepo.GetNoteByIdAndUser(noteID, userID)
	if err != nil {
		return fmt.Errorf("note not found or does not belong to the user")
	}

	// ส่งค่าที่ได้รับไปยัง Repository Layer
	return s.noteRepo.UpdateNoteStatus(noteID, userID, isTodo, isAllDone)
}


func (s *NoteService) DeleteNoteById(noteID uint, userID uint) error {
	// ตรวจสอบว่า Note เป็นของ User หรือไม่
	_, err := s.noteRepo.GetNoteByIdAndUser(noteID, userID)
	if err != nil {
		return fmt.Errorf("note not found or does not belong to the user")
	}

	// ดำเนินการลบโน้ต
	if err := s.noteRepo.DeleteNoteById(noteID); err != nil {
		return fmt.Errorf("failed to delete note: %v", err)
	}
	return nil
}

func (s *NoteService) RestoreNoteById(noteID uint, userID uint) error {
	// ตรวจสอบว่า Note เป็นของ User หรือไม่
	_, err := s.noteRepo.GetNoteByIdAndUser(noteID, userID)
	if err != nil {
		return fmt.Errorf("note not found or does not belong to the user")
	}

	// ดำเนินการกู้คืนโน้ต
	if err := s.noteRepo.RestoreNoteById(noteID); err != nil {
		return fmt.Errorf("failed to restore note: %v", err)
	}
	return nil
}

func (s *NoteService) AddTagToNote(noteID uint, tagID uint, userID uint) error {
	return s.noteRepo.AddTagToNote(noteID, tagID, userID)
}

func (s *NoteService) RemoveTagFromNote(noteID uint, tagID uint, userID uint) error {
	return s.noteRepo.RemoveTagFromNote(noteID, tagID, userID)
}
