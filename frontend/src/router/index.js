import { createRouter, createWebHistory } from 'vue-router'
import TablesView from '../views/TablesView.vue'
import QueriesView from '../views/QueriesView.vue'
import NotFoundView from '../views/NotFoundView.vue'

const routes = [
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
      icon: 'List'
    }
  },
  {
    path: '/queries',
    name: 'queries',
    component: QueriesView,
    meta: {
      title: 'Запросы',
      icon: 'Document'
    }
  },
  {
    path: '/:pathMatch(.*)*',
    name: 'not-found',
    component: NotFoundView,
    meta: {
      title: 'Страница не найдена',
      icon: 'Warning'
    }
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
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