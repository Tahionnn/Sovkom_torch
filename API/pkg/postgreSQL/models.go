package postgreSQL

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"time"
)

// User представляет пользователя системы.
// @Description  Структура пользователя
type User struct {
	gorm.Model
	// Имя пользователя (уникальное)
	Username string `gorm:"unique;not null" json:"username"`
	// Хэш пароля (скрыт в JSON)
	PasswordHash string `gorm:"not null" json:"password"`
}

// Receipt представляет чек пользователя
type Receipt struct {
	gorm.Model
	Username    string         `gorm:"not null" json:"username"`              // Владелец чека
	UUID        string         `gorm:"type:uuid;not null;unique" json:"uuid"` // Уникальный идентификатор чека
	Name        string         `gorm:"not null" json:"name"`                  // Название чека
	Date        time.Time      `gorm:"not null" json:"date"`                  // Дата покупки
	JSONCheck   datatypes.JSON `gorm:"type:jsonb;not null" json:"json_check"` // Сам чек в формате JSON
	IsJSONReady bool           `gorm:"default:false" json:"is_json_ready"`    // Готовность JSON-данных
}

// ReceiptCategories — связь UUID чека и его категорий в JSON
type ReceiptCategories struct {
	UUID    string         `gorm:"primaryKey" json:"uuid"`
	JSONCat datatypes.JSON `gorm:"type:jsonb;not null" json:"categories"`
}

// JSONCat — тип для хранения JSON с категориями
type JSONCat map[string]int
