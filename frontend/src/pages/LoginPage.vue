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
                localStorage.setItem('username', username.value);
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
    <div class="container">
        <div class="row m-5 no-gutters shadow-lg bg-login">
            <div class="col-md-6 d-none d-md-block px-0">
                <img src="https://wallpapers.com/images/hd/captivating-multicolor-abstract-art-azyj40fu4hodo1ys.jpg"
                    class="img-fluid" style="min-height:100%;" />
            </div>
            <div class="col-md-6 p-5 my-auto">
                <h1 class="pb-3 text-center">
                    <i class="fa fa-envira"></i>
                    EcoGraphix
                </h1>

                <div class="alert alert-danger" role="alert" v-if="error">
                    {{ error }}
                </div>

                <div class="form-style">
                    <form method="POST" @submit="login">
                        <div class="form-group pb-3">
                            <input type="text" placeholder="Логін" required="" v-model="username" class="form-control">
                        </div>
                        <div class="form-group pb-3">
                            <input type="password" placeholder="Пароль" required="" v-model="password" class="form-control">
                        </div>

                        <div class="pb-2">
                            <button @click.prevent="login" :enabled="!loading" type="submit"
                                class="btn w-100 font-weight-bold mt-2">Увійти</button>
                        </div>
                    </form>
                </div>
            </div>
        </div>
        <div class="mx-auto sideline">
            <span>Розробив студент ІПЗ-21008Б Немеш Олександр</span>
        </div>
    </div>
</template>

<style scoped>
.form-style input {
    border: 0;
    height: 50px;
    border-radius: 0;
    border-bottom: 1px solid #ccc;
}

.form-style input:focus {
    border-bottom: 1px solid #007bff;
    box-shadow: none;
    outline: 0;
}

.sideline {
    display: flex;
    width: 100%;
    justify-content: center;
    align-items: center;
    text-align: center;
    color: #ccc;
}

button {
    height: 50px;
    border: 0;
    background-color: #43961c;
    color: white;
    font-weight: bold;
}

button:hover {
    background-color: #3c871b;
}

.sideline:before,
.sideline:after {
    content: '';
    border-top: 1px solid var(--input-border-bottom-color);
    margin: 0 20px 0 0;
    flex: 1 0 20px;
}

.sideline:after {
    margin: 0 0 0 20px;
}

.container>div {
    border-radius: 1rem;
    overflow: hidden;
}

.bg-login {
    background-color: #212529;
}
</style>