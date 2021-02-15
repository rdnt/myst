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
    <div id="sites-list">
      <div class="row header">
        <div class="col icon"></div>
        <div class="col name">
          <div class="label">Site</div>
        </div>
        <div class="col user">
          <div class="label">Username</div>
        </div>
        <div class="col pass">
          <div class="label">Password</div>
        </div>
      </div>
      <!-- <simplebar class="entries" data-simplebar-auto-hide="false"> -->
      <div class="entries scroll-view">
        <div class="scroll-overlay" v-on:scroll="scrollView">
          <div v-for="(entry, i) in entries" :key="i" class="dummy"></div>
        </div>
        <!-- <div
          class="scroll-area"
          v-on:scroll="scrollView"
          v-on:mouseover="mouseover"
          v-on:mouseleave="mouseleave"
        > -->
        <div class="scroll-area">
          <router-link
            v-for="(entry, i) in entries"
            :key="i"
            class="row entry card"
            :to="{ path: '/site/' + i + '/edit' }"
          >
            <div class="col icon">
              <img src="/assets/images/logos/medium-logo.png" />
            </div>
            <div class="col name">
              <div class="label">Medium</div>
            </div>
            <div class="col user">
              <div class="label">tasos@sht.gr</div>
            </div>
            <div class="col pass">
              <div class="label">Ex4mplePassw0rd</div>
            </div>
          </router-link>
        </div>
        <!-- </simplebar> -->
      </div>
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

export default {
  name: "Sites",
  components: {
    // simplebar
  },
  data() {
    return {
      over: false,
      showEditModal: false,
      entries: [
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {},
        {}
      ]
    };
  },
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
