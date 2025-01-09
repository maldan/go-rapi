<template>
  <div :class="$style.main">
    <!-- Label -->
    <div v-if="info.type === 'bitmask'" style="margin-bottom: 5px">
      {{ info.name }}
    </div>
    <div v-else style="margin-bottom: 5px">{{ info.label || info.name }}</div>

    <!-- Drop down -->
    <el-select
      v-if="info.type === 'string' && info.optionList"
      :placeholder="info.name"
      v-model="tempRow[info.name]"
      :disabled="!info.isEdit"
      style="width: 100%"
    >
      <el-option v-for="x in info.optionList" :key="x" :value="x">{{
        x
      }}</el-option>
    </el-select>

    <el-input
      v-if="info.type === 'string' && !info.optionList"
      :placeholder="info.name"
      v-model="tempRow[info.name]"
      :disabled="!info.isEdit"
      :type="info.isTextarea ? 'textarea' : 'text'"
      :input-style="{ height: getHeight(info.height) }"
    />
    <el-input-number
      v-if="info.type === 'int' || info.type === 'float'"
      :placeholder="info.name"
      v-model="tempRow[info.name]"
      style="width: 100%"
      :disabled="!info.isEdit"
    />
    <el-checkbox
      v-if="info.type === 'bool'"
      :placeholder="info.name"
      v-model="tempRow[info.name]"
      :disabled="!info.isEdit"
    />
    <el-date-picker
      v-if="info.type === 'date'"
      v-model="tempRow[info.name]"
      type="date"
      :placeholder="info.name"
      :disabled="!info.isEdit"
      style="width: 100%"
    />
    <el-date-picker
      v-if="info.type === 'datetime'"
      v-model="tempRow[info.name]"
      type="datetime"
      :placeholder="info.name"
      :disabled="!info.isEdit"
      style="width: 100%"
    />

    <!-- Bitmask -->
    <div v-if="info.type === 'bitmask'">
      <el-checkbox
        :label="x"
        v-for="(x, i) in info.label.split(',')"
        :placeholder="x"
        @change="
          tempRow[info.name] = changeBitMask($event, tempRow[info.name], i)
        "
        :disabled="!info.isEdit"
        v-model="tempRow[info.name + '_mask_' + i]"
      />
    </div>

    <!-- DataUrl -->
    <div v-if="info.type === 'dataUrl'">
      <el-upload
        ref="upload"
        :limit="1"
        :auto-upload="false"
        :on-exceed="handleExceed"
        :on-change="
          (e) => {
            fileSelected(e, tempRow, info.name);
          }
        "
        :disabled="!info.isEdit"
      >
        <img
          v-if="info.name === 'image' || info.name === 'photo'"
          :src="tempRow[info.name]"
          alt=""
          style="max-width: 128px; max-height: 128px"
        />
        <template #trigger>
          <el-button type="primary" :disabled="!info.isEdit"
            >Select {{ info.name }}</el-button
          >
        </template>
      </el-upload>
    </div>

    <!-- File -->
    <div v-if="info.type === 'file'">
      <!-- <el-upload
        ref="upload"
        :limit="1"
        :auto-upload="false"
        :on-exceed="handleExceed"
        :on-change="
          (e) => {
            fileSelected(e, tempRow, info.name, true);
          }
        "
        :disabled="!info.isEdit"
      >
        <template #trigger>
          <el-button type="primary" :disabled="!info.isEdit"
            >Select {{ info.name }}</el-button
          >
        </template>
      </el-upload> -->
      <input
        type="file"
        @change="
          (e) => {
            fileSelected(e, tempRow, info.name, true);
          }
        "
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import dayjs from "dayjs";
import { ref } from "vue";
import type { UploadInstance, UploadProps, UploadRawFile } from "element-plus";
const upload = ref<UploadInstance>();

const props = defineProps<{
  info: {
    label: string;
    isEdit: boolean;
    isTextarea: boolean;
    type: string;
    name: string;
    height: string;
    optionList: any[];
  };
  tempRow: any;
}>();

// const tempRow = ref({});
const handleExceed: UploadProps["onExceed"] = (files) => {
  upload.value!.clearFiles();
  const file = files[0] as UploadRawFile;
  file.uid = Math.random();
  upload.value!.handleStart(file);
};

function changeBitMask(isSet: boolean, current: number, pos: number): number {
  if (isSet) {
    return current | (1 << pos);
  }
  return current & ~(1 << pos);
}

function fileSelected(e: any, out: any, name: string, isReal: boolean) {
  if (isReal) {
    console.log(e);
    out[name] = e.target.files[0];
    //console.log(out[name]);
  } else {
    const reader = new FileReader();
    reader.addEventListener(
      "load",
      () => {
        out[name] = reader.result;
      },
      false
    );

    reader.readAsDataURL(e.raw);
  }
}

function getHeight(h: string) {
  if (!h || h === "0") return "auto";
  if (!Number(h)) return h;
  return h + "px";
}
</script>

<style lang="scss" module>
.main {
  // border: 2px solid #fe0000;
  // border: 1px dashed rgba(255, 255, 255, 0.3);
  border-bottom: 1px dashed rgba(255, 255, 255, 0.15);
  padding-bottom: 10px;
  // border-radius: 4px;
}
</style>
