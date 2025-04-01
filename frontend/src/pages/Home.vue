<template>
    <div class="flex flex-col flex-grow gap-8 w-full p-6 max-w-6xl mx-auto">
        <Steps :steps="stepsHeaders" v-model="step" :disable-next-button="disableNextButton" class="bg-white dark:bg-gray-800 p-5 rounded-lg shadow-sm">
            <section v-show="step === 0" class="py-4">
                <div class="flex flex-col gap-5">
                    <div class="flex flex-col gap-3">
                        <div class="flex flex-row justify-center">
                            <small class="text-sm text-gray-600 dark:text-gray-400 italic font-light">*The SQL Script maker processes only the first sheet of .xlsx</small>
                        </div>
                        <div class="flex flex-row justify-center my-3">
                            <Button type="button" @click="importXLSX" class="w-48 bg-indigo-600 hover:bg-indigo-700">Import .XLSX</Button>
                        </div>
                    </div>
                    <div v-if="content.length > 0" class="flex flex-row justify-center mt-4 overflow-x-auto">
                        <Table :data="content" class="border border-gray-200 dark:border-gray-700 rounded-md" />
                    </div>
                </div>
            </section>
            <section v-show="step === 1" class="py-4">
                <div v-if="headers.length > 0" class="flex flex-row justify-center">
                    <VariablesCaster :headers="headers" ref="variblesCasterRef" class="w-full max-w-4xl" />
                </div>
            </section>
            <section v-show="step === 2" class="py-4 ">
                <div class="flex flex-col gap-5">
                    <div v-if="!loading" class="flex flex-col gap-5">
                        <div class="flex flex-col md:flex-row gap-3 items-start">
                            <Dropdown label="Saved Queries" id="select-query" :options="querySelectOptions"
                                v-model="selectedQueryId" class="flex-1" />
                        </div>
                        
                        <div v-if="selectedQueryDescription" class="bg-gray-100 dark:bg-gray-700 p-3 rounded-md">
                            <span class="font-medium text-gray-700 dark:text-gray-300">Query Description:</span>
                            <p class="text-gray-600 dark:text-gray-400">{{ selectedQueryDescription }}</p>
                        </div>
                        
                        <div v-if="hasVariablesAssigned" class="flex flex-row justify-center mt-2">
                            <Editor v-model="query" :show-binded-sql="true" :variables="variables" :data="content"
                                :minify="minify" ref="editorRef" class="w-full border border-gray-200 dark:border-gray-700 rounded-md">
                                <div class="flex md:flex-row flex-col justify-between gap-3 w-full mt-4">
                                    <div class="flex flex-row items-center space-x-2">
                                        <input type="checkbox" id="minify" v-model="minify" class="w-4 h-4 text-indigo-600 bg-gray-100 border-gray-300 rounded 
                                            focus:ring-indigo-500 dark:bg-gray-700 dark:border-gray-600 
                                            dark:ring-offset-gray-800 focus:ring-2 dark:focus:ring-indigo-600 
                                            checked:bg-indigo-500 hover:border-indigo-400 transition-colors" />
                                        <label for="minify"
                                            class="text-sm font-medium text-gray-700 dark:text-gray-300">
                                            Minify Output
                                        </label>
                                    </div>
                                    
                                    <div class="flex flex-wrap md:flex-row flex-col gap-2 justify-end">
                                        <Button type="button" @click="openSqlVisualizer" class="bg-purple-600 hover:bg-purple-700">
                                            Visualize SQL
                                        </Button>
                                        <div v-if="!loadingDatabaseConnection" class="flex flex-wrap md:flex-row flex-col gap-2">
                                            <Button v-if="databaseConnection && databaseConnection.ID !== 0" type="button"
                                                @click="testInputSql" class="bg-blue-600 hover:bg-blue-700">Test Input SQL</Button>
                                            <Button v-if="databaseConnection && databaseConnection.ID !== 0" type="button"
                                                @click="testOutputSql" class="bg-blue-700 hover:bg-blue-800">Test Output SQL</Button>
                                        </div>
                                        <div v-if="loadingDatabaseConnection" class="flex justify-center">
                                            <Loader />
                                        </div>
                                        <div v-if="hasEditor">
                                            <Button type="button" @click="createSqlFile" class="bg-green-600 hover:bg-green-700">Save .SQL</Button>
                                        </div>
                                    </div>
                                </div>
                            </Editor>
                        </div>
                        <div v-else class="flex justify-center py-8">
                            <Loader />
                        </div>
                    </div>
                </div>
            </section>
        </Steps>
        <SqlResultTable :data="responseTest" v-model="showSqlTable" />
        <SqlVisualizerModal :isOpen="showSqlVisualizer" :initialQuery="query" @close="closeSqlVisualizer" />
    </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue';
import { asyncComputed } from '@vueuse/core';

import { ReadXLSXFile, CreateSQLFile, GetQueriesList, GetDatabaseConnection, TestQueryInDatabase, TestBatchQueryInDatabase, MakeBindedSQL } from '../../wailsjs/go/main/App'
import { main } from '../../wailsjs/go/models';

import Table from '../components/Table.vue';
import VariablesCaster from '../components/VariablesCaster.vue';
import Editor from '../components/Editor.vue';
import Steps from '../components/Steps.vue';
import Button from '../components/Button.vue';
import Dropdown from '../components/Dropdown.vue';
import Loader from '../components/Loader.vue';
import SqlResultTable from '../components/SqlResultTable.vue';
import SqlVisualizerModal from '../components/SqlVisualizerModal.vue';

const loading = ref<boolean>(false);
const loadingDatabaseConnection = ref<boolean>(false);
const content = ref<{ [k: string]: string }[]>([]);
const variblesCasterRef = ref<typeof VariablesCaster | null>(null);
const editorRef = ref<typeof Editor | null>(null);
const query = ref<string>('SELECT * from family limit 1;');
const showSqlTable = ref<boolean>(false);
const showSqlVisualizer = ref<boolean>(false);
const minify = ref<boolean>(false);

const selectedQueryId = ref<string>("TESTE")
const queries = ref<Array<main.Query>>([])


const databaseConnection = ref<main.DatabaseConnection | undefined>({
    ID: 0,
    Host: "",
    Port: 0,
    Username: "",
    Password: "",
    Database: ""
});

const responseTest = ref<{ [k: string]: any }[][]>([]);

const stepsHeaders = ref<string[]>(['Import .XLSX', 'Assign Variables', 'Save .SQL']);

const step = ref(0);

const firstInputCasted = asyncComputed(async () => {
    const firstContent = content.value[0]

    return await MakeBindedSQL(query.value, [firstContent], variables.value, false);
}, )


const headers = computed(() => {
    if (content.value.length === 0) {
        return [];
    }

    return Object.keys(content.value[0]);
})

const disableNextButton = computed(() => {
    if (step.value === 0) {
        return content.value.length === 0 && headers.value.length === 0;
    }

    if (step.value === 1) {
        return !hasVariablesAssigned.value;
    }
})

const hasVariablesAssigned = computed(() => {
    if (!variblesCasterRef.value) {
        return false;
    }

    return variblesCasterRef.value.variables.length > 0;
})

const hasEditor = computed(() => {
    return !!editorRef.value;
})

const variables = computed(() => {
    if (!variblesCasterRef.value) {
        return [];
    }

    return variblesCasterRef.value.variables;
})

const selectedQueryDescription = computed(() => {
    const querySelected = queries.value.find((query) => query.ID === Number(selectedQueryId.value));

    if (!querySelected) {
        return;
    }

    return querySelected.Description;
})

const querySelectOptions = computed(() => {
    return queries.value.map((query) => {
        return {
            label: query.Title,
            value: String(query.ID),
        }
    })
})

const importXLSX = async () => {
    content.value = JSON.parse(await ReadXLSXFile());
}

const createSqlFile = async () => {
    if (!hasEditor.value) {
        return;
    }

    const sql = await editorRef.value!.getBindedSQL();

    await CreateSQLFile(sql);
}

const getDatabaseConnection = async () => {
    try {
        loadingDatabaseConnection.value = true
        return await GetDatabaseConnection()
    } catch (error) {
        console.error(error)
    } finally {
        loadingDatabaseConnection.value = false
    }
}

const testInputSql = async () => {
    try {
        const firstContent = content.value[0]
        const bindedSql = await MakeBindedSQL(query.value, [firstContent], variables.value, false);

        await testSQL(bindedSql)
    } catch (error) {
        alert(error)
    }
}

const testOutputSql = async () => {
    try {
        if (!hasEditor.value) {
            return;
        }

        await testBatchSQL(await editorRef.value!.getBindedSQL())
    } catch (error) {
        alert(error)
    }
}

const testSQL = async (query: string) => {
    try {
        if (!databaseConnection.value) {
            alert("Database connection not found")
        }

        responseTest.value = []
        responseTest.value = [(await TestQueryInDatabase(databaseConnection.value!, query))]

        showSqlTable.value = true
    } catch (error) {
        throw error
    }
}

const testBatchSQL = async (query: string) => {
    try {
        if (!databaseConnection.value) {
            alert("Database connection not found")
        }

        const queries = query.replaceAll("\n", "").split(';').filter((query) => query.trim() !== '').map((query) => `${query};`);

        responseTest.value = []
        responseTest.value = (await TestBatchQueryInDatabase(databaseConnection.value!, queries))

        showSqlTable.value = true
    } catch (error) {
        throw error
    }
}

const openSqlVisualizer = () => {
  showSqlVisualizer.value = true;
};

const closeSqlVisualizer = () => {
  showSqlVisualizer.value = false;
};

const mount = async () => {
    try {
        loading.value = true
        queries.value = await GetQueriesList(false)
        databaseConnection.value = await getDatabaseConnection()
    } catch (error) {
        console.error(error)
    } finally {
        loading.value = false
    }
}

onMounted(() => {
    mount();
})

watch(() => selectedQueryId.value, (newValue) => {
    const querySelected = queries.value.find((query) => query.ID === Number(newValue));

    if (!querySelected) {
        return;
    }

    query.value = querySelected.Query;
})
</script>