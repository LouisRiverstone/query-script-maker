<template>
  <div 
    class="relative" 
    :class="[fullWidth ? 'w-full' : '']"
  >
    <!-- Simple divider without text -->
    <hr v-if="!$slots.default" 
        class="transition-all duration-300"
        :class="[
          styleClasses,
          spaceClasses,
          thicknessClasses,
          { 'w-full': fullWidth },
        ]">
    
    <!-- Divider with text -->
    <div v-else class="relative flex items-center" :class="[spaceClasses]">
      <div class="flex-grow border-t transition-colors duration-300" 
          :class="[styleClasses, thicknessClasses]"></div>
      <span 
        class="flex-shrink-0 font-medium px-4 transition-colors duration-300"
        :class="[textColorClasses, textSizeClasses]"
      >
        <slot></slot>
      </span>
      <div class="flex-grow border-t transition-colors duration-300" 
          :class="[styleClasses, thicknessClasses]"></div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';

// Props
const props = withDefaults(defineProps<{
  style?: 'solid' | 'dashed' | 'dotted';
  color?: 'default' | 'primary' | 'subtle' | 'accent';
  spaceY?: 'none' | 'sm' | 'md' | 'lg' | 'xl';
  thickness?: 'thin' | 'normal' | 'thick';
  textSize?: 'xs' | 'sm' | 'md' | 'lg';
  fullWidth?: boolean;
}>(), {
  style: 'solid',
  color: 'default',
  spaceY: 'md',
  thickness: 'normal',
  textSize: 'sm',
  fullWidth: true
});

// Border style classes
const styleClasses = computed(() => {
  // First determine style (solid, dashed, dotted)
  let borderStyle = 'border-';
  
  if (props.style === 'dashed') borderStyle += 'dashed';
  else if (props.style === 'dotted') borderStyle += 'dotted';
  else borderStyle += 'solid';
  
  // Then color
  let borderColor = '';
  
  switch (props.color) {
    case 'primary':
      borderColor = 'border-indigo-500 dark:border-indigo-400';
      break;
    case 'subtle':
      borderColor = 'border-gray-200 dark:border-gray-700';
      break;
    case 'accent':
      borderColor = 'border-purple-500 dark:border-purple-400';
      break;
    default:
      borderColor = 'border-gray-300 dark:border-gray-600';
      break;
  }
  
  return `${borderStyle} ${borderColor}`;
});

// Vertical spacing
const spaceClasses = computed(() => {
  switch (props.spaceY) {
    case 'none': return 'my-0';
    case 'sm': return 'my-2';
    case 'lg': return 'my-8';
    case 'xl': return 'my-12';
    default: return 'my-4'; // md
  }
});

// Border thickness
const thicknessClasses = computed(() => {
  switch (props.thickness) {
    case 'thin': return 'border-t-[0.5px]';
    case 'thick': return 'border-t-2';
    default: return 'border-t'; // normal
  }
});

// Text color
const textColorClasses = computed(() => {
  switch (props.color) {
    case 'primary':
      return 'text-indigo-600 dark:text-indigo-400';
    case 'subtle':
      return 'text-gray-400 dark:text-gray-500';
    case 'accent':
      return 'text-purple-600 dark:text-purple-400';
    default:
      return 'text-gray-500 dark:text-gray-400';
  }
});

// Text size
const textSizeClasses = computed(() => {
  switch (props.textSize) {
    case 'xs': return 'text-xs';
    case 'md': return 'text-base';
    case 'lg': return 'text-lg';
    default: return 'text-sm'; // sm
  }
});
</script>