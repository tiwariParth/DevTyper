package main

var (
	keywords = []string{
		"func", "var", "const", "type", "struct", "interface", "map", "chan", 
		"package", "import", "return", "defer", "go", "select", "case", "break",
		"continue", "default", "else", "fallthrough", "for", "goto", "if", "range",
		"switch", "try", "catch", "finally", "throw", "throws", "class", "extends",
	}

	dataTypes = []string{
		"string", "int", "bool", "float64", "byte", "rune", "uint", "int64",
		"float32", "complex64", "complex128", "uint32", "int32", "uint64",
		"error", "void", "null", "undefined", "number", "object", "array",
	}

	methods = []string{
		"len", "make", "new", "append", "delete", "copy", "close", "panic",
		"recover", "print", "println", "scanf", "sprintf", "printf", "isEmpty",
		"toString", "valueOf", "length", "push", "pop", "shift", "unshift",
	}

	concepts = []string{
		"goroutine", "channel", "mutex", "pointer", "slice", "constructor",
		"inheritance", "polymorphism", "encapsulation", "abstraction",
		"interface", "implementation", "callback", "promise", "async", "await",
	}
)

type WordLibrary struct {
	categories map[string][]string
	allWords   []string
}

func NewWordLibrary() *WordLibrary {
	lib := &WordLibrary{
		categories: map[string][]string{
			"keywords":  keywords,
			"dataTypes": dataTypes,
			"methods":   methods,
			"concepts":  concepts,
		},
		allWords: make([]string, 0),
	}

	// Combine all categories into allWords
	for _, words := range lib.categories {
		lib.allWords = append(lib.allWords, words...)
	}

	return lib
}

func (w *WordLibrary) GetAllWords() []string {
	return w.allWords
}

func (w *WordLibrary) GetCategoryWords(category string) []string {
	return w.categories[category]
}
