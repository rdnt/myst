<script lang="ts">
  import api from "@/api";
  import Link from "@/components/Link.svelte";
  import {hash} from "@/lib/color-hash";
  import {invitations} from "@/stores/invitations";
  import {currentUser} from "@/stores/user";
  import {onMount} from "svelte";

  export let keystores;
  export let showCreateKeystoreModal: boolean;
  export let showAuthModal: boolean;

  let username, password: string;

  $: newInvitationsCount = $invitations.filter(inv => inv.invitee.id === $currentUser.id && inv.status === 'pending').length;

  onMount(() => {

  })

  function login() {
    if (!username || !password) {
      return;
    }

    api.login({
      requestBody: {
        username,
        password
      }
    })
  }
</script>

<div class="sidebar">
  <h4>Myst</h4>


  <div class="list">
    <h5 style="display: flex">
      Keystores
      <span style="position: relative;margin-left: auto;font-size:1.4rem;top:-9px;font-weight: bold;"
            on:click={() => {showCreateKeystoreModal = true}}>＋</span>
    </h5>

    {#each keystores as keystore}
      <Link path="/keystore/{keystore.id}">{keystore.name}</Link>
    {/each}
  </div>


  <div class="list bottom">

    {#if $currentUser}
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
    {/if}
    {#if $currentUser === undefined}
      <h6>
        <span class="auth-link btn" on:click={() => showAuthModal = true}>Sign in</span>
        <span class="auth-link">or</span>
        <span class="auth-link btn" on:click={() => showAuthModal = true}>Register</span>
      </h6>
    {:else if $currentUser === null}
      Disconnected
    {:else}
      <h6>Signed in as <strong style="color: {hash($currentUser.username)}">{$currentUser.username}</strong></h6>
    {/if}

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
      margin-bottom: 32px;
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
      margin-bottom: 40px;
    }

    .version {
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

  .auth-link {
    cursor: default;
    font-weight: 600;
    color: #cacfdf;

    &.btn {
      color: #fff !important;
      text-decoration: underline;
    }
  }

  $accent: #00edb1;
  button {
    outline: none;
    border: none;
    height: 48px;
    font-size: 1.1rem;
    font-weight: 500;
    padding: 0 40px;
    border-radius: 5px;
    margin: 0 5px;
    background-color: rgba(#202228, 1);
    color: #fff;

    &.disabled {
      background-color: #161819;
    }

    &.green {
      background-color: #002e23;
      color: $accent;

      &.disabled {
        background-color: #0c1d19;
      }
    }
  }
</style>
