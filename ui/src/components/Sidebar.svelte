<script lang="ts">
  import api from "@/api";
  import type {Invitation} from "@/api";
  import AcceptInvitationModal from "@/components/AcceptInvitationModal.svelte";
  import Link from "@/components/Link.svelte";

  export let keystores;

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

<div class="sidebar">
  <h4>Myst</h4>
  <div class="keystores-list">
    <h5>Keystores</h5>
    {#each keystores as keystore}
      <Link path="/keystore/{keystore.id}">{keystore.name}</Link>
    {/each}
  </div>


  {#await api.getInvitations() then invitations}
    {JSON.stringify(invitations, null, 2)}

    {#each invitations as inv}
      <button on:click={() => {showAcceptInvitationModalFunc(inv)}}>accept</button>
    {/each}

  {/await}
</div>

{#if invitation}
  <AcceptInvitationModal bind:show={showAcceptInvitationModal} {invitation} on:submit={() => {acceptInvitation()}}/>
{/if}
<style lang="scss">
  .sidebar {
    background-color: #0a0e11;
    height: 100%;
    padding: 20px;
    box-sizing: border-box;
    flex-basis: 250px;

    h4 {
      font-weight: 700;
      font-size: 2rem;
      margin: 10px 0 60px;
      padding: 0 20px;
    }

    h5 {
      height: 20px;
      padding: 0 20px;
      margin: 0 0 10px;
      color: #8a8f9f;
      text-transform: uppercase;
      font-size: .85rem;
      font-weight: 600;
      letter-spacing: .5px;
    }

    .keystores-list {
      display: flex;
      flex-direction: column;

      :global(.link a) {
        color: inherit;
        text-decoration: none;
        display: flex;
        align-items: center;
        border-radius: 5px;
        position: relative;
        color: #fff;
        cursor: pointer;
        height: 22px;
        padding: 10px 20px 10px 20px;
        font-size: 1.1rem;
        white-space: nowrap;
        text-overflow: ellipsis;
        margin-bottom: 2px;
        text-decoration: none;

        &:hover {
          background-color: rgba(#1e2328, .75);
        }
      }

      :global(.link.active a) {
        font-weight: 500;
        color: #00edb1;
        background-color: #0c1d19;

        &:after {
          content: '';
          position: absolute;
          top: 0;
          width: 10px;
          height: 100%;
          left: -20px;
          border-top-right-radius: 5px;
          border-bottom-right-radius: 5px;
        }
      }
    }
  }
</style>
