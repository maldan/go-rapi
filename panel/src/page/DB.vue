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

    <!-- Pagination -->
    <el-pagination
      background
      layout="prev, pager, next"
      :total="dbStore.search.total"
      :page-size="dbStore.limit"
      style="margin-bottom: 10px; width: 100%"
      @current-change="changePage"
    />

    <el-input
      placeholder="Filter..."
      @change="refresh"
      v-model="dbStore.filter"
      style="margin-bottom: 10px"
    />

    <!-- Data -->
    <el-table
      :data="dbStore.search.result"
      stripe
      :border="true"
      style="width: 100%"
      :height="tableHeight"
      :empty-text="dbStore.error"
    >
      <!--      <el-table-column
        v-for="v in dbStore.settings.fieldList.filter((x) => !x.isHide)"
        :prop="v.name"
        :label="v.name"
      />-->

      <!-- Custom -->
      <el-table-column
        v-for="v in dbStore.settings.fieldList.filter((x) => !x.isHide)"
        :label="v.name"
      >
        <template #default="scope">
          <div v-if="v.type === 'bool'">
            <el-checkbox size="large" :model-value="scope.row[v.name]" />
          </div>
          <div v-else>{{ scope.row[v.name] }}</div>
        </template>
      </el-table-column>

      <!-- Edit -->
      <el-table-column
        v-if="dbStore.settings.isEditable"
        label="Edit"
        width="70"
      >
        <template #default="scope">
          <el-button @click="editById(scope.row.id)" style="width: 100%"
            ><el-icon><EditPen /></el-icon
          ></el-button>
        </template>
      </el-table-column>

      <!-- Delete -->
      <el-table-column
        v-if="dbStore.settings.isDeletable"
        label="Delete"
        width="70"
      >
        <template #default="scope">
          <el-button @click="deleteById(scope.row.id)" style="width: 100%"
            ><el-icon><DeleteFilled /></el-icon
          ></el-button>
        </template>
      </el-table-column>
    </el-table>

    <!-- Modal Edit -->
    <el-dialog v-model="dialogVisible" title="Edit" width="30%" draggable>
      <div
        v-for="v in dbStore.settings.fieldList"
        :key="v.name"
        style="margin-bottom: 10px"
      >
        <div v-if="v.type === 'bitmask'" style="margin-bottom: 5px">
          {{ v.name }}
        </div>
        <div v-else style="margin-bottom: 5px">{{ v.label || v.name }}</div>

        <!-- Type -->
        <el-input
          v-if="!v.isEdit"
          :disabled="!v.isEdit"
          :placeholder="v.name"
          v-model="tempRow[v.name]"
        />
        <el-input
          v-if="v.isEdit && v.type === 'string'"
          :placeholder="v.name"
          v-model="tempRow[v.name]"
        />
        <el-input-number
          v-if="v.isEdit && v.type === 'int'"
          :placeholder="v.name"
          v-model="tempRow[v.name]"
          style="width: 100%"
        />
        <el-checkbox
          v-if="v.isEdit && v.type === 'bool'"
          :placeholder="v.name"
          v-model="tempRow[v.name]"
        />
        <div v-if="v.isEdit && v.type === 'bitmask'">
          <el-checkbox
            :label="x"
            v-for="(x, i) in v.label.split(',')"
            :placeholder="x"
            @change="
              tempRow[v.name] = changeBitMask($event, tempRow[v.name], i)
            "
            v-model="tempRow[v.name + '_mask_' + i]"
          />
          <!--  :checked="(tempRow[v.name] & (1 << i)) === 1 << i" -->
          <div>
            Bitmask: {{ tempRow[v.name]?.toString(2) }} {{ tempRow[v.name] }}
          </div>
        </div>
      </div>

      <template #footer>
        <span class="dialog-footer">
          <el-button type="primary" @click="save"> Save </el-button>
        </span>
      </template>
    </el-dialog>
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
const editId = ref(0);
const dialogVisible = ref(false);
const tempRow = ref({});

// Hooks
onMounted(async () => {
  tableHeight.value = window.innerHeight - 205;

  await dbStore.getTableList();
  dbStore.table = dbStore.tableList[0];
  await dbStore.getSettings();
  await dbStore.getSearch();
});

// Methods
async function refresh() {
  await dbStore.getSearch();
}

async function deleteById(id: string) {
  if (confirm("Sure?")) {
    await dbStore.deleteById(id);
    await dbStore.getSearch();
  }
}

async function editById(id: number) {
  editId.value = id;
  tempRow.value = {};
  tempRow.value = await dbStore.getById(id);

  for (let i = 0; i < dbStore.settings.fieldList.length; i++) {
    const field = dbStore.settings.fieldList[i];
    if (field.type === "bitmask") {
      for (let j = 0; j < 64; j++) {
        // @ts-ignore
        tempRow.value[field.name + "_mask_" + j] =
          (tempRow.value[field.name] & (1 << j)) === 1 << j;
      }
    }
  }

  console.log(tempRow.value);

  dialogVisible.value = true;
}

async function changeTab(tab: string) {
  dbStore.table = tab;
  dbStore.search.result = [];
  await dbStore.getSettings();
  await dbStore.getSearch();
}

async function changePage(page: number) {
  dbStore.offset = (page - 1) * dbStore.limit;
  dbStore.search.result = [];
  await dbStore.getSearch();
}

async function save() {
  await dbStore.update(editId.value, tempRow.value);
  await dbStore.getSearch();
  dialogVisible.value = false;
}

function changeBitMask(isSet: boolean, current: number, pos: number): number {
  if (isSet) {
    return current | (1 << pos);
  }
  return current & ~(1 << pos);
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
