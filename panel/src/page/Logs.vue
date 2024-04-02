<template>
  <div style="padding: 10px">
    <!-- <Menu
      :list="[
        { url: '/', label: 'Logs' },
        { url: '/methods', label: 'Methods' },
      ]"
      style="margin-bottom: 10px"
    /> -->

    <div style="margin-bottom: 5px; font-size: 16px">Select logs category:</div>
    <div style="margin-bottom: 5px; display: flex; gap: 5px">
      <button
        @click="category = x.name"
        v-for="x in categoryList"
        :key="x.name"
        :class="x.name === category ? $style.selected : null"
      >
        {{ x.name }}
      </button>
    </div>

    <div v-if="category">
      <div style="margin-bottom: 5px; font-size: 16px">Select date range:</div>
      <div
        style="
          margin-bottom: 10px;
          display: flex;
          gap: 10px;
          align-items: center;
        "
      >
        <div style="font-size: 16px">From:</div>
        <input type="date" v-model="fromDate" style="padding: 5px" />
        <div style="font-size: 16px">To:</div>
        <input type="date" v-model="toDate" style="padding: 5px" />
      </div>
      <button @click="download">Download</button>
    </div>

    <!-- List -->
    <!-- <Table
      :format="{
        time: 'date:YYYY-MM-DD HH:mm:ss.SSS',
        kind: (x) => {
          const y = [$style.kind, $style[x]].join(' ');
          return `<div class='${y}'>${x}</div>`;
        },
      }"
      :list="logs.items"
    /> -->
  </div>
</template>

<script setup lang="ts">
import { useLogStore } from "@/store/log";
import Log from "../component/Log.vue";
import { Button, Menu, Table, Input } from "../gam-lib-ui/vue/component/ui";
import { onMounted, ref } from "vue";
import dayjs from "dayjs";
import { HOST } from "@/const";

const logs = useLogStore();

let fromDate = ref(dayjs().add(-7, "day").format("YYYY-MM-DD"));
let toDate = ref(dayjs().format("YYYY-MM-DD"));
let categoryList = ref([]);
let category = ref("");

onMounted(async () => {
  categoryList.value = await logs.getCategoryList();
});

function download() {
  let from = dayjs(fromDate.value).format("YYYY-MM-DD");
  let to = dayjs(toDate.value).format("YYYY-MM-DD");
  window.open(
    `${HOST}/debug/log/download?fromDate=${from}&toDate=${to}&category=${category.value}`,
    "_blank"
  );
}
</script>

<style module lang="scss">
button {
  border: 0;
  border-radius: 4px;
  outline: none;
  padding: 8px 15px;
  cursor: pointer;

  &:hover {
    opacity: 0.8;
  }

  &.selected {
    background: #387ad0;
  }
}

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
