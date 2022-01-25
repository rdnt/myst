<template>
  <div>
    {{ this.keystoreIds.length}}
    {{ error ? 'Request failed: ' + error : undefined }}
    <InitializeKeystoreFullscreenModal :show="ready && onboarding" />
    <div
      id="entries"
      v-if="keystores.length > 0 && !onboarding"
    >
      ENTRIES CONTAINER
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from "vue";
import { mapActions, mapMutations, mapState } from "vuex";
import api from "./api";
import InitializeKeystoreFullscreenModal from "./components/InitializeKeystoreFullscreenModal.vue";

export default defineComponent({
  name: "App",
  components: { InitializeKeystoreFullscreenModal },
  data: () => ({
    error: undefined,
    onboarding: false,
  }),
  computed: {
    ...mapState({
      keystoreIds: (state) => state.keystoreIds,
      keystore: (state) => state.keystore.keystore,
      keystores: (state) => state.keystore.keystores,
      ready: (state) => state.ready,
    }),
  },
  mounted() {
    api.get(`/keystores`).then((response) => {
      console.log(response);
      if (response.status === 200) {
        if (response.data.length == 0) {
          this.onboarding = true;
        }

        this.ready = true;
        // todo: show login form with just master password
      }
    });

    this.getKeystoreIds()
      .then(() => {
        if (this.keystoreIds.length == 0) {
          this.setOnboarding(true);
          this.setReady(true);
        }
      })
      .catch((err) => {
        console.log("err", err);
        this.error = err;
      });
  },
  methods: {
    ...mapActions({
      getKeystoreIds: "keystore/getKeystoreIds",
    }),
    ...mapMutations({
      setOnboarding: "setOnboarding",
      setReady: "setReady",
    }),
  },
});
</script>

<style>
body {
  margin: 0;
}
</style>
