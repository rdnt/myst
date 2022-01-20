<template>
  <div id="edit-site-page">
    <div id="header">
      <div class="navigation-breadcrumbs">
        <router-link class="breadcrumb" :to="{ path: '/sites' }"
          >Sites</router-link
        >
        <router-link class="breadcrumb" :to="{ path: '/sites' }">
          Site {{ $route.params.id }}
        </router-link>
        <!-- <router-link class="breadcrumb" :to="{ path: '/sites' }"
          >Properties</router-link
        > -->
      </div>
    </div>

    <!-- {{  }} -->
  </div>
</template>

<script>
// @ is an alias to /src

export default {
  name: "Sites",
  components: {},
  data() {
    return {
      entries: [{}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}, {}]
    };
  }
};
</script>
