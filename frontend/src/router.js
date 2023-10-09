import { createRouter, createWebHistory } from 'vue-router'

const HomePage = () => import('./pages/HomePage.vue')
const LoginPage = () => import('./pages/LoginPage.vue')

const routes = [
    { path: '/', component: HomePage },
    { path: '/login', component: LoginPage },
]

let router = null

export function getRouter() {
    if (router) return router

    router = createRouter({
        history: createWebHistory(),
        routes,
    })
    return router
}