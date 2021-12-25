<template>
  <div id="sites-page">
    <!-- <div id="header">
      <div id="search">
        <input
          type="text"
          name="query"
          value="test search"
          placeholder="Search…"
          autocomplete="off"
          autocorrect="off"
          autocapitalize="off"
          spellcheck="false"
          required
        />
        <img
          src="data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjQiIGhlaWdodD0iMjQiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyIgZmlsbC1ydWxlPSJldmVub2RkIiBjbGlwLXJ1bGU9ImV2ZW5vZGQiPjxwYXRoIGQ9Ik0xNS44NTMgMTYuNTZjLTEuNjgzIDEuNTE3LTMuOTExIDIuNDQtNi4zNTMgMi40NC01LjI0MyAwLTkuNS00LjI1Ny05LjUtOS41czQuMjU3LTkuNSA5LjUtOS41IDkuNSA0LjI1NyA5LjUgOS41YzAgMi40NDItLjkyMyA0LjY3LTIuNDQgNi4zNTNsNy40NCA3LjQ0LS43MDcuNzA3LTcuNDQtNy40NHptLTYuMzUzLTE1LjU2YzQuNjkxIDAgOC41IDMuODA5IDguNSA4LjVzLTMuODA5IDguNS04LjUgOC41LTguNS0zLjgwOS04LjUtOC41IDMuODA5LTguNSA4LjUtOC41eiIvPjwvc3ZnPg=="
        />
      </div>
    </div> -->
    <div id="header">
      <router-link :to="{ path: '' }" class="breadcrumb">Passwords</router-link>
      <router-link :to="{ path: '' }" class="breadcrumb">All</router-link>
    </div>
    <div id="entries">
      <div class="entry header">
        <div class="name">
          Domain
        </div>
        <div class="user">
          Username
        </div>
        <div class="pass">
          Password
        </div>
      </div>
      <!--      <div class="entries" v-if="keystore">-->
      <router-link
        v-for="entry in keystore.entries"
        :key="entry.id"
        class="entry"
        :to="{ path: '/entry/' + entry.id + '/edit' }"
      >
        <!--        <div class="icon">-->
        <!--          <img src="/assets/images/favicon.ico" />-->
        <!--        </div>-->
        <span class="icon">
          <!--<img src="/assets/images/favicon.ico" alt="" />-->
          <img :src="`http://${entry.label}/favicon.ico`" alt="" />
        </span>
        <span class="name">
          {{ entry.label }}
        </span>
        <span class="user">
          {{ entry.username }}
        </span>
        <span class="pass">
          {{ entry.password }}
        </span>
      </router-link>
      <!--      </div>-->
    </div>
    <!-- <div id="edit-site" :class="{ show: showEditModal }">
      <div class="modal">
        <router-view></router-view>
      </div>
      <router-link :to="{ path: '/sites' }" class="overlay"> </router-link>
    </div> -->
  </div>
</template>

<script>
// @ is an alias to /src
// import simplebar from "simplebar-vue";

import { mapState } from "vuex";

export default {
  name: "Sites",
  components: {
    // simplebar
  },
  data() {
    return {
      over: false,
      showEditModal: false
      // entries: [
      //   {},
      //   {},
      //   {},
      //   {},
      //   {},
      //   {},
      //   {},
      //   {},
      //   {},
      //   {},
      //   {},
      //   {},
      //   {},
      //   {},
      //   {},
      //   {},
      //   {},
      //   {},
      //   {},
      //   {},
      //   {},
      //   {},
      //   {},
      //   {},
      //   {},
      //   {},
      //   {},
      //   {},
      //   {},
      //   {},
      //   {},
      //   {}
      // ]
    };
  },
  computed: mapState({
    keystore: state => state.keystore.keystore
  }),
  watch: {
    $route: function() {
      if (this.$route.name == "EditSite") {
        this.showEditModal = true;
      } else {
        this.showEditModal = false;
      }
    }
  },
  methods: {
    scrollView(event) {
      if (event.target.classList.contains("scroll-overlay")) {
        if (!this.over) {
          event.target.parentNode.querySelector(".scroll-area").scrollTop =
            event.target.scrollTop;
        }
      } else {
        if (this.over) {
          event.target.parentNode.querySelector(".scroll-overlay").scrollTop =
            event.target.scrollTop;
        }
      }
      // console.log(event);
      // console.log();
      // event.target.parentNode.dispatchEvent(new Event("scroll", event));
      // console.log(event.target.parentNode.querySelector(".scroll-area"));
      // console.log(event.target);
    },
    scrolls() {
      console.log("Scroll triggered");
    },
    mouseover() {
      this.over = true;
    },
    mouseleave() {
      this.over = false;
    }
  }
};
</script>
