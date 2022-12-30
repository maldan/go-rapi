import { createApp } from "vue";
import ElementPlus from "element-plus";
import { createPinia } from "pinia";
import "element-plus/dist/index.css";
import "element-plus/theme-chalk/dark/css-vars.css";
import "./main.scss";

import App from "./App.vue";
import router from "./router";
import UI from "@/gam-lib-ui/vue/ui";

const app = createApp(App);

app.use(createPinia());
app.use(router);
app.use(UI);
app.use(ElementPlus);

app.mount("#app");
