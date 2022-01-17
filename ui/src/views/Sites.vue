<template>
  <div id="container" class="entries-container">
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

    <div id="entries-list" :class="{ expanded: !showEditModal }">
      <div id="header">
        <router-link :to="{ path: '' }" class="breadcrumb"
          >Keystore
        </router-link>
        <!--      <router-link :to="{ path: '' }" class="breadcrumb">Passwords</router-link>-->
        <div
          class="create-entry-icon"
          v-on:click="
            () => {
              this.showCreateEntryModal = !this.showCreateEntryModal;
            }
          "
        >
          Create New
        </div>
      </div>
      <div id="entries" v-if="keystore">
        <div class="entry header">
          <div class="name">
            Domain
            <button><img alt="" src="/assets/images/sort-arrow.svg" /></button>
          </div>
          <div class="user">
            Username
            <button><img alt="" src="/assets/images/sort-arrow.svg" /></button>
          </div>
          <div class="pass">
            Password
            <button><img alt="" src="/assets/images/sort-arrow.svg" /></button>
          </div>
        </div>
        <!--      <div class="entries" v-if="keystore">-->
        <router-link
          v-for="entry in keystore.entries"
          :key="entry.id"
          class="entry"
          :to="{
            name: 'entry',
            params: { entryId: entry.id, entry: entry }
          }"
        >
          <span class="icon">
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
    <!--    <div id="entry-container" :class="{ show: showEditModal }">-->
    <!--      <router-view></router-view>-->
    <!--    </div>-->
    <modal :show="showCreateEntryModal" v-on:hide="hideCreateEntryModal">
      <template v-slot:header>
        Add Entry
      </template>
      <div class="field-new">
        <input
          name="test"
          placeholder=""
          type="text"
          v-model="createEntryModal.label"
        />
        <label>Domain</label>
      </div>
      <div class="field-new">
        <input
          name="test"
          placeholder=""
          type="text"
          v-model="createEntryModal.username"
        />
        <label>Username</label>
      </div>
      <div class="field-new">
        <input
          name="test"
          placeholder=""
          type="text"
          v-model="createEntryModal.password"
        />
        <label>Password</label>
      </div>
      <template v-slot:footer>
        <!--        <button class="button create-entry-cancel-button">-->
        <!--          Cancel-->
        <!--        </button>-->
        <button class="button create-entry-button" v-on:click="createEntry">
          Create
        </button>
      </template>
    </modal>
  </div>
</template>

<script>
import Modal from "@/components/modal.vue";

import { mapState } from "vuex";

export default {
  name: "Sites",
  components: {
    Modal
  },
  data() {
    return {
      over: false,
      showEditModal: false,
      showCreateEntryModal: false,
      createEntryModal: {
        label: "",
        username: "",
        password: ""
      }
    };
  },
  computed: mapState({
    keystore: state => state.keystore.keystore
  }),
  watch: {
    $route: function() {
      if (this.$route.name == "entry") {
        this.showEditModal = true;
      } else {
        this.showEditModal = false;
      }
    }
  },
  methods: {
    hideCreateEntryModal() {
      this.showCreateEntryModal = false;
    },
    createEntry() {
      this.$store
        .dispatch("keystore/createEntry", {
          keystoreId: "0000000000000000000000",
          entry: this.createEntryModal
        })
        .then(resp => {
          console.log(resp);
          this.$set(this.createEntryModal, "label", "");
          this.$set(this.createEntryModal, "username", "");
          this.$set(this.createEntryModal, "password", "");
          this.showCreateEntryModal = false;
        });
    },
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

<style lang="scss">
.field-new {
  position: relative;
  display: flex;
  flex-direction: column;

  &:not(:last-child) {
    margin-bottom: 12px;
  }

  input {
    display: block;
    width: 100%;
    font-size: 1rem;
    background-color: transparent;
    caret-color: #fff;
    height: 60px;
    padding: 10px;
    position: relative;
    font-weight: 400;
    transition: 0.1s ease;
    color: rgba(215, 216, 219, 0.75);
    order: 2;

    &::placeholder {
      color: #63666d;
    }

    &:focus {
      color: rgba(215, 216, 219, 1);

      & + label {
        color: rgba(215, 216, 219, 0.9);

        &:after {
          box-shadow: inset 0 0 0 1px lighten(rgba(48, 50, 54, 1), 5%);
        }
      }
    }
  }

  label {
    display: block;
    pointer-events: none;
    transition: 0.1s ease;
    padding: 0 8px;
    color: rgba(215, 216, 219, 0.75);
    order: 1;

    &:after {
      content: "";
      position: absolute;
      display: block;
      bottom: 8px;
      left: 0;
      width: 100%;
      height: 44px;
      box-shadow: inset 0 0 0 1px rgba(48, 50, 54, 1);
      border-radius: 5px;
      pointer-events: none;
      transition: 0.1s ease;
    }
  }
}

.button {
  cursor: default;
  height: 36px;
  display: flex;
  align-items: center;
  color: #d7d8db;
  padding: 8px 16px;
  border-radius: 4px;
  margin: 0;
  border: none;
  font-weight: normal;
  font-size: 1rem;
  background-color: rgba(#393c49, 0.25);

  &:hover {
    background-color: rgba(#393c49, 0.2);
    color: #fff;
  }
}

.create-entry-button {
  margin-left: 12px;
  background-color: #393c49;

  &:hover {
    background-color: rgba(#393c49, 0.75);
  }
}
</style>
