import { defineStore } from "pinia";
import Axios from "axios";
import type { IMethod } from "@/types";

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
        await Axios.get(`http://localhost:16000/debug/api/methodList`)
      ).data.response;
    },
  },
});
