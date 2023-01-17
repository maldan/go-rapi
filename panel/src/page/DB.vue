<template>
  <div :class="$style.main">
    <div :class="$style.header">
      <el-button
        @click="changeTab(x)"
        v-for="x in dbStore.tableList"
        :type="x === dbStore.table ? 'primary' : ''"
        :key="x"
        >{{ x }}</el-button
      >
    </div>

    <el-table
      :data="dbStore.list"
      stripe
      :border="true"
      style="width: 100%"
      :height="tableHeight"
    >
      <el-table-column v-for="(v, k) in dbStore.struct" :prop="k" :label="k" />

      <!-- Edit -->
      <el-table-column label="Edit" width="70">
        <template #default="scope">
          <el-button @click="editById(scope.row.id)" style="width: 100%"
            >Ed</el-button
          >
        </template>
      </el-table-column>

      <!-- Delete -->
      <el-table-column label="Delete" width="70">
        <template #default="scope">
          <el-button @click="deleteById(scope.row.id)" style="width: 100%"
            >x</el-button
          >
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup lang="ts">
import { computed, h, onMounted, ref } from "vue";
import { useDbStore } from "@/store/db";

// Stores
const dbStore = useDbStore();

// Vars
const offset = ref("0");
const tableHeight = ref(400);

// Hooks
onMounted(async () => {
  tableHeight.value = window.innerHeight - 120;

  await dbStore.getTableList();
  dbStore.table = dbStore.tableList[0];
  await dbStore.getStruct();
  await dbStore.getSearch();
});

// Methods
async function deleteById(id: string) {
  await dbStore.deleteById(id);
  await dbStore.getSearch();
}

async function editById(id: string) {}

async function changeTab(tab: string) {
  dbStore.table = tab;
  dbStore.list = [];
  await dbStore.getStruct();
  await dbStore.getSearch();
}
</script>

<style module lang="scss">
.main {
  padding: 10px;
  font-size: 14px;
  height: calc(100% - 80px);

  .header {
    margin-bottom: 10px;
  }
}
</style>
