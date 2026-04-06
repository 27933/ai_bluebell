<template>
  <div class="forbidden-page">
    <div class="forbidden-content">
      <div class="forbidden-icon">🚫</div>
      <h1>403</h1>
      <p>无权访问此页面</p>
      <router-link to="/" class="back-btn">返回首页</router-link>
    </div>
  </div>
</template>

<script setup lang="ts"></script>

<style scoped>
.forbidden-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f8fafc;
}

.forbidden-content {
  text-align: center;
}

.forbidden-icon {
  font-size: 4rem;
  margin-bottom: 1rem;
}

h1 {
  font-size: 5rem;
  font-weight: 800;
  color: #1e293b;
  margin: 0 0 0.5rem;
}

p {
  color: #64748b;
  font-size: 1.1rem;
  margin-bottom: 1.5rem;
}

.back-btn {
  display: inline-block;
  padding: 0.6rem 1.5rem;
  background: #4f46e5;
  color: white;
  border-radius: 8px;
  text-decoration: none;
  font-weight: 600;
  transition: background 0.2s;
}

.back-btn:hover {
  background: #4338ca;
}
</style>
