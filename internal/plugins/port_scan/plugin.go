package portscan

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/Stratify-Systems/ThreatSIM/internal/core"
)

// Plugin simulates a port scanning attack.
//
// What it does:
//   - Generates port_probe events across a range of ports
//   - Simulates an attacker performing network reconnaissance
//   - Marks common service ports (22, 80, 443, etc.) as "open"
//
// What this tests in your security pipeline:
//   - Network scan detection
//   - Unusual traffic pattern alerting
//   - IDS/IPS effectiveness
type Plugin struct{}

var _ core.Plugin = (*Plugin)(nil)

// Common ports that are typically "open" on servers
var commonOpenPorts = map[int]string{
	22:   "ssh",
	80:   "http",
	443:  "https",
	3306: "mysql",
	5432: "postgresql",
	6379: "redis",
	8080: "http-alt",
	8443: "https-alt",
	9090: "prometheus",
}

func (p *Plugin) ID() string   { return "port_scan" }
func (p *Plugin) Name() string { return "Port Scanning Attack" }
func (p *Plugin) Description() string {
	return "Simulates network port scanning to discover open services"
}

func (p *Plugin) DefaultConfig() core.PluginConfig {
	return core.PluginConfig{
		Target:   "10.0.0.1",
		Service:  "network",
		SourceIP: "10.1.2.3",
		Duration: "20s",
		Rate:     20, // 20 port probes per second
		Params: map[string]any{
			"port_start": 1,
			"port_end":   1024,
		},
	}
}

func (p *Plugin) Execute(ctx context.Context, config core.PluginConfig, sink core.EventSink) error {
	duration, err := time.ParseDuration(config.Duration)
	if err != nil {
		duration = 20 * time.Second
	}

	rate := config.Rate
	if rate <= 0 {
		rate = 20
	}
	interval := time.Second / time.Duration(rate)

	// Parse port range
	portStart := 1
	portEnd := 1024
	if params := config.Params; params != nil {
		if v, ok := params["port_start"].(int); ok {
			portStart = v
		}
		if v, ok := params["port_end"].(int); ok {
			portEnd = v
		}
	}

	sourceIP := config.SourceIP
	if sourceIP == "" {
		sourceIP = fmt.Sprintf("10.%d.%d.%d", rand.Intn(256), rand.Intn(256), rand.Intn(256))
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	deadline := time.After(duration)
	currentPort := portStart
	eventCount := 0

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-deadline:
			return nil

		case <-ticker.C:
			if currentPort > portEnd {
				return nil // Scanned all ports
			}

			// Determine if this port is "open"
			portStatus := "closed"
			serviceName := "unknown"
			if svc, isOpen := commonOpenPorts[currentPort]; isOpen {
				portStatus = "open"
				serviceName = svc
			}

			event := core.Event{
				ID:        fmt.Sprintf("ps-%d-%d", time.Now().UnixNano(), eventCount),
				Type:      "port_probe",
				SourceIP:  sourceIP,
				Target:    config.Target,
				Service:   config.Service,
				Timestamp: time.Now(),
				PluginID:  p.ID(),
				Metadata: map[string]any{
					"port":         currentPort,
					"port_status":  portStatus,
					"service_name": serviceName,
					"protocol":     "tcp",
				},
			}

			if err := sink(event); err != nil {
				return fmt.Errorf("failed to send event: %w", err)
			}

			currentPort++
			eventCount++
		}
	}
}
