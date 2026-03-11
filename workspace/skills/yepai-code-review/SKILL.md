---
name: yepai-code-review
description: >-
  MUST be used when user sends a GitHub Pull Request URL (github.com/*/pull/*).
  Uses Claude Code CLI to perform expert-level code review and post comments.
  Do NOT use web_fetch for GitHub PR URLs — always use this skill instead.
user-invocable: true
metadata: {"openclaw":{"requires":{"bins":["gh","claude"],"env":["GITHUB_TOKEN","ANTHROPIC_API_KEY"]},"emoji":"🔍"}}
---

# YepAI Code Review

When the user sends a message containing a GitHub PR URL, run the review script.

## Step 1: Extract PR URL

- Find a URL matching: `https://github.com/{owner}/{repo}/pull/{number}`
- Accept Slack-formatted links like `<URL|text>` or `<URL>` — strip the angle brackets and display text.
- **Sanitize**: reject URLs containing shell metacharacters (`;`, `|`, `&`, `$`, `` ` ``, `(`, `)`, `\n`). Only pass URLs matching the strict pattern above.
- If no valid PR URL is found, respond:
  "Usage: @yepai-code-review https://github.com/org/repo/pull/123"

## Step 2: Run Review Script

Run this command using the `exec` tool:

```
bash /home/picoclaw/.picoclaw/workspace/skills/yepai-code-review/review.sh <PR_URL>
```

Replace `<PR_URL>` with the extracted URL from Step 1.

This script uses Claude Code CLI to:
- Fetch the PR diff and metadata via `gh`
- Classify files by language (Java, Python, React/TypeScript)
- Apply language-specific review criteria
- Post the review as a GitHub PR comment

## Step 3: Relay Output

Forward the script output as your response. Do not add extra commentary or modify the output.

If the script fails, respond: "Code review failed. Please check that GITHUB_TOKEN and ANTHROPIC_API_KEY are configured and the PR URL is accessible."
