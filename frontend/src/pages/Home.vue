<template>
    <div class="flex flex-col flex-grow gap-5 w-full p-5">
        <Steps :steps="stepsHeaders" v-model="step" :disable-next-button="disableNextButton">
            <section v-show="step === 0">
                <div class="flex flex-col gap-3">
                    <div class="flex flex-row justify-center ">
                        <Button type="button" @click="importXLSX">Import .XLSX</Button>
                    </div>
                    <div v-if="content.length > 0" class="flex flex-row justify-center">
                        <Table :data="content" />
                    </div>
                </div>
            </section>
            <section v-show="step === 1">
                <div v-if="headers.length > 0" class="flex flex-row justify-center">
                    <VariablesCaster :headers="headers" ref="variblesCasterRef" />
                </div>
            </section>
            <section v-show="step === 2">
                <div class="flex flex-col gap-3">
                    <div v-if="!loading" class="flex flex-col gap-5">
                        <Dropdown label="Saved Queries" id="select-query" :options="querySelectOptions"
                            v-model="selectedQueryId" />
                        <div v-if="selectedQueryDescription">
                            <small>Query Description: {{ selectedQueryDescription }}</small>
                        </div>
                        <div v-if="hasVariablesAssigned" class="flex flex-row justify-center">
                            <Editor v-model="query" :show-binded-sql="true" :variables="variables" :data="content"
                                ref="editorRef">
                                <div class="flex flex-row justify-end gap-3 w-full">
                                    <div v-if="!loadingDatabaseConnection" class="flex flex-row gap-3 justify-end">
                                        <Button v-if="databaseConnection && databaseConnection.ID !== 0" type="button"
                                            @click="testInputSql">Test Input SQL</Button>
                                        <Button v-if="databaseConnection && databaseConnection.ID !== 0" type="button"
                                          @click="testOutputSql">Test Output SQL</Button>
                                    </div>
                                    <div v-if="loadingDatabaseConnection">
                                        <Loader />
                                    </div>
                                    <div v-if="hasEditor" class="flex flex-row justify-center">
                                        <Button type="button" @click="createSqlFile">Save .SQL</Button>
                                    </div>
                                </div>
                            </Editor>
                        </div>
                        <div v-else>
                            <Loader />
                        </div>
                    </div>
                </div>
            </section>
        </Steps>
        <SqlResultTable :data="responseTest" v-model="showSqlTable" />
    </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref, watch } from 'vue';

import { ReadXLSXFile, CreateSQLFile, GetQueriesList, GetDatabaseConnection, TestQueryInDatabase, TestBatchQueryInDatabase } from '../../wailsjs/go/main/App'
import { main } from '../../wailsjs/go/models';

import Table from '../components/Table.vue';
import VariablesCaster from '../components/VariablesCaster.vue';
import Editor from '../components/Editor.vue';
import Steps from '../components/Steps.vue';
import Button from '../components/Button.vue';
import Dropdown from '../components/Dropdown.vue';
import Loader from '../components/Loader.vue';
import SqlResultTable from '../components/SqlResultTable.vue';

const loading = ref<boolean>(false);
const loadingDatabaseConnection = ref<boolean>(false);
const content = ref<{ [k: string]: string }[]>([]);
const variblesCasterRef = ref<typeof VariablesCaster | null>(null);
const editorRef = ref<typeof Editor | null>(null);
const query = ref<string>('SELECT * from family limit 1;');
const showSqlTable = ref<boolean>(false);

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
    const querySelected = queries.value.find((query) => query.ID === Number (selectedQueryId.value));

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
        await testSQL(query.value)
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