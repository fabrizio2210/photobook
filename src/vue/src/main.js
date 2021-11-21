import Vue from "vue";
import router from "./router";
import store from "./store/index.js";
import App from "./App.vue";
import './assets/style.css';

Vue.config.productionTip = false;

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount("#app");
