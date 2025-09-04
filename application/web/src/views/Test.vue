<template>
  <div>
    <UserCard :user="user" />
  </div>
</template>

<script setup lang="ts">
import UserCard from '../components/UserCard.vue';

const user = ref({
  id: 1,
  name: 'John Doe',
  email: 'john.doe@example.com',
  org: "平台运营方",
  avatar: 'http://localhost:8888/public/images/9f7b1307-95ae-4ab7-9368-ed80e0a49d0dda.jpg',
});
</script>

<style scoped>
</style>