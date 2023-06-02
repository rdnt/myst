<script lang="ts">
  import LoginForm from "@/components/LoginForm.svelte";
  import Messages from "@/components/Messages.svelte";
  import OnboardingForm from "@/components/OnboardingForm.svelte";
  import Main from "@/pages/Main.svelte";
  import {getKeystores, keystores} from "@/stores/keystores.ts";
  import {onMount} from 'svelte';
  import {Router, useNavigate} from "svelte-navigator";
  import {Link, useLocation} from "svelte-navigator";
  import api from "@/api";

  const location = useLocation();
  const navigate = useNavigate();

  let onboarding = false;
  let ready = false;
  let login = false;

  const initialize = async () => {
    onboarding = false
    login = false

    try {
      await api.enclave().catch(err => {
        if (err.status == 404) {
          onboarding = true;
          return Promise.resolve()
        } else if (err.status == 401) {
          login = true;
          return Promise.resolve()
        } else {
          return Promise.reject(err)
        }
      })

      if ((onboarding || login) && $location.pathname !== '/') {
        navigate('/', {replace: true})
        await initialize()
      }

      ready = true
      // getCurrentUser()
      // ready = true;
    } catch (err) {
      console.log('caught error', err)
    }
  }

  onMount(async () => {
    await initialize()
  });
</script>

<Router>
  {#if !ready}
<!--    <span>Loading...</span>-->
  {:else}
    {#if onboarding}
      <OnboardingForm on:initialized={initialize}/>
    {:else if login}
      <LoginForm on:login={async () => {
        onboarding = false
        login = false
      }}/>
    {:else}
      <Main/>
    {/if}
  {/if}
</Router>

<Messages/>


<style lang="scss">
  $bg: #0a0e11;
  //$bg: #111519;
  $accent: #00edb1;
  $text-color: #fff;

  @import url('//rsms.me/inter/inter.css');

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
</style>
