<script lang="ts">
  import Modal from "./Modal.svelte";
  // import * as models from "../api/generated/models";
  import InputField from "./InputField.svelte";
  import Field from "./Field.svelte";
  import {createEventDispatcher, onMount} from 'svelte';
  // import api from "../api/index";
  import type {Invitation, Keystore} from "../api/generated";
  import {showError, showMessage} from "../stores/messages";
  import api from "../api";

  const dispatchCreated = createEventDispatcher<{ created: { id: string } }>();

  export let show: boolean = false;
  export let keystore: Keystore;

  $: user = '';

  $: userValid = user.trim() !== '';

  $: allowSubmit = userValid;

  $: {
    if (!show) {
      user = ''
    }
  }

  onMount(async () => {
    user = ''
  });

  const submit = () => {
    if (!allowSubmit) {
      return;
    }

    api.createInvitation({
      keystoreId: keystore.id,
      createInvitationRequest: {inviteeId: user}
    }).then((inv: Invitation) => {
      showMessage("Invitation sent");
      dispatchCreated('created', {id: inv.id})
    }).catch((err) => {
      showError("Create Invitation Failed");
      console.error(err)
    })
  };
</script>

<form class="create-invitation-modal" on:submit|preventDefault={submit}>
  <Modal bind:show>
    <div class="modal-header" slot="header">
      <div class="title">
        <span>Share with...</span>
      </div>
    </div>

    <div class="modal-content">
      <Field label="Keystore" value={keystore.name}/>
      <InputField bind:value={user} label="User"/>
    </div>

    <div class="modal-footer" slot="footer">
      <button class="button transparent" on:click={() => show = false} type="button">Cancel</button>
      <button class="button green" type="submit">Invite User</button>
    </div>
  </Modal>
</form>

<style lang="scss">
  .create-invitation-modal {

    .modal-header {
      display: flex;
      flex-direction: row;
      //margin-top: 10px;

      .image {
        width: 64px;
        height: 64px;
        padding-right: 20px;

        img {
          width: 64px;
          height: 64px;
        }
      }

      .title {
        font-weight: 600;
        font-size: 1.8rem;
        margin: 0;

        .keystore-name {
          font-weight: 400;
          font-style: italic;
        }
      }
    }


    .modal-content {
      padding-top: 40px;
      box-sizing: border-box;
      box-sizing: border-box;
    }

    .modal-footer {
      display: flex;
      flex-direction: row;
      justify-content: flex-end;
      margin-top: 12px;
      box-sizing: border-box;
    }

  }


  $accent: #00edb1;
  .button {
    outline: none;
    border: none;
    height: 40px;
    font-size: 1.1rem;
    font-weight: 500;
    padding: 0 16px;
    border-radius: 5px;
    background-color: rgba(#202228, 1);
    color: #fff;
    margin-left: 10px;

    &.left {
      //margin-right: auto;
    }

    &.disabled {
      //background-color: #161819;
      opacity: .5;
    }

    &.green {
      background-color: rgba(#002e23, .9);
      color: $accent;

      &.disabled {
        background-color: #0c1d19;
      }
    }

    &.transparent {
      background-color: transparent;
      padding: 0 12px;

      &.disabled {

      }
    }

    &.red {
      background-color: #2e2020;
      color: #ff9999;

      &.disabled {
        background-color: rgba(29, 29, 12, 0.99);
      }
    }
  }
</style>
