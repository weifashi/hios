<template>
    12

    

</template>

<script lang="ts" setup>
import { ref } from 'vue'
import { useRoute } from 'vue-router';

const route = useRoute()
const OkrFollowRef = ref(null)

watch(route,(newValue)=>{
    nextTick(()=>{
        if(newValue.query.active == undefined && OkrFollowRef.value != null){
            OkrFollowRef.value.getList('search')
        }
    })
},{immediate:true})

</script>

<style lang="less" scoped>
.page-okr {
    @apply absolute top-0 bottom-0 left-0 right-0 flex flex-col bg-page-bg  px-16 py-20 md:px-20;

    .okr-title {
        @apply h-42 flex justify-between items-center relative mt-12 mb-14;

        .icon-return {
            @apply block md:hidden mr-16 text-20 z-[2];
        }

        h2 {
            @apply text-title-color text-28 font-semibold;
        }

        .title-active {
            @apply hidden md:block;
        }

        .okr-right {
            @apply flex items-center gap-4 md:gap-6 z-[2];

            .add-button,
            .search-button {
                @apply bg-[#f2f3f5] w-36 h-36 rounded-full flex items-center justify-center cursor-pointer;

                i {
                    @apply text-20 text-emoji-users-color;
                }
            }

            .search-button {
                @apply flex items-center;

                i,
                span {
                    @apply flex-initial;
                }

                .search-button-span {
                    @apply text-14  text-emoji-users-color pr-8 border-solid border-0 border-r;
                }

                :deep(.n-input) {
                    @apply flex-1 bg-transparent border-0;

                    .n-input__border,
                    .n-input__state-border,
                    .n-input__border:focus,
                    .n-input__state-border:focus {
                        @apply border-0 shadow-none;
                    }
                }
            }

            .search-active {
                @apply w-auto flex-1 md:w-320 px-14;

                i {
                    @apply pl-8;
                }
            }
        }
    }

    .okr-tabs {
        @apply flex-1 relative;

        :deep(.n-tabs) {
            @apply absolute left-0 right-0 top-0 bottom-0;

            .n-tabs-pane-wrapper {
                @apply flex-1 relative;

                .n-tab-pane {
                    @apply max-h-full absolute left-0 right-0 top-0 bottom-0;

                    .okr-scrollbar {
                        @apply absolute left-0 right-0 top-0 bottom-0;
                    }
                }
            }
        }
    }
}

//
body.window-portrait {
    .page-okr {
        .okr-title {
            margin: 4px 0 14px -4px;
        }
    }
}
</style>
