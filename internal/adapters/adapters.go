package adapters

import (
	"time"

	"github.com/google/uuid"
	"github.com/vishnusunil243/Job-Portal-Payment-service/entities"
	"gorm.io/gorm"
)

type PaymentAdapter struct {
	DB *gorm.DB
}

func NewPaymentAdapter(db *gorm.DB) *PaymentAdapter {
	return &PaymentAdapter{
		DB: db,
	}
}
func (p *PaymentAdapter) AddPayment(req entities.Payment) error {
	id := uuid.New()
	insertQuery := `INSERT INTO payments (id,user_id,payment_ref,time) VALUES ($1,$2,$3,NOW())`
	if err := p.DB.Exec(insertQuery, id, req.UserId, req.PaymentRef).Error; err != nil {
		return err
	}
	return nil
}
func (p *PaymentAdapter) GetPayment(userId string) (entities.Payment, error) {
	return entities.Payment{}, nil
}
func (p *PaymentAdapter) AddSubscriptionPlan(req entities.Subscription) error {
	id := uuid.New()
	addSubPlan := `INSERT INTO subscriptions (id,amount,duration) VALUES ($1,$2,$3)`
	if err := p.DB.Exec(addSubPlan, id, req.Amount, req.Duration).Error; err != nil {
		return err
	}
	return nil
}
func (p *PaymentAdapter) GetAllSubscriptionPlans() ([]entities.Subscription, error) {
	selectQuery := `SELECT * FROM subscriptions`
	var res []entities.Subscription
	if err := p.DB.Raw(selectQuery).Scan(&res).Error; err != nil {
		return []entities.Subscription{}, err
	}
	return res, nil
}
func (p *PaymentAdapter) GetSubscriptionPlanById(Id string) (entities.Subscription, error) {
	selectQuery := `SELECT * FROM subscriptions WHERE id=?`
	var res entities.Subscription
	if err := p.DB.Raw(selectQuery, Id).Scan(&res).Error; err != nil {
		return entities.Subscription{}, err
	}
	return res, nil
}
func (p *PaymentAdapter) AddUserSubscription(userId, subId, duration string) error {
	var subscribedTill time.Time
	selectSubscribedTill := `SELECT subscribed_till FROM user_subscriptions WHERE user_id=?`
	if err := p.DB.Raw(selectSubscribedTill, userId).Scan(&subscribedTill).Error; err != nil {
		return err
	}
	deleteQuery := `DELETE FROM user_subscriptions WHERE user_id=?`
	if err := p.DB.Exec(deleteQuery, userId).Error; err != nil {
		return err
	}
	id := uuid.New()
	insertNewSubscription := `INSERT INTO user_subscriptions (id, user_id, sub_id, subscribed_till) VALUES (?, ?, ?, NOW()+INTERVAL '` + duration + `')`
	if !subscribedTill.IsZero() && subscribedTill.After(time.Now()) {
		newDuration, err := time.ParseDuration(duration)
		if err != nil {
			return err
		}
		insertNewSubscription = `INSERT INTO user_subscriptions (id,user_id,sub_id,subscribed_till) VALUES (?,?,?,?)`
		time := subscribedTill.Add(newDuration)
		if err := p.DB.Exec(insertNewSubscription, id, userId, subId, time).Error; err != nil {
			return err
		}
		return nil
	}
	if err := p.DB.Exec(insertNewSubscription, id, userId, subId).Error; err != nil {
		return err
	}
	return nil
}
func (p *PaymentAdapter) GetUserSubscription(userId string) (entities.UserSubscription, error) {
	var res entities.UserSubscription
	selectQuery := `SELECT * FROM user_subscriptions WHERE user_id=?`
	if err := p.DB.Raw(selectQuery, userId).Scan(&res).Error; err != nil {
		return entities.UserSubscription{}, err
	}
	return res, nil
}
