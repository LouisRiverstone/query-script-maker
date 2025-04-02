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
                                    class="px-2 py-1 text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700 text-sm"
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
                        </div>
                    </div>
                </div>
                
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
import { ref, watch, onMounted, computed, onUnmounted } from 'vue';
import CodeMirror from 'vue-codemirror6';
import { sql, MySQL } from "@codemirror/lang-sql";
import { oneDarkTheme } from '@codemirror/theme-one-dark';
import { MakeBindedSQL, GetLatestDatabaseStructure } from '../../wailsjs/go/main/App';
import { main } from '../../wailsjs/go/models';
import { computedAsync } from '@vueuse/core';
import { autocompletion, CompletionContext, CompletionResult } from '@codemirror/autocomplete';
import { EditorView, Decoration, ViewPlugin, ViewUpdate, DecorationSet } from '@codemirror/view';

import Divider from './Divider.vue';

const variableText = ref("{{ variable }}");
const isDarkMode = ref(false);
let darkModeObserver: MutationObserver;
const copySuccess = ref(false);
const showSnippets = ref(false);
const showHistory = ref(false);
const showVariables = ref(false);

// Array para armazenar o histórico de consultas
interface QueryHistoryItem {
  text: string;
  timestamp: number;
}

const HISTORY_STORAGE_KEY = 'sql_query_history';
const queryHistory = ref<QueryHistoryItem[]>([]);

// Carregar histórico do localStorage
const loadQueryHistory = () => {
  try {
    const savedHistory = localStorage.getItem(HISTORY_STORAGE_KEY);
    if (savedHistory) {
      queryHistory.value = JSON.parse(savedHistory);
    }
  } catch (error) {
    console.error("Error loading query history:", error);
    queryHistory.value = [];
  }
};

// Salvar consulta atual no histórico
const saveToHistory = () => {
  if (value.value.trim()) {
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
      
      queryHistory.value = [newItem, ...queryHistory.value].slice(0, 20); // Manter no máximo 20 itens
      localStorage.setItem(HISTORY_STORAGE_KEY, JSON.stringify(queryHistory.value));
    }
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

// Variables for table and column autocomplete
const databaseStructure = ref<DatabaseStructure | null>(null);
const tableCompletions = ref<any[]>([]);
const columnCompletions = ref<any[]>([]);

// Load database structure on component mount
const loadDatabaseStructure = async () => {
  try {
    const structure = await GetLatestDatabaseStructure();
    if (structure) {
      databaseStructure.value = JSON.parse(structure) as DatabaseStructure;
      
      // Process structure to create completions
      if (databaseStructure.value && databaseStructure.value.tables) {
        // Create table completions
        tableCompletions.value = databaseStructure.value.tables.map(table => ({
          label: table.name,
          type: "table",
          detail: `Table with ${table.columns.length} columns`,
          info: `Table: ${table.name}`,
          boost: 99
        }));
        
        // Create column completions grouped by table
        const columns: any[] = [];
        databaseStructure.value.tables.forEach(table => {
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
            
            // Also add tablename.columnname format
            columns.push({
              label: `${table.name}.${column.name}`,
              type: "qualified-column",
              detail: `${table.name}.${column.name}`,
              info: `Column: ${column.name} (${column.type})${column.isPrimary ? ' [PK]' : ''}`,
              apply: `${table.name}.${column.name}`,
              boost: column.isPrimary ? 85 : 80
            });
          });
        });
        columnCompletions.value = columns;
      }
    }
  } catch (error) {
    console.error("Error loading database structure:", error);
  }
};

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
    
    // Watch for changes in dark mode
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
    window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', updateDarkMode);
    
    // Load database structure
    loadDatabaseStructure();
    
    // Load query history
    loadQueryHistory();
    
    // Close menus when clicking outside
    document.addEventListener('click', handleClickOutside);
    
    // Configure global keyboard shortcuts
    document.addEventListener('keydown', handleKeyboardShortcuts);
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

// SQL keywords for autocomplete
const sqlKeywords = [
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
  
  // Operators
  { label: "AND", type: "operator", info: "Logical AND operator" },
  { label: "OR", type: "operator", info: "Logical OR operator" },
  { label: "NOT", type: "operator", info: "Logical NOT operator" },
  { label: "IN", type: "operator", info: "Specifies multiple possible values for a column" },
  { label: "NOT IN", type: "operator", info: "Specifies values that are not in a list of values" },
  { label: "BETWEEN", type: "operator", info: "Selects values within a given range" },
  { label: "NOT BETWEEN", type: "operator", info: "Selects values outside a given range" },
  { label: "LIKE", type: "operator", info: "Searches for a specified pattern in a column" },
  { label: "NOT LIKE", type: "operator", info: "Searches for values that don't match a pattern" },
  { label: "REGEXP", type: "operator", info: "Matches a string against a regular expression pattern" },
  { label: "IS NULL", type: "operator", info: "Tests for NULL values" },
  { label: "IS NOT NULL", type: "operator", info: "Tests for non-NULL values" },
  { label: "EXISTS", type: "operator", info: "Tests for the existence of records in a subquery" },
  { label: "ANY", type: "operator", info: "Compares a value to any value in a list" },
  { label: "ALL", type: "operator", info: "Compares a value to all values in a list" },
  { label: "SOME", type: "operator", info: "Same as ANY" },
  
  // Aggregate Functions
  { label: "COUNT()", type: "function", info: "Returns the number of rows" },
  { label: "COUNT(DISTINCT)", type: "function", info: "Returns the number of distinct values" },
  { label: "SUM()", type: "function", info: "Returns the sum of values" },
  { label: "AVG()", type: "function", info: "Returns the average value" },
  { label: "MIN()", type: "function", info: "Returns the minimum value" },
  { label: "MAX()", type: "function", info: "Returns the maximum value" },
  { label: "GROUP_CONCAT()", type: "function", info: "Returns a concatenated string of values" },
  { label: "STD()", type: "function", info: "Returns the standard deviation" },
  { label: "STDDEV()", type: "function", info: "Returns the standard deviation" },
  { label: "STDDEV_POP()", type: "function", info: "Returns the population standard deviation" },
  { label: "STDDEV_SAMP()", type: "function", info: "Returns the sample standard deviation" },
  { label: "VAR_POP()", type: "function", info: "Returns the population variance" },
  { label: "VAR_SAMP()", type: "function", info: "Returns the sample variance" },
  { label: "VARIANCE()", type: "function", info: "Returns the variance" },
  { label: "BIT_AND()", type: "function", info: "Returns the bitwise AND of all bits in a group" },
  { label: "BIT_OR()", type: "function", info: "Returns the bitwise OR of all bits in a group" },
  { label: "BIT_XOR()", type: "function", info: "Returns the bitwise XOR of all bits in a group" },
  { label: "JSON_ARRAYAGG()", type: "function", info: "Returns a JSON array containing values from a group" },
  { label: "JSON_OBJECTAGG()", type: "function", info: "Returns a JSON object containing key-value pairs from a group" },
  
  // String Functions
  { label: "CONCAT()", type: "function", info: "Adds two or more strings together" },
  { label: "CONCAT_WS()", type: "function", info: "Adds two or more strings together with a separator" },
  { label: "SUBSTRING()", type: "function", info: "Extracts a string of characters from a string" },
  { label: "SUBSTR()", type: "function", info: "Extracts a substring from a string (alias for SUBSTRING)" },
  { label: "TRIM()", type: "function", info: "Removes leading and trailing spaces from a string" },
  { label: "LTRIM()", type: "function", info: "Removes leading spaces from a string" },
  { label: "RTRIM()", type: "function", info: "Removes trailing spaces from a string" },
  { label: "UPPER()", type: "function", info: "Converts a string to upper case" },
  { label: "LOWER()", type: "function", info: "Converts a string to lower case" },
  { label: "LCASE()", type: "function", info: "Converts a string to lower case (alias for LOWER)" },
  { label: "UCASE()", type: "function", info: "Converts a string to upper case (alias for UPPER)" },
  { label: "LENGTH()", type: "function", info: "Returns the length of a string" },
  { label: "CHAR_LENGTH()", type: "function", info: "Returns the number of characters in a string" },
  { label: "CHARACTER_LENGTH()", type: "function", info: "Returns the number of characters in a string" },
  { label: "REPLACE()", type: "function", info: "Replaces all occurrences of a substring within a string" },
  { label: "REVERSE()", type: "function", info: "Reverses a string" },
  { label: "REPEAT()", type: "function", info: "Repeats a string a specified number of times" },
  { label: "INSERT()", type: "function", info: "Inserts a substring at the specified position" },
  { label: "LEFT()", type: "function", info: "Returns the leftmost characters from a string" },
  { label: "RIGHT()", type: "function", info: "Returns the rightmost characters from a string" },
  { label: "LOCATE()", type: "function", info: "Returns the position of a substring in a string" },
  { label: "POSITION()", type: "function", info: "Returns the position of a substring in a string" },
  { label: "INSTR()", type: "function", info: "Returns the position of a substring in a string" },
  { label: "LPAD()", type: "function", info: "Left-pads a string with another string" },
  { label: "RPAD()", type: "function", info: "Right-pads a string with another string" },
  { label: "SPACE()", type: "function", info: "Returns a string of the specified number of spaces" },
  { label: "ELT()", type: "function", info: "Returns the string at the specified position" },
  { label: "FIELD()", type: "function", info: "Returns the index position of a value in a list" },
  { label: "FIND_IN_SET()", type: "function", info: "Returns the position of a string in a comma-separated list" },
  { label: "FORMAT()", type: "function", info: "Formats a number to a format like '#,###,###.##'" },
  { label: "HEX()", type: "function", info: "Converts a value to hexadecimal" },
  { label: "UNHEX()", type: "function", info: "Converts a hexadecimal value to a string" },
  { label: "BIN()", type: "function", info: "Returns a binary representation of a number" },
  { label: "OCT()", type: "function", info: "Returns an octal representation of a number" },
  { label: "ASCII()", type: "function", info: "Returns the ASCII value of the leftmost character" },
  { label: "ORD()", type: "function", info: "Returns the character code for the leftmost character" },
  { label: "SOUNDEX()", type: "function", info: "Returns a soundex string for a given string" },
  
  // Date Functions
  { label: "NOW()", type: "function", info: "Returns the current date and time" },
  { label: "SYSDATE()", type: "function", info: "Returns the current date and time" },
  { label: "CURDATE()", type: "function", info: "Returns the current date" },
  { label: "CURRENT_DATE()", type: "function", info: "Returns the current date" },
  { label: "CURTIME()", type: "function", info: "Returns the current time" },
  { label: "CURRENT_TIME()", type: "function", info: "Returns the current time" },
  { label: "CURRENT_TIMESTAMP()", type: "function", info: "Returns the current date and time" },
  { label: "LOCALTIME()", type: "function", info: "Returns the current date and time" },
  { label: "LOCALTIMESTAMP()", type: "function", info: "Returns the current date and time" },
  { label: "UTC_DATE()", type: "function", info: "Returns the current UTC date" },
  { label: "UTC_TIME()", type: "function", info: "Returns the current UTC time" },
  { label: "UTC_TIMESTAMP()", type: "function", info: "Returns the current UTC date and time" },
  { label: "YEAR()", type: "function", info: "Returns the year part of a date" },
  { label: "MONTH()", type: "function", info: "Returns the month part of a date" },
  { label: "DAY()", type: "function", info: "Returns the day part of a date" },
  { label: "HOUR()", type: "function", info: "Returns the hour part of a time" },
  { label: "MINUTE()", type: "function", info: "Returns the minute part of a time" },
  { label: "SECOND()", type: "function", info: "Returns the second part of a time" },
  { label: "MICROSECOND()", type: "function", info: "Returns the microsecond part of a time" },
  { label: "DATE_FORMAT()", type: "function", info: "Formats a date as specified" },
  { label: "TIME_FORMAT()", type: "function", info: "Formats a time as specified" },
  { label: "DATE_ADD()", type: "function", info: "Adds a time/date interval to a date" },
  { label: "ADDDATE()", type: "function", info: "Adds a time/date interval to a date" },
  { label: "DATE_SUB()", type: "function", info: "Subtracts a time/date interval from a date" },
  { label: "SUBDATE()", type: "function", info: "Subtracts a time/date interval from a date" },
  { label: "ADDTIME()", type: "function", info: "Adds a time interval to a time" },
  { label: "SUBTIME()", type: "function", info: "Subtracts a time interval from a time" },
  { label: "DATEDIFF()", type: "function", info: "Returns the difference in days between two dates" },
  { label: "TIMEDIFF()", type: "function", info: "Returns the difference between two times" },
  { label: "DATE()", type: "function", info: "Extracts the date part of a date or datetime" },
  { label: "TIME()", type: "function", info: "Extracts the time part of a time or datetime" },
  { label: "TIMESTAMP()", type: "function", info: "Returns a datetime value" },
  { label: "CONVERT_TZ()", type: "function", info: "Converts a datetime from one timezone to another" },
  { label: "EXTRACT()", type: "function", info: "Extracts a part of a date" },
  { label: "MAKEDATE()", type: "function", info: "Creates a date from a year and day value" },
  { label: "MAKETIME()", type: "function", info: "Creates a time from hour, minute, and second values" },
  { label: "TO_DAYS()", type: "function", info: "Converts a date to the number of days since year 0" },
  { label: "FROM_DAYS()", type: "function", info: "Converts a day number to a date" },
  { label: "TO_SECONDS()", type: "function", info: "Converts a date to the number of seconds since year 0" },
  { label: "FROM_UNIXTIME()", type: "function", info: "Converts a Unix timestamp to a date" },
  { label: "UNIX_TIMESTAMP()", type: "function", info: "Returns the Unix timestamp for a date" },
  { label: "SEC_TO_TIME()", type: "function", info: "Converts seconds to a time" },
  { label: "TIME_TO_SEC()", type: "function", info: "Converts a time to seconds" },
  { label: "DAYOFWEEK()", type: "function", info: "Returns the day of the week for a date" },
  { label: "WEEKDAY()", type: "function", info: "Returns the weekday index for a date" },
  { label: "DAYOFMONTH()", type: "function", info: "Returns the day of the month for a date" },
  { label: "DAYOFYEAR()", type: "function", info: "Returns the day of the year for a date" },
  { label: "DAYNAME()", type: "function", info: "Returns the name of the day for a date" },
  { label: "MONTHNAME()", type: "function", info: "Returns the name of the month for a date" },
  { label: "QUARTER()", type: "function", info: "Returns the quarter for a date" },
  { label: "WEEK()", type: "function", info: "Returns the week number for a date" },
  { label: "WEEKOFYEAR()", type: "function", info: "Returns the week number for a date" },
  { label: "YEARWEEK()", type: "function", info: "Returns year and week for a date" },
  { label: "LAST_DAY()", type: "function", info: "Returns the last day of the month for a date" },
  
  // Control Flow Functions
  { label: "CASE", type: "keyword", info: "Evaluates a list of conditions and returns a value when the first condition is met" },
  { label: "WHEN", type: "keyword", info: "Used in a CASE statement to specify a condition" },
  { label: "THEN", type: "keyword", info: "Used in a CASE statement to specify a result" },
  { label: "ELSE", type: "keyword", info: "Used in a CASE statement to specify a value to return if all conditions are false" },
  { label: "END", type: "keyword", info: "Used in a CASE statement to end the list of conditions" },
  { label: "IF()", type: "function", info: "Returns one value if a condition is TRUE, or another value if a condition is FALSE" },
  { label: "IFNULL()", type: "function", info: "Returns a specified value if the expression is NULL" },
  { label: "NULLIF()", type: "function", info: "Returns NULL if two expressions are equal" },
  { label: "COALESCE()", type: "function", info: "Returns the first non-NULL value in a list" },
  
  // Mathematical Functions
  { label: "ABS()", type: "function", info: "Returns the absolute value of a number" },
  { label: "ACOS()", type: "function", info: "Returns the arc cosine of a number" },
  { label: "ASIN()", type: "function", info: "Returns the arc sine of a number" },
  { label: "ATAN()", type: "function", info: "Returns the arc tangent of a number" },
  { label: "ATAN2()", type: "function", info: "Returns the arc tangent of two numbers" },
  { label: "CEIL()", type: "function", info: "Returns the smallest integer value greater than or equal to a number" },
  { label: "CEILING()", type: "function", info: "Returns the smallest integer value greater than or equal to a number" },
  { label: "COS()", type: "function", info: "Returns the cosine of a number" },
  { label: "COT()", type: "function", info: "Returns the cotangent of a number" },
  { label: "DEGREES()", type: "function", info: "Converts a value in radians to degrees" },
  { label: "EXP()", type: "function", info: "Returns e raised to the power of a number" },
  { label: "FLOOR()", type: "function", info: "Returns the largest integer value less than or equal to a number" },
  { label: "LN()", type: "function", info: "Returns the natural logarithm of a number" },
  { label: "LOG()", type: "function", info: "Returns the natural logarithm of a number, or the logarithm of a number to a specified base" },
  { label: "LOG10()", type: "function", info: "Returns the base-10 logarithm of a number" },
  { label: "LOG2()", type: "function", info: "Returns the base-2 logarithm of a number" },
  { label: "MOD()", type: "function", info: "Returns the remainder of a number divided by another number" },
  { label: "PI()", type: "function", info: "Returns the value of PI" },
  { label: "POW()", type: "function", info: "Returns the value of a number raised to the power of another number" },
  { label: "POWER()", type: "function", info: "Returns the value of a number raised to the power of another number" },
  { label: "RADIANS()", type: "function", info: "Converts a value in degrees to radians" },
  { label: "RAND()", type: "function", info: "Returns a random floating-point value between 0 and 1" },
  { label: "ROUND()", type: "function", info: "Rounds a number to a specified number of decimal places" },
  { label: "SIGN()", type: "function", info: "Returns the sign of a number" },
  { label: "SIN()", type: "function", info: "Returns the sine of a number" },
  { label: "SQRT()", type: "function", info: "Returns the square root of a number" },
  { label: "TAN()", type: "function", info: "Returns the tangent of a number" },
  { label: "TRUNCATE()", type: "function", info: "Truncates a number to a specified number of decimal places" },
  
  // Window Functions
  { label: "ROW_NUMBER()", type: "function", info: "Returns the row number of the current row" },
  { label: "RANK()", type: "function", info: "Returns the rank of the current row" },
  { label: "DENSE_RANK()", type: "function", info: "Returns the dense rank of the current row" },
  { label: "NTILE()", type: "function", info: "Returns the ntile group number of the current row" },
  { label: "LAG()", type: "function", info: "Returns the value of the expression evaluated at the row previous to the current row" },
  { label: "LEAD()", type: "function", info: "Returns the value of the expression evaluated at the row following the current row" },
  { label: "FIRST_VALUE()", type: "function", info: "Returns the value of the expression evaluated at the first row" },
  { label: "LAST_VALUE()", type: "function", info: "Returns the value of the expression evaluated at the last row" },
  { label: "NTH_VALUE()", type: "function", info: "Returns the value of the expression evaluated at the nth row" },
  { label: "PERCENT_RANK()", type: "function", info: "Returns the percent rank of the current row" },
  { label: "CUME_DIST()", type: "function", info: "Returns the cumulative distribution of the current row" },
  
  // JSON Functions (MySQL 5.7+)
  { label: "JSON_ARRAY()", type: "function", info: "Creates a JSON array" },
  { label: "JSON_OBJECT()", type: "function", info: "Creates a JSON object" },
  { label: "JSON_QUOTE()", type: "function", info: "Quotes a string as a JSON value" },
  { label: "JSON_CONTAINS()", type: "function", info: "Returns whether a JSON document contains a specific value" },
  { label: "JSON_CONTAINS_PATH()", type: "function", info: "Returns whether a JSON document contains a specific path" },
  { label: "JSON_EXTRACT()", type: "function", info: "Extracts a value from a JSON document" },
  { label: "JSON_KEYS()", type: "function", info: "Returns the keys from a JSON object" },
  { label: "JSON_SEARCH()", type: "function", info: "Searches a JSON document for a value" },
  { label: "JSON_ARRAY_APPEND()", type: "function", info: "Appends a value to a JSON array" },
  { label: "JSON_ARRAY_INSERT()", type: "function", info: "Inserts a value into a JSON array" },
  { label: "JSON_INSERT()", type: "function", info: "Inserts values into a JSON document" },
  { label: "JSON_REPLACE()", type: "function", info: "Replaces values in a JSON document" },
  { label: "JSON_REMOVE()", type: "function", info: "Removes values from a JSON document" },
  { label: "JSON_SET()", type: "function", info: "Sets values in a JSON document" },
  { label: "JSON_MERGE()", type: "function", info: "Merges JSON documents" },
  { label: "JSON_MERGE_PATCH()", type: "function", info: "Merges JSON documents using JSON Merge Patch" },
  { label: "JSON_MERGE_PRESERVE()", type: "function", info: "Merges JSON documents preserving duplicate keys" },
  { label: "JSON_TYPE()", type: "function", info: "Returns the type of a JSON value" },
  { label: "JSON_VALID()", type: "function", info: "Returns whether a value is valid JSON" },
  { label: "JSON_DEPTH()", type: "function", info: "Returns the maximum depth of a JSON document" },
  { label: "JSON_LENGTH()", type: "function", info: "Returns the length of a JSON document" },
  { label: "JSON_PRETTY()", type: "function", info: "Formats a JSON document for readability" },
  { label: "JSON_STORAGE_SIZE()", type: "function", info: "Returns the storage size of a JSON document" },
  { label: "JSON_TABLE()", type: "function", info: "Returns a relational table from JSON data" },
  { label: "JSON_UNQUOTE()", type: "function", info: "Unquotes a JSON value" },
  
  // MySQL 8.0+ Functions
  { label: "GROUPING()", type: "function", info: "Indicates whether a specified column expression in a GROUP BY clause is aggregated" },
  { label: "LATERAL", type: "keyword", info: "Used with derived tables to refer to preceding tables in the FROM clause" },
  { label: "OVER()", type: "function", info: "Defines a window for a window function" },
  { label: "PARTITION BY", type: "keyword", info: "Divides the result set into partitions" },
  { label: "WITH ROLLUP", type: "keyword", info: "Adds extra rows to the result set of GROUP BY to represent subtotals" },
  { label: "RECURSIVE", type: "keyword", info: "Used in a CTE to define a recursive query" },
  
  // Deprecated functions (included for compatibility)
  { label: "DATABASE()", type: "function", info: "Returns the name of the current database" },
  { label: "SCHEMA()", type: "function", info: "Returns the name of the current database (alias for DATABASE)" },
  { label: "USER()", type: "function", info: "Returns the current MySQL user name and host name" },
  { label: "VERSION()", type: "function", info: "Returns the current version of the MySQL server" },
  { label: "CURRENT_USER()", type: "function", info: "Returns the user name and host name for the MySQL account" },
  { label: "LAST_INSERT_ID()", type: "function", info: "Returns the AUTO_INCREMENT value generated by the last INSERT statement" },
  { label: "PASSWORD()", type: "function", info: "Calculates and returns a password string" },
  { label: "BENCHMARK()", type: "function", info: "Executes an expression repeatedly" },
  { label: "CONVERT()", type: "function", info: "Converts a value to a different data type" },
  { label: "CAST()", type: "function", info: "Converts a value to a different data type" },
];

// Create a completion source for variables
const variableCompletions = computed(() => {
    if (!props.variables) return [];
    
    return props.variables.map(variable => ({
        label: `{{ ${variable.Value} }}`,
        type: "variable",
        detail: `Field from .xlsx`,
        info: `Field from .xlsx: ${variable.Field}`,
        apply: `{{ ${variable.Value} }}`,
        boost: 99, // Give variables higher priority
    }));
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

// Autocomplete function for SQL
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
    
    // First check for table names after FROM or JOIN
    const textBefore = context.state.doc.sliceString(0, context.pos).toLowerCase();
    const lastFromOrJoin = Math.max(
        textBefore.lastIndexOf("from "), 
        textBefore.lastIndexOf("join ")
    );
    const lastComma = textBefore.lastIndexOf(",", context.pos);
    
    // If we're after FROM or JOIN, prioritize table completions
    if (lastFromOrJoin > -1 && lastFromOrJoin > lastComma) {
        // Only filter for tables that match the current word
        const filteredTables = tableCompletions.value.filter(table => 
            table.label.toLowerCase().startsWith(word.text.toLowerCase())
        );
        
        if (filteredTables.length > 0 || context.explicit) {
            return {
                from: word.from,
                options: filteredTables,
                validFor: /\w+/,
            };
        }
    }
    
    // Check for column completions
    // This should work after SELECT, WHERE, ORDER BY, etc.
    // Look for patterns like "table." to provide columns for that specific table
    const periodMatch = context.matchBefore(/(\w+)\.\w*/);
    if (periodMatch) {
        const tableName = periodMatch.text.split('.')[0];
        const columnWord = context.matchBefore(/\w*$/);
        
        if (columnWord) {
            // Get columns for the specific table
            const tableSpecificColumns = columnCompletions.value.filter(column => 
                column.table === tableName && 
                column.type === "column" && 
                column.label.toLowerCase().startsWith(columnWord.text.toLowerCase())
            );
            
            if (tableSpecificColumns.length > 0 || context.explicit) {
                return {
                    from: columnWord.from,
                    options: tableSpecificColumns,
                    validFor: /\w+/,
                };
            }
        }
    }
    
    // Default to showing all columns and SQL keywords
    // Combine SQL keywords with column completions
    const allOptions = [...sqlKeywords, ...columnCompletions.value];
    
    // Filter options based on input
    const filteredOptions = allOptions.filter(option => 
        option.label.toLowerCase().startsWith(word.text.toLowerCase())
    );
    
    if (filteredOptions.length === 0 && !context.explicit) {
        return null;
    }
    
    return {
        from: word.from,
        options: filteredOptions,
        validFor: /\w+/,
    };
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
            if (update.docChanged || update.viewportChanged) {
                this.decorations = this.buildDecorations(update.view);
            }
        }
        
        buildDecorations(view: EditorView) {
            const decorations = [];
            
            for (const { from, to } of view.visibleRanges) {
                const text = view.state.doc.sliceString(from, to);
                const matches = [...text.matchAll(variableRegex)];
                
                for (const match of matches) {
                    const start = from + match.index!;
                    const end = start + match[0].length;
                    
                    // Improved variable highlighting with a more visible style
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

// Combine all extensions
const extensions = computed(() => {
    const exts = [
        myTheme,
        createVariableHighlighter(),
        autocompletion({
            override: [completeVariables, completeSql],
            icons: true,
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
                        
                        // Show different info based on completion type
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
    if (isDarkMode.value) {
        exts.unshift(oneDarkTheme);
    }
    
    return exts;
});

// Output editor extensions (simpler, no autocomplete)
const outputExtensions = computed(() => {
    const exts = [myTheme];
    
    // Add oneDarkTheme only if in dark mode
    if (isDarkMode.value) {
        exts.unshift(oneDarkTheme);
    }
    
    return exts;
});

const linesBinded = computedAsync(async () => {
    if(!props.showBindedSql) {
        return "";
    }

    return await MakeBindedSQL(value.value, props.data!, props.variables!, props.minify) ?? "";
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

// Check if dark mode is enabled
const updateDarkMode = () => {
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
    // Priority: .dark class > data-theme > custom class > media query
    const isDark = hasDarkClass || hasDataThemeDark || hasCustomDarkClass || prefersDark;
    
    // Log the determination factors (for debugging)
    console.log("Dark mode detection:", { 
        hasDarkClass, 
        prefersDark, 
        hasDataThemeDark, 
        hasCustomDarkClass,
        result: isDark
    });
    
    isDarkMode.value = isDark;
};
</script>

<style scoped>
.editor-dark-mode :deep(.cm-editor) {
    background-color: #1f2937; /* dark:bg-gray-800 */
    color: #f9fafb; /* dark:text-gray-100 */
}

/* Custom dark styles that can be applied directly */
.editor-dark-mode :deep(.cm-gutters) {
    background-color: #111827; /* dark:bg-gray-900 */
    color: #6b7280; /* dark:text-gray-500 */
    border-right-color: #374151; /* dark:border-gray-700 */
}

.editor-dark-mode :deep(.cm-activeLineGutter) {
    background-color: #374151; /* dark:bg-gray-700 */
    color: #9ca3af; /* dark:text-gray-400 */
}

.editor-dark-mode :deep(.cm-activeLine) {
    background-color: rgba(55, 65, 81, 0.3); /* dark:bg-gray-700/30 */
}

.editor-dark-mode :deep(.cm-selectionMatch) {
    background-color: rgba(37, 99, 235, 0.2); /* dark:bg-blue-600/20 */
}

.editor-dark-mode :deep(.cm-tooltip) {
    background: #1f2937; /* dark:bg-gray-800 */
    border: 1px solid #374151; /* dark:border-gray-700 */
    color: #f9fafb; /* dark:text-gray-100 */
}

/* Enhanced styling for variables */
:deep(.cm-variable-template) {
    color: #0284c7; /* text-sky-600 */
    background: rgba(2, 132, 199, 0.15); /* bg-sky-50/15 - more visible */
    border-radius: 4px;
    padding: 1px 2px;
    margin: 0 1px;
    font-weight: 500;
    box-shadow: 0 0 0 1px rgba(2, 132, 199, 0.2);
}

.editor-dark-mode :deep(.cm-variable-template) {
    color: #38bdf8; /* dark:text-sky-400 */
    background: rgba(56, 189, 248, 0.15); /* dark:bg-sky-400/15 */
    box-shadow: 0 0 0 1px rgba(56, 189, 248, 0.3);
}

/* Editor style */
:deep(.cm-editor) {
    border-radius: 0.375rem;
    overflow: hidden;
    transition: all 0.2s;
    height: var(--editor-height, auto);
    min-height: 100px;
    resize: vertical; /* Allow vertical resizing */
}

:deep(.cm-editor:focus-within) {
    box-shadow: 0 0 0 2px rgba(2, 132, 199, 0.2);
}

.editor-dark-mode :deep(.cm-editor:focus-within) {
    box-shadow: 0 0 0 2px rgba(56, 189, 248, 0.2);
}

/* Style for autocomplete tooltips */
:deep(.cm-tooltip.cm-tooltip-autocomplete) {
    box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
    border-radius: 6px;
    padding: 4px 0;
    max-height: 300px;
    overflow-y: auto;
}

:deep(.cm-tooltip-autocomplete ul li) {
    padding: 4px 8px;
    border-radius: 4px;
    margin: 0 4px;
}

:deep(.cm-completionMatchedText) {
    text-decoration: none;
    font-weight: 500;
    color: #0284c7; /* text-sky-600 */
}

.editor-dark-mode :deep(.cm-completionMatchedText) {
    color: #38bdf8; /* dark:text-sky-400 */
}

/* Visual indication of resize capability */
:deep(.cm-editor::after) {
    content: "";
    position: absolute;
    bottom: 0;
    left: 0;
    right: 0;
    height: 5px;
    background: linear-gradient(to bottom, transparent, rgba(0, 0, 0, 0.05));
    pointer-events: none;
    opacity: 0;
    transition: opacity 0.2s;
}

:deep(.cm-editor:hover::after) {
    opacity: 1;
}

.editor-dark-mode :deep(.cm-editor::after) {
    background: linear-gradient(to bottom, transparent, rgba(255, 255, 255, 0.05));
}
</style>