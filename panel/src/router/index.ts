import { createRouter, createWebHashHistory } from "vue-router";
import Logs from "../page/Logs.vue";
import Methods from "../page/Methods.vue";
import Settings from "../page/Settings.vue";
import Tests from "../page/Tests.vue";

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
    {
      path: "/tests",
      name: "Tests",
      component: Tests,
    },
  ],
});

export default router;
