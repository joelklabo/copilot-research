package knowledge

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"time"

	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

// RuleEngine manages and applies user-defined rules
type RuleEngine struct {
	rules      []Rule
	km         *KnowledgeManager
	rulesFile  string
	mu         sync.RWMutex
}

// RulesConfig represents the YAML structure for rules
type RulesConfig struct {
	Rules []Rule `yaml:"rules"`
}

// NewRuleEngine creates a new rule engine
func NewRuleEngine(km *KnowledgeManager) (*RuleEngine, error) {
	rulesFile := filepath.Join(km.baseDir, "rules.yaml")
	
	re := &RuleEngine{
		rules:     make([]Rule, 0),
		km:        km,
		rulesFile: rulesFile,
	}
	
	// Load existing rules if file exists
	if err := re.load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to load rules: %w", err)
	}
	
	return re, nil
}

// load reads rules from YAML file
func (re *RuleEngine) load() error {
	data, err := os.ReadFile(re.rulesFile)
	if err != nil {
		return err
	}
	
	var config RulesConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return fmt.Errorf("failed to parse rules YAML: %w", err)
	}
	
	re.mu.Lock()
	re.rules = config.Rules
	re.mu.Unlock()
	
	return nil
}

// save writes rules to YAML file
func (re *RuleEngine) save() error {
	re.mu.RLock()
	config := RulesConfig{Rules: re.rules}
	re.mu.RUnlock()
	
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal rules: %w", err)
	}
	
	if err := os.WriteFile(re.rulesFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write rules file: %w", err)
	}
	
	// Note: Git commit should be done separately by the caller if needed
	// We don't auto-commit here to avoid test hangs
	
	return nil
}

// AddRule adds a new rule
func (re *RuleEngine) AddRule(rule Rule) error {
	// Validate rule
	if err := re.Validate(rule); err != nil {
		return err
	}
	
	// Generate ID if not set
	if rule.ID == "" {
		rule.ID = uuid.New().String()
	}
	
	// Set created time
	if rule.CreatedAt.IsZero() {
		rule.CreatedAt = time.Now()
	}
	
	re.mu.Lock()
	re.rules = append(re.rules, rule)
	re.mu.Unlock()
	
	return re.save()
}

// RemoveRule removes a rule by ID
func (re *RuleEngine) RemoveRule(id string) error {
	re.mu.Lock()
	
	found := false
	newRules := make([]Rule, 0, len(re.rules))
	for _, rule := range re.rules {
		if rule.ID != id {
			newRules = append(newRules, rule)
		} else {
			found = true
		}
	}
	
	if !found {
		re.mu.Unlock()
		return fmt.Errorf("rule not found: %s", id)
	}
	
	re.rules = newRules
	re.mu.Unlock()
	
	return re.save()
}

// ListRules returns all rules
func (re *RuleEngine) ListRules() []Rule {
	re.mu.RLock()
	defer re.mu.RUnlock()
	
	// Return a copy
	rules := make([]Rule, len(re.rules))
	copy(rules, re.rules)
	return rules
}

// Validate validates a rule
func (re *RuleEngine) Validate(rule Rule) error {
	// Check type
	validTypes := map[string]bool{
		"exclude":        true,
		"prefer":         true,
		"always_mention": true,
		"never_mention":  true,
	}
	
	if !validTypes[rule.Type] {
		return fmt.Errorf("invalid rule type: %s", rule.Type)
	}
	
	// Check pattern
	if rule.Pattern == "" {
		return fmt.Errorf("pattern cannot be empty")
	}
	
	// Validate regex
	if _, err := regexp.Compile(rule.Pattern); err != nil {
		return fmt.Errorf("invalid regex pattern: %w", err)
	}
	
	// Type-specific validation
	if rule.Type == "prefer" && rule.Replacement == "" {
		return fmt.Errorf("prefer rule requires replacement")
	}
	
	return nil
}

// Apply applies all rules to content
func (re *RuleEngine) Apply(content string) (string, error) {
	re.mu.RLock()
	rules := make([]Rule, len(re.rules))
	copy(rules, re.rules)
	re.mu.RUnlock()
	
	result := content
	
	for _, rule := range rules {
		var err error
		switch rule.Type {
		case "exclude":
			result, err = re.applyExclude(result, rule)
		case "prefer":
			result, err = re.applyPrefer(result, rule)
		case "never_mention":
			result, err = re.applyNeverMention(result, rule)
		case "always_mention":
			result, err = re.applyAlwaysMention(result, rule)
		}
		
		if err != nil {
			return result, fmt.Errorf("failed to apply rule %s: %w", rule.ID, err)
		}
	}
	
	return result, nil
}

// applyExclude removes matching content
func (re *RuleEngine) applyExclude(content string, rule Rule) (string, error) {
	regex, err := regexp.Compile(rule.Pattern)
	if err != nil {
		return content, err
	}
	
	// Replace matching patterns with empty string
	// This removes just the matched text, not entire sentences
	return regex.ReplaceAllString(content, ""), nil
}

// applyPrefer replaces pattern with replacement
func (re *RuleEngine) applyPrefer(content string, rule Rule) (string, error) {
	regex, err := regexp.Compile(rule.Pattern)
	if err != nil {
		return content, err
	}
	
	return regex.ReplaceAllString(content, rule.Replacement), nil
}

// applyNeverMention removes any mention of pattern
func (re *RuleEngine) applyNeverMention(content string, rule Rule) (string, error) {
	return re.applyExclude(content, rule)
}

// applyAlwaysMention ensures pattern is mentioned (placeholder)
func (re *RuleEngine) applyAlwaysMention(content string, rule Rule) (string, error) {
	// This is complex - would need to understand context
	// For now, just check if it's present
	regex, err := regexp.Compile(rule.Pattern)
	if err != nil {
		return content, err
	}
	
	if !regex.MatchString(content) {
		// Add a note about it
		content += fmt.Sprintf("\n\nNote: Consider %s.", rule.Pattern)
	}
	
	return content, nil
}
