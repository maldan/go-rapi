import {
  createRouter,
  createWebHashHistory,
  createWebHistory,
} from "vue-router";
import Main from "../page/Main.vue";

const router = createRouter({
  history: createWebHashHistory(),
  //history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: "/",
      name: "Main",
      component: Main,
    },
    /*{
      path: '/about',
      name: 'about',
      // route level code-splitting
      // this generates a separate chunk (About.[hash].js) for this route
      // which is lazy-loaded when the route is visited.
      component: () => import('../page/AboutView.vue')
    }*/
  ],
});

export default router;
