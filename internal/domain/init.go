package domain

import (
	"git.ice.global/packages/beeorm/v4"
)

func Init(registry *beeorm.Registry) {
	registry.RegisterEntity(
		&TodoItem{},
		&OutboxEvent{},
	)
}
