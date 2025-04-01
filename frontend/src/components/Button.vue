<template>
    <button
        class="relative inline-flex items-center justify-center font-medium focus:outline-none focus:ring-2 focus:ring-offset-2 transition-all duration-200"
        :class="[
            sizeClasses,
            variantClasses,
            roundedClasses,
            props.disabled ? 'opacity-60 cursor-not-allowed pointer-events-none' : 'transform hover:-translate-y-0.5 hover:shadow-md active:translate-y-0',
            props.fullWidth ? 'w-full' : '',
            props.class
        ]"
        :disabled="props.disabled || props.loading"
        :type="props.type"
        @click="emitClickEvent"
    >
        <!-- Loading spinner overlay -->
        <span v-if="props.loading" class="absolute inset-0 flex items-center justify-center">
            <svg class="animate-spin h-5 w-5" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
        </span>
        
        <!-- Button content -->
        <span :class="{'opacity-0': props.loading}" class="flex items-center gap-2">
            <slot name="icon-left"></slot>
            <slot></slot>
            <slot name="icon-right"></slot>
        </span>
    </button>
</template>

<script lang="ts" setup>
import { useDebounceFn } from '@vueuse/core'
import { computed } from 'vue'

const props = withDefaults(defineProps<{
    type?: 'button' | 'submit' | 'reset'
    variant?: 'primary' | 'secondary' | 'outline' | 'ghost' | 'danger' | 'success'
    size?: 'xs' | 'sm' | 'md' | 'lg'
    rounded?: 'none' | 'sm' | 'md' | 'lg' | 'full'
    disabled?: boolean
    loading?: boolean
    fullWidth?: boolean
    class?: string
}>(), {
    type: 'button',
    variant: 'primary',
    size: 'md',
    rounded: 'md',
    disabled: false,
    loading: false,
    fullWidth: false,
    class: ''
})

const emit = defineEmits(['on-click'])

const emitClickEvent = useDebounceFn(() => {
    if (!props.loading && !props.disabled) {
        emit('on-click')
    }
}, 150)

// Compute classes based on variant
const variantClasses = computed(() => {
    switch (props.variant) {
        case 'primary':
            return 'bg-indigo-600 text-white border border-transparent hover:bg-indigo-700 focus:ring-indigo-500 dark:bg-indigo-500 dark:hover:bg-indigo-600 dark:focus:ring-indigo-400'
        case 'secondary':
            return 'bg-gray-200 text-gray-800 border border-transparent hover:bg-gray-300 focus:ring-gray-500 dark:bg-gray-700 dark:text-gray-200 dark:hover:bg-gray-600 dark:focus:ring-gray-400'
        case 'outline':
            return 'bg-transparent border border-indigo-600 text-indigo-600 hover:bg-indigo-50 focus:ring-indigo-500 dark:border-indigo-400 dark:text-indigo-400 dark:hover:bg-indigo-900/20'
        case 'ghost':
            return 'bg-transparent border border-transparent text-indigo-600 hover:bg-indigo-50 focus:ring-indigo-500 dark:text-indigo-400 dark:hover:bg-indigo-900/20'
        case 'danger':
            return 'bg-red-600 text-white border border-transparent hover:bg-red-700 focus:ring-red-500 dark:bg-red-500 dark:hover:bg-red-600 dark:focus:ring-red-400'
        case 'success':
            return 'bg-green-600 text-white border border-transparent hover:bg-green-700 focus:ring-green-500 dark:bg-green-500 dark:hover:bg-green-600 dark:focus:ring-green-400'
        default:
            return 'bg-indigo-600 text-white border border-transparent hover:bg-indigo-700 focus:ring-indigo-500 dark:bg-indigo-500 dark:hover:bg-indigo-600 dark:focus:ring-indigo-400'
    }
})

// Compute classes based on size
const sizeClasses = computed(() => {
    switch (props.size) {
        case 'xs':
            return 'text-xs px-2.5 py-1.5'
        case 'sm':
            return 'text-sm px-3 py-2'
        case 'lg':
            return 'text-base px-6 py-3'
        default: // md
            return 'text-sm px-4 py-2'
    }
})

// Compute classes based on rounded
const roundedClasses = computed(() => {
    switch (props.rounded) {
        case 'none':
            return 'rounded-none'
        case 'sm':
            return 'rounded'
        case 'lg':
            return 'rounded-lg'
        case 'full':
            return 'rounded-full'
        default: // md
            return 'rounded-md'
    }
})
</script>