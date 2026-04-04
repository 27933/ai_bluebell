<script setup lang="ts">
import Navbar from './components/Navbar.vue'
</script>

<template>
  <el-container class="app-container">
    <Navbar />
    <el-main>
      <router-view />
    </el-main>
  </el-container>
</template>

<style scoped>
.app-container {
  height: 100vh;
}
</style>
