package models

import (
  "time"
  "gorm.io/gorm"
)

type Gallery struct {
  ID uint `gorm:"primaryKey" json:"id"`
  FileName string `json:"file_name"`
  Url string `json:"url"`
  FileSize uint32 `json:"file_size"`
  IsPrivate bool `json:"is_private"`
  Description string `json:"description"`
  CreatedAt time.Time `json:"created_at"`
  UpdatedAt time.Time `json:"updated_at"`
  DeletedAt gorm.DeletedAt `json:"deleted_at"`
}

func (Gallery) TableName() string {
  return "galleries"
}
