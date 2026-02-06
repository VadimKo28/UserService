package subscription

import "time"

type Subscription struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	ServiceName string    `json:"service_name"`
	Price       int       `json:"price"`
	StartDate   time.Time `json:"start_date"`
	EndDate     time.Time `json:"end_date"`
}

type CreateSubscriptionDTO struct {
	UserID      int
	ServiceName string
	Price       int
	StartDate   time.Time
	EndDate     time.Time
}

type UpdateSubscriptionDTO struct {
	ID          int
	UserID      int
	ServiceName string
	Price       int
	StartDate   time.Time
	EndDate     time.Time
}
