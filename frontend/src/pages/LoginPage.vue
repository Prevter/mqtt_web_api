<script setup>
import { ref } from 'vue';
import { getRouter } from '../router';
const router = getRouter();

if (localStorage.getItem('token')) router.replace('/');

const username = ref('');
const password = ref('');
const loading = ref(false);
const error = ref('');

function login() {
    loading.value = true;
    error.value = '';

    // www-form-urlencoded
    const formData = new FormData();
    formData.append('username', username.value);
    formData.append('password', password.value);

    fetch(`${API_URL}/auth`, {
        method: 'POST',
        body: formData,
        credentials: 'include'
    })
        .then(response => response.json())
        .then(data => {
            if (data.error) {
                error.value = data.error;
            } else {
                localStorage.setItem('token', data.token);
                router.replace('/');
            }
        })
        .catch(error => {
            console.error(error);
            error.value = 'Помилка з\'єднання з сервером';
        })
        .finally(() => {
            loading.value = false;
        });
}
</script>

<template>
    <div class="flex min-h-full flex-1 flex-col justify-center px-6 py-12 lg:px-8">
        <div class="sm:mx-auto sm:w-full sm:max-w-sm">
            <img class="mx-auto h-10 w-auto" src="https://tailwindui.com/img/logos/mark.svg?color=green&shade=600"
                alt="Your Company" />
            <h2 class="mt-10 text-center text-2xl font-bold leading-9 tracking-tight">
                Вхід
            </h2>
        </div>

        <div class="mt-4 sm:mx-auto sm:w-full sm:max-w-sm">
            <div v-if="error" class="flex items-center bg-red-500 text-white text-sm font-bold px-4 py-3 mb-3" role="alert">
                <svg class="fill-current h-6 w-6 mr-2" xmlns="http://www.w3.org/2000/svg" height="24"
                    viewBox="0 -960 960 960" width="24">
                    <path
                        d="M480-280q17 0 28.5-11.5T520-320q0-17-11.5-28.5T480-360q-17 0-28.5 11.5T440-320q0 17 11.5 28.5T480-280Zm-40-160h80v-240h-80v240Zm40 360q-83 0-156-31.5T197-197q-54-54-85.5-127T80-480q0-83 31.5-156T197-763q54-54 127-85.5T480-880q83 0 156 31.5T763-763q54 54 85.5 127T880-480q0 83-31.5 156T763-197q-54 54-127 85.5T480-80Zm0-80q134 0 227-93t93-227q0-134-93-227t-227-93q-134 0-227 93t-93 227q0 134 93 227t227 93Zm0-320Z" />
                </svg>
                <p>{{ error }}</p>
            </div>

            <form class="space-y-6" action="#" method="POST" @submit="login">
                <div>
                    <label for="username" class="block text-sm font-medium leading-6">Логін</label>
                    <div class="mt-2">
                        <input id="username" name="username" type="text" autocomplete="email" required="" v-model="username"
                            class="block w-full rounded-md border-1 py-1.5 shadow-sm sm:text-sm sm:leading-6 dark:bg-gray-900 dark:border-gray-700" />
                    </div>
                </div>

                <div>
                    <label for="password" class="block text-sm font-medium leading-6">Пароль</label>
                    <div class="mt-2">
                        <input id="password" name="password" type="password" autocomplete="current-password" required="" v-model="password"
                            class="block w-full rounded-md border-1 py-1.5 shadow-sm sm:text-sm sm:leading-6 dark:bg-gray-900 dark:border-gray-700" />
                    </div>
                </div>

                <div>
                    <button type="submit" @click.prevent="login" :enabled="!loading"
                        class="flex w-full justify-center rounded-md bg-indigo-600 px-3 py-1.5 text-sm font-semibold leading-6 text-white shadow-sm hover:bg-indigo-500 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600">
                        Вхід
                    </button>
                </div>
            </form>
        </div>
    </div>
</template>

<style scoped></style>