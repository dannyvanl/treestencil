[![CI](https://github.com/dannyvanl/treestencil/actions/workflows/ci.yaml/badge.svg)](https://github.com/dannyvanl/treestencil/actions/workflows/ci.yaml)

# Treestencil

Treestencil is a command line utility that renders all templates in a directory to one or more target directories, keeping the template directory tree intact.
Variable values to apply when rendering can be defined per target.

A single yaml file is used to define the targets and their variable values, as well as global variable values and some other configuration options.

See this [example](./example/) to get an idea of how it works.


## Motivation

The original use cases that motivated the development of this tool was a terraform project using a directory per customer environment to provision. Each customer directory container multiple sub directories for different stages with their own terraform states.

Although terraform modules were heavily used to keep things DRY, the customer-specific files had a lot of duplication going on. Several wrappers for terraform were considered, but were not deemed fitting for the specific project.

Treestencil solved the duplication issue by providing a single yaml file to define customer-specific values for a bunch of variables that were then used to render a single directory tree with template files to a directory tree per customer. Although a debatable practice, checking in the generated files ensured that no changes were needed to the github workflows running terraform to provision the environments.


## TODO 

Although this tool can be considered as fully working for the intended purpose, it is still under development.

Some things that still need to be improved:

- Documentation how to use, with examples
- Tests
- Automated release action producing a binary executable for multiple architectures
- Reading the template directory once instead of rereading it for each target

