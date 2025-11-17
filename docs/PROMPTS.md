# Copilot Research - Prompt System Documentation

This document explains how the prompt system in `copilot-research` works, how to create and manage custom prompts, and best practices for prompt engineering.

## Table of Contents
- [Overview](#overview)
- [Prompt File Structure](#prompt-file-structure)
- [Template Variables](#template-variables)
- [Research Modes](#research-modes)
- [Managing Prompts](#managing-prompts)
- [Best Practices for Prompt Engineering](#best-practices-for-prompt-engineering)

---

## Overview

`copilot-research` uses a flexible prompt system that allows you to define and customize the instructions given to the AI provider. Prompts are stored as Markdown files with YAML frontmatter, making them easy to read, edit, and version control.

The tool comes with several built-in prompts optimized for different research modes (e.g., `quick`, `deep`, `compare`, `synthesis`). You can also create your own custom prompts to tailor the AI's behavior to your specific needs.

## Prompt File Structure

Prompts are Markdown files located in the `prompts/` directory within your `copilot-research` configuration directory (default: `~/.copilot-research/prompts/`). Each prompt file consists of two main parts:

1.  **YAML Frontmatter**: A block of YAML at the beginning of the file, enclosed by `---` delimiters. This contains metadata about the prompt.
2.  **Prompt Template Content**: The rest of the Markdown content, which serves as the actual instructions for the AI.

### Example Prompt File (`prompts/default.md`)

```markdown
---
name: default
description: Default research prompt for comprehensive queries
version: 1.0.0
mode: quick
---

You are an expert research assistant specializing in {{mode}} research.
Your task is to provide a comprehensive and accurate report based on the user's query.

**Instructions:**
- Provide a clear and concise overview of the topic.
- Include detailed explanations, examples, and relevant context.
- Cite all sources used, preferably with links.
- Format your response using Markdown, with clear headings and bullet points.
- Ensure the information is up-to-date and factual.

**Research Query:** {{query}}
```

### Frontmatter Fields

The following fields are supported in the YAML frontmatter:

-   `name` (string, **required**): A unique identifier for the prompt (e.g., `default`, `deep-dive`). This is used to select the prompt via the CLI.
-   `description` (string): A brief description of what the prompt is designed for.
-   `version` (string): The version of the prompt.
-   `mode` (string, optional): The default research mode associated with this prompt. If a mode is specified in the prompt's frontmatter, it will override the global `--mode` flag unless explicitly overridden by the user.

## Template Variables

Prompts can include template variables, which are placeholders that `copilot-research` replaces with dynamic content before sending the prompt to the AI provider. Variables are enclosed in double curly braces `{{variable_name}}`.

Currently supported template variables:

-   `{{query}}`: Replaced with the user's research query.
-   `{{mode}}`: Replaced with the active research mode (e.g., `quick`, `deep`).

### Example Usage of Template Variables

```markdown
Your primary goal is to answer the following {{mode}} research question: "{{query}}".
```

## Research Modes

Research modes are predefined strategies that influence how the AI approaches a query. They are often tied to specific prompt templates but can also be overridden.

-   `quick`: Designed for brief overviews and summaries.
-   `deep`: Focuses on in-depth analysis, detailed explanations, and examples.
-   `compare`: Structured to compare and contrast multiple subjects.
-   `synthesis`: Aims to integrate information from various sources into a cohesive report.

You can select a research mode using the `--mode` flag (e.g., `copilot-research "topic" --mode deep`). If a prompt has a `mode` defined in its frontmatter, that mode will be used unless the `--mode` flag is explicitly provided by the user.

## Managing Prompts

### Listing Available Prompts
To see all available prompts (both built-in and custom), use the `prompts list` command:

```bash
copilot-research prompts list
```

### Using a Specific Prompt
You can specify which prompt template to use for a research query with the `--prompt` or `-p` flag:

```bash
copilot-research "Explain microservices" --prompt deep-dive
```

### Setting a Default Prompt
You can set a custom prompt as your default in the configuration:

```bash
copilot-research config set active_prompt my-custom-prompt
```

### Creating Custom Prompts
To create a new custom prompt:
1.  Create a new Markdown file (e.g., `my-custom-prompt.md`) in your `~/.copilot-research/prompts/` directory.
2.  Add the YAML frontmatter and your prompt template content.
3.  You can then use it with `copilot-research "query" --prompt my-custom-prompt`.

## Best Practices for Prompt Engineering

Effective prompt engineering is key to getting high-quality results from AI.

-   **Be Specific**: Clearly define the AI's role, the task, and the desired output format.
-   **Provide Context**: Give the AI enough background information to understand the query.
-   **Specify Format**: Request output in a structured format (e.g., Markdown with headings, bullet points, code blocks).
-   **Emphasize Accuracy**: Instruct the AI to prioritize factual correctness and cite sources.
-   **Use Examples**: If possible, provide examples of good responses.
-   **Iterate**: Experiment with different prompts and modes to find what works best for your research needs.
-   **Leverage Template Variables**: Use `{{query}}` and `{{mode}}` to make your prompts dynamic and reusable.
