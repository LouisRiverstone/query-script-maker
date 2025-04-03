package operations

import (
	"fmt"
	"regexp"
	"strings"

	"sql_script_maker/sqlai/models"
	"sql_script_maker/sqlai/util"
)

// QueryBuilder builds SQL queries based on entity information and operation type
type QueryBuilder struct {
	language string
	dbTables []models.TableForAI
}

// NewQueryBuilder creates a new query builder
func NewQueryBuilder(language string, dbTables []models.TableForAI) *QueryBuilder {
	return &QueryBuilder{
		language: language,
		dbTables: dbTables,
	}
}

// BuildSelectQuery creates a SELECT query
func (qb *QueryBuilder) BuildSelectQuery(tables []models.TableInfo, columns []models.ColumnInfo,
	conditions []models.Condition, subOperations []string, prompt string) string {

	// Build column list
	var columnList []string
	for _, col := range columns {
		if col.Name == "*" {
			columnList = append(columnList, "*")
		} else if col.Function != "" {
			// Get the matching format for the aggregation function
			format := qb.getAggregationFunctionFormat(col.Function)
			if format != "" {
				columnList = append(columnList, fmt.Sprintf(format, col.TableName+"."+col.Name))
			} else {
				columnList = append(columnList, fmt.Sprintf("%s(%s.%s)", col.Function, col.TableName, col.Name))
			}
		} else if col.TableName != "" {
			columnList = append(columnList, col.TableName+"."+col.Name)
		} else {
			columnList = append(columnList, col.Name)
		}
	}

	// If no columns, use *
	if len(columnList) == 0 {
		columnList = append(columnList, "*")
	}

	// Build table list with proper join syntax if needed
	var fromClause string
	if len(tables) == 1 {
		fromClause = tables[0].Name
	} else if len(tables) > 1 {
		// This is a join query, build it properly
		fromClause = qb.buildJoinClause(tables, prompt)
	}

	// Build WHERE clause
	whereClause := qb.buildWhereClause(conditions)

	// Handle sub-operations (additional clauses)
	var additionalClauses []string

	for _, subOp := range subOperations {
		switch subOp {
		case "order":
			orderClause := qb.generateOrderByClause(tables, prompt)
			if orderClause != "" {
				additionalClauses = append(additionalClauses, "ORDER BY "+orderClause)
			}

		case "group":
			groupClause := qb.generateGroupByClause(tables, prompt)
			if groupClause != "" {
				additionalClauses = append(additionalClauses, "GROUP BY "+groupClause)
			}

		case "limit":
			limitClause := qb.generateLimitClause(prompt)
			if limitClause != "" {
				additionalClauses = append(additionalClauses, limitClause)
			}
		}
	}

	// Build the final query
	query := "SELECT " + strings.Join(columnList, ", ") + " FROM " + fromClause

	if whereClause != "" {
		query += " WHERE " + whereClause
	}

	for _, clause := range additionalClauses {
		query += " " + clause
	}

	return query
}

// BuildCountQuery creates a COUNT query
func (qb *QueryBuilder) BuildCountQuery(tables []models.TableInfo, columns []models.ColumnInfo, conditions []models.Condition, prompt string) string {
	// Determine what to count
	countColumn := "*"
	if len(columns) > 0 && columns[0].Name != "*" {
		if columns[0].TableName != "" {
			countColumn = columns[0].TableName + "." + columns[0].Name
		} else {
			countColumn = columns[0].Name
		}
	}

	// Check if we need DISTINCT
	if util.ContainsAny(prompt, []string{"unique", "distinct", "different", "único", "distintos", "diferentes"}) {
		countColumn = "DISTINCT " + countColumn
	}

	// Build table reference
	var fromClause string
	if len(tables) == 1 {
		fromClause = tables[0].Name
	} else if len(tables) > 1 {
		fromClause = qb.buildJoinClause(tables, prompt)
	} else {
		return "SELECT COUNT(*) FROM users" // Fallback to a common table
	}

	// Build the query
	query := fmt.Sprintf("SELECT COUNT(%s) FROM %s", countColumn, fromClause)

	// Add WHERE clause if needed
	whereClause := qb.buildWhereClause(conditions)
	if whereClause != "" {
		query += " WHERE " + whereClause
	}

	return query
}

// BuildJoinQuery creates a JOIN query
func (qb *QueryBuilder) BuildJoinQuery(tables []models.TableInfo, columns []models.ColumnInfo,
	conditions []models.Condition, subOperations []string, prompt string) string {

	// Handle case with insufficient tables
	if len(tables) < 2 {
		// Try to find a related table
		if len(tables) == 1 {
			relatedTable := qb.findRelatedTable(tables[0].Name)
			if relatedTable != "" {
				tables = append(tables, models.TableInfo{
					Name:       relatedTable,
					Confidence: 0.7,
				})
			}
		}

		// Still not enough tables? Default to basic select
		if len(tables) < 2 {
			return qb.BuildSelectQuery(tables, columns, conditions, subOperations, prompt)
		}
	}

	// Determine join type
	joinType := qb.determineJoinType(prompt)

	// Build column list
	var columnList []string
	if len(columns) == 0 {
		// Use primary keys and name columns from both tables
		for _, table := range tables {
			for _, dbTable := range qb.dbTables {
				if dbTable.Name == table.Name {
					for _, col := range dbTable.Columns {
						if col.IsPrimary || strings.ToLower(col.Name) == "id" ||
							util.ContainsAny(strings.ToLower(col.Name), []string{
								"name", "title", "description",
								"nome", "titulo", "descrição",
							}) {
							columnList = append(columnList, table.Name+"."+col.Name)
						}
					}
				}
			}
		}
	} else {
		for _, col := range columns {
			if col.Name == "*" {
				columnList = append(columnList, "*")
			} else if col.Function != "" {
				columnList = append(columnList, fmt.Sprintf("%s(%s.%s)",
					col.Function, col.TableName, col.Name))
			} else if col.TableName != "" {
				columnList = append(columnList, col.TableName+"."+col.Name)
			} else {
				columnList = append(columnList, col.Name)
			}
		}
	}

	// If still no columns, use *
	if len(columnList) == 0 {
		columnList = append(columnList, "*")
	}

	// Find join condition
	joinCondition := FindJoinCondition(tables[0].Name, tables[1].Name, qb.dbTables)

	// Build the query
	query := fmt.Sprintf("SELECT %s FROM %s %s %s ON %s",
		strings.Join(columnList, ", "),
		tables[0].Name,
		joinType,
		tables[1].Name,
		joinCondition)

	// Add WHERE clause if needed
	whereClause := qb.buildWhereClause(conditions)
	if whereClause != "" {
		query += " WHERE " + whereClause
	}

	// Add additional clauses from sub-operations
	for _, subOp := range subOperations {
		switch subOp {
		case "order":
			orderClause := qb.generateOrderByClause(tables, prompt)
			if orderClause != "" {
				query += " ORDER BY " + orderClause
			}

		case "group":
			groupClause := qb.generateGroupByClause(tables, prompt)
			if groupClause != "" {
				query += " GROUP BY " + groupClause
			}

		case "limit":
			limitClause := qb.generateLimitClause(prompt)
			if limitClause != "" {
				query += " " + limitClause
			}
		}
	}

	return query
}

// BuildGroupByQuery creates a GROUP BY query
func (qb *QueryBuilder) BuildGroupByQuery(tables []models.TableInfo, columns []models.ColumnInfo, conditions []models.Condition, prompt string) string {
	if len(tables) == 0 {
		return "SELECT COUNT(*) FROM users GROUP BY status" // Fallback
	}

	// Determine grouping column(s)
	groupColumns := qb.findGroupByColumns(tables, prompt)
	if len(groupColumns) == 0 {
		// Fall back to basic select if no group columns found
		return qb.BuildSelectQuery(tables, columns, conditions, nil, prompt)
	}

	// Prepare SELECT columns: group columns + aggregates
	var selectColumns []string

	// Add group columns to SELECT
	for _, col := range groupColumns {
		selectColumns = append(selectColumns, col)
	}

	// Look for aggregation functions in the prompt
	var aggregatedCols []string

	// Get aggregation functions based on language
	var aggregationFunctions map[string]string
	if qb.language == "pt" {
		aggregationFunctions = map[string]string{
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
	} else {
		aggregationFunctions = map[string]string{
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
	}

	for funcName, funcFormat := range aggregationFunctions {
		if strings.Contains(prompt, funcName) {
			// Find numeric columns to aggregate
			for _, table := range qb.dbTables {
				if table.Name == tables[0].Name {
					for _, col := range table.Columns {
						// Don't aggregate the grouping columns
						if util.ContainsAny(col.Name, groupColumns) {
							continue
						}

						if strings.Contains(col.Type, "int") ||
							strings.Contains(col.Type, "float") ||
							strings.Contains(col.Type, "decimal") ||
							strings.Contains(col.Type, "double") {
							aggregatedCols = append(aggregatedCols,
								fmt.Sprintf(funcFormat, col.Name))
							break
						}
					}
				}
			}
		}
	}

	// If no specific aggregations found, add COUNT(*)
	if len(aggregatedCols) == 0 {
		aggregatedCols = append(aggregatedCols, "COUNT(*) as count")
	}

	// Combine all columns for SELECT
	selectColumns = append(selectColumns, aggregatedCols...)

	// Build the query
	query := fmt.Sprintf("SELECT %s FROM %s",
		strings.Join(selectColumns, ", "),
		tables[0].Name)

	// Add WHERE clause if needed
	whereClause := qb.buildWhereClause(conditions)
	if whereClause != "" {
		query += " WHERE " + whereClause
	}

	// Add GROUP BY clause
	query += " GROUP BY " + strings.Join(groupColumns, ", ")

	// Check for HAVING clause based on language
	if qb.language == "pt" {
		if util.ContainsAny(prompt, []string{"tendo", "com", "onde a contagem", "onde o total"}) {
			// Try to extract having conditions
			havingRegex := regexp.MustCompile(`(?:tendo|com)\s+(\w+)\s+(>|<|=|>=|<=)\s+(\d+)`)
			if matches := havingRegex.FindStringSubmatch(prompt); len(matches) > 3 {
				query += fmt.Sprintf(" HAVING %s %s %s", matches[1], matches[2], matches[3])
			} else {
				// Default HAVING for count queries
				query += " HAVING COUNT(*) > 1"
			}
		}
	} else {
		if util.ContainsAny(prompt, []string{"having", "where the count", "where the total"}) {
			// Try to extract having conditions
			havingRegex := regexp.MustCompile(`having\s+(\w+)\s+(>|<|=|>=|<=)\s+(\d+)`)
			if matches := havingRegex.FindStringSubmatch(prompt); len(matches) > 3 {
				query += fmt.Sprintf(" HAVING %s %s %s", matches[1], matches[2], matches[3])
			} else {
				// Default HAVING for count queries
				query += " HAVING COUNT(*) > 1"
			}
		}
	}

	// Add ORDER BY if needed based on language
	if qb.language == "pt" {
		if util.ContainsAny(prompt, []string{"ordenar", "classificar"}) {
			// Try to determine if we should order by an aggregate
			if len(aggregatedCols) > 0 &&
				(util.ContainsAny(prompt, []string{"maior", "mais", "decrescente"}) ||
					strings.Contains(prompt, "desc")) {
				query += " ORDER BY " + strings.TrimLeft(aggregatedCols[0], "as ") + " DESC"
			} else if len(aggregatedCols) > 0 {
				query += " ORDER BY " + strings.TrimLeft(aggregatedCols[0], "as ") + " ASC"
			}
		}
	} else {
		if util.ContainsAny(prompt, []string{"order", "sort"}) {
			// Try to determine if we should order by an aggregate
			if len(aggregatedCols) > 0 &&
				(util.ContainsAny(prompt, []string{"highest", "most", "largest", "descending"}) ||
					strings.Contains(prompt, "desc")) {
				query += " ORDER BY " + strings.TrimLeft(aggregatedCols[0], "as ") + " DESC"
			} else if len(aggregatedCols) > 0 {
				query += " ORDER BY " + strings.TrimLeft(aggregatedCols[0], "as ") + " ASC"
			}
		}
	}

	return query
}

// BuildOrderByQuery creates an ORDER BY query
func (qb *QueryBuilder) BuildOrderByQuery(tables []models.TableInfo, columns []models.ColumnInfo, conditions []models.Condition, prompt string) string {
	// This is essentially a SELECT query with ORDER BY
	orderByClause := qb.generateOrderByClause(tables, prompt)

	// Build the main SELECT query
	query := qb.BuildSelectQuery(tables, columns, conditions, nil, prompt)

	// Add ORDER BY if we found a suitable clause
	if orderByClause != "" && !strings.Contains(query, "ORDER BY") {
		query += " ORDER BY " + orderByClause
	}

	return query
}

// BuildLimitQuery creates a SELECT query with LIMIT
func (qb *QueryBuilder) BuildLimitQuery(tables []models.TableInfo, columns []models.ColumnInfo, conditions []models.Condition, prompt string) string {
	// This is essentially a SELECT query with LIMIT
	limitClause := qb.generateLimitClause(prompt)

	// Build the main SELECT query
	query := qb.BuildSelectQuery(tables, columns, conditions, nil, prompt)

	// Add LIMIT if we found a suitable clause
	if limitClause != "" && !strings.Contains(query, "LIMIT") {
		query += " " + limitClause
	}

	return query
}

// BuildDistinctQuery creates a DISTINCT query
func (qb *QueryBuilder) BuildDistinctQuery(tables []models.TableInfo, columns []models.ColumnInfo, conditions []models.Condition, prompt string) string {
	if len(tables) == 0 {
		return "SELECT DISTINCT * FROM users" // Fallback
	}

	// Ensure we have columns
	if len(columns) == 0 {
		columns = SelectRelevantColumns(tables, "distinct", prompt, qb.dbTables)
	}

	// Build column list
	var columnList []string
	for _, col := range columns {
		if col.Name == "*" {
			columnList = append(columnList, "*")
		} else if col.TableName != "" {
			columnList = append(columnList, col.TableName+"."+col.Name)
		} else {
			columnList = append(columnList, col.Name)
		}
	}

	// Build the query
	query := fmt.Sprintf("SELECT DISTINCT %s FROM %s",
		strings.Join(columnList, ", "),
		tables[0].Name)

	// Add WHERE clause if needed
	whereClause := qb.buildWhereClause(conditions)
	if whereClause != "" {
		query += " WHERE " + whereClause
	}

	// Add ORDER BY if present in prompt
	if util.ContainsAny(prompt, []string{"order", "sort", "ordenar", "classificar"}) {
		orderClause := qb.generateOrderByClause(tables, prompt)
		if orderClause != "" {
			query += " ORDER BY " + orderClause
		}
	}

	return query
}

// BuildInsertQuery creates an INSERT query
func (qb *QueryBuilder) BuildInsertQuery(tableName string, prompt string) string {
	// Find the table structure
	var tableColumns []models.ColumnForAI
	for _, table := range qb.dbTables {
		if table.Name == tableName {
			tableColumns = table.Columns
			break
		}
	}

	if len(tableColumns) == 0 {
		return fmt.Sprintf("INSERT INTO %s (column1, column2) VALUES (value1, value2)", tableName)
	}

	// Filter out auto-increment columns
	var columns []string
	for _, col := range tableColumns {
		// Skip auto-increment columns
		if strings.Contains(strings.ToLower(col.Extra), "auto_increment") {
			continue
		}
		columns = append(columns, col.Name)
	}

	// Create placeholders
	var placeholders []string
	for range columns {
		placeholders = append(placeholders, "?")
	}

	return fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)",
		tableName,
		strings.Join(columns, ", "),
		strings.Join(placeholders, ", "))
}

// BuildUpdateQuery creates an UPDATE query
func (qb *QueryBuilder) BuildUpdateQuery(tableName string, conditions []models.Condition, prompt string) string {
	// Find the table structure
	var tableColumns []models.ColumnForAI
	for _, table := range qb.dbTables {
		if table.Name == tableName {
			tableColumns = table.Columns
			break
		}
	}

	if len(tableColumns) == 0 {
		return fmt.Sprintf("UPDATE %s SET column1 = value1 WHERE id = ?", tableName)
	}

	// Create SET clause with placeholders
	var setClause []string
	for _, col := range tableColumns {
		// Skip primary key and auto-increment columns
		if col.IsPrimary || strings.Contains(strings.ToLower(col.Extra), "auto_increment") {
			continue
		}

		// Add up to 3 columns only for better usability
		if len(setClause) < 3 {
			setClause = append(setClause, col.Name+" = ?")
		}
	}

	// Build WHERE clause
	var whereClause string
	if len(conditions) > 0 {
		whereClause = qb.buildWhereClause(conditions)
	} else {
		// Default to primary key condition
		for _, col := range tableColumns {
			if col.IsPrimary {
				whereClause = col.Name + " = ?"
				break
			}
		}

		// If no primary key found, use id
		if whereClause == "" {
			whereClause = "id = ?"
		}
	}

	return fmt.Sprintf("UPDATE %s SET %s WHERE %s",
		tableName,
		strings.Join(setClause, ", "),
		whereClause)
}

// BuildDeleteQuery creates a DELETE query
func (qb *QueryBuilder) BuildDeleteQuery(tableName string, conditions []models.Condition, prompt string) string {
	// Build WHERE clause
	var whereClause string
	if len(conditions) > 0 {
		whereClause = qb.buildWhereClause(conditions)
	} else {
		// Default to id condition for safety
		whereClause = "id = ?"
	}

	return fmt.Sprintf("DELETE FROM %s WHERE %s", tableName, whereClause)
}
