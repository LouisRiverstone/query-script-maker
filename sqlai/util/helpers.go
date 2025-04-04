package util

import (
	"fmt"
	"regexp"
	"strings"
)

// ContainsAny checks if str contains any of the keywords
func ContainsAny(str string, keywords []string) bool {
	for _, keyword := range keywords {
		if strings.Contains(str, keyword) {
			return true
		}
	}
	return false
}

// FormatValue formats values for SQL conditions
func FormatValue(value string, operator string) (string, string) {
	// Strip quotes if present
	if (strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) ||
		(strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) {
		value = value[1 : len(value)-1]
	}

	// Check if it's a number
	if regexp.MustCompile(`^\d+(\.\d+)?$`).MatchString(value) {
		return value, operator
	}

	// Handle boolean values
	if value == "true" || value == "false" {
		return value, operator
	}

	// Add quotes for string values
	return "'" + value + "'", operator
}

// FormatLikeValue formats values for LIKE conditions
func FormatLikeValue(prefix, suffix string) func(string, string) (string, string) {
	return func(value string, operator string) (string, string) {
		// Strip quotes if present
		if (strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) ||
			(strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) {
			value = value[1 : len(value)-1]
		}

		return "'" + prefix + value + suffix + "'", "LIKE"
	}
}

// FormatBetweenValue formats values for BETWEEN conditions
func FormatBetweenValue(value string, operator string) (string, string) {
	// This is more complex as BETWEEN needs two values
	// For now, just return placeholder
	return "value1 AND value2", "BETWEEN"
}

// FormatInValue formats values for IN conditions
func FormatInValue(value string, operator string) (string, string) {
	// Clean up and format the list of values
	values := strings.Split(value, ",")
	for i, v := range values {
		v = strings.TrimSpace(v)
		// Strip quotes if present
		if (strings.HasPrefix(v, "'") && strings.HasSuffix(v, "'")) ||
			(strings.HasPrefix(v, "\"") && strings.HasSuffix(v, "\"")) {
			v = v[1 : len(v)-1]
		}

		// Add quotes for string values
		if !regexp.MustCompile(`^\d+(\.\d+)?$`).MatchString(v) {
			v = "'" + v + "'"
		}

		values[i] = v
	}

	return "(" + strings.Join(values, ", ") + ")", "IN"
}

// FormatNullValue handles IS NULL and IS NOT NULL conditions
func FormatNullValue(value string, operator string) (string, string) {
	return "", operator
}

// FormatTimeRange formats relative time ranges
func FormatTimeRange(timeExpr string) (string, string) {
	// Extract the number and the time unit from the expression
	var num int
	var unit string
	parts := strings.Split(timeExpr, " ")
	if len(parts) >= 3 {
		fmt.Sscanf(parts[0], "%d", &num)
		unit = parts[1]

		// Normalize the unit
		switch strings.ToLower(unit) {
		case "day", "days":
			unit = "DAY"
		case "week", "weeks":
			unit = "WEEK"
		case "month", "months":
			unit = "MONTH"
		case "year", "years":
			unit = "YEAR"
		}

		return fmt.Sprintf("DATE_SUB(CURRENT_DATE(), INTERVAL %d %s)", num, unit), ">="
	}

	// Default
	return "DATE_SUB(CURRENT_DATE(), INTERVAL 30 DAY)", ">="
}

// FormatDateRange formats date ranges like "between last week and today"
func FormatDateRange(dateRange string) (string, string) {
	parts := strings.Split(dateRange, " and ")
	if len(parts) != 2 {
		// Default
		return "DATE_SUB(CURRENT_DATE(), INTERVAL 7 DAY) AND CURRENT_DATE()", "BETWEEN"
	}

	from := FormatDateExpression(strings.TrimSpace(parts[0]))
	to := FormatDateExpression(strings.TrimSpace(parts[1]))

	return from + " AND " + to, "BETWEEN"
}

// FormatSpecificDate formats specific date expressions like "today", "last month"
func FormatSpecificDate(dateExpr string) (string, string) {
	expr := FormatDateExpression(dateExpr)
	return expr, "="
}

// FormatDateExpression converts common date expressions to SQL functions
func FormatDateExpression(expr string) string {
	// Map common date expressions to SQL functions
	dateMap := map[string]string{
		"today":      "CURRENT_DATE()",
		"yesterday":  "DATE_SUB(CURRENT_DATE(), INTERVAL 1 DAY)",
		"tomorrow":   "DATE_ADD(CURRENT_DATE(), INTERVAL 1 DAY)",
		"this week":  "DATE_TRUNC('week', CURRENT_DATE())",
		"last week":  "DATE_SUB(DATE_TRUNC('week', CURRENT_DATE()), INTERVAL 1 WEEK)",
		"this month": "DATE_TRUNC('month', CURRENT_DATE())",
		"last month": "DATE_SUB(DATE_TRUNC('month', CURRENT_DATE()), INTERVAL 1 MONTH)",
		"this year":  "DATE_TRUNC('year', CURRENT_DATE())",
		"last year":  "DATE_SUB(DATE_TRUNC('year', CURRENT_DATE()), INTERVAL 1 YEAR)",
		"next week":  "DATE_ADD(DATE_TRUNC('week', CURRENT_DATE()), INTERVAL 1 WEEK)",
		"next month": "DATE_ADD(DATE_TRUNC('month', CURRENT_DATE()), INTERVAL 1 MONTH)",
		"next year":  "DATE_ADD(DATE_TRUNC('year', CURRENT_DATE()), INTERVAL 1 YEAR)",
	}

	// Portuguese date expressions
	ptDateMap := map[string]string{
		"hoje":           "CURRENT_DATE()",
		"ontem":          "DATE_SUB(CURRENT_DATE(), INTERVAL 1 DAY)",
		"amanhã":         "DATE_ADD(CURRENT_DATE(), INTERVAL 1 DAY)",
		"esta semana":    "DATE_TRUNC('week', CURRENT_DATE())",
		"semana passada": "DATE_SUB(DATE_TRUNC('week', CURRENT_DATE()), INTERVAL 1 WEEK)",
		"este mês":       "DATE_TRUNC('month', CURRENT_DATE())",
		"mês passado":    "DATE_SUB(DATE_TRUNC('month', CURRENT_DATE()), INTERVAL 1 MONTH)",
		"este ano":       "DATE_TRUNC('year', CURRENT_DATE())",
		"ano passado":    "DATE_SUB(DATE_TRUNC('year', CURRENT_DATE()), INTERVAL 1 YEAR)",
		"próxima semana": "DATE_ADD(DATE_TRUNC('week', CURRENT_DATE()), INTERVAL 1 WEEK)",
		"próximo mês":    "DATE_ADD(DATE_TRUNC('month', CURRENT_DATE()), INTERVAL 1 MONTH)",
		"próximo ano":    "DATE_ADD(DATE_TRUNC('year', CURRENT_DATE()), INTERVAL 1 YEAR)",
	}

	// First check English expressions
	if sql, ok := dateMap[strings.ToLower(expr)]; ok {
		return sql
	}

	// Then check Portuguese expressions
	if sql, ok := ptDateMap[strings.ToLower(expr)]; ok {
		return sql
	}

	return "CURRENT_DATE()"
}

// SortTablesByConfidence sorts tables by confidence score (descending)
func SortTablesByConfidence(tables []interface{}) {
	// Generic implementation assuming tables have a Confidence field
	for i := 0; i < len(tables)-1; i++ {
		for j := i + 1; j < len(tables); j++ {
			table1 := tables[i].(map[string]interface{})
			table2 := tables[j].(map[string]interface{})

			if table1["Confidence"].(float64) < table2["Confidence"].(float64) {
				tables[i], tables[j] = tables[j], tables[i]
			}
		}
	}
}

// ContainsString checks if a string slice contains a specific string
func ContainsString(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}

// ExtractTablesFromSQL extracts table names from a SQL query
func ExtractTablesFromSQL(sql string) []string {
	var tables []string

	// Extract from FROM clause
	fromRegex := regexp.MustCompile(`FROM\s+(\w+)`)
	if matches := fromRegex.FindStringSubmatch(sql); len(matches) > 1 {
		tables = append(tables, matches[1])
	}

	// Extract from JOIN clauses
	joinRegex := regexp.MustCompile(`JOIN\s+(\w+)`)
	joinMatches := joinRegex.FindAllStringSubmatch(sql, -1)

	for _, match := range joinMatches {
		if len(match) > 1 {
			tables = append(tables, match[1])
		}
	}

	return tables
}

// JaroWinklerSimilarity calculates the Jaro-Winkler similarity between two strings
// Returns a value between 0 (no similarity) and 1 (exact match)
func JaroWinklerSimilarity(s1, s2 string) float64 {
	// If both strings are empty, they're identical
	if len(s1) == 0 && len(s2) == 0 {
		return 1.0
	}

	// If one is empty, no similarity
	if len(s1) == 0 || len(s2) == 0 {
		return 0.0
	}

	// Calculate Jaro similarity first
	matchDistance := max(len(s1), len(s2))/2 - 1
	if matchDistance < 0 {
		matchDistance = 0
	}

	// Check for matches
	s1Matches := make([]bool, len(s1))
	s2Matches := make([]bool, len(s2))
	var matches float64 = 0

	for i := 0; i < len(s1); i++ {
		start := max(0, i-matchDistance)
		end := min(i+matchDistance+1, len(s2))

		for j := start; j < end; j++ {
			// If s2 already matched or characters don't match, skip
			if s2Matches[j] || s1[i] != s2[j] {
				continue
			}

			// Mark characters at this position as matched
			s1Matches[i] = true
			s2Matches[j] = true

			matches++
			break
		}
	}

	// If no matches found, return 0
	if matches == 0 {
		return 0.0
	}

	// Count transpositions
	var transpositions float64 = 0
	var point int = 0

	for i := 0; i < len(s1); i++ {
		if !s1Matches[i] {
			continue
		}

		// Find the next match in s2
		for !s2Matches[point] {
			point++
		}

		// If the characters don't match, count a transposition
		if s1[i] != s2[point] {
			transpositions++
		}

		point++
	}

	// Calculate Jaro similarity
	transpositions /= 2.0
	jaro := (matches/float64(len(s1)) + matches/float64(len(s2)) + (matches-transpositions)/matches) / 3.0

	// Calculate Jaro-Winkler similarity
	prefixLength := 0
	for i := 0; i < min(len(s1), len(s2), 4); i++ {
		if s1[i] == s2[i] {
			prefixLength++
		} else {
			break
		}
	}

	// Apply Winkler's prefix adjustment
	scalingFactor := 0.1 // standard scaling factor
	return jaro + (float64(prefixLength) * scalingFactor * (1.0 - jaro))
}

// IsCommonWord checks if the given word is a common word that should be ignored
// during entity extraction based on the detected language
func IsCommonWord(word string, language string) bool {
	// Common English words to ignore
	commonWordsEN := map[string]bool{
		"the": true, "a": true, "an": true, "and": true, "but": true,
		"or": true, "for": true, "nor": true, "on": true, "at": true,
		"to": true, "from": true, "by": true, "with": true, "about": true,
		"against": true, "between": true, "into": true, "through": true,
		"during": true, "before": true, "after": true, "above": true,
		"below": true, "under": true, "over": true, "of": true, "in": true,
		"out": true, "off": true, "up": true, "down": true, "all": true,
		"any": true, "each": true, "every": true, "some": true, "such": true,
		"that": true, "this": true, "these": true, "those": true, "which": true,
		"who": true, "whom": true, "whose": true, "what": true, "whatever": true,
		"when": true, "where": true, "why": true, "how": true, "because": true,
		"so": true, "very": true, "too": true, "not": true, "only": true,
		"just": true, "more": true, "most": true, "less": true, "least": true,
		"as": true, "if": true, "then": true, "else": true, "therefore": true,
		"thus": true, "however": true, "moreover": true, "nevertheless": true,
		"regardless": true, "furthermore": true, "whereas": true, "indeed": true,
		"already": true, "always": true, "never": true, "sometimes": true,
		"often": true, "seldom": true, "usually": true, "get": true, "select": true,
		"find": true, "show": true, "tell": true, "give": true,
		"list": true, "me": true, "us": true, "mine": true, "ours": true,
		"you": true, "yours": true, "he": true, "his": true, "she": true,
		"hers": true, "they": true, "them": true, "theirs": true, "it": true,
		"its": true, "table": true, "column": true, "row": true, "value": true,
		"field": true, "record": true, "database": true, "data": true,
		"sql": true, "query": true,
		"having": true, "join": true, "left": true, "right": true,
		"inner": true, "outer": true, "full": true, "count": true, "sum": true,
		"avg": true, "min": true, "max": true, "like": true,
		"group": true, "order": true, "limit": true, "offset": true, "can": true,
		"could": true, "would": true, "should": true, "may": true, "might": true,
	}

	// Common Portuguese words to ignore
	commonWordsPT := map[string]bool{
		"o": true, "a": true, "os": true, "as": true, "um": true, "uma": true,
		"uns": true, "umas": true, "de": true, "do": true, "da": true, "dos": true,
		"das": true, "no": true, "na": true, "nas": true, "ao": true,
		"aos": true, "à": true, "às": true, "pelo": true, "pela": true, "pelos": true,
		"pelas": true, "num": true, "numa": true, "nuns": true, "numas": true,
		"dum": true, "duma": true, "duns": true, "dumas": true, "com": true,
		"sem": true, "para": true, "por": true, "em": true, "sobre": true,
		"sob": true, "ante": true, "após": true, "até": true, "desde": true,
		"entre": true, "contra": true, "perante": true, "trás": true, "e": true,
		"mas": true, "ou": true, "que": true, "porque": true, "pois": true,
		"porém": true, "todavia": true, "contudo": true, "entretanto": true,
		"enquanto": true, "quando": true, "como": true, "onde": true, "se": true,
		"caso": true, "senão": true, "embora": true, "apesar": true, "eu": true,
		"tu": true, "ele": true, "ela": true, "nós": true, "vós": true, "eles": true,
		"elas": true, "vocês": true, "meu": true, "minha": true, "meus": true,
		"minhas": true, "teu": true, "tua": true, "teus": true, "tuas": true,
		"seu": true, "sua": true, "seus": true, "suas": true, "nosso": true,
		"nossa": true, "nossos": true, "nossas": true, "vosso": true, "vossa": true,
		"vossos": true, "vossas": true, "este": true, "esta": true, "estes": true,
		"estas": true, "esse": true, "essa": true, "esses": true, "essas": true,
		"aquele": true, "aquela": true, "aqueles": true, "aquelas": true,
		"isto": true, "isso": true, "aquilo": true, "já": true, "ainda": true,
		"sempre": true, "nunca": true, "jamais": true, "agora": true, "depois": true,
		"antes": true, "sim": true, "não": true, "talvez": true, "muito": true,
		"pouco": true, "mais": true, "menos": true, "demais": true, "tão": true,
		"quão": true, "quanto": true, "quanta": true, "quantos": true, "quantas": true,
		"todo": true, "toda": true, "todos": true, "todas": true, "algum": true,
		"alguma": true, "alguns": true, "algumas": true, "nenhum": true, "nenhuma": true,
		"nenhuns": true, "nenhumas": true, "qualquer": true, "quaisquer": true,
		"me": true, "te": true, "lhe": true, "vos": true, "lhes": true,
		"tabela": true, "coluna": true, "linha": true, "valor": true, "campo": true,
		"registro": true, "banco": true, "dados": true, "dado": true, "quero": true,
		"mostre": true, "exiba": true, "liste": true, "selecione": true, "busque": true,
		"encontre": true, "obtenha": true, "traga": true, "consulte": true,
	}

	// Check if it's a common word based on the language
	if language == "pt" {
		return commonWordsPT[word]
	}

	// Default to English
	return commonWordsEN[word]
}

// Helper functions for JaroWinklerSimilarity
func min(values ...int) int {
	result := values[0]
	for _, v := range values[1:] {
		if v < result {
			result = v
		}
	}
	return result
}

func max(values ...int) int {
	result := values[0]
	for _, v := range values[1:] {
		if v > result {
			result = v
		}
	}
	return result
}

// Contains checks if a string slice contains a specific string
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
