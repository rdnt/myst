<script lang="ts">
  import api from "@/api";
  import type {Invitation} from "@/api";
  import AcceptInvitationModal from "@/components/AcceptInvitationModal.svelte";
  import {myInvitations} from "@/stores/invitations";

  export let show: boolean;

  let invitation: Invitation;
  let showAcceptInvitationModal: boolean = false;

  function acceptInvitation() {
    api.acceptInvitation({
      invitationId: invitation.id
    }).then(() => {
      showAcceptInvitationModal = false;
    });
  }

  function showAcceptInvitationModalFunc(inv: Invitation) {
    invitation = inv;
    showAcceptInvitationModal = true;
  }
</script>

<div class="invitations-list" class:show={show}>
  {#each $myInvitations as inv}
    <div class="invitation">
      <div class="details">
        <div>{inv.inviterId}</div>
        <div>{inv.keystoreId}</div>
      </div>
      <button class="button" on:click={() => {showAcceptInvitationModalFunc(inv)}}>accept</button>
    </div>
  {/each}
</div>

{#if invitation}
  <AcceptInvitationModal bind:show={showAcceptInvitationModal} {invitation} on:submit={() => {acceptInvitation()}}/>
{/if}

<style lang="scss">
  $accent: #00edb1;

  .invitations-list {
    //position: absolute;
    //left: 100%;
    //bottom: 0;
    //width: 400px;
    //margin-left: 40px;
    background-color: #111519;
    border-radius: 5px;
    min-height: 400px;
    flex-basis: 50%;
    z-index: 1;
    //opacity: 0;
    //transform: scale(0.98);
    //pointer-events: none;
    //transition: .15s ease;
    transform-origin: bottom left;
    padding: 20px;
    box-sizing: border-box;

    //&.show {
    //  opacity: 1;
    //  transform: scale(1);
    //  pointer-events: all;
    //}
  }

  .invitation {
    border-radius: 4px;
    padding: 20px;
    display: flex;
    justify-content: space-between;

    .details {
      display: flex;
      flex-direction: column;
    }

    &:hover {
      background-color: #25282e;
      color: #f7f8f8;
    }
  }

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
