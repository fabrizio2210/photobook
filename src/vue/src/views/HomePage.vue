<template>
    <div>
        <h3>Photos:</h3>
        <em v-if="photos.loading">Loading photos...</em>
        <img v-show="photos.loading" src="../assets/loading.gif" />
        <span v-if="photos.error" class="text-danger">ERROR: {{photos.error}}</span>
        <div class="photo-list" v-if="photos.photos_list">
            <div class="photo-element" v-for="photo in photos.photos_list" :key="photo.id">
                <img loading=lazy :src=photo.location />
                <div class="photo-description">{{photo.description}}</div>
                <div class="photo-author">{{photo.author}}</div>
            </div>
        </div>
    </div>
</template>

<script>
export default {
    data () {
        return {
            photoname: '',
            submitted: false
        }
    },
    computed: {
        user () {
            return this.$store.state.authentication.user;
        },
        photos () {
            return this.$store.state.photos.all;
        }
    },
    methods: {},
    created () {
        this.$store.dispatch('photos/getAll');
    }
};
</script>
