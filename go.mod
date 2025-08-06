// This module is used by GO to track program dependencies and import the modules that are required for the project to function.
// change this later to be in terms of github.com/modname. https://go.dev/doc/modules/layout
module TraceRoute

go 1.24.4

require golang.org/x/net v0.42.0

require golang.org/x/sys v0.34.0 // indirect
