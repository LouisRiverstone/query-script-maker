<template>
    <button
        class="group relative inline-block text-sm font-medium text-indigo-600 focus:outline-none focus:ring active:text-indigo-500"
        :class="{ 'cursor-not-allowed opacity-50': props.disabled }"
        
        :disabled="props.disabled" :type="props.type" @click.prevent="emitClickEvent">
        <span
            class="absolute inset-0 translate-x-0 translate-y-0 bg-indigo-600 transition-transform group-hover:translate-x-0.5 group-hover:translate-y-0.5"></span>
        <span class="relative block border border-current bg-white px-8 py-3">
            <slot></slot>
        </span>
    </button>
</template>

<script lang="ts" setup>
import { useDebounceFn } from '@vueuse/core'

const props = defineProps<{
    type: 'button' | 'submit' | 'reset'
    disabled?: boolean
}>()

const emit = defineEmits(['on-click'])

const emitClickEvent = useDebounceFn(() => {
   emit('on-click')
}, 150)
</script>