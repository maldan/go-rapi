import { defineStore } from "pinia";
import Axios from "axios";
import type { IMethod } from "@/types";
import { HOST } from "@/const";

export type MethodState = {
  items: IMethod[];
};

export const useMethodStore = defineStore({
  id: "method",
  state: () =>
    ({
      items: [],
    } as MethodState),
  actions: {
    async getList() {
      this.items = (
        await Axios.get(`${HOST}/debug/api/methodList`)
      ).data.response;
    },
  },
});
