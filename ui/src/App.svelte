<script lang="ts">
  import api from "@/api";
  import CreateKeystoreModal from "@/components/CreateKeystoreModal.svelte";
  import LoginForm from "@/components/LoginForm.svelte";
  import Messages from "@/components/Messages.svelte";
  import OnboardingForm from "@/components/OnboardingForm.svelte";
  import Sidebar from "@/components/Sidebar.svelte";
  import Invitations from "@/pages/Invitations.svelte";
  import Keystores from "@/pages/Keystores.svelte";
  import {getKeystores, keystores} from "@/stores/keystores.ts";
  import {currentUser, getCurrentUser} from "@/stores/user";
  import {onMount} from 'svelte';
  import {Route, Router} from "svelte-navigator";

  let onboarding = false;
  let ready = false;
  let login = false;

  let keystore = null;
  let showCreateKeystoreModal: boolean = false;

  const healthCheck = () => {
    api.healthCheck().then(() => {
      console.log("Health check passed");
    }).catch(() => {
      console.log("Health check failed");
    });
  }

  // TODO: re-enable healthcheck
  // const interval = setInterval(healthCheck, 1000);
  // onDestroy(() => clearInterval(interval));

  const initialize = () => {
    // return $keystores;

    getKeystores().then((response) => {

      onboarding = response.length == 0;
      login = false

      if (response.length > 0) {
        //
        // keystores = response.sort((a, b) => {
        //   return a.id < b.id ? 1 : -1;
        // });

        keystore = response[0];
      }

    }).catch((error: Response) => {
      if (error.status == 401) {
        login = true;
        return
      }

      console.log(error)
    }).finally(() => {
      ready = true;
      getCurrentUser()
    });

  }

  onMount(() => {
    initialize()
  });
</script>

<Router>
  {#if !ready}
    <span>Loading...</span>
  {:else}
    {#if onboarding}
      <OnboardingForm on:created={initialize}/>
    {:else if login}
      <LoginForm on:login={initialize}/>
    {:else if $currentUser}
      <Sidebar keystores={$keystores} bind:showCreateKeystoreModal={showCreateKeystoreModal}/>
      <main>
        <Route>
          <Keystores keystores={$keystores}/>
        </Route>

        <Route path="/keystore/:keystoreId">
          <Keystores keystores={$keystores}/>
        </Route>

        <Route path="/keystore/:keystoreId/entry/:entryId">
          <Keystores keystores={$keystores}/>
        </Route>

        <Route path="/invitations">
          <Invitations/>
        </Route>
      </main>
    {/if}
  {/if}
</Router>

<Messages/>

<CreateKeystoreModal bind:show={showCreateKeystoreModal} on:created={() => {getKeystores()}}/>

<style lang="scss">
  $bg: #0a0e11;
  //$bg: #111519;
  $accent: #00edb1;
  $text-color: #fff;

  @import url('https://rsms.me/inter/inter.css');

  :global {

    #app {
      height: 100%;
      display: flex;
      flex-direction: row;
      flex-grow: 1;
    }

    body {
      margin: 0;
      background-color: $bg;

      height: 100vh;
      max-height: 100vh;
      //max-height: -webkit-fill-available;
      //max-height: 100%;
      overflow-x: hidden;
    }

    * {
      font-family: 'Inter', sans-serif;
      font-weight: 300;
      font-size: 100%;
      //color: $text-color;
      line-height: 1.4;
    }

    //::-webkit-scrollbar{
    //	width: 4px;
    //}
    //
    //::-webkit-scrollbar-track-piece{
    //	background-color: transparent;
    //}
    //
    //::-webkit-scrollbar-thumb{
    //	background-color: #363a41;
    //	border-radius: 2px;
    //}
    //
    //::-webkit-scrollbar-thumb:hover{
    //	background-color: #909090;
    //}
  }

  :root {
    color-scheme: dark;
  }

  main {
    display: flex;
    align-items: stretch;
    background-color: $bg;
    width: 100%;
  }
</style>
