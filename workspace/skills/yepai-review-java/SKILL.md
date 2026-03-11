---
name: yepai-review-java
description: >-
  Java/Spring Boot code reviewer. Checks transaction boundaries, null safety,
  Optional misuse, concurrency, exception handling, and API design.
user-invocable: false
disable-model-invocation: false
metadata: {"openclaw":{"requires":{"bins":["gh"],"env":["GITHUB_TOKEN"]}}}
---

# Java Review Criteria

When reviewing Java files in a PR diff, check for these issues. Report each finding with severity, file, line, issue, and recommendation.

## Focus Areas

- **Spring Boot best practices**: Missing `@Service`/`@Repository` stereotypes, field injection instead of constructor injection, missing `@Configuration` on config classes
- **Transaction boundaries**: `@Transactional` on private methods (no effect), missing `@Transactional` on service methods that modify data, incorrect rollback rules
- **Null safety / Optional misuse**: `Optional.get()` without `isPresent()` or `orElse`, returning `null` from methods declaring `Optional`, missing `@Nullable`/`@NonNull` annotations on public API
- **Concurrency**: `synchronized` on non-final fields, shared mutable state without synchronization, race conditions in lazy initialization
- **Exception handling**: Catching generic `Exception`/`Throwable`, swallowed exceptions (empty catch blocks), missing error propagation in async flows
- **API design**: Missing `@Valid` on `@RequestBody` parameters, hardcoded config values instead of `@Value`/`@ConfigurationProperties`, non-RESTful endpoint naming

## Severity Guide

| Anti-Pattern | Severity |
|---|---|
| `synchronized` on non-final fields | BLOCKER |
| `Optional.get()` without `isPresent()`/`orElse` | MAJOR |
| `@Transactional` on private methods (no effect) | MAJOR |
| Missing `@Valid` on `@RequestBody` parameters | MINOR |
| Catching `Exception` instead of specific types | MINOR |
| Hardcoded config values instead of `@Value` | MINOR |
