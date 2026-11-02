---
name: bk-apigateway-framework-patch
description: "Patch the bk-apigateway-framework cookiecutter templates by first rendering a real project from templates/golang or templates/python, making and verifying the desired changes there, then feeding the intentional diff back into the matching template while preserving cookiecutter variables and Jinja raw blocks. Use this whenever the user asks to update, upgrade, test, or validate generated projects from bk-apigateway-framework templates, mentions demo_go/demo_go2/template-generated projects, asks to sync generated-project changes back to templates/golang or templates/python, or wants to prove a fresh project created from the newest template still builds/tests."
compatibility: Requires local shell access, cookiecutter, git, and the language runtime needed by the selected template. For Go template work, source the generated project's or reference project's .envrc when provided and use GOTOOLCHAIN=local when Go version stability matters.
---

# bk-apigateway-framework-patch

Use this skill for the recurring workflow in `bk-apigateway-framework`: make a
change in a generated project, then backport that change into the cookiecutter
template and prove a newly rendered project still works.

This is not a generic dependency-update skill. Its point is to avoid the common
mistake of fixing only `demo_go` or another generated sandbox while leaving
`templates/golang/{{cookiecutter.project_name}}` or `templates/python/...` stale.

## Core Principle

Treat generated projects as disposable proof sandboxes. Treat `templates/golang`
and `templates/python` as the source of truth.

Every template patch must preserve cookiecutter syntax:
- keep placeholders such as `{{cookiecutter.project_name}}`
- keep Jinja control blocks such as `{% raw %}` and `{% endraw %}`
- keep template-specific defaults from `cookiecutter.json`
- do not replace template placeholders with concrete generated-project names

## Trigger Examples

Use this skill for prompts like:
- "create demo_go2 from templates/golang and check the newest template"
- "upgrade this generated project dependency, then sync it back to templates/golang"
- "we changed demo_go; patch the cookiecutter template too"
- "render from templates/python, modify generated files, then update the template"
- "make sure a fresh project from the latest template builds"

Do not use it for ordinary application changes that are not intended to flow
back into `templates/golang` or `templates/python`.

## Required Workflow

### 1. Identify the template and working project

Determine:
- target template: `templates/golang` or `templates/python`
- generated project to modify or inspect, such as `demo_go`, `demo_go2`, or a
  temporary render
- reference environment, such as `demo_go/.envrc` for Go
- requested change, such as dependency upgrade, content edit, or new file

Read first:
- current `git status --short --untracked-files=all`
- the relevant `templates/<lang>/cookiecutter.json`
- generated project files named by the user
- matching template files under `templates/<lang>/{{cookiecutter.project_name}}`

Distinguish validation-only from template mutation. A request to prove a fresh render builds does not itself authorize template fixes. Render missing named outputs only when their creation is requested; otherwise use a temporary sandbox.

### 2. Render a real project before changing the template

Use cookiecutter to produce a concrete project in a new directory. Do not overwrite an existing generated project or treat its user edits as disposable. Record the render inputs so the proof can be repeated.

Example for Go:

```bash
cookiecutter --no-input templates/golang -o . \
  project_name=demo_go2 \
  init_admin=admin
```

Example for a temporary validation render:

```bash
tmp=$(mktemp -d)
cookiecutter --no-input templates/golang -o "$tmp" \
  project_name=template_probe \
  init_admin=admin
```

Prefer a temporary render when the user only wants validation. Use a named
project such as `demo_go2` only when the user explicitly asks for that artifact.

### 3. Activate the right runtime

For Go template work:
- prefer the generated project's `.envrc` if it exists
- otherwise use the known reference env the user named, for example
  `source demo_go/.envrc`
- when preserving Go version compatibility, run Go commands with
  `GOTOOLCHAIN=local` so Go cannot silently download a newer toolchain
- verify with `go version` and `go env GOVERSION GOTOOLCHAIN GOMOD`

For Python template work:
- read the generated project's environment docs and Makefile first
- activate the documented virtualenv or `.envrc` if present
- verify with `python --version` and the relevant package manager command

Stop and ask only if the runtime cannot be discovered from project files or the
user's instruction.

### 4. Make the change in the generated project first

Apply the requested change in the rendered project and verify it there.

Examples:
- Go dependency changes: run `go get`, then `go mod tidy`
- Python dependency changes: use the project-supported dependency manager and
  lock/update command
- file content changes: edit the generated file, then run the smallest relevant
  build/lint/test target
- new files: create them in the generated project first so names, imports, and
  generated values are concrete

Keep the generated change narrow. Avoid unrelated tidy, format, or cleanup churn
unless the project toolchain requires it.

### 5. Translate the generated diff back into the template

Patch the corresponding template file under:
- `templates/golang/{{cookiecutter.project_name}}/...`
- `templates/python/{{cookiecutter.project_name}}/...`

When copying from a generated project back to the template:
- replace concrete project names with `{{cookiecutter.project_name}}`
- restore Jinja blocks that cookiecutter consumed during render
- keep raw template files raw; do not compile rendered output back over a raw
  Jinja template file
- preserve template-only dependencies or files unless a fresh template render
  and project verification proves they are no longer used
- if using `go mod tidy` on a rendered copy, copy `go.mod`/`go.sum` back only
  after restoring the placeholder module path in `go.mod`

For Go modules, a safe pattern is:

```bash
tmp=$(mktemp -d)
cookiecutter --no-input templates/golang -o "$tmp" \
  project_name=template_probe \
  init_admin=admin
cd "$tmp/template_probe"
source /path/to/reference/.envrc
GOTOOLCHAIN=local go get example.com/mod@vX.Y.Z
GOTOOLCHAIN=local go mod tidy
cp go.sum /repo/templates/golang/'{{cookiecutter.project_name}}'/go.sum
cp go.mod /repo/templates/golang/'{{cookiecutter.project_name}}'/go.mod
perl -0pi -e 's/module github\.com\/TencentBlueKing\/template_probe/module github.com\/TencentBlueKing\/\{\{cookiecutter.project_name\}\}/' \
  /repo/templates/golang/'{{cookiecutter.project_name}}'/go.mod
```

Then inspect the template diff manually. Do not trust the copy step by itself.

### 6. Render the patched template again

After patching the template, render a fresh project from the newest template.
Use a fresh output directory so old files cannot hide a template problem. A different project name is useful when testing placeholder restoration, but is not required for isolation.

```bash
tmp=$(mktemp -d)
cookiecutter --no-input templates/golang -o "$tmp" \
  project_name=template_probe \
  init_admin=admin
```

For template files that contain raw Jinja sections, validate only a fully
rendered cookiecutter project. Partial string replacement is not enough and can
leave invalid `{% raw %}` syntax in Go/Python files.

### 7. Verify the fresh render

Run the project's permitted verification commands from the fresh render, and confirm runtime/module paths point to that render. If setup or `go mod tidy` changes dependencies or generated files, account for those changes in the template and render again; do not silently repair the validation copy and call the original template verified.

For Go projects, these are candidate checks only when the generated project supports and permits them:

```bash
source /path/to/reference/.envrc
GOTOOLCHAIN=local go mod tidy
GOTOOLCHAIN=local make build
GOTOOLCHAIN=local make vet
GOTOOLCHAIN=local make test
GOTOOLCHAIN=local go list -m all | rg 'modules-you-changed'
```

An example target is not automatically required. Follow the real Makefile/project contract and command prohibitions; report any missing required check rather than inventing replacements.

For Python projects, use the generated project's Makefile or documented commands
for install, lint, test, and build. If there is no documented command, say so
and run only the checks that the project exposes.

### 8. Report the result

Final reports should include:
- which generated project was used as the sandbox
- which template files changed
- which cookiecutter render was used for final validation
- exact verification commands and pass/fail results
- any version ceiling discovered, especially when a requested library version
  cannot support the required runtime

If verification fails because the template itself has a pre-existing unrelated
issue, preserve the scoped diff and establish baseline evidence before calling the failure pre-existing; report the failure,
and do not broaden into unrelated template repair unless the user asks.

## Common Pitfalls

- Updating `demo_go/go.mod` but not
  `templates/golang/{{cookiecutter.project_name}}/go.mod`
- Copying rendered imports like `github.com/TencentBlueKing/demo_go2/...` back
  into the template instead of restoring `{{cookiecutter.project_name}}`
- Validating by manually replacing only `project_name`; raw Jinja blocks can
  still break compilation unless cookiecutter renders the project
- Allowing Go to auto-switch toolchains and accidentally raising the template's
  Go version
- Running verification in the repo root when the generated project has its own
  Makefile
- Deleting template-only dependencies because one generated sandbox did not use
  them, without proving a fresh render no longer needs them
