package knowledge

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRuleEngine(t *testing.T) {
	tmpDir := t.TempDir()
	km, err := NewKnowledgeManager(tmpDir)
	require.NoError(t, err)
	
	re, err := NewRuleEngine(km)
	require.NoError(t, err)
	require.NotNil(t, re)
}

func TestRuleEngine_AddRule(t *testing.T) {
	tmpDir := t.TempDir()
	km, err := NewKnowledgeManager(tmpDir)
	require.NoError(t, err)
	
	re, err := NewRuleEngine(km)
	require.NoError(t, err)
	
	rule := Rule{
		Type:    "exclude",
		Pattern: "Model View Controller|MVC",
		Reason:  "Using MV architecture instead",
	}
	
	err = re.AddRule(rule)
	require.NoError(t, err)
	
	rules := re.ListRules()
	assert.Len(t, rules, 1)
	assert.Equal(t, "exclude", rules[0].Type)
}

func TestRuleEngine_RemoveRule(t *testing.T) {
	tmpDir := t.TempDir()
	km, err := NewKnowledgeManager(tmpDir)
	require.NoError(t, err)
	
	re, err := NewRuleEngine(km)
	require.NoError(t, err)
	
	rule := Rule{
		Type:    "exclude",
		Pattern: "test pattern",
		Reason:  "test",
	}
	
	err = re.AddRule(rule)
	require.NoError(t, err)
	
	rules := re.ListRules()
	require.Len(t, rules, 1)
	ruleID := rules[0].ID
	
	err = re.RemoveRule(ruleID)
	require.NoError(t, err)
	
	rules = re.ListRules()
	assert.Len(t, rules, 0)
}

func TestRuleEngine_Validate(t *testing.T) {
	tmpDir := t.TempDir()
	km, err := NewKnowledgeManager(tmpDir)
	require.NoError(t, err)
	
	re, err := NewRuleEngine(km)
	require.NoError(t, err)
	
	tests := []struct {
		name    string
		rule    Rule
		wantErr bool
	}{
		{
			name: "valid exclude rule",
			rule: Rule{
				Type:    "exclude",
				Pattern: "test",
				Reason:  "testing",
			},
			wantErr: false,
		},
		{
			name: "invalid type",
			rule: Rule{
				Type:    "invalid_type",
				Pattern: "test",
				Reason:  "testing",
			},
			wantErr: true,
		},
		{
			name: "empty pattern",
			rule: Rule{
				Type:    "exclude",
				Pattern: "",
				Reason:  "testing",
			},
			wantErr: true,
		},
		{
			name: "invalid regex",
			rule: Rule{
				Type:    "exclude",
				Pattern: "[invalid((",
				Reason:  "testing",
			},
			wantErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := re.Validate(tt.rule)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRuleEngine_ApplyExclude(t *testing.T) {
	tmpDir := t.TempDir()
	km, err := NewKnowledgeManager(tmpDir)
	require.NoError(t, err)
	
	re, err := NewRuleEngine(km)
	require.NoError(t, err)
	
	rule := Rule{
		Type:    "exclude",
		Pattern: "MVVM|Model View ViewModel",
		Reason:  "Not using MVVM",
	}
	err = re.AddRule(rule)
	require.NoError(t, err)
	
	content := "You should use MVVM architecture for this project. Model View ViewModel is popular."
	result, err := re.Apply(content)
	require.NoError(t, err)
	
	assert.NotContains(t, result, "MVVM")
	assert.NotContains(t, result, "Model View ViewModel")
}

func TestRuleEngine_ApplyPrefer(t *testing.T) {
	tmpDir := t.TempDir()
	km, err := NewKnowledgeManager(tmpDir)
	require.NoError(t, err)
	
	re, err := NewRuleEngine(km)
	require.NoError(t, err)
	
	rule := Rule{
		Type:        "prefer",
		Pattern:     "XCTest",
		Replacement: "Swift Testing",
		Reason:      "Modern testing framework",
	}
	err = re.AddRule(rule)
	require.NoError(t, err)
	
	content := "Use XCTest for your tests. XCTest is the standard framework."
	result, err := re.Apply(content)
	require.NoError(t, err)
	
	assert.Contains(t, result, "Swift Testing")
	assert.NotContains(t, result, "XCTest")
}

func TestRuleEngine_ApplyNeverMention(t *testing.T) {
	tmpDir := t.TempDir()
	km, err := NewKnowledgeManager(tmpDir)
	require.NoError(t, err)
	
	re, err := NewRuleEngine(km)
	require.NoError(t, err)
	
	rule := Rule{
		Type:    "never_mention",
		Pattern: "Objective-C",
		Reason:  "Swift-only codebase",
	}
	err = re.AddRule(rule)
	require.NoError(t, err)
	
	content := "You can use Swift or Objective-C for iOS development. Objective-C has been around longer."
	result, err := re.Apply(content)
	require.NoError(t, err)
	
	assert.NotContains(t, result, "Objective-C")
}

func TestRuleEngine_ApplyMultipleRules(t *testing.T) {
	tmpDir := t.TempDir()
	km, err := NewKnowledgeManager(tmpDir)
	require.NoError(t, err)
	
	re, err := NewRuleEngine(km)
	require.NoError(t, err)
	
	// Add multiple rules
	rules := []Rule{
		{
			Type:    "exclude",
			Pattern: "MVVM",
			Reason:  "Not using MVVM",
		},
		{
			Type:        "prefer",
			Pattern:     "XCTest",
			Replacement: "Swift Testing",
			Reason:      "Modern framework",
		},
	}
	
	for _, rule := range rules {
		err = re.AddRule(rule)
		require.NoError(t, err)
	}
	
	content := "Use XCTest and MVVM for your iOS app."
	result, err := re.Apply(content)
	require.NoError(t, err)
	
	assert.NotContains(t, result, "MVVM")
	assert.NotContains(t, result, "XCTest")
	assert.Contains(t, result, "Swift Testing")
}

func TestRuleEngine_Persistence(t *testing.T) {
	tmpDir := t.TempDir()
	km, err := NewKnowledgeManager(tmpDir)
	require.NoError(t, err)
	
	// Create and add rules
	re1, err := NewRuleEngine(km)
	require.NoError(t, err)
	
	rule := Rule{
		Type:    "exclude",
		Pattern: "test",
		Reason:  "testing",
	}
	err = re1.AddRule(rule)
	require.NoError(t, err)
	
	// Create new engine and verify rules persisted
	re2, err := NewRuleEngine(km)
	require.NoError(t, err)
	
	rules := re2.ListRules()
	assert.Len(t, rules, 1)
	assert.Equal(t, "exclude", rules[0].Type)
}

func TestRuleEngine_CaseSensitivity(t *testing.T) {
	tmpDir := t.TempDir()
	km, err := NewKnowledgeManager(tmpDir)
	require.NoError(t, err)
	
	re, err := NewRuleEngine(km)
	require.NoError(t, err)
	
	rule := Rule{
		Type:    "exclude",
		Pattern: "(?i)mvvm", // Case insensitive
		Reason:  "Not using MVVM",
	}
	err = re.AddRule(rule)
	require.NoError(t, err)
	
	content := "Use MVVM or mvvm or Mvvm in your code."
	result, err := re.Apply(content)
	require.NoError(t, err)
	
	assert.NotContains(t, result, "MVVM")
	assert.NotContains(t, result, "mvvm")
	assert.NotContains(t, result, "Mvvm")
}
