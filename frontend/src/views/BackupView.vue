<template>
  <div class="backup-container">
    <div class="banner">
      <h1>Управление бэкапами</h1>
      <p>Создавайте, восстанавливайте и удаляйте резервные копии базы данных</p>
    </div>

    <div class="backup-actions">
      <el-button type="primary" @click="handleCreateBackup" :loading="creatingBackup">
        <el-icon><Download /></el-icon>
        Создать бэкап
      </el-button>
    </div>

    <el-card class="backup-list">
      <template #header>
        <div class="card-header">
          <span>Список бэкапов</span>
          <el-button type="text" @click="refreshBackups">
            <el-icon><Refresh /></el-icon>
            Обновить
          </el-button>
        </div>
      </template>

      <el-table :data="backups" v-loading="loading" style="width: 100%">
        <el-table-column prop="filename" label="Имя файла" />
        <el-table-column prop="created_at" label="Создан">
          <template #default="{ row }">
            {{ new Date(row.created_at).toLocaleString('ru-RU') }}
          </template>
        </el-table-column>
        <el-table-column prop="size" label="Размер">
          <template #default="{ row }">
            {{ formatFileSize(row.size) }}
          </template>
        </el-table-column>
        <el-table-column label="Действия" width="300">
          <template #default="{ row }">
            <el-button-group>
              <el-button type="primary" @click="handleDownloadBackup(row.filename)">
                <el-icon><Download /></el-icon>
                Скачать
              </el-button>
              <el-button type="success" @click="handleRestoreBackup(row.filename)">
                <el-icon><Upload /></el-icon>
                Восстановить
              </el-button>
              <el-button type="danger" @click="handleDeleteBackup(row.filename)">
                <el-icon><Delete /></el-icon>
                Удалить
              </el-button>
            </el-button-group>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <!-- Диалог подтверждения восстановления -->
    <el-dialog v-model="restoreDialogVisible" title="Восстановление из бэкапа" width="30%">
      <p>Вы уверены, что хотите восстановить базу данных из бэкапа "{{ selectedBackup }}"?</p>
      <p class="warning-text">Внимание: Это действие перезапишет текущую базу данных!</p>
      <template #footer>
        <span class="dialog-footer">
          <el-button @click="restoreDialogVisible = false">Отмена</el-button>
          <el-button type="success" @click="confirmRestore" :loading="restoring">
            Восстановить
          </el-button>
        </span>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Download, Upload, Delete, UploadFilled, Refresh } from '@element-plus/icons-vue'

console.log('BackupView mounted')

const backups = ref([])
const loading = ref(false)
const creatingBackup = ref(false)
const restoreDialogVisible = ref(false)
const selectedBackup = ref('')
const restoring = ref(false)

const fetchBackups = async () => {
  try {
    loading.value = true
    console.log('Загрузка списка бэкапов...')
    const response = await fetch('/api/backup/list')
    console.log('Ответ сервера:', response)
    if (!response.ok) throw new Error('Ошибка при загрузке списка бэкапов')
    const data = await response.json()
    console.log('Данные бэкапов:', data)
    backups.value = data.backups || []
    console.log('Обновленный список бэкапов:', backups.value)
  } catch (error) {
    console.error('Ошибка при загрузке бэкапов:', error)
    ElMessage.error('Не удалось загрузить список бэкапов')
  } finally {
    loading.value = false
  }
}

const handleCreateBackup = async () => {
  try {
    creatingBackup.value = true
    const response = await fetch('/api/backup/create', { method: 'POST' })
    if (!response.ok) throw new Error('Ошибка при создании бэкапа')
    const data = await response.json()
    ElMessage.success(data.message || 'Бэкап успешно создан')
    await fetchBackups()
  } catch (error) {
    console.error('Ошибка при создании бэкапа:', error)
    ElMessage.error('Не удалось создать бэкап')
  } finally {
    creatingBackup.value = false
  }
}

const handleRestoreBackup = (filename) => {
  selectedBackup.value = filename
  restoreDialogVisible.value = true
}

const handleDownloadBackup = async (filename) => {
  try {
    const response = await fetch(`/api/backup/download/${filename}`)
    if (!response.ok) throw new Error('Ошибка при скачивании бэкапа')
    
    const blob = await response.blob()
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = filename
    document.body.appendChild(a)
    a.click()
    window.URL.revokeObjectURL(url)
    document.body.removeChild(a)
  } catch (error) {
    console.error('Ошибка при скачивании бэкапа:', error)
    ElMessage.error('Не удалось скачать бэкап')
  }
}

const handleDeleteBackup = async (filename) => {
  try {
    await ElMessageBox.confirm(
      'Вы уверены, что хотите удалить этот бэкап?',
      'Подтверждение удаления',
      {
        confirmButtonText: 'Да',
        cancelButtonText: 'Нет',
        type: 'warning'
      }
    )
    
    const response = await fetch(`/api/backup/delete/${filename}`, { method: 'DELETE' })
    if (!response.ok) throw new Error('Ошибка при удалении бэкапа')
    const data = await response.json()
    ElMessage.success(data.message || 'Бэкап успешно удален')
    await fetchBackups()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Ошибка при удалении бэкапа:', error)
      ElMessage.error('Не удалось удалить бэкап')
    }
  }
}

const confirmRestore = async () => {
  try {
    restoring.value = true
    const response = await fetch(`/api/backup/restore/${selectedBackup.value}`, { method: 'POST' })
    if (!response.ok) throw new Error('Ошибка при восстановлении бэкапа')
    const data = await response.json()
    ElMessage.success(data.message || 'База данных успешно восстановлена из бэкапа')
    restoreDialogVisible.value = false
  } catch (error) {
    console.error('Ошибка при восстановлении бэкапа:', error)
    ElMessage.error('Не удалось восстановить базу данных')
  } finally {
    restoring.value = false
  }
}

const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

const refreshBackups = () => {
  fetchBackups()
}

onMounted(() => {
  fetchBackups()
})
</script>

<style scoped>
.backup-container {
  padding: 20px;
}

.backup-actions {
  margin-bottom: 20px;
  display: flex;
  gap: 10px;
}

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

.backup-list {
  margin-top: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.upload-demo {
  width: 100%;
}

.el-upload__tip {
  margin-top: 8px;
  color: #909399;
  font-size: 12px;
}

.warning-text {
  color: #e6a23c;
  margin-top: 10px;
  font-weight: 500;
}
</style> 