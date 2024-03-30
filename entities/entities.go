package entities

import (
	"time"

	"github.com/google/uuid"
)

type Payment struct {
	Id         uuid.UUID `gorm:"primayKey"`
	UserId     uuid.UUID
	PaymentRef string
	Time       time.Time
}
type Subscription struct {
	Id       uuid.UUID
	Duration string
	Amount   float64
}
type UserSubscription struct {
	Id             uuid.UUID `gorm:"primaryKey"`
	UserId         uuid.UUID
	SubId          uuid.UUID
	Subscription   Subscription `gorm:"foreignKey:SubId"`
	SubscribedTill time.Time
}
