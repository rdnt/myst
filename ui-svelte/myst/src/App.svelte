<script lang="ts">
  import {Route, Router} from "svelte-navigator";
  import Sidebar from "./components/Sidebar.svelte";
  import Keystores from "./pages/Keystores.svelte";
  import {onDestroy, onMount} from 'svelte';
  import api from "./api";
  import LoginForm from "./components/LoginForm.svelte";
  import OnboardingForm from "./components/OnboardingForm.svelte";

  let onboarding = false;
  let ready = false;
  let login = false;

  let keystores = [];
  let keystore = null;

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

  function getKeystores() {
    api.keystores().then((response) => {
      onboarding = response.length == 0;
      login = false

      if (response.length > 0) {

        keystores = response.sort((a, b) => {
          return a.id < b.id ? 1 : -1;
        });

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
    });
  }

  onMount(async () => {
    getKeystores()
  });
</script>

<Router>
  {#if !ready}
    <span>Loading...</span>
  {:else}
    {#if onboarding}
      <OnboardingForm on:created={getKeystores} />
    {:else if login}
      <LoginForm on:login={getKeystores}/>
    {:else}
      <Sidebar {keystores}/>
      <main>
        <Route>
          <Keystores {keystores}/>
        </Route>

        <Route path="/keystore/:keystoreId">
          <Keystores {keystores}/>
        </Route>

        <Route path="/keystore/:keystoreId/entry/:entryId">
          <Keystores {keystores}/>
        </Route>
      </main>
    {/if}
  {/if}
</Router>

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
    }

    * {
      font-family: 'Inter', sans-serif;
      font-weight: 300;
      font-size: 100%;
      color: $text-color;
      line-height: 1.4;
    }
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
