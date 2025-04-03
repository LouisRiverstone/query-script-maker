package operations

import (
	"regexp"
	"strings"

	"sql_script_maker/sqlai/models"
	"sql_script_maker/sqlai/util"
)

// ExtractConditions extracts WHERE conditions from a prompt
func ExtractConditions(prompt string, tables []models.TableInfo, columns []models.ColumnInfo, language string) []models.Condition {
	var conditions []models.Condition

	// Common condition patterns (English)
	enConditionPatterns := []struct {
		pattern   string
		operator  string
		formatter func(string, string) (string, string)
	}{
		{`(\w+)\s+equals?\s+([\w\d]+|"[^"]*"|'[^']*')`, "=", util.FormatValue},
		{`(\w+)\s+is\s+([\w\d]+|"[^"]*"|'[^']*')`, "=", util.FormatValue},
		{`(\w+)\s+>(?:=)?\s+([\w\d]+|"[^"]*"|'[^']*')`, ">=", util.FormatValue},
		{`(\w+)\s+<(?:=)?\s+([\w\d]+|"[^"]*"|'[^']*')`, "<=", util.FormatValue},
		{`(\w+)\s+>\s+([\w\d]+|"[^"]*"|'[^']*')`, ">", util.FormatValue},
		{`(\w+)\s+<\s+([\w\d]+|"[^"]*"|'[^']*')`, "<", util.FormatValue},
		{`(\w+)\s+!=\s+([\w\d]+|"[^"]*"|'[^']*')`, "!=", util.FormatValue},
		{`(\w+)\s+not\s+equals?\s+([\w\d]+|"[^"]*"|'[^']*')`, "!=", util.FormatValue},
		{`(\w+)\s+contains\s+([\w\d]+|"[^"]*"|'[^']*')`, "LIKE", util.FormatLikeValue("%", "%")},
		{`(\w+)\s+starts\s+with\s+([\w\d]+|"[^"]*"|'[^']*')`, "LIKE", util.FormatLikeValue("", "%")},
		{`(\w+)\s+ends\s+with\s+([\w\d]+|"[^"]*"|'[^']*')`, "LIKE", util.FormatLikeValue("%", "")},
		{`(\w+)\s+between\s+([\w\d]+|"[^"]*"|'[^']*')\s+and\s+([\w\d]+|"[^"]*"|'[^']*')`, "BETWEEN", util.FormatBetweenValue},
		{`(\w+)\s+in\s+\(([\w\d",'\s]+)\)`, "IN", util.FormatInValue},
		{`(\w+)\s+is\s+null`, "IS NULL", util.FormatNullValue},
		{`(\w+)\s+is\s+not\s+null`, "IS NOT NULL", util.FormatNullValue},
	}

	// Portuguese condition patterns
	ptConditionPatterns := []struct {
		pattern   string
		operator  string
		formatter func(string, string) (string, string)
	}{
		{`(\w+)\s+(?:é|igual a)\s+([\w\d]+|"[^"]*"|'[^']*')`, "=", util.FormatValue},
		{`(\w+)\s+(?:está|sendo)\s+([\w\d]+|"[^"]*"|'[^']*')`, "=", util.FormatValue},
		{`(\w+)\s+(?:maior ou igual a|>=)\s+([\w\d]+|"[^"]*"|'[^']*')`, ">=", util.FormatValue},
		{`(\w+)\s+(?:menor ou igual a|<=)\s+([\w\d]+|"[^"]*"|'[^']*')`, "<=", util.FormatValue},
		{`(\w+)\s+(?:maior que|>)\s+([\w\d]+|"[^"]*"|'[^']*')`, ">", util.FormatValue},
		{`(\w+)\s+(?:menor que|<)\s+([\w\d]+|"[^"]*"|'[^']*')`, "<", util.FormatValue},
		{`(\w+)\s+(?:!=|<>|diferente de|não igual a)\s+([\w\d]+|"[^"]*"|'[^']*')`, "!=", util.FormatValue},
		{`(\w+)\s+(?:contém|contem|possui)\s+([\w\d]+|"[^"]*"|'[^']*')`, "LIKE", util.FormatLikeValue("%", "%")},
		{`(\w+)\s+começa\s+com\s+([\w\d]+|"[^"]*"|'[^']*')`, "LIKE", util.FormatLikeValue("", "%")},
		{`(\w+)\s+termina\s+com\s+([\w\d]+|"[^"]*"|'[^']*')`, "LIKE", util.FormatLikeValue("%", "")},
		{`(\w+)\s+entre\s+([\w\d]+|"[^"]*"|'[^']*')\s+e\s+([\w\d]+|"[^"]*"|'[^']*')`, "BETWEEN", util.FormatBetweenValue},
		{`(\w+)\s+em\s+\(([\w\d",'\s]+)\)`, "IN", util.FormatInValue},
		{`(\w+)\s+(?:é nulo|é null|não existe)`, "IS NULL", util.FormatNullValue},
		{`(\w+)\s+(?:não é nulo|não é null|existe)`, "IS NOT NULL", util.FormatNullValue},
	}

	// Choose the appropriate pattern set based on language
	var conditionPatterns []struct {
		pattern   string
		operator  string
		formatter func(string, string) (string, string)
	}

	if language == "pt" {
		conditionPatterns = ptConditionPatterns
	} else {
		conditionPatterns = enConditionPatterns
	}

	// Extract basic conditions
	for _, pattern := range conditionPatterns {
		re := regexp.MustCompile(pattern.pattern)
		matches := re.FindAllStringSubmatch(prompt, -1)

		for _, match := range matches {
			if len(match) >= 3 {
				colName := match[1]
				colValue := match[2]

				// Find table for this column
				tableName := ""
				for _, col := range columns {
					if strings.EqualFold(col.Name, colName) {
						tableName = col.TableName
						break
					}
				}

				// If table not found but we have tables, use the first one
				if tableName == "" && len(tables) > 0 {
					tableName = tables[0].Name
				}

				// Format the value based on the condition type
				formattedValue, operator := pattern.formatter(colValue, pattern.operator)

				// Determine conjunction (AND/OR)
				conjunction := "AND"
				if language == "pt" {
					if regexp.MustCompile(`\b(ou|um dos)\b`).MatchString(prompt) {
						conjunction = "OR"
					}
				} else {
					if regexp.MustCompile(`\b(or|either)\b`).MatchString(prompt) {
						conjunction = "OR"
					}
				}

				conditions = append(conditions, models.Condition{
					ColumnName:  colName,
					TableName:   tableName,
					Operator:    operator,
					Value:       formattedValue,
					Conjunction: conjunction,
				})
			}
		}
	}

	// Extract date-specific conditions
	dateConditions := extractDateConditions(prompt, tables, columns, language)
	conditions = append(conditions, dateConditions...)

	// Look for complex conditions like "where status is active and (price > 100 or quantity < 5)"
	var complexPatterns []string
	if language == "pt" {
		complexPatterns = []string{
			`onde\s+(.+?)\s+(e|ou)\s+\((.+?)\)`,
			`com\s+(.+?)\s+(e|ou)\s+\((.+?)\)`,
		}
	} else {
		complexPatterns = []string{
			`where\s+(.+?)\s+(and|or)\s+\((.+?)\)`,
			`with\s+(.+?)\s+(and|or)\s+\((.+?)\)`,
		}
	}

	for _, pattern := range complexPatterns {
		re := regexp.MustCompile(pattern)
		if matches := re.FindStringSubmatch(prompt); len(matches) > 3 {
			// First extract the simple condition
			simpleCondExpr := matches[1]
			conjunction := strings.ToUpper(matches[2])
			complexExpr := matches[3]

			conditions = append(conditions, models.Condition{
				IsComplex:   true,
				ComplexExpr: "(" + simpleCondExpr + ") " + conjunction + " (" + complexExpr + ")",
			})
		}
	}

	return conditions
}

// extractDateConditions extracts date-specific conditions from the prompt
func extractDateConditions(prompt string, tables []models.TableInfo, columns []models.ColumnInfo, language string) []models.Condition {
	var conditions []models.Condition

	// Date-specific condition patterns
	var datePatterns []struct {
		pattern   string
		operator  string
		formatter func(string) (string, string)
	}

	if language == "pt" {
		datePatterns = []struct {
			pattern   string
			operator  string
			formatter func(string) (string, string)
		}{
			{`(\w+)\s+nos\s+últimos\s+(\d+)\s+(dia|dias|semana|semanas|mês|meses|ano|anos)`, ">=", util.FormatTimeRange},
			{`(\w+)\s+entre\s+(hoje|ontem|esta semana|semana passada|este mês|mês passado|este ano|ano passado)\s+e\s+(hoje|ontem|esta semana|semana passada|este mês|mês passado|este ano|ano passado)`, "BETWEEN", util.FormatDateRange},
			{`(\w+)\s+(?:é|igual a)\s+(hoje|ontem|esta semana|semana passada|este mês|mês passado|este ano|ano passado)`, "=", util.FormatSpecificDate},
		}
	} else {
		datePatterns = []struct {
			pattern   string
			operator  string
			formatter func(string) (string, string)
		}{
			{`(\w+)\s+in\s+the\s+last\s+(\d+)\s+(day|days|week|weeks|month|months|year|years)`, ">=", util.FormatTimeRange},
			{`(\w+)\s+between\s+(today|yesterday|this week|last week|this month|last month|this year|last year)\s+and\s+(today|yesterday|this week|last week|this month|last month|this year|last year)`, "BETWEEN", util.FormatDateRange},
			{`(\w+)\s+is\s+(today|yesterday|this week|last week|this month|last month|this year|last year)`, "=", util.FormatSpecificDate},
		}
	}

	// Extract date-specific conditions
	for _, pattern := range datePatterns {
		re := regexp.MustCompile(pattern.pattern)
		matches := re.FindAllStringSubmatch(prompt, -1)

		for _, match := range matches {
			if len(match) >= 3 {
				colName := match[1]

				// Find table for this column
				tableName := ""
				for _, col := range columns {
					if strings.EqualFold(col.Name, colName) {
						tableName = col.TableName
						break
					}
				}

				// If table not found but we have tables, use the first one
				if tableName == "" && len(tables) > 0 {
					tableName = tables[0].Name
				}

				// Format the value based on the date pattern
				formattedValue, operator := pattern.formatter(strings.Join(match[2:], " "))

				// Determine conjunction (AND/OR)
				conjunction := "AND"
				if language == "pt" {
					if regexp.MustCompile(`\b(ou|um dos)\b`).MatchString(prompt) {
						conjunction = "OR"
					}
				} else {
					if regexp.MustCompile(`\b(or|either)\b`).MatchString(prompt) {
						conjunction = "OR"
					}
				}

				conditions = append(conditions, models.Condition{
					ColumnName:  colName,
					TableName:   tableName,
					Operator:    operator,
					Value:       formattedValue,
					Conjunction: conjunction,
				})
			}
		}
	}

	return conditions
}
