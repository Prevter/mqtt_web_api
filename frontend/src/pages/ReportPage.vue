<script setup>
import BarChart from '../components/BarChart.vue';
import { getRouter } from '../router';
const router = getRouter();
</script>

<script>
import { getRouter } from '../router';
import { ref } from 'vue';
const router = getRouter();

const stations = ref([]);
const loading = ref(false);
const error = ref('');

function print(name, element) {
    console.log(element)
    const printWindow = window.open('', '', 'height=400,width=800');
    printWindow.document.write(`<html><head><title>${name}</title>`);
    printWindow.document.write('</head><body >');
    printWindow.document.write('<style>td,tr,th { border: 1px solid black; }</style>');
    printWindow.document.write(`<h2>${name}</h2>`);
    printWindow.document.write(element.outerHTML);
    printWindow.document.write('</body></html>');
    printWindow.document.close();
    printWindow.print();
    printWindow.close();
}

function loadStations() {
    loading.value = true;
    error.value = '';
    fetch(`${API_URL}/select/station?limit=100`, {
        method: 'GET',
        credentials: 'include'
    })
        .then(response => response.json())
        .then(data => {
            if (data.error) {
                if (data.status === 401) router.replace('/logout');
            } else {
                stations.value = data.rows;
                selectedStation.value = stations.value[0].station_name;
            }
        })
        .catch(err => {
            console.error(err);
            error.value = err;
        })
        .finally(() => {
            loading.value = false;
        });
}

const report1 = ref([]);

function loadReport1() {
    loading.value = true;
    error.value = '';
    fetch(`${API_URL}/report/station`, {
        method: 'GET',
        credentials: 'include'
    })
        .then(response => response.json())
        .then(data => {
            if (data.error) {
                if (data.status === 401) router.replace('/logout');
                error.value = data.error;
            } else {
                report1.value = data.data ?? [];
            }
        })
        .catch(err => {
            console.error(err);
            error.value = err;
        })
        .finally(() => {
            loading.value = false;
        });
}

const selectedStation = ref(null);
const startDate = ref('2022-01-01T00:00');
const endDate = ref('2024-01-01T00:00');
const report2 = ref([]);

function loadReport2() {
    loading.value = true;
    error.value = '';
    fetch(`${API_URL}/report/measurement?station=${selectedStation.value}&start=${startDate.value.replace('T', ' ')}&end=${endDate.value.replace('T', ' ')}`, {
        method: 'GET',
        credentials: 'include'
    })
        .then(response => response.json())
        .then(data => {
            if (data.error) {
                if (data.status === 401) router.replace('/logout');
                error.value = data.error;
            } else {
                report2.value = data.data ?? [];
            }
        })
        .catch(err => {
            console.error(err);
            error.value = err;
        })
        .finally(() => {
            loading.value = false;
        });
}

const report3 = ref({ labels: [], datasets: [] });

function parseReport3(data) {
    const result = {
        labels: [],
        datasets: []
    };
    console.log(data)
    const cities = [...new Set(data.map(d => d.city))];
    const units = [...new Set(data.map(d => d.unit))];
    units.forEach(unit => {
        result.datasets.push({
            label: unit,
            data: [],
            backgroundColor: '#' + Math.floor(Math.random() * 16777215).toString(16)
        });
    });
    cities.forEach(city => {
        result.labels.push(city);
        units.forEach(unit => {
            const max = data.find(d => d.city === city && d.unit === unit)?.max ?? 0;
            result.datasets.find(d => d.label === unit).data.push(max);
        });
    });
    report3.value = result;
}

function loadReport3() {
    loading.value = true;
    error.value = '';
    fetch(`${API_URL}/report/maxparticles?start=${startDate.value.replace('T', ' ')}&end=${endDate.value.replace('T', ' ')}`, {
        method: 'GET',
        credentials: 'include'
    })
        .then(response => response.json())
        .then(data => {
            if (data.error) {
                if (data.status === 401) router.replace('/logout');
                error.value = data.error;
            } else {
                parseReport3(data.data ?? []);
            }
        })
        .catch(err => {
            console.error(err);
            error.value = err;
        })
        .finally(() => {
            loading.value = false;
        });

}

const report4empty = ref(false);
const report4 = ref({ labels: [], datasets: [] });

const qualityTranslate = {
    'Excellent': 'Відмінний',
    'Fine': 'Добрий',
    'Moderate': 'Середній',
    'Poor': 'Поганий',
    'Very Poor': 'Дуже поганий',
    'Severe': 'Жахливий'
}

function parseReport456(data, target, targetEmpty) {
    targetEmpty.value = data === null || data.length === 0;
    const result = {
        labels: [],
        datasets: [
            {
                label: 'Кількість фіксацій',
                data: [],
                backgroundColor: '#' + Math.floor(Math.random() * 16777215).toString(16)
            }
        ]
    };
    data.forEach(d => {
        result.labels.push(qualityTranslate[d.designation]);
        result.datasets[0].data.push(d.count);
    });
    target.value = result;
}

function loadReport4() {
    loading.value = true;
    error.value = '';
    fetch(`${API_URL}/report/badparticles?station=${selectedStation.value}`, {
        method: 'GET',
        credentials: 'include'
    })
        .then(response => response.json())
        .then(data => {
            if (data.error) {
                if (data.status === 401) router.replace('/logout');
                error.value = data.error;
            } else {
                parseReport456(data.data ?? [], report4, report4empty);
            }
        })
        .catch(err => {
            console.error(err);
            error.value = err;
        })
        .finally(() => {
            loading.value = false;
        });
}

const report5 = ref({ labels: [], datasets: [] });
const report5empty = ref(false);

function loadReport5() {
    loading.value = true;
    error.value = '';
    fetch(`${API_URL}/report/sulfur?station=${selectedStation.value}`, {
        method: 'GET',
        credentials: 'include'
    })
        .then(response => response.json())
        .then(data => {
            if (data.error) {
                if (data.status === 401) router.replace('/logout');
                error.value = data.error;
            } else {
                parseReport456(data.data ?? [], report5, report5empty);
            }
        })
        .catch(err => {
            console.error(err);
            error.value = err;
        })
        .finally(() => {
            loading.value = false;
        });
}

const report6 = ref({ labels: [], datasets: [] });
const report6empty = ref(false);

function loadReport6() {
    loading.value = true;
    error.value = '';
    fetch(`${API_URL}/report/carbon?station=${selectedStation.value}`, {
        method: 'GET',
        credentials: 'include'
    })
        .then(response => response.json())
        .then(data => {
            if (data.error) {
                if (data.status === 401) router.replace('/logout');
                error.value = data.error;
            } else {
                parseReport456(data.data ?? [], report6, report6empty);
            }
        })
        .catch(err => {
            console.error(err);
            error.value = err;
        })
        .finally(() => {
            loading.value = false;
        });
}

function reroute(option) {
    if (option == 'stations') {
        loadReport1();
    } else if (option == 'measurements' || option == 'badparticles' || option == 'carbon' || option == 'sulfur') {
        loadStations();
    }
}

// monitor page changes
router.afterEach((to, from) => {
    reroute(to.params.options);
});

reroute(router.currentRoute.value.params.options);
</script>

<template>
    <div class="container mx-auto px-4 py-3" v-if="router.currentRoute.value.params.options == 'stations'">
        <div class="d-flex justify-content-between">
            <h2>Список підключених станцій</h2>
            <div class="d-flex">
                <button class="btn btn-success ms-2" :disabled="report1.length == 0"
                    @click="print(`Список підключених станцій`, $refs.report1table)">
                    <i class="fa fa-print"></i>
                </button>
            </div>
        </div>

        <div v-if="loading" class="d-flex justify-content-center my-5">
            <div class="spinner-border text-primary" role="status">
                <span class="visually-hidden">Loading...</span>
            </div>
        </div>

        <div v-else-if="error" class="alert alert-danger" role="alert">
            {{ error }}
        </div>

        <div v-else-if="report1.length" class="table-responsive">
            <table class="table table-fixed table-hover bg-transparent" ref="report1table">
                <thead>
                    <tr>
                        <th scope="col" class="col-3">
                            <i class="fa fa-map-marker"></i>
                            Назва станції
                        </th>
                        <th scope="col" class="col-3">
                            <i class="fa fa-clock-o"></i>
                            Перший вимір
                        </th>
                        <th scope="col" class="col-6">
                            <i class="fa fa-thermometer-half"></i>
                            Список параметрів
                        </th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="row in report1" :key="row.unit">
                        <td>{{ row['station_name'] }}</td>
                        <td v-if="row['first_date']">
                            {{ new Date(row['first_date']).toLocaleString() }}
                        </td>
                        <td v-else>
                            <i class="fa fa-times"></i>
                            Н/Д
                        </td>
                        <td v-if="row['units']">
                            {{ row['units'] }}
                        </td>
                        <td v-else>
                            <i class="fa fa-times"></i>
                            Н/Д
                        </td>
                    </tr>
                </tbody>
            </table>
        </div>

    </div>
    <div class="container mx-auto px-4 py-3" v-else-if="router.currentRoute.value.params.options == 'measurements'">
        <h2>Результати вимірювань станції за часовий період</h2>

        <div class="row g-3 align-items-center mb-2">
            <div class="col-auto">
                Станція:
            </div>
            <div class="col-auto">
                <select class="form-select" v-model="selectedStation">
                    <option selected disabled>Оберіть станцію</option>
                    <option v-for="station in stations" :value="station.station_name">{{ station.station_name }}</option>
                </select>
            </div>
            <div class="col-auto">
                Часовий період:
            </div>
            <div class="col-auto">
                <input type="datetime-local" class="form-control" v-model="startDate" />
            </div>
            <div class="col-auto">
                <input type="datetime-local" class="form-control" v-model="endDate" />
            </div>
            <div class="col-auto">
                <button class="btn btn-primary" @click="loadReport2" :disabled="loading">Застосувати</button>
                <button class="btn btn-success ms-2" :disabled="report2.length == 0"
                    @click="print(`Результати вимірювань станції ${selectedStation} з ${startDate.replace('T', ' ')} по ${endDate.replace('T', ' ')}`, $refs.report2table)">
                    <i class="fa fa-print"></i>
                </button>
            </div>
        </div>

        <div v-if="loading" class="d-flex justify-content-center my-5">
            <div class="spinner-border text-primary" role="status">
                <span class="visually-hidden">Loading...</span>
            </div>
        </div>

        <div v-else-if="error" class="alert alert-danger" role="alert">
            {{ error }}
        </div>

        <div v-else-if="report2.length" class="table-responsive">
            <table class="table table-hover bg-transparent" ref="report2table">
                <thead>
                    <tr>
                        <th>Назва (одиниці)</th>
                        <th>Мінімальне значення</th>
                        <th>Середнє значення</th>
                        <th>Максимальне значення</th>
                        <th>Розмір вибірки</th>
                    </tr>
                </thead>
                <tbody>
                    <tr v-for="row in report2" :key="row.unit">
                        <td>{{ row['unit'] }}</td>
                        <td>{{ row['min'] }}</td>
                        <td>{{ Math.round(row['avg'] * 100) / 100 }}</td>
                        <td>{{ row['max'] }}</td>
                        <td>{{ row['count'] }}</td>
                    </tr>
                </tbody>
            </table>
        </div>
    </div>
    <div class="container mx-auto px-4 py-3" v-if="router.currentRoute.value.params.options == 'maxparticles'">
        <h2>Максимальне значення шкідливих частинок PM2.5, PM10</h2>
        <div class="row g-3 align-items-center mb-2">
            <div class="col-auto">
                Часовий період:
            </div>
            <div class="col-auto">
                <input type="datetime-local" class="form-control" v-model="startDate" />
            </div>
            <div class="col-auto">
                <input type="datetime-local" class="form-control" v-model="endDate" />
            </div>
            <div class="col-auto">
                <button class="btn btn-primary" @click="loadReport3" :disabled="loading">Застосувати</button>
            </div>
        </div>

        <div v-if="loading" class="d-flex justify-content-center my-5">
            <div class="spinner-border text-primary" role="status">
                <span class="visually-hidden">Loading...</span>
            </div>
        </div>

        <div v-else-if="error" class="alert alert-danger" role="alert">
            {{ error }}
        </div>

        <div v-else-if="report3.labels.length">
            <div>
                <table class="table table-hover bg-transparent">
                    <thead>
                        <tr>
                            <th>Місто</th>
                            <th>Максимальне значення PM2.5</th>
                            <th>Максимальне значення PM10</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="city, i of report3.labels" :key="city">
                            <td>{{ city }}</td>
                            <td>{{ report3.datasets[1].data[i] }}</td>
                            <td>{{ report3.datasets[0].data[i] }}</td>
                        </tr>
                    </tbody>
                </table>
            </div>
            <div>
                <BarChart :data="report3" />
            </div>
        </div>
    </div>
    <div class="container mx-auto px-4 py-3" v-if="router.currentRoute.value.params.options == 'badparticles'">
        <h2>Кількість фіксацій шкідливого рівня частинок PM2.5</h2>

        <div class="row g-3 align-items-center mb-2">
            <div class="col-auto">
                Станція:
            </div>
            <div class="col-auto">
                <select class="form-select" v-model="selectedStation">
                    <option selected disabled>Оберіть станцію</option>
                    <option v-for="station in stations" :value="station.station_name">{{ station.station_name }}</option>
                </select>
            </div>
            <div class="col-auto">
                <button class="btn btn-primary" @click="loadReport4" :disabled="loading">Застосувати</button>
            </div>
        </div>

        <div v-if="loading" class="d-flex justify-content-center my-5">
            <div class="spinner-border text-primary" role="status">
                <span class="visually-hidden">Loading...</span>
            </div>
        </div>

        <div v-else-if="error" class="alert alert-danger" role="alert">
            {{ error }}
        </div>

        <div v-else-if="report4empty" class="alert alert-warning" role="alert">
            Немає даних
        </div>

        <div v-else-if="report4.labels.length" class="row">
            <div class="col">
                <table class="table table-hover bg-transparent">
                    <thead>
                        <tr>
                            <th>Категорія</th>
                            <th>Кількість фіксацій</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="category, i of report4.labels" :key="category">
                            <td>{{ category }}</td>
                            <td>{{ report4.datasets[0].data[i] }}</td>
                        </tr>
                    </tbody>
                </table>
            </div>
            <div class="col">
                <BarChart :data="report4" />
            </div>
        </div>
    </div>
    <div class="container mx-auto px-4 py-3" v-if="router.currentRoute.value.params.options == 'sulfur'">
        <h2>Кількість фіксацій оптимальних значень діоксиду сірки</h2>

        <div class="row g-3 align-items-center mb-2">
            <div class="col-auto">
                Станція:
            </div>
            <div class="col-auto">
                <select class="form-select" v-model="selectedStation">
                    <option selected disabled>Оберіть станцію</option>
                    <option v-for="station in stations" :value="station.station_name">{{ station.station_name }}</option>
                </select>
            </div>
            <div class="col-auto">
                <button class="btn btn-primary" @click="loadReport5" :disabled="loading">Застосувати</button>
            </div>
        </div>

        <div v-if="loading" class="d-flex justify-content-center my-5">
            <div class="spinner-border text-primary" role="status">
                <span class="visually-hidden">Loading...</span>
            </div>
        </div>

        <div v-else-if="error" class="alert alert-danger" role="alert">
            {{ error }}
        </div>

        <div v-else-if="report5empty" class="alert alert-warning" role="alert">
            Немає даних
        </div>

        <div v-else-if="report5.labels.length" class="row">
            <div class="col">
                <table class="table table-hover bg-transparent">
                    <thead>
                        <tr>
                            <th>Категорія</th>
                            <th>Кількість фіксацій</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="category, i of report5.labels" :key="category">
                            <td>{{ category }}</td>
                            <td>{{ report5.datasets[0].data[i] }}</td>
                        </tr>
                    </tbody>
                </table>
            </div>
            <div class="col">
                <BarChart :data="report5" />
            </div>
        </div>
    </div>
    <div class="container mx-auto px-4 py-3" v-if="router.currentRoute.value.params.options == 'carbon'">
        <h2>Кількість фіксацій оптимальних значень чадного газу</h2>

        <div class="row g-3 align-items-center mb-2">
            <div class="col-auto">
                Станція:
            </div>
            <div class="col-auto">
                <select class="form-select" v-model="selectedStation">
                    <option selected disabled>Оберіть станцію</option>
                    <option v-for="station in stations" :value="station.station_name">{{ station.station_name }}</option>
                </select>
            </div>
            <div class="col-auto">
                <button class="btn btn-primary" @click="loadReport6" :disabled="loading">Застосувати</button>
            </div>
        </div>

        <div v-if="loading" class="d-flex justify-content-center my-5">
            <div class="spinner-border text-primary" role="status">
                <span class="visually-hidden">Loading...</span>
            </div>
        </div>

        <div v-else-if="error" class="alert alert-danger" role="alert">
            {{ error }}
        </div>

        <div v-else-if="report6empty" class="alert alert-warning" role="alert">
            Немає даних
        </div>

        <div v-else-if="report6.labels.length" class="row">
            <div class="col">
                <table class="table table-hover bg-transparent">
                    <thead>
                        <tr>
                            <th>Категорія</th>
                            <th>Кількість фіксацій</th>
                        </tr>
                    </thead>
                    <tbody>
                        <tr v-for="category, i of report6.labels" :key="category">
                            <td>{{ category }}</td>
                            <td>{{ report6.datasets[0].data[i] }}</td>
                        </tr>
                    </tbody>
                </table>
            </div>
            <div class="col">
                <BarChart :data="report6" />
            </div>
        </div>
    </div>
</template>