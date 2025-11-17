package research

import (
	"context"
	"fmt"
	"time"

	"github.com/joelklabo/copilot-research/internal/db"
	"github.com/joelklabo/copilot-research/internal/prompts"
	"github.com/joelklabo/copilot-research/internal/provider"
)

// Engine coordinates the research process
type Engine struct {
	db              *db.SQLiteDB
	promptLoader    *prompts.PromptLoader
	providerManager *provider.ProviderManager
}

// ResearchOptions contains options for a research query
type ResearchOptions struct {
	Query      string
	Mode       string
	PromptName string
	NoStore    bool
}

// ResearchResult contains the result of a research query
type ResearchResult struct {
	Query     string
	Mode      string
	Content   string
	Duration  time.Duration
	SessionID int64
}

// NewEngine creates a new research engine
func NewEngine(database *db.SQLiteDB, loader *prompts.PromptLoader, providerMgr *provider.ProviderManager) *Engine {
	return &Engine{
		db:              database,
		promptLoader:    loader,
		providerManager: providerMgr,
	}
}

// Research executes a research query
func (e *Engine) Research(ctx context.Context, opts ResearchOptions, progress chan<- string) (*ResearchResult, error) {
	start := time.Now()

	// Check context first
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// Send progress: Loading prompt
	if progress != nil {
		progress <- "Loading prompt..."
	}

	// Load the prompt template
	promptName := opts.PromptName
	if promptName == "" {
		promptName = "default"
	}

	prompt, err := e.promptLoader.Load(promptName)
	if err != nil {
		return nil, fmt.Errorf("failed to load prompt: %w", err)
	}

	// Render the prompt with variables
	mode := opts.Mode
	if mode == "" {
		mode = "quick"
	}

	renderedPrompt := e.promptLoader.Render(prompt, map[string]string{
		"query": opts.Query,
		"mode":  mode,
	})

	// Check context again
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	// Send progress: Querying provider
	if progress != nil {
		progress <- "Querying AI provider..."
	}

	// Query the provider
	response, err := e.providerManager.Query(ctx, renderedPrompt, provider.QueryOptions{})
	if err != nil {
		return nil, fmt.Errorf("provider query failed: %w", err)
	}

	// Send progress: Processing results
	if progress != nil {
		progress <- "Processing results..."
	}

	duration := time.Since(start)

	// Create result
	result := &ResearchResult{
		Query:    opts.Query,
		Mode:     mode,
		Content:  response.Content,
		Duration: duration,
	}

	// Store in database if not disabled
	if !opts.NoStore {
		if progress != nil {
			progress <- "Storing in database..."
		}

		session := &db.ResearchSession{
			Query:      opts.Query,
			Mode:       mode,
			PromptUsed: promptName,
			Result:     response.Content,
			CreatedAt:  time.Now(),
		}

		if err := e.db.SaveSession(session); err != nil {
			// Don't fail the entire operation if storage fails
			// Just log and continue
			if progress != nil {
				progress <- fmt.Sprintf("Warning: Failed to store session: %v", err)
			}
		} else {
			result.SessionID = session.ID
		}
	}

	// Send completion progress
	if progress != nil {
		progress <- "Complete!"
	}

	return result, nil
}
