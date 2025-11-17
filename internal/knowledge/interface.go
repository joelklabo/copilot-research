package knowledge

// KnowledgeManagerInterface defines the interface for KnowledgeManager operations.
type KnowledgeManagerInterface interface {
	Add(k *Knowledge) error
	Update(id string, k *Knowledge) error
	Get(id string) (*Knowledge, error)
	Delete(id string) error
	List() ([]*Knowledge, error)
	Search(query string) ([]*Knowledge, error)
	Deduplicate(topicPrefix string) error
	Consolidate() error
	GetRelevantKnowledge(query string, maxSize int) (string, error)
	History(topic string) ([]GitCommit, error)
	Diff(from, to string) (string, error)
	Commit(message string) error
}
