<template>
  <div>
    <button>Methods</button>
    <button>Logs</button>
  </div>

  <div>
    <input type="date" v-model="date" @change="refresh" />
  </div>

  <!-- List -->
  <div style="padding: 5px">
    <Log
      v-for="x in logs.items"
      :key="x.message"
      :kind="x.kind"
      :message="x.message"
      :time="x.time"
      style="margin-bottom: 5px"
    />
  </div>
</template>

<script setup lang="ts">
import { useLogStore } from "../store/log";
import Log from "../component/Log.vue";
import { ref } from "vue";

const logs = useLogStore();

let date = ref("");

function refresh() {
  console.log(date);
  logs.getList(date.value);
}

refresh();
</script>

<style module lang="scss"></style>
