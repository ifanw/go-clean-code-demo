package application

import (
	"time"
)

//goland:noinspection SpellCheckingInspection
type UploadFileCommand struct {
	FileSize        int64
	FileName        string
	Label           string
	Description     string
	TransferredTime time.Time
	UploadedTime    time.Time
}
