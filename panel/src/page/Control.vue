<template>
  <div :class="$style.main" style="padding: 10px">
    <!-- Left menu -->
    <el-tabs tab-position="left" class="demo-tabs">
      <el-tab-pane v-for="v in tabs" :label="v"> {{ v }} </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup lang="ts">
import { h, onMounted, ref } from "vue";
import { useMethodStore } from "@/store/method";
import { useModalStore } from "@/gam-lib-ui/vue/store/modal";

// Stores
const methodStore = useMethodStore();
const modalStore = useModalStore();

// Vars
const tabs = ["backup", "push", "mail"];

// Hooks
onMounted(async () => {});

// Methods
</script>

<style module lang="scss">
.main {
  display: flex;
  flex-direction: column;
  font-size: 14px;
}
</style>
