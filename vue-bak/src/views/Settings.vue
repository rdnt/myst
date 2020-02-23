<template>
  <div class="settings">
    <!-- <Tabs></Tabs> -->
    <div id="header" class="tabs">
      <router-link
        v-for="route in routes"
        :to="{ path: route.path }"
        :key="route.path"
        class="tab"
        :class="{ active: $route.path == route.path }"
        >{{ route.name }}</router-link
      >
    </div>
    <!-- <button class="tab active">Appearance</button>
    <button class="tab">Security</button>
    <button class="tab">Backup / Restore</button> -->
  </div>
</template>

<script>
// import Tabs from "@/components/Tabs.vue";

export default {
  components: {
    // Tabs
  },
  data() {
    return {
      routes: [],
      currentTab: null
    };
  },
  created: function() {
    this.$router.options.routes.forEach(route => {
      if (route.path === "/settings") {
        route.children.forEach(route => {
          this.routes.push(route);
        });
      }
    });
  }
};
</script>
