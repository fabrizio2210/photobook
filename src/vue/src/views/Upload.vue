<template>
  <div class="upload">
    <h2>Upload here your photo</h2>
    <div class="buttons-upload">
      <UploadImage
        class="btn-primary"
        :post-action="'/api/new_photo?ticket_id=' + ticket_id"
        extensions="gif,jpg,jpeg,png,webp"
        accept="image/png,image/gif,image/jpeg,image/webp"
        :multiple="false"
        :size="1024 * 1024 * max_size"
        ref="upload"
        @input-filter="inputFilter"
        @input-file="inputFile"
        v-model="files"
      >
        <button v-if="!files.length" class="btn">Select an image</button>
        <div v-for="file in files" :key="file.id">
          <div>
            <img class="image-thumb" v-if="file.thumb" :src="file.thumb" />
          </div>
          <div>{{ formatSize(file.size) }}</div>
          <div class="transfer-status-error" v-if="status.error">
            {{ status.error }}
          </div>
          <div class="transfer-status-error" v-else-if="file.error">
            {{ errorMessage(file) }}
          </div>
          <div class="transfer-status-complete" v-else-if="file.success">
            <div v-if="show_progress">
              done, the picture will be published in a few seconds. Click on the
              image to change it
            </div>
          </div>
          <div class="transfer-status" v-else-if="file.active">
            <div v-if="show_progress">
              {{ file.progress }}% transfered
              <img src="../assets/loading.gif" />
            </div>
          </div>
          <div class="transfer-status" v-else>
            click on the image to change it
          </div>
        </div>
      </UploadImage>
    </div>
    <div class="text-input-group">
      <div class="form-group">
        <label
          >Description:
          <textarea
            cols="40"
            rows="5"
            maxlength="200"
            v-model="description"
            class="form-description"
          />
        </label>
      </div>
      <div class="form-group">
        <label
          >Your name:
          <input
            type="text"
            v-model="author"
            maxlength="25"
            class="form-author"
          />
        </label>
      </div>
    </div>
    <div>
      <button
        v-show="
          (!$refs.upload || !$refs.upload.active || !show_progress) &&
            files.length
        "
        @click.prevent="uploadMetadata()"
        class="btn"
        type="btn"
      >
        Upload
      </button>
      <button
        v-show="$refs.upload && $refs.upload.active && show_progress"
        @click.prevent="$refs.upload.active = false"
        class="btn-stop"
        type="btn"
      >
        Stop
      </button>
    </div>
  </div>
</template>

<script>
// @ is an alias to /src
import UploadImage from "vue-upload-component";

export default {
  name: "UploadPage",
  components: {
    UploadImage
  },
  data() {
    return {
      max_size: 20,
      show_progress: false,
      description: "",
      author: "",
      files: []
    };
  },
  watch: {
    files() {
      this.resetError();
    }
  },
  computed: {
    uploaded() {
      if (this.files.length > 0) {
        return this.files[0].success;
      }
      return false;
    },
    stored_author() {
      if (this.$store.state.authentication.user) {
        if (this.$store.state.authentication.user.name) {
          return this.$store.state.authentication.user.name;
        }
      }
      return "";
    },
    status() {
      return this.$store.state.photos.status;
    },
    ticket_id() {
      if (this.$store.state.photos) {
        if (this.$store.state.photos.ticket_id) {
          return this.$store.state.photos.ticket_id;
        }
      }
      return "12345";
    },
    uid() {
      if (this.$store.state.authentication.user) {
        if (this.$store.state.authentication.user.uid) {
          return this.$store.state.authentication.user.uid;
        }
      }
      return "12345";
    }
  },

  methods: {
    uploadMetadata() {
      const author_id = this.uid;
      const author = this.$sanitize(this.author);
      const description = this.$sanitize(this.description);
      const ticket_id = this.ticket_id;
      this.$store.dispatch("photos/putMetadata", {
        author_id,
        author,
        description,
        ticket_id
      });
      this.show_progress = true;
    },
    inputFile(newFile, oldFile) {
      // Automatic upload
      if (
        Boolean(newFile) !== Boolean(oldFile) ||
        oldFile.error !== newFile.error
      ) {
        if (!this.$refs.upload.active) {
          var vm = this;
          if (vm.files.length > 0) {
            vm.files[0].data = {
              author_id: vm.uid
            };
            this.$refs.upload.active = true;
          }
        }
      }
    },
    inputFilter(newFile, oldFile) {
      if (
        newFile &&
        newFile.error === "" &&
        newFile.file &&
        (!oldFile || newFile.file !== oldFile.file)
      ) {
        this.show_progress = false;
        // Create a blob field
        newFile.blob = "";
        let URL = window.URL || window.webkitURL;
        if (URL) {
          newFile.blob = URL.createObjectURL(newFile.file);
        }
        // Thumbnails
        newFile.thumb = "";
        if (newFile.blob && newFile.type.substr(0, 6) === "image/") {
          newFile.thumb = newFile.blob;
        }
      }
    },
    errorMessage(file) {
      var msg = file.error;
      if (file.error == "size") {
        msg =
          "The image is too big to be upload, try with an image smaller than " +
          this.max_size +
          "MB.";
      }
      if (file.error == "server") {
        msg =
          "The server could not handle the request due to an unexpected error.";
      }
      if (!this.isObjectEmpty(file.response)) {
        msg = file.response.message;
        if (!msg) {
          msg = "Error, it was not possible to upload.";
        }
      }
      return msg;
    },
    isObjectEmpty(objectName) {
      return Object.keys(objectName).length === 0;
    },
    formatSize(size) {
      if (size > 1024 * 1024 * 1024 * 1024) {
        return (size / 1024 / 1024 / 1024 / 1024).toFixed(2) + " TB";
      } else if (size > 1024 * 1024 * 1024) {
        return (size / 1024 / 1024 / 1024).toFixed(2) + " GB";
      } else if (size > 1024 * 1024) {
        return (size / 1024 / 1024).toFixed(2) + " MB";
      } else if (size > 1024) {
        return (size / 1024).toFixed(2) + " KB";
      }
      return size.toString() + " B";
    },
    resetError() {
      this.$store.dispatch("photos/resetError");
    }
  },
  mounted() {
    this.author = this.stored_author;
    this.$store.dispatch("photos/getTicket");
  },
  destroyed() {
    const author = this.author;
    this.$store.dispatch("authentication/set_name", author);
  }
};
</script>
