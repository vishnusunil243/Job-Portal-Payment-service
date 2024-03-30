package concurrency

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/vishnusunil243/Job-Portal-Payment-service/entities"
	"github.com/vishnusunil243/Job-Portal-Payment-service/internal/service"
	"github.com/vishnusunil243/Job-Portal-proto-files/pb"
	"gorm.io/gorm"
)

type Concurrency struct {
	DB      *gorm.DB
	service *service.PaymentService
	mu      sync.Mutex
}

func NewConcurrency(DB *gorm.DB, service *service.PaymentService) *Concurrency {
	return &Concurrency{
		DB:      DB,
		service: service,
	}
}
func (c *Concurrency) Concurrency() {
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for range ticker.C {
			c.mu.Lock()
			c.UpdateSubscribed()
		}
		c.mu.Unlock()
	}()
}

func (c *Concurrency) UpdateSubscribed() {
	var res []entities.UserSubscription
	selectUsersQuery := `SELECT * FROM user_subscriptions WHERE 
	NOW() > subscribed_till`
	if err := c.DB.Raw(selectUsersQuery).Scan(&res).Error; err != nil {
		log.Print("error performing concurrency ", err)
	}
	for _, user := range res {
		if _, err := c.service.UserConn.UpdateSubscription(context.Background(), &pb.UpdateSubscriptionRequest{
			UserId:       user.UserId.String(),
			Subscription: false,
		}); err != nil {
			log.Print("error while performing concurrecny ", err)
		}
	}
}
