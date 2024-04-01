package usecase

import "github.com/vishnusunil243/Job-Portal-Payment-service/entities"

type PaymentUsecaseInterface interface {
	AddPayment(entities.Payment) error
	AddSubscriptionPlan(entities.Subscription) error
	UpdateSubscriptionPlan(req entities.Subscription) error
	GetAllSubscriptionPlans() ([]entities.Subscription, error)
	GetSubscriptionPlanById(Id string) (entities.Subscription, error)
	AddUserSubscription(userId, subId, duration string) error
	GetUserSubscription(userId string) (entities.UserSubscription, error)
	GetSubscriptionByDuration(duration string) (entities.Subscription, error)
}
