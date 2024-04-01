package usecase

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/vishnusunil243/Job-Portal-Payment-service/entities"
	"github.com/vishnusunil243/Job-Portal-Payment-service/internal/adapters"
)

type PaymentUsecase struct {
	adapters adapters.PaymentAdapterInterface
}

func NewPaymentUsecase(adapters adapters.PaymentAdapterInterface) *PaymentUsecase {
	return &PaymentUsecase{
		adapters: adapters,
	}
}
func (p *PaymentUsecase) AddPayment(req entities.Payment) error {

	return p.adapters.AddPayment(req)
}
func (p *PaymentUsecase) AddSubscriptionPlan(req entities.Subscription) error {
	if err := p.adapters.AddSubscriptionPlan(req); err != nil {
		return err
	}
	return nil
}
func (p *PaymentUsecase) UpdateSubscriptionPlan(req entities.Subscription) error {
	return p.adapters.UpdateSubscriptionPlan(req)
}
func (p *PaymentUsecase) GetAllSubscriptionPlans() ([]entities.Subscription, error) {
	return p.adapters.GetAllSubscriptionPlans()
}
func (p *PaymentUsecase) GetSubscriptionPlanById(Id string) (entities.Subscription, error) {
	return p.adapters.GetSubscriptionPlanById(Id)
}
func (p *PaymentUsecase) AddUserSubscription(userId, subId, duration string) error {
	s := duration[0]
	userSub, err := p.adapters.GetUserSubscription(userId)
	if err != nil {
		return err
	}
	if !userSub.SubscribedTill.IsZero() && userSub.SubscribedTill.After(time.Now()) {
		num, err := strconv.Atoi(string(s))
		if err != nil {
			return err
		}
		if strings.HasSuffix(duration, "years") {
			num = 365 * num
		} else {
			num = 30 * num
		}
		duration = fmt.Sprintf("%dd", num)
	}
	return p.adapters.AddUserSubscription(userId, subId, duration)
}
func (p *PaymentUsecase) GetUserSubscription(userId string) (entities.UserSubscription, error) {
	return p.adapters.GetUserSubscription(userId)
}
func (p *PaymentUsecase) GetSubscriptionByDuration(duration string) (entities.Subscription, error) {
	return p.adapters.GetSubscriptionByDuration(duration)
}
