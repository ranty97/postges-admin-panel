package domain

import "time"

type Backup struct {
	Filename  string    `json:"filename"`
	CreatedAt time.Time `json:"created_at"`
	Size      int64     `json:"size"`
}

type BackupCreated struct {
	Filename string `json:"filename"`
	Message  string `json:"message"`
	Success  bool   `json:"success"`
}

type BackupDeleted struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

type BackupRestored struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}
