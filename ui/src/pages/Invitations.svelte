<script lang="ts">
  import api from "@/api";
  import AcceptInvitationModal from "@/components/AcceptInvitationModal.svelte";
  import DeclineInvitationModal from "@/components/DeclineInvitationModal.svelte";
  import Invitation from "@/components/Invitation.svelte";
  import {hash} from "@/lib/color-hash";
  import {format} from "@/lib/date";
  import {getInvitations, invitations} from "@/stores/invitations";
  import {currentUser} from "@/stores/user";
  import {useFocus} from "svelte-navigator";
  import {onMount} from "svelte";

  const registerFocus = useFocus();

  $: incomingInvitations = $invitations.filter((inv) => inv.invitee.id === $currentUser.id && inv.status === 'pending');
  $: outgoingInvitations = $invitations.filter((inv) => inv.inviter.id === $currentUser.id && (inv.status === 'pending'));
  $: pastInvitations = $invitations.filter((inv) => inv.status !== 'pending')

  let invitation: Invitation;
  let showAcceptInvitationModal: boolean = false;
  let showDeclineInvitationModal: boolean = false;

  function acceptInvitation() {
    api.acceptInvitation({
      invitationId: invitation.id
    }).then(() => {
      showAcceptInvitationModal = false;
      getInvitations();
    });
  }

  function declineInvitation() {
    api.declineOrCancelInvitation({
      invitationId: invitation.id
    }).then(() => {
      showDeclineInvitationModal = false;
      getInvitations();
    });
  }

  function showAcceptInvitationModalFunc(inv: Invitation) {
    invitation = inv;
    showAcceptInvitationModal = true;
  }

  function showDeclineInvitationModalFunc(inv: Invitation) {
    invitation = inv;
    showDeclineInvitationModal = true;
  }

  onMount(() => {
    getInvitations()
  })
</script>

<div class="invitations-list" use:registerFocus>
  {#if incomingInvitations.length > 0}
    <div class="section">
      <h5>Incoming Invitations</h5>
      {#each incomingInvitations as inv}
        <div class="invitation">
          <span class="icon">
            <span style="background-color: {hash(inv.inviter.username)}">
              {inv.inviter.username.slice(0, 2).toUpperCase()}
            </span>
          </span>
          <div class="info">
            <span class="name">
              {inv.keystore.name}
            </span>
            <span class="user">
              Invited by <strong>{inv.inviter.username}</strong> {format(inv.createdAt)}
            </span>
          </div>
          <div class="actions">
            <button class="button red" on:click={() => {showDeclineInvitationModalFunc(inv)}}>Decline</button>
            <button class="button green" on:click={() => {showAcceptInvitationModalFunc(inv)}}>Accept</button>
          </div>
        </div>
      {/each}
    </div>
  {/if}

  {#if outgoingInvitations.length > 0}
    <div class="section">
      <h5>Outgoing Invitations</h5>

      {#each outgoingInvitations as inv}
        <div class="invitation">
          <span class="icon">
            <span style="background-color: {hash(inv.invitee.username)}">
              {inv.invitee.username.slice(0, 2).toUpperCase()}
            </span>
          </span>
          <div class="info">
            <span class="name">
              {inv.keystore.name}
            </span>
            <span class="user">
              Invited <strong>{inv.invitee.username}</strong> {format(inv.createdAt)}
            </span>
          </div>
          <div class="actions">
            <button class="button red" on:click={() => {showDeclineInvitationModalFunc(inv)}}>Delete Invitation</button>
          </div>
        </div>
      {/each}
    </div>
  {/if}

  {#if pastInvitations.length > 0}
    <div class="section">
      <h5>Past Invitations</h5>
      {#each pastInvitations as inv}
        <div class="invitation">
        <span class="icon">
          <span style="background-color: {hash(inv.inviter.id === $currentUser.id ? inv.invitee.username : inv.inviter.username)}">
            {#if inv.inviter.id === $currentUser.id}
              {inv.invitee.username.slice(0, 2).toUpperCase()}
            {:else}
              {inv.inviter.username.slice(0, 2).toUpperCase()}
            {/if}
          </span>
        </span>
          <div class="info">
          <span class="name">
            {inv.keystore.name}
          </span>
            <span class="user">
            {#if inv.status === 'accepted'}
              {#if inv.inviter.id === $currentUser.id}
                Shared with <strong style="color: {hash(inv.invitee.username)}">{inv.invitee.username}</strong> since {format(inv.acceptedAt)}
              {:else}
                Invitation accepted {format(inv.acceptedAt)}.
                Waiting for <strong style="color: {hash(inv.inviter.username)}">{inv.inviter.username}</strong> to come online.
              {/if}
            {:else if inv.status === 'finalized'}
              {#if inv.inviter.id === $currentUser.id}
                Shared with <strong style="color: {hash(inv.invitee.username)}">{inv.invitee.username}</strong> since {format(inv.acceptedAt)}
              {:else}
                Being shared by <strong style="color: {hash(inv.inviter.username)}">{inv.inviter.username}</strong> since {format(inv.acceptedAt)}
              {/if}
            {:else if inv.status === 'deleted'}
              {#if inv.inviter.id === $currentUser.id}
                Invitation to <strong style="color: {hash(inv.invitee.username)}">{inv.invitee.username}</strong> deleted {format(inv.deletedAt)}
              {/if}
            {:else if inv.status === 'declined'}
              {#if inv.inviter.id === $currentUser.id}
                Invitation declined by <strong style="color: {hash(inv.invitee.username)}">{inv.invitee.username}</strong> {format(inv.declinedAt)}
              {:else}
                Declined invitation from <strong style="color: {hash(inv.inviter.username)}">{inv.inviter.username}</strong> {format(inv.declinedAt)}
              {/if}
            {/if}
          </span>
          </div>
          <div class="actions">
            {#if inv.status === 'finalized' || inv.status === 'accepted'}
              {#if inv.inviter.id === $currentUser.id}
                <button class="button red" on:click={() => {/*todo*/}}>Revoke Access</button>
              {:else}
                <button class="button red" on:click={() => {/*todo*/}}>Deny Access</button>
              {/if}
            {/if}
          </div>
        </div>
      {/each}
    </div>
  {/if}
</div>

{#if invitation}
  <AcceptInvitationModal bind:show={showAcceptInvitationModal} {invitation} on:submit={() => {acceptInvitation()}}/>
  <DeclineInvitationModal bind:show={showDeclineInvitationModal} {invitation} on:submit={() => {declineInvitation()}}/>
{/if}

{#if invitation}
  <!--  <Invitation {invitation}/>-->
{/if}


<style lang="scss">
  $accent: #00edb1;

  h5 {
    height: 20px;
    padding: 0 16px;
    margin: 12px 0 12px;
    color: #8a8f9f;
    text-transform: uppercase;
    font-size: .85rem;
    font-weight: 600;
    letter-spacing: .5px;
  }

  .section {
    margin-bottom: 36px;

    &:last-of-type {
      margin-bottom: 0;
    }
  }

  .invitations-list {
    //position: absolute;
    //left: 100%;
    //bottom: 0;
    //width: 400px;
    //margin-left: 40px;
    background-color: #111519;
    border-radius: 5px;
    min-height: 400px;
    flex-basis: 100%;
    z-index: 1;
    //opacity: 0;
    //transform: scale(0.98);
    //pointer-events: none;
    //transition: .15s ease;
    transform-origin: bottom left;
    padding: 20px;
    box-sizing: border-box;
    height: 100%;
    overflow-y: auto;

    //&.show {
    //  opacity: 1;
    //  transform: scale(1);
    //  pointer-events: all;
    //}
  }

  .invitation {
    position: relative;
    display: flex;
    flex-direction: row;
    flex-wrap: nowrap;
    align-items: center;
    padding: 10px 14px;
    box-sizing: border-box;
    border-radius: 5px;
    min-height: 20px;
    text-decoration: none;
    margin-bottom: 2px;
    cursor: default;

    .info {
      display: flex;
      flex-direction: column;
    }

    .actions {
      margin-left: auto;
    }

    &.header {
      color: rgb(138, 143, 152);
      height: 60px;
      padding: 0 34px;
      margin-top: 20px;

      .name {
        flex-basis: calc(30% + 34px);
      }

      span {
        color: rgba(#8a8f9f, .75);
        text-transform: uppercase;
        font-size: 0.85rem;
        font-weight: 600;
        letter-spacing: 0.5px;
        pointer-events: none;
      }

      button {
        background-color: transparent;
        border: none;
        padding: 0;
        margin: 0;

        img {
          width: 10px;
          height: 10px;
          margin-right: 10px;
          opacity: .5;
        }
      }
    }

    .icon {
      padding-right: 16px;

      img {
        display: block;
        width: 48px;
        height: 48px;
        vertical-align: baseline;
      }

      span {
        width: 48px;
        height: 48px;
        display: inline-flex;
        justify-content: center;
        align-items: center;
        font-weight: bold;
        border-radius: 50%;
        background-color: rgba(#aabbcc, .1);
        font-size: 1rem;
      }
    }

    .name {
      flex-basis: calc(30%);
      font-weight: 500;
      margin: 2px 0;
      font-size: 1.1rem;
      color: #fff;
    }

    .name, .user, .pass {
      flex-grow: 1;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
      padding-right: 20px;
      min-height: 20px;
      box-sizing: border-box;
      vertical-align: middle;
      display: inline-block;
      line-height: 1.3;
    }

    .user {
      color: darken(#8a8f9f, 5%);
      font-weight: 500;
      font-size: 1rem;

      strong {
        font-weight: 600;
      }
    }

    .user, .pass {
      flex-basis: 35%;
    }

    &:not(.header) {

      &:hover {
        background-color: #191e23;
        color: #f7f8f8;
      }

      &:active {
        background-color: lighten(#191e23, 2%);
        color: rgb(215, 216, 219);
        opacity: 1;
      }

      &.active {
        background-color: lighten(#1e2328, 3%);
      }

      .pass {
        position: relative;

        button {
          position: absolute;
          right: 0;
          background-color: transparent;
          border: none;
          padding: 0;
          outline: none;
          display: inline-block;
          opacity: 0;

          img {
            height: 20px;
            display: block;
          }
        }
      }
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
