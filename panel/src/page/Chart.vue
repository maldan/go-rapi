<template>
  <div :class="$style.main" style="padding: 10px">
    <!-- Left menu -->
    <el-tabs tab-position="left" class="demo-tabs" @tab-click="changeTab">
      <el-tab-pane v-for="v in tabs" :label="v" :name="v">
        <!-- Command list -->
        <div
          v-for="cmd in chartStore.chartList.filter((x) => x.folder === v)"
          :key="cmd"
        >
          <el-button
            type="primary"
            :loading="chartStore.status[cmd.folder + '/' + cmd.name]"
            @click="chartStore.execute(cmd.folder, cmd.name)"
            >{{ cmd.name }}</el-button
          >

          <!-- Chart -->
          <Chart
            :size="{ width: 1200, height: 300 }"
            :data="chartStore.response[cmd.folder + '/' + cmd.name]"
            :margin="{ left: 0, top: 20, right: 20, bottom: 0 }"
            :direction="'horizontal'"
          >
            <template #layers>
              <Grid strokeDasharray="2,2" />

              <Bar
                :dataKeys="['value', 'value']"
                :barStyle="{ fill: '#90e0ef' }"
              />
            </template>
          </Chart>

          <!--          <fusioncharts
            :type="'column2d'"
            :width="'100%'"
            :height="300"
            :dataFormat="'json'"
            :dataSource="{
              chart: {
                caption: 'Chart',
                theme: 'candy',
              },
              data: chartStore.response[cmd.folder + '/' + cmd.name],
            }"
          ></fusioncharts>-->
        </div>

        <!--        &lt;!&ndash; Data &ndash;&gt;
        <div>{{ chartStore.status[] }}</div>-->
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup lang="ts">
import { computed, h, onMounted, ref } from "vue";
import { useControlStore } from "@/store/control";
import { useChartStore } from "@/store/chart";
import { Chart, Grid, Line, Bar } from "vue3-charts";

// Stores
const chartStore = useChartStore();

// Vars
const tabs = computed(() => {
  return [...new Set(chartStore.chartList.map((x) => x.folder))];
});

// Hooks
onMounted(async () => {
  await chartStore.getCommandList();
  if (chartStore.chartList.length > 0) {
    changeTab(chartStore.chartList[0].folder);
  }
});

// Methods
function changeTab(x: any) {
  chartStore.executeFolder(x.paneName);
}
</script>

<style module lang="scss">
.main {
  display: flex;
  flex-direction: column;
  font-size: 14px;
}
</style>
