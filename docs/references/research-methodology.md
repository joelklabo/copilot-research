# Research Methodology for AI Agents

## Core Principles

1. **Breadth then Depth**: Start wide, narrow focus based on findings
2. **Multiple Sources**: Corroborate across 3+ independent sources
3. **Recency Matters**: Prioritize recent information for fast-moving topics
4. **Authority Hierarchy**: Official docs > Academic papers > Expert blogs > General articles
5. **Synthesis over Aggregation**: Create coherent narrative, not just bullet points

## Multi-Stage Research Process

### Stage 1: Landscape Survey (5-10 queries)
**Goal**: Understand the topic boundaries and key subtopics

```
Query patterns:
- "What is [topic]?"
- "[Topic] overview 2025"
- "Latest developments in [topic]"
- "[Topic] best practices"
- "Common [topic] patterns"
```

**Output**: Topic map with key areas to explore

### Stage 2: Deep Dive (10-20 queries)
**Goal**: Detailed understanding of each key area

```
Query patterns:
- "[Subtopic] implementation guide"
- "[Subtopic] examples and use cases"
- "[Subtopic] common mistakes"
- "[Subtopic] vs [Alternative]"
- "How to [specific task] with [topic]"
```

**Output**: Detailed notes per subtopic with examples

### Stage 3: Synthesis (3-5 queries)
**Goal**: Fill gaps and resolve contradictions

```
Query patterns:
- "How does [concept A] relate to [concept B]?"
- "[Edge case] in [topic]"
- "When to use [approach A] vs [approach B]"
```

**Output**: Cohesive understanding with nuance

### Stage 4: Validation (2-3 queries)
**Goal**: Verify conclusions with authoritative sources

```
Query patterns:
- "Official [topic] documentation"
- "[Topic] release notes 2025"
- "[Topic] migration guide"
```

**Output**: Fact-checked, cited report

## Query Optimization

### Good Queries (Specific, Actionable)
✅ "Swift 6 strict concurrency checking migration guide"
✅ "Bubble Tea async operations best practices"
✅ "SQLite WAL mode vs DELETE mode performance"

### Bad Queries (Vague, Too Broad)
❌ "Tell me about Swift"
❌ "How do I make UIs?"
❌ "Databases"

### Query Templates

**For New Technologies:**
```
1. "[Tech] official documentation 2025"
2. "[Tech] getting started tutorial"
3. "[Tech] vs [Alternative] comparison"
4. "[Tech] real-world examples"
5. "[Tech] common pitfalls"
```

**For Design Patterns:**
```
1. "[Pattern] definition and purpose"
2. "[Pattern] implementation in [language]"
3. "[Pattern] use cases and examples"
4. "[Pattern] advantages and disadvantages"
5. "When to use [Pattern] vs [Alternative]"
```

**For Troubleshooting:**
```
1. "[Error message]" exact match
2. "[Technology] [symptom] solutions"
3. "[Technology] debugging guide"
4. "[Technology] known issues"
```

## Source Evaluation

### Authority Levels

**Tier 1 (Highest Authority)**:
- Official documentation
- Language/framework specifications
- Release notes and changelogs
- Official blogs and announcements

**Tier 2 (High Authority)**:
- Academic papers (peer-reviewed)
- Established expert blogs (authors with credentials)
- Conference talks from maintainers
- Well-maintained open source examples

**Tier 3 (Moderate Authority)**:
- Technical blogs from practitioners
- Stack Overflow answers (high votes)
- GitHub discussions
- Medium articles with citations

**Tier 4 (Lower Authority)**:
- General tutorials without attribution
- Forum posts
- Reddit discussions
- Uncited articles

### Red Flags
🚩 No publication date
🚩 No author credentials
🚩 Contradicts official docs without explanation
🚩 No code examples or evidence
🚩 Overly promotional language

## Deduplication Strategy

### Identifying Duplicates
1. **Exact matches**: Same facts, same examples
2. **Derivative content**: One clearly copying another
3. **Outdated versions**: Superseded by newer information

### Consolidation Rules
- Keep most recent version
- Prefer official over unofficial
- Preserve unique examples from each
- Combine complementary information
- Note when sources disagree (with citations)

## Synthesis Techniques

### Pattern: Hierarchical Summary
```
# Main Topic

## Overview (2-3 sentences)

## Key Concepts
### Concept 1
- Definition
- Example
- Best practices

### Concept 2
...

## Common Patterns
1. Pattern A
2. Pattern B

## Pitfalls and Solutions

## Resources
- [Official Docs](url)
- [Tutorial](url)
```

### Pattern: Comparative Analysis
```
| Aspect | Approach A | Approach B |
|--------|-----------|-----------|
| Performance | Fast | Moderate |
| Complexity | High | Low |
| Use Case | X, Y | Z |
```

### Pattern: Decision Tree
```
Start with: Do you need X?
  → Yes: Use approach A
    → Requires Y? Use variant A1
    → Requires Z? Use variant A2
  → No: Use approach B
```

## Citation Standards

### Inline Citations
Use superscript numbers: "Swift 6 introduces strict concurrency checking¹"

### Source List Format
```
## Sources
1. [Swift 6 Documentation](url) - Official Apple Docs
2. [WWDC 2024: What's New in Swift](url) - Apple Developer Video
3. [Swift Concurrency Migration Guide](url) - Swift.org
```

### When to Cite
- **Always**: Specific claims, statistics, code examples
- **Optional**: General knowledge, widely-known facts
- **Multiple sources**: Controversial or critical information

## Knowledge Persistence

### Database Schema
```sql
CREATE TABLE research_sessions (
    id INTEGER PRIMARY KEY,
    query TEXT NOT NULL,
    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP,
    result TEXT,
    sources TEXT, -- JSON array
    confidence_score REAL
);

CREATE TABLE facts (
    id INTEGER PRIMARY KEY,
    content TEXT NOT NULL,
    topic TEXT,
    source_url TEXT,
    authority_level INTEGER,
    date_verified DATETIME,
    UNIQUE(content, source_url)
);
```

### Reuse Strategy
1. Check cache for similar queries (fuzzy match)
2. Return cached results if < 7 days old (adjust by topic velocity)
3. Augment old results with new findings
4. Update confidence scores based on corroboration

## Quality Metrics

### Completeness
- [ ] Covers all major aspects of topic
- [ ] Includes examples
- [ ] Addresses common questions
- [ ] Notes limitations/trade-offs

### Accuracy
- [ ] 3+ sources for key claims
- [ ] No contradictions unresolved
- [ ] Recent information (< 1 year for tech)
- [ ] Authoritative sources

### Clarity
- [ ] Logical structure
- [ ] Technical terms defined
- [ ] Progressive complexity
- [ ] Actionable takeaways

### Usefulness
- [ ] Answers the question asked
- [ ] Provides next steps
- [ ] Includes working examples
- [ ] Cites resources for deeper learning

## Iterative Refinement

### Self-Review Checklist
1. Did I answer the core question?
2. Are there obvious gaps?
3. Do sources contradict each other?
4. Is this still accurate (check dates)?
5. Could someone use this to accomplish a task?

### Follow-up Query Triggers
- "More details on [specific aspect]"
- "Examples of [specific use case]"
- "How to [specific task]"
- "Comparison of [alternatives]"
- "Latest developments in [topic]"

## Advanced Techniques

### Cross-Domain Synthesis
When topic spans multiple domains (e.g., "iOS app performance optimization"):
1. Research each domain independently (iOS, performance, optimization)
2. Find intersection points
3. Synthesize domain-specific approaches
4. Identify unique considerations at intersection

### Historical Context
For established technologies, include evolution:
1. Original design and rationale
2. Major version changes
3. Paradigm shifts
4. Current best practices
5. Future direction (if known)

### Practical Application
Always include "How to use this" section:
- Setup steps
- Basic example
- Common variations
- Testing approach
- Production considerations

## Tools and Resources

### Research Tools
- GitHub code search (real implementations)
- Official documentation sites
- Stack Overflow (solutions)
- Academic paper databases (theory)
- Conference talk archives (cutting-edge)

### Quality Indicators
- Recent commit activity (if open source)
- Issue count and resolution rate
- Community size and engagement
- Documentation completeness
- Example projects

---

*This methodology evolves based on what works in practice. Update as you discover better approaches.*
