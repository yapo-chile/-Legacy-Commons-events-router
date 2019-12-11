package repository

import (
	"encoding/json"
	"fmt"

	"github.mpi-internal.com/Yapo/events-router/pkg/domain"
	"github.mpi-internal.com/Yapo/events-router/pkg/usecases"
)

// router connects with remote configuration to obtain destination topic and
// make routing for each event
type router struct {
	remoteConfig RConfig
}

// RConfig allows operations for remote configuration
type RConfig interface {
	Get(key string) string
}

// MakeRouter creates a new instance of Router
func MakeRouter(rConfig RConfig) usecases.Router {
	return &router{
		remoteConfig: rConfig,
	}
}

// GetTopics gets destination topic for incoming event
func (r *router) GetTopics(event domain.Event) (topics []string, err error) {
	config := r.remoteConfig.Get(fmt.Sprintf("event.%s.topics", event.Type))
	if config == "" {
		return []string{},
			fmt.Errorf("topic configuration not set for event of type: %s", event.Type)
	}
	if err = json.Unmarshal([]byte(config), &topics); err != nil {
		return []string{},
			fmt.Errorf("unable to unmarshal topics using remote config: %+v", err)
	}
	return
}
