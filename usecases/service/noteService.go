package service

import (
	"miw/entities"
	"miw/usecases/repository"
	"time"
	"fmt"
)

type NoteUseCase interface {
	CreateNote(note *entities.Note) error
	GetAllNote(userid uint) ([]entities.Note, error)
	UpdateNote(noteID uint, note *entities.Note, userID uint) (*entities.Note, error)
	DeleteNoteById(noteID uint, userID uint) error
	RestoreNoteById(noteID uint, userID uint) error
	AddTagToNote(noteID uint, tagID uint, userID uint) error
	RemoveTagFromNote(noteID uint, tagID uint, userID uint) error 
}

type NoteService struct {
	repo repository.NoteRepository
}

func NewNoteService(repo repository.NoteRepository) *NoteService {
	return &NoteService{repo: repo}
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

	return s.repo.CreateNote(note)
}


func (s *NoteService) GetAllNote(userid uint) ([]entities.Note, error) {
	return s.repo.GetAllNoteByUserId(userid)
}


func (s *NoteService) UpdateNote(noteID uint, note *entities.Note, userID uint) (*entities.Note, error) {
	existingNote, err := s.repo.GetNoteByIdAndUser(noteID, userID)
	if err != nil {
		return nil, fmt.Errorf("note not found or does not belong to the user")
	}

	// อัปเดตฟิลด์ที่ไม่ใช่ค่าว่าง
	if note.Title != "" {
		existingNote.Title = note.Title
	}
	if note.Content != "" {
		existingNote.Content = note.Content
	}
	if note.Color != "" {
		existingNote.Color = note.Color
	}
	if note.Priority != 0 {
		existingNote.Priority = note.Priority
	}
	existingNote.TodoItems = note.TodoItems

	// คำนวณ IsAllDone
	existingNote.IsAllDone = true
	for _, todo := range note.TodoItems {
		if !todo.IsDone {
			existingNote.IsAllDone = false
			break
		}
	}

	existingNote.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")

	if err := s.repo.UpdateNote(existingNote); err != nil {
		return nil, err
	}

	return existingNote, nil
}


func (s *NoteService) DeleteNoteById(noteID uint, userID uint) error {
    // ตรวจสอบว่า Note เป็นของ User หรือไม่
    _, err := s.repo.GetNoteByIdAndUser(noteID, userID)
    if err != nil {
        return fmt.Errorf("note not found or does not belong to the user")
    }

    // ดำเนินการลบโน้ต
    if err := s.repo.DeleteNoteById(noteID); err != nil {
        return fmt.Errorf("failed to delete note: %v", err)
    }
    return nil
}

func (s *NoteService) RestoreNoteById(noteID uint, userID uint) error {
    // ตรวจสอบว่า Note เป็นของ User หรือไม่
    _, err := s.repo.GetNoteByIdAndUser(noteID, userID)
    if err != nil {
        return fmt.Errorf("note not found or does not belong to the user")
    }

    // ดำเนินการกู้คืนโน้ต
    if err := s.repo.RestoreNoteById(noteID); err != nil {
        return fmt.Errorf("failed to restore note: %v", err)
    }
    return nil
}


func (s *NoteService) AddTagToNote(noteID uint, tagID uint, userID uint) error {
    return s.repo.AddTagToNote(noteID, tagID, userID)
}


func (s *NoteService) RemoveTagFromNote(noteID uint, tagID uint, userID uint) error {
    return s.repo.RemoveTagFromNote(noteID, tagID, userID)
}


