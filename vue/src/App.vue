<template>
  <div id="app" class="show" :class="{ electron: client == 'electron' }">
    <system-bar v-if="client == 'electron'"></system-bar>
    <preloader></preloader>
    <login></login>
    <main class="transition">
      <navigation></navigation>
      <div id="container">
        <router-view></router-view>
      </div>
    </main>
  </div>
</template>

<script>
import SystemBar from "./components/SystemBar.vue";
import Preloader from "./components/Preloader.vue";
import Navigation from "./components/Navigation.vue";
import Login from "./components/Login.vue";

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
  created() {
    // this.$store.dispatch('initialize')
    this.$socket.onopen = () => {
      this.$store.dispatch("get", "sites");
      // this.$store.send('getSites', '')
      // this.$store.commit('initialize')
    };

    // setTimeout(() => {
    //
    // }, 1000)
  },
  mounted() {
    // console.log(this.$wasm.test());
  }
};
</script>

<style lang="scss">
@import "./styles/App.scss";
</style>
