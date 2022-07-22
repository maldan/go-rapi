import { createApp } from "vue";
import { createPinia } from "pinia";
import "./main.scss";

import App from "./App.vue";
import router from "./router";
import UI from "@/gam-lib-ui/vue/ui";

const app = createApp(App);

app.use(createPinia());
app.use(router);
app.use(UI);

app.mount("#app");
