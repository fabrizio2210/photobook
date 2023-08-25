<template>
  <div>
    <em v-if="photos.loading">Loading photos...</em>
    <img v-show="photos.loading" src="../assets/loading.gif" />
    <span v-if="photos.error" class="text-danger"
      >ERROR: {{ photos.error }}</span
    >
    <div class="photo-list" v-if="photos.photos_list">
      <div
        class="photo-element"
        v-for="photo in photos.photos_list"
        :key="photo.photo_id"
      >
        <img
          class="photo-img-element"
          loading="lazy"
          @error="imgError(photo.photo_id)"
          :src="photo.location"
        />
        <div class="photo-description">{{ photo.description }}</div>
        <div v-show="photo.author" class="photo-author">
          -- {{ photo.author }} --
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {};
  },
  computed: {
    user() {
      return this.$store.state.authentication.user;
    },
    photos() {
      return this.$store.state.photos.all;
    },
  },
  methods: {
    imgError(id) {
      this.$store.dispatch("photos/get", { id });
    }
  },
  mounted() {
    this.vueInsomnia().on();
  },
  destroyed() {
    this.vueInsomnia().off();
  }
};
</script>
