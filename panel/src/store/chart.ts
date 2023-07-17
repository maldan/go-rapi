import { defineStore } from "pinia";
import Axios from "axios";
import { HOST } from "@/const";
import { escape } from "lodash";

export interface ChartState {
  chartList: { folder: string; name: string }[];
  response: Record<string, any>;
  status: Record<string, boolean>;
}

export const useChartStore = defineStore({
  id: "chart",
  state: () =>
    ({
      chartList: [],
      response: {},
      status: {},
    } as ChartState),
  actions: {
    async getCommandList() {
      this.chartList = (await Axios.get(`${HOST}/debug/chart/list`)).data;
    },
    async executeFolder(folder: string) {
      let list = this.chartList.filter((x) => x.folder === folder);
      for (let i = 0; i < list.length; i++) {
        await this.execute(list[i].folder, list[i].name);
      }
    },
    async execute(folder: string, name: string) {
      this.response[folder + "/" + name] = [];
      this.status[folder + "/" + name] = true;
      this.response[folder + "/" + name] = (
        await Axios.post(
          `${HOST}/debug/chart/execute?folder=${folder}&name=${name}`
        )
      ).data;
      this.status[folder + "/" + name] = false;
    },
  },
});
