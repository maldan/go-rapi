<template>
  <div style="padding: 10px">
    <!-- Data -->
    <el-table
      :data="backupStore.taskList"
      stripe
      :border="true"
      style="width: 100%"
      :height="tableHeight"
      empty-text="None"
    >
      <el-table-column label="Id" width="180">
        <template #default="scope"> {{ scope.row.id }} </template>
      </el-table-column>
      <el-table-column label="Src">
        <template #default="scope"> {{ scope.row.src }} </template>
      </el-table-column>
      <el-table-column label="Src Temp">
        <template #default="scope"> {{ scope.row.srcTemp }} </template>
      </el-table-column>
      <el-table-column label="Dst">
        <template #default="scope"> {{ scope.row.dst }} </template>
      </el-table-column>
      <el-table-column label="Period">
        <template #default="scope"> {{ scope.row.period }} </template>
      </el-table-column>
      <el-table-column label="Status">
        <template #default="scope"> {{ scope.row.status }} </template>
      </el-table-column>
      <el-table-column label="Error">
        <template #default="scope"> {{ scope.row.error }} </template>
      </el-table-column>
      <el-table-column label="Last Run">
        <template #default="scope">
          {{ dayjs(scope.row.lastRun).format("DD MMM YYYY HH:mm:ss") }}
        </template>
      </el-table-column>
      <el-table-column label="Last Run">
        <template #default="scope">
          {{ dayjs(scope.row.nextRun).format("DD MMM YYYY HH:mm:ss") }}
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup lang="ts">
import { useLogStore } from "@/store/log";
import { onMounted, ref } from "vue";
import { useBackupStore } from "@/store/backup";
import dayjs from "dayjs";

const backupStore = useBackupStore();
const tableHeight = ref(400);

onMounted(async () => {
  tableHeight.value = window.innerHeight - 80;

  await backupStore.getTaskList();
});
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
