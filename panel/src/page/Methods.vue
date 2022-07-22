<template>
  <div :class="$style.main">
    <Menu
      :list="[
        { url: '/', label: 'Logs' },
        { url: '/methods', label: 'Methods' },
      ]"
    />

    <!-- List -->
    <Table
      :list="methodStore.items"
      :component="component"
      :format="format"
      style="margin-top: 10px"
    />
  </div>
</template>

<script setup lang="ts">
import { Button, Menu, Table, Tree } from "../gam-lib-ui/vue/component/ui";
import { h, ref } from "vue";
import { useMethodStore } from "@/store/method";
import { useModalStore } from "@/gam-lib-ui/vue/store/modal";

const methodStore = useMethodStore();
const modalStore = useModalStore();
methodStore.getList();

const format = {
  input: undefined,
  url: undefined,
};
const component = {
  run: (x: any) => {
    return h(Button, {
      innerHTML: "run",
      onClick: () => {
        openModal(x.url, x.httpMethod, x.input);
      },
    });
  },
};

// Methods
function openModal(url: string, method: string, args: any) {
  const input = {} as Record<string, any>;
  const fl = args.fieldList;
  for (let i = 0; i < fl.length; i++) {
    if (fl[i].kind === "string") input[fl[i].name] = "";
    if (fl[i].kind === "bool") input[fl[i].name] = false;
    if (fl[i].kind === "slice") input[fl[i].name] = [];
    if (fl[i].kind === "map") input[fl[i].name] = {};
    if (fl[i].kind.match(/^int|^float/)) input[fl[i].name] = 0;
  }

  modalStore.show(
    "method/run",
    {
      url,
      method,
      struct: args,
      input,
    },
    () => {}
  );
}
</script>

<style module lang="scss">
.main {
  overflow-y: auto;
  height: 100%;
  padding: 10px;
}
</style>
