<template>
  <div>
    <span>{{ error ? 'Request failed: ' + error : undefined }}</span>
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
import {defineComponent} from "vue";
import {mapActions, mapState, mapStores} from "pinia";
import {mainStore} from "./store";
import InitializeKeystoreFullscreenModal from "./components/InitializeKeystoreFullscreenModal.vue";
import api from "./api";

export default defineComponent({
	name: "App",
	components: {InitializeKeystoreFullscreenModal},
	data: () => ({
		error: undefined,
		onboarding: false,
		ready: false,
    keystore: undefined,
    keystoreIds: [],
    keystores: [],
	}),
	computed: {
		...mapStores(mainStore),
		...mapState(mainStore, ["ready", "onboarding"]),
		// ...mapState({
		//   // keystoreIds: (state) => state.keystore.keystoreIds,
		//   // keystore: (state) => state.keystore.keystore,
		//   ready: (state) => state.ready,
		// }),
	},
	created() {
		console.log('created')
		this.load()
	},
	mounted() {
		console.log(this.mainStore);
		// api.get(`/keystores`).then((response) => {
		// 	console.log(response);
		// 	if (response.status === 200) {
		// 		if (response.data.length == 0) {
		// 			this.onboarding = true;
		// 		}
		//
		// 		this.ready = true;
    //     // todo: show login form with just master password
    //   }
    // });

    // this.$store
    //   .dispatch("keystore/getKeystoreIds")
    //   .then(() => {
    //     if (this.$store.state.keystore.keystoreIds.length == 0) {
    //       this.$store.commit("setOnboarding", true);
    //       this.$store.commit("setReady", true);
    //     }
    //   })
    //   .catch((err) => {
    //     console.log("err", err);
    //     this.error = err;
    //   });
  },
  methods: {
		...mapActions(mainStore, ["load"]),
	},
});
</script>

<style>
body {
  margin: 0;
}
</style>
