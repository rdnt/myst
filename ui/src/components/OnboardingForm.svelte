<script lang="ts">
  import api from "@/api";
  import InputField from "@/components/InputField.svelte";
  import {createEventDispatcher, onMount} from 'svelte';

  const dispatch = createEventDispatcher();

  let warnings: boolean = false;

  let step: number = 1;

  let name: string = '';
  $: nameValid = name.length > 0;

  let password: string = '';
  $: passwordValid = password.length >= 8;

  $: valid = (step === 1 && nameValid) || (step === 2 && passwordValid)

  let nameInput: HTMLInputElement;
  let passwordInput: HTMLInputElement;

  onMount(() => nameInput.focus())

  const submit = () => {
    if (!valid) {
      return;
    }

    if (step === 1) {
      step = 2
      warnings = false
      setTimeout(() => passwordInput.focus(), 0)
    } else if (step === 2) {
      api.createKeystore({
        createKeystoreRequest: {
          name,
          password,
        }
      }).then((keystore) => {
        dispatch('created', keystore)
      }).catch((err) => {
        console.error(err)
      })
    }
  }
</script>

<form class="form" on:submit|preventDefault={submit}>
  <div class="content">
    <h4>Myst</h4>

    <h5>First time setup</h5>
    <div class="separator"></div>

    <div class="step-1">
      {#if (step === 1)}
        <InputField
          bind:value={name}
          class="name-input"
          error={!nameValid && warnings ? 'Cannot be empty' : ''}
          label="Let's choose a name for your first keystore. Keystores help organize your secrets into groups."
          placeholder="Keystore Name"
          readonly={step !== 1}
          on:input={() => warnings = true}
          bind:ref={nameInput}
        />
      {:else}
        <h6>Selected keystore name</h6>
        <span class="selected-name">{name}</span>
      {/if}
    </div>

    {#if step === 2}
      <div class="step-2">
        <InputField
          bind:value={password}
          class="password-input"
          error={!passwordValid && warnings ? 'Password too weak' : ''}
          readonly={step !== 1}
          label="Your keystores will be encrypted using a master password. The security of your secrets will
          depend on its complexity. Choose one wisely and make sure you remember it."
          placeholder="Master password"
          on:input={() => warnings = true}
          bind:ref={passwordInput}
        />
      </div>
    {/if}

    <div class="footer">
      <span class="step-label">Step {step} of 2</span>
      {#if step === 1}
        <button class:disabled={!valid} class="button" type="submit">Next</button>
      {:else}
        <button class:disabled={!valid} class="button green" type="submit">Create</button>
      {/if}
    </div>

  </div>
</form>

<style lang="scss">
  $bg: #0a0e11;
  $accent: #00edb1;
  //$text-color: #f4f8fb;
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

  .selected-name {
    position: relative;
    background: transparent;
    font-weight: 600;
    color: $accent;
    padding: 0;
    font-size: 1.4rem;
    pointer-events: none;
    margin-bottom: 4px;
  }

  .separator {
    width: 50px;
    height: 2px;
    background-color: #31363d;
    margin: 30px auto;
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

      h4 {
        font-size: 2.4rem;
        text-align: center;
        font-weight: 700;
        margin: 0 0 10px;
      }

      h5 {
        font-weight: 400;
        font-size: 1.4rem;
        text-align: center;
        margin: 0;
      }

      h6 {
        font-size: 1.1rem;
        font-weight: 300;
        margin: 0;
      }

      .step-1 {
        margin-top: 100px;
        margin-bottom: 70px;
      }

      .footer {
        display: flex;
        flex-direction: column;
        align-items: center;
        margin-top: 70px;

        .step-label {
          font-size: .9rem;
          opacity: .95;
          margin-bottom: 12px;
        }
      }
    }
  }
</style>
