# rehearse
Markdown-native test framework that turns specifications into executable scenarios. Write human-readable steps, compose acceptance criteria, run with `rehearse run`. No DSL, no glue code — your specs are your tests.

## How is Rehearse different?

Good tools exist in this space. [Gauge](https://gauge.org/) brought markdown to specification testing. [Cucumber](https://cucumber.io/) proved that human-readable tests matter. [Concordion](https://concordion.org/) pioneered living documentation. We respect this work — Rehearse builds on the ideas they validated.

But every one of them splits your test into two artifacts: a readable spec *here*, executable code *over there*. A product owner reads the spec. A developer maintains the fixtures. They drift apart. Nobody notices until something breaks.

Rehearse refuses that split. A scenario is one markdown file where the description *is* the test. Bash, Python, SQL, and Starlark blocks execute directly — no fixture classes, no step definitions, no glue layer. Acceptance criteria are reusable verification units composed across workflows, not reimplemented in each test.

The result: a format that a product owner reads on GitHub, a developer runs in CI, and an AI agent authors fluently — all the same file, all at once. Best experience for humans and AI agents in one package: easy to craft with AI help, easy to understand without it.

See the [full competitive analysis](docs/competitors.md) for a detailed comparison.
