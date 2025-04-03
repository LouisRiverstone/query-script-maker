<template>
    <div class="flex flex-col gap-3 w-full" :class="{ 'editor-dark-mode': isDarkMode }">
        <slot></slot>
        <div class="flex flex-col lg:flex-row gap-6">
            <div class="flex flex-col w-full gap-3">
                <div class="flex flex-col justify-between items-center">
                    <Divider v-if="showBindedSql">Input</Divider>
                    <div class="flex items-center gap-2">
                        <!-- Barra de ferramentas do editor -->
                        <div class="flex items-center border border-gray-200 dark:border-gray-700 rounded-md bg-white dark:bg-gray-800 shadow-sm">
                            <button 
                                @click="formatSQL" 
                                title="Format SQL Query (Ctrl+Shift+F)"
                                class="px-2 py-1 text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 border-r border-gray-200 dark:border-gray-700 text-sm"
                            >
                                <span class="text-xs">⚙️ Format</span>
                            </button>
                            <button 
                                @click="copyToClipboard" 
                                title="Copy to Clipboard"
                                class="px-2 py-1 text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 border-r border-gray-200 dark:border-gray-700 text-sm"
                            >
                                <span class="text-xs">📋 Copy</span>
                            </button>
                            <button 
                                @click="clearEditor" 
                                title="Clear Editor"
                                class="px-2 py-1 text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 border-r border-gray-200 dark:border-gray-700 text-sm"
                            >
                                <span class="text-xs">🧹 Clear</span>
                            </button>
                            <div class="relative snippets-menu">
                                <button 
                                    @click="showSnippets = !showSnippets" 
                                    title="SQL Snippets"
                                    class="px-2 py-1 text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 border-r border-gray-200 dark:border-gray-700 text-sm"
                                >
                                    <span class="text-xs">📝 Snippets</span>
                                </button>
                                <div v-if="showSnippets" class="absolute top-full left-0 mt-1 z-10 w-64 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-md shadow-lg">
                                    <div class="p-2 border-b border-gray-200 dark:border-gray-700 text-xs font-medium text-gray-700 dark:text-gray-300">
                                        Common Snippets
                                    </div>
                                    <div class="max-h-60 overflow-y-auto">
                                        <button 
                                            v-for="(snippet, index) in sqlSnippets" 
                                            :key="index"
                                            @click="insertSnippet(snippet.code)"
                                            class="w-full text-left px-3 py-2 text-xs text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700"
                                        >
                                            {{ snippet.name }}
                                        </button>
                                    </div>
                                </div>
                            </div>
                            <div class="relative history-menu">
                                <button 
                                    @click="showHistory = !showHistory" 
                                    title="Query History"
                                    class="px-2 py-1 text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 border-r border-gray-200 dark:border-gray-700 text-sm"
                                >
                                    <span class="text-xs">🕒 History</span>
                                </button>
                                <div v-if="showHistory" class="absolute top-full left-0 mt-1 z-10 w-64 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-md shadow-lg">
                                    <div class="p-2 border-b border-gray-200 dark:border-gray-700 text-xs font-medium text-gray-700 dark:text-gray-300">
                                        Query History
                                    </div>
                                    <div v-if="queryHistory.length === 0" class="p-3 text-xs text-gray-500 dark:text-gray-400 text-center">
                                        No queries in history
                                    </div>
                                    <div v-else class="max-h-60 overflow-y-auto">
                                        <button 
                                            v-for="(query, index) in queryHistory" 
                                            :key="index"
                                            @click="loadFromHistory(query)"
                                            class="w-full text-left p-2 text-xs text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 border-b border-gray-100 dark:border-gray-700"
                                        >
                                            <div class="truncate mb-1">{{ query.text.substring(0, 50) }}{{ query.text.length > 50 ? '...' : '' }}</div>
                                            <div class="text-gray-400 dark:text-gray-500 text-[10px]">{{ formatDate(query.timestamp) }}</div>
                                        </button>
                                    </div>
                                    <div class="p-2 border-t border-gray-200 dark:border-gray-700">
                                        <button 
                                            @click="clearHistory"
                                            class="w-full text-xs text-red-500 dark:text-red-400 hover:bg-gray-100 dark:hover:bg-gray-700 p-1 rounded"
                                        >
                                            Clear history
                                        </button>
                                    </div>
                                </div>
                            </div>
                            <!-- AI Assistant Button -->
                            <button 
                                @click="toggleSQLAssistant" 
                                title="SQL AI Assistant"
                                class="px-2 py-1 text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 text-sm"
                            >
                                <span class="text-xs">🤖 AI Assistant</span>
                            </button>
                        </div>
                    </div>
                </div>
                
                <!-- SQL AI Assistant -->
                <SQLAssistant 
                    v-if="showSQLAssistant" 
                    :isVisible="showSQLAssistant"
                    @update:isVisible="showSQLAssistant = $event"
                    @sqlGenerated="handleSQLGenerated"
                    class="mb-3"
                />
                
                <div class="flex flex-col md:flex-row items-start justify-center md:items-center gap-2">
                    <div class="relative variables-menu">
                        <button 
                            @click="showVariables = !showVariables"
                            class="inline-flex items-center text-xs text-gray-600 dark:text-gray-400 italic border border-gray-200 dark:border-gray-700 rounded px-2 py-1 hover:bg-gray-50 dark:hover:bg-gray-800"
                        >
                            <span>Variables available</span>
                            <svg xmlns="http://www.w3.org/2000/svg" class="h-3 w-3 ml-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                            </svg>
                        </button>
                        <div v-if="showVariables && props.variables && props.variables.length > 0" class="absolute top-full left-0 mt-1 z-10 w-64 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-md shadow-lg">
                            <div class="p-2 border-b border-gray-200 dark:border-gray-700 text-xs font-medium text-gray-700 dark:text-gray-300">
                                Insert Variable
                            </div>
                            <div class="max-h-60 overflow-y-auto">
                                <button 
                                    v-for="variable in props.variables" 
                                    :key="variable.Field"
                                    @click="insertVariable(variable.Value)"
                                    class="w-full text-left px-3 py-2 text-xs text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 flex justify-between items-center"
                                >
                                    <span>{{ variable.Value }}</span>
                                    <span class="text-gray-400 dark:text-gray-500">{{ variable.Field }}</span>
                                </button>
                            </div>
                        </div>
                        <div v-else-if="showVariables" class="absolute top-full left-0 mt-1 z-10 w-64 bg-white dark:bg-gray-800 border border-gray-200 dark:border-gray-700 rounded-md shadow-lg p-3 text-xs text-gray-500 dark:text-gray-400 text-center">
                            No variables available
                        </div>
                    </div>
                    
                    <div v-if="copySuccess" class="text-xs text-green-500 px-2 py-1 transition-opacity duration-1000" :class="{'opacity-0': !copySuccess}">
                        Copied! ✓
                    </div>
                </div>
                
                <!-- Editor principal -->
                <code-mirror v-model="value" :lang="lang" :extensions="extensions" :linter="null" basic wrap tab class="w-full rounded-md border border-gray-200 dark:border-gray-700 shadow-sm" />
                
                <!-- Botão para salvar consulta no histórico quando não está no modo visualização -->
                <div class="flex justify-end">
                    <button 
                        @click="saveToHistory"
                        class="text-xs text-gray-600 dark:text-gray-400 hover:text-sky-600 dark:hover:text-sky-400 px-2 py-1 rounded"
                    >
                        Save in history
                    </button>
                </div>
            </div>
            <div v-if="showBindedSql" class="flex flex-col w-full gap-3">
                <Divider>Output</Divider>
                <code-mirror v-model="linesBinded" :lang="lang" :extensions="outputExtensions" :linter="null" basic wrap tab class="w-full rounded-md border border-gray-200 dark:border-gray-700 shadow-sm" />
            </div>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { ref, watch, onMounted, computed, onUnmounted, shallowRef } from 'vue';
import CodeMirror from 'vue-codemirror6';
import { sql, MySQL } from "@codemirror/lang-sql";
import { oneDarkTheme } from '@codemirror/theme-one-dark';
import { MakeBindedSQL, GetLatestDatabaseStructure } from '../../wailsjs/go/main/App';
import { computedAsync } from '@vueuse/core';
import { autocompletion, CompletionContext, CompletionResult } from '@codemirror/autocomplete';
import { EditorView, Decoration, ViewPlugin, ViewUpdate, DecorationSet } from '@codemirror/view';
import { useDebounceFn } from '@vueuse/core';

import Divider from './Divider.vue';
import SQLAssistant from './SQLAssistant.vue';
import { main } from '../../wailsjs/go/models';

const variableText = ref("{{ variable }}");
const isDarkMode = ref(false);
let darkModeObserver: MutationObserver;
const copySuccess = ref(false);
const showSnippets = ref(false);
const showHistory = ref(false);
const showVariables = ref(false);
const showSQLAssistant = ref(false);

// Array para armazenar o histórico de consultas
interface QueryHistoryItem {
  text: string;
  timestamp: number;
}

const HISTORY_STORAGE_KEY = 'sql_query_history';
const queryHistory = ref<QueryHistoryItem[]>([]);

// SQL keywords for autocomplete - Moved outside of component setup to prevent recreating on every render
const sqlKeywordsSource = [
  // Common SQL Commands
  { label: "SELECT", type: "keyword", info: "Retrieves data from a database" },
  { label: "FROM", type: "keyword", info: "Specifies which table to select or delete data from" },
  { label: "WHERE", type: "keyword", info: "Filters records based on a condition" },
  { label: "INSERT INTO", type: "keyword", info: "Inserts new data into a database" },
  { label: "UPDATE", type: "keyword", info: "Modifies existing database data" },
  { label: "DELETE", type: "keyword", info: "Deletes data from a database" },
  { label: "CREATE", type: "keyword", info: "Creates a new database object (table, view, etc.)" },
  { label: "ALTER", type: "keyword", info: "Modifies an existing database object" },
  { label: "DROP", type: "keyword", info: "Deletes an existing database object" },
  { label: "TRUNCATE", type: "keyword", info: "Removes all records from a table, but not the table itself" },
  { label: "RENAME", type: "keyword", info: "Renames a database object" },
  { label: "SHOW", type: "keyword", info: "Shows information about databases, tables, columns, or status" },
  { label: "DESCRIBE", type: "keyword", info: "Shows the structure of a table" },
  { label: "EXPLAIN", type: "keyword", info: "Shows the execution plan for a query" },
  { label: "USE", type: "keyword", info: "Selects a database" },
  
  // Join Types
  { label: "JOIN", type: "keyword", info: "Combines rows from two or more tables" },
  { label: "LEFT JOIN", type: "keyword", info: "Returns all records from the left table, and matched records from the right table" },
  { label: "RIGHT JOIN", type: "keyword", info: "Returns all records from the right table, and matched records from the left table" },
  { label: "INNER JOIN", type: "keyword", info: "Returns records that have matching values in both tables" },
  { label: "CROSS JOIN", type: "keyword", info: "Returns the Cartesian product of two tables" },
  { label: "NATURAL JOIN", type: "keyword", info: "Joins tables by automatically finding matching column names" },
  { label: "FULL JOIN", type: "keyword", info: "Returns all records when there is a match in either left or right table" },
  { label: "SELF JOIN", type: "keyword", info: "Joins a table to itself" },
  
  // Clauses
  { label: "GROUP BY", type: "keyword", info: "Groups rows that have the same values into summary rows" },
  { label: "HAVING", type: "keyword", info: "Filters records after GROUP BY is applied" },
  { label: "ORDER BY", type: "keyword", info: "Sorts the result set in ascending or descending order" },
  { label: "LIMIT", type: "keyword", info: "Limits the number of records returned" },
  { label: "OFFSET", type: "keyword", info: "Specifies where to start selecting records" },
  { label: "UNION", type: "keyword", info: "Combines the result sets of two or more SELECT statements" },
  { label: "UNION ALL", type: "keyword", info: "Combines the result sets of two or more SELECT statements (allows duplicates)" },
  { label: "INTERSECT", type: "keyword", info: "Returns the records that both queries have" },
  { label: "EXCEPT", type: "keyword", info: "Returns the records from the first query that are not in the second query" },
  { label: "WITH", type: "keyword", info: "Specifies temporary named result sets (Common Table Expressions)" },
  
  // Common Functions and Keywords - Adding just a subset for performance
  { label: "COUNT()", type: "function", info: "Returns the number of rows" },
  { label: "SUM()", type: "function", info: "Returns the sum of values" },
  { label: "AVG()", type: "function", info: "Returns the average value" },
  { label: "MIN()", type: "function", info: "Returns the minimum value" },
  { label: "MAX()", type: "function", info: "Returns the maximum value" },
  { label: "CASE", type: "keyword", info: "Evaluates conditions and returns a value when the first condition is met" },
  { label: "WHEN", type: "keyword", info: "Used in a CASE statement to specify a condition" },
  { label: "THEN", type: "keyword", info: "Used in a CASE statement to specify a result" },
  { label: "ELSE", type: "keyword", info: "Used in a CASE statement to specify a value to return if all conditions are false" },
  { label: "END", type: "keyword", info: "Used in a CASE statement to end the list of conditions" }
];

// Memoize expensive operations
const loadDatabaseStructure = async () => {
    try {
        // Check if we already have the structure
        if (databaseStructure.value && databaseStructure.value.tables.length > 0) {
            console.log("Using cached database structure");
            return;
        }
        
        console.log("Loading database structure...");
        const structure = await GetLatestDatabaseStructure();
        if (structure) {
            const parsed = JSON.parse(structure) as DatabaseStructure;
            databaseStructure.value = parsed;
            
            // Process structure to create completions only once
            if (parsed && parsed.tables) {
                // Create table completions
                const tables = parsed.tables.map(table => ({
                    label: table.name,
                    type: "table",
                    detail: `Table with ${table.columns.length} columns`,
                    info: `Table: ${table.name}`,
                    boost: 99
                }));
                
                // Create column completions grouped by table
                const columns: any[] = [];
                parsed.tables.forEach(table => {
                    table.columns.forEach(column => {
                        columns.push({
                            label: column.name,
                            type: "column",
                            detail: `${table.name}.${column.name}`,
                            info: `Column: ${column.name} (${column.type})${column.isPrimary ? ' [PK]' : ''}`,
                            apply: column.name, 
                            boost: column.isPrimary ? 95 : 90,
                            table: table.name
                        });
                        
                        // Only add qualified columns for primary keys and frequently used columns to reduce size
                        if (column.isPrimary || column.name.toLowerCase().includes('id') || column.name.toLowerCase().includes('name')) {
                            columns.push({
                                label: `${table.name}.${column.name}`,
                                type: "qualified-column",
                                detail: `${table.name}.${column.name}`,
                                info: `Column: ${column.name} (${column.type})${column.isPrimary ? ' [PK]' : ''}`,
                                apply: `${table.name}.${column.name}`,
                                boost: column.isPrimary ? 85 : 80
                            });
                        }
                    });
                });
                
                // Use shallowRef to avoid deep reactivity costs
                tableCompletions.value = tables;
                columnCompletions.value = columns;
            }
        }
    } catch (error) {
        console.error("Error loading database structure:", error);
    }
};

// Carregar histórico do localStorage com memoização
let cachedHistory: QueryHistoryItem[] | null = null;

const loadQueryHistory = () => {
    if (cachedHistory !== null) {
        queryHistory.value = cachedHistory;
        return;
    }
    
    try {
        const savedHistory = localStorage.getItem(HISTORY_STORAGE_KEY);
        if (savedHistory) {
            const parsed = JSON.parse(savedHistory);
            queryHistory.value = parsed;
            cachedHistory = parsed;
        }
    } catch (error) {
        console.error("Error loading query history:", error);
        queryHistory.value = [];
        cachedHistory = [];
    }
};

// Salvar consulta atual no histórico - com otimização para evitar atualizações desnecessárias
const saveToHistory = () => {
    if (!value.value.trim()) return;
    
    // Evitar duplicatas recentes
    const isDuplicate = queryHistory.value.some(item => 
        item.text === value.value && 
        (Date.now() - item.timestamp) < 60000 // 1 minuto
    );
    
    if (!isDuplicate) {
        const newItem: QueryHistoryItem = {
            text: value.value,
            timestamp: Date.now()
        };
        
        // Limitar o histórico a 20 itens para melhor desempenho
        const newHistory = [newItem, ...queryHistory.value].slice(0, 20);
        queryHistory.value = newHistory;
        cachedHistory = newHistory;
        
        // Usar localStorage em um setTimeout para não bloquear a UI
        setTimeout(() => {
            localStorage.setItem(HISTORY_STORAGE_KEY, JSON.stringify(newHistory));
        }, 0);
    }
};

// Carregar consulta do histórico
const loadFromHistory = (item: QueryHistoryItem) => {
  value.value = item.text;
  showHistory.value = false;
};

// Limpar histórico
const clearHistory = () => {
  if (confirm("Are you sure you want to clear the entire query history?")) {
    queryHistory.value = [];
    localStorage.removeItem(HISTORY_STORAGE_KEY);
    showHistory.value = false;
  }
};

// Formatar data para exibição no histórico
const formatDate = (timestamp: number): string => {
  const date = new Date(timestamp);
  return date.toLocaleDateString() + ' ' + date.toLocaleTimeString();
};

// Inserir variável no editor
const insertVariable = (variableName: string) => {
  const variableText = `{{ ${variableName} }}`;
  
  // Inserir no cursor ou substituir seleção
  // A implementação específica depende da API do editor
  
  value.value += variableText; // Método simples por enquanto
  showVariables.value = false;
};

// Lista de snippets SQL
const sqlSnippets = [
  { 
    name: 'Basic SELECT', 
    code: 'SELECT * FROM table_name WHERE condition;' 
  },
  { 
    name: 'Basic INSERT', 
    code: 'INSERT INTO table_name (column1, column2) VALUES (value1, value2);' 
  },
  { 
    name: 'Basic UPDATE', 
    code: 'UPDATE table_name SET column1 = value1 WHERE condition;' 
  },
  { 
    name: 'Basic DELETE', 
    code: 'DELETE FROM table_name WHERE condition;' 
  },
  { 
    name: 'Simple JOIN', 
    code: 'SELECT t1.column1, t2.column2\nFROM table1 t1\nINNER JOIN table2 t2 ON t1.id = t2.table1_id\nWHERE condition;' 
  },
  { 
    name: 'GROUP BY with aggregation', 
    code: 'SELECT column1, COUNT(*) as count\nFROM table_name\nGROUP BY column1\nHAVING count > 1;' 
  },
  { 
    name: 'Subquery', 
    code: 'SELECT column1\nFROM table1\nWHERE column2 IN (SELECT column2 FROM table2 WHERE condition);' 
  },
  { 
    name: 'CASE WHEN', 
    code: 'SELECT column1,\n  CASE\n    WHEN condition1 THEN result1\n    WHEN condition2 THEN result2\n    ELSE result3\n  END AS column_alias\nFROM table_name;' 
  },
  { 
    name: 'WITH CTE', 
    code: 'WITH cte_name AS (\n  SELECT column1, column2\n  FROM table_name\n  WHERE condition\n)\nSELECT * FROM cte_name;' 
  },
  { 
    name: 'UNION', 
    code: 'SELECT column1 FROM table1\nUNION\nSELECT column1 FROM table2;' 
  }
];

// Ações da barra de ferramentas
const formatSQL = () => {
  try {
    // Implementação básica de formatação
    // Esta é uma abordagem simples, no futuro pode-se usar uma biblioteca de formatação SQL
    const formatted = value.value
      .replace(/SELECT/gi, 'SELECT\n  ')
      .replace(/FROM/gi, '\nFROM\n  ')
      .replace(/WHERE/gi, '\nWHERE\n  ')
      .replace(/GROUP BY/gi, '\nGROUP BY\n  ')
      .replace(/HAVING/gi, '\nHAVING\n  ')
      .replace(/ORDER BY/gi, '\nORDER BY\n  ')
      .replace(/JOIN/gi, '\nJOIN\n  ')
      .replace(/LEFT JOIN/gi, '\nLEFT JOIN\n  ')
      .replace(/RIGHT JOIN/gi, '\nRIGHT JOIN\n  ')
      .replace(/INNER JOIN/gi, '\nINNER JOIN\n  ')
      .replace(/OUTER JOIN/gi, '\nOUTER JOIN\n  ')
      .replace(/UNION/gi, '\nUNION\n')
      .replace(/,/g, ',\n  ');
    
    value.value = formatted;
  } catch (error) {
    console.error("Error formatting SQL:", error);
  }
};

const copyToClipboard = async () => {
  try {
    await navigator.clipboard.writeText(value.value);
    copySuccess.value = true;
    setTimeout(() => {
      copySuccess.value = false;
    }, 2000);
  } catch (error) {
    console.error("Error copying to clipboard:", error);
  }
};

const clearEditor = () => {
  if (confirm("Limpar todo o conteúdo do editor?")) {
    value.value = "";
  }
};

const insertSnippet = (snippetCode: string) => {
  value.value = snippetCode;
  showSnippets.value = false;
};

// Fechar menu de snippets ao clicar fora
onMounted(() => {
  document.addEventListener('click', (event) => {
    if (showSnippets.value && !(event.target as Element).closest('.snippets-menu')) {
      showSnippets.value = false;
    }
  });
});

// Database structure types
interface DatabaseColumn {
  name: string;
  type: string;
  nullable: string;
  key: string;
  default: string;
  extra: string;
  isPrimary: boolean;
}

interface DatabaseForeignKey {
  columnName: string;
  referencedTable: string;
  referencedColumn: string;
  constraintName: string;
}

interface DatabaseTable {
  name: string;
  columns: DatabaseColumn[];
  foreignKeys: DatabaseForeignKey[];
}

interface DatabaseStructure {
  tables: DatabaseTable[];
}

// Variables for table and column autocomplete - using shallowRef for better performance with large objects
const databaseStructure = shallowRef<DatabaseStructure | null>(null);
const tableCompletions = shallowRef<any[]>([]);
const columnCompletions = shallowRef<any[]>([]);
// Use shallowRef for large static array to prevent deep reactivity costs
const sqlKeywords = shallowRef(sqlKeywordsSource);

// Função para fechar menus quando clica fora
const handleClickOutside = (event: MouseEvent): void => {
    const target = event.target as Element;
    
    if (showSnippets.value && !target.closest('.snippets-menu')) {
        showSnippets.value = false;
    }
    
    if (showHistory.value && !target.closest('.history-menu')) {
        showHistory.value = false;
    }
    
    if (showVariables.value && !target.closest('.variables-menu')) {
        showVariables.value = false;
    }
};

onMounted(() => {
    // Initial check
    updateDarkMode();
    
    // Watch for changes in dark mode - use passive listeners for better performance
    darkModeObserver = new MutationObserver(updateDarkMode);
    
    // Observe changes to the class attribute of both html and body elements
    darkModeObserver.observe(document.documentElement, { 
        attributes: true, 
        attributeFilter: ['class'] 
    });
    
    darkModeObserver.observe(document.body, { 
        attributes: true, 
        attributeFilter: ['class'] 
    });
    
    // Listen for system preference changes
    window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', updateDarkMode, { passive: true });
    
    // Load database structure
    loadDatabaseStructure();
    
    // Load query history
    loadQueryHistory();
    
    // Close menus when clicking outside - use passive and capture for performance
    document.addEventListener('click', handleClickOutside, { passive: true });
    
    // Configure global keyboard shortcuts
    document.addEventListener('keydown', handleKeyboardShortcuts, { passive: false });
});

onUnmounted(() => {
    // Clean up observers and event listeners
    if (darkModeObserver) {
        darkModeObserver.disconnect();
    }
    
    window.matchMedia('(prefers-color-scheme: dark)').removeEventListener('change', updateDarkMode);
    
    // Remove event listeners
    document.removeEventListener('click', handleClickOutside);
    document.removeEventListener('keydown', handleKeyboardShortcuts);
});

// Tratador de atalhos de teclado
const handleKeyboardShortcuts = (event: KeyboardEvent): void => {
    // Ctrl+Shift+F or Cmd+Shift+F to format SQL
    if ((event.ctrlKey || event.metaKey) && event.shiftKey && event.key === 'F') {
        event.preventDefault();
        formatSQL();
    }
    
    // Ctrl+Enter or Cmd+Enter to execute (if implemented)
    if ((event.ctrlKey || event.metaKey) && event.key === 'Enter') {
        event.preventDefault();
        // If there's a function to execute the query, call it here
        saveToHistory(); // For now, just save to history
    }
    
    // Esc to close open menus
    if (event.key === 'Escape') {
        showSnippets.value = false;
        showHistory.value = false;
        showVariables.value = false;
    }
};

const props = defineProps<{
    variables?: Array<main.Variable>
    data?: { [key: string]: any }[]
    modelValue: string,
    showBindedSql?: boolean
    minify?: boolean
}>()

const emit = defineEmits(['update:modelValue']);

const value = ref(props.modelValue);

const lang = sql({
    dialect: MySQL
});

// Autocomplete function for variables
const completeVariables = (context: CompletionContext): CompletionResult | null => {
    // Get the word at cursor
    const word = context.matchBefore(/\{\{\s*\w*\s*\}?\}?/);
    
    // If we're not typing a variable, don't show completions
    if (!word || (word.from === word.to && !context.explicit && !word.text.startsWith("{{"))) {
        return null;
    }
    
    return {
        from: word.from,
        options: variableCompletions.value,
        validFor: /\{\{\s*\w*\s*\}?\}?/,
    };
};

// Autocomplete function for SQL - optimized for performance
const completeSql = (context: CompletionContext): CompletionResult | null => {
    // Skip if we're inside a variable template
    const variablePattern = /\{\{\s*\w*\s*\}?\}?/;
    const beforeCursor = context.state.doc.sliceString(
        Math.max(0, context.pos - 20), 
        context.pos
    );
    
    if (variablePattern.test(beforeCursor)) {
        return null;
    }
    
    // Get the word at cursor
    const word = context.matchBefore(/\w+/);
    
    if (!word) {
        return null;
    }
    
    const wordLower = word.text.toLowerCase();
    
    // First check for table names after FROM or JOIN
    const textBefore = context.state.doc.sliceString(
        Math.max(0, context.pos - 100), // Only check the last 100 characters for better performance
        context.pos
    ).toLowerCase();
    
    const lastFromOrJoin = Math.max(
        textBefore.lastIndexOf("from "), 
        textBefore.lastIndexOf("join ")
    );
    const lastComma = textBefore.lastIndexOf(",", context.pos);
    
    // If we're after FROM or JOIN, prioritize table completions
    if (lastFromOrJoin > -1 && lastFromOrJoin > lastComma) {
        // Only filter for tables that match the current word - using faster array methods
        const filteredTables = tableCompletions.value.filter(table => 
            table.label.toLowerCase().indexOf(wordLower) === 0
        );
        
        if (filteredTables.length > 0 || context.explicit) {
            return {
                from: word.from,
                options: filteredTables.slice(0, 100), // Limit to 100 results for performance
                validFor: /\w+/,
            };
        }
    }
    
    // Check for column completions
    const periodMatch = context.matchBefore(/(\w+)\.\w*/);
    if (periodMatch) {
        const tableName = periodMatch.text.split('.')[0];
        const columnWord = context.matchBefore(/\w*$/);
        
        if (columnWord) {
            const columnWordLower = columnWord.text.toLowerCase();
            // Get columns for the specific table - using filter with early return for performance
            const tableSpecificColumns = columnCompletions.value.filter(column => {
                if (column.table !== tableName || column.type !== "column") return false;
                return column.label.toLowerCase().indexOf(columnWordLower) === 0;
            }).slice(0, 100); // Limit results
            
            if (tableSpecificColumns.length > 0 || context.explicit) {
                return {
                    from: columnWord.from,
                    options: tableSpecificColumns,
                    validFor: /\w+/,
                };
            }
        }
    }
    
    // For performance, first check common keywords that might match
    const commonKeywords = ["SELECT", "FROM", "WHERE", "JOIN", "GROUP", "ORDER", "HAVING", "LIMIT"];
    if (commonKeywords.some(kw => kw.toLowerCase().indexOf(wordLower) === 0)) {
        const filteredCommonKeywords = sqlKeywords.value.filter((kw: { label: string }) => 
            kw.label.toLowerCase().indexOf(wordLower) === 0 && 
            commonKeywords.includes(kw.label)
        );
        
        if (filteredCommonKeywords.length > 0) {
            return {
                from: word.from,
                options: filteredCommonKeywords,
                validFor: /\w+/,
            };
        }
    }
    
    // Only if the word has at least 2 characters or explicit request, search all options
    if (wordLower.length >= 2 || context.explicit) {
        // Use a faster filtering approach with maximum limits
        const filteredOptions = [];
        const maxResults = 100;
        
        // First prioritize keywords for better UX
        for (const option of sqlKeywords.value) {
            if (option.label.toLowerCase().indexOf(wordLower) === 0) {
                filteredOptions.push(option);
                if (filteredOptions.length >= maxResults / 2) break;
            }
        }
        
        // Then add columns if we have space
        if (filteredOptions.length < maxResults) {
            for (const option of columnCompletions.value) {
                if (option.label.toLowerCase().indexOf(wordLower) === 0) {
                    filteredOptions.push(option);
                    if (filteredOptions.length >= maxResults) break;
                }
            }
        }
        
        if (filteredOptions.length === 0 && !context.explicit) {
            return null;
        }
        
        return {
            from: word.from,
            options: filteredOptions,
            validFor: /\w+/,
        };
    }
    
    return null;
};

// Add custom styling for autocomplete items
const myTheme = EditorView.baseTheme({
    ".cm-tooltip.cm-tooltip-autocomplete": {
        border: "1px solid #ddd",
        background: "white",
        fontSize: "90%",
        borderRadius: "4px"
    },
    ".cm-tooltip.cm-tooltip-autocomplete.cm-tooltip-autocomplete-dark": {
        backgroundColor: "#1f2937", // dark:bg-gray-800
        border: "1px solid #374151", // dark:border-gray-700
        color: "white"
    },
    ".cm-tooltip-autocomplete .cm-completionIcon.cm-completionIcon-variable": {
        "&:after": { content: "'{{}}'" },
        color: "#0284c7", // text-sky-600
        fontWeight: "bold"
    },
    ".cm-completionLabel.cm-completionLabel-variable": {
        color: "#0284c7" // text-sky-600
    },
    ".cm-tooltip-autocomplete-dark .cm-completionIcon.cm-completionIcon-variable": {
        "&:after": { content: "'{{}}'" },
        color: "#38bdf8", // dark:text-sky-400
    },
    ".cm-tooltip-autocomplete-dark .cm-completionLabel.cm-completionLabel-variable": {
        color: "#38bdf8" // dark:text-sky-400
    },
    ".cm-tooltip-autocomplete .cm-completionIcon.cm-completionIcon-keyword": {
        "&:after": { content: "'k'" },
        color: "#7c3aed", // text-purple-600
        fontWeight: "bold" 
    },
    ".cm-completionLabel.cm-completionLabel-keyword": {
        color: "#7c3aed", // text-purple-600 
        fontWeight: "bold"
    },
    ".cm-tooltip-autocomplete-dark .cm-completionIcon.cm-completionIcon-keyword": {
        "&:after": { content: "'k'" },
        color: "#a78bfa", // dark:text-purple-400
    },
    ".cm-tooltip-autocomplete-dark .cm-completionLabel.cm-completionLabel-keyword": {
        color: "#a78bfa" // dark:text-purple-400
    },
    ".cm-tooltip-autocomplete .cm-completionIcon.cm-completionIcon-function": {
        "&:after": { content: "'f'" },
        color: "#0891b2", // text-cyan-600
        fontWeight: "bold"
    },
    ".cm-completionLabel.cm-completionLabel-function": {
        color: "#0891b2" // text-cyan-600
    },
    ".cm-tooltip-autocomplete-dark .cm-completionIcon.cm-completionIcon-function": {
        "&:after": { content: "'f'" },
        color: "#22d3ee", // dark:text-cyan-400
    },
    ".cm-tooltip-autocomplete-dark .cm-completionLabel.cm-completionLabel-function": {
        color: "#22d3ee" // dark:text-cyan-400
    },
    ".cm-tooltip-autocomplete .cm-completionIcon.cm-completionIcon-operator": {
        "&:after": { content: "'op'" },
        color: "#ea580c", // text-orange-600
        fontWeight: "bold"
    },
    ".cm-completionLabel.cm-completionLabel-operator": {
        color: "#ea580c" // text-orange-600
    },
    ".cm-tooltip-autocomplete-dark .cm-completionIcon.cm-completionIcon-operator": {
        "&:after": { content: "'op'" },
        color: "#fb923c", // dark:text-orange-400
    },
    ".cm-tooltip-autocomplete-dark .cm-completionLabel.cm-completionLabel-operator": {
        color: "#fb923c" // dark:text-orange-400
    },
    ".cm-variable-template": {
        color: "#0284c7", // text-sky-600
        background: "rgba(2, 132, 199, 0.1)", // bg-sky-50/10
        borderRadius: "3px",
        padding: "1px 0"
    },
    "&dark .cm-variable-template": {
        color: "#38bdf8", // dark:text-sky-400
        background: "rgba(56, 189, 248, 0.1)" // dark:bg-sky-400/10
    },
    
    // Selected item style
    ".cm-tooltip-autocomplete ul li[aria-selected=true]": {
        backgroundColor: "#e0f2fe", // bg-sky-100
        color: "#0c4a6e" // text-sky-900
    },
    ".cm-tooltip-autocomplete-dark ul li[aria-selected=true]": {
        backgroundColor: "#075985", // dark:bg-sky-800
        color: "#e0f7ff" // dark:text-sky-50
    },
    ".cm-tooltip-autocomplete .cm-completionIcon.cm-completionIcon-table": {
        "&:after": { content: "'T'" },
        color: "#059669", // text-emerald-600
        fontWeight: "bold"
    },
    ".cm-completionLabel.cm-completionLabel-table": {
        color: "#059669" // text-emerald-600
    },
    ".cm-tooltip-autocomplete-dark .cm-completionIcon.cm-completionIcon-table": {
        "&:after": { content: "'T'" },
        color: "#10b981", // dark:text-emerald-500
    },
    ".cm-tooltip-autocomplete-dark .cm-completionLabel.cm-completionLabel-table": {
        color: "#10b981" // dark:text-emerald-500
    },
    ".cm-tooltip-autocomplete .cm-completionIcon.cm-completionIcon-column": {
        "&:after": { content: "'C'" },
        color: "#0284c7", // text-sky-600
        fontWeight: "bold"
    },
    ".cm-completionLabel.cm-completionLabel-column": {
        color: "#0284c7" // text-sky-600
    },
    ".cm-tooltip-autocomplete-dark .cm-completionIcon.cm-completionIcon-column": {
        "&:after": { content: "'C'" },
        color: "#38bdf8", // dark:text-sky-400
    },
    ".cm-tooltip-autocomplete-dark .cm-completionLabel.cm-completionLabel-column": {
        color: "#38bdf8" // dark:text-sky-400
    },
    ".cm-tooltip-autocomplete .cm-completionIcon.cm-completionIcon-qualified-column": {
        "&:after": { content: "'QC'" },
        color: "#0891b2", // text-cyan-600
        fontWeight: "bold"
    },
    ".cm-completionLabel.cm-completionLabel-qualified-column": {
        color: "#0891b2" // text-cyan-600
    },
    ".cm-tooltip-autocomplete-dark .cm-completionIcon.cm-completionIcon-qualified-column": {
        "&:after": { content: "'QC'" },
        color: "#22d3ee", // dark:text-cyan-400
    },
    ".cm-tooltip-autocomplete-dark .cm-completionLabel.cm-completionLabel-qualified-column": {
        color: "#22d3ee" // dark:text-cyan-400
    },
});

// Create a decoration for highlighting variables
const createVariableHighlighter = () => {
    const variableRegex = /\{\{\s*[a-zA-Z0-9_]+\s*\}\}/g;
    
    return ViewPlugin.fromClass(class {
        decorations: DecorationSet;
        
        constructor(view: EditorView) {
            this.decorations = this.buildDecorations(view);
        }
        
        update(update: ViewUpdate) {
            // Only rebuild decorations if document or viewport changed
            if (update.docChanged || update.viewportChanged) {
                this.decorations = this.buildDecorations(update.view);
            }
        }
        
        buildDecorations(view: EditorView) {
            const decorations = [];
            
            // Only process visible ranges to improve performance
            for (const { from, to } of view.visibleRanges) {
                const text = view.state.doc.sliceString(from, to);
                const matches = [...text.matchAll(variableRegex)];
                
                for (const match of matches) {
                    const start = from + match.index!;
                    const end = start + match[0].length;
                    
                    decorations.push(Decoration.mark({
                        class: "cm-variable-template"
                    }).range(start, end));
                }
            }
            
            return Decoration.set(decorations);
        }
    }, {
        decorations: v => v.decorations
    });
};

// Cache extensions for better performance
const createExtensions = (darkMode: boolean) => {
    const exts = [
        myTheme,
        createVariableHighlighter(),
        autocompletion({
            override: [completeVariables, completeSql],
            icons: true,
            maxRenderedOptions: 100, // Limit rendered options for performance
            defaultKeymap: true,     // Use default keymap for better performance
            optionClass: option => option.type ? `cm-completion-${option.type}` : "",
            addToOptions: [
                {
                    render(completion, state) {
                        const dom = document.createElement("div");
                        dom.style.padding = "2px 0";
                        dom.style.display = "flex";
                        dom.style.alignItems = "center";
                        dom.style.opacity = "0.7";
                        dom.style.fontSize = "0.85em";
                        
                        if (completion.type === "variable" && completion.detail) {
                            dom.textContent = completion.detail.length > 20 
                                ? completion.detail.substring(0, 20) + "..." 
                                : completion.detail;
                        } else if (completion.info) {
                            dom.textContent = typeof completion.info === 'string' 
                                ? completion.info 
                                : '';
                        }
                        
                        return dom;
                    },
                    position: 80
                }
            ]
        })
    ];
    
    // Add oneDarkTheme only if in dark mode
    if (darkMode) {
        exts.unshift(oneDarkTheme);
    }
    
    return exts;
};

// Separate cache for dark and light mode
let lightModeExtensionsCache: any[] | null = null;
let darkModeExtensionsCache: any[] | null = null;
let lightModeOutputExtensionsCache: any[] | null = null;
let darkModeOutputExtensionsCache: any[] | null = null;

// Combine all extensions - with memoization
const extensions = computed(() => {
    if (isDarkMode.value) {
        if (!darkModeExtensionsCache) {
            darkModeExtensionsCache = createExtensions(true);
        }
        return darkModeExtensionsCache;
    } else {
        if (!lightModeExtensionsCache) {
            lightModeExtensionsCache = createExtensions(false);
        }
        return lightModeExtensionsCache;
    }
});

// Output editor extensions (simpler, no autocomplete) - with memoization
const outputExtensions = computed(() => {
    if (isDarkMode.value) {
        if (!darkModeOutputExtensionsCache) {
            darkModeOutputExtensionsCache = isDarkMode.value ? [oneDarkTheme, myTheme] : [myTheme];
        }
        return darkModeOutputExtensionsCache;
    } else {
        if (!lightModeOutputExtensionsCache) {
            lightModeOutputExtensionsCache = [myTheme];
        }
        return lightModeOutputExtensionsCache;
    }
});

// Optimize the SQL binding with debounce to prevent excessive computations
const debouncedMakeBindedSQL = useDebounceFn(async (sql: string, data: any[], variables: any[], minify?: boolean) => {
    // Ensure minify is always a boolean when passed to the Go function
    return await MakeBindedSQL(sql, data, variables, minify || false) ?? "";
}, 300); // 300ms debounce

const linesBinded = computedAsync(async () => {
    if(!props.showBindedSql) {
        return "";
    }

    return await debouncedMakeBindedSQL(value.value, props.data!, props.variables!, props.minify);
}, "");

const getBindedSQL = async (): Promise<string> => {
    if(!props.showBindedSql) {
        return "";
    }

    return await MakeBindedSQL(value.value, props.data!, props.variables!, props.minify) ?? "";
}

watch(() => props.modelValue, (newValue) => {
    value.value = newValue;
});

watch(() => value.value, (val) => {
    emit('update:modelValue', val);
});

// Manual toggle for dark mode (for testing and debugging)
const toggleDarkMode = () => {
    isDarkMode.value = !isDarkMode.value;
    console.log("Manual dark mode toggle:", isDarkMode.value);
};

// Watch for dark mode changes
watch(() => isDarkMode.value, (newValue) => {
    // The computed extensions will automatically update
    console.log("Dark mode changed to:", newValue);
});

defineExpose({
    getBindedSQL,
    toggleDarkMode,
    isDarkMode
});

// Debounced dark mode checking to prevent excessive DOM operations
const updateDarkMode = useDebounceFn(() => {
    // Check for .dark class on html or body
    const hasDarkClass = 
        document.documentElement.classList.contains('dark') || 
        document.body.classList.contains('dark');
    
    // Check for prefers-color-scheme media query if available
    const prefersDark = window.matchMedia && 
        window.matchMedia('(prefers-color-scheme: dark)').matches;
    
    // Check for data-theme="dark" attribute
    const hasDataThemeDark = 
        document.documentElement.getAttribute('data-theme') === 'dark' || 
        document.body.getAttribute('data-theme') === 'dark';
    
    // Check for any custom dark mode implementation (like .tw-dark)
    const hasCustomDarkClass = 
        document.documentElement.classList.contains('tw-dark') || 
        document.body.classList.contains('tw-dark');
    
    // Determine final dark mode state based on checks
    const isDark = hasDarkClass || hasDataThemeDark || hasCustomDarkClass || prefersDark;
    
    // Avoid unnecessary updates if the dark mode state hasn't changed
    if (isDarkMode.value !== isDark) {
        isDarkMode.value = isDark;
    }
}, 100); // 100ms debounce is sufficient for UI theme changes

// Create a completion source for variables
const variableCompletions = computed(() => {
    if (!props.variables) return [];
    
    return props.variables.map(variable => ({
        label: `{{ ${variable.Value} }}`,
        value: variable.Value
    }));
});

// Toggle SQL Assistant
const toggleSQLAssistant = () => {
  showSQLAssistant.value = !showSQLAssistant.value;
  // Close other menus when opening the assistant
  showSnippets.value = false;
  showHistory.value = false;
  showVariables.value = false;
};

// Handle SQL generated by the assistant
const handleSQLGenerated = (sql: string) => {
  // Insert the generated SQL into the editor
  value.value = sql;
  // Close the assistant
  showSQLAssistant.value = false;
};

</script>

<style scoped>
/* Add your styles here */
</style>
