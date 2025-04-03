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

	// Check if we have tables
	if len(tables) == 0 {
		// Try to find a relevant table from context
		mostRelevantTable := FindMostRelevantTable(prompt, qb.dbTables)
		if mostRelevantTable != "" {
			tables = append(tables, models.TableInfo{
				Name:       mostRelevantTable,
				Confidence: 0.6,
			})
		} else if len(qb.dbTables) > 0 {
			// If we can't determine a table, return a comment
			// informing the user about the problem
			if qb.language == "pt" {
				return "-- Não foi possível determinar qual tabela consultar. Por favor, especifique a tabela em sua consulta."
			} else {
				return "-- Unable to determine which table to query. Please specify the table in your query."
			}
		}
	}

	// If we still don't have tables, we can't build the query
	if len(tables) == 0 {
		if qb.language == "pt" {
			return "-- Não foi possível determinar qual tabela consultar"
		} else {
			return "-- Unable to determine which table to query"
		}
	}

	// Directly extract specific table and column mentions from the prompt
	// This is to handle cases where the user explicitly mentions them
	explicitTableRegex := regexp.MustCompile(`(?i)(?:da|de|from)\s+tabela\s+(\w+)`)
	explicitColRegex := regexp.MustCompile(`(?i)(?:a coluna|as colunas|selecione|select|column[s]?)\s+(\w+)`)

	if matches := explicitTableRegex.FindStringSubmatch(prompt); len(matches) > 1 {
		explicitTable := matches[1]
		// Verify this table exists in our schema
		tableExists := false
		for _, dbTable := range qb.dbTables {
			if strings.EqualFold(dbTable.Name, explicitTable) {
				// Replace the tables array entirely with this explicit mention
				tables = []models.TableInfo{{
					Name:       dbTable.Name, // Use correct case from schema
					Confidence: 1.0,
				}}
				tableExists = true
				break
			}
		}

		// If we found an explicit table mention but it doesn't exist in schema,
		// we should still use it (the user might know better than us)
		if !tableExists && explicitTable != "" {
			tables = []models.TableInfo{{
				Name:       explicitTable,
				Confidence: 0.9,
			}}
		}
	}

	// Check for explicit column mentions too
	if matches := explicitColRegex.FindStringSubmatch(prompt); len(matches) > 1 && len(tables) > 0 {
		explicitCol := matches[1]
		// If we have an explicit column mentioned, override any previously detected columns
		columns = []models.ColumnInfo{{
			Name:       explicitCol,
			TableName:  tables[0].Name,
			Confidence: 1.0,
		}}
	}

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

	// If subOperations is empty but we have a prompt, try to determine what suboperations might be needed
	if len(subOperations) == 0 && prompt != "" {
		promptLower := strings.ToLower(prompt)

		// Check for ordering needs
		if util.ContainsAny(promptLower, []string{"order", "sort", "ordenar", "ordenado", "decrescente", "crescente", "asc", "desc"}) {
			subOperations = append(subOperations, "order")
		}

		// Check for grouping needs
		if util.ContainsAny(promptLower, []string{"group", "aggregate", "agrupar", "agrupado", "média", "average", "count", "sum", "soma", "total"}) {
			subOperations = append(subOperations, "group")
		}

		// Add limit by default
		subOperations = append(subOperations, "limit")
	}

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
			} else {
				// Always add a limit if none is specified
				additionalClauses = append(additionalClauses, "LIMIT 100")
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
		// No default table here - we need to find a table from context
		fromClause = FindMostRelevantTable(prompt, qb.dbTables)
		if fromClause == "" && len(qb.dbTables) > 0 {
			// If still nothing found, use the first available table
			fromClause = qb.dbTables[0].Name
		}
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
	// Check if we have tables
	if len(tables) == 0 {
		// Try to find a relevant table from context
		mostRelevantTable := FindMostRelevantTable(prompt, qb.dbTables)
		if mostRelevantTable != "" {
			tables = append(tables, models.TableInfo{
				Name:       mostRelevantTable,
				Confidence: 0.6,
			})
		} else if len(qb.dbTables) > 0 {
			// If we can't determine a table, return a comment
			if qb.language == "pt" {
				return "-- Não foi possível determinar qual tabela agrupar. Por favor, especifique a tabela em sua consulta."
			} else {
				return "-- Unable to determine which table to group. Please specify the table in your query."
			}
		}
	}

	// If we still don't have tables, we can't build the query
	if len(tables) == 0 {
		if qb.language == "pt" {
			return "-- Não foi possível determinar qual tabela agrupar"
		} else {
			return "-- Unable to determine which table to group"
		}
	}

	// Find groupable columns
	groupByColumns := qb.findGroupByColumns(tables, prompt)
	if len(groupByColumns) == 0 {
		// If no explicit GROUP BY columns, try to infer from the prompt
		for _, table := range qb.dbTables {
			if table.Name == tables[0].Name {
				for _, col := range table.Columns {
					colLower := strings.ToLower(col.Name)
					if util.ContainsAny(colLower, []string{
						"category", "type", "status", "group", "state", "country", "region",
						"categoria", "tipo", "estado", "país", "região", "grupo", "classe",
					}) {
						groupByColumns = append(groupByColumns, col.Name)
						break
					}
				}
			}
		}
	}

	// Still no groupable columns, check date columns
	if len(groupByColumns) == 0 {
		for _, table := range qb.dbTables {
			if table.Name == tables[0].Name {
				for _, col := range table.Columns {
					if strings.Contains(strings.ToLower(col.Type), "date") ||
						strings.Contains(strings.ToLower(col.Type), "time") {
						groupByColumns = append(groupByColumns, col.Name)
						break
					}
				}
			}
		}
	}

	// If we still have no columns to group by, we can't build the query
	if len(groupByColumns) == 0 {
		if qb.language == "pt" {
			return "-- Não foi possível identificar colunas para agrupar em " + tables[0].Name
		} else {
			return "-- Unable to identify columns to group by in " + tables[0].Name
		}
	}

	// Now find aggregatable columns
	var aggregatableColumns []string
	var aggregatedCols []string

	for _, table := range qb.dbTables {
		if table.Name == tables[0].Name {
			for _, col := range table.Columns {
				if strings.Contains(col.Type, "int") ||
					strings.Contains(col.Type, "float") ||
					strings.Contains(col.Type, "decimal") ||
					strings.Contains(col.Type, "double") ||
					strings.Contains(col.Type, "number") {
					// Check if column is already in the group (replaces util.ContainsString)
					isInGroup := false
					for _, groupCol := range groupByColumns {
						if groupCol == col.Name {
							isInGroup = true
							break
						}
					}

					if !isInGroup {
						aggregatableColumns = append(aggregatableColumns, col.Name)
					}
				}
			}
		}
	}

	// Select aggregation function based on prompt
	aggregationFunc := "COUNT"
	if util.ContainsAny(prompt, []string{"average", "mean", "média"}) {
		aggregationFunc = "AVG"
	} else if util.ContainsAny(prompt, []string{"sum", "total", "soma"}) {
		aggregationFunc = "SUM"
	} else if util.ContainsAny(prompt, []string{"maximum", "highest", "max", "máximo", "maior"}) {
		aggregationFunc = "MAX"
	} else if util.ContainsAny(prompt, []string{"minimum", "lowest", "min", "mínimo", "menor"}) {
		aggregationFunc = "MIN"
	}

	// Build the query
	var selectClause string
	if len(aggregatableColumns) > 0 && aggregationFunc != "COUNT" {
		// Use an aggregation other than COUNT with a numeric column
		selectClause = strings.Join(groupByColumns, ", ") + ", " + aggregationFunc + "(" + aggregatableColumns[0] + ") as aggregated_value"
		aggregatedCols = append(aggregatedCols, "aggregated_value")
	} else {
		// Default to COUNT(*) aggregation
		selectClause = strings.Join(groupByColumns, ", ") + ", COUNT(*) as count"
		aggregatedCols = append(aggregatedCols, "count")
	}

	// Build the query with GROUP BY
	query := "SELECT " + selectClause + " FROM " + tables[0].Name

	// Add WHERE clause if needed
	whereClause := qb.buildWhereClause(conditions)
	if whereClause != "" {
		query += " WHERE " + whereClause
	}

	// Add GROUP BY clause
	query += " GROUP BY " + strings.Join(groupByColumns, ", ")

	// Add HAVING clause if requested
	if util.ContainsAny(prompt, []string{"having", "where the", "tendo", "onde o"}) {
		havingThreshold := "1" // Default threshold

		// Try to extract a threshold value from the prompt
		thresholdRegex := regexp.MustCompile(`(?i)(?:than|above|over|greater than|more than|acima de|maior que|mais que)\s+(\d+)`)
		if matches := thresholdRegex.FindStringSubmatch(prompt); len(matches) > 1 {
			havingThreshold = matches[1]
		}

		// Determine the appropriate aggregation function for HAVING
		if util.ContainsAny(prompt, []string{"at least", "minimum", "pelo menos", "mínimo"}) {
			query += " HAVING " + aggregatedCols[0] + " >= " + havingThreshold
		} else if util.ContainsAny(prompt, []string{"at most", "maximum", "no máximo", "máximo"}) {
			query += " HAVING " + aggregatedCols[0] + " <= " + havingThreshold
		} else {
			query += " HAVING " + aggregatedCols[0] + " > " + havingThreshold
		}
	}

	// Check if we should order by the aggregated value
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
	// Find table if none provided
	if len(tables) == 0 {
		mostRelevantTable := FindMostRelevantTable(prompt, qb.dbTables)
		if mostRelevantTable != "" {
			tables = append(tables, models.TableInfo{
				Name:       mostRelevantTable,
				Confidence: 0.7,
			})
		} else if len(qb.dbTables) > 0 {
			tables = append(tables, models.TableInfo{
				Name:       qb.dbTables[0].Name,
				Confidence: 0.5,
			})
		}
	}

	// If we still don't have tables, we can't build a query
	if len(tables) == 0 {
		return "-- Unable to determine which table to query"
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
		if qb.language == "pt" {
			return fmt.Sprintf("-- Não foi possível encontrar a estrutura da tabela %s para INSERT", tableName)
		} else {
			return fmt.Sprintf("-- Unable to find the structure of table %s for INSERT", tableName)
		}
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

	// If we don't have columns, we can't create the INSERT
	if len(columns) == 0 {
		if qb.language == "pt" {
			return fmt.Sprintf("-- A tabela %s não possui colunas para INSERT", tableName)
		} else {
			return fmt.Sprintf("-- Table %s has no columns available for INSERT", tableName)
		}
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
		if qb.language == "pt" {
			return fmt.Sprintf("-- Não foi possível encontrar a estrutura da tabela %s para UPDATE", tableName)
		} else {
			return fmt.Sprintf("-- Unable to find the structure of table %s for UPDATE", tableName)
		}
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

	// If we don't have updatable columns, we can't create the UPDATE
	if len(setClause) == 0 {
		if qb.language == "pt" {
			return fmt.Sprintf("-- A tabela %s não possui colunas atualizáveis", tableName)
		} else {
			return fmt.Sprintf("-- Table %s has no updatable columns", tableName)
		}
	}

	// Build WHERE clause
	var whereClause string
	if len(conditions) > 0 {
		whereClause = qb.buildWhereClause(conditions)
	} else {
		// Default to primary key condition
		primaryKeyFound := false
		for _, col := range tableColumns {
			if col.IsPrimary {
				whereClause = col.Name + " = ?"
				primaryKeyFound = true
				break
			}
		}

		// If no primary key found, try id
		if !primaryKeyFound {
			for _, col := range tableColumns {
				if strings.ToLower(col.Name) == "id" {
					whereClause = col.Name + " = ?"
					break
				}
			}
		}

		// If we still don't have a WHERE clause, alert the user
		if whereClause == "" {
			if qb.language == "pt" {
				return fmt.Sprintf("-- UPDATE em %s requer uma condição WHERE. Especifique quais registros deseja atualizar.", tableName)
			} else {
				return fmt.Sprintf("-- UPDATE on %s requires a WHERE condition. Please specify which records to update.", tableName)
			}
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
		// Don't provide a default WHERE condition, as this can be dangerous
		// Alert the user that a WHERE condition is required
		if qb.language == "pt" {
			return fmt.Sprintf("-- DELETE em %s requer uma condição WHERE. Especifique quais registros deseja excluir.", tableName)
		} else {
			return fmt.Sprintf("-- DELETE on %s requires a WHERE condition. Please specify which records to delete.", tableName)
		}
	}

	return fmt.Sprintf("DELETE FROM %s WHERE %s", tableName, whereClause)
}
