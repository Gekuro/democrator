package tools

// explicit import for build tools, prevents go mod tidy from removing their dependencies from go.mod and go.sum

import (
	_ "github.com/99designs/gqlgen"
)