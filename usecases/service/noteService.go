package service

import (
	"miw/entities"
	"miw/usecases/repository"
	"time"
)

type NoteUseCase interface {
	CreateNote(note *entities.Note) error
	GetAllNote(userid uint) ([]entities.Note, error)
	UpdateNote(noteid uint, note *entities.Note) (*entities.Note, error)
	AddTagToNote(noteID uint, tagIDs uint) error // รองรับหลายแท็ก
}

type NoteService struct {
	repo repository.NoteRepository
}

func NewNoteService(repo repository.NoteRepository) *NoteService {
	return &NoteService{repo: repo}
}

func (s *NoteService) CreateNote(note *entities.Note) error {
	// ตั้งค่าเวลา CreatedAt ให้เป็นเวลาประเทศไทย
	timeCreate, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return err
	}
	note.CreatedAt = time.Now().In(timeCreate).Format("2006-01-02 15:04:05")
	return s.repo.CreateNote(note)
}

func (s *NoteService) GetAllNote(userid uint) ([]entities.Note, error) {
	return s.repo.GetAllNoteByUserId(userid)
}

func (s *NoteService) AddTagToNote(noteID uint, tagID uint) error {
	return s.repo.AddTagToNote(noteID, tagID)
}

func (s *NoteService) UpdateNote(noteID uint, note *entities.Note) (*entities.Note, error) {
	// ดึงโน้ตที่ต้องการแก้ไขจากฐานข้อมูล
	existingNote, err := s.repo.GetNoteById(noteID)
	if err != nil {
		return nil, err
	}

	// ตรวจสอบและอัปเดตฟิลด์ที่ไม่ใช่ค่าว่างหรือค่าเริ่มต้น
	if note.Title != "" {
		existingNote.Title = note.Title
	}
	if note.Content != "" {
		existingNote.Content = note.Content
	}
	if note.Color != "" {
		existingNote.Color = note.Color
	}
	if note.Priority != 0 { // สมมติว่าค่า Priority = 0 หมายถึงไม่ได้เปลี่ยนแปลง
		existingNote.Priority = note.Priority
	}
	if note.IsTodo {
		existingNote.IsTodo = note.IsTodo
	}
	if note.TodoStatus {
		existingNote.TodoStatus = note.TodoStatus
	}

	// ตั้งเวลา UpdatedAt ใหม่ให้เป็นเวลาประเทศไทย
	timeUpdate, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return nil, err
	}
	existingNote.UpdatedAt = time.Now().In(timeUpdate).Format("2006-01-02 15:04:05")

	// บันทึกโน้ตที่แก้ไขแล้ว
	if err := s.repo.UpdateNote(existingNote); err != nil {
		return nil, err
	}

	// ดึงข้อมูลโน้ตที่อัปเดตกลับมาเพื่อส่งต่อ
	return s.repo.GetNoteById(noteID)
}
