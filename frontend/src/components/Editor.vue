<template>
    <div class="flex flex-col gap-3 w-full" :class="{ 'editor-dark-mode': isDarkMode }">
        <slot></slot>
        <div class="flex flex-col lg:flex-row gap-6">
            <div class="flex flex-col w-full gap-3">
                <div class="flex justify-between items-center">
                    <Divider v-if="showBindedSql">Input</Divider>
                    <button 
                        @click="toggleDarkMode" 
                        title="Toggle dark mode (debug only)"
                        class="ml-auto text-xs text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 px-2 py-0.5 rounded"
                    >
                        {{ isDarkMode ? '☀️' : '🌙' }}
                    </button>
                </div>
                <small class="text-gray-600 dark:text-gray-400 text-xs italic mb-1 px-1">To use variables, use: {{ variableText }}</small>
                <code-mirror v-model="value" :lang="lang" :extensions="extensions" :linter="null" basic wrap tab class="w-full" />
            </div>
            <div v-if="showBindedSql" class="flex flex-col w-full gap-3">
                <Divider>Output</Divider>
                <code-mirror v-model="linesBinded" :lang="lang" :extensions="outputExtensions" :linter="null" basic wrap tab class="w-full" />
            </div>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { ref, watch, onMounted, computed, onUnmounted } from 'vue';
import CodeMirror from 'vue-codemirror6';
import { sql, MySQL } from "@codemirror/lang-sql";
import { oneDarkTheme } from '@codemirror/theme-one-dark';
import { MakeBindedSQL } from '../../wailsjs/go/main/App';
import { main } from '../../wailsjs/go/models';
import { computedAsync } from '@vueuse/core';
import { autocompletion, CompletionContext, CompletionResult } from '@codemirror/autocomplete';
import { EditorView, Decoration, ViewPlugin, ViewUpdate, DecorationSet } from '@codemirror/view';

import Divider from './Divider.vue';

const variableText = ref("{{ variable }}");
const isDarkMode = ref(false);
let darkModeObserver: MutationObserver;

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
});

onUnmounted(() => {
    // Clean up observers and event listeners
    if (darkModeObserver) {
        darkModeObserver.disconnect();
    }
    
    window.matchMedia('(prefers-color-scheme: dark)').removeEventListener('change', updateDarkMode);
});

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
  { label: "JOIN", type: "keyword", info: "Combines rows from two or more tables" },
  { label: "LEFT JOIN", type: "keyword", info: "Returns all records from the left table, and matched records from the right table" },
  { label: "RIGHT JOIN", type: "keyword", info: "Returns all records from the right table, and matched records from the left table" },
  { label: "INNER JOIN", type: "keyword", info: "Returns records that have matching values in both tables" },
  { label: "FULL JOIN", type: "keyword", info: "Returns all records when there is a match in either left or right table" },
  { label: "GROUP BY", type: "keyword", info: "Groups rows that have the same values into summary rows" },
  { label: "HAVING", type: "keyword", info: "Filters records after GROUP BY is applied" },
  { label: "ORDER BY", type: "keyword", info: "Sorts the result set in ascending or descending order" },
  { label: "LIMIT", type: "keyword", info: "Limits the number of records returned" },
  
  // Operators
  { label: "AND", type: "operator", info: "Logical AND operator" },
  { label: "OR", type: "operator", info: "Logical OR operator" },
  { label: "NOT", type: "operator", info: "Logical NOT operator" },
  { label: "IN", type: "operator", info: "Specifies multiple possible values for a column" },
  { label: "BETWEEN", type: "operator", info: "Selects values within a given range" },
  { label: "LIKE", type: "operator", info: "Searches for a specified pattern in a column" },
  { label: "IS NULL", type: "operator", info: "Tests for NULL values" },
  { label: "IS NOT NULL", type: "operator", info: "Tests for non-NULL values" },
  
  // Aggregate Functions
  { label: "COUNT()", type: "function", info: "Returns the number of rows" },
  { label: "SUM()", type: "function", info: "Returns the sum of values" },
  { label: "AVG()", type: "function", info: "Returns the average value" },
  { label: "MIN()", type: "function", info: "Returns the minimum value" },
  { label: "MAX()", type: "function", info: "Returns the maximum value" },
  
  // String Functions
  { label: "CONCAT()", type: "function", info: "Adds two or more strings together" },
  { label: "SUBSTRING()", type: "function", info: "Extracts a string of characters from a string" },
  { label: "TRIM()", type: "function", info: "Removes leading and trailing spaces from a string" },
  { label: "UPPER()", type: "function", info: "Converts a string to upper case" },
  { label: "LOWER()", type: "function", info: "Converts a string to lower case" },
  { label: "LENGTH()", type: "function", info: "Returns the length of a string" },
  
  // Date Functions
  { label: "NOW()", type: "function", info: "Returns the current date and time" },
  { label: "CURDATE()", type: "function", info: "Returns the current date" },
  { label: "YEAR()", type: "function", info: "Returns the year part of a date" },
  { label: "MONTH()", type: "function", info: "Returns the month part of a date" },
  { label: "DAY()", type: "function", info: "Returns the day part of a date" },
  { label: "DATE_FORMAT()", type: "function", info: "Formats a date as specified" },
  
  // MySQL Specific
  { label: "IFNULL()", type: "function", info: "Returns a specified value if the expression is NULL" },
  { label: "COALESCE()", type: "function", info: "Returns the first non-NULL value in a list" }
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
    
    // Filter SQL keywords based on input
    const filteredKeywords = sqlKeywords.filter(keyword => 
        keyword.label.toLowerCase().startsWith(word.text.toLowerCase())
    );
    
    if (filteredKeywords.length === 0 && !context.explicit) {
        return null;
    }
    
    return {
        from: word.from,
        options: filteredKeywords,
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
    }
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
</style>