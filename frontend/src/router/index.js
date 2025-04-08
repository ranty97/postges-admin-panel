import { createRouter, createWebHistory } from 'vue-router'
import TablesView from '../views/TablesView.vue'
import QueriesView from '../views/QueriesView.vue'
import { List, Document, Warning, Upload } from '@element-plus/icons-vue'

// Импортируем BackupView напрямую для отладки
import BackupView from '../views/BackupView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/tables'
    },
    {
      path: '/tables',
      name: 'tables',
      component: TablesView,
      meta: {
        title: 'Таблицы',
        icon: List
      }
    },
    {
      path: '/queries',
      name: 'queries',
      component: QueriesView,
      meta: {
        title: 'Запросы',
        icon: Document
      }
    },
    {
      path: '/backup',
      name: 'backup',
      component: BackupView,
      meta: {
        title: 'Бэкапы',
        icon: Upload
      }
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'not-found',
      component: () => import('../views/NotFoundView.vue'),
      meta: {
        title: 'Страница не найдена',
        icon: Warning
      }
    }
  ]
})

// Добавим логирование для отладки
router.beforeEach((to, from, next) => {
  console.log('Навигация:', { from: from.path, to: to.path, params: to.params, query: to.query })
  next()
})

// Глобальный навигационный хук
router.beforeEach((to, from, next) => {
  console.log('Навигация:', {
    from: from.path,
    to: to.path,
    params: to.params,
    query: to.query
  })
  
  // Установка заголовка страницы
  document.title = to.meta.title ? `${to.meta.title} | PostgreSQL Admin Panel` : 'PostgreSQL Admin Panel'
  
  // Здесь можно добавить проверку авторизации
  if (to.meta.requiresAuth) {
    // Проверка авторизации
    next()
  } else {
    next()
  }
})

// Обработка ошибок навигации
router.onError((error) => {
  console.error('Ошибка навигации:', error)
})

export default router 