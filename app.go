package main

import (
	"log"

	"gorm.io/gorm"
)

type app struct {
	config config
	logger *log.Logger
	db     *gorm.DB
}
