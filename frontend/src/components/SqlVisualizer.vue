<template>
  <div class="sql-visualizer">
    <div class="flex flex-col md:flex-row gap-4">
      <div class="w-full md:w-1/2">
        <h3 class="text-lg font-semibold text-black dark:text-white mb-2">SQL Query</h3>
        <Editor v-model="sqlQuery" :showBindedSql="false" />
        <div class="mt-4 flex gap-2">
          <Button type="button" @click="visualize" :disabled="!sqlQuery.trim()">Visualize</Button>
          <Button type="button" @click="resetToOriginal" v-if="props.initialQuery">Reset to Original</Button>
        </div>
      </div>
      
      <div class="w-full md:w-1/2">
        <SqlDiagram :query="diagramQuery" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue';
import Editor from './Editor.vue';
import SqlDiagram from './SqlDiagram.vue';
import Button from './Button.vue';

const props = defineProps<{
  initialQuery?: string;
}>();

const sqlQuery = ref(props.initialQuery || '');
const diagramQuery = ref('');
const originalQuery = ref(props.initialQuery || '');

const visualize = () => {
  diagramQuery.value = sqlQuery.value;
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

// Watch for changes to initialQuery
watch(() => props.initialQuery, (newQuery) => {
  if (newQuery) {
    originalQuery.value = newQuery;
    sqlQuery.value = newQuery;
    // Immediately visualize when initialQuery changes
    visualize();
  }
});
</script>

<style scoped>
.sql-visualizer {
  width: 100%;
  padding: 1rem;
}
</style> 