import { createWebHistory, createRouter } from "vue-router";
import Home from "./pages/Home.vue";
import Config from "./pages/Config.vue";

const routes = [
  {
    path: "/",
    name: "Home",
    component: Home,
  },
  {
    path: "/config",
    name: "Config",
    component: Config,
  }
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;