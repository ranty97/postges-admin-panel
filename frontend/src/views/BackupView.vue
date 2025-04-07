<template>
  <div class="backup-container">
    <div class="banner">
      <h1>Управление бэкапами</h1>
      <p>Создавайте, восстанавливайте и удаляйте резервные копии базы данных</p>
    </div>

    <div class="backup-actions">
      <el-button type="primary" @click="handleCreateBackup">
        <el-icon><Download /></el-icon>
        Создать бэкап
      </el-button>
      <el-button type="success" @click="handleRestoreBackup">
        <el-icon><Upload /></el-icon>
        Восстановить из бэкапа
      </el-button>
    </div>

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
      <el-table-column label="Действия" width="200">
        <template #default="{ row }">
          <el-button-group>
            <el-button type="primary" @click="handleDownloadBackup(row.filename)">
              <el-icon><Download /></el-icon>
              Скачать
            </el-button>
            <el-button type="danger" @click="handleDeleteBackup(row.filename)">
              <el-icon><Delete /></el-icon>
              Удалить
            </el-button>
          </el-button-group>
        </template>
      </el-table-column>
    </el-table>

    <!-- Диалог восстановления из бэкапа -->
    <el-dialog v-model="restoreDialogVisible" title="Восстановление из бэкапа" width="30%">
      <el-form>
        <el-form-item label="Выберите файл бэкапа">
          <el-upload
            class="upload-demo"
            drag
            action="/api/backup/restore"
            :on-success="handleRestoreSuccess"
            :on-error="handleRestoreError"
            :before-upload="beforeRestoreUpload"
          >
            <el-icon class="el-icon--upload"><upload-filled /></el-icon>
            <div class="el-upload__text">
              Перетащите файл бэкапа сюда или <em>нажмите для загрузки</em>
            </div>
          </el-upload>
        </el-form-item>
      </el-form>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Download, Upload, Delete, UploadFilled } from '@element-plus/icons-vue'

const backups = ref([])
const loading = ref(false)
const restoreDialogVisible = ref(false)

const fetchBackups = async () => {
  try {
    loading.value = true
    const response = await fetch('/api/backup/list')
    if (!response.ok) throw new Error('Ошибка при загрузке списка бэкапов')
    backups.value = await response.json()
  } catch (error) {
    console.error('Ошибка при загрузке бэкапов:', error)
    ElMessage.error('Не удалось загрузить список бэкапов')
  } finally {
    loading.value = false
  }
}

const handleCreateBackup = async () => {
  try {
    loading.value = true
    const response = await fetch('/api/backup/create', { method: 'POST' })
    if (!response.ok) throw new Error('Ошибка при создании бэкапа')
    ElMessage.success('Бэкап успешно создан')
    await fetchBackups()
  } catch (error) {
    console.error('Ошибка при создании бэкапа:', error)
    ElMessage.error('Не удалось создать бэкап')
  } finally {
    loading.value = false
  }
}

const handleRestoreBackup = () => {
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
    
    ElMessage.success('Бэкап успешно удален')
    await fetchBackups()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Ошибка при удалении бэкапа:', error)
      ElMessage.error('Не удалось удалить бэкап')
    }
  }
}

const handleRestoreSuccess = () => {
  ElMessage.success('База данных успешно восстановлена из бэкапа')
  restoreDialogVisible.value = false
}

const handleRestoreError = () => {
  ElMessage.error('Ошибка при восстановлении базы данных')
}

const beforeRestoreUpload = (file) => {
  if (!file.name.endsWith('.sql')) {
    ElMessage.error('Файл должен иметь расширение .sql')
    return false
  }
  return true
}

const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
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

.upload-demo {
  width: 100%;
}
</style> 