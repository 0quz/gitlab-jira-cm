package model

import "gorm.io/gorm"

// dependency part
type MarcoModel struct {
	db *gorm.DB
}

func NewMarcoModel(db *gorm.DB) *MarcoModel {
	return &MarcoModel{
		db: db,
	}
}
