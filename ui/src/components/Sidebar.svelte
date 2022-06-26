<script lang="ts">
  import Link from "@/components/Link.svelte";
  import {hash} from "@/lib/color-hash";
  import {invitations} from "@/stores/invitations";
  import {currentUser} from "@/stores/user";
  import {onMount} from "svelte";

  export let keystores;
  export let showCreateKeystoreModal: boolean;

  $: newInvitationsCount = $invitations.filter(inv => inv.inviteeId === $currentUser.id && inv.status === 'pending').length;

  onMount(() => {

  })
</script>

<div class="sidebar">
  <h4>Myst</h4>
  {#if $currentUser}
    <h6>Signed in as <strong style="color: {hash($currentUser.id)}">{$currentUser.id}</strong></h6>
  {/if}

  <div class="list">
    <h5 style="display: flex">
      Keystores
      <span style="position: relative;margin-left: auto;font-size:1.4rem;top:-9px;font-weight: bold;" on:click={() => {showCreateKeystoreModal = true}}>＋</span>
    </h5>

    {#each keystores as keystore}
      <Link path="/keystore/{keystore.id}">{keystore.name}</Link>
    {/each}
  </div>

  <div class="list bottom">
    <h5>Sync</h5>
    <div class="rel">
      <Link path="/invitations">
        Invitations
        {#if newInvitationsCount > 0}
          <div class="badge">{newInvitationsCount}</div>
        {/if}
      </Link>


<!--      <Link active={showInvitations} on:click={() => showInvitations = !showInvitations}>-->
<!--        Invitations-->
<!--        {#if $invitations.length > 0 || true}-->
<!--          <div class="badge">{$invitations.length + 2}</div>-->
<!--        {/if}-->
<!--      </Link>-->
    </div>
    <h5 class="version">v0.0.0-0123456</h5>
  </div>

</div>

<style lang="scss">
  .sidebar {
    position: relative;
    background-color: #0a0e11;
    height: 100%;
    padding: 12px 18px;
    box-sizing: border-box;
    flex-basis: 300px;
    display: flex;
    flex-direction: column;
    overflow-y: auto;

    h4 {
      font-weight: 700;
      font-size: 2rem;
      padding: 0 12px;
      margin: 0;
      margin-top: 12px;
    }

    h5 {
      height: 20px;
      padding: 0 12px;
      margin: 0 0 12px;
      color: #8a8f9f;
      text-transform: uppercase;
      font-size: .85rem;
      font-weight: 600;
      letter-spacing: .5px;

      strong {
        font-weight: 700;
      }
    }

    h6 {
      height: 20px;
      padding: 0 12px;
      color: #8a8f9f;
      margin: 10px 0 10px;
      font-size: .9rem;
      font-weight: 500;
      margin-bottom: 40px;

      strong {
        font-weight: 600;
      }
    }

    .rel {
      position: relative;
      width: 100%;
    }

    .version {
      margin-top: 40px;
      opacity: .75;
      font-weight: 500;
      text-transform: none;
    }

    .list {
      display: flex;
      flex-direction: column;
      flex-wrap: wrap;
      margin: 0;
      margin-top: 20px;

      &.bottom {
        bottom: 0;
        padding-top: 48px;
        margin-top: auto;
      }

      :global(.link a), :global(.link button) {
        margin: 0;
        line-height: 1.4;
        background-color: transparent;
        border: none;
        color: inherit;
        text-decoration: none;
        display: flex;
        width: 100%;
        align-items: center;
        border-radius: 5px;
        position: relative;
        cursor: pointer;
        padding: 10px 12px;
        font-size: 1.1rem;
        white-space: nowrap;
        text-overflow: ellipsis;
        font-weight: 500;
        box-sizing: border-box;
        margin-bottom: 4px;

        &:hover {
          background-color: rgba(#1e2328, .75);
        }
      }

      :global(.link.active a), :global(.link.active button) {
        font-weight: 500;
        color: #00edb1;
        background-color: #0c1d19;
        background-color: rgba(#1e2328, .75);

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

  .badge {
    margin-left: auto;
    font-weight: bold;
  }
</style>
