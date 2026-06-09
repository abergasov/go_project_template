# AGENTS.md — Optimum Gateway Developer Guide

## Operating rules

### 1. Stay inside scope

Before editing, identify the requested behavior, the likely touched packages, and the non-goals.

Change only what is required to satisfy the task.
Do not add side refactors.
Do not rename unrelated things.
Do not move files because it feels cleaner.
Do not “improve” architecture unless the task explicitly includes it.

If a correct fix needs broader changes, say that before editing.

### 2. Respect existing conventions

Look at the nearest production code and tests in the same package first.

Match the repository’s established patterns for:
- naming
- package layout
- error handling
- logging via `logger.With*`
- context plumbing
- tests
- configuration access through `AppConfig` getters/runtime helpers
- dependency injection
- concurrency control

Do not introduce a new style when an existing one already solves the problem.

### 3. Prefer minimal change

Choose the smallest change that solves the problem.

Prefer:
- modifying existing flow
- reusing existing helpers
- extending existing types carefully
- deleting complexity when possible

Avoid:
- speculative abstractions
- new layers
- generic helpers added “for future use”
- broad rewrites

### 4. Make code readable

Write code that another engineer can follow quickly.

Prefer:
- clear names
- short control flow
- explicit error handling
- narrow responsibilities
- small helpers only when they remove duplication or clarify a single concept
- comments only where intent is not obvious

Avoid:
- clever compactness
- hidden side effects
- unnecessary indirection
- dense nested logic

### 5. Preserve invariants

Before changing behavior, identify invariants in the touched path.

Do not silently break:
- protocol assumptions
- lifecycle guarantees
- cleanup behavior
- locking discipline
- context cancellation semantics
- persistence assumptions
- existing API behavior

If an invariant must change, state the old behavior, the new behavior, and why the change is required.

### 6. Distinguish required change from optional improvement

You may notice other problems nearby.
That does not make them part of this task.

If they matter, report them in risks or follow-up notes.
Do not expand scope on your own.

### 7. Be honest about uncertainty

If code behavior is unclear or conflicts with research:

- stop
- report the conflict

Do not patch blindly.

### 8. Avoid unnecessary surface growth

Do not add new Go module dependencies, new top-level packages, or new public APIs unless the task requires them.
Prefer existing packages, config fields, and interfaces.

## Code quality rules

Your implementation should be:
* simple
* explicit
* locally understandable
* behavior-safe
* easy to review
* easy to test

Prefer:
* existing interfaces, config getters, and logger helpers over new ones
* explicit branching over magic helpers
* small helper extraction only when it reduces local complexity
* stable behavior over elegant rewrite
* focused tests close to the changed behavior

Avoid:
* generic wrappers with single use
* helper explosion
* comments that restate code
* hidden global state
* unnecessary concurrency changes
* new dependencies without task pressure
* premature optimization

---

## Failure conditions

Your work is considered poor if you:
* change unrelated files
* introduce new abstraction without pressure
* fail to explain changed behavior
* violate non-goals
* silently alter existing semantics
* produce vague report instead of exact changes
* hide uncertainty
* make code more complex than task requires

---

## Response style

* Be exact.
* Be concise.
* Be technical.
* Do not pad.
* Do not sell the code.
* Do not claim elegance.
* Do not hide weak points.

Use plain engineering language.

---

## Final instruction

Your implementation must make the codebase better in the narrowest way needed to solve the task.

Write less.
Change less.
Break less.
Explain exactly what changed.
