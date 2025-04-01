<template>
  <div class="relative" :class="fullWidth ? 'w-full' : ''">
    <label v-if="label" :for="selectId" class="block mb-2 text-sm font-medium text-gray-700 dark:text-gray-300">
      {{ label }}
      <span v-if="required" class="text-red-500 ml-1">*</span>
    </label>
    
    <div class="relative">
      <!-- Left icon -->
      <div v-if="$slots['icon-left']" class="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none text-gray-500 dark:text-gray-400">
        <slot name="icon-left"></slot>
      </div>
      
      <select
        :id="selectId"
        :value="modelValue"
        :required="required"
        :disabled="disabled"
        @change="$emit('update:modelValue', ($event.target as HTMLSelectElement).value)"
        @focus="isFocused = true"
        @blur="isFocused = false"
        class="block appearance-none bg-white dark:bg-gray-800 text-gray-900 dark:text-white border transition-all duration-200 focus:ring-2 focus:outline-none pr-10"
        :class="[
          error ? 'border-red-500 focus:ring-red-200 dark:focus:ring-red-800' : 'border-gray-300 dark:border-gray-600 focus:ring-indigo-200 dark:focus:ring-indigo-800 focus:border-indigo-500 dark:focus:border-indigo-500',
          disabled ? 'bg-gray-100 dark:bg-gray-700 cursor-not-allowed opacity-70' : 'hover:border-gray-400 dark:hover:border-gray-500',
          $slots['icon-left'] ? 'pl-10' : 'pl-4',
          isFocused && !error ? 'border-indigo-500 shadow-sm' : '',
          roundedClasses,
          sizeClasses,
          fullWidth ? 'w-full' : ''
        ]"
      >
        <option v-if="placeholder" value="" disabled selected>{{ placeholder }}</option>
        <slot></slot>
      </select>
      
      <!-- Custom dropdown arrow icon -->
      <div class="absolute inset-y-0 right-0 flex items-center pr-3 pointer-events-none">
        <svg 
          xmlns="http://www.w3.org/2000/svg" 
          class="h-5 w-5 transition-transform duration-200" 
          :class="[
            isFocused ? 'transform rotate-180' : '',
            error ? 'text-red-500' : 'text-gray-500 dark:text-gray-400'
          ]"
          fill="none" 
          viewBox="0 0 24 24" 
          stroke="currentColor"
        >
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
        </svg>
      </div>
    </div>
    
    <!-- Error message with transition -->
    <transition 
      enter-active-class="transition duration-200 ease-out" 
      enter-from-class="opacity-0 -translate-y-1" 
      enter-to-class="opacity-100 translate-y-0"
      leave-active-class="transition duration-150 ease-in" 
      leave-from-class="opacity-100 translate-y-0" 
      leave-to-class="opacity-0 -translate-y-1"
    >
      <p v-if="error" class="mt-1 text-sm text-red-600 dark:text-red-400">
        {{ error }}
      </p>
    </transition>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue';

// Unique ID for select
const selectId = computed(() => `select-${Math.random().toString(36).substr(2, 9)}`);

// Track focus state
const isFocused = ref(false);

// Props
const props = withDefaults(defineProps<{
  modelValue: string;
  label?: string;
  placeholder?: string;
  required?: boolean;
  disabled?: boolean;
  error?: string;
  size?: 'sm' | 'md' | 'lg';
  rounded?: 'none' | 'sm' | 'md' | 'lg' | 'full';
  fullWidth?: boolean;
}>(), {
  placeholder: '',
  required: false,
  disabled: false,
  size: 'md',
  rounded: 'md',
  fullWidth: true
});

// Emits
defineEmits(['update:modelValue']);

// Size classes
const sizeClasses = computed(() => {
  switch (props.size) {
    case 'sm': return 'py-1.5 text-sm';
    case 'lg': return 'py-3 text-base';
    default: return 'py-2 text-sm'; // md
  }
});

// Rounded classes
const roundedClasses = computed(() => {
  switch (props.rounded) {
    case 'none': return 'rounded-none';
    case 'sm': return 'rounded';
    case 'lg': return 'rounded-lg';
    case 'full': return 'rounded-full';
    default: return 'rounded-md'; // md
  }
});
</script>