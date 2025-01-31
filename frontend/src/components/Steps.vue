<template>
    <div class="w-full max-w-2xl mx-auto">
      <h2 class="sr-only">Steps</h2>
      <div class="after:mt-4 after:block after:h-1 after:w-full after:rounded-lg after:bg-gray-200">
        <ol class="grid grid-cols-3 text-sm font-medium text-gray-500">
          <li v-for="(step, i) in steps" :key="i" class="relative flex text-indigo-600" :class="{
            'justify-start': i < currentStep,
            'justify-center': i === currentStep,
            'justify-end': i > currentStep
          }">
            <span class="absolute -bottom-[1.75rem] start-0 rounded-full text-white"
              :class="{ 
                'bg-indigo-600': i <= currentStep,
                'bg-gray-600': i === currentStep,
                'bg-gray-200': i > currentStep,
                'start-0': i === 0,
                '-bottom-[1.75rem] left-1/2 -translate-x-1/2' : i > 0
              }">
              <svg class="size-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd"
                  d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z"
                  clip-rule="evenodd" />
              </svg>
            </span>
  
            <span class="hidden sm:block">{{ step.title }}</span>
  
            <svg class="size-6 sm:hidden" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
              stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round"
                d="M10 6H5a2 2 0 00-2 2v9a2 2 0 002 2h14a2 2 0 002-2V8a2 2 0 00-2-2h-5m-4 0V5a2 2 0 114 0v1m-4 0a2 2 0 104 0m-5 8a2 2 0 100-4 2 2 0 000 4zm0 0c1.306 0 2.417.835 2.83 2M9 14a3.001 3.001 0 00-2.83 2M15 11h3m-3 4h2" />
            </svg>
          </li>
        </ol>
      </div>
    </div>
  </template>
  
  <script lang="ts" setup>
  import { ref } from "vue"
  
  type TStep = {
    title: string
  }
  
  const props = defineProps<{
    steps: TStep[]
  }>()
  
  const currentStep = ref(0)
  
  const nextStep = () => {
    if (currentStep.value < props.steps.length - 1) {
      currentStep.value++
    }
  }
  
  const prevStep = () => {
    if (currentStep.value > 0) {
      currentStep.value--
    }
  }
  
  defineExpose({
    nextStep,
    prevStep
  })
  </script>