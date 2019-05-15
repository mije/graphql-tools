[![Build Status](https://travis-ci.org/mije/graphql-tools.svg?branch=master)](https://travis-ci.org/mije/graphql-tools) [![Go Report Card](https://goreportcard.com/badge/github.com/mije/graphql-tools)](https://goreportcard.com/report/github.com/mije/graphql-tools) [![codecov](https://codecov.io/gh/mije/graphql-tools/branch/master/graph/badge.svg)](https://codecov.io/gh/mije/graphql-tools)

# GraphQL Tools
Set of GraphQL tools to be used either as a Go libraries within your projects or from command line as executable binaries.

## Schema
A set of tools to process GrephQL schema.  

### Compare
Computes a diff between to schemas serialised using SDL and category and assign severity to each change. This may be useful to guard GraphQL schema evolutions and prevent breaking changes. 
It is inspired by similar [Ruby gem](https://rubygems.org/gems/graphql-schema_comparator). Which provides a great categorization of possible changes but does not provide a structured output.
