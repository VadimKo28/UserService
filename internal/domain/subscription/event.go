package subscription

import "time"

const (
	EventTypeCreated = "subscription.created"
	EventTypeUpdated = "subscription.updated"
)

type Event struct {
	Type       string       `json:"type"`
	OccurredAt time.Time    `json:"occurred_at"`
	Payload    EventPayload `json:"payload"`
}

type EventPayload struct {
	ID          int        `json:"id"`
	UserID      int        `json:"user_id"`
	ServiceName string     `json:"service_name"`
	Price       int        `json:"price"`
	StartDate   time.Time  `json:"start_date"`
	EndDate     *time.Time `json:"end_date,omitempty"`
}

func NewCreatedEvent(id int, dto CreateSubscriptionDTO) Event {
	return Event{
		Type:       EventTypeCreated,
		OccurredAt: time.Now().UTC(),
		Payload: EventPayload{
			ID:          id,
			UserID:      dto.UserID,
			ServiceName: dto.ServiceName,
			Price:       dto.Price,
			StartDate:   dto.StartDate,
			EndDate:     toNullableTime(dto.EndDate),
		},
	}
}

func NewUpdatedEvent(dto UpdateSubscriptionDTO) Event {
	return Event{
		Type:       EventTypeUpdated,
		OccurredAt: time.Now().UTC(),
		Payload: EventPayload{
			ID:          dto.ID,
			UserID:      dto.UserID,
			ServiceName: dto.ServiceName,
			Price:       dto.Price,
			StartDate:   dto.StartDate,
			EndDate:     toNullableTime(dto.EndDate),
		},
	}
}

func toNullableTime(value time.Time) *time.Time {
	if value.IsZero() {
		return nil
	}
	normalized := value.UTC()
	return &normalized
}
