import Vue from "vue";
import VueRouter from "vue-router";

import HomePage from "../views/HomePage.vue";
import Upload from "../views/Upload.vue";
import Edit from "../views/Edit.vue";
import Admin from "../views/Admin.vue";

Vue.use(VueRouter);

const routes = [
  {
    path: "/",
    name: "Home",
    component: HomePage
  },
  {
    path: "/upload",
    name: "Upload",
    component: Upload
  },
  {
    path: "/edit",
    name: "Edit",
    component: Edit
  },
  {
    path: "/admin",
    name: "Admin",
    component: Admin
  },
  // otherwise redirect to home
  { path: "*", redirect: "/" }
];

const router = new VueRouter({
  mode: "history",
  routes
});

export default router;
