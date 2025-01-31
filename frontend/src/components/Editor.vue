<template>
    <div class="flex lg:flex-row flex-col gap-3 w-full">
        <div class="flex flex-col lg:w-1/2 w-full">
            <div>Entrada</div>
            <code-mirror v-model="value" :lang="lang" :extensions="[oneDarkTheme]" :linter="null" basic wrap tab
                class="w-full" />
        </div>
        <div class="flex flex-col lg:w-1/2 w-fulll">
            <div>Saida</div>
            <code-mirror v-model="firstLineBinded" :lang="lang" :extensions="[oneDarkTheme]" :linter="null" basic wrap tab />
        </div>
    </div>
</template>

<script lang="ts" setup>
import { computed, ref } from 'vue';
import CodeMirror from 'vue-codemirror6';
import { sql, MySQL } from "@codemirror/lang-sql";
import { oneDarkTheme } from '@codemirror/theme-one-dark';

const props = defineProps<{
    variables: { field: string, value: string, position: number }[]
    data: { [key: string]: any }[]
    modelValue: string
}>()

const value = ref(props.modelValue);

const lang = sql({
    dialect: MySQL,
});

const firstLineBinded = computed(() => {
    return makeBindedSQL(value.value) ?? "";
});

const makeBindedSQL = (query: string) => {
    return props.data.reduce((acc, row, index) => {
        query = query.replace(/{{\w+}}/g, (match) => {
            const variable = props.variables.find((variable) => {
                return match === `{{${variable.value}}}`;
            });

            if (!variable) {
                return match;
            }

            return row[variable.field];
        });

        acc += query;

        if (index !== props.data.length - 1) {
            acc += "\n";
        }

        return acc;
    }, "");
}

const getBindedSQL = (): string => {
    return makeBindedSQL(value.value);
}

defineExpose({
    getBindedSQL
});

</script>