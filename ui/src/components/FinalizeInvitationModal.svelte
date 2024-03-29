<script lang="ts">
  import type {Invitation} from "@/api";
  import Modal from "@/components/Modal.svelte";
  import {createEventDispatcher} from "svelte";
  import {hash} from "@/lib/color-hash.js";

  const dispatch = createEventDispatcher();
  const submit = () => {
    dispatch("submit");
  }

  export let show: boolean = false;
  export let invitation: Invitation;
</script>

<form on:submit|preventDefault={submit}>
  <Modal bind:show={show}>
    <div class="delete-title" slot="header">
      Are you sure you want to verify this invitation?
    </div>

    <div class="modal-content">
      Make sure the invitee's identity icon below matches the one on their side, and then proceed to verify.

      <div class="identity">
        <strong style="color: {hash(invitation.invitee.username)}">{invitation.invitee.username}</strong>'s identity:
        <div>
          {#if invitation.inviter.icon}
            <img style="width: 64px; height: 64px;" src={'data:image/svg+xml,'+encodeURIComponent(invitation.inviter.icon)} alt="">
          {:else if invitation.invitee.icon}
            <img style="width: 64px; height: 64px;" src={'data:image/svg+xml,'+encodeURIComponent(invitation.invitee.icon)} alt="">
          {/if}
        </div>
      </div>

    </div>

    <div class="modal-footer" slot="footer">
      <button class="button transparent" on:click={() => show = false} type="button">Cancel</button>
      <button class="button green" type="submit">Verify</button>
    </div>
  </Modal>
</form>

<style lang="scss">
  .modal-footer {
    display: flex;
    flex-direction: row;
    justify-content: flex-end;
    margin-top: 22px;
  }

  .delete-title {
    padding: 4px;
    box-sizing: border-box;
    font-size: 1.1rem;
  }

  .modal-content {
    padding: 4px;
  }

  .identity {
    padding-top: 20px;
    img {
      margin-top: 10px;
    }
  }

  strong {
    font-weight: bold;
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
