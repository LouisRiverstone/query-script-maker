package sqlai

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"sql_script_maker/sqlai/language"
	"sql_script_maker/sqlai/models"
	"sql_script_maker/sqlai/operations"
	"sql_script_maker/sqlai/util"
)

// Singleton pattern for the SQL AI model
var (
	instance     *SQLAssistant
	instanceOnce sync.Once
)

// SQLAssistant is the main AI assistant for SQL generation
type SQLAssistant struct {
	dbStructure          models.DatabaseStructureForAI
	templates            map[string]string
	history              []models.QueryHistory
	dialect              models.SQLDialect
	operationPatterns    map[string][]string
	aggregationFunctions map[string]string
	dateFunctions        map[string]string
	stringFunctions      map[string]string
	confidenceThreshold  float64
	maxJoinDepth         int
	learningMode         bool
	promptAnalyzer       *language.PromptAnalyzer
	detectedLanguage     string
}

// GetSQLAssistant returns the singleton instance of the SQL assistant
func GetSQLAssistant() *SQLAssistant {
	instanceOnce.Do(func() {
		instance = &SQLAssistant{
			templates:            language.KeywordMap,
			dbStructure:          models.DatabaseStructureForAI{QueryCache: make(map[string]string)},
			operationPatterns:    language.OperationPatternsEN,
			aggregationFunctions: language.AggregationFunctionsEN,
			dateFunctions:        language.DateFunctionsEN,
			stringFunctions:      language.StringFunctionsEN,
			confidenceThreshold:  0.7,
			maxJoinDepth:         3,
			learningMode:         true,
			promptAnalyzer:       language.NewPromptAnalyzer(),
			detectedLanguage:     "en", // Default to English
		}
	})
	return instance
}

// Init initializes the SQL assistant with the database structure
func (s *SQLAssistant) Init(structureJSON string) error {
	if structureJSON == "" {
		return fmt.Errorf("empty database structure")
	}

	var structure models.DatabaseStructureForAI
	err := json.Unmarshal([]byte(structureJSON), &structure)
	if err != nil {
		return err
	}

	s.dbStructure = structure
	s.dbStructure.QueryCache = make(map[string]string)
	s.dbStructure.LastUpdated = time.Now()

	// Set dialect based on DB type
	s.setDialect(structure.DBType)

	return nil
}

// setDialect configures the SQL dialect specifics
func (s *SQLAssistant) setDialect(dbType string) {
	switch strings.ToLower(dbType) {
	case "mysql":
		s.dialect = models.SQLDialect{
			Name:                "MySQL",
			LimitSyntax:         "LIMIT %d OFFSET %d",
			SupportsCTE:         true,
			SupportsWindowFuncs: true,
		}
	case "postgresql", "postgres":
		s.dialect = models.SQLDialect{
			Name:                "PostgreSQL",
			LimitSyntax:         "LIMIT %d OFFSET %d",
			SupportsCTE:         true,
			SupportsWindowFuncs: true,
		}
	case "sqlite":
		s.dialect = models.SQLDialect{
			Name:                "SQLite",
			LimitSyntax:         "LIMIT %d OFFSET %d",
			SupportsCTE:         true,
			SupportsWindowFuncs: true,
		}
	case "mssql", "sqlserver":
		s.dialect = models.SQLDialect{
			Name:                "SQL Server",
			LimitSyntax:         "OFFSET %d ROWS FETCH NEXT %d ROWS ONLY",
			SupportsCTE:         true,
			SupportsWindowFuncs: true,
		}
	case "oracle":
		s.dialect = models.SQLDialect{
			Name:                "Oracle",
			LimitSyntax:         "OFFSET %d ROWS FETCH NEXT %d ROWS ONLY",
			SupportsCTE:         true,
			SupportsWindowFuncs: true,
		}
	default:
		// Default to a generic SQL dialect
		s.dialect = models.SQLDialect{
			Name:                "Generic SQL",
			LimitSyntax:         "LIMIT %d OFFSET %d",
			SupportsCTE:         false,
			SupportsWindowFuncs: false,
		}
	}
}

// GenerateSQL generates SQL from a natural language query
func (s *SQLAssistant) GenerateSQL(prompt string) (string, error) {
	if s.dbStructure.Tables == nil || len(s.dbStructure.Tables) == 0 {
		return "", fmt.Errorf("database structure not initialized")
	}

	// Analyze and normalize the prompt and detect language
	analyzedPrompt, detectedLang := s.promptAnalyzer.AnalyzePrompt(prompt)
	s.detectedLanguage = detectedLang

	// Update function maps based on detected language
	s.updateLanguageFunctions()

	// Check the query cache first (use the original prompt for caching)
	normalizedPrompt := strings.TrimSpace(strings.ToLower(prompt))
	if cachedSQL, exists := s.dbStructure.QueryCache[normalizedPrompt]; exists {
		return cachedSQL, nil
	}

	// Extract key information with semantic understanding
	tables, columns, conditions := s.extractEntitiesAdvanced(analyzedPrompt)

	// Match intent to SQL operation with confidence scoring
	operation, subOperations := s.detectOperationAdvanced(analyzedPrompt)

	// Generate SQL based on intent, tables, columns and conditions
	sql := s.buildSQLQueryAdvanced(operation, subOperations, tables, columns, conditions, analyzedPrompt)

	// Validate the generated SQL
	sql = s.validateAndOptimizeSQL(sql)

	// Cache the result
	s.dbStructure.QueryCache[normalizedPrompt] = sql

	// Add to query history for learning
	if s.learningMode {
		s.history = append(s.history, models.QueryHistory{
			Query:     sql,
			CreatedAt: time.Now(),
			Success:   true,
		})
	}

	return sql, nil
}

// updateLanguageFunctions updates the function maps based on detected language
func (s *SQLAssistant) updateLanguageFunctions() {
	s.operationPatterns = language.GetOperationPatterns(s.detectedLanguage)
	s.aggregationFunctions = language.GetAggregationFunctions(s.detectedLanguage)
	s.dateFunctions = language.GetDateFunctions(s.detectedLanguage)
	s.stringFunctions = language.GetStringFunctions(s.detectedLanguage)
}

// detectOperationAdvanced determines the SQL operation based on the prompt with improved pattern matching
func (s *SQLAssistant) detectOperationAdvanced(prompt string) (string, []string) {
	// Map to store operation match confidence scores
	operationScores := make(map[string]float64)

	// Check for operation patterns with regex
	for operation, patterns := range s.operationPatterns {
		for _, pattern := range patterns {
			if operations.RegexMatch(pattern, prompt) {
				operationScores[operation] += 0.5
			}
		}
	}

	// Identify keyword presence
	keywordSets := operations.GetOperationKeywords(s.detectedLanguage)

	for operation, keywords := range keywordSets {
		for _, keyword := range keywords {
			if strings.Contains(prompt, keyword) {
				operationScores[operation] += 0.3
			}
		}
	}

	// Find operation with highest score
	var bestOperation string
	var highestScore float64

	for operation, score := range operationScores {
		if score > highestScore {
			highestScore = score
			bestOperation = operation
		}
	}

	// Get sub-operations (additional operations that should be included)
	var subOperations []string
	for operation, score := range operationScores {
		if operation != bestOperation && score >= s.confidenceThreshold {
			subOperations = append(subOperations, operation)
		}
	}

	// Default to select if no clear operation detected
	if bestOperation == "" || highestScore < s.confidenceThreshold {
		bestOperation = "select"
	}

	return bestOperation, subOperations
}

// extractEntitiesAdvanced extracts tables, columns and conditions with improved semantic understanding
func (s *SQLAssistant) extractEntitiesAdvanced(prompt string) ([]models.TableInfo, []models.ColumnInfo, []models.Condition) {
	var tables []models.TableInfo
	var columns []models.ColumnInfo
	var conditions []models.Condition

	// Identify table names with confidence scores
	tableScores := make(map[string]float64)
	for _, table := range s.dbStructure.Tables {
		tableName := strings.ToLower(table.Name)

		// Exact match has highest confidence
		if strings.Contains(prompt, " "+tableName+" ") ||
			strings.HasPrefix(prompt, tableName+" ") ||
			strings.HasSuffix(prompt, " "+tableName) {
			tableScores[table.Name] = 1.0
			continue
		}

		// Try singular/plural variations
		singular := strings.TrimSuffix(tableName, "s")
		plural := tableName + "s"

		if strings.Contains(prompt, " "+singular+" ") ||
			strings.HasPrefix(prompt, singular+" ") ||
			strings.HasSuffix(prompt, " "+singular) {
			tableScores[table.Name] = 0.9
			continue
		}

		if strings.Contains(prompt, " "+plural+" ") ||
			strings.HasPrefix(prompt, plural+" ") ||
			strings.HasSuffix(prompt, " "+plural) {
			tableScores[table.Name] = 0.9
			continue
		}

		// Check for partial matches with word boundaries
		pattern := `\b` + tableName + `\b`
		if operations.RegexMatch(pattern, prompt) {
			tableScores[table.Name] = 0.8
			continue
		}

		// Check for partial word match
		if strings.Contains(prompt, tableName) {
			tableScores[table.Name] = 0.6
		}
	}

	// Extract tables with confidence above threshold
	for tableName, score := range tableScores {
		if score >= s.confidenceThreshold {
			tables = append(tables, models.TableInfo{
				Name:       tableName,
				Confidence: score,
			})
		}
	}

	// Sort tables by confidence score (descending)
	operations.SortTablesByConfidence(tables)

	// Extract columns from the most probable tables
	for _, tableInfo := range tables {
		for _, table := range s.dbStructure.Tables {
			if table.Name == tableInfo.Name {
				for _, col := range table.Columns {
					colName := strings.ToLower(col.Name)
					colScore := 0.0

					// Exact match with word boundaries
					pattern := `\b` + colName + `\b`
					if operations.RegexMatch(pattern, prompt) {
						colScore = 1.0
					} else if strings.Contains(prompt, colName) {
						// Partial match
						colScore = 0.7
					}

					// Check for natural language variations of column names
					wordColName := strings.ReplaceAll(colName, "_", " ")
					if colName != wordColName && strings.Contains(prompt, wordColName) {
						colScore = 0.9
					}

					if colScore >= s.confidenceThreshold {
						columns = append(columns, models.ColumnInfo{
							Name:       col.Name,
							TableName:  table.Name,
							Type:       col.Type,
							IsPrimary:  col.IsPrimary,
							Confidence: colScore,
						})
					}
				}
			}
		}
	}

	// If no tables found, try harder with fuzzy matching
	if len(tables) == 0 {
		for _, table := range s.dbStructure.Tables {
			tableName := strings.ToLower(table.Name)

			// Break table name into words for partial matching
			tableWords := strings.Split(tableName, "_")
			for _, word := range tableWords {
				if len(word) > 3 && strings.Contains(prompt, word) {
					tables = append(tables, models.TableInfo{
						Name:       table.Name,
						Confidence: 0.5,
					})
					break
				}
			}
		}
	}

	// Extract conditions from the prompt
	conditions = operations.ExtractConditions(prompt, tables, columns, s.detectedLanguage)

	// If no columns found but tables were found, include all primary key and commonly used columns
	if len(columns) == 0 && len(tables) > 0 {
		for _, tableInfo := range tables {
			for _, table := range s.dbStructure.Tables {
				if table.Name == tableInfo.Name {
					// Add primary key columns
					for _, col := range table.Columns {
						if col.IsPrimary || strings.ToLower(col.Name) == "id" ||
							util.ContainsAny(strings.ToLower(col.Name), []string{"name", "title", "description", "date"}) {
							columns = append(columns, models.ColumnInfo{
								Name:       col.Name,
								TableName:  table.Name,
								Type:       col.Type,
								IsPrimary:  col.IsPrimary,
								Confidence: 0.7,
							})
						}
					}
				}
			}
		}
	}

	return tables, columns, conditions
}

// buildSQLQueryAdvanced builds complex SQL queries based on the detected entities and operations
func (s *SQLAssistant) buildSQLQueryAdvanced(operation string, subOperations []string,
	tables []models.TableInfo, columns []models.ColumnInfo, conditions []models.Condition, prompt string) string {

	// Handle case when no tables were detected
	if len(tables) == 0 {
		// Try to find the most relevant table based on the query context
		mostRelevantTable := operations.FindMostRelevantTable(prompt, s.dbStructure.Tables)
		if mostRelevantTable != "" {
			tables = append(tables, models.TableInfo{
				Name:       mostRelevantTable,
				Confidence: 0.6,
			})
		} else if len(s.dbStructure.Tables) > 0 {
			// If still no tables, take the first table
			tables = append(tables, models.TableInfo{
				Name:       s.dbStructure.Tables[0].Name,
				Confidence: 0.5,
			})
		}
	}

	// If no columns were detected, select appropriate columns
	if len(columns) == 0 {
		columns = operations.SelectRelevantColumns(tables, operation, prompt, s.dbStructure.Tables)
	}

	// Handle different types of operations
	queryBuilder := operations.NewQueryBuilder(s.detectedLanguage, s.dbStructure.Tables)

	switch operation {
	case "select":
		return queryBuilder.BuildSelectQuery(tables, columns, conditions, subOperations, prompt)

	case "count":
		return queryBuilder.BuildCountQuery(tables, columns, conditions, prompt)

	case "join":
		return queryBuilder.BuildJoinQuery(tables, columns, conditions, subOperations, prompt)

	case "group":
		return queryBuilder.BuildGroupByQuery(tables, columns, conditions, prompt)

	case "order":
		return queryBuilder.BuildOrderByQuery(tables, columns, conditions, prompt)

	case "limit":
		return queryBuilder.BuildLimitQuery(tables, columns, conditions, prompt)

	case "distinct":
		return queryBuilder.BuildDistinctQuery(tables, columns, conditions, prompt)

	case "insert":
		return queryBuilder.BuildInsertQuery(tables[0].Name, prompt)

	case "update":
		return queryBuilder.BuildUpdateQuery(tables[0].Name, conditions, prompt)

	case "delete":
		return queryBuilder.BuildDeleteQuery(tables[0].Name, conditions, prompt)

	default:
		// Default to a SELECT query
		return queryBuilder.BuildSelectQuery(tables, columns, conditions, subOperations, prompt)
	}
}

// validateAndOptimizeSQL performs simple validation and optimization on the generated SQL
func (s *SQLAssistant) validateAndOptimizeSQL(sql string) string {
	// Simple validation checks
	if !strings.Contains(sql, "FROM") && strings.Contains(sql, "SELECT") {
		sql = strings.Replace(sql, "SELECT", "SELECT * FROM", 1)
	}

	// Check for missing table in simple queries
	fromRegex := operations.CompileRegex(`FROM\s+(\w+)`)
	if !fromRegex.MatchString(sql) && strings.Contains(sql, "SELECT") {
		// Add a default table if possible
		if len(s.dbStructure.Tables) > 0 {
			tableName := s.dbStructure.Tables[0].Name
			if !strings.Contains(sql, "FROM") {
				sql += " FROM " + tableName
			}
		}
	}

	// Optimize: Replace SELECT * when joining multiple tables
	if strings.Contains(sql, "JOIN") && strings.Contains(sql, "SELECT *") {
		// Replace with specific columns from each table
		selectRegex := operations.CompileRegex(`SELECT\s+\*`)
		tablesInQuery := util.ExtractTablesFromSQL(sql)

		if len(tablesInQuery) >= 2 {
			var columns []string

			for _, tableName := range tablesInQuery {
				for _, table := range s.dbStructure.Tables {
					if table.Name == tableName {
						// Add primary key and a few important columns
						for _, col := range table.Columns {
							if col.IsPrimary ||
								strings.ToLower(col.Name) == "id" ||
								util.ContainsAny(strings.ToLower(col.Name),
									[]string{"name", "title", "created_at"}) {
								columns = append(columns, tableName+"."+col.Name)
							}
						}
					}
				}
			}

			if len(columns) > 0 {
				newSelect := "SELECT " + strings.Join(columns, ", ")
				sql = selectRegex.ReplaceAllString(sql, newSelect)
			}
		}
	}

	return sql
}

// RecordQueryFeedback records the outcome of a query execution for learning
func (s *SQLAssistant) RecordQueryFeedback(feedback models.FeedbackResult) {
	if s.learningMode {
		// Find the query in history and update results
		for i, hist := range s.history {
			if hist.Query == feedback.Query {
				s.history[i].Success = feedback.WasSuccessful
				s.history[i].ResultCount = feedback.RowCount
				s.history[i].ExecutionTime = feedback.ExecutionTime
				break
			}
		}
	}
}

// Reset clears the assistant's cached data
func (s *SQLAssistant) Reset() {
	s.dbStructure = models.DatabaseStructureForAI{QueryCache: make(map[string]string)}
	s.history = nil
	s.detectedLanguage = "en"
}
