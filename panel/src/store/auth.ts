import { defineStore } from "pinia";
import Axios from "axios";
import { HOST } from "@/const";
import type { IMethod } from "@/store/method";

export type AuthState = {
  isAuth: boolean;
};

export const useAuthStore = defineStore({
  id: "test",
  state: () =>
    ({
      isAuth: false,
    } as AuthState),
  actions: {
    async auth(login: string, password: string) {
      let authKey = (
        await Axios.post(`${HOST}/debug/api/auth`, { login, password })
      ).data;
      localStorage.setItem("debugAuthKey", authKey);
    },
    async check(key: string) {
      this.isAuth = (
        await Axios.get(`${HOST}/debug/api/checkAuth?key=${key}`)
      ).data;
    },
  },
});
