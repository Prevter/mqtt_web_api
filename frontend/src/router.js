import { createRouter, createWebHistory } from 'vue-router'

const HomePage = () => import('./pages/HomePage.vue')
const LoginPage = () => import('./pages/LoginPage.vue')
const LogoutPage = () => import('./pages/LogoutPage.vue')

const routes = [
    { path: '/logout', component: LogoutPage },
    { path: '/login', component: LoginPage },
    { path: '/', component: HomePage },
    { path: '/:page', component: HomePage },
    { path: '/:page/:options', component: HomePage },
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