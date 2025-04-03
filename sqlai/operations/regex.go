package operations

import (
	"regexp"
)

// CompileRegex compiles a regular expression and returns it
func CompileRegex(pattern string) *regexp.Regexp {
	return regexp.MustCompile(pattern)
}

// RegexMatch checks if a pattern matches a string
func RegexMatch(pattern, str string) bool {
	re := regexp.MustCompile(pattern)
	return re.MatchString(str)
}

// GetOperationKeywords returns keyword sets for detecting operations based on language
func GetOperationKeywords(language string) map[string][]string {
	if language == "pt" {
		return map[string][]string{
			"select": {
				"mostrar", "exibir", "listar", "buscar", "selecionar", "obter", "encontrar",
				"recuperar", "pesquisar", "consultar", "qual", "quais", "quem",
			},
			"count": {
				"contar", "quantos", "quantas", "número total", "total de", "soma de", "quantidade", "quantificar",
			},
			"insert": {
				"adicionar", "criar", "inserir", "novo", "armazenar", "colocar", "cadastrar", "incluir", "salvar",
			},
			"update": {
				"atualizar", "modificar", "mudar", "editar", "alterar", "trocar", "corrigir",
			},
			"delete": {
				"deletar", "remover", "excluir", "eliminar", "apagar", "retirar",
			},
			"join": {
				"juntar", "combinar", "mesclar", "relacionado", "vinculado", "conectado", "junto com", "relacionar",
			},
			"group": {
				"agrupar", "agregar", "resumir", "categorizar", "grupo", "classificar por grupo",
			},
			"order": {
				"ordenar", "classificar", "organizar", "sequenciar", "ascendente", "descendente", "ordem crescente", "ordem decrescente",
			},
			"limit": {
				"limitar", "limite", "topo", "primeiro", "primeiros", "últimos", "recentes", "mais recentes", "somente", "apenas",
			},
			"distinct": {
				"único", "distintos", "diferentes", "separados", "sem duplicatas", "sem repetição",
			},
		}
	}

	// Default to English
	return map[string][]string{
		"select": {
			"show", "display", "list", "get", "find", "select", "retrieve",
			"search", "query", "what", "which", "who",
		},
		"count": {
			"count", "how many", "total number", "tally", "sum up", "quantity", "quantify",
		},
		"insert": {
			"add", "create", "insert", "new", "store", "put", "save",
		},
		"update": {
			"update", "modify", "change", "edit", "alter", "adjust", "correct",
		},
		"delete": {
			"delete", "remove", "drop", "eliminate", "erase", "clear",
		},
		"join": {
			"join", "combine", "merge", "related", "linked", "connected", "with", "and their",
		},
		"group": {
			"group", "aggregate", "summarize", "cluster", "categorize", "by type", "by category",
		},
		"order": {
			"order", "sort", "arrange", "rank", "sequence", "ascending", "descending",
		},
		"limit": {
			"limit", "top", "first", "latest", "recent", "newest", "only", "just",
		},
		"distinct": {
			"unique", "distinct", "different", "separate", "without duplicates",
		},
	}
}
