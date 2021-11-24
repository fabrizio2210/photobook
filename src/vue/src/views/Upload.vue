<template>
  <div class="upload">
    <h1>Carica qui la tua foto</h1>
    <div class="buttons-upload">
      <UploadImage 
          class="btn-primary"
          post-action="/api/new_photo"
          :data="{author_id: '1234', description: this.$sanitize(description), author: this.$sanitize(author)}"
          extensions="gif,jpg,jpeg,png,webp"
          accept="image/png,image/gif,image/jpeg,image/webp"
          :multiple="false"
          :size="1024 * 1024 * 10"
          ref="upload"
          @input-filter="inputFilter"
          v-model="files" >
         <button v-if="!files.length" class="btn">Seleziona un'immagine</button>
         <div v-for="file in files" :key="file.id">
           <div><img class="image-thumb" v-if="file.thumb" :src="file.thumb" /></div>
           <div>{{formatSize(file.size)}}</div>
           <div v-if="file.error">{{file.error}}</div>
           <div v-else-if="file.success">fatto, clicca sull'immagine per cambiare</div>
           <div v-else-if="file.active">trasferimento</div>
           <div v-else-if="!!file.error">{{file.error}}</div>
           <div v-else>clicca sull'immagine per cambiare</div>
         </div>
      </UploadImage>
    </div>
    <div class="text-input-group">
      <div class="form-group">
          <label for="description">Descrizione:</label>
          <textarea cols="40" rows="5" v-model="description" name="description" class="form-description"/>
      </div>
      <div class="form-group">
          <label for="author">Il tuo nome:</label>
          <input type="text" v-model="author" name="author" class="form-author" />
      </div>
      </div>
    <div>
      <button v-show="(!$refs.upload || !$refs.upload.active) && files.length" @click.prevent="$refs.upload.active = true" class="btn" type="btn">Carica</button>
      <button v-show="$refs.upload && $refs.upload.active" @click.prevent="$refs.upload.active = false" class="btn-stop" type="btn">Ferma</button>
    </div>
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
    inputFilter(newFile, oldFile) {
    
      if (newFile && newFile.error === "" && newFile.file && (!oldFile || newFile.file !== oldFile.file)) {
        // Create a blob field
        newFile.blob = ''
        let URL = (window.URL || window.webkitURL)
        if (URL) {
          newFile.blob = URL.createObjectURL(newFile.file)
        }
        // Thumbnails
        newFile.thumb = ''
        if (newFile.blob && newFile.type.substr(0, 6) === 'image/') {
          newFile.thumb = newFile.blob
        }
      }
    },
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
