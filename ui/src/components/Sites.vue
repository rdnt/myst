<template>
  <div id="sites" class="cards">
    <div class="header">
      <div class="cards-label">
        <span>Sites</span>
      </div>
      <div class="label">
        <span>Username</span>
      </div>
      <div class="label">
        <span>Password</span>
      </div>
      <div class="label">
        <span>URL</span>
      </div>
      <div class="label">
        <span>Notes</span>
      </div>
      <div id="entries-switch">
        <button
          class="grid"
          @click="rowsToCards()"
          :class="{
            active: mode == 'rows'
          }"
        >
          <img src="/assets/images/entries/grid.svg" />
        </button>
        <button
          class="rows"
          @click="cardsToRows()"
          :class="{
            active: mode == 'cards'
          }"
        >
          <img src="/assets/images/entries/rows.svg" />
        </button>
      </div>
    </div>
    <div class="sites">
      <!-- <div class="row"> -->
      <div
        v-for="site in sites"
        class="site card"
        :tabindex="site.id + 1"
        :key="site.id"
      >
        {{ site.id }}
        <!-- <div class="drag">
                    <img src="data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjQiIGhlaWdodD0iMjQiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyIgZmlsbC1ydWxlPSJldmVub2RkIiBjbGlwLXJ1bGU9ImV2ZW5vZGQiPjxwYXRoIGQ9Ik0yNCAxOHYxaC0yNHYtMWgyNHptMC02djFoLTI0di0xaDI0em0wLTZ2MWgtMjR2LTFoMjR6IiBmaWxsPSIjMTA0MGUyIi8+PHBhdGggZD0iTTI0IDE5aC0yNHYtMWgyNHYxem0wLTZoLTI0di0xaDI0djF6bTAtNmgtMjR2LTFoMjR2MXoiLz48L3N2Zz4=">
                </div> -->
        <div class="logo"></div>
        <div class="title">
          {{ site.name }}
        </div>
        <div class="username">
          {{ site.email }}
        </div>
        <div class="password">
          •••••••
        </div>
        <!-- <div class="field username">
                    <span>{{.Username}}</span>
                </div>
                <div class="field password">
                    <span>{{.Password}}</span>
                </div>
                <div class="field url">
                    <span>{{.URL}}</span>
                </div>
                <div class="field notes">
                    <span>{{.Notes}}</span>
                </div> -->
      </div>
      <!-- </div> -->
    </div>
  </div>
</template>

<script>
import { mapState } from "vuex";

export default {
  // computed: {
  //     entries() {
  //         return this.$store.sites
  //     }
  // },
  // computed: mapState([
  //     // arrow functions can make the code very succinct!
  //     'sites'
  // ]),
  // computed: mapState({
  //     sites: state => state.sites.all
  // }),
  computed: mapState({
    sites: state => state.sites.all
  }),
  // computed: {
  //     sites: {
  //         get() {
  //             return this.$store.state.sites.all
  //         },
  //         set(value) {
  //             this.$store.commit('setSites', value)
  //         }
  //     }
  // },
  data() {
    return {
      mode: "rows"
    };
  },
  methods: {
    rowsToCards: function() {
      this.compress = true;
      setTimeout(() => {
        this.expand = true;
      }, 500);
    },
    cardsToRows: function() {
      this.shrink = true;
      setTimeout(() => {
        this.decompress = true;
      }, 500);
      setTimeout(() => {
        this.compress = false;
        this.expand = false;
        this.shrink = false;
        this.decompress = false;
      }, 1000);
    }
  }
};
</script>
