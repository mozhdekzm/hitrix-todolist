package service

//func ServiceProviderTodoService() *service.DefinitionGlobal {
//	return &service.DefinitionGlobal{
//		Name: "todo_service",
//		Build: func(ctn di.Container) (interface{}, error) {
//			db := ctn.Get("orm_engine").(*mysql.BeeORMEngine)
//			redisClient := ctn.Get("redis_pool_default").(*redis.Client)
//
//			todoRepo := mysql.NewTodoRepository(*db)
//			outboxRepo := mysql.NewOutboxRepository(*db)
//			streamPublisher := redis.NewStreamPublisher(*redisClient, ctn.(hitrix.DIContainer).App().Config())
//
//			todoService := usecase.NewTodoService(todoRepo, outboxRepo, streamPublisher, *db)
//			return todoService, nil
//		},
//	}
//}
//
//// OutboxWorker - Global
//func ServiceProviderOutboxWorker() *hitrix.DefinitionGlobal {
//	return &hitrix.DefinitionGlobal{
//		Name: "outbox_worker",
//		Build: func(ctn di.Container) (interface{}, error) {
//			todoService := ctn.Get("todo_service").(*usecase.TodoService)
//			outboxWorker := worker.NewOutboxWorker(todoService.OutboxRepo, todoService.StreamPublisher)
//			return outboxWorker, nil
//		},
//	}
//}
