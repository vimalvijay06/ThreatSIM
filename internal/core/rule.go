package core

// Rule defines a detection rule that the detection engine evaluates.
//
// Example: "If 20+ login_failed events from the same IP occur within 30s,
// trigger a brute_force_attack alert with HIGH severity."
//
// Rules are loaded from YAML files, not hardcoded — so security teams
// can add/modify rules without touching code.
type Rule struct {
	// Unique name (e.g., "brute_force_attack")
	Name string `yaml:"name" json:"name"`

	// Human-readable description
	Description string `yaml:"description" json:"description"`

	// What conditions trigger this rule
	Condition RuleCondition `yaml:"condition" json:"condition"`

	// How severe is this detection
	Severity Severity `yaml:"severity" json:"severity"`

	// Base risk score to assign when this rule triggers
	RiskScore int `yaml:"risk_score" json:"risk_score"`
}

// RuleCondition specifies when a rule should fire.
type RuleCondition struct {
	// Which event type to watch (e.g., "login_failed", "port_probe")
	EventType string `yaml:"event_type" json:"event_type"`

	// Group events by this field before counting (e.g., "source_ip")
	GroupBy string `yaml:"group_by" json:"group_by"`

	// How many events needed to trigger (e.g., 20)
	Threshold int `yaml:"threshold" json:"threshold"`

	// Time window to count events in (e.g., "30s", "5m")
	Window string `yaml:"window" json:"window"`
}
