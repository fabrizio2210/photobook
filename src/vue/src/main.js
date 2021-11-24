import Vue from "vue";
import router from "./router";
import store from "./store/index.js";
import App from "./App.vue";
import VueSSE from "vue-sse";
import vueInsomnia from "vue-insomnia";
import VueSanitize from "vue-sanitize";
import './assets/style.css';

Vue.use(vueInsomnia);
Vue.use(VueSanitize);
Vue.use(VueSSE);

Vue.config.productionTip = false;

new Vue({
  router,
  store,
  render: h => h(App)
}).$mount("#app");
