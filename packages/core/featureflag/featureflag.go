package featureflag

import (
	"os"
	"strings"
	"sync"
)

type Flag string

const (
	AIEnabled     Flag = "AI_ENABLED"
	SecurityScan  Flag = "SECURITY_SCAN"
	Graph         Flag = "GRAPH"
	MCP           Flag = "MCP"
	Cache         Flag = "CACHE"
	Telemetry     Flag = "TELEMETRY"
	Canvas        Flag = "CANVAS"
	PluginSystem  Flag = "PLUGIN_SYSTEM"
	BackgroundWk  Flag = "BACKGROUND_WORKERS"
	Dashboard     Flag = "DASHBOARD"
)

type FeatureFlags struct {
	mu     sync.RWMutex
	flags  map[Flag]bool
	envKey string
}

var Default = New("AUTODEV_FEATURE")

func New(envKey string) *FeatureFlags {
	ff := &FeatureFlags{
		flags:  make(map[Flag]bool),
		envKey: envKey,
	}
	ff.initFromEnv()
	return ff
}

func (ff *FeatureFlags) initFromEnv() {
	val := os.Getenv(ff.envKey)
	if val == "" {
		return
	}
	// Format: AI_ENABLED=true,SECURITY_SCAN=false
	pairs := strings.Split(val, ",")
	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			continue
		}
		flag := Flag(strings.TrimSpace(parts[0]))
		enabled := strings.TrimSpace(parts[1]) == "true"
		ff.flags[flag] = enabled
	}
}

func (ff *FeatureFlags) IsEnabled(flag Flag) bool {
	// Individual env var always takes highest precedence
	envVal := os.Getenv("AUTODEV_" + string(flag))
	if envVal != "" {
		return envVal == "true" || envVal == "1"
	}

	ff.mu.RLock()
	defer ff.mu.RUnlock()
	v, ok := ff.flags[flag]
	if ok {
		return v
	}
	return true
}

func (ff *FeatureFlags) Set(flag Flag, enabled bool) {
	ff.mu.Lock()
	defer ff.mu.Unlock()
	ff.flags[flag] = enabled
}

func (ff *FeatureFlags) All() map[Flag]bool {
	ff.mu.RLock()
	defer ff.mu.RUnlock()
	result := make(map[Flag]bool, len(ff.flags))
	for k, v := range ff.flags {
		result[k] = v
	}
	// Include defaults for unset flags
	for _, f := range allFlags() {
		if _, ok := result[f]; !ok {
			result[f] = true
		}
	}
	return result
}

func allFlags() []Flag {
	return []Flag{AIEnabled, SecurityScan, Graph, MCP, Cache, Telemetry, Canvas, PluginSystem, BackgroundWk, Dashboard}
}

// Package-level convenience
func IsEnabled(flag Flag) bool { return Default.IsEnabled(flag) }
func Set(flag Flag, enabled bool) { Default.Set(flag, enabled) }
