package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormDBHandler interface {
	GetDB() (*gorm.DB, error)
}

type gormDBHandler struct {
	URL string
}

func (dh *gormDBHandler) GetDB() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dh.URL))
	if err != nil {
		return nil, err
	}
	return db, err
}

func NewDBHandler(URL string) *gormDBHandler {
	return &gormDBHandler{
		URL: URL,
	}
}
