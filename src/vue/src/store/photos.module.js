import { photoService } from "../services";
import Vue from "vue";
import router from '../router/';

export const photos = {
  namespaced: true,
  state: {
    all: { photos_list : [] },
    last_timestamp: 0,
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
    getSince({ commit }, { last_timestamp }) {
      commit("getSinceRequest");

      photoService.getSince(last_timestamp).then(
        photos => commit("getSinceSuccess", photos["photos"]),
        error => commit("getSinceFailure", error)
      );
    }
  },
  mutations: {
    creatingRequest(state) {
      state.status = { creating: true };
    },
    creatingSuccess(state, photo) {
      state.status = { created: true };
      state.all.photos_list.push(photo);
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
    getSinceRequest(state) {
      state.all.loading = true ;
    },
    getAllSuccess(state, photos) {
      const photos_list = photos.slice().reverse();
      Vue.set(state, 'last_timestamp', photos_list[0].timestamp);
      state.all = { photos_list };
    },
    getSinceSuccess(state, photos) {
      console.log(photos);
      // Vue.Set???
      state.all.photos_list = photos.slice().reverse().concat(state.all.photos_list);
      if (state.all.photos_list.length > 0){
        state.last_timestamp = state.all.photos_list[0].timestamp;
      }
      Vue.delete(state.all, 'loading');
    },
    getAllFailure(state, error) {
      state.all = { error };
    },
    getSinceFailure(state, error) {
      state.all = { error };
    },
    getFailure(state, error) {
      state.all = { error };
    }
  }
};
