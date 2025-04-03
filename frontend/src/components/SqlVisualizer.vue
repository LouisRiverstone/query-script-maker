<template>
  <div class="sql-visualizer">
    <div class="flex flex-col md:flex-row gap-4">
      <div class="w-full md:w-1/2">
        <h3 class="text-lg font-semibold text-black dark:text-white mb-2">SQL Query</h3>
        <Editor v-model="sqlQuery" :showBindedSql="false" @update:model-value="debouncedUpdate" />
        <div class="mt-4 flex gap-2">
          <Button type="button" @click="visualize" :disabled="!sqlQuery.trim()">Visualize</Button>
          <Button type="button" @click="resetToOriginal" v-if="props.initialQuery">Reset to Original</Button>
        </div>
      </div>
      
      <div class="w-full md:w-1/2">
        <Suspense>
          <SqlDiagram :query="diagramQuery" :databaseStructure="databaseStructure" />
          <template #fallback>
            <div class="flex items-center justify-center h-64 bg-gray-100 dark:bg-gray-800 rounded-lg">
              <div class="text-center">
                <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-500 mx-auto"></div>
                <p class="mt-2 text-gray-600 dark:text-gray-400">Loading diagram...</p>
              </div>
            </div>
          </template>
        </Suspense>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, computed, defineAsyncComponent, onBeforeUnmount } from 'vue';
import Editor from './Editor.vue';
import Button from './Button.vue';

// Lazy load SqlDiagram component
const SqlDiagram = defineAsyncComponent(() => 
  import('./SqlDiagram.vue').then(mod => {
    // Delay returning the component slightly to ensure smooth loading
    return new Promise(resolve => {
      setTimeout(() => resolve(mod.default), 100);
    });
  })
);

const props = defineProps<{
  initialQuery?: string;
  databaseStructure?: string;
}>();

const sqlQuery = ref(props.initialQuery || '');
const diagramQuery = ref('');
const originalQuery = ref(props.initialQuery || '');
const databaseStructure = computed(() => props.databaseStructure || '');

// Debounce function to prevent excessive updates
const debounce = (fn: Function, delay: number) => {
  let timer: ReturnType<typeof setTimeout> | null = null;
  return function(...args: any[]) {
    if (timer) clearTimeout(timer);
    timer = setTimeout(() => {
      fn(...args);
      timer = null;
    }, delay);
  };
};

// Create debounced update function
const debouncedUpdate = debounce((newValue: string) => {
  if (autoVisualize.value) {
    diagramQuery.value = newValue;
  }
}, 500);

// Auto-visualize flag
const autoVisualize = ref(false);

const visualize = () => {
  diagramQuery.value = sqlQuery.value;
  // Enable auto-visualize after manual visualization
  autoVisualize.value = true;
};

const resetToOriginal = () => {
  sqlQuery.value = originalQuery.value;
  visualize();
};

// Initialize with props.initialQuery if provided
onMounted(() => {
  if (props.initialQuery) {
    originalQuery.value = props.initialQuery;
    sqlQuery.value = props.initialQuery;
    // Immediately visualize the query on mount
    visualize();
  }
});

// Watch for changes to initialQuery with throttling
watch(() => props.initialQuery, (newQuery) => {
  if (newQuery) {
    originalQuery.value = newQuery;
    sqlQuery.value = newQuery;
    // Immediately visualize when initialQuery changes
    visualize();
  }
}, { throttle: 200 });

// Clean up event listeners and timers
onBeforeUnmount(() => {
  autoVisualize.value = false;
});
</script>

<style scoped>
.sql-visualizer {
  width: 100%;
  padding: 1rem;
}
</style> 