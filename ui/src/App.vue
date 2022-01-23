<template>
	{{ error ? 'Request failed: ' + error : undefined }}
	<InitializeKeystoreFullscreenModal :show="ready && onboarding"/>
	<div id="entries" v-if="keystores.length > 0 && !onboarding">
		ENTRIES CONTAINER
	</div>
</template>

<script lang="ts">
import {defineComponent} from 'vue'
import {mapState} from 'vuex'
import InitializeKeystoreFullscreenModal from './components/InitializeKeystoreFullscreenModal.vue'
import api from './api'

export default defineComponent({
	name: 'App',
	components: {InitializeKeystoreFullscreenModal},
	data: () => ({
		error: undefined,
		onboarding: false,
		ready: false,
		keystore: undefined,
		keystoreIds: [],
		keystores: []
	}),
	computed: {
		// ...mapState({
		// 	keystoreIds: (state) => state.keystore.keystoreIds,
		// 	keystore: (state) => state.keystore.keystore,
		// 	ready: (state) => state.ready,
		// }),
	},
	mounted() {
		api
			.get(`/keystores`)
			.then(response => {
				console.log(response);
				if (response.status === 200) {
					if (response.data.length == 0) {
						this.onboarding = true
					}

					this.ready = true
					// todo: show login form with just master password
				}
			});

		this.$store.dispatch("keystore/getKeystoreIds").then(() => {
			if (this.$store.state.keystore.keystoreIds.length == 0) {
				this.$store.commit("setOnboarding", true)
				this.$store.commit("setReady", true)
			}
		}).catch((err) => {
			console.log('err', err)
			this.error = err
		})
	},
	methods: {}
})
</script>

<style>
body {
	margin: 0;
}
</style>
