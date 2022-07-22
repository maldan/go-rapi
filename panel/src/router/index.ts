import { createRouter, createWebHashHistory } from "vue-router";
import Logs from "../page/Logs.vue";
import Methods from "../page/Methods.vue";
import Test from "../page/Test.vue";

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
      path: "/test",
      name: "Test",
      component: Test,
    },
  ],
});

export default router;
