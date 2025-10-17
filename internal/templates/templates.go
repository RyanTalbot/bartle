package templates

import _ "embed"

//go:embed conventional.yaml.tmpl
var Conventional string

//go:embed jira.yaml.tmpl
var Jira string

//go:embed custom.yaml.tmpl
var Custom string
