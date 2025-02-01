<template>
  <div class="flex flex-col gap-3">
    <h2 class="sr-only">Steps</h2>
    <div
      class="relative after:absolute after:inset-x-0 after:top-1/2 after:block after:h-0.5 after:-translate-y-1/2 after:rounded-lg after:bg-gray-100">
      <ol class="relative z-10 flex justify-between text-sm font-medium text-gray-500">
        <li class="flex items-center gap-2 bg-gray-200 dark:bg-gray-800 dark:text-white text-gray-900 p-2">
          <span :class="{
            'size-6 rounded-full bg-gray-300 dark:bg-gray-700 dark:text-white text-gray-900 text-center text-[10px]/6 font-bold': value !== 0,
            'size-6 rounded-full bg-indigo-600 text-center text-[10px]/6 font-bold text-white': value === 0
          }">
            1 </span>

          <span class="hidden sm:block"> {{ props.steps[0] }} </span>
        </li>

        <li class="flex items-center gap-2 bg-gray-200 dark:bg-gray-800 dark:text-white text-gray-900 p-2">
          <span :class="{
            'size-6 rounded-full bg-gray-300 dark:bg-gray-700 dark:text-white text-gray-900 text-center text-[10px]/6 font-bold': value !== 1,
            'size-6 rounded-full bg-indigo-600 text-center text-[10px]/6 font-bold text-white': value === 1
          }">
            2
          </span>

          <span class="hidden sm:block"> {{ props.steps[1] }} </span>
        </li>

        <li class="flex items-center gap-2 bg-gray-200 dark:bg-gray-800 dark:text-white text-gray-900 p-2">
          <span :class="{
            'size-6 rounded-full bg-gray-300 dark:bg-gray-700 dark:text-white text-gray-900 text-center text-[10px]/6 font-bold': value !== 2,
            'size-6 rounded-full bg-indigo-600 text-center text-[10px]/6 font-bold text-white': value === 2
          }">
            3 </span>

          <span class="hidden sm:block"> {{ props.steps[2] }} </span>
        </li>
      </ol>
    </div>

    <div class="w-full flex flex-col gap-3">
      <div class=" bg-gray-300 dark:bg-gray-700 px-5 py-3 rounded-md"> 
        <slot></slot>
      </div>
      <div class="flex flex-row sm:justify-end justify-between gap-3">
        <Button :disabled="!hasPrevious" type="button" @click="prevStep">
          Previous
        </Button>

        <Button :disabled="!hasNext" type="button" @click="nextStep">
          Next
        </button>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { ref, watch, computed } from 'vue';
import Button from './Button.vue';

const props = defineProps<{
  steps: string[];
  modelValue: number;
  disableNextButton?: boolean;
}>();

const value = ref(props.modelValue);

const emit = defineEmits<{
  'update:modelValue': (value: number) => void;
}>();

const hasPrevious = computed(() => value.value > 0);
const hasNext = computed(() => value.value < props.steps.length - 1 && !props.disableNextButton);

const nextStep = () => {
  if (value.value < props.steps.length - 1) {
    value.value++;
  }
};

const prevStep = () => {
  if (value.value > 0) {
    value.value--;
  }
};

watch(() => props.modelValue, (value) => {
  emit('update:modelValue', value);
});

watch(() => value.value, (value) => {
  emit('update:modelValue', value);
});

</script>

<style scoped></style>