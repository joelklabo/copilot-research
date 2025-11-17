package knowledge

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// KnowledgeManager handles CRUD operations, Git versioning, and consolidation
type KnowledgeManager struct {
	baseDir string
	cache   map[string]*Knowledge
	mu      sync.RWMutex
}

// GitCommit represents a git commit entry
type GitCommit struct {
	Hash      string
	Author    string
	Date      time.Time
	Message   string
}

// NewKnowledgeManager creates a new knowledge manager and initializes git repo
func NewKnowledgeManager(baseDir string) (*KnowledgeManager, error) {
	// Ensure directory exists
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create knowledge directory: %w", err)
	}

	km := &KnowledgeManager{
		baseDir: baseDir,
		cache:   make(map[string]*Knowledge),
	}

	// Initialize git repo if not exists
	gitDir := filepath.Join(baseDir, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		if err := km.initGit(); err != nil {
			return nil, fmt.Errorf("failed to initialize git: %w", err)
		}
	}

	// Load existing knowledge into cache
	if err := km.loadCache(); err != nil {
		return nil, fmt.Errorf("failed to load cache: %w", err)
	}

	return km, nil
}

// initGit initializes a git repository
func (km *KnowledgeManager) initGit() error {
	cmd := exec.Command("git", "init")
	cmd.Dir = km.baseDir
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git init failed: %w, output: %s", err, output)
	}

	// Configure git
	commands := [][]string{
		{"git", "config", "user.name", "Copilot Research"},
		{"git", "config", "user.email", "research@copilot.local"},
	}

	for _, args := range commands {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = km.baseDir
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("git config failed: %w", err)
		}
	}

	return nil
}

// loadCache loads all knowledge files into memory
func (km *KnowledgeManager) loadCache() error {
	return filepath.Walk(km.baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || !strings.HasSuffix(path, ".md") {
			return nil
		}

		k, err := ParseKnowledge(path)
		if err != nil {
			// Skip files that can't be parsed
			return nil
		}

		km.mu.Lock()
		km.cache[k.Topic] = k
		km.mu.Unlock()

		return nil
	})
}

// Add adds new knowledge and commits to git
func (km *KnowledgeManager) Add(k *Knowledge) error {
	km.mu.Lock()
	defer km.mu.Unlock()

	// Set metadata
	if k.CreatedAt.IsZero() {
		k.CreatedAt = time.Now()
	}
	k.UpdatedAt = time.Now()
	k.Version = 1

	// Write to file
	filename := km.getFilePath(k.Topic)
	if err := k.Save(filename); err != nil {
		return fmt.Errorf("failed to save knowledge: %w", err)
	}

	// Update cache
	km.cache[k.Topic] = k

	// Commit to git
	message := fmt.Sprintf("Add: %s - %s", k.Topic, truncate(k.Content, 50))
	if err := km.commit(filename, message); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}

// Update updates existing knowledge
func (km *KnowledgeManager) Update(id string, k *Knowledge) error {
	km.mu.Lock()
	defer km.mu.Unlock()

	existing, exists := km.cache[id]
	if !exists {
		return fmt.Errorf("knowledge not found: %s", id)
	}

	// Increment version
	k.Version = existing.Version + 1
	k.CreatedAt = existing.CreatedAt
	k.UpdatedAt = time.Now()
	k.Topic = id

	// Write to file
	filename := km.getFilePath(id)
	if err := k.Save(filename); err != nil {
		return fmt.Errorf("failed to save knowledge: %w", err)
	}

	// Update cache
	km.cache[id] = k

	// Commit to git
	message := fmt.Sprintf("Update: %s - %s", id, truncate(k.Content, 50))
	if err := km.commit(filename, message); err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	return nil
}

// Get retrieves knowledge by topic
func (km *KnowledgeManager) Get(id string) (*Knowledge, error) {
	km.mu.RLock()
	defer km.mu.RUnlock()

	k, exists := km.cache[id]
	if !exists {
		return nil, fmt.Errorf("knowledge not found: %s", id)
	}

	return k, nil
}

// Delete removes knowledge
func (km *KnowledgeManager) Delete(id string) error {
	km.mu.Lock()
	defer km.mu.Unlock()

	if _, exists := km.cache[id]; !exists {
		return fmt.Errorf("knowledge not found: %s", id)
	}

	// Remove file
	filename := km.getFilePath(id)
	if err := os.Remove(filename); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove file: %w", err)
	}

	// Remove from cache
	delete(km.cache, id)

	// Commit to git
	message := fmt.Sprintf("Remove: %s", id)
	if err := km.commitDeletion(filename, message); err != nil {
		return fmt.Errorf("failed to commit deletion: %w", err)
	}

	return nil
}

// List returns all knowledge entries
func (km *KnowledgeManager) List() ([]*Knowledge, error) {
	km.mu.RLock()
	defer km.mu.RUnlock()

	list := make([]*Knowledge, 0, len(km.cache))
	for _, k := range km.cache {
		list = append(list, k)
	}

	return list, nil
}

// Search finds knowledge entries matching query
func (km *KnowledgeManager) Search(query string) ([]*Knowledge, error) {
	km.mu.RLock()
	defer km.mu.RUnlock()

	query = strings.ToLower(query)
	results := make([]*Knowledge, 0)

	for _, k := range km.cache {
		if strings.Contains(strings.ToLower(k.Topic), query) ||
			strings.Contains(strings.ToLower(k.Content), query) ||
			containsTag(k.Tags, query) {
			results = append(results, k)
		}
	}

	return results, nil
}

// Deduplicate removes duplicate or very similar entries
func (km *KnowledgeManager) Deduplicate(topicPrefix string) error {
	km.mu.Lock()
	defer km.mu.Unlock()

	// Find all entries matching prefix
	candidates := make([]*Knowledge, 0)
	for _, k := range km.cache {
		if strings.HasPrefix(k.Topic, topicPrefix) {
			candidates = append(candidates, k)
		}
	}

	if len(candidates) < 2 {
		return nil // Nothing to deduplicate
	}

	// Simple deduplication: keep highest confidence, newest version
	toRemove := make(map[string]bool)
	for i := 0; i < len(candidates); i++ {
		if toRemove[candidates[i].Topic] {
			continue // Already marked for removal
		}
		for j := i + 1; j < len(candidates); j++ {
			if toRemove[candidates[j].Topic] {
				continue // Already marked for removal
			}
			similarity := calculateSimilarity(candidates[i].Content, candidates[j].Content)
			if similarity > 0.85 { // Lower threshold to actually find duplicates
				// Keep the one with higher confidence or newer
				var remove string
				if candidates[i].Confidence > candidates[j].Confidence {
					remove = candidates[j].Topic
				} else if candidates[i].Confidence < candidates[j].Confidence {
					remove = candidates[i].Topic
				} else if candidates[i].UpdatedAt.After(candidates[j].UpdatedAt) {
					remove = candidates[j].Topic
				} else {
					remove = candidates[i].Topic
				}
				toRemove[remove] = true
			}
		}
	}

	// Remove duplicates
	for topic := range toRemove {
		filename := km.getFilePath(topic)
		os.Remove(filename)
		delete(km.cache, topic)
	}

	if len(toRemove) > 0 {
		message := fmt.Sprintf("Deduplicate: Removed %d duplicate entries in %s", len(toRemove), topicPrefix)
		if err := km.commitAll(message); err != nil {
			return fmt.Errorf("failed to commit deduplication: %w", err)
		}
	}

	return nil
}

// Consolidate performs cleanup and optimization
func (km *KnowledgeManager) Consolidate() error {
	km.mu.Lock()
	defer km.mu.Unlock()

	// Group by topic prefix (first part before /)
	groups := make(map[string][]*Knowledge)
	for _, k := range km.cache {
		prefix := strings.Split(k.Topic, "/")[0]
		groups[prefix] = append(groups[prefix], k)
	}

	consolidated := false
	for _, entries := range groups {
		if len(entries) > 1 {
			// Simple consolidation: merge similar entries
			// This is a placeholder for more sophisticated logic
			consolidated = true
		}
	}

	if consolidated {
		message := "Consolidate: Merged and optimized knowledge entries"
		if err := km.commitAll(message); err != nil {
			return fmt.Errorf("failed to commit consolidation: %w", err)
		}
	}

	return nil
}

// GetRelevantKnowledge retrieves knowledge relevant to a query
func (km *KnowledgeManager) GetRelevantKnowledge(query string, maxSize int) (string, error) {
	results, err := km.Search(query)
	if err != nil {
		return "", err
	}

	if len(results) == 0 {
		return "", nil
	}

	var sb strings.Builder
	totalSize := 0

	for _, k := range results {
		content := fmt.Sprintf("## %s\n\n%s\n\n", k.Topic, strings.TrimSpace(k.Content))
		if totalSize+len(content) > maxSize {
			break
		}
		sb.WriteString(content)
		totalSize += len(content)
	}

	return sb.String(), nil
}

// History returns git commit history for a topic
func (km *KnowledgeManager) History(topic string) ([]GitCommit, error) {
	filename := km.getFilePath(topic)
	
	cmd := exec.Command("git", "log", "--pretty=format:%H|%an|%at|%s", "--", filepath.Base(filename))
	cmd.Dir = km.baseDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("git log failed: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	commits := make([]GitCommit, 0, len(lines))

	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) != 4 {
			continue
		}

		var timestamp int64
		_, err := fmt.Sscanf(parts[2], "%d", &timestamp) // Added error check
		if err != nil {
			// Log the error or handle it appropriately, for now, skip this commit
			continue
		}

		commits = append(commits, GitCommit{
			Hash:    parts[0],
				Author:  parts[1],
				Date:    time.Unix(timestamp, 0),
				Message: parts[3],
		})
	}

	return commits, nil
}

// Diff returns the diff between two commits
func (km *KnowledgeManager) Diff(from, to string) (string, error) {
	cmd := exec.Command("git", "diff", from, to)
	cmd.Dir = km.baseDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("git diff failed: %w", err)
	}
	return string(output), nil
}

// commit commits a single file to git
func (km *KnowledgeManager) commit(filename, message string) error {
	commands := [][]string{
		{"git", "add", filepath.Base(filename)},
		{"git", "commit", "-m", message},
	}

	for _, args := range commands {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = km.baseDir
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("git command failed: %w, output: %s", err, output)
		}
	}

	return nil
}

// commitDeletion commits a file deletion
func (km *KnowledgeManager) commitDeletion(filename, message string) error {
	commands := [][]string{
		{"git", "rm", filepath.Base(filename)},
		{"git", "commit", "-m", message},
	}

	for _, args := range commands {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = km.baseDir
		if output, err := cmd.CombinedOutput(); err != nil {
			// File might already be deleted
			if !strings.Contains(string(output), "did not match any files") {
				return fmt.Errorf("git command failed: %w", err)
			}
		}
	}

	return nil
}

// commitAll commits all changes
func (km *KnowledgeManager) commitAll(message string) error {
	commands := [][]string{
		{"git", "add", "-A"},
		{"git", "commit", "-m", message},
	}

	for _, args := range commands {
		cmd := exec.Command(args[0], args[1:]...)
		cmd.Dir = km.baseDir
		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("git command failed: %w, output: %s", err, output)
		}
	}

	return nil
}

// Commit manually commits changes with a message
func (km *KnowledgeManager) Commit(message string) error {
	return km.commitAll(message)
}

// getFilePath returns the full file path for a topic
func (km *KnowledgeManager) getFilePath(topic string) string {
	// Replace / with - for filesystem safety and remove invalid chars
	safeTopic := strings.ReplaceAll(topic, "/", "-")
	safeTopic = strings.ReplaceAll(safeTopic, " ", "_")
	// Remove other potentially problematic characters
	safeTopic = strings.Map(func(r rune) rune {
		if r == '-' || r == '_' || r == '.' || (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			return r
		}
		return '_'
	}, safeTopic)
	return filepath.Join(km.baseDir, safeTopic+".md")
}

// Helper functions

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

func containsTag(tags []string, query string) bool {
	for _, tag := range tags {
		if strings.Contains(strings.ToLower(tag), query) {
			return true
		}
	}
	return false
}

// calculateSimilarity is a simple similarity metric (0.0 to 1.0)
// This is a placeholder - could use Levenshtein distance or embeddings
func calculateSimilarity(s1, s2 string) float64 {
	s1 = strings.ToLower(s1)
	s2 = strings.ToLower(s2)

	if s1 == s2 {
		return 1.0
	}

	// Simple word overlap metric
	words1 := strings.Fields(s1)
	words2 := strings.Fields(s2)

	if len(words1) == 0 || len(words2) == 0 {
		return 0.0
	}

	wordSet1 := make(map[string]bool)
	for _, w := range words1 {
		wordSet1[w] = true
	}

	overlap := 0
	for _, w := range words2 {
		if wordSet1[w] {
			overlap++
		}
	}

	maxLen := float64(max(len(words1), len(words2)))
	return float64(overlap) / maxLen
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
