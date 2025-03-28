<script setup>
import PageNavigator from '../components/PageNavigator.vue';
import { ref } from 'vue';

const loading = ref(true);
const error = ref('');
const result = ref([]);
const page = ref(1);
const total_pages = ref(1);

function changePage(new_page) {
    if (new_page < 1) return;
    if (new_page > total_pages.value) return;
    if (new_page === page.value) return;
    page.value = new_page;
    loadStations();
}

function loadStations() {
    loading.value = true;
    error.value = '';
    result.value = [];

    fetch(`${API_URL}/select/station?page=${page.value - 1}`, {
        method: 'GET',
        credentials: 'include'
    })
        .then(response => response.json())
        .then(data => {
            if (data.error) {
                error.value = data.error;
                if (data.status === 401) router.replace('/logout');
            } else {
                total_pages.value = data.total_pages;
                if (data.rows === null) {
                    result.value = [];
                    return
                }
                else result.value = data.rows;
            }
        })
        .catch(error => {
            console.error(error);
        })
        .finally(() => {
            loading.value = false;
        });
}

loadStations();

</script>

<template>
    <div class="container mx-auto px-4 py-3">
        <h2>Станції моніторингу</h2>

        <div class="container">
            <div v-if="loading" class="d-flex justify-content-center my-5">
                <div class="spinner-border text-primary" role="status">
                    <span class="visually-hidden">Loading...</span>
                </div>
            </div>

            <div v-else-if="error" class="alert alert-danger" role="alert">
                {{ error }}
            </div>

            <div v-else-if="result.length === 0" class="alert alert-warning" role="alert">
                Немає даних
            </div>

            <div v-else class="table-responsive">
                <table class="table table-hover bg-transparent">
                    <thead>
                        <tr>
                            <th>#</th>
                            <th>Назва станції</th>
                            <th>Розташування</th>
                            <th>Деталі</th>
                            <th>Опції</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="row in result" :key="row.id">
                            <!-- <td v-for="key in Object.keys(row)" :key="key">{{ row[key] }}</td> -->
                            <td>{{ row['id_station'] }}</td>
                            <td>{{ row['station_name'] }}</td>
                            <td>
                                {{ row['city'] }}<br/>
                                <a :href="`https://www.google.com/maps/search/?api=1&query=${row['coordinates']['lon']},${row['coordinates']['lat']}`"
                                    target="_blank" class="text-decoration-none">
                                    <i class="fa fa-map-marker"></i>
                                    {{ row['coordinates']['lon'] }}, {{ row['coordinates']['lat'] }}
                                </a>
                            </td>
                            <td>
                                <span>
                                    <i :class="{
                                        'fa fa-check-circle': row['station_status'] === 'ACTIVE',
                                        'fa fa-times-circle': row['station_status'] === 'INACTIVE'
                                    }"></i>
                                    {{ row['station_status'] === 'ACTIVE' ? 'Активна' : 'Неактивна' }}
                                    <br />
                                </span>
                                <span v-if="row['server_url']">
                                    <i class="fa fa-server"></i>
                                    {{ row['server_url'] }} ({{ row['server_status'] === 'ACTIVE' ? 'активний' :
                                        'неактивний' }})
                                    <br />
                                </span>
                                <span v-if="row['id_saveecobot'] !== 'NULL'">
                                    <i class="fa fa-globe"></i>
                                    {{ row['id_saveecobot'] }}
                                </span>
                            </td>
                            <td>
                                <button class="btn btn-sm btn-text-success" @click="$emit('navigate', `measurements/station:${row['id_station']}`)">
                                    <i class="fa fa-eye"></i>
                                    Вимірювання
                                </button>
                            </td>
                        </tr>
                    </tbody>
                </table>

                <PageNavigator :page="page" :total_pages="total_pages" @change-page="changePage" />
            </div>
        </div>
    </div>
</template>

<style>
table, th, td {
    background-color: transparent !important;
}

.btn-text-success {
    color: #28a745;
}

.btn-text-success:hover {
    color: #1e7e34;
}
</style>