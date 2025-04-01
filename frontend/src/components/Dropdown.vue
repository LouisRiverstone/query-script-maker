<template>
    <div class="relative" :class="fullWidth ? 'w-full' : ''">
        <label v-if="label" :for="id" class="block mb-2 text-sm font-medium text-gray-700 dark:text-gray-300">
            {{ label }}
            <span v-if="required" class="text-red-500">*</span>
        </label>
        
        <div class="relative">
            <button 
                @click="toggleDropdown"
                :id="id"
                :disabled="disabled"
                :aria-expanded="isOpen ? 'true' : 'false'"
                aria-haspopup="listbox"
                class="relative w-full transition-all duration-200 flex items-center justify-between text-left border bg-white dark:bg-gray-800 text-gray-900 dark:text-white focus:outline-none focus:ring-2"
                :class="[
                    disabled ? 'cursor-not-allowed bg-gray-100 dark:bg-gray-700 opacity-70' : 'cursor-pointer hover:bg-gray-50 dark:hover:bg-gray-700',
                    error ? 'border-red-500 focus:ring-red-200 dark:focus:ring-red-800' : 'border-gray-300 dark:border-gray-600 focus:ring-indigo-200 dark:focus:ring-indigo-800 focus:border-indigo-500 dark:focus:border-indigo-500',
                    roundedClasses,
                    sizeClasses,
                    fullWidth ? 'w-full' : ''
                ]"
            >
                <span class="block truncate pl-3" :class="{'text-gray-500 dark:text-gray-400': !selectedOption}">
                    {{ selectedOption ? selectedOption.label : placeholder }}
                </span>
                
                <span class="pointer-events-none absolute inset-y-0 right-0 flex items-center pr-2">
                    <svg 
                        class="h-5 w-5 transition-transform duration-200" 
                        :class="[
                            isOpen ? 'transform rotate-180' : '',
                            error ? 'text-red-500' : 'text-gray-400 dark:text-gray-500'
                        ]" 
                        xmlns="http://www.w3.org/2000/svg" 
                        viewBox="0 0 20 20" 
                        fill="currentColor" 
                        aria-hidden="true"
                    >
                        <path fill-rule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clip-rule="evenodd" />
                    </svg>
                </span>
            </button>

            <!-- Dropdown options -->
            <div 
                v-if="isOpen" 
                class="absolute z-10 mt-1 w-full overflow-auto rounded-md bg-white dark:bg-gray-800 shadow-lg border border-gray-200 dark:border-gray-700 ring-1 ring-black ring-opacity-5 transition-all duration-100 max-h-60"
                :class="[roundedClasses]"
            >
                <ul 
                    tabindex="-1" 
                    role="listbox" 
                    :aria-labelledby="id" 
                    class="py-1 text-base overflow-auto focus:outline-none"
                >
                    <li 
                        v-for="(option, i) in options" 
                        :key="i" 
                        @click="selectOption(option)"
                        role="option"
                        :aria-selected="selectedOption && selectedOption.value === option.value ? 'true' : 'false'"
                        class="cursor-pointer transition-colors duration-150 relative py-2 pl-3 pr-9 hover:bg-indigo-50 dark:hover:bg-indigo-900/30"
                        :class="[
                            selectedOption && selectedOption.value === option.value 
                                ? 'bg-indigo-100 dark:bg-indigo-900/20 text-indigo-900 dark:text-indigo-200' 
                                : 'text-gray-900 dark:text-gray-200'
                        ]"
                    >
                        <span 
                            class="block truncate"
                            :class="selectedOption && selectedOption.value === option.value ? 'font-medium' : 'font-normal'"
                        >
                            {{ option.label }}
                        </span>

                        <!-- Selected indicator -->
                        <span 
                            v-if="selectedOption && selectedOption.value === option.value" 
                            class="absolute inset-y-0 right-0 flex items-center pr-4 text-indigo-600 dark:text-indigo-400"
                        >
                            <svg class="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                                <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                            </svg>
                        </span>
                    </li>
                </ul>
            </div>
        </div>
        
        <!-- Mensagem de erro -->
        <p v-if="error" class="mt-1 text-sm text-red-600 dark:text-red-400">
            {{ error }}
        </p>
    </div>
</template>

<script setup lang="ts">
import { ref, watch, computed, onMounted, onBeforeUnmount } from 'vue';

// Props
const props = withDefaults(defineProps<{
    options: {
        value: string
        label: string
    }[]
    id: string
    label?: string
    placeholder?: string
    required?: boolean
    disabled?: boolean
    error?: string
    modelValue: string
    size?: 'sm' | 'md' | 'lg'
    rounded?: 'none' | 'sm' | 'md' | 'lg' | 'full'
    fullWidth?: boolean
}>(), {
    placeholder: 'Select an option',
    disabled: false,
    required: false,
    size: 'md',
    rounded: 'md',
    fullWidth: true
});

// Emits
const emit = defineEmits(['update:modelValue']);

const isOpen = ref(false);
const selectedOption = ref<{ label: string; value: string } | null>(null);

const toggleDropdown = () => {
    if (!props.disabled) {
        isOpen.value = !isOpen.value;
    }
};

const selectOption = (option: { label: string; value: string }) => {
    selectedOption.value = option;
    isOpen.value = false;
    emit('update:modelValue', option.value);
};

const handleClickOutside = (event: MouseEvent) => {
    const target = event.target as HTMLElement;
    if (!target.closest(`#${props.id}`)) {
        isOpen.value = false;
    }
};

onMounted(() => {
    document.addEventListener('click', handleClickOutside);
});

onBeforeUnmount(() => {
    document.removeEventListener('click', handleClickOutside);
});

watch(() => props.modelValue, (newValue) => {
    const option = props.options.find((option) => option.value === newValue);
    selectedOption.value = option || null;
}, { immediate: true });

const sizeClasses = computed(() => {
    switch (props.size) {
        case 'sm': return 'py-1.5 text-sm';
        case 'lg': return 'py-3 text-base';
        default: return 'py-2 text-sm';
    }
});

const roundedClasses = computed(() => {
    switch (props.rounded) {
        case 'none': return 'rounded-none';
        case 'sm': return 'rounded';
        case 'lg': return 'rounded-lg';
        case 'full': return 'rounded-full';
        default: return 'rounded-md';
    }
});
</script>