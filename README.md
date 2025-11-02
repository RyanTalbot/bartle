# Bartle

**Bartle** helps teams write clean, consistent commit messages across repositories.

A lightweight Git companion built for developers and teams who care about clean history.

## Features

-  **Single source of truth** 🧩 `.bartle.yaml`  defines commit rules for your team
-  **Git-integrated** ⚙️ runs automatically via `commit-msg` hook
-  **Conventional or JIRA styles** 🔍 out of the box
-  **Fast failure** 🚫 reject invalid commits before they hit your repo

---

## Installation (Temporary)

```bash
# build from source
git clone https://github.com/RyanTalbot/bartle.git
cd bartle
go build -o bartle .
```

---

## Quick Start

Get up and running with Bartle in under a minute.

### 1. Initialize Bartle in your repository

```bash
bartle init
```

This will create a `.bartle.yaml` file in your project root.

```yaml
# .bartle.yaml
style: conventional
ai:
  enabled: false
  provider: openai
  model: gpt-5
  api_key: env:OPENAI_API_KEY
  temperature: 0.2
rules:
  scope_required: true
  max_line_length: 72
  lowercase_start: false
  types: [feat, fix, docs, refactor, test, chore]
hook:
  auto_apply: false
  block_on_fail: true
```

### 2. Install the Git hook

```bash
bartle install-hook
```

This adds a commit-msg hook to your repository, so every commit is validated automatically.


### 3. Commit your changes

Bad commits will be rejected, and valid commits will be accepted.

```bash
git commit -m "this is a bad commit"
# ❌ Invalid commit message:
# - missing ':' separator (e.g., type(scope): subject)

git commit -m "feat(ui): add button"
# ✅ Commit message is valid!
```

### 4. (Optional) Lint your commits manually

You can run the linter manually to check your commit messages first.

```bash
bartle lint -m "feat(api): add pagination"
# ✅ Commit message is valid!

bartle lint -m "wrong format"
# ❌ Invalid commit message:
#  - not conventional format (e.g., type(scope): subject)
```

### 5. (Optional) Uninstall the git hook

You can uninstall the hook at any time.

```bash
bartle uninstall-hook
```

