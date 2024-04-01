package entities

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	Id         uuid.UUID `gorm:"primayKey" json:"id,omitempty"`
	UserId     uuid.UUID `json:"user_id,omitempty"`
	PaymentRef string    `json:"payment_ref,omitempty"`
	Time       time.Time `json:"time,omitempty"`
}

type Subscription struct {
	Id       uuid.UUID `json:"id,omitempty"`
	Duration string    `json:"duration,omitempty"`
	Amount   float64   `json:"amount,omitempty"`
}

type UserSubscription struct {
	Id             uuid.UUID    `gorm:"primaryKey" json:"Id,omitempty"`
	UserId         uuid.UUID    `json:"user_id,omitempty"`
	SubId          uuid.UUID    `json:"sub_id,omitempty"`
	Subscription   Subscription `gorm:"foreignKey:SubId"`
	SubscribedTill time.Time    `json:"subscribed_till,omitempty"`
}
