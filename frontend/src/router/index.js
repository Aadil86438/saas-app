import Vue from 'vue'
import VueRouter from 'vue-router'
import LoginView from '../views/LoginView.vue'
import SignupView from '../views/SignupView.vue'
import TodoView from '../views/TodoView.vue'

Vue.use(VueRouter)

const routes = [
  {
    path: '/',
    redirect: '/login'
  },
  {
    path: '/login',
    name: 'Login',
    component: LoginView
  },
  {
    path: '/signup',
    name: 'Signup',
    component: SignupView
  },
  {
    path: '/todos',
    name: 'Todos',
    component: TodoView,
    meta: {
      requiresAuth: true
    }
  }
]

const router = new VueRouter({
  mode: 'history',
  base: process.env.BASE_URL,
  routes
})

router.beforeEach((to, from, next) => {
  const lToken = localStorage.getItem('token')
  if (to.matched.some(record => record.meta.requiresAuth)) {
    if (!lToken) {
      next('/login')
    } else {
      next()
    }
  } else {
    if (lToken && (to.path === '/login' || to.path === '/signup')) {
      next('/todos')
    } else {
      next()
    }
  }
})

export default router

