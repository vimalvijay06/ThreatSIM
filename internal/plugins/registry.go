package plugins

import (
	"fmt"
	"sync"

	"github.com/Stratify-Systems/ThreatSIM/internal/core"
)

// Registry holds all available attack plugins.
// Plugins register themselves at startup, and the scenario engine
// or CLI looks them up by ID when it's time to execute.
type Registry struct {
	mu      sync.RWMutex
	plugins map[string]core.Plugin
}

// NewRegistry creates an empty plugin registry.
func NewRegistry() *Registry {
	return &Registry{
		plugins: make(map[string]core.Plugin),
	}
}

// Register adds a plugin to the registry.
// Returns an error if a plugin with the same ID is already registered.
func (r *Registry) Register(plugin core.Plugin) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := plugin.ID()
	if _, exists := r.plugins[id]; exists {
		return fmt.Errorf("plugin already registered: %s", id)
	}

	r.plugins[id] = plugin
	return nil
}

// Get returns a plugin by its ID.
// Returns an error if the plugin is not found.
func (r *Registry) Get(id string) (core.Plugin, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	plugin, exists := r.plugins[id]
	if !exists {
		return nil, fmt.Errorf("plugin not found: %s", id)
	}

	return plugin, nil
}

// List returns all registered plugins.
func (r *Registry) List() []core.Plugin {
	r.mu.RLock()
	defer r.mu.RUnlock()

	result := make([]core.Plugin, 0, len(r.plugins))
	for _, p := range r.plugins {
		result = append(result, p)
	}
	return result
}

// IDs returns a list of all registered plugin IDs.
func (r *Registry) IDs() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids := make([]string, 0, len(r.plugins))
	for id := range r.plugins {
		ids = append(ids, id)
	}
	return ids
}
