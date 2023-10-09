import { $t } from "@/lang/index"
import Layout from '@/layout/index.vue'
import Login from '@/views/login.vue'

export const routes = [
    {
        name: "login",
        path: "/login",
        meta: { title: $t('登陆'), login: false },
        component: Login
    },
    {
        name: "/",
        path: "/",
        redirect: "/home",
        meta: { title: $t('优笔记') },
        component: Layout,
        children: [
            {
                name: "home",
                meta: { title: $t('优笔记') },
                path: 'home',
                component: () => import('@/views/index.vue'),
            },
        ],
    }
]