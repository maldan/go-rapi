import { createRouter, createWebHashHistory } from "vue-router";
import Logs from "../page/Logs.vue";
import Methods from "../page/Methods.vue";
import Settings from "../page/Settings.vue";

const router = createRouter({
  history: createWebHashHistory(),
  routes: [
    {
      path: "/",
      name: "Logs",
      component: Logs,
    },
    {
      path: "/methods",
      name: "Methods",
      component: Methods,
    },
    {
      path: "/settings",
      name: "Settings",
      component: Settings,
    },
  ],
});

export default router;
