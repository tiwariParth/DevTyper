package rust

var (
	Templates = []string{
		"fn %s(%s: %s) -> %s {",
		"let mut %s: %s = %s;",
		"struct %s<%s> { %s: %s }",
		"impl %s for %s {",
		"match %s {",
		"if let Some(%s) = %s {",
		"pub fn %s(&self) -> Result<%s, %s> {",
	}

	Variables  = []string{"err", "val", "data", "result", "item", "cfg", "ctx"}
	Types      = []string{"String", "i32", "bool", "Option", "Result", "Vec"}
	Operations = []string{"None", "Some", "Ok", "Err", "true", "false", "0", "\"\""}
)
