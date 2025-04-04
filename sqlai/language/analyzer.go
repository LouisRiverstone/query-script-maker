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

	// Adicionando etapa para realizar identificação explícita de tabelas e colunas
	normalizedPrompt = p.enhanceTableColumnIdentifiers(normalizedPrompt)

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
	prompt = regexp.MustCompile(`(?i)quanto[s]? de`).ReplaceAllString(prompt, "contar")
	prompt = regexp.MustCompile(`(?i)qual (a|o) (quantidade|numero|total) de`).ReplaceAllString(prompt, "contar")
	prompt = regexp.MustCompile(`(?i)quantidad[e]? de`).ReplaceAllString(prompt, "contar")
	prompt = regexp.MustCompile(`(?i)mostre-me`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)me mostre`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)mostr[e|ar]`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)encontre todos`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)encontr[e|ar]`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)obter todos`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)obtenha`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)obt[e|er]`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)listar todos`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)list[e|ar]`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)exibir todos`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)exib[a|ir]`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)visualiz[e|ar]`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)quero ver`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)preciso ver`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)preciso saber`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)gostaria de ver`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)gostaria de saber`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)quais são`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)quais os`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)quais as`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)quais foram`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)me traga`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)traga-me`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)traz[er] para mim`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)consult[e|ar]`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)busc[a|ar]`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)pesquis[e|ar]`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)retorn[e|ar]`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)extrair`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)traz[er] os dados de`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)apresent[e|ar]`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)relacion[e|ar]`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)informar`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)indicar`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)fornecer`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)favor (listar|mostrar|exibir)`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)quero (uma lista|uma relação)`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)necessito (de|dos|das|ver)`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)ver (os|as|todos|todas)`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)detalhar`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)conhecer`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)identificar`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)localizar`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)verificar`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)checar`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)conferir`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)analisar`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)apurar`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)levantar`).ReplaceAllString(prompt, "selecionar")
	prompt = regexp.MustCompile(`(?i)recuperar`).ReplaceAllString(prompt, "selecionar")

	// Additional counting expressions
	prompt = regexp.MustCompile(`(?i)contar quantos`).ReplaceAllString(prompt, "contar")
	prompt = regexp.MustCompile(`(?i)totais de`).ReplaceAllString(prompt, "contar")
	prompt = regexp.MustCompile(`(?i)somatória de`).ReplaceAllString(prompt, "contar")
	prompt = regexp.MustCompile(`(?i)somar (os|as|todos|todas)`).ReplaceAllString(prompt, "contar")
	prompt = regexp.MustCompile(`(?i)totalizar`).ReplaceAllString(prompt, "contar")
	prompt = regexp.MustCompile(`(?i)computar`).ReplaceAllString(prompt, "contar")
	prompt = regexp.MustCompile(`(?i)calcular (o número|a quantidade)`).ReplaceAllString(prompt, "contar")
	prompt = regexp.MustCompile(`(?i)enumerar`).ReplaceAllString(prompt, "contar")
	prompt = regexp.MustCompile(`(?i)contabilizar`).ReplaceAllString(prompt, "contar")
	prompt = regexp.MustCompile(`(?i)número total de`).ReplaceAllString(prompt, "contar")
	prompt = regexp.MustCompile(`(?i)quantidade total de`).ReplaceAllString(prompt, "contar")
	prompt = regexp.MustCompile(`(?i)(qual|quais) (o|a) (número|quantidade) total`).ReplaceAllString(prompt, "contar")

	// Additional update expressions
	prompt = regexp.MustCompile(`(?i)mudar`).ReplaceAllString(prompt, "atualizar")
	prompt = regexp.MustCompile(`(?i)alteração`).ReplaceAllString(prompt, "atualizar")
	prompt = regexp.MustCompile(`(?i)modificação`).ReplaceAllString(prompt, "atualizar")
	prompt = regexp.MustCompile(`(?i)atualização`).ReplaceAllString(prompt, "atualizar")
	prompt = regexp.MustCompile(`(?i)revis(ar|ão)`).ReplaceAllString(prompt, "atualizar")
	prompt = regexp.MustCompile(`(?i)corrigir`).ReplaceAllString(prompt, "atualizar")
	prompt = regexp.MustCompile(`(?i)reescrever`).ReplaceAllString(prompt, "atualizar")
	prompt = regexp.MustCompile(`(?i)substituir`).ReplaceAllString(prompt, "atualizar")
	prompt = regexp.MustCompile(`(?i)trocar`).ReplaceAllString(prompt, "atualizar")
	prompt = regexp.MustCompile(`(?i)emendar`).ReplaceAllString(prompt, "atualizar")
	prompt = regexp.MustCompile(`(?i)renovar`).ReplaceAllString(prompt, "atualizar")
	prompt = regexp.MustCompile(`(?i)reformular`).ReplaceAllString(prompt, "atualizar")
	prompt = regexp.MustCompile(`(?i)melhorar`).ReplaceAllString(prompt, "atualizar")

	// Additional delete expressions
	prompt = regexp.MustCompile(`(?i)deletar`).ReplaceAllString(prompt, "excluir")
	prompt = regexp.MustCompile(`(?i)remover`).ReplaceAllString(prompt, "excluir")
	prompt = regexp.MustCompile(`(?i)eliminar`).ReplaceAllString(prompt, "excluir")
	prompt = regexp.MustCompile(`(?i)apagar`).ReplaceAllString(prompt, "excluir")
	prompt = regexp.MustCompile(`(?i)retirar`).ReplaceAllString(prompt, "excluir")
	prompt = regexp.MustCompile(`(?i)suprimir`).ReplaceAllString(prompt, "excluir")
	prompt = regexp.MustCompile(`(?i)descartar`).ReplaceAllString(prompt, "excluir")
	prompt = regexp.MustCompile(`(?i)anular`).ReplaceAllString(prompt, "excluir")
	prompt = regexp.MustCompile(`(?i)cancelar`).ReplaceAllString(prompt, "excluir")
	prompt = regexp.MustCompile(`(?i)limpar`).ReplaceAllString(prompt, "excluir")
	prompt = regexp.MustCompile(`(?i)excluir definitivamente`).ReplaceAllString(prompt, "excluir")
	prompt = regexp.MustCompile(`(?i)se livrar de`).ReplaceAllString(prompt, "excluir")

	// Additional insert expressions
	prompt = regexp.MustCompile(`(?i)inserir`).ReplaceAllString(prompt, "adicionar")
	prompt = regexp.MustCompile(`(?i)incluir`).ReplaceAllString(prompt, "adicionar")
	prompt = regexp.MustCompile(`(?i)adicionar`).ReplaceAllString(prompt, "adicionar")
	prompt = regexp.MustCompile(`(?i)acrescentar`).ReplaceAllString(prompt, "adicionar")
	prompt = regexp.MustCompile(`(?i)cadastrar`).ReplaceAllString(prompt, "adicionar")
	prompt = regexp.MustCompile(`(?i)registrar`).ReplaceAllString(prompt, "adicionar")
	prompt = regexp.MustCompile(`(?i)gravar`).ReplaceAllString(prompt, "adicionar")
	prompt = regexp.MustCompile(`(?i)anexar`).ReplaceAllString(prompt, "adicionar")
	prompt = regexp.MustCompile(`(?i)incorporar`).ReplaceAllString(prompt, "adicionar")
	prompt = regexp.MustCompile(`(?i)integrar`).ReplaceAllString(prompt, "adicionar")
	prompt = regexp.MustCompile(`(?i)salvar`).ReplaceAllString(prompt, "adicionar")
	prompt = regexp.MustCompile(`(?i)guardar`).ReplaceAllString(prompt, "adicionar")
	prompt = regexp.MustCompile(`(?i)armazenar`).ReplaceAllString(prompt, "adicionar")
	prompt = regexp.MustCompile(`(?i)criar ((um|uma) )?(nov[oa])`).ReplaceAllString(prompt, "adicionar")
	prompt = regexp.MustCompile(`(?i)novo registro`).ReplaceAllString(prompt, "adicionar")
	prompt = regexp.MustCompile(`(?i)nova entrada`).ReplaceAllString(prompt, "adicionar")

	// Standardize condition expressions
	prompt = regexp.MustCompile(`(?i)maior ou igual a`).ReplaceAllString(prompt, ">=")
	prompt = regexp.MustCompile(`(?i)maior ou igual`).ReplaceAllString(prompt, ">=")
	prompt = regexp.MustCompile(`(?i)no mínimo`).ReplaceAllString(prompt, ">=")
	prompt = regexp.MustCompile(`(?i)pelo menos`).ReplaceAllString(prompt, ">=")
	prompt = regexp.MustCompile(`(?i)a partir de`).ReplaceAllString(prompt, ">=")
	prompt = regexp.MustCompile(`(?i)menor ou igual a`).ReplaceAllString(prompt, "<=")
	prompt = regexp.MustCompile(`(?i)menor ou igual`).ReplaceAllString(prompt, "<=")
	prompt = regexp.MustCompile(`(?i)no máximo`).ReplaceAllString(prompt, "<=")
	prompt = regexp.MustCompile(`(?i)até no máximo`).ReplaceAllString(prompt, "<=")
	prompt = regexp.MustCompile(`(?i)não (ultrapassa|exceda)`).ReplaceAllString(prompt, "<=")
	prompt = regexp.MustCompile(`(?i)maior que`).ReplaceAllString(prompt, ">")
	prompt = regexp.MustCompile(`(?i)maior do que`).ReplaceAllString(prompt, ">")
	prompt = regexp.MustCompile(`(?i)acima de`).ReplaceAllString(prompt, ">")
	prompt = regexp.MustCompile(`(?i)superior a`).ReplaceAllString(prompt, ">")
	prompt = regexp.MustCompile(`(?i)mais de`).ReplaceAllString(prompt, ">")
	prompt = regexp.MustCompile(`(?i)menor que`).ReplaceAllString(prompt, "<")
	prompt = regexp.MustCompile(`(?i)menor do que`).ReplaceAllString(prompt, "<")
	prompt = regexp.MustCompile(`(?i)abaixo de`).ReplaceAllString(prompt, "<")
	prompt = regexp.MustCompile(`(?i)inferior a`).ReplaceAllString(prompt, "<")
	prompt = regexp.MustCompile(`(?i)menos de`).ReplaceAllString(prompt, "<")
	prompt = regexp.MustCompile(`(?i)igual a`).ReplaceAllString(prompt, "=")
	prompt = regexp.MustCompile(`(?i)exatamente`).ReplaceAllString(prompt, "=")
	prompt = regexp.MustCompile(`(?i)é igual a`).ReplaceAllString(prompt, "=")
	prompt = regexp.MustCompile(`(?i)equivalente a`).ReplaceAllString(prompt, "=")
	prompt = regexp.MustCompile(`(?i)seja`).ReplaceAllString(prompt, "=")
	prompt = regexp.MustCompile(`(?i)diferente de`).ReplaceAllString(prompt, "!=")
	prompt = regexp.MustCompile(`(?i)não igual a`).ReplaceAllString(prompt, "!=")
	prompt = regexp.MustCompile(`(?i)não é`).ReplaceAllString(prompt, "!=")
	prompt = regexp.MustCompile(`(?i)exceto`).ReplaceAllString(prompt, "!=")
	prompt = regexp.MustCompile(`(?i)sem ser`).ReplaceAllString(prompt, "!=")
	prompt = regexp.MustCompile(`(?i)com exceção de`).ReplaceAllString(prompt, "!=")
	prompt = regexp.MustCompile(`(?i)excluindo`).ReplaceAllString(prompt, "!=")
	prompt = regexp.MustCompile(`(?i)contendo`).ReplaceAllString(prompt, "LIKE")
	prompt = regexp.MustCompile(`(?i)que contém`).ReplaceAllString(prompt, "LIKE")
	prompt = regexp.MustCompile(`(?i)que contenha`).ReplaceAllString(prompt, "LIKE")
	prompt = regexp.MustCompile(`(?i)tem em seu conteúdo`).ReplaceAllString(prompt, "LIKE")
	prompt = regexp.MustCompile(`(?i)inclui`).ReplaceAllString(prompt, "LIKE")
	prompt = regexp.MustCompile(`(?i)possui`).ReplaceAllString(prompt, "LIKE")
	prompt = regexp.MustCompile(`(?i)apresenta`).ReplaceAllString(prompt, "LIKE")
	prompt = regexp.MustCompile(`(?i)começa com`).ReplaceAllString(prompt, "LIKE")
	prompt = regexp.MustCompile(`(?i)iniciando com`).ReplaceAllString(prompt, "LIKE")
	prompt = regexp.MustCompile(`(?i)inicia com`).ReplaceAllString(prompt, "LIKE")
	prompt = regexp.MustCompile(`(?i)começa por`).ReplaceAllString(prompt, "LIKE")
	prompt = regexp.MustCompile(`(?i)no início`).ReplaceAllString(prompt, "LIKE")
	prompt = regexp.MustCompile(`(?i)termina com`).ReplaceAllString(prompt, "LIKE")
	prompt = regexp.MustCompile(`(?i)terminando com`).ReplaceAllString(prompt, "LIKE")
	prompt = regexp.MustCompile(`(?i)finaliza com`).ReplaceAllString(prompt, "LIKE")
	prompt = regexp.MustCompile(`(?i)acaba com`).ReplaceAllString(prompt, "LIKE")
	prompt = regexp.MustCompile(`(?i)no final`).ReplaceAllString(prompt, "LIKE")
	prompt = regexp.MustCompile(`(?i)ao final`).ReplaceAllString(prompt, "LIKE")
	prompt = regexp.MustCompile(`(?i)que está entre`).ReplaceAllString(prompt, "BETWEEN")
	prompt = regexp.MustCompile(`(?i)entre os valores`).ReplaceAllString(prompt, "BETWEEN")
	prompt = regexp.MustCompile(`(?i)no intervalo`).ReplaceAllString(prompt, "BETWEEN")
	prompt = regexp.MustCompile(`(?i)dentro do intervalo`).ReplaceAllString(prompt, "BETWEEN")
	prompt = regexp.MustCompile(`(?i)na faixa de`).ReplaceAllString(prompt, "BETWEEN")
	prompt = regexp.MustCompile(`(?i)cai (no|entre)`).ReplaceAllString(prompt, "BETWEEN")

	// Add more null condition expressions
	prompt = regexp.MustCompile(`(?i)está vazio`).ReplaceAllString(prompt, "IS NULL")
	prompt = regexp.MustCompile(`(?i)é vazio`).ReplaceAllString(prompt, "IS NULL")
	prompt = regexp.MustCompile(`(?i)não foi preenchido`).ReplaceAllString(prompt, "IS NULL")
	prompt = regexp.MustCompile(`(?i)é nulo`).ReplaceAllString(prompt, "IS NULL")
	prompt = regexp.MustCompile(`(?i)não existe`).ReplaceAllString(prompt, "IS NULL")
	prompt = regexp.MustCompile(`(?i)está preenchido`).ReplaceAllString(prompt, "IS NOT NULL")
	prompt = regexp.MustCompile(`(?i)não está vazio`).ReplaceAllString(prompt, "IS NOT NULL")
	prompt = regexp.MustCompile(`(?i)foi preenchido`).ReplaceAllString(prompt, "IS NOT NULL")
	prompt = regexp.MustCompile(`(?i)não é nulo`).ReplaceAllString(prompt, "IS NOT NULL")
	prompt = regexp.MustCompile(`(?i)existe`).ReplaceAllString(prompt, "IS NOT NULL")
	prompt = regexp.MustCompile(`(?i)está definido`).ReplaceAllString(prompt, "IS NOT NULL")

	// Normalize join expressions
	prompt = regexp.MustCompile(`(?i)junto com seus`).ReplaceAllString(prompt, "juntar")
	prompt = regexp.MustCompile(`(?i)junto com suas`).ReplaceAllString(prompt, "juntar")
	prompt = regexp.MustCompile(`(?i)junto com os`).ReplaceAllString(prompt, "juntar")
	prompt = regexp.MustCompile(`(?i)junto com as`).ReplaceAllString(prompt, "juntar")
	prompt = regexp.MustCompile(`(?i)juntamente com`).ReplaceAllString(prompt, "juntar")
	prompt = regexp.MustCompile(`(?i)combinado com`).ReplaceAllString(prompt, "juntar")
	prompt = regexp.MustCompile(`(?i)relacionado com`).ReplaceAllString(prompt, "juntar")
	prompt = regexp.MustCompile(`(?i)relacionado a`).ReplaceAllString(prompt, "juntar")
	prompt = regexp.MustCompile(`(?i)relacionando`).ReplaceAllString(prompt, "juntar")
	prompt = regexp.MustCompile(`(?i)associado com`).ReplaceAllString(prompt, "juntar")
	prompt = regexp.MustCompile(`(?i)associado a`).ReplaceAllString(prompt, "juntar")
	prompt = regexp.MustCompile(`(?i)associando`).ReplaceAllString(prompt, "juntar")
	prompt = regexp.MustCompile(`(?i)conectado com`).ReplaceAllString(prompt, "juntar")
	prompt = regexp.MustCompile(`(?i)conectado a`).ReplaceAllString(prompt, "juntar")
	prompt = regexp.MustCompile(`(?i)que estão em`).ReplaceAllString(prompt, "juntar")
	prompt = regexp.MustCompile(`(?i)vinculado a`).ReplaceAllString(prompt, "juntar")
	prompt = regexp.MustCompile(`(?i)e seus respectivos`).ReplaceAllString(prompt, "juntar")
	prompt = regexp.MustCompile(`(?i)e suas respectivas`).ReplaceAllString(prompt, "juntar")
	prompt = regexp.MustCompile(`(?i)e todas as`).ReplaceAllString(prompt, "juntar")
	prompt = regexp.MustCompile(`(?i)e todos os`).ReplaceAllString(prompt, "juntar")
	prompt = regexp.MustCompile(`(?i)incluindo seus`).ReplaceAllString(prompt, "juntar")
	prompt = regexp.MustCompile(`(?i)incluindo suas`).ReplaceAllString(prompt, "juntar")
	prompt = regexp.MustCompile(`(?i)fazendo join com`).ReplaceAllString(prompt, "juntar")
	prompt = regexp.MustCompile(`(?i)ligando com`).ReplaceAllString(prompt, "juntar")
	prompt = regexp.MustCompile(`(?i)ligado a`).ReplaceAllString(prompt, "juntar")
	prompt = regexp.MustCompile(`(?i)conectando com`).ReplaceAllString(prompt, "juntar")
	prompt = regexp.MustCompile(`(?i)relacione com`).ReplaceAllString(prompt, "juntar")

	// Identify and normalize time-related expressions
	var ptTimeExpressions = map[string]string{
		"(?i)últimos (\\d+) dias":      "data >= DATE_SUB(CURRENT_DATE(), INTERVAL $1 DAY)",
		"(?i)últimas (\\d+) semanas":   "data >= DATE_SUB(CURRENT_DATE(), INTERVAL $1 WEEK)",
		"(?i)últimos (\\d+) meses":     "data >= DATE_SUB(CURRENT_DATE(), INTERVAL $1 MONTH)",
		"(?i)últimos (\\d+) anos":      "data >= DATE_SUB(CURRENT_DATE(), INTERVAL $1 YEAR)",
		"(?i)nas últimas (\\d+) horas": "data >= DATE_SUB(CURRENT_TIMESTAMP(), INTERVAL $1 HOUR)",
		"(?i)no último (\\d+) dia":     "data >= DATE_SUB(CURRENT_DATE(), INTERVAL $1 DAY)",
		"(?i)no último (\\d+) mês":     "data >= DATE_SUB(CURRENT_DATE(), INTERVAL $1 MONTH)",
		"(?i)no último (\\d+) ano":     "data >= DATE_SUB(CURRENT_DATE(), INTERVAL $1 YEAR)",
		"(?i)nos últimos (\\d+) dias":  "data >= DATE_SUB(CURRENT_DATE(), INTERVAL $1 DAY)",
		"(?i)nos últimos (\\d+) meses": "data >= DATE_SUB(CURRENT_DATE(), INTERVAL $1 MONTH)",
		"(?i)nos últimos (\\d+) anos":  "data >= DATE_SUB(CURRENT_DATE(), INTERVAL $1 YEAR)",
		"(?i)dentro de (\\d+) dias":    "data <= DATE_ADD(CURRENT_DATE(), INTERVAL $1 DAY)",
		"(?i)dentro de (\\d+) semanas": "data <= DATE_ADD(CURRENT_DATE(), INTERVAL $1 WEEK)",
		"(?i)dentro de (\\d+) meses":   "data <= DATE_ADD(CURRENT_DATE(), INTERVAL $1 MONTH)",
		"(?i)dentro de (\\d+) anos":    "data <= DATE_ADD(CURRENT_DATE(), INTERVAL $1 YEAR)",
		"(?i)próximos (\\d+) dias":     "data <= DATE_ADD(CURRENT_DATE(), INTERVAL $1 DAY)",
		"(?i)próximas (\\d+) semanas":  "data <= DATE_ADD(CURRENT_DATE(), INTERVAL $1 WEEK)",
		"(?i)próximos (\\d+) meses":    "data <= DATE_ADD(CURRENT_DATE(), INTERVAL $1 MONTH)",
		"(?i)próximos (\\d+) anos":     "data <= DATE_ADD(CURRENT_DATE(), INTERVAL $1 YEAR)",
		"(?i)há (\\d+) dias atrás":     "data = DATE_SUB(CURRENT_DATE(), INTERVAL $1 DAY)",
		"(?i)há (\\d+) semanas atrás":  "data = DATE_SUB(CURRENT_DATE(), INTERVAL $1 WEEK)",
		"(?i)há (\\d+) meses atrás":    "data = DATE_SUB(CURRENT_DATE(), INTERVAL $1 MONTH)",
		"(?i)há (\\d+) anos atrás":     "data = DATE_SUB(CURRENT_DATE(), INTERVAL $1 YEAR)",
		"(?i)daqui a (\\d+) dias":      "data = DATE_ADD(CURRENT_DATE(), INTERVAL $1 DAY)",
		"(?i)daqui a (\\d+) semanas":   "data = DATE_ADD(CURRENT_DATE(), INTERVAL $1 WEEK)",
		"(?i)daqui a (\\d+) meses":     "data = DATE_ADD(CURRENT_DATE(), INTERVAL $1 MONTH)",
		"(?i)daqui a (\\d+) anos":      "data = DATE_ADD(CURRENT_DATE(), INTERVAL $1 YEAR)",
		"(?i)hoje":                     "data = CURRENT_DATE()",
		"(?i)agora":                    "data = CURRENT_TIMESTAMP()",
		"(?i)dia atual":                "data = CURRENT_DATE()",
		"(?i)data atual":               "data = CURRENT_DATE()",
		"(?i)momento atual":            "data = CURRENT_TIMESTAMP()",
		"(?i)ontem":                    "data = DATE_SUB(CURRENT_DATE(), INTERVAL 1 DAY)",
		"(?i)dia anterior":             "data = DATE_SUB(CURRENT_DATE(), INTERVAL 1 DAY)",
		"(?i)anteontem":                "data = DATE_SUB(CURRENT_DATE(), INTERVAL 2 DAY)",
		"(?i)dois dias atrás":          "data = DATE_SUB(CURRENT_DATE(), INTERVAL 2 DAY)",
		"(?i)amanhã":                   "data = DATE_ADD(CURRENT_DATE(), INTERVAL 1 DAY)",
		"(?i)dia seguinte":             "data = DATE_ADD(CURRENT_DATE(), INTERVAL 1 DAY)",
		"(?i)próximo dia":              "data = DATE_ADD(CURRENT_DATE(), INTERVAL 1 DAY)",
		"(?i)depois de amanhã":         "data = DATE_ADD(CURRENT_DATE(), INTERVAL 2 DAY)",
		"(?i)esta semana":              "data >= DATE_TRUNC('week', CURRENT_DATE())",
		"(?i)semana atual":             "data >= DATE_TRUNC('week', CURRENT_DATE())",
		"(?i)semana corrente":          "data >= DATE_TRUNC('week', CURRENT_DATE())",
		"(?i)este mês":                 "data >= DATE_TRUNC('month', CURRENT_DATE())",
		"(?i)mês atual":                "data >= DATE_TRUNC('month', CURRENT_DATE())",
		"(?i)mês corrente":             "data >= DATE_TRUNC('month', CURRENT_DATE())",
		"(?i)este ano":                 "data >= DATE_TRUNC('year', CURRENT_DATE())",
		"(?i)ano atual":                "data >= DATE_TRUNC('year', CURRENT_DATE())",
		"(?i)ano corrente":             "data >= DATE_TRUNC('year', CURRENT_DATE())",
		"(?i)semana passada":           "data BETWEEN DATE_SUB(DATE_TRUNC('week', CURRENT_DATE()), INTERVAL 1 WEEK) AND DATE_SUB(DATE_TRUNC('week', CURRENT_DATE()), INTERVAL 1 DAY)",
		"(?i)última semana":            "data BETWEEN DATE_SUB(DATE_TRUNC('week', CURRENT_DATE()), INTERVAL 1 WEEK) AND DATE_SUB(DATE_TRUNC('week', CURRENT_DATE()), INTERVAL 1 DAY)",
		"(?i)mês passado":              "data BETWEEN DATE_SUB(DATE_TRUNC('month', CURRENT_DATE()), INTERVAL 1 MONTH) AND DATE_SUB(DATE_TRUNC('month', CURRENT_DATE()), INTERVAL 1 DAY)",
		"(?i)último mês":               "data BETWEEN DATE_SUB(DATE_TRUNC('month', CURRENT_DATE()), INTERVAL 1 MONTH) AND DATE_SUB(DATE_TRUNC('month', CURRENT_DATE()), INTERVAL 1 DAY)",
		"(?i)mês anterior":             "data BETWEEN DATE_SUB(DATE_TRUNC('month', CURRENT_DATE()), INTERVAL 1 MONTH) AND DATE_SUB(DATE_TRUNC('month', CURRENT_DATE()), INTERVAL 1 DAY)",
		"(?i)ano passado":              "data BETWEEN DATE_SUB(DATE_TRUNC('year', CURRENT_DATE()), INTERVAL 1 YEAR) AND DATE_SUB(DATE_TRUNC('year', CURRENT_DATE()), INTERVAL 1 DAY)",
		"(?i)último ano":               "data BETWEEN DATE_SUB(DATE_TRUNC('year', CURRENT_DATE()), INTERVAL 1 YEAR) AND DATE_SUB(DATE_TRUNC('year', CURRENT_DATE()), INTERVAL 1 DAY)",
		"(?i)ano anterior":             "data BETWEEN DATE_SUB(DATE_TRUNC('year', CURRENT_DATE()), INTERVAL 1 YEAR) AND DATE_SUB(DATE_TRUNC('year', CURRENT_DATE()), INTERVAL 1 DAY)",
	}

	for pattern, replacement := range ptTimeExpressions {
		re := regexp.MustCompile(pattern)
		prompt = re.ReplaceAllString(prompt, replacement)
	}

	// Identify and normalize group and order operations
	groupOrderExpressions := map[string]string{
		"(?i)agrupados? por":        "GROUP BY",
		"(?i)agrupar por":           "GROUP BY",
		"(?i)agrupando por":         "GROUP BY",
		"(?i)agrupad[ao]s? em":      "GROUP BY",
		"(?i)agrupamento por":       "GROUP BY",
		"(?i)separados? por":        "GROUP BY",
		"(?i)divididos? por":        "GROUP BY",
		"(?i)categorizad[ao]s? por": "GROUP BY",
		"(?i)classificad[ao]s? por": "GROUP BY",
		"(?i)por grupos? de":        "GROUP BY",
		"(?i)em grupos? de":         "GROUP BY",
		"(?i)por categorias? de":    "GROUP BY",
		"(?i)ordenados? por":        "ORDER BY",
		"(?i)ordenar por":           "ORDER BY",
		"(?i)ordenação por":         "ORDER BY",
		"(?i)ordenando por":         "ORDER BY",
		"(?i)classificado por":      "ORDER BY",
		"(?i)classificar por":       "ORDER BY",
		"(?i)classificando por":     "ORDER BY",
		"(?i)organizado por":        "ORDER BY",
		"(?i)organizar por":         "ORDER BY",
		"(?i)organizando por":       "ORDER BY",
		"(?i)em ordem crescente":    "ORDER BY ASC",
		"(?i)de forma crescente":    "ORDER BY ASC",
		"(?i)crescentemente":        "ORDER BY ASC",
		"(?i)em ordem ascendente":   "ORDER BY ASC",
		"(?i)ascendentemente":       "ORDER BY ASC",
		"(?i)em ordem decrescente":  "ORDER BY DESC",
		"(?i)de forma decrescente":  "ORDER BY DESC",
		"(?i)decrescentemente":      "ORDER BY DESC",
		"(?i)em ordem descendente":  "ORDER BY DESC",
		"(?i)descendentemente":      "ORDER BY DESC",
		"(?i)do menor para o maior": "ORDER BY ASC",
		"(?i)do maior para o menor": "ORDER BY DESC",
		"(?i)do mais recente":       "ORDER BY DESC",
		"(?i)do mais antigo":        "ORDER BY ASC",
		"(?i)alfabeticamente":       "ORDER BY ASC",
		"(?i)iniciando pelo menor":  "ORDER BY ASC",
		"(?i)iniciando pelo maior":  "ORDER BY DESC",
		"(?i)começando pelo menor":  "ORDER BY ASC",
		"(?i)começando pelo maior":  "ORDER BY DESC",
		"(?i)mostrando primeiro":    "ORDER BY ASC",
		"(?i)prioridade para":       "ORDER BY ASC",
		"(?i)pelos mais novos":      "ORDER BY DESC",
		"(?i)pelos mais antigos":    "ORDER BY ASC",
		"(?i)com limite":            "LIMIT",
		"(?i)limitando a":           "LIMIT",
		"(?i)limitar a":             "LIMIT",
		"(?i)limitado a":            "LIMIT",
		"(?i)restringir a":          "LIMIT",
		"(?i)restringido a":         "LIMIT",
		"(?i)primeiros? (\\d+)":     "LIMIT $1",
		"(?i)somente (\\d+)":        "LIMIT $1",
		"(?i)apenas (\\d+)":         "LIMIT $1",
		"(?i)tão somente (\\d+)":    "LIMIT $1",
		"(?i)no máximo (\\d+)":      "LIMIT $1",
		"(?i)até (\\d+)":            "LIMIT $1",
		"(?i)não mais que (\\d+)":   "LIMIT $1",
		"(?i)mostrar (\\d+)":        "LIMIT $1",
		"(?i)exibir (\\d+)":         "LIMIT $1",
		"(?i)trazer (\\d+)":         "LIMIT $1",
		"(?i)listar (\\d+)":         "LIMIT $1",
		"(?i)obter (\\d+)":          "LIMIT $1",
	}

	for pattern, replacement := range groupOrderExpressions {
		re := regexp.MustCompile(pattern)
		prompt = re.ReplaceAllString(prompt, replacement)
	}

	// Normalizações avançadas para contagens e agregações
	prompt = regexp.MustCompile(`(?i)média de`).ReplaceAllString(prompt, "AVG")
	prompt = regexp.MustCompile(`(?i)média dos`).ReplaceAllString(prompt, "AVG")
	prompt = regexp.MustCompile(`(?i)média das`).ReplaceAllString(prompt, "AVG")
	prompt = regexp.MustCompile(`(?i)valor médio`).ReplaceAllString(prompt, "AVG")
	prompt = regexp.MustCompile(`(?i)mediana de`).ReplaceAllString(prompt, "MEDIAN")
	prompt = regexp.MustCompile(`(?i)mediana dos`).ReplaceAllString(prompt, "MEDIAN")
	prompt = regexp.MustCompile(`(?i)mediana das`).ReplaceAllString(prompt, "MEDIAN")
	prompt = regexp.MustCompile(`(?i)valor mediano`).ReplaceAllString(prompt, "MEDIAN")
	prompt = regexp.MustCompile(`(?i)soma de`).ReplaceAllString(prompt, "SUM")
	prompt = regexp.MustCompile(`(?i)soma dos`).ReplaceAllString(prompt, "SUM")
	prompt = regexp.MustCompile(`(?i)soma das`).ReplaceAllString(prompt, "SUM")
	prompt = regexp.MustCompile(`(?i)somatória de`).ReplaceAllString(prompt, "SUM")
	prompt = regexp.MustCompile(`(?i)total dos`).ReplaceAllString(prompt, "SUM")
	prompt = regexp.MustCompile(`(?i)total das`).ReplaceAllString(prompt, "SUM")
	prompt = regexp.MustCompile(`(?i)máximo de`).ReplaceAllString(prompt, "MAX")
	prompt = regexp.MustCompile(`(?i)máximo dos`).ReplaceAllString(prompt, "MAX")
	prompt = regexp.MustCompile(`(?i)máximo das`).ReplaceAllString(prompt, "MAX")
	prompt = regexp.MustCompile(`(?i)maior valor de`).ReplaceAllString(prompt, "MAX")
	prompt = regexp.MustCompile(`(?i)maior valor dos`).ReplaceAllString(prompt, "MAX")
	prompt = regexp.MustCompile(`(?i)maior valor das`).ReplaceAllString(prompt, "MAX")
	prompt = regexp.MustCompile(`(?i)valor máximo`).ReplaceAllString(prompt, "MAX")
	prompt = regexp.MustCompile(`(?i)mínimo de`).ReplaceAllString(prompt, "MIN")
	prompt = regexp.MustCompile(`(?i)mínimo dos`).ReplaceAllString(prompt, "MIN")
	prompt = regexp.MustCompile(`(?i)mínimo das`).ReplaceAllString(prompt, "MIN")
	prompt = regexp.MustCompile(`(?i)menor valor de`).ReplaceAllString(prompt, "MIN")
	prompt = regexp.MustCompile(`(?i)menor valor dos`).ReplaceAllString(prompt, "MIN")
	prompt = regexp.MustCompile(`(?i)menor valor das`).ReplaceAllString(prompt, "MIN")
	prompt = regexp.MustCompile(`(?i)valor mínimo`).ReplaceAllString(prompt, "MIN")
	prompt = regexp.MustCompile(`(?i)conte quantos`).ReplaceAllString(prompt, "COUNT")
	prompt = regexp.MustCompile(`(?i)conte as`).ReplaceAllString(prompt, "COUNT")
	prompt = regexp.MustCompile(`(?i)conte os`).ReplaceAllString(prompt, "COUNT")
	prompt = regexp.MustCompile(`(?i)contagem de`).ReplaceAllString(prompt, "COUNT")
	prompt = regexp.MustCompile(`(?i)contagem dos`).ReplaceAllString(prompt, "COUNT")
	prompt = regexp.MustCompile(`(?i)contagem das`).ReplaceAllString(prompt, "COUNT")
	prompt = regexp.MustCompile(`(?i)quantidade total`).ReplaceAllString(prompt, "COUNT")
	prompt = regexp.MustCompile(`(?i)número total`).ReplaceAllString(prompt, "COUNT")

	// Expressões avançadas para junções
	prompt = regexp.MustCompile(`(?i)cruzar dados com`).ReplaceAllString(prompt, "juntar")
	prompt = regexp.MustCompile(`(?i)cruzar com`).ReplaceAllString(prompt, "juntar")
	prompt = regexp.MustCompile(`(?i)unir com`).ReplaceAllString(prompt, "juntar")
	prompt = regexp.MustCompile(`(?i)unindo com`).ReplaceAllString(prompt, "juntar")
	prompt = regexp.MustCompile(`(?i)fazer junção com`).ReplaceAllString(prompt, "juntar")
	prompt = regexp.MustCompile(`(?i)relacionando a`).ReplaceAllString(prompt, "juntar")
	prompt = regexp.MustCompile(`(?i)fazer relacionamento com`).ReplaceAllString(prompt, "juntar")
	prompt = regexp.MustCompile(`(?i)juntamente aos`).ReplaceAllString(prompt, "juntar")
	prompt = regexp.MustCompile(`(?i)juntamente às`).ReplaceAllString(prompt, "juntar")
	prompt = regexp.MustCompile(`(?i)correlacionar com`).ReplaceAllString(prompt, "juntar")
	prompt = regexp.MustCompile(`(?i)correlacionado com`).ReplaceAllString(prompt, "juntar")
	prompt = regexp.MustCompile(`(?i)combinar dados de`).ReplaceAllString(prompt, "juntar")
	prompt = regexp.MustCompile(`(?i)trazer dados relacionados de`).ReplaceAllString(prompt, "juntar")
	prompt = regexp.MustCompile(`(?i)fazer inner join com`).ReplaceAllString(prompt, "INNER JOIN")
	prompt = regexp.MustCompile(`(?i)fazer left join com`).ReplaceAllString(prompt, "LEFT JOIN")
	prompt = regexp.MustCompile(`(?i)fazer right join com`).ReplaceAllString(prompt, "RIGHT JOIN")
	prompt = regexp.MustCompile(`(?i)fazer full join com`).ReplaceAllString(prompt, "FULL JOIN")
	prompt = regexp.MustCompile(`(?i)mesmo que não tenha`).ReplaceAllString(prompt, "LEFT JOIN")
	prompt = regexp.MustCompile(`(?i)mesmo sem correspondência`).ReplaceAllString(prompt, "LEFT JOIN")
	prompt = regexp.MustCompile(`(?i)incluindo registros sem correspondência`).ReplaceAllString(prompt, "LEFT JOIN")
	prompt = regexp.MustCompile(`(?i)mesmo que não exista`).ReplaceAllString(prompt, "LEFT JOIN")

	// Expressões para filtragem avançada
	prompt = regexp.MustCompile(`(?i)existam registros em`).ReplaceAllString(prompt, "EXISTS")
	prompt = regexp.MustCompile(`(?i)exista algum em`).ReplaceAllString(prompt, "EXISTS")
	prompt = regexp.MustCompile(`(?i)tenha algum registro em`).ReplaceAllString(prompt, "EXISTS")
	prompt = regexp.MustCompile(`(?i)não existam registros em`).ReplaceAllString(prompt, "NOT EXISTS")
	prompt = regexp.MustCompile(`(?i)não exista nenhum em`).ReplaceAllString(prompt, "NOT EXISTS")
	prompt = regexp.MustCompile(`(?i)não tenha nenhum registro em`).ReplaceAllString(prompt, "NOT EXISTS")
	prompt = regexp.MustCompile(`(?i)pertença ao conjunto`).ReplaceAllString(prompt, "IN")
	prompt = regexp.MustCompile(`(?i)esteja na lista`).ReplaceAllString(prompt, "IN")
	prompt = regexp.MustCompile(`(?i)faça parte de`).ReplaceAllString(prompt, "IN")
	prompt = regexp.MustCompile(`(?i)seja um dos`).ReplaceAllString(prompt, "IN")
	prompt = regexp.MustCompile(`(?i)seja uma das`).ReplaceAllString(prompt, "IN")
	prompt = regexp.MustCompile(`(?i)encontre-se em`).ReplaceAllString(prompt, "IN")
	prompt = regexp.MustCompile(`(?i)não pertença ao conjunto`).ReplaceAllString(prompt, "NOT IN")
	prompt = regexp.MustCompile(`(?i)não esteja na lista`).ReplaceAllString(prompt, "NOT IN")
	prompt = regexp.MustCompile(`(?i)não faça parte de`).ReplaceAllString(prompt, "NOT IN")
	prompt = regexp.MustCompile(`(?i)não seja um dos`).ReplaceAllString(prompt, "NOT IN")
	prompt = regexp.MustCompile(`(?i)não seja uma das`).ReplaceAllString(prompt, "NOT IN")
	prompt = regexp.MustCompile(`(?i)não encontre-se em`).ReplaceAllString(prompt, "NOT IN")
	prompt = regexp.MustCompile(`(?i)caso em que`).ReplaceAllString(prompt, "CASE WHEN")
	prompt = regexp.MustCompile(`(?i)nos casos em que`).ReplaceAllString(prompt, "CASE WHEN")
	prompt = regexp.MustCompile(`(?i)na situação onde`).ReplaceAllString(prompt, "CASE WHEN")
	prompt = regexp.MustCompile(`(?i)quando for`).ReplaceAllString(prompt, "CASE WHEN")
	prompt = regexp.MustCompile(`(?i)caso contrário`).ReplaceAllString(prompt, "ELSE")
	prompt = regexp.MustCompile(`(?i)do contrário`).ReplaceAllString(prompt, "ELSE")
	prompt = regexp.MustCompile(`(?i)se não`).ReplaceAllString(prompt, "ELSE")
	prompt = regexp.MustCompile(`(?i)senão`).ReplaceAllString(prompt, "ELSE")

	// Expressões adicionais para datas e períodos
	prompt = regexp.MustCompile(`(?i)primeiro dia do mês`).ReplaceAllString(prompt, "DATE_TRUNC('month', CURRENT_DATE())")
	prompt = regexp.MustCompile(`(?i)último dia do mês`).ReplaceAllString(prompt, "LAST_DAY(CURRENT_DATE())")
	prompt = regexp.MustCompile(`(?i)início do mês`).ReplaceAllString(prompt, "DATE_TRUNC('month', CURRENT_DATE())")
	prompt = regexp.MustCompile(`(?i)final do mês`).ReplaceAllString(prompt, "LAST_DAY(CURRENT_DATE())")
	prompt = regexp.MustCompile(`(?i)começo do mês`).ReplaceAllString(prompt, "DATE_TRUNC('month', CURRENT_DATE())")
	prompt = regexp.MustCompile(`(?i)fim do mês`).ReplaceAllString(prompt, "LAST_DAY(CURRENT_DATE())")
	prompt = regexp.MustCompile(`(?i)início do ano`).ReplaceAllString(prompt, "DATE_TRUNC('year', CURRENT_DATE())")
	prompt = regexp.MustCompile(`(?i)final do ano`).ReplaceAllString(prompt, "DATE_ADD(DATE_TRUNC('year', CURRENT_DATE()), INTERVAL 1 YEAR - 1 DAY)")
	prompt = regexp.MustCompile(`(?i)começo do ano`).ReplaceAllString(prompt, "DATE_TRUNC('year', CURRENT_DATE())")
	prompt = regexp.MustCompile(`(?i)fim do ano`).ReplaceAllString(prompt, "DATE_ADD(DATE_TRUNC('year', CURRENT_DATE()), INTERVAL 1 YEAR - 1 DAY)")
	prompt = regexp.MustCompile(`(?i)início da semana`).ReplaceAllString(prompt, "DATE_TRUNC('week', CURRENT_DATE())")
	prompt = regexp.MustCompile(`(?i)final da semana`).ReplaceAllString(prompt, "DATE_ADD(DATE_TRUNC('week', CURRENT_DATE()), INTERVAL 6 DAY)")
	prompt = regexp.MustCompile(`(?i)começo da semana`).ReplaceAllString(prompt, "DATE_TRUNC('week', CURRENT_DATE())")
	prompt = regexp.MustCompile(`(?i)fim da semana`).ReplaceAllString(prompt, "DATE_ADD(DATE_TRUNC('week', CURRENT_DATE()), INTERVAL 6 DAY)")
	prompt = regexp.MustCompile(`(?i)primeiro trimestre`).ReplaceAllString(prompt, "data BETWEEN DATE_TRUNC('year', CURRENT_DATE()) AND DATE_ADD(DATE_TRUNC('year', CURRENT_DATE()), INTERVAL 3 MONTH - 1 DAY)")
	prompt = regexp.MustCompile(`(?i)segundo trimestre`).ReplaceAllString(prompt, "data BETWEEN DATE_ADD(DATE_TRUNC('year', CURRENT_DATE()), INTERVAL 3 MONTH) AND DATE_ADD(DATE_TRUNC('year', CURRENT_DATE()), INTERVAL 6 MONTH - 1 DAY)")
	prompt = regexp.MustCompile(`(?i)terceiro trimestre`).ReplaceAllString(prompt, "data BETWEEN DATE_ADD(DATE_TRUNC('year', CURRENT_DATE()), INTERVAL 6 MONTH) AND DATE_ADD(DATE_TRUNC('year', CURRENT_DATE()), INTERVAL 9 MONTH - 1 DAY)")
	prompt = regexp.MustCompile(`(?i)quarto trimestre`).ReplaceAllString(prompt, "data BETWEEN DATE_ADD(DATE_TRUNC('year', CURRENT_DATE()), INTERVAL 9 MONTH) AND DATE_ADD(DATE_TRUNC('year', CURRENT_DATE()), INTERVAL 1 YEAR - 1 DAY)")
	prompt = regexp.MustCompile(`(?i)primeiro semestre`).ReplaceAllString(prompt, "data BETWEEN DATE_TRUNC('year', CURRENT_DATE()) AND DATE_ADD(DATE_TRUNC('year', CURRENT_DATE()), INTERVAL 6 MONTH - 1 DAY)")
	prompt = regexp.MustCompile(`(?i)segundo semestre`).ReplaceAllString(prompt, "data BETWEEN DATE_ADD(DATE_TRUNC('year', CURRENT_DATE()), INTERVAL 6 MONTH) AND DATE_ADD(DATE_TRUNC('year', CURRENT_DATE()), INTERVAL 1 YEAR - 1 DAY)")
	prompt = regexp.MustCompile(`(?i)janeiro`).ReplaceAllString(prompt, "MONTH(data) = 1")
	prompt = regexp.MustCompile(`(?i)fevereiro`).ReplaceAllString(prompt, "MONTH(data) = 2")
	prompt = regexp.MustCompile(`(?i)março`).ReplaceAllString(prompt, "MONTH(data) = 3")
	prompt = regexp.MustCompile(`(?i)abril`).ReplaceAllString(prompt, "MONTH(data) = 4")
	prompt = regexp.MustCompile(`(?i)maio`).ReplaceAllString(prompt, "MONTH(data) = 5")
	prompt = regexp.MustCompile(`(?i)junho`).ReplaceAllString(prompt, "MONTH(data) = 6")
	prompt = regexp.MustCompile(`(?i)julho`).ReplaceAllString(prompt, "MONTH(data) = 7")
	prompt = regexp.MustCompile(`(?i)agosto`).ReplaceAllString(prompt, "MONTH(data) = 8")
	prompt = regexp.MustCompile(`(?i)setembro`).ReplaceAllString(prompt, "MONTH(data) = 9")
	prompt = regexp.MustCompile(`(?i)outubro`).ReplaceAllString(prompt, "MONTH(data) = 10")
	prompt = regexp.MustCompile(`(?i)novembro`).ReplaceAllString(prompt, "MONTH(data) = 11")
	prompt = regexp.MustCompile(`(?i)dezembro`).ReplaceAllString(prompt, "MONTH(data) = 12")
	prompt = regexp.MustCompile(`(?i)último (\d+) dias úteis`).ReplaceAllString(prompt, "data >= DATE_SUB(CURRENT_DATE(), INTERVAL $1 * 7/5 DAY)")
	prompt = regexp.MustCompile(`(?i)próximos (\d+) dias úteis`).ReplaceAllString(prompt, "data <= DATE_ADD(CURRENT_DATE(), INTERVAL $1 * 7/5 DAY)")
	prompt = regexp.MustCompile(`(?i)dia útil`).ReplaceAllString(prompt, "DAYOFWEEK(data) BETWEEN 2 AND 6")
	prompt = regexp.MustCompile(`(?i)dias úteis`).ReplaceAllString(prompt, "DAYOFWEEK(data) BETWEEN 2 AND 6")
	prompt = regexp.MustCompile(`(?i)final de semana`).ReplaceAllString(prompt, "DAYOFWEEK(data) IN (1, 7)")
	prompt = regexp.MustCompile(`(?i)fim de semana`).ReplaceAllString(prompt, "DAYOFWEEK(data) IN (1, 7)")
	prompt = regexp.MustCompile(`(?i)segunda[- ]feira`).ReplaceAllString(prompt, "DAYOFWEEK(data) = 2")
	prompt = regexp.MustCompile(`(?i)terça[- ]feira`).ReplaceAllString(prompt, "DAYOFWEEK(data) = 3")
	prompt = regexp.MustCompile(`(?i)quarta[- ]feira`).ReplaceAllString(prompt, "DAYOFWEEK(data) = 4")
	prompt = regexp.MustCompile(`(?i)quinta[- ]feira`).ReplaceAllString(prompt, "DAYOFWEEK(data) = 5")
	prompt = regexp.MustCompile(`(?i)sexta[- ]feira`).ReplaceAllString(prompt, "DAYOFWEEK(data) = 6")
	prompt = regexp.MustCompile(`(?i)sábado`).ReplaceAllString(prompt, "DAYOFWEEK(data) = 7")
	prompt = regexp.MustCompile(`(?i)domingo`).ReplaceAllString(prompt, "DAYOFWEEK(data) = 1")
	prompt = regexp.MustCompile(`(?i)data de (hoje|agora)`).ReplaceAllString(prompt, "data = CURRENT_DATE()")
	prompt = regexp.MustCompile(`(?i)hora (atual|corrente)`).ReplaceAllString(prompt, "HOUR(CURRENT_TIMESTAMP())")
	prompt = regexp.MustCompile(`(?i)minuto (atual|corrente)`).ReplaceAllString(prompt, "MINUTE(CURRENT_TIMESTAMP())")

	// Adicionar padronizações mais específicas para estrutura "onde coluna igual a valor"
	// Isso ajudará a identificar condições mais facilmente
	prompt = regexp.MustCompile(`(?i)onde\s+([a-zA-Z0-9_\.]+)\s+é\s+igual\s+a\s+(['"]?)([^'"]+)(['"]?)`).
		ReplaceAllString(prompt, "onde $1 = $2$3$4")
	prompt = regexp.MustCompile(`(?i)onde\s+([a-zA-Z0-9_\.]+)\s+é\s+(['"]?)([^'"]+)(['"]?)`).
		ReplaceAllString(prompt, "onde $1 = $2$3$4")
	prompt = regexp.MustCompile(`(?i)onde\s+([a-zA-Z0-9_\.]+)\s+igual\s+a\s+(['"]?)([^'"]+)(['"]?)`).
		ReplaceAllString(prompt, "onde $1 = $2$3$4")
	prompt = regexp.MustCompile(`(?i)onde\s+([a-zA-Z0-9_\.]+)\s+vale\s+(['"]?)([^'"]+)(['"]?)`).
		ReplaceAllString(prompt, "onde $1 = $2$3$4")
	prompt = regexp.MustCompile(`(?i)onde\s+(?:a|o)\s+(?:coluna|campo)\s+([a-zA-Z0-9_\.]+)\s+(?:é|vale|contém)\s+(['"]?)([^'"]+)(['"]?)`).
		ReplaceAllString(prompt, "onde $1 = $2$3$4")

	// Padronizar outros operadores de comparação
	prompt = regexp.MustCompile(`(?i)onde\s+([a-zA-Z0-9_\.]+)\s+(?:é\s+)?maior\s+que\s+(['"]?)([^'"]+)(['"]?)`).
		ReplaceAllString(prompt, "onde $1 > $2$3$4")
	prompt = regexp.MustCompile(`(?i)onde\s+([a-zA-Z0-9_\.]+)\s+(?:é\s+)?maior\s+do\s+que\s+(['"]?)([^'"]+)(['"]?)`).
		ReplaceAllString(prompt, "onde $1 > $2$3$4")
	prompt = regexp.MustCompile(`(?i)onde\s+([a-zA-Z0-9_\.]+)\s+(?:é\s+)?maior\s+ou\s+igual\s+(?:a|que)\s+(['"]?)([^'"]+)(['"]?)`).
		ReplaceAllString(prompt, "onde $1 >= $2$3$4")
	prompt = regexp.MustCompile(`(?i)onde\s+([a-zA-Z0-9_\.]+)\s+(?:é\s+)?menor\s+que\s+(['"]?)([^'"]+)(['"]?)`).
		ReplaceAllString(prompt, "onde $1 < $2$3$4")
	prompt = regexp.MustCompile(`(?i)onde\s+([a-zA-Z0-9_\.]+)\s+(?:é\s+)?menor\s+do\s+que\s+(['"]?)([^'"]+)(['"]?)`).
		ReplaceAllString(prompt, "onde $1 < $2$3$4")
	prompt = regexp.MustCompile(`(?i)onde\s+([a-zA-Z0-9_\.]+)\s+(?:é\s+)?menor\s+ou\s+igual\s+(?:a|que)\s+(['"]?)([^'"]+)(['"]?)`).
		ReplaceAllString(prompt, "onde $1 <= $2$3$4")
	prompt = regexp.MustCompile(`(?i)onde\s+([a-zA-Z0-9_\.]+)\s+(?:é\s+)?diferente\s+de\s+(['"]?)([^'"]+)(['"]?)`).
		ReplaceAllString(prompt, "onde $1 != $2$3$4")
	prompt = regexp.MustCompile(`(?i)onde\s+([a-zA-Z0-9_\.]+)\s+(?:é\s+)?não\s+é\s+(['"]?)([^'"]+)(['"]?)`).
		ReplaceAllString(prompt, "onde $1 != $2$3$4")

	// Padronizar expressões com "for"
	prompt = regexp.MustCompile(`(?i)onde\s+([a-zA-Z0-9_\.]+)\s+for\s+(['"]?)([^'"]+)(['"]?)`).
		ReplaceAllString(prompt, "onde $1 = $2$3$4")
	prompt = regexp.MustCompile(`(?i)onde\s+([a-zA-Z0-9_\.]+)\s+for\s+(?:maior|superior)\s+(?:a|que)\s+(['"]?)([^'"]+)(['"]?)`).
		ReplaceAllString(prompt, "onde $1 > $2$3$4")
	prompt = regexp.MustCompile(`(?i)onde\s+([a-zA-Z0-9_\.]+)\s+for\s+(?:maior|superior)\s+ou\s+igual\s+(?:a|que)\s+(['"]?)([^'"]+)(['"]?)`).
		ReplaceAllString(prompt, "onde $1 >= $2$3$4")
	prompt = regexp.MustCompile(`(?i)onde\s+([a-zA-Z0-9_\.]+)\s+for\s+(?:menor|inferior)\s+(?:a|que)\s+(['"]?)([^'"]+)(['"]?)`).
		ReplaceAllString(prompt, "onde $1 < $2$3$4")
	prompt = regexp.MustCompile(`(?i)onde\s+([a-zA-Z0-9_\.]+)\s+for\s+(?:menor|inferior)\s+ou\s+igual\s+(?:a|que)\s+(['"]?)([^'"]+)(['"]?)`).
		ReplaceAllString(prompt, "onde $1 <= $2$3$4")

	// Padronizar coluna a direita do operador
	prompt = regexp.MustCompile(`(?i)onde\s+(?:a|o)\s+(?:coluna|campo)\s+([a-zA-Z0-9_\.]+)\s+(?:é\s+)?maior\s+que\s+(['"]?)([^'"]+)(['"]?)`).
		ReplaceAllString(prompt, "onde $1 > $2$3$4")
	prompt = regexp.MustCompile(`(?i)onde\s+(?:a|o)\s+(?:coluna|campo)\s+([a-zA-Z0-9_\.]+)\s+(?:é\s+)?maior\s+ou\s+igual\s+(?:a|que)\s+(['"]?)([^'"]+)(['"]?)`).
		ReplaceAllString(prompt, "onde $1 >= $2$3$4")
	prompt = regexp.MustCompile(`(?i)onde\s+(?:a|o)\s+(?:coluna|campo)\s+([a-zA-Z0-9_\.]+)\s+(?:é\s+)?menor\s+que\s+(['"]?)([^'"]+)(['"]?)`).
		ReplaceAllString(prompt, "onde $1 < $2$3$4")
	prompt = regexp.MustCompile(`(?i)onde\s+(?:a|o)\s+(?:coluna|campo)\s+([a-zA-Z0-9_\.]+)\s+(?:é\s+)?menor\s+ou\s+igual\s+(?:a|que)\s+(['"]?)([^'"]+)(['"]?)`).
		ReplaceAllString(prompt, "onde $1 <= $2$3$4")

	// Padronizar menções a tabelas e colunas
	prompt = regexp.MustCompile(`(?i)da\s+tabela\s+([a-zA-Z0-9_]+)`).
		ReplaceAllString(prompt, "da tabela $1")
	prompt = regexp.MustCompile(`(?i)na\s+tabela\s+([a-zA-Z0-9_]+)`).
		ReplaceAllString(prompt, "na tabela $1")
	prompt = regexp.MustCompile(`(?i)para\s+a\s+tabela\s+([a-zA-Z0-9_]+)`).
		ReplaceAllString(prompt, "para a tabela $1")
	prompt = regexp.MustCompile(`(?i)dados\s+da\s+tabela\s+([a-zA-Z0-9_]+)`).
		ReplaceAllString(prompt, "dados da tabela $1")
	prompt = regexp.MustCompile(`(?i)registros\s+da\s+tabela\s+([a-zA-Z0-9_]+)`).
		ReplaceAllString(prompt, "registros da tabela $1")
	prompt = regexp.MustCompile(`(?i)entidade\s+([a-zA-Z0-9_]+)`).
		ReplaceAllString(prompt, "tabela $1")
	prompt = regexp.MustCompile(`(?i)entidades\s+([a-zA-Z0-9_]+)`).
		ReplaceAllString(prompt, "tabela $1")
	prompt = regexp.MustCompile(`(?i)dentro\s+da\s+tabela\s+([a-zA-Z0-9_]+)`).
		ReplaceAllString(prompt, "na tabela $1")

	prompt = regexp.MustCompile(`(?i)a\s+coluna\s+([a-zA-Z0-9_]+)`).
		ReplaceAllString(prompt, "a coluna $1")
	prompt = regexp.MustCompile(`(?i)o\s+campo\s+([a-zA-Z0-9_]+)`).
		ReplaceAllString(prompt, "a coluna $1")
	prompt = regexp.MustCompile(`(?i)o\s+atributo\s+([a-zA-Z0-9_]+)`).
		ReplaceAllString(prompt, "a coluna $1")
	prompt = regexp.MustCompile(`(?i)os\s+dados\s+de\s+([a-zA-Z0-9_]+)`).
		ReplaceAllString(prompt, "a coluna $1")
	prompt = regexp.MustCompile(`(?i)as\s+informações\s+de\s+([a-zA-Z0-9_]+)`).
		ReplaceAllString(prompt, "a coluna $1")
	prompt = regexp.MustCompile(`(?i)o\s+valor\s+de\s+([a-zA-Z0-9_]+)`).
		ReplaceAllString(prompt, "a coluna $1")
	prompt = regexp.MustCompile(`(?i)a\s+propriedade\s+([a-zA-Z0-9_]+)`).
		ReplaceAllString(prompt, "a coluna $1")

	// Padronizar referências a coluna.tabela e tabela.coluna
	prompt = regexp.MustCompile(`(?i)(?:a|o)\s+(?:coluna|campo)\s+([a-zA-Z0-9_]+)\s+da\s+tabela\s+([a-zA-Z0-9_]+)`).
		ReplaceAllString(prompt, "a coluna $2.$1")
	prompt = regexp.MustCompile(`(?i)(?:a|o)\s+(?:coluna|campo)\s+([a-zA-Z0-9_]+)\s+na\s+tabela\s+([a-zA-Z0-9_]+)`).
		ReplaceAllString(prompt, "a coluna $2.$1")
	prompt = regexp.MustCompile(`(?i)(?:a|o)\s+(?:coluna|campo)\s+([a-zA-Z0-9_]+)\s+em\s+([a-zA-Z0-9_]+)`).
		ReplaceAllString(prompt, "a coluna $2.$1")
	prompt = regexp.MustCompile(`(?i)(?:a|o)\s+(?:coluna|campo|atributo)\s+([a-zA-Z0-9_]+)\s+(?:pertencente|pertence)\s+(?:a|à)\s+([a-zA-Z0-9_]+)`).
		ReplaceAllString(prompt, "a coluna $2.$1")
	prompt = regexp.MustCompile(`(?i)(?:os|as)\s+(?:dados|informações)\s+de\s+([a-zA-Z0-9_]+)\s+(?:da|na)\s+(?:tabela)\s+([a-zA-Z0-9_]+)`).
		ReplaceAllString(prompt, "a coluna $2.$1")
	prompt = regexp.MustCompile(`(?i)(?:a|o)\s+([a-zA-Z0-9_]+)\s+(?:da|do)\s+([a-zA-Z0-9_]+)`).
		ReplaceAllString(prompt, "a coluna $2.$1")

	// Melhorar padronização para operadores de comparação
	// Usar grupo de captura para garantir que "é igual a valor" seja normalizado para "= valor"
	prompt = regexp.MustCompile(`(?i)\s+é\s+igual\s+a\s+(['"]?)([^'"]+)(['"]?)`).
		ReplaceAllString(prompt, " = $1$2$3")
	prompt = regexp.MustCompile(`(?i)\s+é\s+igual\s+(['"]?)([^'"]+)(['"]?)`).
		ReplaceAllString(prompt, " = $1$2$3")
	prompt = regexp.MustCompile(`(?i)\s+igual\s+a\s+(['"]?)([^'"]+)(['"]?)`).
		ReplaceAllString(prompt, " = $1$2$3")
	prompt = regexp.MustCompile(`(?i)\s+vale\s+(['"]?)([^'"]+)(['"]?)`).
		ReplaceAllString(prompt, " = $1$2$3")

	// Padronizar "onde coluna = valor" para garantir que o operador esteja correto
	prompt = regexp.MustCompile(`(?i)onde\s+(?:a\s+)?(?:coluna\s+)?([a-zA-Z0-9_\.]+)\s+é\s+(['"]?)([^'"]+)(['"]?)`).
		ReplaceAllString(prompt, "onde $1 = $2$3$4")

	return prompt
}

// enhanceTableColumnIdentifiers processa o prompt para destacar referências explícitas a tabelas e colunas
func (p *PromptAnalyzer) enhanceTableColumnIdentifiers(prompt string) string {
	// Expressões para identificar tabelas mencionadas explicitamente
	tableExpressions := map[string][]string{
		"pt": {
			`\b(tabela|table)\s+["']?([a-zA-Z0-9_]+)["']?`,
			`\b(na|da|na|das|dos|nas|nos) tabelas?\s+["']?([a-zA-Z0-9_]+)["']?`,
			`\btabelas?\s+(de|dos|das)?\s+["']?([a-zA-Z0-9_]+)["']?`,
			`\bdados\s+(de|da|do)\s+["']?([a-zA-Z0-9_]+)["']?`,
			`\bregistros?\s+(de|da|do)\s+["']?([a-zA-Z0-9_]+)["']?`,
			`\bentidades?\s+(de|da|do)?\s+["']?([a-zA-Z0-9_]+)["']?`,
			`\bpara\s+(a|o)\s+([a-zA-Z0-9_]+)\b`,
			`\bentidades?\s+["']?([a-zA-Z0-9_]+)["']?`,
			`\bobtenha\s+(?:dados|registros|informações)\s+(?:de|da|do)\s+["']?([a-zA-Z0-9_]+)["']?`,
			`\bselecione\s+(?:da|na|de|do)\s+["']?([a-zA-Z0-9_]+)["']?`,
			`\b(consulte|liste|mostre|exiba)\s+(?:a|o)?\s+["']?([a-zA-Z0-9_]+)["']?`,
			`\b(cadastro|base|banco)\s+(?:de|da|do)?\s+["']?([a-zA-Z0-9_]+)["']?`,
		},
		"en": {
			`\b(table|tables)\s+["']?([a-zA-Z0-9_]+)["']?`,
			`\b(in|from|to|on|of|the) tables?\s+["']?([a-zA-Z0-9_]+)["']?`,
			`\btables?\s+(of|for|with|containing)?\s+["']?([a-zA-Z0-9_]+)["']?`,
			`\bdata\s+(from|of|in)\s+["']?([a-zA-Z0-9_]+)["']?`,
			`\brecords?\s+(from|of|in)\s+["']?([a-zA-Z0-9_]+)["']?`,
			`\bentit(y|ies)\s+(of|named)?\s+["']?([a-zA-Z0-9_]+)["']?`,
			`\bfor\s+(?:the)?\s+([a-zA-Z0-9_]+)\b`,
			`\bentit(?:y|ies)\s+["']?([a-zA-Z0-9_]+)["']?`,
			`\bget\s+(?:data|records|information)\s+(?:from|of)\s+["']?([a-zA-Z0-9_]+)["']?`,
			`\bselect\s+(?:from|in)\s+["']?([a-zA-Z0-9_]+)["']?`,
			`\b(query|list|show|display)\s+(?:the)?\s+["']?([a-zA-Z0-9_]+)["']?`,
			`\b(catalog|database|store)\s+(?:of)?\s+["']?([a-zA-Z0-9_]+)["']?`,
		},
	}

	// Expressões para identificar colunas mencionadas explicitamente
	columnExpressions := map[string][]string{
		"pt": {
			`\b(coluna|campo|column|field|atributo)\s+["']?([a-zA-Z0-9_]+)["']?`,
			`\b(na|da|do|das|dos|nas|nos) colunas?\s+["']?([a-zA-Z0-9_]+)["']?`,
			`\bcolunas?\s+(de|da|do)?\s+["']?([a-zA-Z0-9_]+)["']?`,
			`\bcampos?\s+(de|da|do)?\s+["']?([a-zA-Z0-9_]+)["']?`,
			`\batributos?\s+(de|da|do)?\s+["']?([a-zA-Z0-9_]+)["']?`,
			`\bpropriedades?\s+(de|da|do)?\s+["']?([a-zA-Z0-9_]+)["']?`,
			`\bvalores?\s+(de|da|do|para)\s+["']?([a-zA-Z0-9_]+)["']?`,
			`\binformações?\s+(sobre|de|do|da)\s+["']?([a-zA-Z0-9_]+)["']?`,
			`\bdados\s+(de|da|do)\s+["']?([a-zA-Z0-9_]+)["']?`,
			`\bitem\s+["']?([a-zA-Z0-9_]+)["']?`,
			`\bvariável\s+["']?([a-zA-Z0-9_]+)["']?`,
			`\bentrada\s+["']?([a-zA-Z0-9_]+)["']?`,
			`\bparâmetro\s+["']?([a-zA-Z0-9_]+)["']?`,
			`\b(?:o|a)\s+["']?([a-zA-Z0-9_]+)["']?\s+(?:da|do)\s+(?:tabela|entidade)\s+`,
			`\bonde\s+["']?([a-zA-Z0-9_]+)["']?\s+(?:é|igual|=)`,
			`\bcategorizado\s+por\s+["']?([a-zA-Z0-9_]+)["']?`,
			`\bagrupado\s+por\s+["']?([a-zA-Z0-9_]+)["']?`,
			`\bordenado\s+por\s+["']?([a-zA-Z0-9_]+)["']?`,
		},
		"en": {
			`\b(column|field|attribute|property)\s+["']?([a-zA-Z0-9_]+)["']?`,
			`\b(in|from|with|having|on|of|the) columns?\s+["']?([a-zA-Z0-9_]+)["']?`,
			`\bcolumns?\s+(of|named|called)?\s+["']?([a-zA-Z0-9_]+)["']?`,
			`\bfields?\s+(of|named|called)?\s+["']?([a-zA-Z0-9_]+)["']?`,
			`\battributes?\s+(of|named|called)?\s+["']?([a-zA-Z0-9_]+)["']?`,
			`\bproperties?\s+(of|named|called)?\s+["']?([a-zA-Z0-9_]+)["']?`,
			`\bvalues?\s+(of|for|in)\s+["']?([a-zA-Z0-9_]+)["']?`,
			`\binformation\s+(about|of|for)\s+["']?([a-zA-Z0-9_]+)["']?`,
			`\bdata\s+(of|for)\s+["']?([a-zA-Z0-9_]+)["']?`,
			`\bitem\s+["']?([a-zA-Z0-9_]+)["']?`,
			`\bvariable\s+["']?([a-zA-Z0-9_]+)["']?`,
			`\bentry\s+["']?([a-zA-Z0-9_]+)["']?`,
			`\bparameter\s+["']?([a-zA-Z0-9_]+)["']?`,
			`\bthe\s+["']?([a-zA-Z0-9_]+)["']?\s+(?:of|from)\s+(?:table|entity)\s+`,
			`\bwhere\s+["']?([a-zA-Z0-9_]+)["']?\s+(?:is|equals|=)`,
			`\bcategorized\s+by\s+["']?([a-zA-Z0-9_]+)["']?`,
			`\bgrouped\s+by\s+["']?([a-zA-Z0-9_]+)["']?`,
			`\bordered\s+by\s+["']?([a-zA-Z0-9_]+)["']?`,
		},
	}

	// Expressões para identificar relações de tabela.coluna
	qualifiedColumnExpressions := map[string][]string{
		"pt": {
			`\b([a-zA-Z0-9_]+)\.([a-zA-Z0-9_]+)\b`,
			`\bcolunas?\s+([a-zA-Z0-9_]+)\.([a-zA-Z0-9_]+)\b`,
			`\bcampos?\s+([a-zA-Z0-9_]+)\.([a-zA-Z0-9_]+)\b`,
		},
		"en": {
			`\b([a-zA-Z0-9_]+)\.([a-zA-Z0-9_]+)\b`,
			`\bcolumns?\s+([a-zA-Z0-9_]+)\.([a-zA-Z0-9_]+)\b`,
			`\bfields?\s+([a-zA-Z0-9_]+)\.([a-zA-Z0-9_]+)\b`,
		},
	}

	// Aplicar expressões para o idioma atual
	lang := p.language
	if _, exists := tableExpressions[lang]; !exists {
		lang = "en" // Fallback para inglês
	}

	// Marcar tabelas no prompt
	for _, expr := range tableExpressions[lang] {
		re := regexp.MustCompile(expr)
		prompt = re.ReplaceAllString(prompt, " TABELA:$2 ")
	}

	// Marcar colunas no prompt
	for _, expr := range columnExpressions[lang] {
		re := regexp.MustCompile(expr)
		prompt = re.ReplaceAllString(prompt, " COLUNA:$2 ")
	}

	// Marcar colunas qualificadas no prompt (tabela.coluna)
	for _, expr := range qualifiedColumnExpressions[lang] {
		re := regexp.MustCompile(expr)
		prompt = re.ReplaceAllString(prompt, " TABELA:$1 COLUNA:$2 ")
	}

	// Adicionar expressões específicas para o caso do exemplo
	// "selecione a coluna X da tabela Y onde a coluna Z = valor"
	specificExpressions := map[string][]string{
		"pt": {
			`(?i)selecione\s+(?:a\s+)?(?:coluna\s+)?([a-zA-Z0-9_]+)\s+d[aeo]\s+(?:tabela\s+)?([a-zA-Z0-9_]+)`,
			`(?i)mostre\s+(?:a\s+)?(?:coluna\s+)?([a-zA-Z0-9_]+)\s+d[aeo]\s+(?:tabela\s+)?([a-zA-Z0-9_]+)`,
			`(?i)liste\s+(?:a\s+)?(?:coluna\s+)?([a-zA-Z0-9_]+)\s+d[aeo]\s+(?:tabela\s+)?([a-zA-Z0-9_]+)`,
			`(?i)exiba\s+(?:a\s+)?(?:coluna\s+)?([a-zA-Z0-9_]+)\s+d[aeo]\s+(?:tabela\s+)?([a-zA-Z0-9_]+)`,
		},
		"en": {
			`(?i)select\s+(?:the\s+)?(?:column\s+)?([a-zA-Z0-9_]+)\s+from\s+(?:table\s+)?([a-zA-Z0-9_]+)`,
			`(?i)show\s+(?:the\s+)?(?:column\s+)?([a-zA-Z0-9_]+)\s+from\s+(?:table\s+)?([a-zA-Z0-9_]+)`,
			`(?i)list\s+(?:the\s+)?(?:column\s+)?([a-zA-Z0-9_]+)\s+from\s+(?:table\s+)?([a-zA-Z0-9_]+)`,
			`(?i)display\s+(?:the\s+)?(?:column\s+)?([a-zA-Z0-9_]+)\s+from\s+(?:table\s+)?([a-zA-Z0-9_]+)`,
		},
	}

	// Aplicar as expressões específicas
	if _, exists := specificExpressions[lang]; !exists {
		lang = "en" // Fallback para inglês
	}

	for _, expr := range specificExpressions[lang] {
		re := regexp.MustCompile(expr)
		matches := re.FindAllStringSubmatch(prompt, -1)
		for _, match := range matches {
			if len(match) > 2 {
				colName := match[1]
				tableName := match[2]
				// Substituir por marcações explícitas
				prompt += " COLUNA:" + colName + " TABELA:" + tableName + " "
			}
		}
	}

	// Aplicar expressões para condições WHERE explícitas
	whereExpressions := map[string][]string{
		"pt": {
			`(?i)onde\s+(?:a\s+)?(?:coluna\s+)?([a-zA-Z0-9_\.]+)\s+=\s+(['"]?)([^'"]+)(['"]?)`,
			`(?i)onde\s+(?:a\s+)?(?:coluna\s+)?([a-zA-Z0-9_\.]+)\s+é\s+(['"]?)([^'"]+)(['"]?)`,
			`(?i)onde\s+(?:a\s+)?(?:coluna\s+)?([a-zA-Z0-9_\.]+)\s+>\s+(['"]?)([^'"]+)(['"]?)`,
			`(?i)onde\s+(?:a\s+)?(?:coluna\s+)?([a-zA-Z0-9_\.]+)\s+<\s+(['"]?)([^'"]+)(['"]?)`,
			`(?i)onde\s+(?:a\s+)?(?:coluna\s+)?([a-zA-Z0-9_\.]+)\s+>=\s+(['"]?)([^'"]+)(['"]?)`,
			`(?i)onde\s+(?:a\s+)?(?:coluna\s+)?([a-zA-Z0-9_\.]+)\s+<=\s+(['"]?)([^'"]+)(['"]?)`,
			`(?i)onde\s+(?:a\s+)?(?:coluna\s+)?([a-zA-Z0-9_\.]+)\s+!=\s+(['"]?)([^'"]+)(['"]?)`,
			`(?i)onde\s+(?:a\s+)?(?:coluna\s+)?([a-zA-Z0-9_\.]+)\s+(?:maior|superior)(?:\s+ou\s+igual)?\s+(?:a|que)\s+(['"]?)([^'"]+)(['"]?)`,
			`(?i)onde\s+(?:a\s+)?(?:coluna\s+)?([a-zA-Z0-9_\.]+)\s+(?:menor|inferior)(?:\s+ou\s+igual)?\s+(?:a|que)\s+(['"]?)([^'"]+)(['"]?)`,
			`(?i)onde\s+(?:a\s+)?(?:coluna\s+)?([a-zA-Z0-9_\.]+)\s+(?:diferente|não\s+é|distinto)\s+(?:de|que)?\s+(['"]?)([^'"]+)(['"]?)`,
			`(?i)onde\s+(?:a\s+)?(?:coluna\s+)?([a-zA-Z0-9_\.]+)\s+for\s+(['"]?)([^'"]+)(['"]?)`,
			`(?i)onde\s+(?:a\s+)?(?:coluna\s+)?([a-zA-Z0-9_\.]+)\s+for\s+(?:maior|superior)(?:\s+ou\s+igual)?\s+(?:a|que)\s+(['"]?)([^'"]+)(['"]?)`,
			`(?i)onde\s+(?:a\s+)?(?:coluna\s+)?([a-zA-Z0-9_\.]+)\s+for\s+(?:menor|inferior)(?:\s+ou\s+igual)?\s+(?:a|que)\s+(['"]?)([^'"]+)(['"]?)`,
			`(?i)com\s+(?:a\s+)?(?:coluna\s+)?([a-zA-Z0-9_\.]+)\s+=\s+(['"]?)([^'"]+)(['"]?)`,
			`(?i)que\s+(?:tenha|possua|contenha)\s+(?:a\s+)?(?:coluna\s+)?([a-zA-Z0-9_\.]+)\s+(?:igual a|como|=)\s+(['"]?)([^'"]+)(['"]?)`,
			`(?i)cujo\s+(?:a\s+)?(?:coluna\s+)?([a-zA-Z0-9_\.]+)\s+(?:seja|é|=)\s+(['"]?)([^'"]+)(['"]?)`,
			`(?i)filtrando\s+(?:por|onde)\s+(?:a\s+)?(?:coluna\s+)?([a-zA-Z0-9_\.]+)\s+(?:é|=)\s+(['"]?)([^'"]+)(['"]?)`,
			`(?i)para\s+(?:a\s+)?(?:coluna\s+)?([a-zA-Z0-9_\.]+)\s+(?:igual a|=)\s+(['"]?)([^'"]+)(['"]?)`,
			`(?i)quando\s+(?:a\s+)?(?:coluna\s+)?([a-zA-Z0-9_\.]+)\s+(?:for|é|=)\s+(['"]?)([^'"]+)(['"]?)`,
		},
		"en": {
			`(?i)where\s+(?:the\s+)?(?:column\s+)?([a-zA-Z0-9_\.]+)\s+=\s+(['"]?)([^'"]+)(['"]?)`,
			`(?i)where\s+(?:the\s+)?(?:column\s+)?([a-zA-Z0-9_\.]+)\s+is\s+(['"]?)([^'"]+)(['"]?)`,
			`(?i)where\s+(?:the\s+)?(?:column\s+)?([a-zA-Z0-9_\.]+)\s+>\s+(['"]?)([^'"]+)(['"]?)`,
			`(?i)where\s+(?:the\s+)?(?:column\s+)?([a-zA-Z0-9_\.]+)\s+<\s+(['"]?)([^'"]+)(['"]?)`,
			`(?i)where\s+(?:the\s+)?(?:column\s+)?([a-zA-Z0-9_\.]+)\s+>=\s+(['"]?)([^'"]+)(['"]?)`,
			`(?i)where\s+(?:the\s+)?(?:column\s+)?([a-zA-Z0-9_\.]+)\s+<=\s+(['"]?)([^'"]+)(['"]?)`,
			`(?i)where\s+(?:the\s+)?(?:column\s+)?([a-zA-Z0-9_\.]+)\s+!=\s+(['"]?)([^'"]+)(['"]?)`,
			`(?i)where\s+(?:the\s+)?(?:column\s+)?([a-zA-Z0-9_\.]+)\s+(?:greater|higher|more)\s+(?:than|than\s+or\s+equal\s+to)\s+(['"]?)([^'"]+)(['"]?)`,
			`(?i)where\s+(?:the\s+)?(?:column\s+)?([a-zA-Z0-9_\.]+)\s+(?:less|lower|smaller)\s+(?:than|than\s+or\s+equal\s+to)\s+(['"]?)([^'"]+)(['"]?)`,
			`(?i)where\s+(?:the\s+)?(?:column\s+)?([a-zA-Z0-9_\.]+)\s+(?:different|not|distinct)\s+(?:from|than)?\s+(['"]?)([^'"]+)(['"]?)`,
			`(?i)with\s+(?:the\s+)?(?:column\s+)?([a-zA-Z0-9_\.]+)\s+=\s+(['"]?)([^'"]+)(['"]?)`,
			`(?i)that\s+(?:has|have|contains)\s+(?:the\s+)?(?:column\s+)?([a-zA-Z0-9_\.]+)\s+(?:equal to|as|=)\s+(['"]?)([^'"]+)(['"]?)`,
			`(?i)whose\s+(?:the\s+)?(?:column\s+)?([a-zA-Z0-9_\.]+)\s+(?:is|=)\s+(['"]?)([^'"]+)(['"]?)`,
			`(?i)filtering\s+(?:by|where)\s+(?:the\s+)?(?:column\s+)?([a-zA-Z0-9_\.]+)\s+(?:is|=)\s+(['"]?)([^'"]+)(['"]?)`,
			`(?i)for\s+(?:the\s+)?(?:column\s+)?([a-zA-Z0-9_\.]+)\s+(?:equal to|=)\s+(['"]?)([^'"]+)(['"]?)`,
			`(?i)when\s+(?:the\s+)?(?:column\s+)?([a-zA-Z0-9_\.]+)\s+(?:is|=)\s+(['"]?)([^'"]+)(['"]?)`,
		},
	}

	for _, expr := range whereExpressions[lang] {
		re := regexp.MustCompile(expr)
		matches := re.FindAllStringSubmatch(prompt, -1)
		for _, match := range matches {
			if len(match) > 3 {
				colName := match[1]
				// Se tem um ponto, separar tabela e coluna
				if strings.Contains(colName, ".") {
					parts := strings.Split(colName, ".")
					if len(parts) == 2 {
						prompt += " TABELA:" + parts[0] + " COLUNA:" + parts[1] + " "
					}
				} else {
					prompt += " COLUNA:" + colName + " "
				}
			}
		}
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
			`(?i)(?:listar|mostrar|selecionar|exibir|visualizar)\s+(?:os|as)?\s+([a-z0-9_]+)`,
			`(?i)(?:no|na|nos|nas)\s+([a-z0-9_]+)`,
			`(?i)(?:dados|informações|registros)\s+(?:de|da|do|dos|das)\s+([a-z0-9_]+)`,
			`(?i)(?:quero|preciso|gostaria)\s+(?:ver|saber|consultar)\s+(?:os|as)?\s+([a-z0-9_]+)`,
			`(?i)(?:me\s+(?:mostre|traga|informe))\s+(?:os|as)?\s+([a-z0-9_]+)`,
			`(?i)(?:consultar)\s+(?:os|as)?\s+([a-z0-9_]+)`,
			`(?i)(?:analisar)\s+(?:os|as)?\s+([a-z0-9_]+)`,
			`(?i)(?:buscar)\s+(?:os|as)?\s+([a-z0-9_]+)`,
			`(?i)(?:de\s+tipo|da\s+categoria)\s+([a-z0-9_]+)`,
			`(?i)(?:de|da|do|dos|das)\s+(?:tabela|entidade)\s+([a-z0-9_]+)`,
			`(?i)(?:em\s+relação\s+(?:a|aos|às))\s+([a-z0-9_]+)`,
			`(?i)(?:cadastrados|registrados)\s+(?:em|na|no|nas|nos)\s+([a-z0-9_]+)`,
		}
	} else {
		patterns = []string{
			`(?i)(?:table|entity|from)\s+([a-z0-9_]+)`,
			`(?i)(?:all|the)\s+([a-z0-9_]+)`,
			`(?i)(?:list|show|select|display|view)\s+(?:the|all)?\s+([a-z0-9_]+)`,
			`(?i)(?:in|from|of)\s+the\s+([a-z0-9_]+)`,
			`(?i)(?:data|information|records)\s+(?:of|from|about)\s+([a-z0-9_]+)`,
			`(?i)(?:I\s+(?:want|need|would like)\s+to\s+(?:see|know|query))\s+(?:the|all)?\s+([a-z0-9_]+)`,
			`(?i)(?:show\s+me|give\s+me|tell\s+me\s+about)\s+(?:the|all)?\s+([a-z0-9_]+)`,
			`(?i)(?:query|search|analyze|analyze)\s+(?:the|all)?\s+([a-z0-9_]+)`,
			`(?i)(?:of\s+type|in\s+category)\s+([a-z0-9_]+)`,
			`(?i)(?:from\s+the\s+(?:table|entity))\s+([a-z0-9_]+)`,
			`(?i)(?:related\s+to)\s+([a-z0-9_]+)`,
			`(?i)(?:stored|registered)\s+in\s+(?:the)?\s+([a-z0-9_]+)`,
		}
	}

	// Try each pattern, starting with the most specific ones
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

	// Additional method: scan for terms that appear in the context and in the query
	words := strings.Fields(strings.ToLower(prompt))
	for _, word := range words {
		// Only consider potential entity names (ignore short words and common prepositions)
		if len(word) > 3 && !isCommonWord(word, p.language) {
			cleanWord := strings.Trim(word, ",.;:!?()[]{}\"'")
			if cleanWord != "" && (strings.Contains(context, "\""+cleanWord+"\"") ||
				strings.Contains(context, "'"+cleanWord+"'") ||
				strings.Contains(context, " "+cleanWord+" ")) {
				return cleanWord
			}
		}
	}

	return ""
}

// isCommonWord checks if a word is a common word that should be ignored for entity detection
func isCommonWord(word string, language string) bool {
	commonWords := map[string][]string{
		"en": {"the", "and", "for", "that", "with", "this", "from", "have", "are", "were", "been", "has", "their", "what", "when", "where", "which", "who", "whose", "whom"},
		"pt": {"de", "da", "do", "das", "dos", "para", "com", "que", "por", "como", "uma", "mais", "esse", "essa", "este", "esta", "aquele", "aquela", "seu", "sua", "seus", "suas", "meu", "minha", "nosso", "nossa", "qual", "quais", "onde", "quando", "quem"},
	}

	if words, exists := commonWords[language]; exists {
		for _, commonWord := range words {
			if word == commonWord {
				return true
			}
		}
	}
	return false
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
			`(?i)(?:mostrar|selecionar|exibir|obter|listar)\s+([^,\.]+(?:,\s*[^,\.]+)*)\s+(?:de|da|dos|das)`,
			`(?i)(?:me mostre|me traga|mostre-me|informe-me|me informe)\s+(?:o|os|a|as)?\s+([^,\.]+(?:,\s*[^,\.]+)*)\s+(?:de|da|dos|das)`,
			`(?i)(?:quero|preciso|gostaria de)\s+(?:ver|saber|consultar)\s+(?:o|os|a|as)?\s+([^,\.]+(?:,\s*[^,\.]+)*)\s+(?:de|da|dos|das)`,
			`(?i)(?:contendo|incluindo|com)\s+(?:os campos|as colunas|os atributos)\s+([^,\.]+(?:,\s*[^,\.]+)*)`,
			`(?i)(?:retorne|retornar|recupere|recuperar)\s+(?:o|os|a|as)?\s+([^,\.]+(?:,\s*[^,\.]+)*)\s+(?:de|da|dos|das)`,
			`(?i)(?:visualizar|consultar|analisar)\s+(?:o|os|a|as)?\s+([^,\.]+(?:,\s*[^,\.]+)*)\s+(?:de|da|dos|das)`,
			`(?i)(?:dados|informações)\s+(?:sobre|de|da|do)\s+([^,\.]+(?:,\s*[^,\.]+)*)`,
			`(?i)(?:valores|registros)\s+(?:de|da|do)\s+([^,\.]+(?:,\s*[^,\.]+)*)`,
		}
	} else {
		patterns = []string{
			`(?i)fields\s+(?:like|such as)\s+([^,\.]+(?:,\s*[^,\.]+)*)`,
			`(?i)columns\s+(?:like|such as)\s+([^,\.]+(?:,\s*[^,\.]+)*)`,
			`(?i)attributes\s+(?:like|such as)\s+([^,\.]+(?:,\s*[^,\.]+)*)`,
			`(?i)(?:show|select|get|retrieve|list|display)\s+([^,\.]+(?:,\s*[^,\.]+)*)\s+(?:from|of)`,
			`(?i)(?:show me|give me|tell me about)\s+(?:the)?\s+([^,\.]+(?:,\s*[^,\.]+)*)\s+(?:from|of)`,
			`(?i)(?:I want|I need|I would like to)\s+(?:see|know|query)\s+(?:the)?\s+([^,\.]+(?:,\s*[^,\.]+)*)\s+(?:from|of)`,
			`(?i)(?:containing|including|with)\s+(?:the fields|the columns|the attributes)\s+([^,\.]+(?:,\s*[^,\.]+)*)`,
			`(?i)(?:return|retrieve|get me)\s+(?:the)?\s+([^,\.]+(?:,\s*[^,\.]+)*)\s+(?:from|of)`,
			`(?i)(?:view|query|analyze)\s+(?:the)?\s+([^,\.]+(?:,\s*[^,\.]+)*)\s+(?:from|of)`,
			`(?i)(?:data|information)\s+(?:about|of|from)\s+([^,\.]+(?:,\s*[^,\.]+)*)`,
			`(?i)(?:values|records)\s+(?:of|from)\s+([^,\.]+(?:,\s*[^,\.]+)*)`,
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
			if len(attributes) > 0 {
				break
			}
		}
	}

	// If no attributes were found with patterns, try to extract from specific keywords
	if len(attributes) == 0 {
		// Keywords that might indicate what comes next is an attribute
		var attributeIndicators []string
		if p.language == "pt" {
			attributeIndicators = []string{
				"coluna", "colunas", "campo", "campos", "atributo", "atributos",
				"ordenado por", "ordenados por", "agrupado por", "agrupados por",
				"filtrar por", "filtrados por", "valor de", "valores de",
			}
		} else {
			attributeIndicators = []string{
				"column", "columns", "field", "fields", "attribute", "attributes",
				"ordered by", "grouped by", "filtered by", "value of", "values of",
			}
		}

		words := strings.Fields(prompt)
		for i, word := range words {
			wordLower := strings.ToLower(word)
			for _, indicator := range attributeIndicators {
				if wordLower == indicator && i+1 < len(words) {
					// The word after the indicator might be an attribute
					possibleAttribute := words[i+1]
					// Remove any punctuation
					possibleAttribute = strings.Trim(possibleAttribute, ",.;:!?()[]{}\"'")
					if len(possibleAttribute) > 1 && !isCommonWord(possibleAttribute, p.language) {
						attributes = append(attributes, possibleAttribute)
					}
				}
			}
		}
	}

	// Remove duplicates
	uniqueAttributes := make([]string, 0, len(attributes))
	seen := make(map[string]bool)
	for _, attr := range attributes {
		attrLower := strings.ToLower(attr)
		if !seen[attrLower] {
			seen[attrLower] = true
			uniqueAttributes = append(uniqueAttributes, attr)
		}
	}

	return uniqueAttributes
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
