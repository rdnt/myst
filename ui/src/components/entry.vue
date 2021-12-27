<template>
  <div id="entry">
    {{ entry }}

    <modal :link="{ name: 'entry' }">
      <div class="modal-header">
        Header
      </div>
      <div class="modal-body">
        Body<br />
        Body<br />
        Body<br />
        Body<br />
        Body<br />
        Body<br />
        Body<br />
        Body<br />
        Body<br />
        Body<br />
        Body<br />
        Body<br />
        Body<br />
        Body<br />
        Body<br />
        Body<br />
        Body<br />
        Body<br />
        Body<br />
        Body<br />
        Body<br />
        Body<br />
        Body<br />
        Body<br />
        Body<br />
        Body<br />
        Body<br />
        Body<br />
        Body<br />
        Body<br />
        Body<br />
        Body<br />
        Body<br />
        Body<br />
        Body<br />
        Body<br />
      </div>
      <div class="modal-footer">
        Footer
      </div>
    </modal>
  </div>
</template>

<script>
import Modal from "@/components/modal.vue";
import { mapState } from "vuex";

export default {
  components: {
    Modal
  },
  computed: mapState({
    keystore: state => state.keystore.keystore,
    entry(state) {
      if (!state.keystore.keystore) {
        return null;
      }

      return state.keystore.keystore.entries.filter(e => {
        return e.id === this.$route.params.entryId;
      });
    }
  }),
  data() {
    return {};
  }
};
</script>
