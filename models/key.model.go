package models

import (
	"ssi-gitlab.teda.th/ssi/core/utils"
	"time"
)

type Key struct {
	ID                  string     `json:"id" gorm:"id"`
	PublicKey           string     `json:"public_key" gorm:"public_key"`
	PrivateKeyEncrypted string     `json:"private_key_encrypted" gorm:"private_key_encrypted"`
	Type                string     `json:"type" gorm:"type"`
	CreatedAt           *time.Time `json:"created_at" gorm:"created_at"`
	UpdatedAt           *time.Time `json:"updated_at" gorm:"updated_at"`
	DeletedAt           *time.Time `json:"deleted_at,omitempty" gorm:"deleted_at"`
}

func (m Key) TableName() string {
	return "keys"
}

func NewKey(publicKey string, encryptedPrivateKey string, keyType string) *Key {
	return &Key{
		ID:                  utils.GetUUID(),
		PublicKey:           publicKey,
		PrivateKeyEncrypted: encryptedPrivateKey,
		Type:                keyType,
		CreatedAt:           utils.GetCurrentDateTime(),
		UpdatedAt:           utils.GetCurrentDateTime(),
	}
}
