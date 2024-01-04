package routes

import "log"

type Handlers struct {
	Logger *log.Logger
}

func NewHandlers( logger *log.Logger) *Handlers {
	return &Handlers{
		Logger:      logger,
	}
}