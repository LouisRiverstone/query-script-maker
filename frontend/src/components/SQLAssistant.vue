<template>
  <div class="sql-assistant flex flex-col gap-3 border border-gray-200 dark:border-gray-700 rounded-md p-4 bg-white dark:bg-gray-800">
    <div class="flex items-center justify-between mb-2">
      <h3 class="text-sm font-medium text-gray-700 dark:text-gray-300">SQL AI Assistant</h3>
      <button
        @click="toggleAssistant"
        class="text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-200 p-1 rounded-full"
      >
        <svg
          xmlns="http://www.w3.org/2000/svg"
          class="h-4 w-4"
          fill="none"
          viewBox="0 0 24 24"
          stroke="currentColor"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M6 18L18 6M6 6l12 12"
          />
        </svg>
      </button>
    </div>

    <div class="flex flex-col gap-3">
      <div class="flex flex-col">
        <label class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
          Describe the SQL query in natural language
        </label>
        <textarea
          v-model="prompt"
          class="w-full rounded-md border border-gray-200 dark:border-gray-700 bg-white dark:bg-gray-800 text-gray-900 dark:text-gray-100 text-sm p-2 focus:ring-blue-500 focus:border-blue-500 min-h-[80px] resize-vertical"
          placeholder="Example: Show all customers who made a purchase in the last month ordered by total amount spent"
          @keydown.ctrl.enter="generateSQL"
        ></textarea>
        <div v-if="errorMessage" class="mt-1 text-xs text-red-500 dark:text-red-400">
          {{ errorMessage }}
        </div>
      </div>
      <p class="text-xs text-gray-500 dark:text-gray-400 mt-1">
        Press Ctrl+Enter to generate SQL or click the button below
      </p>

      <div class="flex justify-end">
        <button
          @click="generateSQL"
          class="inline-flex items-center justify-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed"
          :disabled="isLoading || !prompt.trim()"
        >
          <svg v-if="isLoading" class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          {{ isLoading ? 'Generating...' : 'Generate SQL' }}
        </button>
      </div>

      <div v-if="generatedSQL" class="mt-2 border border-gray-200 dark:border-gray-700 rounded-md overflow-hidden">
        <div class="flex justify-between items-center p-2 bg-gray-50 dark:bg-gray-900 border-b border-gray-200 dark:border-gray-700">
          <span class="text-xs font-medium text-gray-600 dark:text-gray-400">Generated SQL</span>
          <button 
            @click="copyToClipboard" 
            class="flex items-center text-xs text-blue-600 dark:text-blue-400 hover:text-blue-800 dark:hover:text-blue-300 px-2 py-1 rounded"
            :class="copySuccess ? 'text-green-600 dark:text-green-400' : ''"
          >
            <svg v-if="copySuccess" xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
            </svg>
            <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" viewBox="0 0 20 20" fill="currentColor">
              <path d="M8 3a1 1 0 011-1h2a1 1 0 110 2H9a1 1 0 01-1-1z" />
              <path d="M6 3a2 2 0 00-2 2v11a2 2 0 002 2h8a2 2 0 002-2V5a2 2 0 00-2-2 3 3 0 01-3 3H9a3 3 0 01-3-3z" />
            </svg>
            {{ copySuccess ? 'Copied!' : 'Copy' }}
          </button>
        </div>
        <div class="p-3 bg-gray-50 dark:bg-gray-900 text-xs font-mono whitespace-pre-wrap text-gray-800 dark:text-gray-200">{{ generatedSQL }}</div>
      </div>

      <div class="flex justify-end mt-1">
        <button
          @click="useGeneratedSQL"
          class="text-xs text-gray-600 dark:text-gray-400 hover:text-indigo-600 dark:hover:text-indigo-400 px-2 py-1 rounded"
          :disabled="!generatedSQL"
        >
          Use this SQL
        </button>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref } from 'vue';
import { GenerateSQLFromPrompt } from '../../wailsjs/go/main/App';

const props = defineProps<{
  isVisible: boolean
}>();

const emit = defineEmits(['update:isVisible', 'sqlGenerated']);

const prompt = ref('');
const generatedSQL = ref('');
const isLoading = ref(false);
const errorMessage = ref('');
const copySuccess = ref(false);

// Generate SQL from the prompt
const generateSQL = async () => {
  if (!prompt.value.trim()) {
    errorMessage.value = 'Please enter a description of the SQL query you need.';
    return;
  }

  try {
    isLoading.value = true;
    errorMessage.value = '';
    
    const sql = await GenerateSQLFromPrompt(prompt.value);
    generatedSQL.value = sql;
  } catch (error) {
    console.error('Error generating SQL:', error);
    errorMessage.value = error instanceof Error ? error.message : 'Failed to generate SQL. Please try again.';
    generatedSQL.value = '';
  } finally {
    isLoading.value = false;
  }
};

// Copy the generated SQL to clipboard
const copyToClipboard = async () => {
  try {
    await navigator.clipboard.writeText(generatedSQL.value);
    copySuccess.value = true;
    setTimeout(() => {
      copySuccess.value = false;
    }, 2000);
  } catch (err) {
    console.error('Failed to copy:', err);
  }
};

// Use the generated SQL in the editor
const useGeneratedSQL = () => {
  if (generatedSQL.value) {
    emit('sqlGenerated', generatedSQL.value);
  }
};

// Toggle the assistant visibility
const toggleAssistant = () => {
  emit('update:isVisible', false);
};
</script>

<style scoped>
.sql-assistant {
  max-width: 100%;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
}
</style> 