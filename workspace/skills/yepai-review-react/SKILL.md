---
name: yepai-review-react
description: >-
  React/TypeScript code reviewer. Checks hooks correctness, dependency arrays,
  component boundaries, state management, rendering performance, accessibility,
  and TypeScript typing.
user-invocable: false
disable-model-invocation: false
metadata: {"openclaw":{"requires":{"bins":["gh"],"env":["GITHUB_TOKEN"]}}}
---

# React/TypeScript Review Criteria

When reviewing `.tsx`, `.jsx`, and `.ts` files in a PR diff, check for these issues. Report each finding with severity, file, line, issue, and recommendation.

## Focus Areas

- **Hooks correctness**: Hooks called inside conditions, loops, or nested functions (Rules of Hooks violation), custom hooks not prefixed with `use`
- **Dependency arrays**: Missing dependencies in `useEffect`/`useMemo`/`useCallback` arrays, stale closures from incorrect deps, unnecessary dependencies causing excessive re-runs
- **Component boundaries**: Prop drilling through 3+ component levels (suggest context or composition), components exceeding ~300 lines, mixing concerns (data fetching + rendering in one component)
- **State management**: Derived state stored in `useState` (should be computed), state updates during render (infinite loop risk), unnecessary state that could be a ref or computed value
- **Rendering performance**: Missing `key` prop or using array index as `key` in dynamic lists, large inline objects/functions in JSX causing re-renders, missing memoization on expensive computations
- **Accessibility**: Missing `alt` attribute on `<img>`, non-semantic elements used as buttons without ARIA roles, missing keyboard event handlers alongside click handlers
- **TypeScript typing**: `any` type usage, missing return type annotations on exported functions, unsafe type assertions (`as unknown as T` chains)

## Severity Guide

| Anti-Pattern | Severity |
|---|---|
| Hooks called inside conditions or loops | BLOCKER |
| State update in render body (infinite loop risk) | BLOCKER |
| `useEffect` with missing dependency in array | MAJOR |
| Using array index as `key` in dynamic lists | MINOR |
| Using `any` type | MINOR |
| Missing `alt` attribute on `<img>` | MINOR |
| Prop drilling through 3+ levels | MINOR |
