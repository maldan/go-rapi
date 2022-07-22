<template>
  <ComponentPage />
</template>

<script setup lang="ts">
import { useLogStore } from "@/store/log";
import Log from "../component/Log.vue";
import ComponentPage from "../gam-lib-ui/vue/example/ComponentPage.vue";
import { ref } from "vue";
</script>

<style module lang="scss"></style>
