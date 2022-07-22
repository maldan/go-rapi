import { defineStore } from "pinia";
import Axios from "axios";

export type LogState = {
  items: { kind: string; message: string; time: string }[];
};

export const useLogStore = defineStore({
  id: "log",
  state: () =>
    ({
      items: [],
    } as LogState),
  actions: {
    async getList(date: string) {
      this.items = (
        await Axios.get(`http://localhost:16000/debug/log/search?date=${date}`)
      ).data.response;
    },
  },
});
