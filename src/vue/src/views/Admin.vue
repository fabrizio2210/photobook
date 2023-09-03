<template>
  <div>
    <em v-if="admin.loading">Loading...</em>
    <img v-show="admin.loading" src="../assets/loading.gif" />
    <span v-if="admin.error" class="text-danger"
      >ERROR: {{ admin.error }}</span
    >
    <div class="admin-grid" >
      <div class="admin-row" >
        <div class="btn-container">
          <button
            class="admin-col-btn btn-print"
            @click.prevent="askPrint(uid)"
          >
            Print
          </button>
        </div>
      </div>
      <div class="admin-row" >
        <h3 class="label-btn">Is upload blocked:</h3>
        <Switchbox
          class="switchbox"
          checkboxId="a"
          :checked="admin.upload"
          @changeCheck="toggleUpload(uid)"
        />
      </div>
    </div>
  </div>
</template>

<script>
import Switch from "@/components/Switch.vue";
export default {
  components: {
    Switchbox: Switch
  },
  data() {
    return {};
  },
  computed: {
    admin() {
      return this.$store.state.admin;
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
    toggleUpload(uid) {
      this.$store.dispatch("admin/toggleUpload", { uid });
    },
    askPrint(uid) {
      this.$store.dispatch("admin/askPrint", { uid });
    }
  },
  mounted() {
    this.$store.dispatch("admin/getUpload");
  },
  destroyed() {}
};
</script>
