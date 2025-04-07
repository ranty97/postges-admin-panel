import { createRouter, createWebHistory } from 'vue-router'
import TablesView from '../views/TablesView.vue'
import QueriesView from '../views/QueriesView.vue'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      redirect: '/tables'
    },
    {
      path: '/tables',
      name: 'tables',
      component: TablesView
    },
    {
      path: '/queries',
      name: 'queries',
      component: QueriesView
    }
  ]
})

export default router 