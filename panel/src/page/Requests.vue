<template>
  <div :class="$style.main">
    <div :class="$style.header">
      <el-input
        placeholder="Offset..."
        v-model="offset"
        style="margin-right: 5px"
        @change="refresh"
      />
      <el-input placeholder="Limit..." v-model="limit" @change="refresh" />
    </div>
    <el-table
      :data="requestStore.list"
      stripe
      :border="true"
      style="width: 100%"
      :height="tableHeight"
    >
      <!-- Method tag -->
      <el-table-column label="Method" width="100">
        <template #default="scope">
          <MethodTag :tag="scope.row.httpMethod" />
        </template>
      </el-table-column>
      <el-table-column prop="url" label="Url" />

      <el-table-column label="Args">
        <template #default="scope">
          <pre
            v-html="formatHighlight(scope.row.args || {}, customColorOptions)"
          ></pre>
        </template>
      </el-table-column>

      <el-table-column label="Response">
        <template #default="scope">
          <pre
            v-html="
              formatHighlight(scope.row.response || {}, customColorOptions)
            "
          ></pre>
        </template>
      </el-table-column>

      <el-table-column label="Error">
        <template #default="scope">
          {{ scope.row.error.status ? scope.row.error : "-" }}
        </template>
      </el-table-column>

      <el-table-column label="Created" width="150">
        <template #default="scope">
          {{ dayjs(scope.row.created).format("HH:mm:ss.SSS") }}
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup lang="ts">
import { computed, h, onMounted, ref } from "vue";
import { useRequestStore } from "@/store/request";
import MethodTag from "@/component/MethodTag.vue";
import dayjs from "dayjs";
import formatHighlight from "json-format-highlight";

// Stores
const requestStore = useRequestStore();

// Vars
const offset = ref("0");
const limit = ref("100");
const tableHeight = ref(400);
const customColorOptions = ref({
  keyColor: "#af6ed1",
  numberColor: "#77b0fc",
  stringColor: "#57ab51",
  trueColor: "#ff8080",
  falseColor: "#ff8080",
  nullColor: "#e54b4b",
});

// Hooks
onMounted(async () => {
  tableHeight.value = window.innerHeight - 120;

  await refresh();
});

// Methods
async function refresh() {
  await requestStore.getList(~~offset.value, ~~limit.value);
}
</script>

<style module lang="scss">
.main {
  padding: 10px;
  font-size: 14px;
  height: calc(100% - 80px);

  .header {
    display: flex;
    margin-bottom: 10px;
  }

  pre {
    background: #101010;
    padding: 5px;
    box-sizing: border-box;
    word-break: break-all;
    font-size: 14px;
    margin: 0;
    max-height: 300px;
    overflow-y: auto;
    line-height: 16px;
  }

  .request {
    display: flex;
    margin-bottom: 10px;

    > div {
      flex: 1;
    }
  }
}
</style>
