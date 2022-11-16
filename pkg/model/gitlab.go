package model

import (
	"time"
)

// models for event
type Event struct {
	MrUrl string    `gorm:"primary_key" json:"mr_url"`
	Date  time.Time `gorm:"autoCreateTime" json:"date"`
	Topic string    `json:"topic"`
	//Detail string     `json:"detail"`
}

// models for hot_fix_error
type HotFixError struct {
	MrUrl string    `gorm:"primary_key" json:"mr_url"`
	Date  time.Time `gorm:"autoCreateTime" json:"date"`
	//Detail string     `json:"detail"`
}

type Tabler interface {
	TableName() string
}

// TableName overrides the table name used by Events to `event`
func (Event) TableName() string {
	return "event"
}

// TableName overrides the table name used by HotFixErrors to `hot_fix_error`
func (HotFixError) TableName() string {
	return "hot_fix_error"
}

// Add an event to DB
func (m *MarcoModel) CreateEvent(e *Event) error {
	return m.db.Create(e).Error
}

// Find an event if it exists in DB
func (m *MarcoModel) FindEventByMrUrl(mrUrl string) error {
	return m.db.First(&Event{}, "mr_url = ?", mrUrl).Error
}

// Add an hotfix error to DB
func (m *MarcoModel) CreateHotFixError(h *HotFixError) error {
	return m.db.Create(h).Error
}
