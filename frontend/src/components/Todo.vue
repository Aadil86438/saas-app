<template>
  <v-container fluid class="pa-4">
    <v-app-bar :dark="darkMode" color="primary" elevation="2" rounded class="mb-6">
      <v-toolbar-title class="text-h5 font-weight-light">
        My Todos
      </v-toolbar-title>
      <v-spacer></v-spacer>
      <v-chip color="white" text-color="primary" class="mr-4">
        <v-icon left small>mdi-account</v-icon>
        {{ currentUser.username }}
      </v-chip>
      <v-btn icon @click="toggleTheme">
        <v-icon>{{ darkMode ? 'mdi-weather-sunny' : 'mdi-weather-night' }}</v-icon>
      </v-btn>
      <v-btn text @click="handleLogout">
        <v-icon left>mdi-logout</v-icon>
        Logout
      </v-btn>
    </v-app-bar>

    <v-row>
      <v-col cols="12" md="8" offset-md="2">
        <v-card class="elevation-4 rounded-lg mb-6" :dark="darkMode">
          <v-card-title class="text-h6 pa-6">
            Create New Todo
          </v-card-title>
          <v-card-text class="pa-6 pt-0">
            <v-form ref="todoForm" v-model="todoFormValid">
              <v-text-field
                v-model="newTodoTitle"
                :rules="titleRules"
                label="Title"
                outlined
                rounded
                class="mb-4"
                @keyup.enter="handleCreateTodo"
              ></v-text-field>
              <v-textarea
                v-model="newTodoContent"
                label="Content"
                outlined
                rounded
                rows="3"
                class="mb-4"
              ></v-textarea>
              <v-btn
                :disabled="!todoFormValid || creating"
                :loading="creating"
                color="primary"
                large
                rounded
                block
                @click="handleCreateTodo"
              >
                Add Todo
              </v-btn>
            </v-form>
          </v-card-text>
        </v-card>

        <v-card
          v-for="todo in todosArr"
          :key="todo.id"
          class="elevation-2 rounded-lg mb-4"
          :dark="darkMode"
          :class="{ 'completed-todo': todo.completed }"
        >
          <v-card-text class="pa-6">
            <div class="d-flex align-center mb-4">
              <v-checkbox
                :input-value="todo.completed"
                color="primary"
                @change="handleToggleComplete(todo)"
              ></v-checkbox>
              <v-text-field
                v-if="editingId === todo.id"
                v-model="editTitle"
                :rules="titleRules"
                outlined
                dense
                rounded
                class="mr-2"
              ></v-text-field>
              <span v-else class="text-h6 font-weight-medium" :class="{ 'text-decoration-line-through': todo.completed }">
                {{ todo.title }}
              </span>
              <v-spacer></v-spacer>
              <v-btn
                v-if="editingId === todo.id"
                icon
                small
                color="success"
                @click="handleUpdateTodo(todo)"
              >
                <v-icon>mdi-check</v-icon>
              </v-btn>
              <v-btn
                v-else
                icon
                small
                color="primary"
                @click="startEdit(todo)"
              >
                <v-icon>mdi-pencil</v-icon>
              </v-btn>
              <v-btn
                icon
                small
                color="error"
                @click="handleDeleteTodo(todo.id)"
              >
                <v-icon>mdi-delete</v-icon>
              </v-btn>
            </div>
            <div v-if="editingId === todo.id" class="mb-4">
              <v-textarea
                v-model="editContent"
                outlined
                dense
                rounded
                rows="2"
              ></v-textarea>
            </div>
            <div v-else class="text-body-1 mb-2" :class="{ 'text-decoration-line-through': todo.completed }">
              {{ todo.content || 'No content' }}
            </div>
            <div class="text-caption text--secondary">
              {{ formatDate(todo.created_at) }}
            </div>
          </v-card-text>
        </v-card>

        <v-card v-if="todosArr.length === 0" class="elevation-2 rounded-lg text-center pa-8" :dark="darkMode">
          <v-icon size="64" class="mb-4">mdi-clipboard-text-outline</v-icon>
          <div class="text-h6 mb-2">No todos yet</div>
          <div class="text-body-2 text--secondary">Create your first todo to get started</div>
        </v-card>
      </v-col>
    </v-row>

    <v-snackbar v-model="snackbar" :color="snackbarColor" :timeout="4000" top>
      {{ snackbarText }}
      <template v-slot:action="{ attrs }">
        <v-btn text v-bind="attrs" @click="snackbar = false">Close</v-btn>
      </template>
    </v-snackbar>
  </v-container>
</template>

<script>
import EventService from '../services/EventService'

export default {
  name: 'Todo',
  data() {
    return {
      todosArr: [],
      newTodoTitle: '',
      newTodoContent: '',
      todoFormValid: false,
      creating: false,
      editingId: null,
      editTitle: '',
      editContent: '',
      snackbar: false,
      snackbarText: '',
      snackbarColor: 'error',
      currentUser: {},
      titleRules: [
        v => !!v || 'Title is required',
        v => (v && v.length >= 1) || 'Title must be at least 1 character'
      ]
    }
  },
  computed: {
    darkMode() {
      return this.$vuetify.theme.dark
    }
  },
  mounted() {
    this.loadUser()
    this.loadTodos()
  },
  methods: {
    loadUser() {
      const lUserStr = localStorage.getItem('user')
      if (lUserStr) {
        this.currentUser = JSON.parse(lUserStr)
      }
    },
    loadTodos() {
      const lToken = localStorage.getItem('token')
      if (!lToken) {
        this.$router.push('/login')
        return
      }

      EventService.listTodos(lToken)
        .then((lRes) => {
          if (lRes.data.status === 's') {
            this.todosArr = lRes.data.data || []
          } else {
            this.showSnackbar(lRes.data.message || 'Failed to load todos', 'error')
          }
        })
        .catch((lErr) => {
          if (lErr.response && lErr.response.status === 401) {
            localStorage.removeItem('token')
            localStorage.removeItem('user')
            this.$router.push('/login')
          } else {
            this.showSnackbar('Failed to load todos', 'error')
          }
        })
    },
    handleCreateTodo() {
      if (!this.$refs.todoForm.validate()) {
        return
      }

      this.creating = true
      const lToken = localStorage.getItem('token')
      const lData = {
        title: this.newTodoTitle,
        content: this.newTodoContent
      }

      EventService.createTodo(lData, lToken)
        .then((lRes) => {
          if (lRes.data.status === 's') {
            this.newTodoTitle = ''
            this.newTodoContent = ''
            this.$refs.todoForm.resetValidation()
            this.showSnackbar('Todo created successfully', 'success')
            this.loadTodos()
          } else {
            this.showSnackbar(lRes.data.message || 'Failed to create todo', 'error')
          }
        })
        .catch((lErr) => {
          if (lErr.response && lErr.response.data && lErr.response.data.message) {
            this.showSnackbar(lErr.response.data.message, 'error')
          } else {
            this.showSnackbar('Failed to create todo', 'error')
          }
        })
        .finally(() => {
          this.creating = false
        })
    },
    startEdit(pTodo) {
      this.editingId = pTodo.id
      this.editTitle = pTodo.title
      this.editContent = pTodo.content
    },
    handleUpdateTodo(pTodo) {
      if (!this.editTitle || this.editTitle.trim() === '') {
        this.showSnackbar('Title is required', 'error')
        return
      }

      const lToken = localStorage.getItem('token')
      const lData = {
        title: this.editTitle,
        content: this.editContent,
        completed: pTodo.completed
      }

      EventService.updateTodo(pTodo.id, lData, lToken)
        .then((lRes) => {
          if (lRes.data.status === 's') {
            this.editingId = null
            this.showSnackbar('Todo updated successfully', 'success')
            this.loadTodos()
          } else {
            this.showSnackbar(lRes.data.message || 'Failed to update todo', 'error')
          }
        })
        .catch((lErr) => {
          if (lErr.response && lErr.response.data && lErr.response.data.message) {
            this.showSnackbar(lErr.response.data.message, 'error')
          } else {
            this.showSnackbar('Failed to update todo', 'error')
          }
        })
    },
    handleToggleComplete(pTodo) {
      const lToken = localStorage.getItem('token')
      const lData = {
        title: pTodo.title,
        content: pTodo.content,
        completed: !pTodo.completed
      }

      EventService.updateTodo(pTodo.id, lData, lToken)
        .then((lRes) => {
          if (lRes.data.status === 's') {
            this.loadTodos()
          } else {
            this.showSnackbar(lRes.data.message || 'Failed to update todo', 'error')
          }
        })
        .catch((lErr) => {
          if (lErr.response && lErr.response.data && lErr.response.data.message) {
            this.showSnackbar(lErr.response.data.message, 'error')
          } else {
            this.showSnackbar('Failed to update todo', 'error')
          }
        })
    },
    handleDeleteTodo(pTodoID) {
      if (!confirm('Are you sure you want to delete this todo?')) {
        return
      }

      const lToken = localStorage.getItem('token')
      EventService.deleteTodo(pTodoID, lToken)
        .then((lRes) => {
          if (lRes.data.status === 's') {
            this.showSnackbar('Todo deleted successfully', 'success')
            this.loadTodos()
          } else {
            this.showSnackbar(lRes.data.message || 'Failed to delete todo', 'error')
          }
        })
        .catch((lErr) => {
          if (lErr.response && lErr.response.data && lErr.response.data.message) {
            this.showSnackbar(lErr.response.data.message, 'error')
          } else {
            this.showSnackbar('Failed to delete todo', 'error')
          }
        })
    },
    handleLogout() {
      const lToken = localStorage.getItem('token')
      if (lToken) {
        EventService.logout(lToken)
          .then(() => {
            localStorage.removeItem('token')
            localStorage.removeItem('user')
            this.$router.push('/login')
          })
          .catch(() => {
            localStorage.removeItem('token')
            localStorage.removeItem('user')
            this.$router.push('/login')
          })
      } else {
        this.$router.push('/login')
      }
    },
    toggleTheme() {
      this.$vuetify.theme.dark = !this.$vuetify.theme.dark
      localStorage.setItem('darkMode', this.$vuetify.theme.dark)
    },
    formatDate(pDateStr) {
      if (!pDateStr) return ''
      const lDate = new Date(pDateStr)
      return lDate.toLocaleDateString() + ' ' + lDate.toLocaleTimeString()
    },
    showSnackbar(pText, pColor) {
      this.snackbarText = pText
      this.snackbarColor = pColor
      this.snackbar = true
    }
  }
}
</script>

/* Premium Background Gradient */
.premium-wrapper {
  background: radial-gradient(circle at top left, #1a1a2e, #16213e, #0f3460);
  min-height: 100vh;
}

/* Glassmorphism Effect */
.glass-card {
  background: rgba(255, 255, 255, 0.03) !important;
  backdrop-filter: blur(12px);
  -webkit-backdrop-filter: blur(12px);
  border: 1px solid rgba(255, 255, 255, 0.1) !important;
  border-radius: 24px !important;
  box-shadow: 0 8px 32px 0 rgba(0, 0, 0, 0.37) !important;
}

/* Premium Button with Cinematic Gradient */
.premium-btn {
  background: linear-gradient(45deg, #6366f1, #a855f7) !important;
  color: white !important;
  border-radius: 16px !important;
  text-transform: none !important;
  font-weight: 700 !important;
  letter-spacing: 0.5px;
  transition: all 0.3s ease;
}

.premium-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 10px 20px rgba(99, 102, 241, 0.4) !important;
}

/* Input Styling */
.premium-input {
  background: rgba(0, 0, 0, 0.2) !important;
  border-radius: 16px !important;
  border: 1px solid rgba(255, 255, 255, 0.05);
}

/* Todo Item Styling */
.todo-item-card {
  background: rgba(255, 255, 255, 0.05) !important;
  border-radius: 20px !important;
  border: 1px solid rgba(255, 255, 255, 0.05) !important;
  transition: all 0.4s cubic-bezier(0.175, 0.885, 0.32, 1.275);
}

.todo-item-card:hover {
  background: rgba(255, 255, 255, 0.08) !important;
  transform: scale(1.02);
}

.todo-completed {
  opacity: 0.5;
  filter: grayscale(0.5);
}

/* Text Effects */
.gradient-text {
  background: linear-gradient(to right, #fff, #94a3b8);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.text-strike {
  text-decoration: line-through;
  color: rgba(255, 255, 255, 0.4);
}

