<script lang="ts">
  import api from "@/api";
  import InputField from "@/components/InputField.svelte";
  import Modal from "@/components/Modal.svelte";
  import {showError, showMessage} from "@/stores/messages";
  // import {getCurrentUser} from "@/stores/user";

  export let show: boolean = false;

  let username, password: string;

  $: usernameValid = username.trim() !== '';
  $: passwordValid = password.trim() !== '';

  $: allowSubmit = usernameValid && passwordValid;

  $: {
    if (!show) {
      username = ''
      password = ''
    }
  }

  const reset = () => {
    username = '';
    password = ''
  }

  const submit = () => {
    if (!allowSubmit) {
      return;
    }

    api.register({
      requestBody: {
        username,
        password,
      }
    }).then(() => {
      showMessage("Signed in.");
      reset()
      getCurrentUser()
      show = false;
    }).catch((err) => {
      showError("Signing in failed.");
      console.error(err)
    })
  };
</script>

<form class="create-entry-modal" on:submit|preventDefault={submit}>
  <Modal bind:show>
    <div class="create-title" slot="header">Register</div>

    <div class="modal-content">
      <InputField bind:value={username} label="Username"/>
      <InputField bind:value={password} label="Password"/>
    </div>

    <div class="modal-footer" slot="footer">
      <button class="button transparent" on:click={() => show = false} type="button">Cancel</button>
      <button class="button green" type="submit">Register</button>
    </div>
  </Modal>
</form>

<style lang="scss">
  .create-title {
    padding: 4px;
    box-sizing: border-box;
    font-size: 1.1rem;
  }

  .create-entry-modal {

    .modal-header {
      display: flex;
      flex-direction: row;
      //margin-top: 10px;

      .image {
        width: 64px;
        height: 64px;
        padding-right: 20px;

        img {
          width: 64px;
          height: 64px;
        }
      }

      .title {
        display: flex;
        flex-direction: column;
        flex-grow: 1;

        .website {
          font-weight: 600;
          font-size: 1.8rem;
          margin: 0;
        }

        .username {
          //padding: 5px 0;
        }

      }
    }


    .modal-content {
      padding-top: 40px;
      box-sizing: border-box;
      box-sizing: border-box;
    }

    .modal-footer {
      display: flex;
      flex-direction: row;
      justify-content: flex-end;
      margin-top: 12px;
      box-sizing: border-box;
    }

  }


  $accent: #00edb1;
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
