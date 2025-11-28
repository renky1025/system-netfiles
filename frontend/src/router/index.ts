import { createRouter, createWebHistory } from 'vue-router';
import { useUserStore } from '../store/user';
import { ElMessage } from 'element-plus';
import Login from '../views/Login.vue';
import Register from '../views/Register.vue';
import Home from '../views/Home.vue';
import RecycleBin from '../views/RecycleBin.vue';
import MyShares from '../views/MyShares.vue';

const routes = [
    { path: '/login', component: Login },
    { path: '/register', component: Register },
    { path: '/', component: Home, meta: { requiresAuth: true } },
    { path: '/folder/:folderId', component: Home, meta: { requiresAuth: true } },
    { path: '/recycle', component: RecycleBin, meta: { requiresAuth: true } },
    { path: '/shares', component: MyShares, meta: { requiresAuth: true } },
    { path: '/share/:code', component: () => import('../views/ShareAccess.vue'), meta: { requiresAuth: false } },
    { path: '/forgot-password', component: () => import('../views/ForgotPassword.vue') },
    { path: '/reset-password', component: () => import('../views/ResetPassword.vue') },
    { path: '/change-password', component: () => import('../views/ChangePassword.vue'), meta: { requiresAuth: true } },

    // Admin Routes
    {
        path: '/admin',
        component: () => import('../layouts/AdminLayout.vue'),
        meta: { requiresAuth: true, requiresAdmin: true },
        children: [
            { path: '', redirect: '/admin/dashboard' },
            { path: 'dashboard', component: () => import('../views/admin/Dashboard.vue') },
            { path: 'users', component: () => import('../views/admin/UserList.vue') },
            { path: 'orgs', component: () => import('../views/admin/OrgManagement.vue') },
            { path: 'roles', component: () => import('../views/admin/RoleManagement.vue') },
            { path: 'files', component: () => import('../views/admin/FileList.vue') },
            { path: 'shares', component: () => import('../views/admin/ShareList.vue') },
            { path: 'logs', component: () => import('../views/admin/AuditLogs.vue') },
            { path: 'settings', component: () => import('../views/admin/SystemSettings.vue') },
        ]
    }
];

const router = createRouter({
    history: createWebHistory(),
    routes,
});

router.beforeEach((to, _from, next) => {
    const token = localStorage.getItem('token');
    if (to.meta.requiresAuth && !token) {
        next('/login');
    } else if (to.meta.requiresAdmin) {
        const userStore = useUserStore();
        console.log('Admin check:', {
            user: userStore.user,
            isAdmin: userStore.isAdmin,
            roles: userStore.user?.roles
        });
        if (!userStore.isAdmin) {
            ElMessage.error('Access denied: Admins only');
            next('/');
            return;
        }
        next();
    } else {
        next();
    }
});

export default router;
