package entity

import (
	"encoding/json"
	"github.com/google/uuid"
	"gorm.io/datatypes"
	"time"
)

const (
	OperationTypeUnknown = iota
	OperationTypeRefill
	OperationTypeTransfer
)

var operationTypes = []string{"Refill", "Transfer"}

type Operation struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primarykey"`
	Type         int
	Details      datatypes.JSON
	DateTime     time.Time     `gorm:"default:CURRENT_TIMESTAMP"`
	Transactions []Transaction `gorm:"foreignKey:OperationId"`
}

func (o *Operation) GetOperationType() string {
	i := o.Type - 1
	if len(operationTypes) < i || i < 0 {
		return "Unknown"
	}
	return operationTypes[i]
}

type OperationRefillDetails struct {
	Source string
}

func (details OperationRefillDetails) JSON() (result datatypes.JSON) {
	data, err := json.Marshal(details)
	if err != nil {
		return
	}
	return data
}

type OperationTransferDetails struct {
	Comment     string
	AccountName string
}

func (details OperationTransferDetails) JSON() (result datatypes.JSON) {
	data, err := json.Marshal(details)
	if err != nil {
		return
	}
	return data
}
