import axios from 'axios'

// Функция для чтения запросов из файла
const readQueries = () => {
  try {
    const queries = localStorage.getItem('queries')
    return queries ? JSON.parse(queries) : { queries: [] }
  } catch (error) {
    console.error('Error reading queries:', error)
    return { queries: [] }
  }
}

// Функция для записи запросов в файл
const writeQueries = (data) => {
  localStorage.setItem('queries', JSON.stringify(data))
}

export const queriesApi = {
  // Получение списка всех запросов
  getQueries: async () => {
    return readQueries()
  },

  // Создание нового запроса
  createQuery: async (queryData) => {
    const data = readQueries()
    const newQuery = {
      id: Date.now(), // Используем timestamp как id
      name: queryData.name,
      query: queryData.query,
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString()
    }
    
    data.queries.push(newQuery)
    writeQueries(data)
    return newQuery
  },

  // Обновление существующего запроса
  updateQuery: async (id, queryData) => {
    const data = readQueries()
    const queryIndex = data.queries.findIndex(q => q.id === id)
    
    if (queryIndex === -1) {
      throw new Error('Запрос не найден')
    }
    
    data.queries[queryIndex] = {
      ...data.queries[queryIndex],
      name: queryData.name,
      query: queryData.query,
      updated_at: new Date().toISOString()
    }
    
    writeQueries(data)
    return data.queries[queryIndex]
  },

  // Удаление запроса
  deleteQuery: async (id) => {
    const data = readQueries()
    const queryIndex = data.queries.findIndex(q => q.id === id)
    
    if (queryIndex === -1) {
      throw new Error('Запрос не найден')
    }
    
    data.queries.splice(queryIndex, 1)
    writeQueries(data)
    return { message: 'Запрос успешно удален' }
  },

  // Выполнение запроса
  executeQuery: async (query) => {
    const response = await axios.post('/api/execute', { query })
    return response.data
  }
} 