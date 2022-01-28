<template>
  <div>
		<span>{{!ready ? 'Loading...' : ''}}</span>
		<span>{{ error ? 'Request failed: ' + JSON.stringify(error) : undefined }}</span>
    <InitializeKeystoreFullscreenModal v-if="onboarding" @created="keystoreCreated($event)" />
		<Login v-if="login" @login="loggedIn()" />
		<div
			id="entries"
			v-if="keystores.length > 0"
		>
			KEYSTORES {{JSON.stringify(keystores)}}
		</div>
  </div>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import InitializeKeystoreFullscreenModal from "./components/InitializeKeystoreFullscreenModal.vue";
import Login from "./components/LoginForm.vue";
import api from "./api";
import {Keystore} from "./api/generated";

export default defineComponent({
	name: "App",
	components: {InitializeKeystoreFullscreenModal, Login},
	data(): {
		onboarding: boolean,
		login: boolean,
		error: any,
		ready: boolean,
		keystore: Keystore | undefined,
		keystores: Keystore[],
	} {
		return {
			onboarding: false,
				login: false,
				error: undefined,
				ready: false,
				keystore: undefined,
				keystores: [],
		}
	},
	mounted() {
		this.ready = false
		this.init()
  },
  methods: {
		init() {
			api.keystoresRaw().then((resp) => {
				if (resp.raw.status === 200) {
					resp.value().then(keystores => {
						this.onboarding = keystores.length == 0;
						this.login = keystores.length > 0;
					})
				} else if (resp.raw.status === 401) {
					this.ready = true;
				} else {
					return Promise.reject(resp)
				}
			}).catch(error => {
				this.error = error;
			}).finally(() => {
				this.ready = true;
			});

			this.login = true;
			this.ready = true;
		},
		keystoreCreated(keystore: Keystore) {
			this.onboarding = false;
			this.keystore = keystore;
			this.keystores = [keystore];
		},
		loggedIn() {
			console.log("logged in");

			api.keystores().then((keystores) => {
				this.keystores = keystores
			}).catch(error => {
				this.error = error;
			}).finally(() => {
				this.login = false;
			});
		}
	},
});
</script>

<style lang="scss">
$bg: #0a0e11;
$accent: #00edb1;
$text-color: #fff;

body {
  margin: 0;
	background-color: $bg;
}

@import url('https://rsms.me/inter/inter.css');

* {
	font-family: 'Inter', sans-serif;
	font-weight: 300;
	font-size: 100%;
	color: $text-color;
	line-height: 1.325;
}
</style>
