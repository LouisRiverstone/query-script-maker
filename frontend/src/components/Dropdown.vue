<template>
    <div class="flex flex-col w-full">
        <label :for="id" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">{{ label }}</label>
        <div class="relative inline-block w-full">
            <button @click="toggleDropdown"
                :id="id"
                :disabled="disabled"
                class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-indigo-500 focus:border-indigo-500 block w-full p-2.5 dark:bg-gray-600 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-indigo-500 dark:focus:border-indigo-500 text-start">
                {{ selectedOption ? selectedOption.label : 'Select an option' }}
                <span class="absolute inset-y-0 right-0 top-3 flex pr-2 pointer-events-none">
                    <svg class="h-5 w-5 text-gray-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20"
                        fill="currentColor" aria-hidden="true">
                        <path fill-rule="evenodd"
                            d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z"
                            clip-rule="evenodd" />
                    </svg>
                </span>
            </button>
            <div v-if="isOpen" class="absolute mt-1 w-full rounded-md bg-white dark:bg-gray-600 shadow-2xl z-10 border border-gray-300 dark:border-gray-700">
                <ul tabindex="-1" role="listbox"
                    class="max-h-60 rounded-md py-1 text-base text-black :dark:text-white overflow-auto focus:outline-none sm:text-sm">
                    <li v-for="(option,i) in props.options" :key="i" @click="selectOption(option)"
                        class="cursor-pointer select-none relative py-2 pl-3 pr-9 hover:bg-indigo-600 hover:text-white text-black dark:text-white dark:hover:bg-indigo-600 dark:hover:text-white border-2 border-transparent">
                        <span
                            :class="{ 'font-semibold': selectedOption && selectedOption.value === option.value, 'font-normal': !(selectedOption && selectedOption.value === option.value) }"
                            class="block truncate">
                            {{ option.label }}
                        </span>
                    </li>
                </ul>
            </div>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { ref,watch } from 'vue';

const props = defineProps<{
    options: {
        value: string
        label: string
    }[]
    id: string
    disabled?: boolean
    label: string
    modelValue: string
}>()

const emit = defineEmits(['update:modelValue'])

const isOpen = ref(false);
const selectedOption = ref<{ label: string; value: string } | null>(null);

const toggleDropdown = () => {
    isOpen.value = !isOpen.value;
};

const selectOption = (option: { label: string; value: string }) => {
    selectedOption.value = option;
    isOpen.value = false;
    emit('update:modelValue', selectedOption.value?.value ?? null);
};

watch(() => props.modelValue, (newValue) => {
    const option = props.options.find((option) => option.value === newValue);
    selectedOption.value = option || null;
});
</script>