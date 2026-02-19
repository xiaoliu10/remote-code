/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements.  See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership.  The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License.  You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied.  See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

import { createRouter, createWebHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'

// Route configuration
const routes: RouteRecordRaw[] = [
  {
    path: '/',
    redirect: '/app/sessions'
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/Login.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/app',
    component: () => import('@/layouts/MainLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      {
        path: '',
        redirect: '/app/sessions'
      },
      {
        path: 'sessions',
        name: 'Sessions',
        component: () => import('@/views/SessionsView.vue'),
        meta: { requiresAuth: true }
      },
      {
        path: 'sessions/:name',
        name: 'SessionDetail',
        component: () => import('@/views/SessionsView.vue'),
        meta: { requiresAuth: true },
        props: true
      }
    ]
  },
  // Legacy routes for backward compatibility
  {
    path: '/dashboard',
    redirect: '/app/sessions'
  },
  {
    path: '/session/:name',
    redirect: (to) => `/app/sessions/${to.params.name}`
  }
]

// Create router
const router = createRouter({
  history: createWebHistory(),
  routes
})

// Route guard - check authentication
router.beforeEach((to, _from, next) => {
  const token = localStorage.getItem('token')
  const requiresAuth = to.meta.requiresAuth !== false

  if (requiresAuth && !token) {
    next('/login')
  } else if (to.path === '/login' && token) {
    next('/app/sessions')
  } else {
    next()
  }
})

export default router
