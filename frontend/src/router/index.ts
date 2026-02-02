import { createRouter, createWebHistory } from 'vue-router'
import TicketsList from '@/pages/Tickets/List.vue'
import TicketDetail from '@/pages/Tickets/Detail.vue'
import CreateTicket from '@/pages/Tickets/Create.vue'
import ArticlesIndex from '@/pages/Articles/Index.vue'
import ArticleDetail from '@/pages/Articles/Detail.vue'
import Dashboard from '@/pages/Dashboard/Index.vue'
import Login from '@/pages/Auth/Login.vue'
import Landing from '@/pages/Landing.vue'
import Profile from '@/pages/Profile/Index.vue'
import { getCookie } from '@/utils/cookies'
import { COOKIE_NAMES } from '@/utils/constants'

const isAuthenticated = (): boolean => {
  return !!getCookie(COOKIE_NAMES.ACCESS_TOKEN)
}

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'landing',
      component: Landing,
      meta: { requiresGuest: false },
    },
    {
      path: '/dashboard',
      name: 'dashboard',
      component: Dashboard,
      meta: { requiresAuth: true },
    },
    {
      path: '/login',
      name: 'login',
      component: Login,
      meta: { requiresGuest: true },
    },
    {
      path: '/register',
      name: 'Register',
      component: () => import('@/pages/Auth/Register.vue'),
      meta: { requiresGuest: true },
    },
    {
      path: '/forgot-password',
      name: 'ForgotPassword',
      component: () => import('@/pages/Auth/ForgotPassword.vue'),
      meta: { requiresGuest: true },
    },
    {
      path: '/reset-password',
      name: 'ResetPassword',
      component: () => import('@/pages/Auth/ResetPassword.vue'),
      meta: { requiresGuest: true },
    },
    {
      path: '/profile',
      name: 'profile',
      component: Profile,
      meta: { requiresAuth: true },
    },
    {
      path: '/tickets',
      name: 'tickets',
      component: TicketsList,
      meta: { requiresAuth: true },
    },
    {
      path: '/tickets/create',
      name: 'create-ticket',
      component: CreateTicket,
      meta: { requiresAuth: true },
    },
    {
      path: '/tickets/:id',
      name: 'ticket-detail',
      component: TicketDetail,
      props: true,
      meta: { requiresAuth: true },
    },
    {
      path: '/articles',
      name: 'articles',
      component: ArticlesIndex,
      meta: { requiresAuth: true },
    },
    {
      path: '/articles/:id',
      name: 'article-detail',
      component: ArticleDetail,
      props: true,
      meta: { requiresAuth: false }, // Allow guest access to article detail
    },
    // Admin Routes
    {
      path: '/admin/tenants',
      name: 'admin-tenants',
      component: () => import('@/pages/Admin/Tenants/List.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/admin/tenants/create',
      name: 'create-tenant',
      component: () => import('@/pages/Admin/Tenants/Form.vue'),
      meta: { requiresAuth: true },
    },
    {
      path: '/admin/tenants/:id/edit',
      name: 'edit-tenant',
      component: () => import('@/pages/Admin/Tenants/Form.vue'),
      props: true,
      meta: { requiresAuth: true },
    },
  ],
})

router.beforeEach((to, _from, next) => {
  const authenticated = isAuthenticated()

  if (to.meta.requiresAuth && !authenticated) {
    next({ name: 'login', query: { redirect: to.fullPath } })
    return
  }

  if (to.meta.requiresGuest && authenticated) {
    next({ name: 'dashboard' })
    return
  }

  next()
})

export default router

