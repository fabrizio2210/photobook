<template>
  <div>
    <div class="alert alert-info">
      Carica la tua foto con la descrizione.
    </div>
    <form @submit.prevent="handleSubmit">
      <div class="form-group">
        <label for="description">Descrizione</label>
        <input
          type="text"
          v-model="description"
          name="description"
          class="form-control"
        />
      </div>
      <div class="form-group">
        <label htmlFor="author">Il tuo nome</label>
        <input
          type="text"
          v-model="author"
          name="author"
          class="form-control"
        />
      </div>
      <div class="form-group">
        <button class="btn btn-primary" :disabled="loggingIn">Login</button>
        <img v-show="loggingIn" src="../assets/loading.gif" />
      </div>
    </form>
  </div>
</template>

<script>
export default {
  data() {
    return {
      username: "",
      password: "",
      submitted: false
    };
  },
  computed: {
    loggingIn() {
      return this.$store.state.authentication.status.loggingIn;
    }
  },
  created() {
    // reset login status
    this.$store.dispatch("authentication/logout");
  },
  methods: {
    handleSubmit() {
      this.submitted = true;
      const { username, password } = this;
      const { dispatch } = this.$store;
      if (username && password) {
        dispatch("authentication/login", { username, password });
      }
    }
  }
};
</script>
