package language

import (
	"fmt"
	"regexp"
	"strings"
)

// PromptAnalyzer processes natural language prompts for SQL generation
type PromptAnalyzer struct {
	language string
}

// NewPromptAnalyzer creates a new prompt analyzer instance
func NewPromptAnalyzer() *PromptAnalyzer {
	return &PromptAnalyzer{
		language: "en", // Default to English
	}
}

// AnalyzePrompt performs deeper semantic analysis of the prompt
func (p *PromptAnalyzer) AnalyzePrompt(prompt string) (string, string) {
	// Detect the language first
	p.language = DetectLanguage(prompt)

	// Normalize prompt based on the detected language
	normalizedPrompt := p.normalizePrompt(prompt)

	return normalizedPrompt, p.language
}

// normalizePrompt standardizes and enhances the prompt for processing
func (p *PromptAnalyzer) normalizePrompt(prompt string) string {
	switch p.language {
	case "pt":
		return p.normalizePortuguesePrompt(prompt)
	default:
		return p.normalizeEnglishPrompt(prompt)
	}
}

// normalizeEnglishPrompt normalizes English prompts
func (p *PromptAnalyzer) normalizeEnglishPrompt(prompt string) string {
	// Normalize variations of common terms
	prompt = regexp.MustCompile(`(?i)how many`).ReplaceAllString(prompt, "count")
	prompt = regexp.MustCompile(`(?i)show me`).ReplaceAllString(prompt, "select")
	prompt = regexp.MustCompile(`(?i)give me`).ReplaceAllString(prompt, "select")
	prompt = regexp.MustCompile(`(?i)find all`).ReplaceAllString(prompt, "select")
	prompt = regexp.MustCompile(`(?i)get all`).ReplaceAllString(prompt, "select")

	// Standardize condition expressions
	prompt = regexp.MustCompile(`(?i)greater than or equal to`).ReplaceAllString(prompt, ">=")
	prompt = regexp.MustCompile(`(?i)less than or equal to`).ReplaceAllString(prompt, "<=")
	prompt = regexp.MustCompile(`(?i)greater than`).ReplaceAllString(prompt, ">")
	prompt = regexp.MustCompile(`(?i)less than`).ReplaceAllString(prompt, "<")
	prompt = regexp.MustCompile(`(?i)equal to`).ReplaceAllString(prompt, "=")
	prompt = regexp.MustCompile(`(?i)not equal to`).ReplaceAllString(prompt, "!=")

	// Normalize join expressions
	prompt = regexp.MustCompile(`(?i)along with their`).ReplaceAllString(prompt, "join")
	prompt = regexp.MustCompile(`(?i)together with`).ReplaceAllString(prompt, "join")
	prompt = regexp.MustCompile(`(?i)combined with`).ReplaceAllString(prompt, "join")

	// Identify and normalize time-related expressions
	timeExpressions := map[string]string{
		"(?i)last (\\d+) days":      "date >= DATE_SUB(CURRENT_DATE(), INTERVAL $1 DAY)",
		"(?i)past (\\d+) weeks":     "date >= DATE_SUB(CURRENT_DATE(), INTERVAL $1 WEEK)",
		"(?i)last (\\d+) months":    "date >= DATE_SUB(CURRENT_DATE(), INTERVAL $1 MONTH)",
		"(?i)previous (\\d+) years": "date >= DATE_SUB(CURRENT_DATE(), INTERVAL $1 YEAR)",
	}

	for pattern, replacement := range timeExpressions {
		re := regexp.MustCompile(pattern)
		prompt = re.ReplaceAllString(prompt, replacement)
	}

	return prompt
}

// normalizePortuguesePrompt normalizes Portuguese prompts
func (p *PromptAnalyzer) normalizePortuguesePrompt(prompt string) string {
	// Normalize variations of common terms
	prompt = regexp.MustCompile(`(?i)quantos`).ReplaceAllString(prompt, "contar")
	prompt = regexp.MustCompile(`(?i)mostre-me`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)me mostre`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)encontre todos`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)obter todos`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)listar todos`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)exibir todos`).ReplaceAllString(prompt, "selecionar")

	// Standardize condition expressions
	prompt = regexp.MustCompile(`(?i)maior ou igual a`).ReplaceAllString(prompt, ">=")
	prompt = regexp.MustCompile(`(?i)menor ou igual a`).ReplaceAllString(prompt, "<=")
	prompt = regexp.MustCompile(`(?i)maior que`).ReplaceAllString(prompt, ">")
	prompt = regexp.MustCompile(`(?i)menor que`).ReplaceAllString(prompt, "<")
	prompt = regexp.MustCompile(`(?i)igual a`).ReplaceAllString(prompt, "=")
	prompt = regexp.MustCompile(`(?i)diferente de`).ReplaceAllString(prompt, "!=")
	prompt = regexp.MustCompile(`(?i)não igual a`).ReplaceAllString(prompt, "!=")

	// Normalize join expressions
	prompt = regexp.MustCompile(`(?i)junto com seus`).ReplaceAllString(prompt, "juntar")
	prompt = regexp.MustCompile(`(?i)junto com seus`).ReplaceAllString(prompt, "juntar")
	prompt = regexp.MustCompile(`(?i)combinado com`).ReplaceAllString(prompt, "juntar")
	prompt = regexp.MustCompile(`(?i)relacionado com`).ReplaceAllString(prompt, "juntar")
	prompt = regexp.MustCompile(`(?i)associado com`).ReplaceAllString(prompt, "juntar")

	// Identify and normalize time-related expressions
	timeExpressions := map[string]string{
		"(?i)últimos (\\d+) dias":    "data >= DATE_SUB(CURRENT_DATE(), INTERVAL $1 DAY)",
		"(?i)últimas (\\d+) semanas": "data >= DATE_SUB(CURRENT_DATE(), INTERVAL $1 WEEK)",
		"(?i)últimos (\\d+) meses":   "data >= DATE_SUB(CURRENT_DATE(), INTERVAL $1 MONTH)",
		"(?i)últimos (\\d+) anos":    "data >= DATE_SUB(CURRENT_DATE(), INTERVAL $1 YEAR)",
	}

	for pattern, replacement := range timeExpressions {
		re := regexp.MustCompile(pattern)
		prompt = re.ReplaceAllString(prompt, replacement)
	}

	return prompt
}

// ExtractEntityNameFromContext extracts the most likely entity name from the context
func (p *PromptAnalyzer) ExtractEntityNameFromContext(prompt, context string) string {
	// Check for entities in quotes first
	quotePattern := `["']([^"']+)["']`
	re := regexp.MustCompile(quotePattern)
	matches := re.FindStringSubmatch(prompt)
	if len(matches) > 1 {
		return matches[1]
	}

	// Look for common entity extraction patterns based on language
	var patterns []string
	if p.language == "pt" {
		patterns = []string{
			`(?i)(?:tabela|entidade|cadastro)\s+de\s+([a-z0-9_]+)`,
			`(?i)(?:os|as)\s+([a-z0-9_]+)`,
			`(?i)(?:todos|todas)\s+(?:os|as)\s+([a-z0-9_]+)`,
			`(?i)(?:listar|mostrar|selecionar)\s+([a-z0-9_]+)`,
			`(?i)(?:no|na|nos|nas)\s+([a-z0-9_]+)`,
		}
	} else {
		patterns = []string{
			`(?i)(?:table|entity|from)\s+([a-z0-9_]+)`,
			`(?i)(?:all|the)\s+([a-z0-9_]+)`,
			`(?i)(?:list|show|select)\s+([a-z0-9_]+)`,
			`(?i)(?:in|from|of)\s+the\s+([a-z0-9_]+)`,
		}
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(prompt)
		if len(matches) > 1 && len(matches[1]) > 2 {
			// Check if this might be a table name by seeing if it exists in context
			if strings.Contains(context, "\""+matches[1]+"\"") ||
				strings.Contains(context, "'"+matches[1]+"'") ||
				strings.Contains(context, " "+matches[1]+" ") {
				return matches[1]
			}
		}
	}

	return ""
}

// ExtractAttributesFromPrompt extracts column names or attributes from the prompt
func (p *PromptAnalyzer) ExtractAttributesFromPrompt(prompt string) []string {
	var attributes []string
	var patterns []string

	// Define patterns based on language
	if p.language == "pt" {
		patterns = []string{
			`(?i)campos\s+(?:como|como os)\s+([^,\.]+(?:,\s*[^,\.]+)*)`,
			`(?i)colunas\s+(?:como|como as)\s+([^,\.]+(?:,\s*[^,\.]+)*)`,
			`(?i)atributos\s+(?:como|como os)\s+([^,\.]+(?:,\s*[^,\.]+)*)`,
			`(?i)(?:mostrar|selecionar|exibir|obter)\s+([^,\.]+(?:,\s*[^,\.]+)*)\s+(?:de|da|dos|das)`,
		}
	} else {
		patterns = []string{
			`(?i)fields\s+(?:like|such as)\s+([^,\.]+(?:,\s*[^,\.]+)*)`,
			`(?i)columns\s+(?:like|such as)\s+([^,\.]+(?:,\s*[^,\.]+)*)`,
			`(?i)attributes\s+(?:like|such as)\s+([^,\.]+(?:,\s*[^,\.]+)*)`,
			`(?i)(?:show|select|get|retrieve)\s+([^,\.]+(?:,\s*[^,\.]+)*)\s+(?:from|of)`,
		}
	}

	// Try to extract attributes using patterns
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(prompt)
		if len(matches) > 1 {
			// Split by comma and clean up
			parts := strings.Split(matches[1], ",")
			for _, part := range parts {
				attribute := strings.TrimSpace(part)
				if len(attribute) > 0 {
					attributes = append(attributes, attribute)
				}
			}
			break
		}
	}

	return attributes
}

// ExtractConditionsFromPrompt extracts filter conditions from the prompt
func (p *PromptAnalyzer) ExtractConditionsFromPrompt(prompt string) map[string]string {
	conditions := make(map[string]string)
	var patterns []string

	// Define patterns based on language
	if p.language == "pt" {
		patterns = []string{
			`(?i)onde\s+([a-z0-9_]+)\s+(=|!=|>|<|>=|<=|igual a|diferente de|maior que|menor que|maior ou igual a|menor ou igual a)\s+([^,\s]+)`,
			`(?i)com\s+([a-z0-9_]+)\s+(=|!=|>|<|>=|<=|igual a|diferente de|maior que|menor que|maior ou igual a|menor ou igual a)\s+([^,\s]+)`,
			`(?i)que\s+(?:tenha|tem|possuem|possui)\s+([a-z0-9_]+)\s+(=|!=|>|<|>=|<=|igual a|diferente de|maior que|menor que|maior ou igual a|menor ou igual a)\s+([^,\s]+)`,
		}
	} else {
		patterns = []string{
			`(?i)where\s+([a-z0-9_]+)\s+(=|!=|>|<|>=|<=|equal to|not equal to|greater than|less than|greater than or equal to|less than or equal to)\s+([^,\s]+)`,
			`(?i)with\s+([a-z0-9_]+)\s+(=|!=|>|<|>=|<=|equal to|not equal to|greater than|less than|greater than or equal to|less than or equal to)\s+([^,\s]+)`,
			`(?i)that\s+(?:have|has|contains)\s+([a-z0-9_]+)\s+(=|!=|>|<|>=|<=|equal to|not equal to|greater than|less than|greater than or equal to|less than or equal to)\s+([^,\s]+)`,
		}
	}

	// Try to extract conditions using patterns
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringSubmatch(prompt, -1)
		for _, match := range matches {
			if len(match) > 3 {
				attribute := strings.TrimSpace(match[1])
				operator := strings.TrimSpace(match[2])
				value := strings.TrimSpace(match[3])

				// Normalize operators
				switch operator {
				case "igual a", "equal to":
					operator = "="
				case "diferente de", "not equal to":
					operator = "!="
				case "maior que", "greater than":
					operator = ">"
				case "menor que", "less than":
					operator = "<"
				case "maior ou igual a", "greater than or equal to":
					operator = ">="
				case "menor ou igual a", "less than or equal to":
					operator = "<="
				}

				conditions[attribute] = operator + " " + value
			}
		}
	}

	return conditions
}

// ExtractLimitFromPrompt extracts a limit value from the prompt
func (p *PromptAnalyzer) ExtractLimitFromPrompt(prompt string) int {
	var patterns []string

	// Define patterns based on language
	if p.language == "pt" {
		patterns = []string{
			`(?i)(?:limite|limitar|apenas|somente|top)\s+(\d+)`,
			`(?i)(?:primeiros|primeiras|últimos|últimas)\s+(\d+)`,
			`(?i)(\d+)\s+(?:resultados|registros|linhas)`,
		}
	} else {
		patterns = []string{
			`(?i)(?:limit|limited to|only|just|top)\s+(\d+)`,
			`(?i)(?:first|last)\s+(\d+)`,
			`(?i)(\d+)\s+(?:results|records|rows)`,
		}
	}

	// Try to extract limit using patterns
	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindStringSubmatch(prompt)
		if len(matches) > 1 {
			limit := strings.TrimSpace(matches[1])
			limitNum := 0
			_, err := fmt.Sscanf(limit, "%d", &limitNum)
			if err == nil && limitNum > 0 {
				return limitNum
			}
		}
	}

	// Default limit
	return 0
}
