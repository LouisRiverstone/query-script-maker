<template>
  <div v-if="isOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black bg-opacity-50">
    <div class="bg-white dark:bg-gray-800 rounded-lg shadow-xl w-full max-w-5xl max-h-[90vh] flex flex-col">
      <!-- Modal Header -->
      <div class="flex justify-between items-center p-4 border-b border-gray-200 dark:border-gray-700">
        <h2 class="text-xl font-semibold text-gray-800 dark:text-white">SQL Visualizer</h2>
        <button 
          @click="close"
          class="text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200 focus:outline-none"
        >
          <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
      
      <!-- Modal Body -->
      <div class="flex-1 overflow-auto p-4">
        <SqlVisualizer :initialQuery="initialQuery" :databaseStructure="databaseStructure" />
      </div>
      
      <!-- Modal Footer -->
      <div class="flex justify-end p-4 border-t border-gray-200 dark:border-gray-700">
        <Button type="button" @click="close">Close</Button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onBeforeUnmount } from 'vue';
import SqlVisualizer from './SqlVisualizer.vue';
import Button from './Button.vue';
import { GetLatestDatabaseStructure } from '../../wailsjs/go/main/App';

const props = defineProps<{
  isOpen: boolean;
  initialQuery?: string;
  databaseStructure?: string;
}>();

const emit = defineEmits(['close']);
const databaseStructure = ref<string>('');

const close = () => {
  emit('close');
};

// Fetch database structure if not provided
const fetchDatabaseStructure = async () => {
  if (!props.databaseStructure && props.isOpen) {
    try {
      const structure = await GetLatestDatabaseStructure();
      if (structure) {
        databaseStructure.value = structure;
      }
    } catch (error) {
      console.error('Error fetching database structure:', error);
    }
  } else if (props.databaseStructure) {
    databaseStructure.value = props.databaseStructure;
  }
};

// Prevent body scrolling when modal is open
watch(() => props.isOpen, (isOpen) => {
  if (isOpen) {
    document.body.style.overflow = 'hidden';
    fetchDatabaseStructure();
  } else {
    document.body.style.overflow = '';
  }
});

// Watch for database structure changes from props
watch(() => props.databaseStructure, (newValue) => {
  if (newValue) {
    databaseStructure.value = newValue;
  }
});

// Handle escape key to close modal
const handleKeyDown = (event: KeyboardEvent) => {
  if (event.key === 'Escape' && props.isOpen) {
    close();
  }
};

onMounted(() => {
  window.addEventListener('keydown', handleKeyDown);
  fetchDatabaseStructure();
});

onBeforeUnmount(() => {
  window.removeEventListener('keydown', handleKeyDown);
});
</script>

<style scoped>
.fixed {
  animation: fadeIn 0.2s ease-in-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}
</style> 