import Vue from "vue";
import Vuex from "vuex";

import { alert } from "./alert.module";
import { authentication } from "./authentication.module";
import { photos } from "./photos.module";

Vue.use(Vuex);

const store = new Vuex.Store({
  modules: {
    alert,
    authentication,
    photos
  }
});

export default store;
