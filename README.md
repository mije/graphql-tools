# GraphQL Tools [![Build Status](https://travis-ci.org/mije/graphql-tools.svg?branch=master)](https://travis-ci.org/mije/graphql-tools) [![Go Report Card](https://goreportcard.com/badge/github.com/mije/graphql-tools)](https://goreportcard.com/report/github.com/mije/graphql-tools) [![codecov](https://codecov.io/gh/mije/graphql-tools/branch/master/graph/badge.svg)](https://codecov.io/gh/mije/graphql-tools)

Set of GraphQL tools to be used either as Go libraries or as CLI tools.

## Schema
Set of tools to process schemas.

### Compare
Computes diff between two schemas, categorize and assign severity to each detected change. This is useful to guard schema evolutions and to prevent introducing breaking changes (inspired by [graphql-schema_comparator](https://rubygems.org/gems/graphql-schema_comparator) Ruby gem).
