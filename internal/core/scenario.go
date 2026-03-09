package core

// Scenario represents an attack chain — a sequence of attack steps.
//
// Real attackers don't run a single attack. They execute chains:
//
//	Step 1 → Port scan (reconnaissance)
//	Step 2 → Credential stuffing (initial access)
//	Step 3 → Privilege escalation (deeper access)
//
// Scenarios model this behavior.
type Scenario struct {
	// Unique name for this scenario
	Name string `yaml:"name" json:"name"`

	// Human-readable description
	Description string `yaml:"description" json:"description"`

	// Ordered list of attack steps
	Steps []ScenarioStep `yaml:"steps" json:"steps"`
}

// ScenarioStep is a single attack within a scenario.
type ScenarioStep struct {
	// Which plugin to run (e.g., "port_scan", "brute_force")
	PluginID string `yaml:"plugin" json:"plugin"`

	// Configuration overrides for this step
	Config PluginConfig `yaml:"config" json:"config"`

	// How long to wait after this step before the next one
	// (e.g., "5s", "1m") — simulates attacker pausing between stages
	Delay string `yaml:"delay" json:"delay"`
}

// ScenarioStatus tracks the execution state of a running scenario
type ScenarioStatus struct {
	ScenarioName string `json:"scenario_name"`
	CurrentStep  int    `json:"current_step"`
	TotalSteps   int    `json:"total_steps"`
	State        string `json:"state"` // "running", "completed", "failed", "cancelled"
	Error        string `json:"error,omitempty"`
}
