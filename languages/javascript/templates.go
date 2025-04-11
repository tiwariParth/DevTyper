package javascript

var (
	Templates = []string{
		"function %s(%s) {",
		"const %s = %s",
		"let %s = %s",
		"class %s extends %s {",
		"if (%s === %s) {",
		"for (let %s of %s) {",
		"async function %s(%s) {",
		"try { %s } catch(%s) {",
	}

	Variables  = []string{"err", "data", "result", "item", "obj", "ctx", "response"}
	Types      = []string{"Array", "Object", "string", "number", "boolean"}
	Operations = []string{"null", "undefined", "true", "false", "0", "''", "[]", "{}"}
)
