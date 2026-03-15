package domain

// Kafka event type constants for agent-core-service domain events
const (
	EventBotSeen           = "agent.bot.seen"
	EventBotOnline         = "agent.bot.online"
	EventBotOffline        = "agent.bot.offline"
	EventBotDegraded       = "agent.bot.degraded"
	EventMemoryCreated     = "agent.memory.created"
	EventMemoryDeleted     = "agent.memory.deleted"
	EventTodoUpdated       = "agent.todo.updated"
	EventTodoCompleted     = "agent.todo.completed"
	EventHeartbeatReceived = "agent.heartbeat.received"
)
