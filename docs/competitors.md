# Competitive Landscape

How Rehearse compares to existing tools in the executable specification and markdown testing space.

## Direct Competitors

### Gauge (ThoughtWorks)

[gauge.org](https://gauge.org/) | [GitHub](https://github.com/getgauge/gauge)

The closest competitor. Open-source, markdown-based specification testing from ThoughtWorks. Uses `.spec` or `.md` files with step-driven structure, structured reporting, parallel execution, and multi-language support (Java, C#, JS, Python, Ruby).

**Key difference:** Gauge separates intent from implementation — step text in markdown maps to step definitions written in code. Rehearse keeps both in one file: the markdown describes the step, the code block executes it. No glue code, no fixture layer.

**Minimal test definition:**

```markdown
# Search specification

## Search for a term
* Navigate to "https://example.com"
* Search for "testing"
* Verify results contain "testing"
```

Each `*` step requires a matching function in a separate code file:

```java
// StepImplementation.java
@Step("Search for <term>")
public void searchFor(String term) {
    // implementation here
}
```

Rehearse equivalent — everything in one file:

```markdown
## search-for-term

```bash
curl -s "https://example.com/search?q=testing" | grep -q "testing"
`` `
```

---

### Concordion

[concordion.org](https://concordion.org/) | [GitHub](https://github.com/concordion/concordion)

Mature Java framework for "living documentation" — executable specifications that double as docs. Originally HTML-only, now supports Markdown. Commands are embedded via markdown link syntax.

**Key difference:** Java-centric. Commands use markdown link notation (`[value](- "command")`), which is clever but non-obvious. Requires Java fixture classes to back every specification.

**Minimal test definition:**

```markdown
# Greeting

The greeting for user [Bob](- "#name") should be
[Hello Bob!](- "?=greetingFor(#name)")
```

Backed by:

```java
public class GreetingFixture {
    public String greetingFor(String name) {
        return "Hello " + name + "!";
    }
}
```

---

### DaSpec

[GitHub](https://github.com/daspec/daspec-js) | [npm](https://www.npmjs.com/package/daspec)

Executable markdown specifications for JavaScript. Validates documents against software, inspired by Concordion. Announced by Gojko Adzic (author of "Specification by Example").

**Key difference:** JavaScript/Node.js only. Small community, appears low-activity. Follows the fixture pattern — markdown references functions defined elsewhere.

**Minimal test definition:**

```markdown
# Check full name

| First name | Last name | Full name     |
| ---------- | --------- | ------------- |
| John       | Lennon    | John Lennon   |
| Paul        | McCartney | Paul McCartney |
```

Backed by:

```javascript
defineStep(/Check full name/, function (firstCol, lastCol, fullCol) {
    return firstCol + ' ' + lastCol;
});
```

---

### Silk

[GitHub](https://github.com/matryer/silk)

Markdown-driven API testing for Go by Mat Ryer. Clean, minimal format — headings describe requests, code blocks define expectations.

**Key difference:** API testing only. Not a general-purpose workflow framework. No acceptance criteria composition, no multi-step data flow, no parallel execution.

**Minimal test definition:**

```markdown
# Check API

## GET /api/users

* Status: 200
* Content-Type: "application/json"

```json
{
  "users": [{"name": "John"}]
}
`` `
```

---

## Code Block Executors (Partial Overlap)

These tools run code blocks from markdown but lack test composition, acceptance criteria, structured reporting, or workflow orchestration.

### mdsh

[GitHub](https://github.com/zimbatm/mdsh) (Rust) | [GitHub](https://github.com/bashup/mdsh) (Bash)

Markdown shell preprocessor. Executes shell blocks in markdown files. The `--frozen` flag verifies output hasn't changed — useful for keeping docs in sync.

**Key difference:** Documentation synchronization tool, not a test framework. No structured pass/fail reporting, no step composition, no AC framework.

**Minimal usage:**

````markdown
```bash
echo "Hello, world!"
```

```
Hello, world!
```
````

mdsh runs the bash block and verifies the output block matches.

---

### txm (Tests eX Markdown)

[GitHub](https://github.com/anko/txm)

Language-agnostic markdown code block tester. Uses HTML comment annotations to define test expectations.

**Key difference:** Annotation-driven (HTML comments), not structure-driven (headings/tables). No workflows, no composition, no reporting beyond pass/fail per block.

**Minimal test definition:**

```markdown
<!-- !test program bash -->

<!-- !test check greeting -->

    echo "hello world"

<!-- !test out greeting -->

    hello world
```

---

### MDX (Go)

[GitHub](https://github.com/mjbozo/mdx)

Recent Go tool for executing markdown code blocks. Early-stage, minimal feature set.

**Key difference:** Simple executor with no test framework features. No structured reporting, composition, or AC model.

---

### markdown-exec (Python)

[GitHub](https://github.com/pawamoy/markdown-exec) | [PyPI](https://pypi.org/project/markdown-exec/)

MkDocs plugin that executes code blocks in markdown during documentation build. Supports Python, bash, and other languages.

**Key difference:** Documentation build tool, not a test framework. Tied to MkDocs ecosystem.

---

## BDD Frameworks (Different Paradigm)

### Cucumber / Gherkin

[cucumber.io](https://cucumber.io/) | [GitHub](https://github.com/cucumber)

The dominant BDD framework. `.feature` files use Given/When/Then syntax (Gherkin DSL). Step definitions live in separate code files.

**Key difference:** Custom DSL, not markdown. Does not render on GitHub as documentation. Requires step definition glue code in a host language. Well-established ecosystem but high ceremony.

**Minimal test definition:**

```gherkin
Feature: Greeting
  Scenario: Say hello
    Given the system is running
    When I request a greeting for "Bob"
    Then the response should be "Hello Bob!"
```

Backed by:

```javascript
Given('the system is running', function () { /* ... */ });
When('I request a greeting for {string}', function (name) { /* ... */ });
Then('the response should be {string}', function (expected) { /* ... */ });
```

---

### FitNesse

[fitnesse.org](http://fitnesse.org/) | [GitHub](https://github.com/unclebob/fitnesse)

Wiki-based acceptance testing framework by Ward Cunningham / Uncle Bob. Table-driven specifications with Java fixture backing.

**Key difference:** Wiki format (not markdown), table-driven, Java-centric. The grandparent of "specs as tests" — proved the concept but hasn't evolved with modern workflows.

**Minimal test definition:**

```
|Greeting|
|name|greeting?|
|Bob|Hello Bob!|
```

Backed by a Java fixture class.

---

## Comparison Matrix

| Capability | Rehearse | Gauge | Concordion | Cucumber | Silk | mdsh | txm |
|---|---|---|---|---|---|---|---|
| Markdown-native | Yes | Yes | Yes | No (Gherkin) | Yes | Yes | Yes |
| Direct code execution | Yes | No (fixtures) | No (fixtures) | No (step defs) | Partial | Yes | Yes |
| Multi-step workflows | Yes | Yes | Yes | Yes | No | No | No |
| AC composition | Yes | No | No | No | No | No | No |
| Structured reporting | Yes | Yes | Yes | Yes | No | No | No |
| Data flow between steps | Yes | Partial | No | Partial | No | No | No |
| Parallel execution | Yes | Yes | No | No | No | No | No |
| Sub-flow includes | Yes | Partial | No | No | No | No | No |
| SQL verification | Yes | Via plugins | Via fixtures | Via step defs | No | No | No |
| No custom DSL | Yes | Yes | Link syntax | Gherkin DSL | Yes | Yes | Comment annotations |
| Renders on GitHub | Yes | Yes | Yes | No | Yes | Yes | Yes |
| Language-agnostic scripts | Bash/Python/SQL/Starlark | Host language | Java | Host language | Go | Shell | Various |

## Rehearse's Unique Position

No existing tool combines these three properties:

1. **Direct execution** — bash, Python, SQL, and Starlark blocks run as-is from the markdown. No fixture classes, no step definitions, no glue code.
2. **Acceptance criteria as composable units** — reusable verification blocks referenced across scenarios. One AC, many workflows.
3. **Pure markdown** — renders on GitHub, readable by non-developers, authored by AI agents. The spec IS the test.

Gauge comes closest but requires a fixture code layer between the spec and execution. Everything else is either a simple code-block executor (no composition or reporting) or a legacy framework (not markdown-native).

## Outstanding Questions

None at this time.
