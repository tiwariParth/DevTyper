package golang

var (
	Templates = []string{
		"func %s() %s {",
		"var %s %s = %s",
		"type %s struct { %s %s }",
		"if %s != nil { return %s }",
		"for %s := range %s {",
		"switch %s := %s.(type) {",
		"map[%s]%s{%s: %s}",
		"func (%s *%s) %s() %s {",
	}

	Variables  = []string{"err", "val", "data", "result", "item", "obj", "ctx"}
	Types      = []string{"string", "int", "bool", "error", "interface{}"}
	Operations = []string{"nil", "true", "false", "0", "1", "\"\""}
)
