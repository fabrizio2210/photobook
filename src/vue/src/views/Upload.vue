<template>
  <div class="upload">
    <h1>This is a upload page</h1>
      <ul>
        <li v-for="file in files" :key="file.id">
          <span>{{file.name}}</span> -
          <span>{{formatSize(file.size)}}</span> -
          <span v-if="file.error">{{file.error}}</span>
          <span v-else-if="file.success">success</span>
          <span v-else-if="file.active">active</span>
          <span v-else-if="!!file.error">{{file.error}}</span>
          <span v-else></span>
        </li>
      </ul>
    <div class="form-group">
	<label for="description">Descrizione</label>
	<input type="text" v-model="description" name="description" class="form-description" />
    </div>
    <div class="form-group">
	<label htmlFor="author">Il tuo nome</label>
	<input type="text" v-model="author" name="author" class="form-author" />
    </div>
    <UploadImage 
        class="btn btn-primary"
        post-action="/api/new_photo"
        :data="{author_id: '1234', description: description, author: author}"
        extensions="gif,jpg,jpeg,png,webp"
        accept="image/png,image/gif,image/jpeg,image/webp"
        :multiple="false"
        :size="1024 * 1024 * 10"
        ref="upload"
        v-model="files" >
       <a class="btn" href="#">Select an image</a>
    </UploadImage>
  <button v-show="!$refs.upload || !$refs.upload.active" @click.prevent="$refs.upload.active = true" type="button">Start upload</button>
  <button v-show="$refs.upload && $refs.upload.active" @click.prevent="$refs.upload.active = false" type="button">Stop upload</button>
  </div>
</template>

<script>
// @ is an alias to /src
import UploadImage from 'vue-upload-component';

export default {
  name: "UploadPage",
  components: {
    UploadImage
  },
  data () {
    return {
      description: '',
      author: '',
      upload: {},
      files: [],
    }
  },
  methods: {
    formatSize (size) {
       if (size > 1024 * 1024 * 1024 * 1024) {
         return (size / 1024 / 1024 / 1024 / 1024).toFixed(2) + ' TB'
       } else if (size > 1024 * 1024 * 1024) {
         return (size / 1024 / 1024 / 1024).toFixed(2) + ' GB'
       } else if (size > 1024 * 1024) {
         return (size / 1024 / 1024).toFixed(2) + ' MB'
       } else if (size > 1024) {
         return (size / 1024).toFixed(2) + ' KB'
       }
       return size.toString() + ' B'
    }
  }
};
</script>
