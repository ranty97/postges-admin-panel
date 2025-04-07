<template>
  <el-container class="app-container">
    <el-header class="app-header">
      <div class="header-content">
        <div class="logo">
          <el-icon><DataLine /></el-icon>
          <h1>PostgreSQL Admin Panel</h1>
        </div>
      </div>
    </el-header>
    <el-container>
      <el-aside width="200px" class="app-sidebar">
        <el-menu
          :default-active="activePath"
          class="el-menu-vertical"
          router
        >
          <el-menu-item index="/tables">
            <el-icon><List /></el-icon>
            <span>Таблицы</span>
          </el-menu-item>
          <el-menu-item index="/queries">
            <el-icon><Document /></el-icon>
            <span>Запросы</span>
          </el-menu-item>
        </el-menu>
      </el-aside>
      <el-main class="app-main">
        <router-view v-slot="{ Component }">
          <keep-alive>
            <component :is="Component" />
          </keep-alive>
        </router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup>
import { DataLine, List, Document } from '@element-plus/icons-vue'
import { useRouter, useRoute } from 'vue-router'
import { computed } from 'vue'

const router = useRouter()
const route = useRoute()

const activePath = computed(() => {
  return route.path === '/' ? '/tables' : route.path
})
</script>

<style>
.app-container {
  min-height: 100vh;
}

.app-header {
  background: linear-gradient(135deg, #2c3e50 0%, #3498db 100%);
  color: white;
  padding: 0 20px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.header-content {
  height: 100%;
  display: flex;
  align-items: center;
}

.logo {
  display: flex;
  align-items: center;
  gap: 10px;
}

.logo .el-icon {
  font-size: 24px;
}

.logo h1 {
  margin: 0;
  font-size: 1.5rem;
  font-weight: 600;
}

.app-sidebar {
  background-color: #f5f7fa;
  border-right: 1px solid #e6e6e6;
}

.app-main {
  background-color: #f5f7fa;
  padding: 20px;
}

.el-menu-vertical {
  border-right: none;
}

.el-menu-item {
  display: flex;
  align-items: center;
  gap: 8px;
}

.el-menu-item .el-icon {
  font-size: 18px;
}

.el-menu-item.is-active {
  background-color: #ecf5ff;
  color: #409EFF;
}

.el-menu-item:hover {
  background-color: #f5f7fa;
}
</style> 