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
