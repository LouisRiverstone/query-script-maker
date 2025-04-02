<template>
    <div class="flex flex-col gap-3 w-full">
        <slot></slot>
        <div class="flex flex-col lg:flex-row gap-6">
            <div class="flex flex-col w-full gap-3">
                <Divider v-if="showBindedSql">Input</Divider>
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
    const isDark = 
        document.documentElement.classList.contains('dark') || 
        document.body.classList.contains('dark') ||
        window.matchMedia('(prefers-color-scheme: dark)').matches;
    
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

// Create a completion source for variables
const variableCompletions = computed(() => {
    if (!props.variables) return [];
    
    return props.variables.map(variable => ({
        label: `{{ ${variable.Value} }}`,
        type: "variable",
        detail: variable.Field,
        info: `Campo: ${variable.Field}`,
        apply: `{{ ${variable.Value} }}`,
        boost: 99, // Give variables higher priority
    }));
});

// Autocomplete function
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
const extensions = computed(() => [
    isDarkMode.value ? oneDarkTheme : [],
    myTheme,
    createVariableHighlighter(),
    autocompletion({
        override: [completeVariables],
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
                    
                    // Only show for variable type completions
                    if (completion.type === "variable" && completion.detail) {
                        dom.textContent = completion.detail.length > 20 
                            ? completion.detail.substring(0, 20) + "..." 
                            : completion.detail;
                    }
                    
                    return dom;
                },
                position: 80
            }
        ]
    })
]);

// Output editor extensions (simpler, no autocomplete)
const outputExtensions = computed(() => [
    isDarkMode.value ? oneDarkTheme : [],
    myTheme
]);

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

// Watch for dark mode changes
watch(() => isDarkMode.value, (newValue) => {
    // The computed extensions will automatically update
    console.log("Dark mode changed to:", newValue);
});

defineExpose({
    getBindedSQL
});
</script>