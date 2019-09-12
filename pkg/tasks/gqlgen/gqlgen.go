package gqlgen

import (
	"github.com/staffano/sonbyg/pkg/builder"
	"github.com/staffano/sonbyg/pkg/tasks/packages/gogo"
	"github.com/staffano/sonbyg/pkg/utils"
)

// Gqlgen runs github.com/99designs/gqlgen in the cwd directory.
// Further configuration of the command execution is done in the
// file ./gqlgen.yml
func Gqlgen(b *builder.Builder, vars *utils.Variables, cwd string) *builder.Task {
	(*vars)["CWD"] = cwd
	delete(*vars, "GOBIN")
	return gogo.Run(b, vars, "github.com/99designs/gqlgen", "-v")
}
