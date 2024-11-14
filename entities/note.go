package entities

type Note struct {
	NoteID      uint      `json:"note_id" gorm:"primaryKey"`
	UserID      uint      `json:"user_id"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Color       string    `json:"color"`
	Priority    int       `json:"priority"`
	IsTodo      bool      `json:"is_todo"`
	TodoStatus  bool      `json:"todo_status"`
	CreatedAt   string    `json:"created_at"`
	UpdatedAt   string    `json:"updated_at"`
	DeletedAt   string    `json:"deleted_at"`
	Tags        []Tag     `gorm:"many2many:note_tags"`
	Reminders   []Reminder `gorm:"foreignKey:NoteID"`
	Event       Event     `gorm:"foreignKey:NoteID;constraint:OnDelete:CASCADE;"`
}

type Reminder struct {
	ReminderID   uint   `json:"reminder_id" gorm:"primaryKey"`
	NoteID       uint   `json:"note_id"`
	ReminderTime string `json:"reminder_time"`
	Recurring    bool   `json:"recurring"`
	Frequency    string `json:"frequency"`
}

type NoteTag struct {
	NoteID uint `json:"note_id"`
	TagID  uint `json:"tag_id"`
}

type Tag struct {
	TagID   uint    `json:"tag_id" gorm:"primaryKey"`
	TagName string  `json:"tag_name"`
	Notes   []Note  `gorm:"many2many:note_tags"`
}

type ShareNote struct {
	ShareNoteID uint   `json:"share_note_id" gorm:"primaryKey"`
	NoteID      uint   `json:"note_id"`
	SharedWith  uint   `json:"shared_with"`
}

type Event struct {
	EventID   uint   `json:"event_id" gorm:"primaryKey"`
	NoteID    uint   `json:"note_id" gorm:"unique"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}
