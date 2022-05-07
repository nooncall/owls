<template>
  <div>
    <router-view v-slot="{ Component }">
      <transition mode="out-in" name="el-fade-in-linear">
        <keep-alive :include="routerStore.keepAliveRouters">
          <component :is="Component" />
        </keep-alive>
      </transition>
    </router-view>
  </div>
</template>

<script>
export default {
  name: 'System'
}
</script>
<script setup>
import { useRouterStore } from '@/pinia/modules/router'
const routerStore = useRouterStore()
</script>
// 申请页，列表、创建、撤销
// 审批页，简单通过或拒绝，无历史页

// 初始化名改为管理员角色，管理员拥有所有权限

// auth： 列表、添加、删除。

// 然后就可以写查询啦！
