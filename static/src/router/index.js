import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

/* Layout */
import Layout from '@/layout'

/**
 * Note: sub-menu only appear when route children.length >= 1
 * Detail see: https://panjiachen.github.io/vue-element-admin-site/guide/essentials/router-and-nav.html
 *
 * hidden: true                   if set true, item will not show in the sidebar(default is false)
 * alwaysShow: true               if set true, will always show the root menu
 *                                if not set alwaysShow, when item has more than one children route,
 *                                it will becomes nested mode, otherwise not show the root menu
 * redirect: noRedirect           if set noRedirect will no redirect in the breadcrumb
 * name:'router-name'             the name is used by <keep-alive> (must set!!!)
 * meta : {
    roles: ['admin','editor']    control the page roles (you can set multiple roles)
    title: 'title'               the name show in sidebar and breadcrumb (recommend set)
    icon: 'svg-name'/'el-icon-x' the icon show in the sidebar
    breadcrumb: false            if set false, the item will hidden in breadcrumb(default is true)
    activeMenu: '/example/list'  if set path, the sidebar will highlight the path you set
  }
 */

/**
 * constantRoutes
 * a base page that does not have permission requirements
 * all roles can be accessed
 */
export const constantRoutes = [
  {
    path: '/login',
    component: () => import('@/views/login/index'),
    hidden: true
  },

  {
    path: '/404',
    component: () => import('@/views/404'),
    hidden: true
  },

  {
    path: '/',
    component: Layout,
    redirect: '/dashboard',
    children: [{
      path: 'dashboard',
      name: 'Dashboard',
      component: () => import('@/views/dashboard/index'),
      meta: { title: 'Dashboard', icon: 'dashboard' }
    }]
  },

  {
    path: '/conf',
    component: Layout,
    redirect: '/conf/job',
    name: 'Conf',
    meta: { title: '配置中心', icon: 'el-icon-s-tools'},
    children: [
      {
        path: 'job',
        name: 'Job',
        component: () => import('@/views/conf/job'),
        meta: { title: '任务配置', icon: 'el-icon-s-operation' }
      },
      {
        path: 'group',
        name: 'Group',
        component: () => import('@/views/conf/group'),
        meta: { title: '集群配置', icon: 'tree' }
      }
    ]
  },

  {
    path: '/schedule',
    component: Layout,
    redirect: '/schedule/node',
    name: 'Example',
    meta: { title: '调度中心', icon: 'el-icon-s-help' },
    children: [
      {
        path: 'node',
        name: 'Node',
        component: () => import('@/views/schedule/node'),
        meta: { title: '调度集群', icon: 'el-icon-cpu' }
      },
      {
        path: 'plan',
        name: 'Plan',
        component: () => import('@/views/schedule/plan'),
        meta: { title: '执行计划', icon: 'el-icon-s-promotion' }
      },{
        path: 'snapshot',
        name: 'Snapshot',
        component: () => import('@/views/schedule/snapshot'),
        meta: { title: '调度快照', icon: 'el-icon-video-camera-solid' }
      },{
        path: 'execute',
        name: 'Execute',
        component: () => import('@/views/schedule/execute'),
        meta: { title: '任务快照', icon: 'el-icon-message-solid' }
      },{
        path: 'execute_history',
        name: 'ExecuteHistory',
        component: () => import('@/views/schedule/execute_history'),
        meta: { title: '任务历史快照', icon: 'el-icon-video-camera-solid' }
      }
    ]
  },
  {
    path: 'external-link',
    component: Layout,
    children: [
      {
        path: 'https://github.com/busgo/pink',
        meta: { title: '项目地址', icon: 'link' }
      }
    ]
  },

  // 404 page must be placed at the end !!!
  { path: '*', redirect: '/404', hidden: true }
]

const createRouter = () => new Router({
  // mode: 'history', // require service support
  scrollBehavior: () => ({ y: 0 }),
  routes: constantRoutes
})

const router = createRouter()

// Detail see: https://github.com/vuejs/vue-router/issues/1234#issuecomment-357941465
export function resetRouter() {
  const newRouter = createRouter()
  router.matcher = newRouter.matcher // reset router
}

export default router
