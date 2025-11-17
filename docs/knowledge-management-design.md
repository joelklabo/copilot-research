# Knowledge Management System Design

## Problem Statement
We need a system to:
1. Prevent duplicate information in the knowledge base
2. Ensure consistency across all knowledge entries
3. Allow users to edit/remove/modify knowledge
4. Keep the knowledge base at a manageable size without losing information
5. Enable both automated and manual curation

## Three Design Approaches

---

## Option A: Rule-Based Knowledge Filter with Declarative Config

### Concept
A `.copilot-research/preferences.yaml` file that declares rules, preferences, and filters that are applied both during research and during knowledge consolidation.

### API Examples

```bash
# Add a preference (anti-pattern)
copilot-research preferences add --type exclude --topic "MVVM" --reason "We use MV architecture"

# Add a preference (preferred pattern)
copilot-research preferences add --type prefer --topic "MV architecture" --weight high

# Add a context rule
copilot-research preferences add --type context --key "ios-version" --value "26"

# List all preferences
copilot-research preferences list

# Remove a preference
copilot-research preferences remove --id 3

# Consolidate knowledge with current preferences
copilot-research knowledge consolidate
```

### Preferences File Structure
```yaml
version: 1.0
preferences:
  exclude:
    - topic: "MVVM"
      reason: "We use MV architecture instead"
      added: "2024-11-17"
    - topic: "UIKit"
      except: ["legacy migration", "interop"]
      reason: "SwiftUI-first approach"
      
  prefer:
    - topic: "MV architecture"
      weight: high
      keywords: ["Model", "View", "SwiftUI state"]
    - topic: "Swift 6 concurrency"
      weight: high
      
  context:
    ios_version: "26"
    macos_version: "26"
    swift_version: "6"
    target_platforms: ["iOS", "macOS", "iPadOS"]
    
  consolidation:
    max_size_kb: 500
    deduplication_threshold: 0.85  # 85% similarity = duplicate
    auto_consolidate: true
    consolidate_after_entries: 10
```

### How It Works
1. **During Research**: Preferences are injected into the research prompt
   - "Never mention MVVM, we use MV architecture"
   - "Prioritize Swift 6 concurrency patterns"
2. **During Consolidation**: AI agent reviews knowledge base with preferences
   - Removes entries matching exclude rules
   - Merges similar entries
   - Prioritizes prefer topics
3. **Knowledge Base Structure**: Single consolidated markdown with sections

### Pros
- Declarative and portable (can share preferences via git)
- Easy to understand and modify
- Preferences survive across sessions
- Can be version controlled
- Clear separation of concerns

### Cons
- Requires parsing YAML
- Need to build preference management CLI
- Preferences might drift from actual knowledge base

---

## Option B: Interactive Knowledge Editor with AI-Assisted Curation

### Concept
An interactive TUI (Terminal User Interface) that shows the knowledge base and lets you edit, merge, or delete entries with AI assistance.

### API Examples

```bash
# Open interactive knowledge editor
copilot-research knowledge edit

# Or directly manipulate from CLI
copilot-research knowledge remove --topic "MVVM"
copilot-research knowledge merge --ids 5,6,7 --into "Swift Concurrency Best Practices"
copilot-research knowledge summarize --topic "SwiftUI" --max-tokens 500
```

### Interactive TUI Flow
```
┌─ Copilot Research Knowledge Editor ───────────────────────────┐
│                                                                │
│  📚 Knowledge Base (127 entries, 2.3 MB)                      │
│                                                                │
│  ▼ Swift Concurrency (23 entries)                            │
│    → Actors and Isolation (8 entries) ..................... 📝 │
│    → Sendable Protocol (7 entries) ........................ 📝 │
│    → MainActor Usage (8 entries) [DUPLICATES DETECTED] .... ⚠️  │
│                                                                │
│  ▼ SwiftUI Architecture (15 entries)                          │
│    → MV Pattern (9 entries) ................................ 📝 │
│    → MVVM Pattern (6 entries) [MARKED FOR REMOVAL] ........ ❌ │
│                                                                │
│  Commands:                                                     │
│  [e]dit [m]erge [d]elete [s]ummarize [f]ind [c]onsolidate [q] │
└────────────────────────────────────────────────────────────────┘
```

### When You Select an Entry
```
┌─ Edit Entry: MVVM Pattern ─────────────────────────────────────┐
│                                                                │
│  MVVM (Model-View-ViewModel) is a design pattern...           │
│  [Full content displayed]                                      │
│                                                                │
│  AI Suggestions:                                               │
│  💡 This contradicts your MV architecture preference          │
│  💡 Similar content found in entries #12, #34                 │
│  💡 Referenced by 3 other entries                             │
│                                                                │
│  Actions:                                                      │
│  [d] Delete entirely                                           │
│  [r] Replace with AI summary of MV architecture               │
│  [m] Merge with related entries                               │
│  [e] Edit manually                                            │
│  [k] Keep as-is                                               │
└────────────────────────────────────────────────────────────────┘
```

### How It Works
1. **Knowledge Base Format**: Structured JSON with metadata
```json
{
  "entries": [
    {
      "id": "kb-001",
      "topic": "Swift Concurrency Actors",
      "content": "...",
      "sources": ["url1", "url2"],
      "added": "2024-11-17",
      "last_used": "2024-11-17",
      "references": ["kb-002", "kb-005"],
      "tags": ["swift", "concurrency", "actors"],
      "size_bytes": 2048,
      "quality_score": 0.92
    }
  ],
  "metadata": {
    "total_entries": 127,
    "last_consolidated": "2024-11-16",
    "version": "1.0"
  }
}
```

2. **AI-Assisted Actions**:
   - Detect duplicates using embeddings/similarity
   - Suggest merges
   - Auto-generate summaries
   - Find contradictions

3. **CLI Commands** call AI with specific instructions:
```bash
# Under the hood, this calls:
# copilot with prompt: "Review knowledge base, find all MVVM entries, 
# check if they conflict with MV architecture preference, suggest removal"
```

### Pros
- Visual and intuitive
- AI helps identify issues
- Precise control over each entry
- Can see relationships between entries
- Undo/redo capability

### Cons
- More complex to build (requires TUI library)
- Need structured knowledge format
- Interactive (not scriptable)
- Requires more sophisticated data structure

---

## Option C: Smart Knowledge Database with Query Language

### Concept
Knowledge stored in SQLite with full-text search, automatic deduplication, and a query language for management. AI operates as a "knowledge librarian" that can be instructed via natural language or SQL-like queries.

### API Examples

```bash
# Add knowledge preferences as metadata
copilot-research config set architecture.style "MV"
copilot-research config set architecture.exclude "MVVM,MVC,VIPER"

# Query knowledge
copilot-research knowledge query "swift concurrency best practices"
copilot-research knowledge find --tag actors --since "2024-11-01"

# AI-assisted management
copilot-research knowledge ask "Remove all entries about MVVM"
copilot-research knowledge ask "Consolidate all actor-related entries into one comprehensive guide"
copilot-research knowledge ask "What entries contradict our MV architecture preference?"

# Direct manipulation
copilot-research knowledge delete --where "topic LIKE '%MVVM%'"
copilot-research knowledge merge --ids 1,2,3 --strategy summarize

# Automatic consolidation
copilot-research knowledge consolidate --auto --max-size 1MB
```

### Database Schema
```sql
CREATE TABLE knowledge (
    id INTEGER PRIMARY KEY,
    topic TEXT NOT NULL,
    content TEXT NOT NULL,
    summary TEXT,
    source_urls TEXT, -- JSON array
    added_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_accessed TIMESTAMP,
    access_count INTEGER DEFAULT 0,
    quality_score REAL DEFAULT 0.5,
    embedding BLOB, -- Vector embedding for similarity search
    tags TEXT, -- JSON array
    metadata TEXT -- JSON for extensibility
);

CREATE TABLE preferences (
    key TEXT PRIMARY KEY,
    value TEXT,
    type TEXT, -- 'exclude', 'prefer', 'context'
    reason TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE relationships (
    from_id INTEGER,
    to_id INTEGER,
    relationship_type TEXT, -- 'references', 'contradicts', 'similar'
    strength REAL,
    FOREIGN KEY (from_id) REFERENCES knowledge(id),
    FOREIGN KEY (to_id) REFERENCES knowledge(id)
);

CREATE VIRTUAL TABLE knowledge_fts USING fts5(topic, content, tags);
```

### How It Works
1. **Natural Language Interface**: 
   - User says: "Remove all MVVM entries"
   - CLI calls Copilot with: "You are a knowledge librarian. User wants to remove all MVVM entries. Generate SQL to find them, show results, ask for confirmation, then delete."
   
2. **Automatic Deduplication**:
   - On insert, generate embedding
   - Find similar entries (cosine similarity > 0.85)
   - Prompt user or auto-merge based on config

3. **Smart Consolidation**:
   - Identify clusters of related topics
   - Generate consolidated summaries
   - Keep originals in archive table

4. **Research Integration**:
   - Before research, load relevant existing knowledge
   - After research, merge with database
   - Preference rules applied via SQL filters

### Example Workflow
```bash
# Set preferences
$ copilot-research config set exclude.patterns "MVVM,UIKit patterns"

# Do research (automatically checks preferences)
$ copilot-research "SwiftUI architecture patterns"
🔍 Researching...
✓ Found 15 sources
⚠️  Filtered out 3 entries matching exclude patterns (MVVM)
✓ Added 12 new knowledge entries
💡 Detected 2 similar entries, merged automatically

# Ask about knowledge base
$ copilot-research knowledge ask "What do we know about architecture?"
📚 Found 23 entries about architecture:
   - MV Pattern (9 entries) - Last used: 2 hours ago
   - State Management (8 entries) - Last used: 5 days ago
   - View Composition (6 entries) - Last used: 1 week ago

# Clean up old knowledge
$ copilot-research knowledge ask "Remove entries not accessed in 30 days with quality score < 0.6"
Found 7 entries matching criteria:
  - "Legacy UIViewController patterns" (added 45 days ago, score: 0.4)
  - "Objective-C bridging tips" (added 60 days ago, score: 0.5)
  ...
Delete these? [y/N]
```

### Pros
- Powerful querying and analytics
- Natural language + SQL flexibility
- Built-in deduplication
- Scalable to large knowledge bases
- Can track usage patterns
- Relationship mapping

### Cons
- Most complex to implement
- Requires SQLite + FTS + embeddings
- Learning curve for query language
- Heavier dependency footprint

---

## Recommendation Matrix

| Feature | Option A (Rules) | Option B (TUI) | Option C (Database) |
|---------|-----------------|----------------|---------------------|
| Easy to implement | ⭐⭐⭐⭐ | ⭐⭐ | ⭐ |
| User-friendly | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐ |
| Powerful | ⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| Scriptable | ⭐⭐⭐⭐⭐ | ⭐⭐ | ⭐⭐⭐⭐⭐ |
| Scalable | ⭐⭐⭐ | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ |
| AI-assisted | ⭐⭐⭐ | ⭐⭐⭐⭐⭐ | ⭐⭐⭐⭐ |

## Hybrid Approach (Recommended)

**Start with Option A, evolve to Option C**

### Phase 1: Rule-Based (MVP)
- Implement preferences.yaml
- Basic consolidation via AI prompts
- Simple CLI commands
- Good enough for initial use

### Phase 2: Add Database (Scale)
- Migrate to SQLite while keeping preferences.yaml
- Add similarity detection
- Add usage tracking
- Keep simple API

### Phase 3: Optional TUI (Polish)
- Build interactive editor for power users
- Keep CLI for scripts
- Best of both worlds

This gives us:
- ✅ Quick to ship
- ✅ Easy to use initially
- ✅ Scales as we need it
- ✅ Can add sophistication incrementally
