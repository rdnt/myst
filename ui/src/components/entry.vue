<template>
  <div id="entry">
    {{ entry }}

    <modal :link="{ name: 'entry' }" class="entry-modal">
      <div class="modal-header">
        Create Entry
      </div>
      <div class="modal-body">
        <div class="field-new">
          <input name="test" type="text" placeholder="" />
          <label>Domain</label>
        </div>
        <div class="field-new">
          <input name="test" type="text" placeholder="" />
          <label>Username</label>
        </div>
        <div class="field-new">
          <input name="test" type="text" placeholder="" />
          <label>Password</label>
        </div>
      </div>
      <div class="modal-footer">
        <button class="button create-entry-cancel-button">
          Cancel
        </button>
        <button class="button create-entry-button">
          Create
        </button>
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

//<style lang="scss">
//.field-new {
//  position: relative;
//  display: flex;
//  flex-direction: column;
//
//  &:not(:last-child) {
//    margin-bottom: 12px;
//  }
//
//  input {
//    display: block;
//    width: 100%;
//    font-size: 1rem;
//    background-color: transparent;
//    caret-color: #fff;
//    height: 60px;
//    padding: 10px;
//    position: relative;
//    //margin-top: 16px;
//    //margin: 4px 0;
//    //&:not(:last-child) {
//    //  margin-bottom: 10px;
//    //}
//    //margin-bottom: 4px;
//    font-weight: 300;
//    transition: 0.1s ease;
//    color: rgba(215, 216, 219, 0.75);
//
//    //box-shadow: inset 0 0 0 1px rgba(45, 47, 54, 1);
//    //border-radius: 5px;
//
//    order: 2;
//
//    &::placeholder {
//      color: #63666d;
//    }
//
//    &:focus {
//      color: rgba(215, 216, 219, 1);
//
//      & + label {
//        color: lighten(#63666d, 15%);
//
//        &:after {
//          box-shadow: inset 0 0 0 1px lighten(rgba(48, 50, 54, 1), 5%);
//        }
//      }
//    }
//  }
//
//  label {
//    display: block;
//    pointer-events: none;
//    transition: 0.1s ease;
//    padding: 0 8px;
//    color: #63666d;
//    order: 1;
//
//    &:after {
//      content: "";
//      position: absolute;
//      display: block;
//      bottom: 8px;
//      left: 0;
//      width: 100%;
//      height: 44px;
//      box-shadow: inset 0 0 0 1px rgba(48, 50, 54, 1);
//      border-radius: 5px;
//      pointer-events: none;
//      transition: 0.1s ease;
//    }
//  }
//}
//
//.button {
//  cursor: default;
//  font-weight: 300;
//  font-size: 14px;
//  height: 36px;
//  display: flex;
//  align-items: center;
//  color: #d7d8db;
//  padding: 8px 16px;
//  border-radius: 4px;
//  background-color: rgba(#2d2f36, 0.5);
//  margin: 0;
//  border: none;
//  font-weight: normal;
//  font-size: 1rem;
//
//  background-color: rgba(#393c49, 0.25);
//
//  &:hover {
//    background-color: rgba(#393c49, 0.2);
//    color: #fff;
//  }
//}
//
//.create-entry-button {
//  margin-left: 12px;
//  background-color: #393c49;
//  //background-color: rgba(#7a7b7f, 0.25);
//  //color: #0a0b0c;
//  //font-weight: 500;
//  //font-weight: bold;
//  &:hover {
//    background-color: rgba(#393c49, 0.75);
//    //background-color: #5acba8;
//    //color: #0a0b0c;
//  }
//}
</style>
