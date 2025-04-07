<template>
  <div class="queries-container">
    <div class="queries-header">
      <h2>Управление запросами</h2>
      <el-button type="primary" @click="dialogVisible = true">
        <el-icon><Plus /></el-icon>
        Создать запрос
      </el-button>
    </div>

    <el-table :data="queries" v-loading="loading" style="width: 100%">
      <el-table-column prop="name" label="Название" />
      <el-table-column prop="query" label="Запрос" />
      <el-table-column prop="created_at" label="Создан">
        <template #default="{ row }">
          {{ new Date(row.created_at).toLocaleString('ru-RU') }}
        </template>
      </el-table-column>
      <el-table-column prop="updated_at" label="Обновлен">
        <template #default="{ row }">
          {{ new Date(row.updated_at).toLocaleString('ru-RU') }}
        </template>
      </el-table-column>
      <el-table-column label="Действия" width="250">
        <template #default="{ row }">
          <el-button-group>
            <el-button type="primary" @click="handleExecute(row.query)">
              <el-icon><VideoPlay /></el-icon>
              Выполнить
            </el-button>
            <el-button type="warning" @click="startEdit(row)">
              <el-icon><Edit /></el-icon>
              Редактировать
            </el-button>
            <el-button type="danger" @click="handleDelete(row.id)">
              <el-icon><Delete /></el-icon>
              Удалить
            </el-button>
          </el-button-group>
        </template>
      </el-table-column>
    </el-table>

    <!-- Диалог создания нового запроса -->
    <el-dialog v-model="dialogVisible" title="Создание нового запроса" width="50%">
      <el-form>
        <el-form-item label="Название">
          <el-input v-model="newQueryName" />
        </el-form-item>
        <el-form-item label="Запрос">
          <el-input v-model="newQueryText" type="textarea" rows="5" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="dialogVisible = false">Отмена</el-button>
          <el-button type="primary" @click="handleCreate">Создать</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- Диалог редактирования запроса -->
    <el-dialog v-model="editDialogVisible" title="Редактирование запроса" width="50%">
      <el-form>
        <el-form-item label="Название">
          <el-input v-model="editingQuery.name" />
        </el-form-item>
        <el-form-item label="Запрос">
          <el-input v-model="editingQuery.query" type="textarea" rows="5" />
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="editDialogVisible = false">Отмена</el-button>
          <el-button type="primary" @click="handleEdit">Сохранить</el-button>
        </span>
      </template>
    </el-dialog>

    <!-- Диалог результатов -->
    <el-dialog v-model="showResultsDialog" title="Результаты запроса" width="80%">
      <div class="results-actions">
        <el-button type="success" @click="exportToExcel">
          <el-icon><Download /></el-icon>
          Экспорт в Excel
        </el-button>
      </div>
      <el-table :data="queryResults" style="width: 100%">
        <el-table-column v-for="col in resultColumns" :key="col" :prop="col" :label="col" />
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { queriesApi } from '../api/queries'
import { onBeforeRouteUpdate } from 'vue-router'
import { Plus, VideoPlay, Edit, Delete, Download } from '@element-plus/icons-vue'
import * as XLSX from 'xlsx'

const queries = ref([])
const loading = ref(false)
const dialogVisible = ref(false)
const editDialogVisible = ref(false)
const showResultsDialog = ref(false)
const newQueryName = ref('')
const newQueryText = ref('')
const editingQuery = ref(null)
const queryResults = ref([])
const resultColumns = ref([])

const startEdit = (query) => {
  editingQuery.value = { ...query }
  editDialogVisible.value = true
}

const fetchQueries = async () => {
  try {
    loading.value = true
    const response = await queriesApi.getQueries()
    queries.value = response
  } catch (error) {
    console.error('Ошибка при загрузке запросов:', error)
    ElMessage.error('Не удалось загрузить запросы')
  } finally {
    loading.value = false
  }
}

const handleCreate = async () => {
  try {
    const newQuery = await queriesApi.createQuery({
      name: newQueryName.value,
      query: newQueryText.value
    })
    queries.value.push(newQuery)
    ElMessage.success('Запрос успешно создан')
    dialogVisible.value = false
    newQueryName.value = ''
    newQueryText.value = ''
  } catch (error) {
    console.error('Ошибка при создании запроса:', error)
    ElMessage.error('Не удалось создать запрос')
  }
}

const handleEdit = async () => {
  try {
    const updatedQuery = await queriesApi.updateQuery(editingQuery.value.id, {
      name: editingQuery.value.name,
      query: editingQuery.value.query
    })
    const index = queries.value.findIndex(q => q.id === updatedQuery.id)
    if (index !== -1) {
      queries.value[index] = updatedQuery
    }
    ElMessage.success('Запрос успешно обновлен')
    editDialogVisible.value = false
  } catch (error) {
    console.error('Ошибка при обновлении запроса:', error)
    ElMessage.error('Не удалось обновить запрос')
  }
}

const handleDelete = async (id) => {
  try {
    await queriesApi.deleteQuery(id)
    queries.value = queries.value.filter(q => q.id !== id)
    ElMessage.success('Запрос успешно удален')
  } catch (error) {
    console.error('Ошибка при удалении запроса:', error)
    ElMessage.error('Не удалось удалить запрос')
  }
}

const handleExecute = async (query) => {
  try {
    loading.value = true
    console.log('Выполнение запроса:', query)
    const result = await queriesApi.executeQuery(query)
    console.log('Результат запроса:', result)
    
    if (result && result.result) {
      try {
        const parsedResult = JSON.parse(result.result)
        console.log('Распарсенный результат:', parsedResult)
        
        if (Array.isArray(parsedResult) && parsedResult.length > 0) {
          resultColumns.value = Object.keys(parsedResult[0])
          queryResults.value = parsedResult
          showResultsDialog.value = true
        } else {
          ElMessage.success('Запрос выполнен успешно, но не вернул результатов')
        }
      } catch (parseError) {
        console.error('Ошибка при парсинге результата:', parseError)
        ElMessage.error('Ошибка при обработке результатов запроса')
      }
    } else {
      ElMessage.success('Запрос выполнен успешно, но не вернул результатов')
    }
  } catch (error) {
    console.error('Ошибка при выполнении запроса:', error)
    ElMessage.error('Не удалось выполнить запрос: ' + (error.response?.data?.message || error.message))
  } finally {
    loading.value = false
  }
}

const exportToExcel = () => {
  try {
    // Создаем новую книгу Excel
    const wb = XLSX.utils.book_new()
    
    // Преобразуем данные в формат, подходящий для Excel
    const ws = XLSX.utils.json_to_sheet(queryResults.value)
    
    // Добавляем лист в книгу
    XLSX.utils.book_append_sheet(wb, ws, "Результаты запроса")
    
    // Генерируем имя файла с текущей датой и временем
    const fileName = `query_results_${new Date().toISOString().replace(/[:.]/g, '-')}.xlsx`
    
    // Сохраняем файл
    XLSX.writeFile(wb, fileName)
    
    ElMessage.success('Результаты успешно экспортированы в Excel')
  } catch (error) {
    console.error('Ошибка при экспорте в Excel:', error)
    ElMessage.error('Не удалось экспортировать результаты в Excel')
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
.queries-container {
  padding: 20px;
}

.queries-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.dialog-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
}

.results-actions {
  margin-bottom: 20px;
  display: flex;
  justify-content: flex-end;
}
</style> 