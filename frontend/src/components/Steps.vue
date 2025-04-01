<template>
  <div class="flex flex-col gap-4">
    <h2 class="sr-only">Steps</h2>
    
    <!-- Barra de progresso -->
    <div class="relative mb-2">
      <div class="flex mb-2 items-center justify-between">
        <div class="w-full flex items-center">
          <!-- Linha de progresso completa (background) -->
          <div class="w-full bg-gray-200 dark:bg-gray-700 rounded-full h-2.5">
            <!-- Barra de progresso preenchida baseada no valor atual -->
            <div 
              class="bg-indigo-600 dark:bg-indigo-500 h-2.5 rounded-full transition-all duration-500 ease-in-out" 
              :style="`width: ${(value / (steps.length - 1)) * 100}%`"
            ></div>
          </div>
        </div>
      </div>
      
      <!-- Steps indicators -->
      <ol class="grid grid-cols-3 mt-4 text-sm font-medium text-gray-600 dark:text-gray-400">
        <li 
          v-for="(stepName, index) in steps" 
          :key="index" 
          class="flex flex-col items-center"
          :class="{'text-indigo-600 dark:text-indigo-400': value >= index, 'text-gray-500 dark:text-gray-500': value < index}"
        >
          <!-- Círculo indicador de step -->
          <div 
            class="flex items-center justify-center w-8 h-8 rounded-full mb-2 transition-all duration-300"
            :class="[
              value === index 
                ? 'bg-indigo-600 text-white dark:bg-indigo-500 dark:text-white ring-4 ring-indigo-100 dark:ring-indigo-900/30' 
                : value > index 
                  ? 'bg-indigo-100 dark:bg-indigo-900/30 text-indigo-600 dark:text-indigo-400' 
                  : 'bg-gray-200 dark:bg-gray-700 text-gray-500 dark:text-gray-500'
            ]"
          >
            <span v-if="value > index" class="text-xl">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
              </svg>
            </span>
            <span v-else class="text-xs font-bold">{{ index + 1 }}</span>
          </div>
          
          <!-- Nome do step -->
          <span class="text-xs text-center font-medium transition-colors duration-300 max-w-[120px]">
            {{ stepName }}
          </span>
        </li>
      </ol>
    </div>
    
    <!-- Conteúdo do step atual -->
    <div class="bg-white dark:bg-gray-800 p-6 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700">
      <slot></slot>
    </div>
    
    <!-- Botões de navegação -->
    <div class="flex justify-between items-center mt-4">
      <button 
        @click="prevStep" 
        :disabled="!hasPrevious" 
        class="flex items-center gap-1 py-2 px-4 rounded-md transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
        :class="[
          hasPrevious 
            ? 'text-gray-700 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700' 
            : 'text-gray-400 dark:text-gray-600 cursor-not-allowed'
        ]"
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
        </svg>
        Previous
      </button>
      
      <div class="text-sm text-gray-500 dark:text-gray-400">
        Step {{ value + 1 }} of {{ steps.length }}
      </div>
      
      <button 
        @click="nextStep" 
        :disabled="!hasNext" 
        class="flex items-center gap-1 py-2 px-4 rounded-md transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
        :class="[
          hasNext 
            ? 'bg-indigo-600 text-white hover:bg-indigo-700 dark:bg-indigo-500 dark:hover:bg-indigo-600' 
            : 'bg-gray-300 dark:bg-gray-700 text-gray-500 dark:text-gray-500 cursor-not-allowed'
        ]"
      >
        Next
        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
        </svg>
      </button>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, watch, computed } from 'vue';

const props = defineProps<{
  steps: string[];
  modelValue: number;
  disableNextButton?: boolean;
}>();

const value = ref(props.modelValue);

const emit = defineEmits(['update:modelValue']);

const hasPrevious = computed(() => value.value > 0);
const hasNext = computed(() => value.value < props.steps.length - 1 && !props.disableNextButton);

const nextStep = () => {
  if (hasNext.value) {
    value.value++;
  }
};

const prevStep = () => {
  if (hasPrevious.value) {
    value.value--;
  }
};

watch(() => props.modelValue, (newValue) => {
  value.value = newValue;
});

watch(() => value.value, (newValue) => {
  emit('update:modelValue', newValue);
});
</script>

<style scoped></style>