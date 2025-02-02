<template>
    <div class="flex lg:flex-row flex-col gap-3 w-full">
        <div class="flex flex-col" :class="{ 'lg:w-1/2 w-full': showBindedSql, 'lg:w-full w-full': !showBindedSql }">
            <div v-if="showBindedSql">Entrada</div>
            <code-mirror v-model="value" :lang="lang" :extensions="[oneDarkTheme]" :linter="null" basic wrap tab class="w-full" />
        </div>
        <div v-if="showBindedSql" class="flex flex-col lg:w-1/2 w-fulll">
            <div>Saida</div>
            <code-mirror v-model="linesBinded" :lang="lang" :extensions="[oneDarkTheme]" :linter="null" basic wrap tab />
        </div>
    </div>
</template>

<script lang="ts" setup>
import { ref, watch } from 'vue';
import CodeMirror from 'vue-codemirror6';
import { sql, MySQL } from "@codemirror/lang-sql";
import { oneDarkTheme } from '@codemirror/theme-one-dark';
import { MakeBindedSQL } from '../../wailsjs/go/main/App';
import { main } from '../../wailsjs/go/models';
import { computedAsync } from '@vueuse/core';

const props = defineProps<{
    variables?: Array<main.Variable>
    data?: { [key: string]: any }[]
    modelValue: string,
    showBindedSql?: boolean
}>()

const emit= defineEmits(['update:modelValue']);

const value = ref(props.modelValue);

const lang = sql({
    dialect: MySQL,
});

const linesBinded = computedAsync(async () => {
    if(!props.showBindedSql) {
        return "";
    }

    return await MakeBindedSQL(value.value, props.data!, props.variables!) ?? "";
}, "");

const getBindedSQL = async (): Promise<string> => {
    if(!props.showBindedSql) {
        return "";
    }

    return await MakeBindedSQL(value.value, props.data!, props.variables!);
}

defineExpose({
    getBindedSQL
});

watch(() => props.modelValue, (newValue) => {
    value.value = newValue;
});

watch(() => value.value, (val) => {
    emit('update:modelValue', val);
});
</script>