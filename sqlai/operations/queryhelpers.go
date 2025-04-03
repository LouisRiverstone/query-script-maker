package operations

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"sql_script_maker/sqlai/models"
	"sql_script_maker/sqlai/util"
)

// Helper methods for QueryBuilder

// buildJoinClause constructs a JOIN clause between tables
func (qb *QueryBuilder) buildJoinClause(tables []models.TableInfo, prompt string) string {
	if len(tables) < 2 {
		return tables[0].Name
	}

	// Determine join type
	joinType := qb.determineJoinType(prompt)

	// Start with the first table
	fromClause := tables[0].Name

	// Add joins for remaining tables
	for i := 1; i < len(tables); i++ {
		// Find the join condition between these tables
		joinCondition := FindJoinCondition(tables[0].Name, tables[i].Name, qb.dbTables)

		if joinCondition != "" {
			fromClause += fmt.Sprintf(" %s %s ON %s", joinType, tables[i].Name, joinCondition)
		} else {
			// No explicit join condition found, fall back to a simpler approach
			fromClause += fmt.Sprintf(" %s %s", joinType, tables[i].Name)
		}
	}

	return fromClause
}

// determineJoinType determines the type of join to use based on the prompt
func (qb *QueryBuilder) determineJoinType(prompt string) string {
	// Check based on language
	if qb.language == "pt" {
		if strings.Contains(prompt, "left join") || strings.Contains(prompt, "junção esquerda") {
			return "LEFT JOIN"
		} else if strings.Contains(prompt, "right join") || strings.Contains(prompt, "junção direita") {
			return "RIGHT JOIN"
		} else if strings.Contains(prompt, "full join") || strings.Contains(prompt, "junção completa") {
			return "FULL JOIN"
		} else if strings.Contains(prompt, "inner join") || strings.Contains(prompt, "junção interna") {
			return "INNER JOIN"
		}
	} else {
		if strings.Contains(prompt, "left join") {
			return "LEFT JOIN"
		} else if strings.Contains(prompt, "right join") {
			return "RIGHT JOIN"
		} else if strings.Contains(prompt, "full join") {
			return "FULL JOIN"
		} else if strings.Contains(prompt, "inner join") {
			return "INNER JOIN"
		}
	}

	// Default join type
	return "JOIN"
}

// buildWhereClause constructs a WHERE clause from conditions
func (qb *QueryBuilder) buildWhereClause(conditions []models.Condition) string {
	if len(conditions) == 0 {
		return ""
	}

	var clauses []string

	for i, cond := range conditions {
		if cond.IsComplex {
			// This is a complex expression that's already formatted
			clauses = append(clauses, cond.ComplexExpr)
			continue
		}

		// Skip first conjunction
		conjunction := ""
		if i > 0 {
			conjunction = cond.Conjunction + " "
		}

		// Format the column reference
		var columnRef string
		if cond.TableName != "" {
			columnRef = cond.TableName + "." + cond.ColumnName
		} else {
			columnRef = cond.ColumnName
		}

		// Special cases for NULL operators
		if cond.Operator == "IS NULL" || cond.Operator == "IS NOT NULL" {
			clauses = append(clauses, conjunction+columnRef+" "+cond.Operator)
		} else {
			clauses = append(clauses, conjunction+columnRef+" "+cond.Operator+" "+cond.Value)
		}
	}

	return strings.Join(clauses, " ")
}

// generateOrderByClause creates an ORDER BY clause
func (qb *QueryBuilder) generateOrderByClause(tables []models.TableInfo, prompt string) string {
	if len(tables) == 0 {
		return ""
	}

	// Check for explicit order by in the prompt based on language
	var orderByRegex *regexp.Regexp
	if qb.language == "pt" {
		orderByRegex = regexp.MustCompile(`(?i)(?:ordenar|classificar)\s+por\s+(\w+)(?:\s+(asc|desc|crescente|decrescente))?`)
	} else {
		orderByRegex = regexp.MustCompile(`(?i)order\s+by\s+(\w+)(\s+(asc|desc))?`)
	}

	if matches := orderByRegex.FindStringSubmatch(prompt); len(matches) > 1 {
		column := matches[1]
		direction := "ASC"
		if len(matches) > 3 && matches[3] != "" {
			// Handle Portuguese directions
			if matches[3] == "decrescente" || matches[3] == "desc" {
				direction = "DESC"
			} else {
				direction = "ASC"
			}
		}

		// Check if this is a valid column
		isValidColumn := false
		for _, table := range qb.dbTables {
			if table.Name == tables[0].Name {
				for _, col := range table.Columns {
					if strings.EqualFold(col.Name, column) {
						isValidColumn = true
						break
					}
				}
			}
		}

		if isValidColumn {
			return column + " " + direction
		}
	}

	// Check for ordering keywords based on language
	var direction string
	if qb.language == "pt" {
		if util.ContainsAny(prompt, []string{"desc", "decrescente", "mais recente", "último", "mais alto", "maior"}) {
			direction = "DESC"
		} else {
			direction = "ASC"
		}
	} else {
		if util.ContainsAny(prompt, []string{"desc", "descending", "newest", "latest", "highest", "largest"}) {
			direction = "DESC"
		} else {
			direction = "ASC"
		}
	}

	// Look for suitable columns to order by
	for _, table := range qb.dbTables {
		if table.Name == tables[0].Name {
			// First check for date columns
			for _, col := range table.Columns {
				colLower := strings.ToLower(col.Name)
				if util.ContainsAny(colLower, []string{"created", "updated", "date", "time", "criado", "atualizado", "data", "hora"}) {
					return col.Name + " " + direction
				}
			}

			// Then check for ID or primary key
			for _, col := range table.Columns {
				if col.IsPrimary || strings.EqualFold(col.Name, "id") {
					return col.Name + " " + direction
				}
			}

			// Then check for name or title
			for _, col := range table.Columns {
				colLower := strings.ToLower(col.Name)
				if util.ContainsAny(colLower, []string{"name", "title", "nome", "titulo"}) {
					return col.Name + " " + direction
				}
			}
		}
	}

	return ""
}

// generateGroupByClause creates a GROUP BY clause
func (qb *QueryBuilder) generateGroupByClause(tables []models.TableInfo, prompt string) string {
	if len(tables) == 0 {
		return ""
	}

	// Check for explicit group by in the prompt based on language
	var groupByRegex *regexp.Regexp
	if qb.language == "pt" {
		groupByRegex = regexp.MustCompile(`(?i)(?:agrupar|agrupar por)\s+(\w+)`)
	} else {
		groupByRegex = regexp.MustCompile(`(?i)group\s+by\s+(\w+)`)
	}

	if matches := groupByRegex.FindStringSubmatch(prompt); len(matches) > 1 {
		column := matches[1]

		// Check if this is a valid column
		isValidColumn := false
		for _, table := range qb.dbTables {
			if table.Name == tables[0].Name {
				for _, col := range table.Columns {
					if strings.EqualFold(col.Name, column) {
						isValidColumn = true
						break
					}
				}
			}
		}

		if isValidColumn {
			return column
		}
	}

	// Look for categorical columns that would make sense for grouping
	for _, table := range qb.dbTables {
		if table.Name == tables[0].Name {
			// Look for categorical columns based on language
			var categoryKeys []string
			if qb.language == "pt" {
				categoryKeys = []string{"tipo", "categoria", "status", "grupo", "classe", "nivel"}
			} else {
				categoryKeys = []string{"type", "category", "status", "group", "class", "level"}
			}

			for _, col := range table.Columns {
				colLower := strings.ToLower(col.Name)
				if util.ContainsAny(colLower, categoryKeys) {
					return col.Name
				}
			}

			// Check for date parts that might be used for grouping
			for _, col := range table.Columns {
				colLower := strings.ToLower(col.Name)
				if util.ContainsAny(colLower, []string{"date", "time", "created", "updated", "data", "hora", "criado", "atualizado"}) {
					// If the prompt mentions specific date parts, use those
					if qb.language == "pt" {
						if util.ContainsAny(prompt, []string{"por ano", "anual", "por ano"}) {
							return fmt.Sprintf("EXTRACT(YEAR FROM %s)", col.Name)
						} else if util.ContainsAny(prompt, []string{"por mês", "mensal", "por mes"}) {
							return fmt.Sprintf("EXTRACT(MONTH FROM %s)", col.Name)
						} else if util.ContainsAny(prompt, []string{"por dia", "diário", "por dia"}) {
							return fmt.Sprintf("EXTRACT(DAY FROM %s)", col.Name)
						}
					} else {
						if util.ContainsAny(prompt, []string{"by year", "yearly", "per year"}) {
							return fmt.Sprintf("EXTRACT(YEAR FROM %s)", col.Name)
						} else if util.ContainsAny(prompt, []string{"by month", "monthly", "per month"}) {
							return fmt.Sprintf("EXTRACT(MONTH FROM %s)", col.Name)
						} else if util.ContainsAny(prompt, []string{"by day", "daily", "per day"}) {
							return fmt.Sprintf("EXTRACT(DAY FROM %s)", col.Name)
						}
					}
					return col.Name
				}
			}
		}
	}

	return ""
}

// generateLimitClause creates a LIMIT clause based on prompt analysis
func (qb *QueryBuilder) generateLimitClause(prompt string) string {
	// Check for explicit limit patterns in different languages
	limitRegexes := []*regexp.Regexp{
		regexp.MustCompile(`(?i)limit\s+(\d+)`),
		regexp.MustCompile(`(?i)limitar\s+(?:a|para|em)?\s*(\d+)`),
		regexp.MustCompile(`(?i)(?:mostrar|exibir|retornar|trazer)\s+(?:apenas\s+)?(\d+)`),
		regexp.MustCompile(`(?i)(?:top|primeiros|primeiras|apenas)\s+(\d+)`),
		regexp.MustCompile(`(?i)(\d+)\s+(?:registros|linhas|resultados|itens)`),
	}

	for _, regex := range limitRegexes {
		if matches := regex.FindStringSubmatch(prompt); len(matches) > 1 {
			limitStr := matches[1]
			if limit, err := strconv.Atoi(limitStr); err == nil {
				// Make sure the limit is reasonable (between 1 and 1000)
				if limit < 1 {
					limit = 1
				} else if limit > 1000 {
					limit = 1000
				}
				return fmt.Sprintf("LIMIT %d", limit)
			}
		}
	}

	// Check for keywords indicating small result sets
	smallLimitKeywords := []string{"pequeno", "poucos", "exemplo", "amostra", "small", "few", "example", "sample"}
	for _, keyword := range smallLimitKeywords {
		if strings.Contains(strings.ToLower(prompt), keyword) {
			return "LIMIT 10"
		}
	}

	// Check for keywords indicating larger result sets but still limited
	mediumLimitKeywords := []string{"vários", "diversos", "alguns", "many", "several", "some"}
	for _, keyword := range mediumLimitKeywords {
		if strings.Contains(strings.ToLower(prompt), keyword) {
			return "LIMIT 50"
		}
	}

	// Check for keywords indicating very large result sets
	largeLimitKeywords := []string{"todos", "todas", "completo", "completa", "all", "complete", "full"}
	for _, keyword := range largeLimitKeywords {
		if strings.Contains(strings.ToLower(prompt), keyword) {
			return "LIMIT 1000"
		}
	}

	// Default to a reasonable limit if no specific indication
	return "LIMIT 100"
}

// findRelatedTable finds a table related to the given table
func (qb *QueryBuilder) findRelatedTable(tableName string) string {
	// First check foreign keys
	for _, table := range qb.dbTables {
		if table.Name == tableName {
			for _, fk := range table.ForeignKeys {
				return fk.ReferencedTable
			}
		}
	}

	// Check if other tables reference this one
	for _, table := range qb.dbTables {
		if table.Name != tableName {
			for _, fk := range table.ForeignKeys {
				if fk.ReferencedTable == tableName {
					return table.Name
				}
			}
		}
	}

	// Look for naming patterns
	singularTable := strings.TrimSuffix(tableName, "s")
	for _, table := range qb.dbTables {
		if table.Name != tableName {
			// Check for columns that look like foreign keys
			for _, col := range table.Columns {
				if col.Name == tableName+"_id" || col.Name == singularTable+"_id" {
					return table.Name
				}
			}
		}
	}

	// Look for common relationships in data models
	relationshipPatterns := map[string][]string{
		"user":     {"profile", "address", "order", "comment", "post"},
		"usuario":  {"perfil", "endereco", "pedido", "comentario", "postagem"},
		"product":  {"category", "order_item", "review", "inventory"},
		"produto":  {"categoria", "item_pedido", "avaliacao", "inventario"},
		"order":    {"order_item", "payment", "shipment", "customer"},
		"pedido":   {"item_pedido", "pagamento", "envio", "cliente"},
		"post":     {"comment", "category", "tag", "author"},
		"postagem": {"comentario", "categoria", "tag", "autor"},
		"customer": {"address", "payment", "order"},
		"cliente":  {"endereco", "pagamento", "pedido"},
	}

	for pattern, relatedTables := range relationshipPatterns {
		if strings.Contains(tableName, pattern) {
			// Check if any of the related tables exist in our schema
			for _, relTable := range relatedTables {
				for _, table := range qb.dbTables {
					if table.Name == relTable ||
						table.Name == relTable+"s" {
						return table.Name
					}
				}
			}
		}
	}

	return ""
}

// findGroupByColumns identifies columns suitable for GROUP BY
func (qb *QueryBuilder) findGroupByColumns(tables []models.TableInfo, prompt string) []string {
	if len(tables) == 0 {
		return nil
	}

	// Look for explicit group by column in the prompt
	var groupByRegex *regexp.Regexp
	if qb.language == "pt" {
		groupByRegex = regexp.MustCompile(`(?i)(?:agrupar|agrupar por)\s+(\w+)`)
	} else {
		groupByRegex = regexp.MustCompile(`(?i)group\s+by\s+(\w+)`)
	}

	if matches := groupByRegex.FindStringSubmatch(prompt); len(matches) > 1 {
		columnName := matches[1]

		// Verify this is a valid column
		for _, table := range qb.dbTables {
			if table.Name == tables[0].Name {
				for _, col := range table.Columns {
					if strings.EqualFold(col.Name, columnName) {
						return []string{columnName}
					}
				}
			}
		}
	}

	// Check for date grouping patterns based on language
	var dateGroupPatterns map[string]string
	if qb.language == "pt" {
		dateGroupPatterns = map[string]string{
			"(?i)por\\s+ano":       "EXTRACT(YEAR FROM %s)",
			"(?i)por\\s+mês":       "EXTRACT(MONTH FROM %s)",
			"(?i)por\\s+dia":       "EXTRACT(DAY FROM %s)",
			"(?i)por\\s+trimestre": "EXTRACT(QUARTER FROM %s)",
			"(?i)por\\s+semana":    "EXTRACT(WEEK FROM %s)",
			"(?i)mensal\\s+":       "EXTRACT(MONTH FROM %s)",
			"(?i)anual\\s+":        "EXTRACT(YEAR FROM %s)",
			"(?i)diário\\s+":       "EXTRACT(DAY FROM %s)",
			"(?i)trimestral\\s+":   "EXTRACT(QUARTER FROM %s)",
			"(?i)semanal\\s+":      "EXTRACT(WEEK FROM %s)",
		}
	} else {
		dateGroupPatterns = map[string]string{
			"(?i)by\\s+year":    "EXTRACT(YEAR FROM %s)",
			"(?i)by\\s+month":   "EXTRACT(MONTH FROM %s)",
			"(?i)by\\s+day":     "EXTRACT(DAY FROM %s)",
			"(?i)by\\s+quarter": "EXTRACT(QUARTER FROM %s)",
			"(?i)by\\s+week":    "EXTRACT(WEEK FROM %s)",
			"(?i)monthly\\s+":   "EXTRACT(MONTH FROM %s)",
			"(?i)yearly\\s+":    "EXTRACT(YEAR FROM %s)",
			"(?i)daily\\s+":     "EXTRACT(DAY FROM %s)",
			"(?i)quarterly\\s+": "EXTRACT(QUARTER FROM %s)",
			"(?i)weekly\\s+":    "EXTRACT(WEEK FROM %s)",
		}
	}

	for pattern, format := range dateGroupPatterns {
		if regexp.MustCompile(pattern).MatchString(prompt) {
			// Look for date columns
			for _, table := range qb.dbTables {
				if table.Name == tables[0].Name {
					for _, col := range table.Columns {
						if strings.Contains(strings.ToLower(col.Type), "date") ||
							strings.Contains(strings.ToLower(col.Type), "time") {
							return []string{fmt.Sprintf(format, col.Name)}
						}
					}
				}
			}
		}
	}

	// Look for categorical columns
	var categoricalColumns []string
	for _, table := range qb.dbTables {
		if table.Name == tables[0].Name {
			// Define category keys based on language
			var categoryKeys []string
			if qb.language == "pt" {
				categoryKeys = []string{"tipo", "categoria", "status", "grupo", "classe", "nivel", "estado", "pais", "regiao"}
			} else {
				categoryKeys = []string{"type", "category", "status", "group", "class", "level", "state", "country", "region"}
			}

			for _, col := range table.Columns {
				colLower := strings.ToLower(col.Name)
				// Text columns with names suggesting categories
				if (strings.Contains(col.Type, "char") || strings.Contains(col.Type, "text")) &&
					util.ContainsAny(colLower, categoryKeys) {
					categoricalColumns = append(categoricalColumns, col.Name)
				}
			}
		}
	}

	if len(categoricalColumns) > 0 {
		return categoricalColumns
	}

	// If no suitable categorical columns, try foreign key references
	for _, table := range qb.dbTables {
		if table.Name == tables[0].Name {
			for _, col := range table.Columns {
				if strings.HasSuffix(col.Name, "_id") && !col.IsPrimary {
					return []string{col.Name}
				}
			}
		}
	}

	return nil
}

// getAggregationFunctionFormat returns the format string for an aggregation function
func (qb *QueryBuilder) getAggregationFunctionFormat(function string) string {
	// Get the matching format from the maps based on language
	function = strings.ToLower(function)

	if qb.language == "pt" {
		switch function {
		case "avg", "média":
			return "AVG(%s)"
		case "sum", "soma":
			return "SUM(%s)"
		case "max", "máximo", "maior":
			return "MAX(%s)"
		case "min", "mínimo", "menor":
			return "MIN(%s)"
		case "count", "contagem", "contar":
			return "COUNT(%s)"
		case "count distinct", "contagem distinta", "contagem única":
			return "COUNT(DISTINCT %s)"
		case "stddev", "desvio padrão":
			return "STDDEV(%s)"
		case "variance", "variância":
			return "VARIANCE(%s)"
		case "median", "mediana":
			return "PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY %s)"
		}
	} else {
		switch function {
		case "avg", "average", "mean":
			return "AVG(%s)"
		case "sum", "total":
			return "SUM(%s)"
		case "max", "maximum", "highest":
			return "MAX(%s)"
		case "min", "minimum", "lowest":
			return "MIN(%s)"
		case "count":
			return "COUNT(%s)"
		case "count distinct", "unique count":
			return "COUNT(DISTINCT %s)"
		case "stddev", "standard deviation":
			return "STDDEV(%s)"
		case "variance":
			return "VARIANCE(%s)"
		case "median":
			return "PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY %s)"
		}
	}

	// If no match found, return empty string to use the original function name
	return ""
}
