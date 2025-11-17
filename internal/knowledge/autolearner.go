package knowledge

import (
	"fmt"

	"github.com/joelklabo/copilot-research/internal/research" // Assuming ResearchResult is here
)

// AutoLearner extracts and stores knowledge from research results.
type AutoLearner struct {
	km KnowledgeManagerInterface // Changed to interface
	// Potentially other dependencies like a research engine for context
}

// NewAutoLearner creates a new AutoLearner instance.
func NewAutoLearner(km KnowledgeManagerInterface) *AutoLearner { // Changed to interface
	return &AutoLearner{
		km: km,
	}
}

// AnalyzeResult analyzes a research result and suggests knowledge entries.
func (al *AutoLearner) AnalyzeResult(result *research.ResearchResult) ([]*Knowledge, error) {
	// Placeholder for actual analysis logic
	// For now, just create a dummy knowledge entry
	
	if result == nil || result.Content == "" {
		return nil, fmt.Errorf("research result is empty or nil")
	}

	// Simple extraction: create one knowledge entry from the result content
	k := &Knowledge{
		Topic:      result.Query, // Use query as topic for simplicity
		Content:    result.Content,
		Source:     "auto-learned",
		Confidence: 0.7, // Default confidence for auto-learned
		Tags:       []string{"auto-learned", result.Mode},
	}

	return []*Knowledge{k}, nil
}