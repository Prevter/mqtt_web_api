<script>
export default {
    name: 'NavSidebar',
    props: {
        pages: {
            type: Array,
            required: true
        },
        title: {
            type: String,
            required: true
        }
    },
    data() {
        return {

        }
    },
    methods: {
        openPage(page) {
            this.$emit('open-page', page);
        },
        logout() {
            this.$emit('logout');
        }
    }
}
</script>

<template>
    <div class="sidebar flex-shrink-0 p-3">
        <span class="d-flex align-items-center pb-3 mb-3 text-decoration-none border-bottom">
            <i class="fa fa-envira fs-3 me-2"></i>
            <span class="fs-5 fw-semibold">{{ title }}</span>
        </span>
        <ul class="list-unstyled ps-0">

            <li class="mb-1" v-for="category in pages">
                <button class="btn btn-toggle align-items-center rounded" data-bs-toggle="collapse"
                    :data-bs-target="`#${category.id}-collapse`" aria-expanded="true">
                    {{ category.title }}
                </button>
                <div :class="({ 'collapse show': category.active, 'collapse': !category.active })"
                    :id="`${category.id}-collapse`" style="">
                    <ul class="btn-toggle-nav list-unstyled fw-normal pb-1 small">
                        <li v-for="page in category.pages">
                            <a :href="`/${page.path}`" class="rounded d-flex" @click.prevent="openPage(page)">
                                <span style="width: 1rem;" class="d-flex">
                                    <i :class="page.icon" class="m-auto"></i>
                                </span>
                                <span class="ms-2">{{ page.title }}</span>
                            </a>
                        </li>
                    </ul>
                </div>
            </li>
        </ul>
    </div>
</template>

<style>
:root {
    --nav-bg: linear-gradient(-45deg, #e3f2fd 0%, #bbdefb 100%);
    --nav-color: rgba(0, 0, 0, .95);
    --nav-hover-bg: rgba(0, 0, 0, .075);
    --nav-link-hover-color: rgba(0, 0, 0, .65);
    --nav-link-hover-bg: transparent;
    --nav-link-active-color: rgba(0, 0, 0, .95);
    --nav-link-active-bg: rgba(0, 0, 0, .075);
}

:root[data-bs-theme='dark'] {
    --nav-bg: linear-gradient(-45deg, #20061d 0%, #200e49 100%);
    --nav-color: rgba(255, 255, 255, .95);
    --nav-hover-bg: rgba(255, 255, 255, .075);
    --nav-link-hover-color: rgba(255, 255, 255, .65);
    --nav-link-hover-bg: transparent;
    --nav-link-active-color: rgba(255, 255, 255, .95);
    --nav-link-active-bg: rgba(255, 255, 255, .075);
}

body {
    min-height: 100vh;
    min-height: -webkit-fill-available;
}

html {
    height: -webkit-fill-available;
}

main {
    display: flex;
    flex-wrap: nowrap;
    height: 100vh;
    height: -webkit-fill-available;
    max-height: 100vh;
    overflow-x: auto;
    overflow-y: hidden;
}

.b-example-divider {
    flex-shrink: 0;
    width: 1.5rem;
    height: 100vh;
    background-color: rgba(0, 0, 0, .1);
    border: solid rgba(0, 0, 0, .15);
    border-width: 1px 0;
    box-shadow: inset 0 .5em 1.5em rgba(0, 0, 0, .1), inset 0 .125em .5em rgba(0, 0, 0, .15);
}

.bi {
    vertical-align: -.125em;
    pointer-events: none;
    fill: currentColor;
}

.dropdown-toggle {
    outline: 0;
}

.nav-flush .nav-link {
    border-radius: 0;
}

.btn-toggle {
    display: inline-flex;
    align-items: center;
    padding: .25rem .5rem;
    font-weight: 600;
    color: var(--nav-color);
    background-color: transparent;
    border: 0;
}

.btn-toggle:hover,
.btn-logout:hover {
    color:  var(--nav-color);
    background-color: var(--nav-hover-bg);
}

.btn-toggle::before {
    width: 1.25em;
    line-height: 0;
    transition: transform .35s ease;
    transform-origin: .5em 50%;
    stroke: rgba(255, 0, 0, .25);
}

.btn-toggle[aria-expanded="true"] {
    color: var(--nav-color);
}

.btn-toggle[aria-expanded="true"]::before {
    transform: rotate(90deg);
}

.btn-toggle-nav a {
    display: inline-flex;
    padding: .1875rem .5rem;
    margin-top: .125rem;
    margin-left: 1.25rem;
    text-decoration: none;
}

.btn-toggle-nav a:hover,
.btn-toggle-nav a:focus {
    background-color: var(--nav-link-hover-bg);
}

.scrollarea {
    overflow-y: auto;
}

.fw-semibold {
    font-weight: 600;
}

.lh-tight {
    line-height: 1.25;
}

.sidebar a {
    color: var(--nav-color);
}

.sidebar a:hover {
    color: var(--nav-link-hover-color);
    background-color: var(--nav-link-hover-bg);
}

.sidebar {
    min-height: 100svh;
    height: 100hmax;
    padding: 48px 0 0;
    box-shadow: inset -1px 0 0 rgba(0, 0, 0, .1);
    background: var(--nav-bg);
    transition: transform .3s ease-in-out;
    overflow-y: auto;
    overflow-x: hidden;
    width: max-content;
}
</style>
