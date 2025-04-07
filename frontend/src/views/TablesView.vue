<template>
  <div>
    <el-row :gutter="20" class="mb-4">
      <el-col :span="24">
        <el-button type="primary" @click="showCreateTableDialog = true">
          Создать таблицу
        </el-button>
      </el-col>
    </el-row>

    <el-table :data="tables" style="width: 100%">
      <el-table-column prop="name" label="Название таблицы" />
      <el-table-column fixed="right" label="Действия" width="200">
        <template #default="scope">
          <el-button-group>
            <el-button type="primary" @click="viewTable(scope.row)">
              Просмотр
            </el-button>
            <el-button type="danger" @click="deleteTable(scope.row)">
              Удалить
            </el-button>
          </el-button-group>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="showCreateTableDialog" title="Создание таблицы">
      <el-form :model="newTable" label-width="120px">
        <el-form-item label="Название таблицы">
          <el-input v-model="newTable.name" />
        </el-form-item>
        <el-form-item label="Поля">
          <div v-for="(field, index) in newTable.fields" :key="index" class="field-item">
            <el-input v-model="field.name" placeholder="Имя поля" style="width: 200px" />
            <el-select v-model="field.type" placeholder="Тип" style="width: 150px">
              <el-option label="VARCHAR" value="VARCHAR" />
              <el-option label="INTEGER" value="INTEGER" />
              <el-option label="BOOLEAN" value="BOOLEAN" />
              <el-option label="DATE" value="DATE" />
            </el-select>
            <el-button type="danger" @click="removeField(index)">Удалить</el-button>
          </div>
          <el-button type="primary" @click="addField">Добавить поле</el-button>
        </el-form-item>
      </el-form>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="showCreateTableDialog = false">Отмена</el-button>
          <el-button type="primary" @click="createTable">Создать</el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { ElMessage, ElMessageBox } from 'element-plus'

const tables = ref([])
const showCreateTableDialog = ref(false)
const newTable = ref({
  name: '',
  fields: []
})

const fetchTables = async () => {
  try {
    const response = await axios.get('/api/tables')
    tables.value = response.data
  } catch (error) {
    ElMessage.error('Ошибка при получении списка таблиц')
  }
}

const createTable = async () => {
  try {
    const fields = newTable.value.fields.map(field => `${field.name} ${field.type}`).join(', ')
    const query = `CREATE TABLE ${newTable.value.name} (${fields})`
    await axios.post('/api/execute', { query })
    ElMessage.success('Таблица успешно создана')
    showCreateTableDialog.value = false
    fetchTables()
  } catch (error) {
    ElMessage.error('Ошибка при создании таблицы')
  }
}

const deleteTable = async (table) => {
  try {
    await ElMessageBox.confirm(
      `Вы уверены, что хотите удалить таблицу ${table.name}?`,
      'Предупреждение',
      {
        confirmButtonText: 'Да',
        cancelButtonText: 'Нет',
        type: 'warning',
      }
    )
    const query = `DROP TABLE ${table.name}`
    await axios.post('/api/execute', { query })
    ElMessage.success('Таблица успешно удалена')
    fetchTables()
  } catch (error) {
    if (error !== 'cancel') {
      ElMessage.error('Ошибка при удалении таблицы')
    }
  }
}

const viewTable = async (table) => {
  try {
    const query = `SELECT * FROM ${table.name}`
    const response = await axios.post('/api/execute', { query })
    // Здесь можно добавить логику для отображения данных таблицы
    console.log(response.data)
  } catch (error) {
    ElMessage.error('Ошибка при получении данных таблицы')
  }
}

const addField = () => {
  newTable.value.fields.push({ name: '', type: '' })
}

const removeField = (index) => {
  newTable.value.fields.splice(index, 1)
}

onMounted(() => {
  fetchTables()
})
</script>

<style scoped>
.field-item {
  margin-bottom: 10px;
  display: flex;
  gap: 10px;
  align-items: center;
}
</style> 