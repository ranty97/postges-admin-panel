<template>
  <div>
    <el-row :gutter="20" class="mb-4">
      <el-col :span="24">
        <el-button type="primary" @click="showCreateTableDialog = true">
          Создать таблицу
        </el-button>
      </el-col>
    </el-row>

    <el-table v-loading="loading" :data="tables" style="width: 100%">
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
              <el-option label="VARCHAR(255)" value="VARCHAR(255)" />
              <el-option label="INTEGER" value="INTEGER" />
              <el-option label="BOOLEAN" value="BOOLEAN" />
              <el-option label="DATE" value="DATE" />
              <el-option label="NUMERIC" value="NUMERIC" />
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

    <el-dialog v-model="showTableDataDialog" title="Данные таблицы" width="80%">
      <el-button type="primary" @click="addNewRow" class="mb-4">
        Добавить запись
      </el-button>
      <el-table :data="tableData" style="width: 100%">
        <el-table-column
          v-for="column in tableColumns"
          :key="column"
          :prop="column"
          :label="column"
        >
          <template #default="scope">
            <el-input
              v-if="scope.row.isEditing"
              v-model="scope.row[column]"
            />
            <span v-else @dblclick="startEditing(scope.row)">{{ scope.row[column] }}</span>
          </template>
        </el-table-column>
        <el-table-column label="Действия" width="200">
          <template #default="scope">
            <el-button-group>
              <el-button 
                v-if="scope.row.isEditing"
                type="success" 
                @click="saveRow(scope.row)"
              >
                Сохранить
              </el-button>
              <el-button 
                v-if="scope.row.isEditing"
                type="warning" 
                @click="cancelEditing(scope.row)"
              >
                Отмена
              </el-button>
              <el-button 
                v-if="!scope.row.isEditing"
                type="danger" 
                @click="deleteRow(scope.row)"
              >
                Удалить
              </el-button>
            </el-button-group>
          </template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import axios from 'axios'
import { ElMessage, ElMessageBox } from 'element-plus'

const tables = ref([])
const loading = ref(false)
const showCreateTableDialog = ref(false)
const showTableDataDialog = ref(false)
const tableData = ref([])
const tableColumns = ref([])
const newTable = ref({
  name: '',
  fields: []
})

const fetchTables = async () => {
  loading.value = true
  try {
    const response = await axios.get('/api/tables')
    console.log('Tables response:', response.data)
    tables.value = response.data.tables.map(name => ({ name }))
  } catch (error) {
    console.error('Error fetching tables:', error)
    ElMessage.error('Ошибка при получении списка таблиц: ' + (error.response?.data?.message || error.message))
  } finally {
    loading.value = false
  }
}

const createTable = async () => {
  if (!newTable.value.name) {
    ElMessage.warning('Введите название таблицы')
    return
  }
  if (newTable.value.fields.length === 0) {
    ElMessage.warning('Добавьте хотя бы одно поле')
    return
  }

  try {
    // Добавляем id как первичный ключ
    const fields = ['id SERIAL PRIMARY KEY', ...newTable.value.fields.map(field => `${field.name} ${field.type}`)].join(', ')
    const query = `CREATE TABLE ${newTable.value.name} (${fields})`
    console.log('Creating table with query:', query)
    await axios.post('/api/execute', { query })
    ElMessage.success('Таблица успешно создана')
    showCreateTableDialog.value = false
    newTable.value = { name: '', fields: [] }
    fetchTables()
  } catch (error) {
    console.error('Error creating table:', error)
    ElMessage.error('Ошибка при создании таблицы: ' + (error.response?.data?.message || error.message))
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
    console.log('Deleting table with query:', query)
    await axios.post('/api/execute', { query })
    ElMessage.success('Таблица успешно удалена')
    fetchTables()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Error deleting table:', error)
      ElMessage.error('Ошибка при удалении таблицы: ' + (error.response?.data?.message || error.message))
    }
  }
}

const decodeBase64 = (str) => {
  try {
    // Проверяем, является ли строка закодированной в base64
    if (str && typeof str === 'string') {
      // Проверяем, является ли строка закодированной в base64
      if (str.match(/^[A-Za-z0-9+/=]+$/)) {
        const decoded = atob(str)
        // Проверяем, является ли декодированное значение числом
        const num = Number(decoded)
        return isNaN(num) ? decoded : num
      }
      // Если строка не закодирована в base64, но выглядит как число, преобразуем в число
      const num = Number(str)
      return isNaN(num) ? str : num
    }
    return str
  } catch (e) {
    return str
  }
}

const viewTable = async (table) => {
  currentTable.value = table.name
  try {
    // Сначала получаем данные таблицы
    const query = `SELECT * FROM ${table.name}`
    const response = await axios.post('/api/execute', { query })
    
    if (response.data.result) {
      const data = JSON.parse(response.data.result)
      if (Array.isArray(data) && data.length > 0) {
        // Если есть данные, берем названия колонок из первой записи
        tableColumns.value = Object.keys(data[0])
        
        // Получаем информацию о типах полей
        const typesQuery = `SELECT column_name, data_type 
                          FROM information_schema.columns 
                          WHERE table_name = '${table.name}'`
        const typesResponse = await axios.post('/api/execute', { query: typesQuery })
        const typesData = JSON.parse(typesResponse.data.result)
        
        const columnTypes = {}
        typesData.forEach(col => {
          columnTypes[col.column_name] = col.data_type
        })
        
        // Преобразуем значения
        tableData.value = data.map(row => {
          const newRow = { ...row, isEditing: false }
          Object.keys(newRow).forEach(key => {
            // Декодируем все поля, которые могут быть закодированы в base64
            if (newRow[key] !== null) {
              newRow[key] = decodeBase64(newRow[key])
            }
          })
          return newRow
        })
      } else {
        // Если таблица пуста, получаем структуру из information_schema
        const columnsQuery = `SELECT column_name FROM information_schema.columns WHERE table_name = '${table.name}' ORDER BY ordinal_position`
        const columnsResponse = await axios.post('/api/execute', { query: columnsQuery })
        const columnsData = JSON.parse(columnsResponse.data.result)
        tableColumns.value = columnsData.map(col => decodeBase64(col.column_name))
        tableData.value = []
        
        // Показываем сообщение с предложением добавить запись
        ElMessageBox.confirm(
          'Таблица пуста. Хотите добавить первую запись?',
          'Пустая таблица',
          {
            confirmButtonText: 'Да',
            cancelButtonText: 'Нет',
            type: 'info',
          }
        ).then(() => {
          addNewRow()
        }).catch(() => {
          // Пользователь отказался добавлять запись
        })
      }
    } else {
      // Если таблица пуста, получаем структуру из information_schema
      const columnsQuery = `SELECT column_name FROM information_schema.columns WHERE table_name = '${table.name}' ORDER BY ordinal_position`
      const columnsResponse = await axios.post('/api/execute', { query: columnsQuery })
      const columnsData = JSON.parse(columnsResponse.data.result)
      tableColumns.value = columnsData.map(col => decodeBase64(col.column_name))
      tableData.value = []
      
      // Показываем сообщение с предложением добавить запись
      ElMessageBox.confirm(
        'Таблица пуста. Хотите добавить первую запись?',
        'Пустая таблица',
        {
          confirmButtonText: 'Да',
          cancelButtonText: 'Нет',
          type: 'info',
        }
      ).then(() => {
        addNewRow()
      }).catch(() => {
        // Пользователь отказался добавлять запись
      })
    }
    
    showTableDataDialog.value = true
  } catch (error) {
    console.error('Error viewing table:', error)
    ElMessage.error('Ошибка при получении данных таблицы: ' + (error.response?.data?.message || error.message))
  }
}

const addField = () => {
  newTable.value.fields.push({ name: '', type: '' })
}

const removeField = (index) => {
  newTable.value.fields.splice(index, 1)
}

const startEditing = (row) => {
  row.isEditing = true
}

const addNewRow = async () => {
  try {
    const tableName = currentTable.value
    
    // Получаем информацию о типах полей
    const typesQuery = `SELECT column_name, data_type, numeric_precision, numeric_scale 
                       FROM information_schema.columns 
                       WHERE table_name = '${tableName}'`
    const typesResponse = await axios.post('/api/execute', { query: typesQuery })
    const typesData = JSON.parse(typesResponse.data.result)
    
    // Создаем объект с пустыми значениями для каждого поля
    const newRow = {}
    typesData.forEach(col => {
      // Пропускаем поле id, так как оно будет автоматически сгенерировано
      if (col.column_name === 'id') return
      
      // Устанавливаем начальные значения в зависимости от типа поля
      switch (col.data_type) {
        case 'integer':
        case 'bigint':
        case 'smallint':
        case 'numeric':
          newRow[col.column_name] = 0
          break
        case 'boolean':
          newRow[col.column_name] = false
          break
        case 'date':
          newRow[col.column_name] = new Date().toISOString().split('T')[0]
          break
        default:
          newRow[col.column_name] = ''
      }
    })
    
    // Добавляем новую строку в таблицу
    tableData.value.push({ ...newRow, isEditing: true })
    
    // Находим индекс только что добавленной строки
    const newRowIndex = tableData.value.length - 1
    
    // Устанавливаем фокус на первое редактируемое поле
    setTimeout(() => {
      const firstInput = document.querySelector(`.el-table__body tr:nth-child(${newRowIndex + 1}) .el-input__inner`)
      if (firstInput) {
        firstInput.focus()
      }
    }, 0)
  } catch (error) {
    console.error('Error adding row:', error)
    ElMessage.error('Ошибка при добавлении записи: ' + (error.response?.data?.message || error.message))
  }
}

const saveRow = async (row) => {
  try {
    const tableName = currentTable.value
    const columns = tableColumns.value.filter(col => col !== 'id') // Исключаем id из списка колонок
    
    // Получаем информацию о типах полей
    const typesQuery = `SELECT column_name, data_type, numeric_precision, numeric_scale 
                       FROM information_schema.columns 
                       WHERE table_name = '${tableName}'`
    const typesResponse = await axios.post('/api/execute', { query: typesQuery })
    const typesData = JSON.parse(typesResponse.data.result)
    
    const columnTypes = {}
    typesData.forEach(col => {
      columnTypes[col.column_name] = {
        type: col.data_type,
        precision: col.numeric_precision,
        scale: col.numeric_scale
      }
    })
    
    // Проверяем, является ли строка новой (без id)
    if (!row.id) {
      // Формируем список колонок и значений для INSERT
      const insertColumns = []
      const insertValues = []
      
      columns.forEach(col => {
        const value = row[col]
        const colType = columnTypes[col]
        
        if (value === '' || value === null || value === undefined) {
          insertColumns.push(col)
          insertValues.push('NULL')
        } else if (colType) {
          insertColumns.push(col)
          if (colType.type === 'integer' || colType.type === 'bigint' || colType.type === 'smallint' || colType.type === 'numeric') {
            insertValues.push(value)
          } else if (colType.type === 'boolean') {
            insertValues.push(value ? 'true' : 'false')
          } else {
            insertValues.push(`'${value}'`)
          }
        } else {
          insertColumns.push(col)
          insertValues.push(`'${value}'`)
        }
      })
      
      // Вставляем новую запись
      const query = `INSERT INTO ${tableName} (${insertColumns.join(', ')}) VALUES (${insertValues.join(', ')}) RETURNING id`
      const response = await axios.post('/api/execute', { query })
      const result = JSON.parse(response.data.result)
      if (result && result.length > 0) {
        row.id = result[0].id
      }
    } else {
      // Проверяем, существует ли запись с таким id
      const checkQuery = `SELECT id FROM ${tableName} WHERE id = ${row.id}`
      const checkResponse = await axios.post('/api/execute', { query: checkQuery })
      const checkResult = JSON.parse(checkResponse.data.result)
      
      if (checkResult && checkResult.length > 0) {
        // Если запись существует, обновляем её
        const updateFields = []
        
        columns.forEach(col => {
          const value = row[col]
          const colType = columnTypes[col]
          
          if (value === '' || value === null || value === undefined) {
            updateFields.push(`${col} = NULL`)
          } else if (colType) {
            if (colType.type === 'integer' || colType.type === 'bigint' || colType.type === 'smallint' || colType.type === 'numeric') {
              updateFields.push(`${col} = ${value}`)
            } else if (colType.type === 'boolean') {
              updateFields.push(`${col} = ${value ? 'true' : 'false'}`)
            } else {
              updateFields.push(`${col} = '${value}'`)
            }
          } else {
            updateFields.push(`${col} = '${value}'`)
          }
        })
        
        const query = `UPDATE ${tableName} SET ${updateFields.join(', ')} WHERE id = ${row.id}`
        await axios.post('/api/execute', { query })
      } else {
        // Если записи нет, вставляем новую
        const insertColumns = []
        const insertValues = []
        
        columns.forEach(col => {
          const value = row[col]
          const colType = columnTypes[col]
          
          if (value === '' || value === null || value === undefined) {
            insertColumns.push(col)
            insertValues.push('NULL')
          } else if (colType) {
            insertColumns.push(col)
            if (colType.type === 'integer' || colType.type === 'bigint' || colType.type === 'smallint' || colType.type === 'numeric') {
              insertValues.push(value)
            } else if (colType.type === 'boolean') {
              insertValues.push(value ? 'true' : 'false')
            } else {
              insertValues.push(`'${value}'`)
            }
          } else {
            insertColumns.push(col)
            insertValues.push(`'${value}'`)
          }
        })
        
        const query = `INSERT INTO ${tableName} (id, ${insertColumns.join(', ')}) VALUES (${row.id}, ${insertValues.join(', ')})`
        await axios.post('/api/execute', { query })
      }
    }
    
    row.isEditing = false
    ElMessage.success('Запись успешно сохранена')
    viewTable({ name: tableName })
  } catch (error) {
    console.error('Error saving row:', error)
    ElMessage.error('Ошибка при сохранении записи: ' + (error.response?.data?.message || error.message))
  }
}

const deleteRow = async (row) => {
  try {
    await ElMessageBox.confirm(
      'Вы уверены, что хотите удалить эту запись?',
      'Предупреждение',
      {
        confirmButtonText: 'Да',
        cancelButtonText: 'Нет',
        type: 'warning',
      }
    )
    const tableName = currentTable.value
    
    // Получаем первичный ключ таблицы
    const pkQuery = `SELECT 
      kcu.column_name
    FROM information_schema.key_column_usage kcu
    JOIN information_schema.table_constraints tc
      ON kcu.constraint_name = tc.constraint_name
    WHERE tc.table_name = '${tableName}'
      AND tc.constraint_type = 'PRIMARY KEY'`
    
    const pkResponse = await axios.post('/api/execute', { query: pkQuery })
    const pkData = JSON.parse(pkResponse.data.result)
    const primaryKey = pkData[0]?.column_name
    
    if (!primaryKey) {
      throw new Error('Не найден первичный ключ таблицы')
    }
    
    const query = `DELETE FROM ${tableName} WHERE ${primaryKey} = ${row[primaryKey]}`
    await axios.post('/api/execute', { query })
    ElMessage.success('Запись успешно удалена')
    viewTable({ name: tableName })
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Error deleting row:', error)
      ElMessage.error('Ошибка при удалении записи: ' + (error.response?.data?.message || error.message))
    }
  }
}

const cancelEditing = (row) => {
  row.isEditing = false
  // Если это новая запись, удаляем её
  if (!row.id) {
    const index = tableData.value.findIndex(r => r === row)
    if (index !== -1) {
      tableData.value.splice(index, 1)
    }
  } else {
    // Если это существующая запись, перезагружаем данные
    viewTable({ name: currentTable.value })
  }
}

const currentTable = ref('')

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