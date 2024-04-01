package cronjobs

import (
	"context"
	"log"
	"time"

	"github.com/robfig/cron"
	"github.com/vishnusunil243/Job-Portal-Payment-service/entities"
	"github.com/vishnusunil243/Job-Portal-Payment-service/internal/service"
	"github.com/vishnusunil243/Job-Portal-Payment-service/kafka"
	"github.com/vishnusunil243/Job-Portal-proto-files/pb"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

type CronJob struct {
	service         *service.PaymentService
	notificatonConn pb.EmailServiceClient
	DB              *gorm.DB
}

func NewCronJob(service *service.PaymentService, db *gorm.DB) *CronJob {
	notificationConn, _ := grpc.Dial(":8087", grpc.WithInsecure())
	return &CronJob{
		service:         service,
		notificatonConn: pb.NewEmailServiceClient(notificationConn),
		DB:              db,
	}
}
func (c *CronJob) Start() {
	cron := cron.New()
	err := cron.AddFunc("12 00 * * *", func() {
		c.CheckSubscriptions()
	})
	if err != nil {
		log.Print("error scheduling cron job ", err)
	}
	cron.Start()
}
func (c *CronJob) CheckSubscriptions() {
	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02")
	var usersubscriptions []entities.UserSubscription
	selectQuery := `SELECT * FROM user_subscriptions WHERE DATE(subscribed_till)=$1`
	if err := c.DB.Raw(selectQuery, tomorrow).Scan(&usersubscriptions).Error; err != nil {
		log.Print("error while performing cron jobs ", err)
	}
	for _, usrSub := range usersubscriptions {
		userData, err := c.service.UserConn.GetUser(context.Background(), &pb.GetUserById{
			Id: usrSub.UserId.String(),
		})
		if err != nil {
			log.Print("error obtaining user info")
		}
		if err := kafka.ProduceWarningSubscriptionEndingMessage(userData.Email); err != nil {
			log.Print("error sending email ", err)
		}
		if _, err := c.notificatonConn.AddNotification(context.Background(), &pb.AddNotificationRequest{
			UserId:  usrSub.UserId.String(),
			Message: `{"message":"your subscription is ending tomorrow please subscribe to continue using our services"}`,
		}); err != nil {
			log.Print("error while sending notifications ", err)
		}
	}
}
