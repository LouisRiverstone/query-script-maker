import { createWebHistory, createRouter } from "vue-router";
import Home from "./pages/Home.vue";
import Config from "./pages/Config.vue";
import DatabaseBrowser from "./pages/DatabaseBrowser.vue";

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
  },
  {
    path: "/database-browser",
    name: "DatabaseBrowser",
    component: DatabaseBrowser,
  }
];

const router = createRouter({
  history: createWebHistory(),
  routes,
});

export default router;