<script setup>
import PageNavigator from '../components/PageNavigator.vue';
import { getRouter } from '../router';
import { ref } from 'vue';
const router = getRouter();

const optionsStr = router.currentRoute.value.params.options;
let options = optionsStr ? optionsStr.split(',').map(o => { const e = o.split(':'); return { [e[0]]: e[1] } }) : [];

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
    loadMeasurements();
}

function loadMeasurements() {
    loading.value = true;
    error.value = '';
    result.value = [];

    const extra = options.map(o => `${Object.keys(o)[0]}=${Object.values(o)[0]}`).join('&');

    fetch(`${API_URL}/select/measurement?page=${page.value - 1}&${extra}`, {
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

function clearFilters() {
    options = [];
    router.replace('/measurements');
    loadMeasurements();
}

loadMeasurements();

</script>

<template>
    <div class="container mx-auto px-4 py-3">
        <div class="d-flex justify-content-between">
            <h2>Вимірювання</h2>
            <div class="d-flex">
                <button class="btn btn-danger" @click="clearFilters" v-if="options.length">
                    <i class="fa fa-trash"></i>
                    Очистити фільтри
                </button>
            </div>
        </div>

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
                            <th>Станція</th>
                            <th>Дата/час</th>
                            <th>Параметр</th>
                            <th>Значення</th>
                            <th>Стан</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="row in result" :key="row.id">
                            <td>{{ row['id_measurment'] }}</td>
                            <td>({{ row['id_station'] }}) {{ row['station_name'] }}</td>
                            <td>{{ new Date(row['measurment_time']).toLocaleString() }}</td>
                            <td>{{ row['unit_title'] }}</td>
                            <td>{{ row['measurment_value'] }} {{ row['unit_unit'] }}</td>
                            <td>
                                <span v-if="row['designation']" :class="{
                                    'excellent': row['designation'] === 'Excellent',
                                    'fine': row['designation'] === 'Fine',
                                    'moderate': row['designation'] === 'Moderate',
                                    'poor': row['designation'] === 'Poor',
                                    'very-poor': row['designation'] === 'Very Poor',
                                    'severe': row['designation'] === 'Severe'
                                }">
                                    <i :class="{
                                        'fa': true,
                                        'fa-smile-o': row['designation'] === 'Excellent' || row['designation'] === 'Fine',
                                        'fa-meh-o': row['designation'] === 'Moderate',
                                        'fa-frown-o': row['designation'] === 'Poor' || row['designation'] === 'Very Poor' || row['designation'] === 'Severe',
                                    }"></i>
                                    {{
                                        {
                                            'Excellent': 'Відмінний',
                                            'Fine': 'Добрий',
                                            'Moderate': 'Середній',
                                            'Poor': 'Поганий',
                                            'Very Poor': 'Дуже поганий',
                                            'Severe': 'Жахливий'
                                        }[row['designation']]
                                    }}
                                </span>
                                <span v-else>
                                    <i class="fa fa-question-circle"></i>
                                    Не визначено
                                </span>
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
.excellent {
    color: #0ecc0e;
}

.fine {
    color: #9cb153;
}

.moderate {
    color: #be991f;
}

.poor {
    color: #e27006;
}

.very-poor {
    color: #bd3e0c;
}

.severe {
    color: #ff0000;
}
</style>