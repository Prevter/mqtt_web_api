<script setup>
import { getRouter } from '../router';
const router = getRouter();

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

logout();
</script>

<template>
    
</template>