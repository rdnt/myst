<script lang="ts">
  import api from "@/api";
  import InputField from "@/components/InputField.svelte";
  import {createEventDispatcher} from 'svelte';

  const dispatch = createEventDispatcher();

  let password: string = '12345678';
  $: passwordValid = password.length >= 8;

  let error: boolean = false;

  function submit() {
    if (!passwordValid) {
      return;
    }

    api.authenticate({
      requestBody: {
        password
      }
    }).then(() => {
      dispatch('login')
      password = ''
    }).catch((err) => {
      if (err.status == 401) {
        error = true
        password = ''

        return
      }

      console.error(err)
    })
  }
</script>

<form class="form" on:submit|preventDefault={submit}>
  <div class="content">
    <h4>Myst</h4>
    <InputField
      bind:value={password}
      class="password-input"
      error={error ? 'Incorrect Password' : ''}
      label="Use your master password to access your secrets."
      placeholder="Master Password"
    />
    <div class="footer">
      <button class="button green" type="submit" class:disabled={!passwordValid}>Login</button>
    </div>
  </div>
</form>

<style lang="scss">
  $bg: #0a0e11;
  $accent: #00edb1;
  $text-color: #fff;

  .button {
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

  .form {
    width: 100%;
    height: 100vh;
    display: flex;
    flex-direction: row;
    align-items: center;
    justify-content: center;
    background: $bg;

    .content {
      position: fixed;
      color: #e1e4e8;
      width: 500px;
      height: 400px;

      h4 {
        font-weight: 700;
        font-size: 3.2rem;
        text-align: left;
        margin: 0;
        text-align: center;
      }

      h6 {
        font-size: 1.1rem;
        font-weight: 300;
        margin: 0;
      }

      :global(.password-input) {
        padding: 90px 0 70px;
      }

      .footer {
        display: flex;
        flex-direction: column;
        align-items: center;
      }
    }
  }
</style>

