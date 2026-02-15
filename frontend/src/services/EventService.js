import axios from 'axios'

// This checks if Vercel has provided a URL. If not, it uses localhost.
const lAPIBaseURL = process.env.VUE_APP_API_URL || 'http://localhost:8080/api'

const lAxiosInstance = axios.create({
  baseURL: lAPIBaseURL,
  headers: {
    'Content-Type': 'application/json'
  }
})

const EventService = {
  signup: function(pData) {
    return lAxiosInstance.post('/auth/signup', pData)
  },

  login: function(pData) {
    return lAxiosInstance.post('/auth/login', pData)
  },

  logout: function(pToken) {
    return lAxiosInstance.post('/auth/logout', null, {
      headers: { 'Authorization': pToken }
    })
  },

  verifyToken: function(pToken) {
    return lAxiosInstance.get('/auth/verify', {
      headers: { 'Authorization': pToken }
    })
  },

  createTodo: function(pData, pToken) {
    return lAxiosInstance.post('/todos', pData, {
      headers: { 'Authorization': pToken }
    })
  },

  listTodos: function(pToken) {
    return lAxiosInstance.get('/todos', {
      headers: { 'Authorization': pToken }
    })
  },

  updateTodo: function(pTodoID, pData, pToken) {
    return lAxiosInstance.put('/todos/' + pTodoID, pData, {
      headers: { 'Authorization': pToken }
    })
  },

  deleteTodo: function(pTodoID, pToken) {
    return lAxiosInstance.delete('/todos/' + pTodoID, {
      headers: { 'Authorization': pToken }
    })
  }
}

export default EventService
