<template>

    <!--  -->
    欢迎使用

    <n-space vertical v-if="0">
        <n-layout class="root-layout" has-sider>
            <!-- <n-layout-header>颐和园路</n-layout-header> -->
            <n-layout-sider content-style="padding: 24px;">
                海淀桥
                <!-- <n-menu
                    v-model:value="activeKey"
                    :collapsed="collapsed"
                    :collapsed-width="64"
                    :collapsed-icon-size="22"
                    :options="menuOptions"
                /> -->
            </n-layout-sider>
            <n-layout style="height: 100%;">
                <n-layout-header>颐和园路</n-layout-header>
                <n-layout-content content-style="padding: 24px;">
                    <router-view />
                </n-layout-content>
                <n-layout-footer>成府路</n-layout-footer>
            </n-layout>
        </n-layout>
    </n-space>
</template>

<script lang="ts" setup>
import { onMounted, watch } from 'vue'
import { useLoadingBar, useMessage } from 'naive-ui'
import { loadingBarApiRef } from "../routes";
import { UserStore } from "../store/user";
import { GlobalStore } from '@/store';

const message = useMessage()
const userStore = UserStore()
const loadingBar = useLoadingBar()
const globalStore = GlobalStore()

nextTick(()=>{
    window.$message = message
})

watch(
    () => userStore.info,
    () => { },
    { immediate: true }
)

onMounted(() => {
    loadingBarApiRef.value = loadingBar
    loadingBar.finish()
    window.addEventListener('scroll', windowScrollListener);
})

const windowScrollListener = () => {
    globalStore.$patch((state) => {
        state.windowScrollY = window.scrollY
    })
}

const otherEvents = () => {
    // 非客户端监听窗口激活
    const hiddenProperty = 'hidden' in document ? 'hidden' : 'webkitHidden' in document ? 'webkitHidden' : 'mozHidden' in document ? 'mozHidden' : '';
    const visibilityChangeEvent = hiddenProperty.replace(/hidden/i, 'visibilitychange');
    document.addEventListener(visibilityChangeEvent, () => {
        globalStore.$patch((state) => {
            state.windowActive = !document[hiddenProperty]
        })
    });
}


otherEvents()
</script>

<style lang="less" scoped>
.root-layout {
    @apply absolute w-full bottom-0 top-0 right-0 left-0 min-h-full;
    .child-view {
        @apply absolute w-full bottom-0 top-0 right-0 left-0 min-h-full;
    }
}
</style>
