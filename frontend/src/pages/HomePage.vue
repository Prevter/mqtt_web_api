<script setup>
import NavSidebar from '../components/NavSidebar.vue';
import StationsPage from './StationsPage.vue';
import MeasurementsPage from './MeasurementsPage.vue';
import ReportPage from './ReportPage.vue';

import { getRouter } from '../router';
import { ref } from 'vue';
const router = getRouter();

if (!localStorage.getItem('token')) router.replace('/login');

const username = ref(localStorage.getItem('username'));

const perms = {
    "postgres": {
        "home": true,
        "reports": true,
    },
    "ecographix": {
        "home": true,
        "reports": true,
    },
}

const role = perms[username.value] || {
    "home": false,
    "reports": true,
}

const pages = [];

if (role.home) pages.push({
    id: 'home',
    title: 'Дані',
    active: true,
    pages: [
        {
            title: 'Станції',
            path: 'stations',
            icon: 'fa fa-building'
        },
        {
            title: 'Вимірювання',
            path: 'measurements',
            icon: 'fa fa-database'
        },
    ]
});

if (role.reports) pages.push({
    id: 'reports',
    title: 'Звіти',
    active: true,
    pages: [
        {
            title: 'Список підключених станцій',
            path: 'report/stations',
            icon: 'fa fa-file-text-o'
        },
        {
            title: 'Результати вимірювань станції',
            path: 'report/measurements',
            icon: 'fa fa-file-text-o'
        },
        {
            title: 'Макс. значення шкідливих частинок',
            path: 'report/maxparticles',
            icon: 'fa fa-pie-chart'
        },
        {
            title: 'Фіксації шкідливого рівня частинок',
            path: 'report/badparticles',
            icon: 'fa fa-pie-chart'
        },
        {
            title: 'Фіксації діоксиду сірки',
            path: 'report/sulfur',
            icon: 'fa fa-pie-chart'
        },
        {
            title: 'Фіксації чадного газу',
            path: 'report/carbon',
            icon: 'fa fa-pie-chart'
        }
    ]
});

const page = ref('');
page.value = router.currentRoute.value.params.page;

if (!page.value) {
    router.replace('/stations');
    page.value = 'stations';
}

function logout() {
    router.replace('/logout')
}

function navigate(path) {
    console.log(path);
    path = path.path || path;
    const splitPath = path.split('/');

    if (page.value === splitPath[0]) {
        router.replace(`/${path}/`);
        router.currentRoute.value.params.options = splitPath.length > 1 ? splitPath[1] : undefined;
        return;
    }

    page.value = splitPath[0];
    router.currentRoute.value.params.options = splitPath.length > 1 ? splitPath[1] : undefined;
    router.replace(`/${path}/`);
}

</script>

<template>
    <div class="d-flex">
        <NavSidebar :pages="pages" title="EcoGraphix" v-on:open-page="navigate" v-on:logout="logout" />
        <div class="w-100">
            <header class="d-flex">
                <div class="ms-auto d-flex">
                    <div class="my-auto me-2">
                        <i class="fa fa-user-circle-o"></i>
                        {{ username }}
                    </div>
                    <button class="btn btn-sm btn-outline-danger my-auto me-2" @click="logout">
                        <i class="fa fa-power-off"></i>
                        Вихід
                    </button>
                </div>
            </header>
            <div class="container">
                <StationsPage v-if="page === 'stations' && role.home" v-on:navigate="navigate" />
                <MeasurementsPage v-if="page === 'measurements' && role.home" v-on:navigate="navigate" />
                <ReportPage v-if="page === 'report' && role.reports" v-on:navigate="navigate" />
            </div>
        </div>
    </div>
</template>

<style>
:root {
    --header-bg: linear-gradient(90deg, #bbdefb 0%, #93bde0 100%);
}

:root[data-bs-theme='dark'] {
    --header-bg: linear-gradient(90deg, #1f0c3d 0%, #2b1166 100%);
}

header {
    display: block;
    height: 60px;
    width: 100%;
    background: var(--header-bg);
    border: 1px solid rgba(0, 0, 0, 0.3);
}
</style>