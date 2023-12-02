import { defineStore } from "pinia";
import Axios from "axios";
import { HOST } from "@/const";
import { escape } from "lodash";

export interface BackupState {
  taskList: { folder: string; name: string }[];
}

export const useBackupStore = defineStore({
  id: "backup",
  state: () =>
    ({
      taskList: [],
    } as BackupState),
  actions: {
    async getTaskList() {
      this.taskList = (await Axios.get(`${HOST}/debug/backup/taskList`)).data;
    },
  },
});
