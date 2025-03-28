<script setup>
import { ref } from 'vue';
import { getRouter } from '../router';
import PageNavigator from '../components/PageNavigator.vue';
const router = getRouter();

if (!localStorage.getItem('token')) router.replace('/login');

function logout() {
    fetch(`${API_URL}/auth`, {
        method: 'DELETE',
        credentials: 'include'
    })
        .then(response => response.json())
        .then(_data => {
            localStorage.removeItem('token');
            router.replace('/login');
        })
        .catch(error => {
            console.error(error);
        });
}

const loading = ref(false);
const consoleInput = ref('SELECT * FROM Station');
const error = ref('');
const result = ref([]);
let keys;

function buildTableOutput() {
    if (typeof result.value === 'string') {
        return `<div class="bg-red-500 w-full text-white font-bold py-2 px-4 rounded">${result.value}</div>`;
    }

    let output = '<table class="table-auto w-full">';
    if (result.value.length > 0) {
        output += '<thead><tr>';
        for (const key in result.value[0]) {
            output += `<th class="px-4 py-2">${key}</th>`;
        }
        output += '</tr></thead>';
    }
    output += '<tbody>';
    for (const row of result.value) {
        output += '<tr>';
        for (const key in row) {
            output += `<td class="border px-4 py-2">${row[key]}</td>`;
        }
        output += '</tr>';
    }
    output += '</tbody></table>';
    return output;
}

function submitConsoleInput() {
    const formdata = new FormData();
    formdata.append('query', consoleInput.value);

    loading.value = true;

    fetch(`${API_URL}/console`, {
        method: 'POST',
        body: formdata,
        credentials: 'include'
    })
        .then(response => response.json())
        .then(data => {
            if (data.error) {
                error.value = data.error;
            } else {
                result.value = data.rows;
                keys = Object.keys(data.rows[0]);
            }
        })
        .catch(error => {
            console.error(error);
        })
        .finally(() => {
            loading.value = false;
        });
}

const showSqlTerminal = ref(false);

function openSqlTerminal() {
    showSqlTerminal.value = !showSqlTerminal.value;
    if (showSqlTerminal.value) {
        document.getElementById('console-input').focus();
    }
}

const lastTable = ref('');
const limit = ref(10);
const page = ref(1);
const total_pages = ref(1);

function updateTable() {
    loading.value = true;
    error.value = '';

    fetch(`${API_URL}/select/${lastTable.value}?limit=${limit.value}&page=${page.value - 1}`, {
        method: 'GET',
        credentials: 'include'
    })
        .then(response => response.json())
        .then(data => {
            if (data.error) {
                error.value = data.error;
            } else {
                total_pages.value = data.total_pages;
                if (data.rows === null) {
                    result.value = [];
                    return
                } 
                else result.value = data.rows;
                keys = Object.keys(data.rows[0]);
            }
        })
        .catch(error => {
            console.error(error);
        })
        .finally(() => {
            loading.value = false;
        });
}

function loadTable(tableName) {
    lastTable.value = tableName;
    error.value = '';
    result.value = [];
    page.value = 1;
    limit.value = 10;
    total_pages.value = 1;
    updateTable();
}

function changePage(new_page) {
    if (lastTable.value === '') return;
    if (new_page < 1) return;
    if (new_page > total_pages.value) return;
    page.value = new_page;
    updateTable();
}

</script>

<template>
    <div class="container mx-auto px-4 py-8">
        <div class="flex justify-between">
            <button class="bg-green-900 px-3 rounded" @click="openSqlTerminal">
                <span v-if="showSqlTerminal">
                    <i class="fa fa-terminal"></i>
                    Close SQL Mode
                </span>
                <span v-else>
                    <i class="fa fa-terminal"></i>
                    SQL Mode
                </span>
            </button>
            <button class="bg-red-500 hover:bg-red-700 text-white font-bold py-2 px-4 rounded"
                @click="logout">Logout</button>
        </div>
        <div class="mt-8">
            <div class="bg-gray-200 dark:bg-gray-800 rounded-lg px-4 py-1">
                <div class="flex items-center mb-4">
                    <div class="bg-red-500 w-full text-white font-bold py-2 px-4 rounded" v-if="error">
                        {{ error }}
                    </div>
                    <table class="table-auto w-full" v-else-if="result.length">
                        <thead>
                            <tr>
                                <th class="px-4 py-2" v-for="key in keys">{{ key }}</th>
                            </tr>
                        </thead>
                        <tbody>
                            <tr v-for="row in result">
                                <td class="border px-4 py-2" v-for="key in keys">{{ row[key] }}</td>
                            </tr>
                        </tbody>
                    </table>
                </div>

                <!-- Pagination -->
                <PageNavigator :total_pages="total_pages" :page="page" @change-page="changePage" v-if="!showSqlTerminal"/>

            </div>
        </div>


        <div class="flex items-center my-3" v-if="showSqlTerminal">
            <input type="text" id="console-input" class="dark:bg-gray-700 dark:text-white rounded-lg py-2 px-4 w-full"
                v-model="consoleInput">
            <button
                class="text-black hover:text-gray-800 dark:hover:text-blue-200 dark:text-white font-bold mt-auto py-2 px-4 rounded"
                @click="submitConsoleInput">Submit</button>
        </div>
        <div class="my-3" v-else>
            <h2 class="text-2xl font-bold">Tables</h2>
            <br />

            <div class="grid grid-cols-3 gap-4">
                <button class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
                    @click="loadTable('category')">Category</button>
                <button class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
                    @click="loadTable('coordinates')">Coordinates</button>
                <button class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
                    @click="loadTable('favorite')">Favorite</button>
                <button class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
                    @click="loadTable('measured_unit')">Measured Unit</button>
                <button class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
                    @click="loadTable('measurement')">Measurement</button>
                <button class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
                    @click="loadTable('mqtt_server')">MQTT Server</button>
                <button class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
                    @click="loadTable('mqtt_unit')">MQTT Unit</button>
                <button class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
                    @click="loadTable('optimal_value')">Optimal Value</button>
                <button class="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
                    @click="loadTable('station')">Station</button>
            </div>
        </div>

    </div>
</template>

<style scoped>
.console-input {
    @apply dark:bg-gray-700 dark:text-white rounded-lg py-2 px-4 w-full;
}

.console-input:focus {
    @apply outline-none;
}

.console-input:focus-visible {
    @apply ring-2 ring-blue-500;
}
</style>