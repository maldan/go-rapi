<template>
  <div :class="$style.main" style="padding: 10px">
    <div>Access token</div>
    <el-input
      placeholder="Access token"
      @change="onSave"
      v-model="accessToken"
    ></el-input>
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
const accessToken = ref("");

// Hooks
onMounted(async () => {
  accessToken.value = localStorage.getItem("debug__accessToken") || "";
});

// Methods
function onSave() {
  localStorage.setItem("debug__accessToken", accessToken.value);
}
</script>

<style module lang="scss">
.main {
  display: flex;
  flex-direction: column;
}
</style>
