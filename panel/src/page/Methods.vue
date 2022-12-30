<template>
  <div :class="$style.main">
    <!-- Left menu -->
    <el-tabs tab-position="left" class="demo-tabs">
      <el-tab-pane v-for="(controller, k) in controllerList" :label="k">
        <el-alert
          style="margin-bottom: 5px"
          :title="k"
          type="info"
          :closable="false"
        />

        <el-table
          @row-click="rowClick"
          :data="controller"
          stripe
          :border="true"
          style="width: 100%"
        >
          <!-- Method tag -->
          <el-table-column label="Method" width="100">
            <template #default="scope">
              <MethodTag :tag="scope.row.httpMethod" />
            </template>
          </el-table-column>

          <el-table-column prop="name" label="Name" width="150" />
          <el-table-column prop="fullPath" label="FullPath" />

          <!-- Params -->
          <el-table-column label="Params" width="220">
            <template #default="scope">
              <div v-if="scope.row.inputMethod === 'json'">
                {{ methodStore.params[scope.row.uid] }}
              </div>
              <div v-if="scope.row.inputMethod === 'multipart'">
                {{ methodStore.formData[scope.row.uid] }}
              </div>
            </template>
          </el-table-column>

          <!-- Show params -->
          <el-table-column label="Params" width="85">
            <template #default="scope">
              <el-button
                @click="
                  rowClick(scope.row);
                  dialogVisible = true;
                "
                >Edit</el-button
              >
            </template>
          </el-table-column>

          <!-- Run -->
          <el-table-column label="Run" width="85">
            <template #default="scope">
              <el-button
                :type="selectedMethodUid === scope.row.uid ? 'primary' : null"
                @click="run(scope.row)"
                >Run</el-button
              >
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>
    </el-tabs>

    <!-- Response -->
    <div :class="$style.response" style="">
      <div :class="$style.header">
        <div>{{ selectedMethodUid }}</div>
        <div style="margin-left: 15px">
          Status: {{ methodStore.responseInfo[selectedMethodUid]?.status }}
        </div>
        <div style="margin-left: 15px">
          Time: {{ methodStore.responseInfo[selectedMethodUid]?.time }} ms
        </div>
      </div>

      <pre
        style="
          overflow-y: auto;
          height: calc(100% - 100px);
          background: #262626;
          padding: 10px;
          box-sizing: border-box;
          word-break: break-all;
        "
        v-html="
          formatHighlight(
            methodStore.response[selectedMethodUid] || {},
            customColorOptions
          )
        "
      ></pre>
    </div>

    <!-- Method params -->
    <el-dialog
      v-model="dialogVisible"
      :title="selectedMethodUid + ' params'"
      width="40%"
      draggable
    >
      <el-input
        v-if="methodStore.methodsMap[selectedMethodUid].inputMethod === 'json'"
        v-model="methodStore.params[selectedMethodUid]"
        :rows="10"
        type="textarea"
        placeholder="Raw json..."
      />

      <!-- Multipart -->
      <div
        v-if="
          methodStore.methodsMap[selectedMethodUid].inputMethod === 'multipart'
        "
      >
        <div
          v-for="x in methodStore.methodsMap[selectedMethodUid].input.fieldList"
        >
          <input
            v-if="x.type === 'rapi_core.File'"
            :placeholder="x.name"
            type="file"
            @change="
              fileSelect(
                $event,
                x.name,
                methodStore.formData[selectedMethodUid]
              )
            "
          />
        </div>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { Button, Menu, Table, Tree } from "../gam-lib-ui/vue/component/ui";
import { h, onMounted, ref } from "vue";
import type { IMethod } from "@/store/method";
import { useMethodStore } from "@/store/method";
import { useModalStore } from "@/gam-lib-ui/vue/store/modal";
import MethodTag from "@/component/MethodTag.vue";
import formatHighlight from "json-format-highlight";

// Stores
const methodStore = useMethodStore();
const modalStore = useModalStore();

// Vars
const controllerList = ref({} as Record<string, any>);
const selectedMethodUid = ref("");
const dialogVisible = ref(false);
const customColorOptions = ref({
  keyColor: "#be31ec",
  numberColor: "#77b0fc",
  stringColor: "#18a70e",
  trueColor: "#ff8080",
  falseColor: "#ff8080",
  nullColor: "#e54b4b",
});

// Hooks
onMounted(async () => {
  await methodStore.getList();
  for (let i = 0; i < methodStore.items.length; i++) {
    if (!controllerList.value[methodStore.items[i].controller])
      controllerList.value[methodStore.items[i].controller] = [];
    controllerList.value[methodStore.items[i].controller].push(
      methodStore.items[i]
    );
  }
});

// Methods
function rowClick(m: IMethod) {
  selectedMethodUid.value = m.uid;
}

function run(m: IMethod) {
  try {
    JSON.parse(methodStore.params[m.uid]);
  } catch (e) {
    methodStore.response[m.uid] = e;
    return;
  }

  if (m.inputMethod === "multipart") {
    const formData = new FormData();

    for (const key in methodStore.formData[m.uid]) {
      formData.set(key, methodStore.formData[m.uid][key]);
    }

    methodStore.run(m.uid, m.httpMethod, m.url, formData);
  } else {
    methodStore.run(
      m.uid,
      m.httpMethod,
      m.url,
      JSON.parse(methodStore.params[m.uid])
    );
  }
}

function fileSelect(e: any, name: string, destination: any) {
  console.log(e.target.files[0]);
  destination[name] = e.target.files[0];
}
</script>

<style module lang="scss">
.main {
  overflow-y: auto;
  height: 100%;
  padding: 10px;
  display: grid;
  grid-template-columns: 1fr 0.5fr;

  .response {
    padding: 10px;

    .header {
      height: 16px;
      display: flex;
    }
  }
}
</style>
