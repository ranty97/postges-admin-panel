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
        <el-button type="primary" @click="executeQuery">
          Выполнить запрос
        </el-button>
      </el-col>
    </el-row>

    <el-table v-if="queryResult.length > 0" :data="queryResult" style="width: 100%">
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

const executeQuery = async () => {
  try {
    const response = await axios.post('/api/execute', { query: query.value })
    if (Array.isArray(response.data)) {
      queryResult.value = response.data
      error.value = ''
    } else {
      ElMessage.success('Запрос выполнен успешно')
      queryResult.value = []
      error.value = ''
    }
  } catch (err) {
    error.value = err.response?.data?.message || 'Произошла ошибка при выполнении запроса'
    queryResult.value = []
  }
}
</script>

<style scoped>
.error-message {
  color: red;
  margin-top: 10px;
}
</style> 