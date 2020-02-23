<template>
  <div id="navigation">
    <!-- @mouseover="hovering = true"
        @mouseleave="hovering = false"
        :class="{hover: hovering}"> -->
    <!-- <div class="logo-container">
      <img class="logo" src="/assets/images/vault.svg" alt="" />
      <img class="label" src="/assets/images/vaultlabel.svg" alt="" />
    </div> -->
    <div class="list">
      <router-link
        class="nav_btn item"
        v-for="route in routes"
        :to="{ path: route.path }"
        :key="route.path"
        :class="{ active: subIsActive(route.path) }"
      >
        <img :src="route.props.icon" alt="" />
        <span>{{ route.name }}</span>
      </router-link>
      <!-- <button
        class="nav_btn item"
        v-for="item in items"
        :key="item.id"
        :class="{ active: item.id == selected }"
        @mousedown="selectTab(item.id)"
        @animationend="click = -1"
      >
        <img :src="item.icon" alt="" />
        <span>{{ item.name }}</span>
      </button> -->
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      routes: [],
      hovering: false,
      tab: 0,
      selected: 0
    };
  },
  methods: {
    selectTab: function(index) {
      this.selected = index;
    },
    subIsActive(input) {
      const paths = Array.isArray(input) ? input : [input];
      return paths.some(path => {
        return this.$route.path.indexOf(path) === 0; // current path starts with this path string
      });
    }
  },
  mounted: function() {
    this.$nextTick(function() {
      this.$router.options.routes.forEach(route => {
        if (route.props && route.props.sidebar === true) {
          this.routes.push(route);
        }
      });
    });
  }
};
</script>
