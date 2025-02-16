import { defineStore } from "pinia";
import Axios from "axios";
import { HOST } from "@/const";
import { escape } from "lodash";

export interface DbState {
  tableList: string[];
  table: string;
  search: {
    total: number;
    count: number;
    page: number;
    result: any[];
  };
  settings: {
    isEditable: boolean;
    isDeletable: boolean;
    fieldList: {
      name: string;
      isEdit: boolean;
      isHide: boolean;
      hasFilter: boolean;
      type: string;
      width: string;
      height: string;
    }[];
  };
  filter: Record<string, string>;
  offset: number;
  limit: number;
  error: string;
}

export const useDbStore = defineStore({
  id: "db",
  state: () =>
    ({
      tableList: [],
      settings: {
        isEditable: false,
        isDeletable: false,
        fieldList: [],
        editFieldList: [],
      },
      table: "",
      search: {
        result: [],
        total: 0,
        page: 0,
        count: 0,
      },
      filter: {},
      offset: 0,
      limit: 20,
      error: "",
    } as DbState),
  actions: {
    async getTableList() {
      this.tableList = (await Axios.get(`${HOST}/debug/data/tableList`)).data;
    },
    async getSettings() {
      this.settings = (
        await Axios.get(`${HOST}/debug/data/settings?table=${this.table}`)
      ).data;
    },
    async getById(id: number) {
      return (
        await Axios.get(`${HOST}/debug/data/byId?table=${this.table}&id=${id}`)
      ).data;
    },
    async create(value: any) {
      let hasFile = Object.values(value).filter((x) => x instanceof File);

      if (hasFile) {
        let formData = new FormData();
        let data: any = {};
        let id = 0;
        for (let key in value) {
          if (value[key] instanceof File) {
            formData.append("file" + id, value[key], value[key].name);
            data[key] = `%%file${id}%%`;
            id += 1;
          } else {
            data[key] = value[key];
          }
        }

        formData.set("data", JSON.stringify(data));

        return (
          await Axios.post(
            `${HOST}/debug/data/create?table=${this.table}`,
            formData
          )
        ).data;
      } else {
        return (
          await Axios.post(`${HOST}/debug/data/create?table=${this.table}`, {
            data: JSON.stringify(value),
          })
        ).data;
      }
    },
    async update(id: number, value: any) {
      return (
        await Axios.post(
          `${HOST}/debug/data/byId?table=${this.table}&id=${id}`,
          {
            data: JSON.stringify(value),
          }
        )
      ).data;
    },
    async getSearch() {
      this.error = "";
      this.search.result = [];

      try {
        this.search = (
          await Axios.get(
            `${HOST}/debug/data/search?table=${this.table}&filter=${btoa(
              JSON.stringify(this.filter)
            )}&offset=${this.offset}&limit=${this.limit}`
          )
        ).data;
      } catch (e: any) {
        this.error = e.response.data.description;
        this.search.total = 0;
      }
    },
    async exportData() {
      try {
        const exportData = (
          await Axios.get(
            `${HOST}/debug/data/export?table=${this.table}&filter=${btoa(
              JSON.stringify(this.filter)
            )}&offset=${this.offset}&limit=${this.limit}`
          )
        ).data;
        let out = Object.keys(exportData[0]).join(",") + "\n";
        for (let i = 0; i < exportData.length; i++) {
          out +=
            Object.values(exportData[i])
              .map((x) => `"${x}"`)
              .join(",") + "\n";
        }

        const file = new Blob([out], { type: "text" });
        const a = document.createElement("a"),
          url = URL.createObjectURL(file);
        a.href = url;
        a.download = this.table + ".csv";
        document.body.appendChild(a);
        a.click();
        setTimeout(function () {
          document.body.removeChild(a);
          window.URL.revokeObjectURL(url);
        }, 0);
      } catch (e: any) {
        console.error(e);
      }
    },
    async deleteById(id: string) {
      await Axios.delete(
        `${HOST}/debug/data/byId?table=${this.table}&id=${id}`
      );
    },
  },
});
