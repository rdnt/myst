<script lang="ts">
  import Modal from "./Modal.svelte";
  import * as models from "../api/generated/models";
  import InputField from "./InputField.svelte";
  import Field from "./Field.svelte";
  import {createEventDispatcher, onMount} from 'svelte';
  const dispatch = createEventDispatcher();

  export let show: boolean = false;
  export let entry: models.Entry;

  $: password = '';
  $: notes = '';

  $: passwordValid = password.trim() !== '';
  $: passwordChanged = entry && password.trim() !== entry.password;
  $: notesChanged = entry && notes.trim() !== entry.notes;

  $: allowSubmit = passwordValid && (passwordChanged || notesChanged);

  $: {
    if (!show && entry) {
      password = entry.password
      notes = entry.notes
    }
  }

  onMount(async () => {
    password = entry.password
    notes = entry.notes
  });

  const submit = () => {
    if (!allowSubmit) {
      return;
    }

    dispatch('submit', {
      password,
      notes,
    });
  };
</script>

<form class="edit-entry-modal" on:submit|preventDefault={submit}>
  <Modal bind:show>
    <div class="modal-header" slot="header">
      <div class="image">
        <img alt={entry.website}
             src="https://www.nicepng.com/png/full/52-520535_free-files-github-github-icon-png-white.png"/>
      </div>

      <div class="title">
        <span class="website">{entry.website}</span>
        <span class="username">{entry.username}</span>
      </div>
    </div>

    <div class="modal-content">
      <Field label="Website" value={entry.website}/>
      <Field label="Username" value={entry.username}/>
      <InputField bind:value={password} error={!passwordValid && 'Password cannot be empty'} label="Password"/>
      <InputField bind:value={notes} label="Notes"/>
    </div>

    <div class="modal-footer" slot="footer">
      <button class="button transparent" on:click={() => show = false} type="button">Cancel</button>
      <button class="button green" type="submit">Save Changes</button>
    </div>
  </Modal>
</form>

<style lang="scss">
  .edit-entry-modal {

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
