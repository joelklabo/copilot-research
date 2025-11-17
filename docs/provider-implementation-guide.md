# Provider Implementation Guide

This guide explains how to implement a new AI provider for the Copilot Research CLI. It's designed to be read and followed by AI agents (like Claude or Copilot) to implement providers without human intervention.

## Overview

The provider abstraction layer allows Copilot Research to support multiple AI backends (GitHub Copilot, OpenAI, Anthropic Claude, etc.) through a unified interface. All providers implement the `AIProvider` interface and are managed by the `ProviderFactory`.

## Core Architecture

### 1. AIProvider Interface

Every provider must implement this interface (located in `internal/provider/provider.go`):

```go
type AIProvider interface {
    // Name returns the provider's unique identifier
    // Example: "github-copilot", "openai", "anthropic"
    Name() string
    
    // Query sends a prompt and returns a response
    // Must respect context cancellation and timeout
    Query(ctx context.Context, prompt string, opts QueryOptions) (*Response, error)
    
    // IsAuthenticated checks if credentials are valid
    // Should be fast (< 1 second) - cache auth status if needed
    IsAuthenticated() bool
    
    // RequiresAuth returns authentication information
    // Used to guide users through setup
    RequiresAuth() AuthInfo
    
    // Capabilities describes what the provider supports
    // Used for feature detection and routing
    Capabilities() ProviderCapabilities
}
```

### 2. Supporting Types

```go
// QueryOptions - Configuration for each query
type QueryOptions struct {
    MaxTokens   int     // Maximum tokens to generate
    Temperature float64 // Randomness (0.0 = deterministic, 2.0 = very random)
    TopP        float64 // Nucleus sampling parameter
    Model       string  // Specific model to use (if provider supports multiple)
    Stream      bool    // Whether to stream responses
}

// Response - Standardized response format
type Response struct {
    Content    string                 // The actual response text
    Provider   string                 // Provider name that generated it
    Model      string                 // Model that was used
    TokensUsed TokenUsage            // Token consumption
    Duration   time.Duration          // How long the query took
    Metadata   map[string]interface{} // Provider-specific data
}

// TokenUsage - Track token consumption
type TokenUsage struct {
    Prompt     int // Tokens in the prompt
    Completion int // Tokens in the completion
    Total      int // Total tokens used
}

// ProviderCapabilities - What the provider can do
type ProviderCapabilities struct {
    Streaming      bool // Can stream responses
    FunctionCall   bool // Supports function calling
    MaxTokens      int  // Maximum tokens supported
    SupportsImages bool // Can process images
}

// AuthInfo - Authentication guidance
type AuthInfo struct {
    Type         string // "oauth", "apikey", "cli"
    IsConfigured bool   // Whether auth is set up
    HelpURL      string // Documentation URL
    Instructions string // Step-by-step setup guide
}
```

## Implementation Steps

### Step 1: Create Provider File

Create a new file: `internal/provider/{provider_name}.go`

Example: `internal/provider/openai.go`

### Step 2: Define Provider Struct

```go
package provider

import (
    "context"
    "fmt"
    "os"
    "time"
)

// OpenAIProvider implements AIProvider for OpenAI
type OpenAIProvider struct {
    apiKey     string
    timeout    time.Duration
    baseURL    string
    model      string
}

// NewOpenAIProvider creates a new OpenAI provider
func NewOpenAIProvider(timeout time.Duration) *OpenAIProvider {
    return &OpenAIProvider{
        timeout: timeout,
        baseURL: "https://api.openai.com/v1",
        model:   "gpt-4", // Default model
    }
}
```

### Step 3: Implement Name()

```go
func (o *OpenAIProvider) Name() string {
    return "openai"
}
```

### Step 4: Implement IsAuthenticated()

Check credentials in order of priority:

1. Environment variables
2. Configuration files
3. CLI tools
4. System keychain

```go
func (o *OpenAIProvider) IsAuthenticated() bool {
    // Check environment variable
    if apiKey := os.Getenv("OPENAI_API_KEY"); apiKey != "" {
        o.apiKey = apiKey
        return true
    }
    
    // Check config file
    if apiKey := o.loadFromConfig(); apiKey != "" {
        o.apiKey = apiKey
        return true
    }
    
    return false
}

func (o *OpenAIProvider) loadFromConfig() string {
    // Load from ~/.copilot-research/config.yaml
    // Return API key if found
    return ""
}
```

### Step 5: Implement RequiresAuth()

Provide clear, actionable instructions:

```go
func (o *OpenAIProvider) RequiresAuth() AuthInfo {
    if o.IsAuthenticated() {
        return AuthInfo{
            Type:         "apikey",
            IsConfigured: true,
        }
    }
    
    return AuthInfo{
        Type:         "apikey",
        IsConfigured: false,
        HelpURL:      "https://platform.openai.com/api-keys",
        Instructions: `OpenAI API authentication required.

Please set your API key using one of these methods:

1. Environment variable (recommended):
   export OPENAI_API_KEY=sk-your-key-here
   
2. Configuration file:
   copilot-research config set openai.api_key sk-your-key-here

Get your API key at: https://platform.openai.com/api-keys

Note: You need an OpenAI account with API access.

Once configured, run your command again.`,
    }
}
```

### Step 6: Implement Capabilities()

Be honest about what your provider supports:

```go
func (o *OpenAIProvider) Capabilities() ProviderCapabilities {
    return ProviderCapabilities{
        Streaming:      true,  // OpenAI supports streaming
        FunctionCall:   true,  // OpenAI supports function calling
        MaxTokens:      128000, // GPT-4 context window
        SupportsImages: true,  // GPT-4 Vision can process images
    }
}
```

### Step 7: Implement Query()

This is the core method. Must:
- Respect context cancellation
- Handle timeouts
- Return errors with helpful messages
- Parse provider response into standard format

```go
func (o *OpenAIProvider) Query(ctx context.Context, prompt string, opts QueryOptions) (*Response, error) {
    // Check authentication
    if !o.IsAuthenticated() {
        return nil, fmt.Errorf("not authenticated: please set OPENAI_API_KEY")
    }
    
    // Create context with timeout
    queryCtx, cancel := context.WithTimeout(ctx, o.timeout)
    defer cancel()
    
    // Prepare request
    start := time.Now()
    reqBody := map[string]interface{}{
        "model": o.model,
        "messages": []map[string]string{
            {"role": "user", "content": prompt},
        },
    }
    
    // Apply options
    if opts.MaxTokens > 0 {
        reqBody["max_tokens"] = opts.MaxTokens
    }
    if opts.Temperature > 0 {
        reqBody["temperature"] = opts.Temperature
    }
    
    // Make API call
    // (Use your preferred HTTP client - net/http, resty, etc.)
    response, err := o.makeRequest(queryCtx, reqBody)
    duration := time.Since(start)
    
    if err != nil {
        // Handle common errors with helpful messages
        if queryCtx.Err() == context.DeadlineExceeded {
            return nil, fmt.Errorf("query timeout after %v", o.timeout)
        }
        
        // Check for auth errors
        if isAuthError(err) {
            return nil, fmt.Errorf("authentication failed: invalid API key")
        }
        
        // Check for rate limiting
        if isRateLimitError(err) {
            return nil, fmt.Errorf("rate limit exceeded: please try again later")
        }
        
        return nil, fmt.Errorf("OpenAI API error: %w", err)
    }
    
    // Parse and return response
    return &Response{
        Content:  response.Content,
        Provider: "openai",
        Model:    response.Model,
        Duration: duration,
        TokensUsed: TokenUsage{
            Prompt:     response.Usage.PromptTokens,
            Completion: response.Usage.CompletionTokens,
            Total:      response.Usage.TotalTokens,
        },
        Metadata: map[string]interface{}{
            "finish_reason": response.FinishReason,
        },
    }, nil
}
```

## Testing Your Provider

### Step 1: Create Test File

Create `internal/provider/{provider_name}_test.go`

### Step 2: Write Comprehensive Tests

```go
package provider

import (
    "context"
    "os"
    "testing"
    "time"

    "github.com/stretchr/testify/assert"
)

func TestNewOpenAIProvider(t *testing.T) {
    provider := NewOpenAIProvider(30 * time.Second)
    assert.NotNil(t, provider)
    assert.Equal(t, "openai", provider.Name())
}

func TestOpenAIProvider_Capabilities(t *testing.T) {
    provider := NewOpenAIProvider(30 * time.Second)
    
    caps := provider.Capabilities()
    assert.True(t, caps.Streaming)
    assert.True(t, caps.FunctionCall)
    assert.Equal(t, 128000, caps.MaxTokens)
    assert.True(t, caps.SupportsImages)
}

func TestOpenAIProvider_IsAuthenticated_WithAPIKey(t *testing.T) {
    os.Setenv("OPENAI_API_KEY", "sk-test-key")
    defer os.Unsetenv("OPENAI_API_KEY")
    
    provider := NewOpenAIProvider(30 * time.Second)
    assert.True(t, provider.IsAuthenticated())
}

func TestOpenAIProvider_IsAuthenticated_NoAPIKey(t *testing.T) {
    os.Unsetenv("OPENAI_API_KEY")
    
    provider := NewOpenAIProvider(30 * time.Second)
    
    // May be true if configured elsewhere, test accordingly
    if !provider.IsAuthenticated() {
        authInfo := provider.RequiresAuth()
        assert.False(t, authInfo.IsConfigured)
        assert.NotEmpty(t, authInfo.Instructions)
        assert.Contains(t, authInfo.Instructions, "OPENAI_API_KEY")
    }
}

func TestOpenAIProvider_Query_NotAuthenticated(t *testing.T) {
    os.Unsetenv("OPENAI_API_KEY")
    
    provider := NewOpenAIProvider(30 * time.Second)
    
    if provider.IsAuthenticated() {
        t.Skip("Skipping - provider is authenticated")
    }
    
    ctx := context.Background()
    opts := QueryOptions{}
    
    _, err := provider.Query(ctx, "test", opts)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "not authenticated")
}

// Note: Testing actual API calls requires:
// 1. Valid API key
// 2. Active account
// 3. Network connectivity
// These should be integration tests, not unit tests
```

### Step 3: Run Tests

```bash
go test ./internal/provider/... -v -run="TestOpenAI"
```

## Registration and Usage

### Registering Your Provider

In your application initialization code:

```go
func initProviders() *provider.ProviderFactory {
    factory := provider.NewProviderFactory()
    
    // Register GitHub Copilot
    ghProvider := provider.NewGitHubCopilotProvider(60 * time.Second)
    factory.Register("github-copilot", ghProvider)
    
    // Register OpenAI
    openaiProvider := provider.NewOpenAIProvider(30 * time.Second)
    factory.Register("openai", openaiProvider)
    
    // Register more providers...
    
    return factory
}
```

### Using the Provider Manager

For automatic fallback:

```go
func main() {
    factory := initProviders()
    
    // Create manager with primary and fallback
    manager := provider.NewProviderManager(factory, "github-copilot", "openai")
    
    // Query will try github-copilot first, fall back to openai
    ctx := context.Background()
    opts := provider.QueryOptions{
        MaxTokens:   2000,
        Temperature: 0.7,
    }
    
    resp, err := manager.Query(ctx, "How do Swift actors work?", opts)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(resp.Content)
}
```

## Best Practices

### 1. Authentication

- **Check credentials fast**: IsAuthenticated() should be < 1 second
- **Cache auth status**: Don't validate on every call
- **Priority order**: Env vars > config files > CLI tools
- **Clear errors**: Tell users exactly how to fix auth issues

### 2. Error Handling

- **Respect context**: Always check `ctx.Done()` and `ctx.Err()`
- **Timeout gracefully**: Return helpful timeout messages
- **Parse provider errors**: Convert API errors to user-friendly messages
- **Rate limiting**: Detect and report rate limit errors clearly

### 3. Response Parsing

- **Clean output**: Remove markdown wrappers, API artifacts
- **Preserve formatting**: Keep code blocks, lists, headers intact
- **Token counting**: If provider gives usage, use it; otherwise estimate
- **Metadata**: Store provider-specific data in Response.Metadata

### 4. Configuration

- **Sensible defaults**: Use common models (gpt-4, claude-3-5-sonnet)
- **Allow overrides**: Respect QueryOptions for per-query customization
- **Document limits**: Be clear about MaxTokens, rate limits, etc.

## Example: Complete Anthropic Provider

```go
package provider

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "os"
    "time"
)

type AnthropicProvider struct {
    apiKey  string
    timeout time.Duration
    baseURL string
    model   string
    client  *http.Client
}

func NewAnthropicProvider(timeout time.Duration) *AnthropicProvider {
    return &AnthropicProvider{
        timeout: timeout,
        baseURL: "https://api.anthropic.com/v1",
        model:   "claude-3-5-sonnet-20241022",
        client:  &http.Client{Timeout: timeout},
    }
}

func (a *AnthropicProvider) Name() string {
    return "anthropic"
}

func (a *AnthropicProvider) IsAuthenticated() bool {
    if apiKey := os.Getenv("ANTHROPIC_API_KEY"); apiKey != "" {
        a.apiKey = apiKey
        return true
    }
    return false
}

func (a *AnthropicProvider) RequiresAuth() AuthInfo {
    if a.IsAuthenticated() {
        return AuthInfo{
            Type:         "apikey",
            IsConfigured: true,
        }
    }
    
    return AuthInfo{
        Type:         "apikey",
        IsConfigured: false,
        HelpURL:      "https://console.anthropic.com/settings/keys",
        Instructions: `Anthropic Claude API authentication required.

Please set your API key:

1. Environment variable (recommended):
   export ANTHROPIC_API_KEY=sk-ant-your-key-here

2. Configuration file:
   copilot-research config set anthropic.api_key sk-ant-your-key-here

Get your API key at: https://console.anthropic.com/settings/keys

Once configured, run your command again.`,
    }
}

func (a *AnthropicProvider) Capabilities() ProviderCapabilities {
    return ProviderCapabilities{
        Streaming:      true,
        FunctionCall:   true,
        MaxTokens:      200000, // Claude 3.5 context window
        SupportsImages: true,
    }
}

func (a *AnthropicProvider) Query(ctx context.Context, prompt string, opts QueryOptions) (*Response, error) {
    if !a.IsAuthenticated() {
        return nil, fmt.Errorf("not authenticated: please set ANTHROPIC_API_KEY")
    }
    
    start := time.Now()
    
    // Build request
    reqBody := map[string]interface{}{
        "model": a.model,
        "messages": []map[string]string{
            {"role": "user", "content": prompt},
        },
        "max_tokens": 4096,
    }
    
    if opts.MaxTokens > 0 {
        reqBody["max_tokens"] = opts.MaxTokens
    }
    if opts.Temperature > 0 {
        reqBody["temperature"] = opts.Temperature
    }
    
    jsonData, err := json.Marshal(reqBody)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request: %w", err)
    }
    
    // Create HTTP request
    req, err := http.NewRequestWithContext(ctx, "POST", a.baseURL+"/messages", bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }
    
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("X-API-Key", a.apiKey)
    req.Header.Set("Anthropic-Version", "2023-06-01")
    
    // Execute request
    resp, err := a.client.Do(req)
    if err != nil {
        if ctx.Err() == context.DeadlineExceeded {
            return nil, fmt.Errorf("query timeout after %v", a.timeout)
        }
        return nil, fmt.Errorf("API request failed: %w", err)
    }
    defer resp.Body.Close()
    
    duration := time.Since(start)
    
    // Read response
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response: %w", err)
    }
    
    // Check status code
    if resp.StatusCode != 200 {
        return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
    }
    
    // Parse response
    var apiResp struct {
        Content []struct {
            Text string `json:"text"`
        } `json:"content"`
        Model string `json:"model"`
        Usage struct {
            InputTokens  int `json:"input_tokens"`
            OutputTokens int `json:"output_tokens"`
        } `json:"usage"`
    }
    
    if err := json.Unmarshal(body, &apiResp); err != nil {
        return nil, fmt.Errorf("failed to parse response: %w", err)
    }
    
    // Extract content
    content := ""
    if len(apiResp.Content) > 0 {
        content = apiResp.Content[0].Text
    }
    
    return &Response{
        Content:  content,
        Provider: "anthropic",
        Model:    apiResp.Model,
        Duration: duration,
        TokensUsed: TokenUsage{
            Prompt:     apiResp.Usage.InputTokens,
            Completion: apiResp.Usage.OutputTokens,
            Total:      apiResp.Usage.InputTokens + apiResp.Usage.OutputTokens,
        },
    }, nil
}
```

## Troubleshooting

### Common Issues

1. **Tests fail due to local auth**: Use `t.Skip()` if provider is already authenticated
2. **Timeout errors**: Ensure context is properly passed through
3. **Memory leaks**: Always close response bodies, defer cancel functions
4. **Rate limiting**: Implement exponential backoff for retries

### Debug Checklist

- [ ] Provider registered in factory
- [ ] Name() returns unique identifier
- [ ] IsAuthenticated() checks all auth methods
- [ ] RequiresAuth() provides clear instructions
- [ ] Capabilities() accurately reflects provider
- [ ] Query() respects context cancellation
- [ ] Query() returns Response in standard format
- [ ] Errors are clear and actionable
- [ ] Tests cover all auth methods
- [ ] Tests handle timeout scenarios

## Summary

To implement a new provider:

1. Create `internal/provider/{name}.go`
2. Implement all AIProvider interface methods
3. Create comprehensive tests
4. Register in factory
5. Document in README.md

The provider system is designed to be simple and predictable. Follow the patterns shown in `github_copilot.go` and you'll have a working provider in under an hour.

For questions or examples, see:
- `internal/provider/github_copilot.go` - Complete CLI-based provider
- `internal/provider/provider.go` - Interface definitions
- `internal/provider/provider_test.go` - Mock provider example
