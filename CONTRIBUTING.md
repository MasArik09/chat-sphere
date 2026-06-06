# Contributing to ChatSphere

Thank you for your interest in contributing to ChatSphere! This document outlines guidelines, branch naming structures, commit message styles, and the pull request workflow to help you get started.

---

## 1. Code of Conduct
By participating in this project, you agree to maintain a respectful, welcoming, and collaborative environment.

---

## 2. Local Development Setup

To begin contributing:
1. **Fork the Repository**: Create a personal fork on GitHub.
2. **Clone the Fork**:
   ```bash
   git clone https://github.com/your-username/chat-sphere.git
   cd chat-sphere
   ```
3. **Set Up Development Environments**:
   - Install **Docker Desktop** (recommended to run databases and dependencies easily via `docker compose up -d`).
   - Install **Go 1.22+** if you plan to write backend API modifications locally.
   - Install **Node.js 20+** if you plan to edit React frontend assets locally.
4. **Run Verification Commands**:
   - Backend tests: Run `go test -v ./internal/...` in the `backend` directory.
   - Frontend verification: Run `npm run lint` and `npm run build` in the `frontend` directory.

---

## 3. Branch Naming Convention

We use a structured branch naming convention. When creating a new branch, prefix it with the type of work being performed:

- `feature/` : New features or UI additions (e.g., `feature/group-chat-creation`).
- `bugfix/`  : Resolving bugs or N+1 query loops (e.g., `bugfix/typing-state-cleanup`).
- `docs/`    : Documentation updates and ADR records (e.g., `docs/add-api-reference`).
- `refactor/`: Refactoring code without functional changes (e.g., `refactor/auth-middleware`).
- `test/`    : Adding or rewriting unit or integration tests (e.g., `test/websocket-manager-mocks`).
- `ci/`      : GitHub actions pipelines modifications (e.g., `ci/add-cache-dependencies`).

---

## 4. Git Commit Message Conventions

We enforce a **Conventional Commits** format. This ensures clean repository history and automates changelog generation.

### Commit Format
```text
<type>(<scope>): <short description>

[optional body describing technical decisions or context]

[optional footer referencing issue numbers e.g. Closes #123]
```

### Supported Types:
- `feat`     : A new user-facing feature.
- `fix`      : A bug fix (e.g. database rollback corrections).
- `docs`     : Changes to markdown files, ADRs, or README guides.
- `style`    : Layout adjustments (indentations, CSS classes, no code logic changes).
- `refactor` : Structural changes that neither fix a bug nor add a feature.
- `test`     : Creating unit tests or integration test runners.
- `chore`    : Configuration updates, dependency upgrades, or repository hygiene.

### Commit Examples:
* *Feature*: `feat(websocket): broadcast typing indicators to conversation participants`
* *Bug Fix*: `fix(database): rollback message creation when timestamp update fails`
* *Docs*: `docs(readme): add docker orchestration commands to quickstart guide`

---

## 5. Pull Request Workflow

1. **Keep Pull Requests Focused**: Each pull request should address a single ticket or feature. Avoid combining unrelated fixes into a single PR.
2. **Push to Your Fork**:
   ```bash
   git push origin feature/your-feature-name
   ```
3. **Open a Pull Request**: Submit a PR to the `develop` (or `main`) branch of the upstream repository.
4. **Code Quality Checklist**:
   - Write unit tests for new service logic.
   - Run linter checks to ensure no formatting errors exist.
   - Fill out the PR description template detailing *what* changes were made, *why* they were made, and *how* they were tested.
5. **Code Review**: A maintainer will review your code. Address any requested changes by pushing commits to the same branch.
6. **Merge**: Once approved and all CI checks pass, a maintainer will squash and merge your PR.

---

## 6. Code Style Expectations

### Go (Backend) Style
- **Formatting**: All files must be formatted using the official `gofmt` standard.
- **Error Handling**: Never ignore returned errors. Handle them explicitly (e.g., logging them and returning early or wrapping them in context).
- **Naming Conventions**: Use `camelCase` for private variables and `PascalCase` for public structs, interfaces, and methods.

### TypeScript / React (Frontend) Style
- **Type Safety**: Avoid using `any`. Declare explicit interfaces or types for component props, API request/response structures, and store models.
- **State Management**: Use Zustand stores for global synchronization states. Avoid over-complicating components with local state if global states apply.
- **Component File Structure**: Group related components inside feature directories (e.g. `src/features/auth`, `src/features/conversations`).
