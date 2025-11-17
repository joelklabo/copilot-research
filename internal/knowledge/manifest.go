package knowledge

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

// LoadManifest loads the MANIFEST.yaml file
func LoadManifest(baseDir string) (*Manifest, error) {
	manifestPath := filepath.Join(baseDir, "MANIFEST.yaml")
	
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read manifest: %w", err)
	}

	var manifest Manifest
	if err := yaml.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse manifest: %w", err)
	}

	return &manifest, nil
}

// SaveManifest saves the MANIFEST.yaml file
func SaveManifest(baseDir string, manifest *Manifest) error {
	manifest.Updated = time.Now()
	manifest.Metadata.TotalTopics = len(manifest.Topics)

	manifestPath := filepath.Join(baseDir, "MANIFEST.yaml")
	
	data, err := yaml.Marshal(manifest)
	if err != nil {
		return fmt.Errorf("failed to marshal manifest: %w", err)
	}

	if err := os.WriteFile(manifestPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write manifest: %w", err)
	}

	return nil
}

// AddTopicToManifest adds a new topic to the manifest
func AddTopicToManifest(baseDir string, topic ManifestTopic) error {
	manifest, err := LoadManifest(baseDir)
	if err != nil {
		return err
	}

	// Check if topic already exists
	for i, t := range manifest.Topics {
		if t.Name == topic.Name {
			// Update existing
			manifest.Topics[i] = topic
			return SaveManifest(baseDir, manifest)
		}
	}

	// Add new topic
	manifest.Topics = append(manifest.Topics, topic)
	return SaveManifest(baseDir, manifest)
}

// RemoveTopicFromManifest removes a topic from the manifest
func RemoveTopicFromManifest(baseDir string, topicName string) error {
	manifest, err := LoadManifest(baseDir)
	if err != nil {
		return err
	}

	// Find and remove topic
	for i, t := range manifest.Topics {
		if t.Name == topicName {
			manifest.Topics = append(manifest.Topics[:i], manifest.Topics[i+1:]...)
			return SaveManifest(baseDir, manifest)
		}
	}

	return fmt.Errorf("topic %s not found in manifest", topicName)
}

// GetTopicFromManifest retrieves a topic from the manifest
func GetTopicFromManifest(baseDir string, topicName string) (*ManifestTopic, error) {
	manifest, err := LoadManifest(baseDir)
	if err != nil {
		return nil, err
	}

	for _, t := range manifest.Topics {
		if t.Name == topicName {
			return &t, nil
		}
	}

	return nil, fmt.Errorf("topic %s not found", topicName)
}
