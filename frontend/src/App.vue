<template>
  <el-container>
    <el-header>
      <h1>PostgreSQL Admin Panel</h1>
    </el-header>
    <el-container>
      <el-aside width="200px">
        <el-menu
          default-active="1"
          class="el-menu-vertical"
          @select="handleSelect"
        >
          <el-menu-item index="1">
            <el-icon><List /></el-icon>
            <span>Таблицы</span>
          </el-menu-item>
          <el-menu-item index="2">
            <el-icon><Edit /></el-icon>
            <span>Запросы</span>
          </el-menu-item>
        </el-menu>
      </el-aside>
      <el-main>
        <router-view></router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { List, Edit } from '@element-plus/icons-vue'

const router = useRouter()

const handleSelect = (key) => {
  if (key === '1') {
    router.push('/tables')
  } else if (key === '2') {
    router.push('/queries')
  }
}
</script>

<style>
.el-header {
  background-color: #409EFF;
  color: white;
  line-height: 60px;
  text-align: center;
}

.el-aside {
  background-color: #f5f7fa;
}

.el-main {
  padding: 20px;
}
</style> 