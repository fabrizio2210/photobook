import Vue from "vue";
import Vuex from "vuex";

import { alert } from "./alert.module";
import { authentication } from "./authentication.module";
import { photos } from "./photos.module";
import { admin } from "./admin.module";

Vue.use(Vuex);

const store = new Vuex.Store({
  modules: {
    admin,
    alert,
    authentication,
    photos
  }
});

export default store;
