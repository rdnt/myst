<template>
  <div id="app" class="show" :class="{ electron: client == 'electron' }">
    <system-bar v-if="client == 'electron'"></system-bar>
    <preloader></preloader>
    <login></login>
    <main class="transition" :class="{ show: loggedIn }">
      <navigation></navigation>
      <div id="container">
        <router-view></router-view>
      </div>
    </main>
  </div>
</template>

<script>
import { mapState } from "vuex";
import SystemBar from "./components/SystemBar.vue";
import Preloader from "./components/Preloader.vue";
import Navigation from "./components/Navigation.vue";
import Login from "./components/Login.vue";
// import Keystore from "@/domain";
// import Search from "./components/Search.vue";
// import Authenticator from "./components/Authenticator.vue";
// import PasswordGenerator from "./components/Generator.vue";
// import Sites from "./components/Sites.vue";

export default {
  name: "App",
  components: {
    SystemBar,
    Preloader,
    Navigation,
    Login
    // Search,
    // Sites,
    // Authenticator,
    // PasswordGenerator,
    // Tabs
  },
  data() {
    return {
      client: "browser"
    };
  },
  computed: mapState({
    loggedIn: state => state.keystore.keystore !== null
  }),
  created() {},
  mounted() {}
};

// eslint-disable-next-line no-unused-vars
</script>

<style lang="scss">
@use "@/styles/App";
</style>
