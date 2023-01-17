<template>
  <div :class="$style.mainApp">
    <el-menu
      :default-active="$route.path"
      class="el-menu-demo"
      mode="horizontal"
      @select="handleSelect"
    >
      <el-menu-item index="/logs">Logs</el-menu-item>
      <el-menu-item index="/requests">Requests</el-menu-item>
      <el-menu-item index="/methods">Methods</el-menu-item>
      <el-menu-item index="/tests">Tests</el-menu-item>
      <el-menu-item index="/db">DB</el-menu-item>
      <el-menu-item index="/settings">Settings</el-menu-item>
    </el-menu>

    <RouterView />
  </div>
</template>

<script lang="ts" setup>
import { RouterView, useRouter } from "vue-router";

const router = useRouter();

const handleSelect = (key: string, keyPath: string[]) => {
  router.push(key);
};
</script>

<style lang="scss" module>
.mainApp {
  height: 100%;
}
</style>
