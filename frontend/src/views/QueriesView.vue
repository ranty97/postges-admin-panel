<template>
  <div>
    <el-row :gutter="20" class="mb-4">
      <el-col :span="24">
        <el-input
          v-model="query"
          type="textarea"
          :rows="4"
          placeholder="Введите SQL запрос"
        />
      </el-col>
    </el-row>

    <el-row :gutter="20" class="mb-4">
      <el-col :span="24">
        <el-button type="primary" @click="executeQuery" :loading="loading">
          Выполнить запрос
        </el-button>
      </el-col>
    </el-row>

    <el-table 
      v-if="queryResult.length > 0" 
      :data="queryResult" 
      style="width: 100%"
      v-loading="loading"
    >
      <el-table-column
        v-for="column in Object.keys(queryResult[0])"
        :key="column"
        :prop="column"
        :label="column"
      />
    </el-table>

    <div v-if="error" class="error-message">
      {{ error }}
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import axios from 'axios'
import { ElMessage } from 'element-plus'

const query = ref('')
const queryResult = ref([])
const error = ref('')
const loading = ref(false)

const executeQuery = async () => {
  if (!query.value.trim()) {
    ElMessage.warning('Введите SQL запрос')
    return
  }

  loading.value = true
  error.value = ''
  queryResult.value = []

  try {
    console.log('Executing query:', query.value)
    const response = await axios.post('/api/execute', { query: query.value })
    console.log('Query response:', response.data)
    
    if (response.data.result) {
      const data = JSON.parse(response.data.result)
      if (Array.isArray(data)) {
        queryResult.value = data
      } else {
        ElMessage.success('Запрос выполнен успешно')
      }
    } else {
      ElMessage.success('Запрос выполнен успешно')
    }
  } catch (err) {
    console.error('Error executing query:', err)
    error.value = err.response?.data?.message || 'Произошла ошибка при выполнении запроса'
    ElMessage.error(error.value)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.error-message {
  color: red;
  margin-top: 10px;
  padding: 10px;
  background-color: #fef0f0;
  border-radius: 4px;
}
</style> 