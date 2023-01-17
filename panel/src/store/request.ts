import { defineStore } from "pinia";
import Axios from "axios";
import { HOST } from "@/const";

export interface IRequest {
  id: string;
  httpMethod: string;
  url: string;
}

export type RequestState = {
  list: IRequest[];
};

export const useRequestStore = defineStore({
  id: "request",
  state: () =>
    ({
      list: [],
    } as RequestState),
  actions: {
    async getList(offset: number, limit: number) {
      this.list = (
        await Axios.get(
          `${HOST}/debug/api/requestList?offset=${offset}&limit=${limit}`
        )
      ).data;
    },
  },
});
