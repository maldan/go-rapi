<template>
  <div style="padding: 10px">
    <Menu
      :list="[
        { url: '/', label: 'Logs' },
        { url: '/methods', label: 'Methods' },
      ]"
      style="margin-bottom: 10px"
    />

    <div style="margin-bottom: 10px">
      <input type="date" v-model="date" @change="refresh" />
    </div>

    <!-- List -->
    <Table
      :format="{
        time: 'date:YYYY-MM-DD HH:mm:ss.SSS',
        kind: (x) => {
          const y = [$style.kind, $style[x]].join(' ');
          return `<div class='${y}'>${x}</div>`;
        },
      }"
      :list="logs.items"
    />
  </div>
</template>

<script setup lang="ts">
import { useLogStore } from "@/store/log";
import Log from "../component/Log.vue";
import { Button, Menu, Table, Input } from "../gam-lib-ui/vue/component/ui";
import { ref } from "vue";

const logs = useLogStore();

let date = ref("");

function refresh() {
  logs.getList(date.value);
}

refresh();
</script>

<style module lang="scss">
.kind {
  padding: 2px 15px 2px 5px;
  background: #a31c1c;
  color: #fefefe;
  // border-radius: 4px;
  font-weight: bold;
  text-transform: uppercase;
  font-size: 10px;
  margin-right: 15px;
  position: relative;
  width: 32px;
  height: 14px;

  &::after {
    content: "";
    display: block;
    position: absolute;
    width: 12px;
    height: 12px;
    transform: rotate(45deg);
    right: -7px;
    top: 3px;
    background: #2b2b2b;
  }

  &.error {
    background: #a31c1c;

    &::after {
      background: #a31c1c;
    }
  }

  &.info {
    background: #198bb4;

    &::after {
      background: #198bb4;
    }
  }
}
</style>
