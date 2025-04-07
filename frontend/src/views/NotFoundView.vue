<template>
  <div class="not-found-container">
    <el-result
      icon="error"
      title="404"
      sub-title="Страница не найдена"
    >
      <template #extra>
        <el-button type="primary" @click="$router.push('/')">
          Вернуться на главную
        </el-button>
      </template>
    </el-result>
  </div>
</template>

<script setup>
import { ElResult, ElButton } from 'element-plus'
</script>

<style scoped>
.not-found-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 80vh;
}
</style> 