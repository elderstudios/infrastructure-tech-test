package main

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type MemberDto struct {
	Name string
	Age  int
}

type Member struct {
	gorm.Model
	ID uuid.UUID `gorm:"primaryKey"`
	Name string
	Age int
	CreatedAt time.Time
}