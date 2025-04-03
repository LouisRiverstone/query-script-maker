package language

import (
	"strings"
)

// KeywordMap contains common SQL keyword templates
var KeywordMap = map[string]string{
	"select":        "SELECT %s FROM %s",
	"select_where":  "SELECT %s FROM %s WHERE %s",
	"select_join":   "SELECT %s FROM %s JOIN %s ON %s.%s = %s.%s",
	"select_limit":  "SELECT %s FROM %s LIMIT %s",
	"count":         "SELECT COUNT(%s) FROM %s",
	"insert":        "INSERT INTO %s (%s) VALUES (%s)",
	"update":        "UPDATE %s SET %s WHERE %s",
	"delete":        "DELETE FROM %s WHERE %s",
	"create_table":  "CREATE TABLE %s (%s)",
	"alter_table":   "ALTER TABLE %s %s",
	"add_column":    "ALTER TABLE %s ADD COLUMN %s %s",
	"drop_column":   "ALTER TABLE %s DROP COLUMN %s",
	"create_index":  "CREATE INDEX %s ON %s (%s)",
	"drop_index":    "DROP INDEX %s",
	"group_by":      "SELECT %s FROM %s GROUP BY %s",
	"order_by":      "SELECT %s FROM %s ORDER BY %s",
	"distinct":      "SELECT DISTINCT %s FROM %s",
	"inner_join":    "SELECT %s FROM %s INNER JOIN %s ON %s",
	"left_join":     "SELECT %s FROM %s LEFT JOIN %s ON %s",
	"right_join":    "SELECT %s FROM %s RIGHT JOIN %s ON %s",
	"full_join":     "SELECT %s FROM %s FULL JOIN %s ON %s",
	"union":         "%s UNION %s",
	"union_all":     "%s UNION ALL %s",
	"case_when":     "CASE WHEN %s THEN %s ELSE %s END",
	"exists":        "SELECT %s FROM %s WHERE EXISTS (SELECT 1 FROM %s WHERE %s)",
	"not_exists":    "SELECT %s FROM %s WHERE NOT EXISTS (SELECT 1 FROM %s WHERE %s)",
	"between":       "%s BETWEEN %s AND %s",
	"in":            "%s IN (%s)",
	"having":        "SELECT %s FROM %s GROUP BY %s HAVING %s",
	"select_cte":    "WITH %s AS (%s) SELECT %s FROM %s",
	"window_func":   "SELECT %s, %s OVER (%s) FROM %s",
	"pivot":         "SELECT %s FROM %s PIVOT (%s FOR %s IN (%s))",
	"subquery":      "SELECT %s FROM (%s) AS %s",
	"offset":        "SELECT %s FROM %s LIMIT %s OFFSET %s",
	"join_using":    "SELECT %s FROM %s JOIN %s USING (%s)",
	"recursive_cte": "WITH RECURSIVE %s AS (%s UNION ALL %s) SELECT %s FROM %s",
	"json_extract":  "SELECT JSON_EXTRACT(%s, '%s') FROM %s",
	"date_part":     "SELECT DATE_PART('%s', %s) FROM %s",
	"regex_match":   "SELECT %s FROM %s WHERE %s REGEXP '%s'",
	"rollup":        "SELECT %s FROM %s GROUP BY ROLLUP(%s)",
	"cube":          "SELECT %s FROM %s GROUP BY CUBE(%s)",
}

// OperationPatternsEN contains English regex patterns to detect SQL operations
var OperationPatternsEN = map[string][]string{
	"select": {
		`(?i)(show|get|find|display|list|fetch|retrieve|search|query)`,
		`(?i)(what|which|who)`,
		`(?i)(select.*from)`,
	},
	"count": {
		`(?i)(count|how many|total number|tally|sum up)`,
		`(?i)(count\(.*\))`,
	},
	"insert": {
		`(?i)(add|create|insert|new|store|put)`,
		`(?i)(insert into)`,
	},
	"update": {
		`(?i)(update|modify|change|edit|alter)`,
		`(?i)(set.*where)`,
	},
	"delete": {
		`(?i)(delete|remove|drop|eliminate|erase)`,
		`(?i)(delete from)`,
	},
	"join": {
		`(?i)(join|combine|merge|related|linked|connected)`,
		`(?i)(together with)`,
		`(?i)(along with its)`,
	},
	"group": {
		`(?i)(group|aggregate|summarize|cluster|categorize)`,
		`(?i)(group by)`,
	},
	"order": {
		`(?i)(order|sort|arrange|rank|sequence)`,
		`(?i)(ascending|descending)`,
		`(?i)(order by)`,
	},
	"limit": {
		`(?i)(limit|top|first|latest|recent|newest|only|just)`,
		`(?i)(limit\s+\d+)`,
	},
	"distinct": {
		`(?i)(unique|distinct|different|separate)`,
		`(?i)(without duplicates)`,
	},
	"complex": {
		`(?i)(compare|correlate|trend|pattern|analyze)`,
		`(?i)(over time|per month|by year|per category)`,
	},
}

// OperationPatternsPT contains Portuguese regex patterns to detect SQL operations
var OperationPatternsPT = map[string][]string{
	"select": {
		`(?i)(mostrar|obter|encontrar|exibir|listar|buscar|recuperar|pesquisar|consultar)`,
		`(?i)(qual|quais|quem)`,
		`(?i)(selecionar.*de)`,
		`(?i)(selecione)`,
	},
	"count": {
		`(?i)(contar|quantos|número total|total de|soma de)`,
		`(?i)(contar\(.*\))`,
		`(?i)(quantidade de)`,
	},
	"insert": {
		`(?i)(adicionar|criar|inserir|novo|armazenar|colocar)`,
		`(?i)(inserir em)`,
		`(?i)(cadastrar)`,
		`(?i)(incluir)`,
	},
	"update": {
		`(?i)(atualizar|modificar|mudar|editar|alterar)`,
		`(?i)(definir.*onde)`,
		`(?i)(trocar)`,
	},
	"delete": {
		`(?i)(deletar|remover|excluir|eliminar|apagar)`,
		`(?i)(excluir de)`,
	},
	"join": {
		`(?i)(juntar|combinar|mesclar|relacionado|vinculado|conectado)`,
		`(?i)(junto com)`,
		`(?i)(junto de seus)`,
		`(?i)(relacionar com)`,
	},
	"group": {
		`(?i)(agrupar|agregar|resumir|categorizar)`,
		`(?i)(agrupar por)`,
		`(?i)(grupo)`,
	},
	"order": {
		`(?i)(ordenar|classificar|organizar|sequenciar)`,
		`(?i)(ascendente|descendente)`,
		`(?i)(ordem crescente|ordem decrescente)`,
		`(?i)(ordenar por)`,
	},
	"limit": {
		`(?i)(limitar|topo|primeiro|último|recente|mais novo|somente|apenas)`,
		`(?i)(limite\s+\d+)`,
		`(?i)(primeiros \d+)`,
	},
	"distinct": {
		`(?i)(único|distinto|diferente|separado)`,
		`(?i)(sem duplicatas)`,
		`(?i)(sem repetição)`,
	},
	"complex": {
		`(?i)(comparar|correlacionar|tendência|padrão|analisar)`,
		`(?i)(ao longo do tempo|por mês|por ano|por categoria)`,
	},
}

// AggregationFunctionsEN maps English natural language to SQL aggregate functions
var AggregationFunctionsEN = map[string]string{
	"average":            "AVG(%s)",
	"mean":               "AVG(%s)",
	"sum":                "SUM(%s)",
	"total":              "SUM(%s)",
	"maximum":            "MAX(%s)",
	"highest":            "MAX(%s)",
	"minimum":            "MIN(%s)",
	"lowest":             "MIN(%s)",
	"count":              "COUNT(%s)",
	"count distinct":     "COUNT(DISTINCT %s)",
	"unique count":       "COUNT(DISTINCT %s)",
	"standard deviation": "STDDEV(%s)",
	"variance":           "VARIANCE(%s)",
	"median":             "PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY %s)",
}

// AggregationFunctionsPT maps Portuguese natural language to SQL aggregate functions
var AggregationFunctionsPT = map[string]string{
	"média":             "AVG(%s)",
	"soma":              "SUM(%s)",
	"total":             "SUM(%s)",
	"máximo":            "MAX(%s)",
	"maior":             "MAX(%s)",
	"mínimo":            "MIN(%s)",
	"menor":             "MIN(%s)",
	"contagem":          "COUNT(%s)",
	"contar":            "COUNT(%s)",
	"contagem distinta": "COUNT(DISTINCT %s)",
	"contagem única":    "COUNT(DISTINCT %s)",
	"desvio padrão":     "STDDEV(%s)",
	"variância":         "VARIANCE(%s)",
	"mediana":           "PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY %s)",
}

// DateFunctionsEN maps English natural language to SQL date functions
var DateFunctionsEN = map[string]string{
	"year":         "EXTRACT(YEAR FROM %s)",
	"month":        "EXTRACT(MONTH FROM %s)",
	"day":          "EXTRACT(DAY FROM %s)",
	"quarter":      "EXTRACT(QUARTER FROM %s)",
	"week":         "EXTRACT(WEEK FROM %s)",
	"hour":         "EXTRACT(HOUR FROM %s)",
	"minute":       "EXTRACT(MINUTE FROM %s)",
	"second":       "EXTRACT(SECOND FROM %s)",
	"current date": "CURRENT_DATE()",
	"today":        "CURRENT_DATE()",
	"yesterday":    "DATE_SUB(CURRENT_DATE(), INTERVAL 1 DAY)",
	"tomorrow":     "DATE_ADD(CURRENT_DATE(), INTERVAL 1 DAY)",
	"this week":    "DATE_TRUNC('week', CURRENT_DATE())",
	"this month":   "DATE_TRUNC('month', CURRENT_DATE())",
	"this year":    "DATE_TRUNC('year', CURRENT_DATE())",
	"last week":    "DATE_SUB(DATE_TRUNC('week', CURRENT_DATE()), INTERVAL 1 WEEK)",
	"last month":   "DATE_SUB(DATE_TRUNC('month', CURRENT_DATE()), INTERVAL 1 MONTH)",
	"last year":    "DATE_SUB(DATE_TRUNC('year', CURRENT_DATE()), INTERVAL 1 YEAR)",
	"next week":    "DATE_ADD(DATE_TRUNC('week', CURRENT_DATE()), INTERVAL 1 WEEK)",
	"next month":   "DATE_ADD(DATE_TRUNC('month', CURRENT_DATE()), INTERVAL 1 MONTH)",
	"next year":    "DATE_ADD(DATE_TRUNC('year', CURRENT_DATE()), INTERVAL 1 YEAR)",
}

// DateFunctionsPT maps Portuguese natural language to SQL date functions
var DateFunctionsPT = map[string]string{
	"ano":            "EXTRACT(YEAR FROM %s)",
	"mês":            "EXTRACT(MONTH FROM %s)",
	"dia":            "EXTRACT(DAY FROM %s)",
	"trimestre":      "EXTRACT(QUARTER FROM %s)",
	"semana":         "EXTRACT(WEEK FROM %s)",
	"hora":           "EXTRACT(HOUR FROM %s)",
	"minuto":         "EXTRACT(MINUTE FROM %s)",
	"segundo":        "EXTRACT(SECOND FROM %s)",
	"data atual":     "CURRENT_DATE()",
	"hoje":           "CURRENT_DATE()",
	"ontem":          "DATE_SUB(CURRENT_DATE(), INTERVAL 1 DAY)",
	"amanhã":         "DATE_ADD(CURRENT_DATE(), INTERVAL 1 DAY)",
	"esta semana":    "DATE_TRUNC('week', CURRENT_DATE())",
	"esse mês":       "DATE_TRUNC('month', CURRENT_DATE())",
	"este mês":       "DATE_TRUNC('month', CURRENT_DATE())",
	"este ano":       "DATE_TRUNC('year', CURRENT_DATE())",
	"semana passada": "DATE_SUB(DATE_TRUNC('week', CURRENT_DATE()), INTERVAL 1 WEEK)",
	"última semana":  "DATE_SUB(DATE_TRUNC('week', CURRENT_DATE()), INTERVAL 1 WEEK)",
	"mês passado":    "DATE_SUB(DATE_TRUNC('month', CURRENT_DATE()), INTERVAL 1 MONTH)",
	"último mês":     "DATE_SUB(DATE_TRUNC('month', CURRENT_DATE()), INTERVAL 1 MONTH)",
	"ano passado":    "DATE_SUB(DATE_TRUNC('year', CURRENT_DATE()), INTERVAL 1 YEAR)",
	"último ano":     "DATE_SUB(DATE_TRUNC('year', CURRENT_DATE()), INTERVAL 1 YEAR)",
	"próxima semana": "DATE_ADD(DATE_TRUNC('week', CURRENT_DATE()), INTERVAL 1 WEEK)",
	"próximo mês":    "DATE_ADD(DATE_TRUNC('month', CURRENT_DATE()), INTERVAL 1 MONTH)",
	"próximo ano":    "DATE_ADD(DATE_TRUNC('year', CURRENT_DATE()), INTERVAL 1 YEAR)",
}

// StringFunctionsEN maps English natural language to SQL string functions
var StringFunctionsEN = map[string]string{
	"length":      "LENGTH(%s)",
	"uppercase":   "UPPER(%s)",
	"lowercase":   "LOWER(%s)",
	"capitalize":  "INITCAP(%s)",
	"trim":        "TRIM(%s)",
	"substring":   "SUBSTRING(%s, %d, %d)",
	"concatenate": "CONCAT(%s, %s)",
	"replace":     "REPLACE(%s, %s, %s)",
	"contains":    "%s LIKE '%%%s%%'",
	"starts with": "%s LIKE '%s%%'",
	"ends with":   "%s LIKE '%%%s'",
}

// StringFunctionsPT maps Portuguese natural language to SQL string functions
var StringFunctionsPT = map[string]string{
	"tamanho":         "LENGTH(%s)",
	"comprimento":     "LENGTH(%s)",
	"maiúsculas":      "UPPER(%s)",
	"maiúsculo":       "UPPER(%s)",
	"minúsculas":      "LOWER(%s)",
	"minúsculo":       "LOWER(%s)",
	"capitalizar":     "INITCAP(%s)",
	"remover espaços": "TRIM(%s)",
	"substring":       "SUBSTRING(%s, %d, %d)",
	"concatenar":      "CONCAT(%s, %s)",
	"substituir":      "REPLACE(%s, %s, %s)",
	"contém":          "%s LIKE '%%%s%%'",
	"começa com":      "%s LIKE '%s%%'",
	"termina com":     "%s LIKE '%%%s'",
}

// LanguageKeywords contains keywords for different languages to help detect the language
var LanguageKeywords = map[string][]string{
	"en": {"show", "get", "find", "what", "which", "how", "many", "where", "when", "list", "select"},
	"pt": {"mostrar", "obter", "encontrar", "qual", "quais", "como", "quantos", "onde", "quando", "listar", "selecionar"},
}

// GetOperationPatterns returns the appropriate operation patterns based on the detected language
func GetOperationPatterns(language string) map[string][]string {
	if language == "pt" {
		return OperationPatternsPT
	}
	return OperationPatternsEN
}

// GetAggregationFunctions returns the appropriate aggregation functions based on the detected language
func GetAggregationFunctions(language string) map[string]string {
	if language == "pt" {
		return AggregationFunctionsPT
	}
	return AggregationFunctionsEN
}

// GetDateFunctions returns the appropriate date functions based on the detected language
func GetDateFunctions(language string) map[string]string {
	if language == "pt" {
		return DateFunctionsPT
	}
	return DateFunctionsEN
}

// GetStringFunctions returns the appropriate string functions based on the detected language
func GetStringFunctions(language string) map[string]string {
	if language == "pt" {
		return StringFunctionsPT
	}
	return StringFunctionsEN
}

// DetectLanguage attempts to detect the language of the input prompt
func DetectLanguage(prompt string) string {
	// Normalize to lowercase
	prompt = strings.ToLower(prompt)

	// Count occurrences of language-specific keywords
	enMatches := 0
	ptMatches := 0

	for _, keyword := range LanguageKeywords["en"] {
		if strings.Contains(prompt, keyword) {
			enMatches++
		}
	}

	for _, keyword := range LanguageKeywords["pt"] {
		if strings.Contains(prompt, keyword) {
			ptMatches++
		}
	}

	// Compare and return the detected language
	if ptMatches > enMatches {
		return "pt"
	}

	// Default to English
	return "en"
}
