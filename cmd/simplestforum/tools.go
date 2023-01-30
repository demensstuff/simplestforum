//go:build tools

package main

import (
	_ "github.com/99designs/gqlgen"
)

// This file is needed to load gqlgen which generates the resolver but is not explicitly used anywhere
