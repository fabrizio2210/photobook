import { photoService } from "../services";
import Vue from "vue";
import router from '../router/';

export const photos = {
  namespaced: true,
  state: {
    all: { photos_dict : {}, photos_list : [] },
    status : {}
  },
  actions: {
    creating({ dispatch, commit }, { photoname }) {
      commit("creatingRequest", {});
        photoService.create(photoname).then(
          photo => {
            commit("creatingSuccess", photo);
            router.push({
              name: 'PhotoSettings',
              params: { photo_id: photo.id }
            });
          },
          error => {
            commit("creatingFailure", error);
            dispatch("alert/error", error, { root: true});
          }
        );
      },
    getAll({ commit }) {
      commit("getAllRequest");

      photoService.getAll().then(
        photos => commit("getAllSuccess", photos["photos"]),
        error => commit("getAllFailure", error)
      );
    },
    get({ commit }, { photo_id}) {
      commit("getRequest");

      photoService.get(photo_id).then(
        photo => commit("getSuccess", photo["photo"]),
        error => commit("getFailure", error)
      );
    }
  },
  mutations: {
    creatingRequest(state) {
      state.status = { creating: true };
    },
    creatingSuccess(state, photo) {
      state.status = { created: true };
      state.all.photos_dict[photo.id] = photo;
    },
    creatingFailure(state) {
      state.status = {};
      state.user = null;
    },
    getAllRequest(state) {
      state.all = { loading: true };
    },
    getRequest(state) {
      state.all.loading = true ;
    },
    getAllSuccess(state, photos) {
      //const photos_dict = {}
      //photos.forEach((element) => photos_dict[element.id] = element);
      //state.all = { photos_dict };
      const photos_list = photos.slice().reverse();
      state.all = { photos_list };
    },
    getSuccess(state, photo) {
      console.log(photo);
      Vue.set(state.all.photos_dict, photo.id, photo);
      Vue.delete(state.all, 'loading');
    },
    getAllFailure(state, error) {
      state.all = { error };
    },
    getFailure(state, error) {
      state.all = { error };
    }
  }
};
