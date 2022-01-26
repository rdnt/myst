<template>
  <div id="app" class="show" :class="{ electron: client == 'electron' }">
    <system-bar v-if="client == 'electron'"></system-bar>
    <preloader></preloader>
    <login></login>
    <main class="transition" :class="{ show: loggedIn }">
      <!--      <navigation></navigation>-->
      <router-view></router-view>
    </main>
  </div>
</template>

<script>
import { mapState } from "vuex";

export default {
  name: "App",
  components: {},
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
