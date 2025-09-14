package domain

import (
	"encoding/json"
	"time"

	"git.ice.global/packages/beeorm/v4"
)

type OutboxEvent struct {
	beeorm.ORM `orm:"table=outbox_events;redisCache"`
	ID         uint64    `orm:"required"`
	EventType  string    `orm:"required"`
	EntityID   uint64    `orm:"required"`
	EntityType string    `orm:"required"`
	Payload    string    `orm:"required"`
	Published  bool      `orm:"required"`
	CreatedAt  time.Time `orm:"time=true"`
	UpdatedAt  time.Time `orm:"time=true"`
}

func (o *OutboxEvent) GetID() uint64 {
	return o.ID
}

func (o *OutboxEvent) SetID(id uint64) {
	o.ID = id
}

func NewOutboxEvent(eventType, entityType string, entityID uint64, payload interface{}) (*OutboxEvent, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &OutboxEvent{
		EventType:  eventType,
		EntityID:   entityID,
		EntityType: entityType,
		Payload:    string(payloadBytes),
		Published:  false,
		CreatedAt:  now,
		UpdatedAt:  now,
	}, nil
}

func (OutboxEvent) TableName() string {
	return "outbox_events"
}
