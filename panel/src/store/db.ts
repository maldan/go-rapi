import { defineStore } from "pinia";
import Axios from "axios";
import { HOST } from "@/const";

export type DbState = {
  tableList: string[];
  struct: Record<string, string>;
  table: string;
  list: any[];
  offset: number;
  limit: number;
};

export const useDbStore = defineStore({
  id: "db",
  state: () =>
    ({
      tableList: [],
      struct: {},
      table: "",
      list: [],
      offset: 0,
      limit: 20,
    } as DbState),
  actions: {
    async getTableList() {
      this.tableList = (await Axios.get(`${HOST}/debug/data/tableList`)).data;
    },
    async getStruct() {
      this.struct = (
        await Axios.get(`${HOST}/debug/data/struct?table=${this.table}`)
      ).data;
    },
    async getSearch() {
      this.list = (
        await Axios.get(
          `${HOST}/debug/data/search?table=${this.table}&offset=${this.offset}&limit=${this.limit}`
        )
      ).data;
    },
    async deleteById(id: string) {
      this.list = (
        await Axios.delete(
          `${HOST}/debug/data/byId?table=${this.table}&id=${id}`
        )
      ).data;
    },
  },
});
