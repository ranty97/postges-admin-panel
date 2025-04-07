<template>
  <div>
    <div class="banner">
      <h1>Управление запросами</h1>
      <p>Создавайте, сохраняйте и выполняйте SQL-запросы</p>
    </div>

    <el-row :gutter="20" class="mb-4">
      <el-col :span="24">
        <el-button type="primary" @click="showNewQueryDialog = true">
          Создать запрос
        </el-button>
      </el-col>
    </el-row>

    <el-table v-loading="loading" :data="queries" style="width: 100%">
      <el-table-column prop="name" label="Название запроса" />
      <el-table-column prop="query" label="SQL-запрос" show-overflow-tooltip />
      <el-table-column fixed="right" label="Действия" width="300">
        <template #default="scope">
          <el-button-group>
            <el-button type="primary" @click="executeQuery(scope.row)">
              Выполнить
            </el-button>
            <el-button type="warning" @click="editQuery(scope.row)">
              Редактировать
            </el-button>
            <el-button type="danger" @click="deleteQuery(scope.row)">
              Удалить
            </el-button>
          </el-button-group>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog
      v-model="showNewQueryDialog"
      :title="editingQuery ? 'Редактирование запроса' : 'Создание нового запроса'"
      width="60%"
      :close-on-click-modal="false"
    >
      <el-form :model="currentQuery" label-position="top">
        <el-form-item label="Название запроса" required>
          <el-input v-model="currentQuery.name" placeholder="Введите название запроса" />
        </el-form-item>
        <el-form-item label="SQL-запрос" required>
          <el-input
            v-model="currentQuery.query"
            type="textarea"
            :rows="6"
            placeholder="Введите SQL-запрос"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="cancelQuery">Отмена</el-button>
          <el-button type="primary" @click="saveQuery">
            {{ editingQuery ? 'Сохранить' : 'Создать' }}
          </el-button>
        </span>
      </template>
    </el-dialog>

    <el-dialog
      v-model="showResultsDialog"
      title="Результаты запроса"
      width="80%"
    >
      <el-table :data="queryResults" style="width: 100%">
        <el-table-column
          v-for="column in resultColumns"
          :key="column"
          :prop="column"
          :label="column"
        />
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { queriesApi } from '../routes/queries'
import { onBeforeRouteUpdate } from 'vue-router'

const queries = ref([])
const loading = ref(false)
const showNewQueryDialog = ref(false)
const showResultsDialog = ref(false)
const editingQuery = ref(null)
const currentQuery = ref({
  name: '',
  query: ''
})
const queryResults = ref([])
const resultColumns = ref([])

const fetchQueries = async () => {
  loading.value = true
  try {
    queries.value = await queriesApi.getQueries()
  } catch (error) {
    console.error('Error fetching queries:', error)
    ElMessage.error('Ошибка при получении списка запросов: ' + (error.response?.data?.message || error.message))
  } finally {
    loading.value = false
  }
}

const saveQuery = async () => {
  if (!currentQuery.value.name || !currentQuery.value.query) {
    ElMessage.warning('Заполните все поля')
    return
  }

  try {
    if (editingQuery.value) {
      await queriesApi.updateQuery(editingQuery.value.id, currentQuery.value)
      ElMessage.success('Запрос успешно обновлен')
    } else {
      await queriesApi.createQuery(currentQuery.value)
      ElMessage.success('Запрос успешно создан')
    }
    showNewQueryDialog.value = false
    fetchQueries()
  } catch (error) {
    console.error('Error saving query:', error)
    ElMessage.error('Ошибка при сохранении запроса: ' + (error.response?.data?.message || error.message))
  }
}

const deleteQuery = async (query) => {
  try {
    await ElMessageBox.confirm(
      `Вы уверены, что хотите удалить запрос "${query.name}"?`,
      'Предупреждение',
      {
        confirmButtonText: 'Да',
        cancelButtonText: 'Нет',
        type: 'warning',
      }
    )
    await queriesApi.deleteQuery(query.id)
    ElMessage.success('Запрос успешно удален')
    fetchQueries()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Error deleting query:', error)
      ElMessage.error('Ошибка при удалении запроса: ' + (error.response?.data?.message || error.message))
    }
  }
}

const executeQuery = async (query) => {
  try {
    const response = await queriesApi.executeQuery(query.query)
    if (response.result) {
      const data = JSON.parse(response.result)
      if (Array.isArray(data) && data.length > 0) {
        resultColumns.value = Object.keys(data[0])
        queryResults.value = data
        showResultsDialog.value = true
      } else {
        ElMessage.success('Запрос выполнен успешно')
      }
    }
  } catch (error) {
    console.error('Error executing query:', error)
    ElMessage.error('Ошибка при выполнении запроса: ' + (error.response?.data?.message || error.message))
  }
}

const editQuery = (query) => {
  editingQuery.value = query
  currentQuery.value = { ...query }
  showNewQueryDialog.value = true
}

const cancelQuery = () => {
  showNewQueryDialog.value = false
  editingQuery.value = null
  currentQuery.value = {
    name: '',
    query: ''
  }
}

onMounted(() => {
  console.log('QueriesView mounted')
  fetchQueries()
})

onBeforeRouteUpdate((to, from, next) => {
  console.log('QueriesView route update:', { to: to.path, from: from.path })
  fetchQueries()
  next()
})
</script>

<style scoped>
.banner {
  background: linear-gradient(135deg, #409EFF 0%, #36cfc9 100%);
  color: white;
  padding: 2rem;
  margin-bottom: 2rem;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.banner h1 {
  margin: 0;
  font-size: 2rem;
  font-weight: 600;
}

.banner p {
  margin: 0.5rem 0 0;
  font-size: 1.1rem;
  opacity: 0.9;
}
</style> 