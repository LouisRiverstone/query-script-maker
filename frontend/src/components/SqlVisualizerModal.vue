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
        <SqlVisualizer :initialQuery="initialQuery" />
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

const props = defineProps<{
  isOpen: boolean;
  initialQuery?: string;
}>();

const emit = defineEmits(['close']);

const close = () => {
  emit('close');
};

// Prevent body scrolling when modal is open
watch(() => props.isOpen, (isOpen) => {
  if (isOpen) {
    document.body.style.overflow = 'hidden';
  } else {
    document.body.style.overflow = '';
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