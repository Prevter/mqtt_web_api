<script setup>
import { ref } from 'vue';
import { getRouter } from '../router';
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
        });
}

</script>

<template>
    <div class="container mx-auto px-4 py-8">
        <div class="flex justify-end">
            <button class="bg-red-500 hover:bg-red-700 text-white font-bold py-2 px-4 rounded"
                @click="logout">Logout</button>
        </div>
        <div class="mt-8">
            <div class="bg-gray-200 dark:bg-gray-800 rounded-lg p-4">
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
                <div class="flex items-center mb-4">
                    <input type="text" id="console-input"
                        class="dark:bg-gray-700 dark:text-white rounded-lg py-2 px-4 w-full" v-model="consoleInput">
                    <button
                        class="text-black hover:text-gray-800 dark:hover:text-blue-200 dark:text-white font-bold mt-auto py-2 px-4 rounded"
                        @click="submitConsoleInput">Submit</button>
                </div>
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