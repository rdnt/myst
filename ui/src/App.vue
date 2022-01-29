<template>
	<span>{{!ready ? 'Loading...' : ''}}</span>
	<InitializeKeystoreFullscreenModal v-if="onboarding" @created="keystoreCreated($event)" />
	<Login v-if="login" @login="loggedIn()" />
	<transition :duration="300" name="show">
		<main v-if="keystore">
			<Sidebar />
			<EntriesList :entries="keystore.entries"></EntriesList>
		</main>
	</transition>
</template>

<script lang="ts">
import {defineComponent} from "vue";
import InitializeKeystoreFullscreenModal from "./components/InitializeKeystoreFullscreenModal.vue";
import Login from "./components/LoginForm.vue";
import api from "./api";
import {Keystore, Entry} from "./api/generated";
import EntriesList from "./components/EntriesList.vue";
import Sidebar from "./components/Sidebar.vue";

export default defineComponent({
	name: "App",
	components: {Sidebar, EntriesList, InitializeKeystoreFullscreenModal, Login},
	data(): {
		onboarding: boolean,
		login: boolean,
		ready: boolean,
		keystore?: Keystore,
		keystores: Keystore[],
	} {
		return {
			onboarding: false,
			login: false,
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
			api.keystores().then((keystores) => {
				this.onboarding = keystores.length == 0;
				this.keystores = keystores
				this.keystore = keystores[0]
			}).catch((error: Response) => {
				if (error.status == 401) {
					this.login = true;
					return
				}

				console.log(error)
			}).finally(() => {
				this.ready = true;
			});
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
				this.keystore = keystores[0]
			}).catch(error => {
				console.log(error)
			}).finally(() => {
				this.login = false;
			});
		}
	},
});
</script>

<style lang="scss">
//$bg: #0a0e11;
$bg: #111519;
$accent: #00edb1;
$text-color: #fff;

:root {
	color-scheme: dark;
}

html, body {
	//height: 100%;
	//overflow: hidden
}

html {
	height: -webkit-fill-available;
}
body {
	margin: 0;
	background-color: $bg;

	height: 100vh;
	max-height: 100vh;
	max-height: -webkit-fill-available;
	//max-height: 100%;
}

main {
	display: flex;
	align-items: stretch;
	background-color: $bg;
	width: 100%;
}

@import url('https://rsms.me/inter/inter.css');

* {
	font-family: 'Inter', sans-serif;
	font-weight: 300;
	font-size: 100%;
	color: $text-color;
	line-height: 1.325;
}

*::-webkit-scrollbar {
	width: 0px;
	display: none;
	background: transparent;
}

* {
	scrollbar-width: none;
	-ms-overflow-style: none;
}

#app {
	height: 100%;
	display: flex;
	flex-direction: row;
	flex-grow: 1;
}

.show-enter-active, .show-leave-active {
	transition: .3s;
}

.show-enter-from {
	opacity: 0;
	transform: scale(.98);
}

.show-enter-to {
	opacity: 1;
	transform: scale(1);
}

.show-leave-from {
	opacity: 1;
	transform: scale(1);
}

.show-leave-to {
	opacity: 0;
	transform: scale(.98);
}
</style>
