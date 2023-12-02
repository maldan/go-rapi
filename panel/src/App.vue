<template>
  <div :class="$style.mainApp">
    <el-menu
      v-if="authStore.isAuth"
      :default-active="$route.path"
      class="el-menu-demo"
      mode="horizontal"
      @select="handleSelect"
    >
      <el-menu-item index="/requests" style="color: #fe6e3d">
        <el-icon><Promotion /></el-icon>Requests
      </el-menu-item>
      <el-menu-item index="/logs" style="color: #c1fe48">
        <el-icon><EditPen /></el-icon>Logs
      </el-menu-item>
      <el-menu-item index="/methods" style="color: #c1fe48">
        <el-icon><EditPen /></el-icon>Methods
      </el-menu-item>
      <el-menu-item index="/tests" style="color: #fea048">
        <el-icon><WarnTriangleFilled /></el-icon>Tests
      </el-menu-item>
      <el-menu-item index="/db" style="color: #ec48fe">
        <el-icon><Document /></el-icon>DB
      </el-menu-item>
      <el-menu-item index="/control" style="color: #fed448">
        <el-icon><Operation /></el-icon>Control
      </el-menu-item>

      <!-- Charts -->
      <el-menu-item index="/chart" style="color: #fed448">
        <el-icon><DataAnalysis /></el-icon>Chart
      </el-menu-item>

      <!-- Backup -->
      <el-menu-item index="/backup" style="color: #fed448">
        <el-icon><DataAnalysis /></el-icon>Backup
      </el-menu-item>

      <el-menu-item index="/settings">
        <el-icon><Tools /></el-icon>Settings
      </el-menu-item>
    </el-menu>

    <RouterView v-if="authStore.isAuth" />

    <div
      v-if="!authStore.isAuth"
      style="padding: 10px; display: flex; flex-direction: column"
    >
      <el-input
        placeholder="Login"
        v-model="login"
        style="margin-bottom: 10px"
      />
      <el-input
        placeholder="Password"
        v-model="password"
        style="margin-bottom: 10px"
      />
      <el-button @click="auth" type="primary" style="flex: 1">Login</el-button>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { RouterView, useRouter } from "vue-router";
import { onMounted, ref } from "vue";
import { useAuthStore } from "@/store/auth";
import Axios from "axios";

const router = useRouter();
const authStore = useAuthStore();
const login = ref("");
const password = ref("");

const handleSelect = (key: string, keyPath: string[]) => {
  router.push(key);
};

onMounted(() => {
  authStore.check(localStorage.getItem("debugAuthKey") + "");
});

async function auth() {
  await authStore.auth(login.value, password.value);
  window.location.reload();
}
</script>

<style lang="scss" module>
.mainApp {
  height: 100%;
}
</style>
