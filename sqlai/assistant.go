package sqlai

import (
	"encoding/json"
	"fmt"
	"math"
	"regexp"
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
			confidenceThreshold:  0.65,
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

	// Apply context enrichment to improve model understanding
	contextEnrichedPrompt := s.enrichPromptWithContext(analyzedPrompt)

	// Check the query cache first (use the original prompt for caching)
	normalizedPrompt := strings.TrimSpace(strings.ToLower(prompt))
	if cachedSQL, exists := s.dbStructure.QueryCache[normalizedPrompt]; exists {
		return cachedSQL, nil
	}

	// Extract key information with semantic understanding
	tables, columns, conditions := s.extractEntitiesAdvanced(contextEnrichedPrompt)

	// Match intent to SQL operation with confidence scoring
	operation, subOperations := s.detectOperationAdvanced(contextEnrichedPrompt)

	// Generate SQL based on intent, tables, columns and conditions
	sql := s.buildSQLQueryAdvanced(operation, subOperations, tables, columns, conditions, contextEnrichedPrompt)

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

// detectOperationAdvanced detects the database operation using a more sophisticated approach
func (s *SQLAssistant) detectOperationAdvanced(prompt string) (string, []string) {
	// Get potential entities from the prompt
	entities := s.identifyPotentialEntities(prompt)
	normalizedPrompt := strings.ToLower(prompt)

	// Define operation context patterns for different languages
	operationPatterns := map[string]map[string][]string{
		"en": {
			"select": {
				`(?i)select`, `(?i)show`, `(?i)list`, `(?i)get`, `(?i)display`, `(?i)fetch`, `(?i)retrieve`,
				`(?i)find`, `(?i)search`, `(?i)query`, `(?i)view`, `(?i)show me`, `(?i)give me`,
				`(?i)what are`, `(?i)which are`, `(?i)who are`, `(?i)where are`,
			},
			"count": {
				`(?i)count`, `(?i)how many`, `(?i)number of`, `(?i)total number`,
				`(?i)quantity of`, `(?i)sum of`, `(?i)tally`,
			},
			"insert": {
				`(?i)insert`, `(?i)add`, `(?i)create`, `(?i)register`, `(?i)put`, `(?i)save`,
				`(?i)store`, `(?i)record`, `(?i)add new`, `(?i)new record`,
			},
			"update": {
				`(?i)update`, `(?i)change`, `(?i)modify`, `(?i)edit`, `(?i)alter`, `(?i)adjust`,
				`(?i)revise`, `(?i)correct`, `(?i)amend`,
			},
			"delete": {
				`(?i)delete`, `(?i)remove`, `(?i)drop`, `(?i)eliminate`, `(?i)discard`, `(?i)erase`,
				`(?i)purge`, `(?i)expunge`, `(?i)clear`, `(?i)wipe`,
			},
			"join": {
				`(?i)join`, `(?i)connect`, `(?i)link`, `(?i)associate`, `(?i)relate`, `(?i)combine`,
				`(?i)merge`, `(?i)along with`, `(?i)together with`, `(?i)attached to`,
				`(?i)linked to`, `(?i)associated with`, `(?i)related to`,
			},
			"group": {
				`(?i)group by`, `(?i)grouped by`, `(?i)organize by`, `(?i)organized by`,
				`(?i)categorize by`, `(?i)categorized by`, `(?i)aggregate by`,
				`(?i)cluster by`, `(?i)clustered by`,
			},
			"sort": {
				`(?i)sort by`, `(?i)sorted by`, `(?i)order by`, `(?i)ordered by`,
				`(?i)arrange by`, `(?i)arranged by`, `(?i)rank by`, `(?i)ranked by`,
			},
		},
		"pt": {
			"select": {
				`(?i)selecionar`, `(?i)selecione`, `(?i)mostrar`, `(?i)mostre`, `(?i)listar`, `(?i)liste`,
				`(?i)exibir`, `(?i)exiba`, `(?i)consultar`, `(?i)consulte`, `(?i)buscar`, `(?i)busque`,
				`(?i)encontrar`, `(?i)encontre`, `(?i)quais são`, `(?i)quem são`, `(?i)onde estão`,
				`(?i)me mostre`, `(?i)me traga`, `(?i)me dê`, `(?i)ver`, `(?i)veja`, `(?i)visualizar`,
				`(?i)visualize`, `(?i)obter`, `(?i)obtenha`, `(?i)trazer`, `(?i)traga`,
				`(?i)apresentar`, `(?i)apresente`,
			},
			"count": {
				`(?i)contar`, `(?i)conte`, `(?i)quantos`, `(?i)quantas`, `(?i)quantidade de`,
				`(?i)número de`, `(?i)total de`, `(?i)somatória de`, `(?i)soma de`,
				`(?i)contagem de`, `(?i)totalizar`,
			},
			"insert": {
				`(?i)inserir`, `(?i)insira`, `(?i)adicionar`, `(?i)adicione`, `(?i)incluir`, `(?i)inclua`,
				`(?i)criar`, `(?i)crie`, `(?i)cadastrar`, `(?i)cadastre`, `(?i)registrar`, `(?i)registre`,
				`(?i)salvar`, `(?i)salve`, `(?i)guardar`, `(?i)guarde`, `(?i)armazenar`, `(?i)armazene`,
				`(?i)novo registro`, `(?i)nova entrada`, `(?i)grave`, `(?i)grave`,
			},
			"update": {
				`(?i)atualizar`, `(?i)atualize`, `(?i)alterar`, `(?i)altere`, `(?i)modificar`, `(?i)modifique`,
				`(?i)editar`, `(?i)edite`, `(?i)mudar`, `(?i)mude`, `(?i)corrigir`, `(?i)corrija`,
				`(?i)revisar`, `(?i)revise`, `(?i)emendar`, `(?i)emende`, `(?i)ajustar`, `(?i)ajuste`,
				`(?i)reescrever`, `(?i)reescreva`, `(?i)atualização`, `(?i)mudança`, `(?i)alteração`, `(?i)modificação`,
			},
			"delete": {
				`(?i)deletar`, `(?i)delete`, `(?i)excluir`, `(?i)exclua`, `(?i)remover`, `(?i)remova`,
				`(?i)eliminar`, `(?i)elimine`, `(?i)apagar`, `(?i)apague`, `(?i)retirar`, `(?i)retire`,
				`(?i)descartar`, `(?i)descarte`, `(?i)anular`, `(?i)anule`, `(?i)cancelar`, `(?i)cancele`,
				`(?i)limpar`, `(?i)limpe`, `(?i)suprimir`, `(?i)suprima`,
			},
			"join": {
				`(?i)juntar`, `(?i)junte`, `(?i)unir`, `(?i)una`, `(?i)conectar`, `(?i)conecte`,
				`(?i)vincular`, `(?i)vincule`, `(?i)associar`, `(?i)associe`, `(?i)relacionar`, `(?i)relacione`,
				`(?i)combinar`, `(?i)combine`, `(?i)ligar`, `(?i)ligue`, `(?i)junto com`, `(?i)juntamente com`,
				`(?i)vinculado a`, `(?i)associado a`, `(?i)relacionado a`, `(?i)combinado com`,
				`(?i)e seus respectivos`, `(?i)e suas respectivas`, `(?i)cruzamento com`, `(?i)cruzar com`,
			},
			"group": {
				`(?i)agrupar por`, `(?i)agrupe por`, `(?i)agrupado por`, `(?i)agrupados por`,
				`(?i)agrupar em`, `(?i)agrupe em`, `(?i)agrupamento por`, `(?i)por grupos de`,
				`(?i)separar por`, `(?i)separe por`, `(?i)separados por`, `(?i)categorizar por`,
				`(?i)categorize por`, `(?i)categorizado por`, `(?i)categorização por`,
			},
			"sort": {
				`(?i)ordenar por`, `(?i)ordene por`, `(?i)ordenado por`, `(?i)ordenados por`,
				`(?i)ordenação por`, `(?i)classificar por`, `(?i)classifique por`, `(?i)classificado por`,
				`(?i)classificados por`, `(?i)organizar por`, `(?i)organize por`, `(?i)organizado por`,
				`(?i)organizados por`, `(?i)em ordem`, `(?i)por ordem`,
			},
			"subquery": {
				`(?i)para cada`, `(?i)que contenha`, `(?i)que possua`, `(?i)que tenha`, `(?i)onde exista`,
				`(?i)que existe em`, `(?i)que não existe em`, `(?i)para os quais`, `(?i)para as quais`,
				`(?i)para aqueles que`, `(?i)para aquelas que`, `(?i)que estejam em`,
				`(?i)que particip(a|am|e|em) de`, `(?i)que faz(em)? parte de`,
				`(?i)cujo`, `(?i)cuja`, `(?i)cujos`, `(?i)cujas`,
			},
		},
	}

	// Define SQL syntax fragments for detection
	sqlFragments := map[string][]string{
		"select": {
			`SELECT`, `FROM`, `WHERE`, `HAVING`, `DISTINCT`, `ALL`,
		},
		"count": {
			`COUNT\(`, `COUNT \(`, `SUM\(`, `SUM \(`, `AVG\(`, `AVG \(`,
			`MIN\(`, `MIN \(`, `MAX\(`, `MAX \(`,
		},
		"insert": {
			`INSERT`, `INTO`, `VALUES`, `DEFAULT`, `SET`,
		},
		"update": {
			`UPDATE`, `SET`, `=`, `:=`,
		},
		"delete": {
			`DELETE`, `FROM`, `TRUNCATE`,
		},
		"join": {
			`JOIN`, `INNER JOIN`, `LEFT JOIN`, `RIGHT JOIN`, `FULL JOIN`,
			`CROSS JOIN`, `NATURAL JOIN`, `USING\(`, `USING \(`, `ON`,
		},
		"group": {
			`GROUP BY`, `HAVING`,
		},
		"sort": {
			`ORDER BY`, `ASC`, `DESC`,
		},
	}

	// Use the patterns based on the detected language
	patterns, exists := operationPatterns[s.detectedLanguage]
	if !exists {
		// Default to English if the language is not supported
		patterns = operationPatterns["en"]
	}

	// Calculate confidence scores for each operation
	scores := map[string]float64{
		"select": 0.1, // Small default weight for select as it's the most common operation
		"count":  0,
		"insert": 0,
		"update": 0,
		"delete": 0,
		"join":   0,
		"group":  0,
		"sort":   0,
	}

	// Check for direct operation mentions
	for op, regexPatterns := range patterns {
		for _, pattern := range regexPatterns {
			if match, _ := regexp.MatchString(pattern, normalizedPrompt); match {
				// Match found, increase confidence score
				scores[op] += 0.5 // Strong indicator when operation pattern is directly matched
			}
		}
	}

	// Check for keywords that suggest operations
	keywords := map[string][]string{
		"select": {"show", "list", "find", "display", "view", "mostrar", "listar", "encontrar", "exibir", "visualizar"},
		"count":  {"how many", "count", "number of", "quantos", "quantas", "conte", "contar", "total", "contar quantos"},
		"insert": {"new", "add", "create", "inserir", "adicionar", "criar", "novo", "nova", "cadastrar"},
		"update": {"change", "modify", "edit", "alterar", "modificar", "editar", "corrigir", "atualizar", "mudar"},
		"delete": {"remove", "delete", "exclude", "remover", "deletar", "excluir", "apagar", "eliminar"},
		"join":   {"combine", "join", "with", "related", "junto", "juntar", "com", "relacionado", "relacionar", "unir"},
		"group":  {"group", "categorize", "cluster", "agrupar", "categorizar", "por categoria", "por grupo"},
		"sort":   {"order", "sort", "arrange", "ordenar", "classificar", "organizar", "em ordem"},
	}

	for op, terms := range keywords {
		for _, term := range terms {
			if strings.Contains(strings.ToLower(prompt), term) {
				// Keyword found, add medium weight
				scores[op] += 0.3
			}
		}
	}

	// Check for SQL syntax fragments
	for op, fragments := range sqlFragments {
		for _, fragment := range fragments {
			if match, _ := regexp.MatchString(fmt.Sprintf(`(?i)%s`, fragment), prompt); match {
				// SQL fragment found, add high weight
				scores[op] += 0.4
			}
		}
	}

	// Context modifiers: certain context clues make operations more likely
	// Presence of specific entities makes operations more likely
	if len(entities) > 0 {
		for _, entity := range entities {
			if entity != "" {
				// Having a clear entity increases select score
				scores["select"] += 0.3
			}
		}
	}

	// Check for question marks (suggests query/select/count)
	if strings.Contains(prompt, "?") {
		scores["select"] += 0.2
		scores["count"] += 0.2
	}

	// Check for conditions (suggest filtering)
	conditionPatterns := []string{
		`(?i)where`, `(?i)when`, `(?i)if`, `(?i)only`, `(?i)except`, `(?i)not`,
		`(?i)with`, `(?i)having`, `(?i)onde`, `(?i)quando`, `(?i)se`, `(?i)apenas`,
		`(?i)exceto`, `(?i)não`, `(?i)com`, `(?i)que tenham`, `(?i)que possuam`,
	}
	for _, pattern := range conditionPatterns {
		if match, _ := regexp.MatchString(pattern, prompt); match {
			scores["select"] += 0.2 // Filtering operations usually involve selection
		}
	}

	// Detailed analysis for composite or complex operations
	hasSubquery := false

	// Check for subquery patterns in Portuguese
	if s.detectedLanguage == "pt" && patterns["subquery"] != nil {
		for _, pattern := range patterns["subquery"] {
			if match, _ := regexp.MatchString(pattern, prompt); match {
				hasSubquery = true
				break
			}
		}
	}

	// Look for patterns that indicate nested queries or complex operations
	nestedPatterns := []string{
		`(?i)within`, `(?i)in the`, `(?i)inside`, `(?i)nested`, `(?i)sub`,
		`(?i)dentro`, `(?i)no interior`, `(?i)aninhado`, `(?i)sub`,
		`(?i)for each`, `(?i)para cada`, `(?i)que estejam em`,
	}

	for _, pattern := range nestedPatterns {
		if match, _ := regexp.MatchString(pattern, prompt); match {
			hasSubquery = true
			break
		}
	}

	// If subquery is detected, adjust scores
	if hasSubquery {
		scores["select"] += 0.4 // Subqueries are often selects
	}

	// Identify the operation with the highest score
	bestOp := "select" // Default
	highestScore := scores["select"]

	for op, score := range scores {
		if score > highestScore {
			highestScore = score
			bestOp = op
		}
	}

	// Sub-operation detection - these are specialized forms of operations
	if bestOp == "select" {
		// Check if it's specifically a count operation
		if scores["count"] > 0.6 {
			bestOp = "count"
		}
	}

	// Combine additional operations with the main operation if they have high scores
	additionalOps := []string{}
	for op, score := range scores {
		// Don't include the best operation again, and only include significant operations
		if op != bestOp && score > 0.7 {
			additionalOps = append(additionalOps, op)
		}
	}

	// For now, return the best operation. In a more advanced implementation,
	// you could return the main operation plus significant additional operations.
	return bestOp, additionalOps
}

// extractEntitiesAdvanced extracts tables, columns and conditions with improved semantic understanding
func (s *SQLAssistant) extractEntitiesAdvanced(prompt string) ([]models.TableInfo, []models.ColumnInfo, []models.Condition) {
	var tables []models.TableInfo
	var columns []models.ColumnInfo
	var conditions []models.Condition

	// Normalize the prompt for more accurate matching
	normalizedPrompt := strings.ToLower(prompt)
	normalizedPrompt = strings.ReplaceAll(normalizedPrompt, "_", " ")
	normalizedPrompt = " " + normalizedPrompt + " " // Add spaces to aid in word boundary detection

	// Identify table names with improved confidence scores
	tableScores := make(map[string]float64)
	for _, table := range s.dbStructure.Tables {
		tableName := strings.ToLower(table.Name)
		normalizedTableName := strings.ReplaceAll(tableName, "_", " ")

		// Exact match with word boundaries has highest confidence
		if strings.Contains(normalizedPrompt, " "+tableName+" ") ||
			strings.Contains(normalizedPrompt, " "+normalizedTableName+" ") {
			tableScores[table.Name] = 1.0
			continue
		}

		// Try singular/plural variations with improved word boundary detection
		singular := strings.TrimSuffix(tableName, "s")
		plural := tableName + "s"
		singularNormalized := strings.ReplaceAll(singular, "_", " ")
		pluralNormalized := strings.ReplaceAll(plural, "_", " ")

		if strings.Contains(normalizedPrompt, " "+singular+" ") ||
			strings.Contains(normalizedPrompt, " "+singularNormalized+" ") {
			tableScores[table.Name] = 0.95
			continue
		}

		if strings.Contains(normalizedPrompt, " "+plural+" ") ||
			strings.Contains(normalizedPrompt, " "+pluralNormalized+" ") {
			tableScores[table.Name] = 0.95
			continue
		}

		// Check for exact match with word boundaries at start or end
		if strings.HasPrefix(normalizedPrompt, " "+tableName+" ") ||
			strings.HasSuffix(normalizedPrompt, " "+tableName+" ") ||
			strings.HasPrefix(normalizedPrompt, " "+normalizedTableName+" ") ||
			strings.HasSuffix(normalizedPrompt, " "+normalizedTableName+" ") {
			tableScores[table.Name] = 0.95
			continue
		}

		// Check for exact match pattern with word boundaries
		pattern := `\b` + regexp.QuoteMeta(tableName) + `\b`
		if operations.RegexMatch(pattern, normalizedPrompt) {
			tableScores[table.Name] = 0.9
			continue
		}

		pattern = `\b` + regexp.QuoteMeta(normalizedTableName) + `\b`
		if operations.RegexMatch(pattern, normalizedPrompt) {
			tableScores[table.Name] = 0.9
			continue
		}

		// Handle tables with underscores by checking each part
		tableNameParts := strings.Split(tableName, "_")
		if len(tableNameParts) > 1 {
			matchCount := 0
			for _, part := range tableNameParts {
				if len(part) > 2 && strings.Contains(normalizedPrompt, " "+part+" ") {
					matchCount++
				}
			}

			if matchCount == len(tableNameParts) {
				tableScores[table.Name] = 0.85
				continue
			} else if matchCount > 0 {
				partialScore := 0.5 + (0.3 * float64(matchCount) / float64(len(tableNameParts)))
				tableScores[table.Name] = partialScore
				continue
			}
		}

		// Check for partial word match as lowest confidence
		if strings.Contains(normalizedPrompt, tableName) ||
			strings.Contains(normalizedPrompt, normalizedTableName) {
			tableScores[table.Name] = 0.5 + (0.1 * float64(len(tableName)) / 20.0) // Longer name matches get slightly higher confidence
		}

		// Table name description matching - if provided
		if table.Description != "" && len(table.Description) > 3 {
			desc := strings.ToLower(table.Description)
			if strings.Contains(normalizedPrompt, desc) {
				tableScores[table.Name] = 0.85
			}
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

	// Extract columns from the most probable tables with improved matching
	columnScores := make(map[string]map[string]float64) // tableName -> colName -> score

	for _, tableInfo := range tables {
		columnScores[tableInfo.Name] = make(map[string]float64)

		for _, table := range s.dbStructure.Tables {
			if table.Name == tableInfo.Name {
				for _, col := range table.Columns {
					colName := strings.ToLower(col.Name)

					// Normalized column name for better matching
					normalizedColName := strings.ReplaceAll(colName, "_", " ")

					// Start with a base confidence derived from table confidence
					baseConfidence := tableInfo.Confidence * 0.8

					// Exact match with word boundaries
					pattern := `\b` + regexp.QuoteMeta(colName) + `\b`
					if operations.RegexMatch(pattern, normalizedPrompt) {
						columnScores[tableInfo.Name][col.Name] = math.Max(1.0, columnScores[tableInfo.Name][col.Name])
						continue
					}

					// Check for normalized column name (spaces instead of underscores)
					pattern = `\b` + regexp.QuoteMeta(normalizedColName) + `\b`
					if operations.RegexMatch(pattern, normalizedPrompt) {
						columnScores[tableInfo.Name][col.Name] = math.Max(0.95, columnScores[tableInfo.Name][col.Name])
						continue
					}

					// Check for match with column description if available
					if col.Description != "" && len(col.Description) > 3 {
						desc := strings.ToLower(col.Description)
						if strings.Contains(normalizedPrompt, desc) {
							columnScores[tableInfo.Name][col.Name] = math.Max(0.9, columnScores[tableInfo.Name][col.Name])
							continue
						}
					}

					// Check for partial matches on columns with underscores
					colParts := strings.Split(colName, "_")
					if len(colParts) > 1 {
						matchCount := 0
						for _, part := range colParts {
							if len(part) > 2 && strings.Contains(normalizedPrompt, " "+part+" ") {
								matchCount++
							}
						}

						if matchCount == len(colParts) {
							columnScores[tableInfo.Name][col.Name] = math.Max(0.85, columnScores[tableInfo.Name][col.Name])
							continue
						} else if matchCount > 0 {
							partialScore := baseConfidence + (0.2 * float64(matchCount) / float64(len(colParts)))
							columnScores[tableInfo.Name][col.Name] = math.Max(partialScore, columnScores[tableInfo.Name][col.Name])
							continue
						}
					}

					// Apply context-based confidence boost for certain column types
					if col.IsPrimary || strings.Contains(colName, "id") {
						// Primary keys and ID fields get a slight boost if any ID-related term is in prompt
						if strings.Contains(normalizedPrompt, "id") ||
							strings.Contains(normalizedPrompt, " identifier") ||
							strings.Contains(normalizedPrompt, " identificador") {
							columnScores[tableInfo.Name][col.Name] = math.Max(baseConfidence+0.15, columnScores[tableInfo.Name][col.Name])
						}
					}

					// Add a slight confidence boost for common columns based on query context
					if util.ContainsAny(colName, []string{"name", "title", "description", "descricao", "nome", "titulo"}) {
						if util.ContainsAny(normalizedPrompt, []string{"name", "title", "description", "nome", "titulo", "descricao"}) {
							columnScores[tableInfo.Name][col.Name] = math.Max(baseConfidence+0.2, columnScores[tableInfo.Name][col.Name])
						}
					}

					// Date columns get a boost when time-related terms are present
					if util.ContainsAny(colName, []string{"date", "time", "data", "hora", "timestamp", "created", "updated"}) {
						if util.ContainsAny(normalizedPrompt, []string{"date", "time", "when", "data", "hora", "quando", "periodo", "today", "yesterday", "hoje", "ontem"}) {
							columnScores[tableInfo.Name][col.Name] = math.Max(baseConfidence+0.2, columnScores[tableInfo.Name][col.Name])
						}
					}

					// Simple partial match as lowest confidence
					if strings.Contains(normalizedPrompt, colName) || strings.Contains(normalizedPrompt, normalizedColName) {
						columnScores[tableInfo.Name][col.Name] = math.Max(baseConfidence, columnScores[tableInfo.Name][col.Name])
					}

					// For sample values matching
					if len(col.SampleValues) > 0 {
						for _, sample := range col.SampleValues {
							if sample != "" && len(sample) > 2 && strings.Contains(normalizedPrompt, strings.ToLower(sample)) {
								columnScores[tableInfo.Name][col.Name] = math.Max(0.85, columnScores[tableInfo.Name][col.Name])
								break
							}
						}
					}
				}
			}
		}
	}

	// Build the final columns list from the confidence scores
	for tableName, colScores := range columnScores {
		for colName, confidence := range colScores {
			if confidence >= s.confidenceThreshold {
				// Find the column type and primary key status
				colType := ""
				isPrimary := false

				for _, table := range s.dbStructure.Tables {
					if table.Name == tableName {
						for _, col := range table.Columns {
							if col.Name == colName {
								colType = col.Type
								isPrimary = col.IsPrimary
								break
							}
						}
						break
					}
				}

				columns = append(columns, models.ColumnInfo{
					Name:       colName,
					TableName:  tableName,
					Type:       colType,
					IsPrimary:  isPrimary,
					Confidence: confidence,
				})
			}
		}
	}

	// If no tables found, try harder with improved fuzzy matching
	if len(tables) == 0 {
		// Use contextual hints to guess the most likely tables
		potentialEntities := s.identifyPotentialEntities(prompt)

		for _, entity := range potentialEntities {
			// Try to find a table with similar name
			bestMatch := ""
			bestScore := 0.0

			for _, table := range s.dbStructure.Tables {
				tableName := strings.ToLower(table.Name)
				similarity := util.JaroWinklerSimilarity(entity, tableName)

				if similarity > 0.8 && similarity > bestScore {
					bestMatch = table.Name
					bestScore = similarity
				}

				// Also check singular/plural forms
				singular := strings.TrimSuffix(tableName, "s")
				if singular != tableName {
					similarity = util.JaroWinklerSimilarity(entity, singular)
					if similarity > 0.8 && similarity > bestScore {
						bestMatch = table.Name
						bestScore = similarity
					}
				}

				plural := tableName + "s"
				similarity = util.JaroWinklerSimilarity(entity, plural)
				if similarity > 0.8 && similarity > bestScore {
					bestMatch = table.Name
					bestScore = similarity
				}
			}

			if bestMatch != "" {
				tables = append(tables, models.TableInfo{
					Name:       bestMatch,
					Confidence: 0.6,
				})
				break
			}
		}

		// If still no tables found, try matching based on columns mentioned
		if len(tables) == 0 {
			columnsInPrompt := s.extractPotentialColumnNames(prompt)
			tableScores := make(map[string]float64)

			for _, colName := range columnsInPrompt {
				for _, table := range s.dbStructure.Tables {
					for _, col := range table.Columns {
						if strings.EqualFold(col.Name, colName) ||
							strings.EqualFold(strings.ReplaceAll(col.Name, "_", ""), strings.ReplaceAll(colName, " ", "")) {
							tableScores[table.Name] += 0.3
						}
					}
				}
			}

			// Find the best table match
			var bestTable string
			var bestScore float64 = 0.0

			for tableName, score := range tableScores {
				if score > bestScore {
					bestScore = score
					bestTable = tableName
				}
			}

			if bestTable != "" && bestScore >= 0.3 {
				tables = append(tables, models.TableInfo{
					Name:       bestTable,
					Confidence: 0.5,
				})
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
						if col.IsPrimary ||
							strings.ToLower(col.Name) == "id" ||
							strings.HasSuffix(strings.ToLower(col.Name), "_id") ||
							util.ContainsAny(strings.ToLower(col.Name), []string{"name", "nome", "title", "titulo", "description", "descricao", "date", "data"}) {
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

// extractPotentialColumnNames extracts possible column names from the prompt
func (s *SQLAssistant) extractPotentialColumnNames(prompt string) []string {
	var possibleColumns []string
	words := strings.Fields(strings.ToLower(prompt))

	// Single words that could be column names
	for _, word := range words {
		// Clean the word from punctuation
		word = strings.Trim(word, ",.;:()\"'!?")
		if len(word) > 2 && !util.IsCommonWord(word, s.detectedLanguage) {
			possibleColumns = append(possibleColumns, word)
		}
	}

	// Look for phrases that might be column names (with underscores)
	for i := 0; i < len(words)-1; i++ {
		phrase := words[i] + "_" + words[i+1]
		possibleColumns = append(possibleColumns, phrase)
	}

	return possibleColumns
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
			// If we still don't have tables but have DB info,
			// check if the query mentions specific columns
			for _, table := range s.dbStructure.Tables {
				// Skip system or metadata tables
				if strings.HasPrefix(table.Name, "sys_") ||
					strings.HasPrefix(table.Name, "meta_") ||
					strings.HasPrefix(table.Name, "pg_") {
					continue
				}

				// Check if any table column is mentioned in the query
				for _, col := range table.Columns {
					if strings.Contains(strings.ToLower(prompt), strings.ToLower(col.Name)) {
						tables = append(tables, models.TableInfo{
							Name:       table.Name,
							Confidence: 0.5,
						})
						break
					}
				}

				if len(tables) > 0 {
					break // We found a candidate table
				}
			}

			// If we still don't have tables, we don't add a default automatically
			// Just continue with empty tables so operation-specific functions handle this case
		}
	}

	// Check for explicit table and column mentions
	explicitTableRegex := regexp.MustCompile(`(?i)(?:da|de|from)\s+tabela\s+(\w+)`)
	if matches := explicitTableRegex.FindStringSubmatch(prompt); len(matches) > 1 {
		explicitTable := matches[1]
		// Verify this table exists in our schema
		tableExists := false
		for _, dbTable := range s.dbStructure.Tables {
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

	// Check for explicit column mentions
	explicitColRegex := regexp.MustCompile(`(?i)(?:a coluna|as colunas|coluna|colunas|the column|the columns|column|columns)\s+(\w+)`)
	if matches := explicitColRegex.FindStringSubmatch(prompt); len(matches) > 1 && len(tables) > 0 {
		explicitCol := matches[1]
		// If we have an explicit column mentioned, override any previously detected columns
		columns = []models.ColumnInfo{{
			Name:       explicitCol,
			TableName:  tables[0].Name,
			Confidence: 1.0,
		}}
	}

	// If no columns were detected, select appropriate columns
	if len(columns) == 0 && len(tables) > 0 {
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
		if len(tables) > 0 {
			return queryBuilder.BuildInsertQuery(tables[0].Name, prompt)
		}
		if s.detectedLanguage == "pt" {
			return "-- Não foi possível identificar a tabela para o INSERT"
		} else {
			return "-- Unable to identify the table for INSERT"
		}

	case "update":
		if len(tables) > 0 {
			return queryBuilder.BuildUpdateQuery(tables[0].Name, conditions, prompt)
		}
		if s.detectedLanguage == "pt" {
			return "-- Não foi possível identificar a tabela para o UPDATE"
		} else {
			return "-- Unable to identify the table for UPDATE"
		}

	case "delete":
		if len(tables) > 0 {
			return queryBuilder.BuildDeleteQuery(tables[0].Name, conditions, prompt)
		}
		if s.detectedLanguage == "pt" {
			return "-- Não foi possível identificar a tabela para o DELETE"
		} else {
			return "-- Unable to identify the table for DELETE"
		}

	default:
		// Default to a SELECT query
		return queryBuilder.BuildSelectQuery(tables, columns, conditions, subOperations, prompt)
	}
}

// validateAndOptimizeSQL performs simple validation and optimization on the generated SQL
func (s *SQLAssistant) validateAndOptimizeSQL(sql string) string {
	// If query starts with a comment, it means there was an error in generation
	// We won't try to "fix" this type of query
	if strings.HasPrefix(strings.TrimSpace(sql), "--") {
		// Translate the comment to English if it's in Portuguese
		if s.detectedLanguage == "pt" {
			sql = translateErrorCommentToEnglish(sql)
		}
		return sql
	}

	// Simple validation checks
	if !strings.Contains(sql, "FROM") && strings.Contains(sql, "SELECT") {
		// In this case, we have a SELECT without FROM
		// We won't add a default table, as this could be problematic
		// Better to return a comment explaining the problem
		if s.detectedLanguage == "pt" {
			return "-- Consulta SELECT incompleta. Falta especificar a tabela (FROM)."
		} else {
			return "-- Incomplete SELECT query. Missing table specification (FROM)."
		}
	}

	// Check for missing table in simple queries
	fromRegex := operations.CompileRegex(`FROM\s+(\w+)`)
	if !fromRegex.MatchString(sql) && strings.Contains(sql, "SELECT") {
		// We won't add a default table anymore
		// Return a comment explaining the problem
		if s.detectedLanguage == "pt" {
			return "-- Consulta SELECT incompleta. Falta especificar a tabela (FROM)."
		} else {
			return "-- Incomplete SELECT query. Missing table specification (FROM)."
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

// translateErrorCommentToEnglish translates Portuguese error comments to English
func translateErrorCommentToEnglish(comment string) string {
	// Common error messages and their translations
	translations := map[string]string{
		"-- Não foi possível determinar qual tabela consultar. Por favor, especifique a tabela em sua consulta.": "-- Unable to determine which table to query. Please specify the table in your query.",
		"-- Não foi possível determinar qual tabela consultar":                                                   "-- Unable to determine which table to query",
		"-- Consulta SELECT incompleta. Falta especificar a tabela (FROM).":                                      "-- Incomplete SELECT query. Missing table specification (FROM).",
		"-- Não foi possível identificar a tabela para o INSERT":                                                 "-- Unable to identify the table for INSERT",
		"-- Não foi possível identificar a tabela para o UPDATE":                                                 "-- Unable to identify the table for UPDATE",
		"-- Não foi possível identificar a tabela para o DELETE":                                                 "-- Unable to identify the table for DELETE",
		"-- Não foi possível determinar qual tabela agrupar. Por favor, especifique a tabela em sua consulta.":   "-- Unable to determine which table to group. Please specify the table in your query.",
		"-- Não foi possível determinar qual tabela agrupar":                                                     "-- Unable to determine which table to group",
		"-- Não foi possível identificar colunas para agrupar em":                                                "-- Unable to identify columns to group by in",
	}

	// Replace known error messages
	for ptError, enError := range translations {
		if strings.HasPrefix(comment, ptError) {
			// Handle the case where there's additional text after the standard message
			remainingText := strings.TrimPrefix(comment, ptError)
			return enError + remainingText
		}
	}

	// Special cases with variable content
	deleteRegex := regexp.MustCompile(`-- DELETE em (\w+) requer uma condição WHERE. Especifique quais registros deseja excluir.`)
	if matches := deleteRegex.FindStringSubmatch(comment); len(matches) > 1 {
		return fmt.Sprintf("-- DELETE on %s requires a WHERE condition. Please specify which records to delete.", matches[1])
	}

	updateRegex := regexp.MustCompile(`-- UPDATE em (\w+) requer uma condição WHERE. Especifique quais registros deseja atualizar.`)
	if matches := updateRegex.FindStringSubmatch(comment); len(matches) > 1 {
		return fmt.Sprintf("-- UPDATE on %s requires a WHERE condition. Please specify which records to update.", matches[1])
	}

	insertRegex := regexp.MustCompile(`-- A tabela (\w+) não possui colunas para INSERT`)
	if matches := insertRegex.FindStringSubmatch(comment); len(matches) > 1 {
		return fmt.Sprintf("-- Table %s has no columns available for INSERT", matches[1])
	}

	noStructureRegex := regexp.MustCompile(`-- Não foi possível encontrar a estrutura da tabela (\w+) para INSERT`)
	if matches := noStructureRegex.FindStringSubmatch(comment); len(matches) > 1 {
		return fmt.Sprintf("-- Unable to find the structure of table %s for INSERT", matches[1])
	}

	// If no match found, return the original comment
	return comment
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

// enrichPromptWithContext enhances the prompt with database structural information
func (s *SQLAssistant) enrichPromptWithContext(prompt string) string {
	// Skip if there are no tables in the structure
	if len(s.dbStructure.Tables) == 0 {
		return prompt
	}

	// Detect potential entities in the prompt
	potentialEntities := s.identifyPotentialEntities(prompt)
	if len(potentialEntities) == 0 {
		return prompt
	}

	// Find relevant tables that match the detected entities
	var relevantTables []models.TableForAI
	for _, entity := range potentialEntities {
		for _, table := range s.dbStructure.Tables {
			// Check for exact match or plural form
			singularEntity := strings.TrimSuffix(entity, "s")
			pluralEntity := entity + "s"

			// Use Jaro-Winkler similarity for fuzzy matching of table names
			similarityThreshold := 0.85

			// Try different forms of the entity name
			entityForms := []string{entity, singularEntity, pluralEntity}
			for _, form := range entityForms {
				similarity := util.JaroWinklerSimilarity(strings.ToLower(form), strings.ToLower(table.Name))

				if similarity >= similarityThreshold ||
					strings.Contains(strings.ToLower(table.Name), strings.ToLower(form)) ||
					strings.Contains(strings.ToLower(form), strings.ToLower(table.Name)) {
					// Found a match
					relevantTables = append(relevantTables, table)
					break
				}
			}
		}
	}

	// If no relevant tables found by entity names, try to find by column names
	if len(relevantTables) == 0 {
		columnEntities := s.extractPotentialColumnNames(prompt)

		for _, column := range columnEntities {
			for _, table := range s.dbStructure.Tables {
				for _, col := range table.Columns {
					// Use Jaro-Winkler similarity for fuzzy matching of column names
					similarity := util.JaroWinklerSimilarity(strings.ToLower(column), strings.ToLower(col.Name))

					if similarity >= 0.8 ||
						strings.Contains(strings.ToLower(col.Name), strings.ToLower(column)) ||
						strings.Contains(strings.ToLower(column), strings.ToLower(col.Name)) {
						// Found a table containing a matching column
						relevantTables = append(relevantTables, table)
						break
					}
				}
			}
		}
	}

	// Limit the number of relevant tables to avoid excessive context
	// Keep only the top 3 tables if there are too many
	if len(relevantTables) > 3 {
		relevantTables = relevantTables[:3]
	}

	// No relevant tables found
	if len(relevantTables) == 0 {
		return prompt
	}

	// Build context information recursively
	var contextBuilder strings.Builder
	tablesSeen := make(map[string]bool)

	contextBuilder.WriteString("\nContext:\n")

	// Process each table with depth tracking to avoid cycles
	for _, table := range relevantTables {
		if !tablesSeen[table.Name] {
			// Process the table and add information to the context builder
			buildTableContext(&contextBuilder, table, s.dbStructure.Tables, tablesSeen, 0, 2) // Set max depth to 2
		}
	}

	// Only add context if we found useful information
	if contextBuilder.Len() > 15 { // More than just "Context:" header
		return prompt + contextBuilder.String()
	}

	return prompt
}

// buildTableContext recursively builds context information about a table, its columns, and relations
func buildTableContext(sb *strings.Builder, table models.TableForAI, allTables []models.TableForAI, tablesSeen map[string]bool, depth int, maxDepth int) {
	if depth > maxDepth {
		return // Prevent infinite recursion
	}

	// Mark this table as seen
	tablesSeen[table.Name] = true

	// Add table information
	sb.WriteString(fmt.Sprintf("Table: %s", table.Name))
	if table.Description != "" {
		sb.WriteString(fmt.Sprintf(" - %s", table.Description))
	}
	sb.WriteString("\n")

	// Add columns
	sb.WriteString("  Columns: ")
	cols := []string{}
	for _, col := range table.Columns {
		colDesc := col.Name
		if col.IsPrimary {
			colDesc += " (PK)"
		}
		// Check for foreign keys
		for _, fk := range table.ForeignKeys {
			if fk.ColumnName == col.Name {
				colDesc += fmt.Sprintf(" (FK -> %s.%s)", fk.ReferencedTable, fk.ReferencedColumn)
				break
			}
		}
		if col.Type != "" {
			colDesc += fmt.Sprintf(" (%s)", col.Type)
		}
		cols = append(cols, colDesc)
	}
	sb.WriteString(strings.Join(cols, ", "))
	sb.WriteString("\n")

	// Find and add relationships
	for _, fk := range table.ForeignKeys {
		if !tablesSeen[fk.ReferencedTable] {
			// Find the related table
			var relatedTable *models.TableForAI
			for i := range allTables {
				if allTables[i].Name == fk.ReferencedTable {
					relatedTable = &allTables[i]
					break
				}
			}

			if relatedTable != nil {
				// Process the related table recursively at the next depth level
				if depth < maxDepth {
					sb.WriteString(fmt.Sprintf("  Related to: %s through %s.%s -> %s.%s\n",
						relatedTable.Name, table.Name, fk.ColumnName, relatedTable.Name, fk.ReferencedColumn))
					buildTableContext(sb, *relatedTable, allTables, tablesSeen, depth+1, maxDepth)
				} else {
					sb.WriteString(fmt.Sprintf("  Related to: %s (not expanded due to depth limit)\n", relatedTable.Name))
				}
			}
		}
	}

	// Add a separator after each table
	sb.WriteString("\n")
}

// identifyPotentialEntities extracts potential entity names from the prompt
func (s *SQLAssistant) identifyPotentialEntities(prompt string) []string {
	var entities []string
	normalizedPrompt := strings.ToLower(prompt)

	// Use PromptAnalyzer to extract likely entities
	contextEntities := s.promptAnalyzer.ExtractEntityNameFromContext(normalizedPrompt, "")
	if contextEntities != "" {
		entities = append(entities, contextEntities)
	}

	// Extract words that might be table names
	words := strings.Fields(normalizedPrompt)
	for _, word := range words {
		// Skip very short words, common words, and those with special chars
		word = strings.Trim(word, ",.;:!?()\"'")
		if len(word) <= 2 || util.IsCommonWord(word, s.detectedLanguage) ||
			strings.ContainsAny(word, "/\\@#$%^&*+=<>{}[]|") {
			continue
		}

		// Normalize singular/plural forms
		singularForm := strings.TrimSuffix(word, "s")
		if singularForm != word && len(singularForm) > 2 {
			entities = append(entities, singularForm)
		}

		// Add original word too
		if !util.Contains(entities, word) {
			entities = append(entities, word)
		}
	}

	// Look for compound entities (words separated by spaces)
	for i := 0; i < len(words)-1; i++ {
		compound := words[i] + "_" + words[i+1]
		compound = strings.Trim(compound, ",.;:!?()\"'")
		if len(compound) > 4 && !util.IsCommonWord(words[i], s.detectedLanguage) &&
			!util.IsCommonWord(words[i+1], s.detectedLanguage) {
			entities = append(entities, compound)
		}
	}

	// Prioritize entity hints
	var prioritizedEntities []string
	entityHints := map[string][]string{
		"pt": {"tabela", "entidade", "cadastro", "registro", "relação"},
		"en": {"table", "entity", "record", "relation"},
	}

	hintsToCheck := entityHints["en"]
	if s.detectedLanguage == "pt" {
		hintsToCheck = entityHints["pt"]
	}

	// Search for patterns like "table users", "registros de clientes", etc.
	for _, hint := range hintsToCheck {
		for i, word := range words {
			if word == hint && i+1 < len(words) {
				nextWord := strings.Trim(words[i+1], ",.;:!?()\"'")
				if len(nextWord) > 2 && !util.IsCommonWord(nextWord, s.detectedLanguage) {
					// Prioritize entities that appear after a hint word
					prioritizedEntities = append(prioritizedEntities, nextWord)
				}

				// Also check for "de/da/do" pattern in Portuguese
				if s.detectedLanguage == "pt" && i+2 < len(words) {
					if (words[i+1] == "de" || words[i+1] == "da" || words[i+1] == "do" || words[i+1] == "dos" || words[i+1] == "das") && i+2 < len(words) {
						nextNextWord := strings.Trim(words[i+2], ",.;:!?()\"'")
						if len(nextNextWord) > 2 && !util.IsCommonWord(nextNextWord, s.detectedLanguage) {
							prioritizedEntities = append(prioritizedEntities, nextNextWord)
						}
					}
				}

				// Also check for "of" pattern in English
				if s.detectedLanguage == "en" && i+2 < len(words) {
					if words[i+1] == "of" && i+2 < len(words) {
						nextNextWord := strings.Trim(words[i+2], ",.;:!?()\"'")
						if len(nextNextWord) > 2 && !util.IsCommonWord(nextNextWord, s.detectedLanguage) {
							prioritizedEntities = append(prioritizedEntities, nextNextWord)
						}
					}
				}
			}
		}
	}

	// Combine prioritized entities with regular ones, removing duplicates
	for _, entity := range prioritizedEntities {
		if !util.Contains(entities, entity) {
			entities = append([]string{entity}, entities...) // Add at the beginning for priority
		}
	}

	// Check words that appear after "from", "da", "de", etc.
	fromHints := map[string][]string{
		"pt": {"de", "da", "do", "dos", "das"},
		"en": {"from", "in", "of"},
	}

	hintsToCheck = fromHints["en"]
	if s.detectedLanguage == "pt" {
		hintsToCheck = fromHints["pt"]
	}

	for _, hint := range hintsToCheck {
		for i, word := range words {
			if word == hint && i+1 < len(words) {
				nextWord := strings.Trim(words[i+1], ",.;:!?()\"'")
				if len(nextWord) > 2 && !util.IsCommonWord(nextWord, s.detectedLanguage) && !util.Contains(entities, nextWord) {
					entities = append([]string{nextWord}, entities...) // Add at the beginning for priority
				}
			}
		}
	}

	return entities
}

// SetLanguageForTesting define o idioma para teste
func (s *SQLAssistant) SetLanguageForTesting(lang string) {
	s.detectedLanguage = lang
	s.updateLanguageFunctions()
}

// DetectOperationForTesting expõe a detecção de operação para testes
func (s *SQLAssistant) DetectOperationForTesting(prompt string) string {
	op, _ := s.detectOperationAdvanced(prompt)
	return op
}
