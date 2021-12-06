<template>
    <div>
        <em v-if="photos.loading">Loading photos...</em>
        <img v-show="photos.loading" src="../assets/loading.gif" />
        <span v-if="photos.error" class="text-danger">ERROR: {{photos.error}}</span>
        <div class="photo-list" v-if="photos.photos_list">
            <div class="photo-element" v-for="photo in photos.photos_list" :key="photo.id">
                <img class="photo-img-element" loading=lazy :src=photo.location />
                <div class="photo-description">{{photo.description}}</div>
                <div v-show="photo.author" class="photo-author">-- {{photo.author}} --</div>
            </div>
        </div>
    </div>
</template>

<script>

export default {
    data () {
        return { }
    },
    computed: {
        photos () {
            return this.$store.state.photos.my;
        },
        last_timestamp () {
            return this.$store.state.photos.last_timestamp;
        },
        uid () {
          if (this.$store.state.authentication.user) {
            if (this.$store.state.authentication.user.uid) {
              return this.$store.state.authentication.user.uid;
            }
          }
          return "12345";
        },
    },
    methods: {
       populatePhotos(uid) {
         if (typeof uid !== 'undefined') {
           this.$store.dispatch('photos/getOwn', {uid});
         }
       },
    },
    mounted () {
      this.populatePhotos(this.uid);
    },
    destroyed () {
    }
};
</script>
