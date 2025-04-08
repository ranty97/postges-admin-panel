import axios from 'axios'

// Функция для чтения запросов из localStorage
const readQueries = () => {
  try {
    const queries = localStorage.getItem('queries')
    return queries ? JSON.parse(queries) : []
  } catch (error) {
    console.error('Error reading queries:', error)
    return []
  }
}

// Функция для сохранения запросов в localStorage
const writeQueries = (queries) => {
  try {
    localStorage.setItem('queries', JSON.stringify(queries))
  } catch (error) {
    console.error('Error writing queries:', error)
  }
}

export const queriesApi = {
  // Получение списка всех запросов
  getQueries: async () => {
    return readQueries()
  },

  // Создание нового запроса
  createQuery: async (queryData) => {
    const queries = readQueries()
    const newQuery = {
      id: Date.now(), // Используем timestamp как id
      name: queryData.name,
      query: queryData.query,
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString()
    }
    
    queries.push(newQuery)
    writeQueries(queries)
    return newQuery
  },

  // Обновление существующего запроса
  updateQuery: async (id, queryData) => {
    const queries = readQueries()
    const queryIndex = queries.findIndex(q => q.id === id)
    
    if (queryIndex === -1) {
      throw new Error('Запрос не найден')
    }
    
    queries[queryIndex] = {
      ...queries[queryIndex],
      name: queryData.name,
      query: queryData.query,
      updated_at: new Date().toISOString()
    }
    
    writeQueries(queries)
    return queries[queryIndex]
  },

  // Удаление запроса
  deleteQuery: async (id) => {
    const queries = readQueries()
    const queryIndex = queries.findIndex(q => q.id === id)
    
    if (queryIndex === -1) {
      throw new Error('Запрос не найден')
    }
    
    queries.splice(queryIndex, 1)
    writeQueries(queries)
    return { message: 'Запрос успешно удален' }
  },

  // Выполнение запроса
  executeQuery: async (query) => {
    const response = await axios.post('/api/execute', { query })
    return response.data
  }
} 