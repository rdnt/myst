<script lang="ts">
    import LoginForm from "@/components/LoginForm.svelte";
    import Messages from "@/components/Messages.svelte";
    import OnboardingForm from "@/components/OnboardingForm.svelte";
    import Main from "@/pages/Main.svelte";
    import {onDestroy, onMount} from 'svelte';
    import {Route, Router, useLocation, useNavigate} from "svelte-navigator";
    import api from "@/api";
    import {AuthState, authState, getAuthState} from "@/stores/authState";
    import {showError} from "@/stores/messages";

    const location = useLocation();
    const navigate = useNavigate();

    $: {
        const state = $authState;

        if ((state === AuthState.Onboarding || state === AuthState.SignedOut) && $location.pathname !== '/') {
            navigate('/', {replace: true})
        }
    }

    const checkState = async () => {
        const oldState = $authState;
        const state = await getAuthState();
        if (state === AuthState.SignedIn) {
            await api.healthCheck()
        } else if (state === AuthState.SignedOut && oldState === AuthState.SignedIn) {
            // was logged out, notify
            showError("Signed out", "You have been signed out due to inactivity.");
        }
    }

    const interval = setInterval(checkState, 20000);

    onDestroy(() => clearInterval(interval));

    onMount(async () => {
        await checkState()
    });
</script>

{#if $authState === undefined}
    <!--        <span>Loading...</span>-->
{:else}
    {#if $authState === AuthState.Onboarding}
        <OnboardingForm on:initialized={async () => {
            await checkState()
        }}/>
    {:else if $authState === AuthState.SignedOut}
        <LoginForm on:login={async () => {
            await checkState()
        } }/>
    {:else if $authState === AuthState.SignedIn}
        <Main/>
    {/if}
{/if}

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
