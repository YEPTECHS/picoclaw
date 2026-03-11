---
name: yepai-review-python
description: >-
  Python code reviewer. Checks async/await correctness, error handling,
  input validation, security issues, logging, and code quality.
user-invocable: false
disable-model-invocation: false
metadata: {"openclaw":{"requires":{"bins":["gh"],"env":["GITHUB_TOKEN"]}}}
---

# Python Review Criteria

When reviewing Python files in a PR diff, check for these issues. Report each finding with severity, file, line, issue, and recommendation.

## Focus Areas

- **Async/await correctness**: Missing `await` on coroutine calls, sync blocking calls inside `async` functions (e.g., `time.sleep` instead of `asyncio.sleep`), unawaited coroutine warnings
- **Error handling**: Bare `except:` clauses, `except Exception as e: pass` (swallowed exceptions), overly broad exception catching hiding bugs, missing error context in re-raises
- **Input validation**: Unsanitized user input used in operations, missing type checks at boundaries, unchecked dictionary key access on external data
- **Security**: f-strings or string formatting with user input in SQL queries or shell commands (injection risk), `subprocess` calls with unsanitized input (`shell=True`), path traversal via unvalidated file paths
- **Logging**: Missing log statements in error handlers, sensitive data (passwords, tokens) in log output, inconsistent log levels
- **Code quality**: `import *` usage, mutable default arguments (`def f(x=[])`), unused imports, overly complex functions (high cyclomatic complexity)

## Severity Guide

| Anti-Pattern | Severity |
|---|---|
| Missing `await` on coroutine call | BLOCKER |
| f-string with user input in SQL/shell command | BLOCKER |
| `subprocess` with unsanitized input (`shell=True`) | BLOCKER |
| `except Exception as e: pass` (swallowed) | MAJOR |
| Bare `except:` clause | MINOR |
| `import *` usage | MINOR |
| Mutable default arguments (`def f(x=[])`) | MINOR |
