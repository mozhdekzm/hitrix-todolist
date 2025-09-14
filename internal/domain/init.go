package domain

import (
	"git.ice.global/packages/beeorm/v4"
	"git.ice.global/packages/hitrix/pkg/entity"
)

func Init(registry *beeorm.Registry) {
	registry.RegisterEntity(
		&TodoItem{},
		&OutboxEvent{},
	)
	registry.RegisterEntity(&entity.RequestLoggerEntity{})

}
