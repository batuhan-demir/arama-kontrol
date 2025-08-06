package dal

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/lib/pq" // PostgreSQL array desteği için
)

type JSONB map[string]interface{}

// GORM için gerekli metodlar
func (j *JSONB) Scan(value interface{}) error {
	if value == nil {
		*j = make(JSONB)
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, j)
	case string:
		return json.Unmarshal([]byte(v), j)
	default:
		return errors.New("cannot scan JSONB")
	}
}

func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}

// Yeni tip: JSONB array için
type JSONBArray []JSONB

func (j *JSONBArray) Scan(value interface{}) error {
	if value == nil {
		*j = make(JSONBArray, 0)
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, j)
	case string:
		return json.Unmarshal([]byte(v), j)
	default:
		return errors.New("cannot scan JSONBArray")
	}
}

func (j JSONBArray) Value() (driver.Value, error) {
	return json.Marshal(j)
}

type Call struct {
	CallId     string         `json:"call_id" gorm:"primarykey;unique;"`
	CallerNum  string         `json:"caller_num"`
	Redirects  pq.StringArray `json:"redirects" gorm:"type:text[]"`
	AddedBy    string         `gorm:"default:'system'" json:"added_by"`
	Events     JSONBArray     `json:"events" gorm:"type:jsonb"`
	CallStatus string         `json:"call_status" gorm:"default:'not_answered'"`
	AnsweredBy string         `json:"answered_by" gorm:"default:''"`
	StartedAt  string         `json:"started_at"`
	EndedAt    string         `json:"ended_at"`
	CallRecord string         `json:"call_record" gorm:"default:''"`
}

type CreateCall struct {
	CallId         string `json:"unique_id"`
	CustomerNum    string `json:"customer_num"`
	InternalNum    string `json:"internal_num"`
	IncomingNumber string `json:"incoming_number"`
	Scenario       string `json:"scenario"`
	Timestamp      string `json:"timestamp"`
	CallRecord     string `json:"call_record"`
}

type CreateCallCDR struct {
	CustomerNum    string `json:"arayan"`
	CallId         string `json:"asteriskId"`
	InternalNum    string `json:"santral"`
	IncomingNumber string `json:"trunk"`
	Scenario       string `json:"scenario"`
	Timestamp      string `json:"bas"`
	CallRecord     string `json:"seskaydi"`
}
