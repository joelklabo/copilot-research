package provider

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// AIProvider is the interface that all AI providers must implement
type AIProvider interface {
	// Name returns the provider's unique identifier
	Name() string
	
	// Query sends a prompt to the provider and returns the response
	Query(ctx context.Context, prompt string, opts QueryOptions) (*Response, error)
	
	// IsAuthenticated checks if the provider is properly authenticated
	IsAuthenticated() bool
	
	// RequiresAuth returns authentication information
	RequiresAuth() AuthInfo
	
	// Capabilities returns the provider's capabilities
	Capabilities() ProviderCapabilities
}

// QueryOptions contains options for querying a provider
type QueryOptions struct {
	MaxTokens   int
	Temperature float64
	TopP        float64
	Model       string
	Stream      bool
}

// Response represents the response from a provider
type Response struct {
	Content    string
	Provider   string
	Model      string
	TokensUsed TokenUsage
	Duration   time.Duration
	Metadata   map[string]interface{}
}

// TokenUsage tracks token consumption
type TokenUsage struct {
	Prompt     int
	Completion int
	Total      int
}

// ProviderCapabilities describes what a provider can do
type ProviderCapabilities struct {
	Streaming      bool
	FunctionCall   bool
	MaxTokens      int
	SupportsImages bool
}

// AuthInfo provides authentication information for a provider
type AuthInfo struct {
	Type         string // "oauth", "apikey", "cli"
	IsConfigured bool
	HelpURL      string
	Instructions string
}

// ProviderFactory manages provider instances
type ProviderFactory struct {
	providers map[string]AIProvider
	mu        sync.RWMutex
}

// NewProviderFactory creates a new provider factory
func NewProviderFactory() *ProviderFactory {
	return &ProviderFactory{
		providers: make(map[string]AIProvider),
	}
}

// Register registers a provider with the factory
func (f *ProviderFactory) Register(name string, provider AIProvider) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	
	if _, exists := f.providers[name]; exists {
		return fmt.Errorf("provider '%s' is already registered", name)
	}
	
	f.providers[name] = provider
	return nil
}

// Get retrieves a provider by name
func (f *ProviderFactory) Get(name string) (AIProvider, error) {
	f.mu.RLock()
	defer f.mu.RUnlock()
	
	provider, exists := f.providers[name]
	if !exists {
		return nil, fmt.Errorf("provider '%s' not found", name)
	}
	
	return provider, nil
}

// List returns all registered provider names
func (f *ProviderFactory) List() []string {
	f.mu.RLock()
	defer f.mu.RUnlock()
	
	names := make([]string, 0, len(f.providers))
	for name := range f.providers {
		names = append(names, name)
	}
	
	return names
}

// Unregister removes a provider from the factory
func (f *ProviderFactory) Unregister(name string) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	
	if _, exists := f.providers[name]; !exists {
		return fmt.Errorf("provider '%s' not found", name)
	}
	
	delete(f.providers, name)
	return nil
}

// ProviderManager manages provider selection and fallback logic
type ProviderManager struct {
	factory              *ProviderFactory
	primary              string
	fallback             string
	autoFallback         bool
	notifyFallback       bool
	notificationHandler  func(string)
}

// NewProviderManager creates a new provider manager
// Updated signature to include autoFallback and notifyFallback
func NewProviderManager(factory *ProviderFactory, primary, fallback string, autoFallback, notifyFallback bool) *ProviderManager {
	return &ProviderManager{
		factory:        factory,
		primary:        primary,
		fallback:       fallback,
		autoFallback:   autoFallback,  // Use provided value
		notifyFallback: notifyFallback,  // Use provided value
		notificationHandler: func(msg string) {
			// Default: print to stdout
			fmt.Println(msg)
		},
	}
}

// GetFactory returns the ProviderFactory associated with the manager
func (pm *ProviderManager) GetFactory() *ProviderFactory {
    return pm.factory
}

// SetAutoFallback enables or disables automatic fallback
func (pm *ProviderManager) SetAutoFallback(enabled bool) {
	pm.autoFallback = enabled
}

// SetNotifyFallback enables or disables fallback notifications
func (pm *ProviderManager) SetNotifyFallback(enabled bool) {
	pm.notifyFallback = enabled
}

// SetNotificationHandler sets a custom notification handler
func (pm *ProviderManager) SetNotificationHandler(handler func(string)) {
	pm.notificationHandler = handler
}

// Query attempts to query the primary provider, falling back if it fails
func (pm *ProviderManager) Query(ctx context.Context, prompt string, opts QueryOptions) (*Response, error) {
	// Try primary provider
	if pm.primary != "" {
		provider, err := pm.factory.Get(pm.primary)
		if err == nil && provider.IsAuthenticated() {
			resp, err := provider.Query(ctx, prompt, opts)
			if err == nil {
				return resp, nil
			}
			// Primary failed, log it
		}
	}
	
	// Try fallback provider if auto-fallback is enabled
	if pm.autoFallback && pm.fallback != "" {
		provider, err := pm.factory.Get(pm.fallback)
		if err == nil && provider.IsAuthenticated() {
			// Notify user about fallback
			if pm.notifyFallback && pm.notificationHandler != nil {
				pm.notificationHandler(fmt.Sprintf("ℹ️  Using %s (primary unavailable)", pm.fallback))
			}
			
			resp, err := provider.Query(ctx, prompt, opts)
			if err == nil {
				return resp, nil
			}
		}
	}
	
	// All providers failed
	return nil, fmt.Errorf("all providers failed: primary=%s, fallback=%s", pm.primary, pm.fallback)
}

// CheckAuthentication returns lists of authenticated and unauthenticated providers
func (pm *ProviderManager) CheckAuthentication() (authenticated []string, unauthenticated []string) {
	authenticated = make([]string, 0)
	unauthenticated = make([]string, 0)
	
	for _, name := range pm.factory.List() {
		provider, err := pm.factory.Get(name)
		if err != nil {
			continue
		}
		
		if provider.IsAuthenticated() {
			authenticated = append(authenticated, name)
		} else {
			unauthenticated = append(unauthenticated, name)
		}
	}
	
	return authenticated, unauthenticated
}

// SetPrimary sets the primary provider
func (pm *ProviderManager) SetPrimary(name string) {
	pm.primary = name
}

// SetFallback sets the fallback provider
func (pm *ProviderManager) SetFallback(name string) {
	pm.fallback = name
}

// GetPrimary returns the primary provider name
func (pm *ProviderManager) GetPrimary() string {
	return pm.primary
}

// GetFallback returns the fallback provider name
func (pm *ProviderManager) GetFallback() string {
	return pm.fallback
}