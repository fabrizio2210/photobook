import Vue from "vue";
import VueRouter from "vue-router";

import HomePage from "../views/HomePage.vue";
import Upload from "../views/Upload.vue";
import Edit from "../views/Edit.vue";

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
  // otherwise redirect to home
  { path: "*", redirect: "/" }
];

const router = new VueRouter({
  mode: "history",
  routes
});

export default router;
