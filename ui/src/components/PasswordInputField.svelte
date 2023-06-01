<script lang="ts">
  import {createEventDispatcher} from 'svelte';

  export let ref = null;
  export let label: string = '';
  export let error: string | undefined = undefined;
  export let placeholder: string = '';
  export let value: string;

  const dispatch = createEventDispatcher();
</script>

<div class={$$props.class} class:field={true}>
  {#if label}
    <label>{label}</label>
  {/if}
  <input bind:this={ref} bind:value={value} on:input={dispatch('input', value)} {placeholder} class:has-error={error} type="password" autocomplete="dont-autocomplete" />
  {#if error}
    <span>{error}</span>
  {/if}
</div>

<style lang="scss">
  .field {
    display: flex;
    flex-direction: row;
    flex-wrap: wrap;

    label {
      flex-basis: 100%;
      margin-bottom: 12px;
      font-size: 1.1rem;
      //padding: 0 14px;

      //color: rgba(#8a8f9f, .75);
      //text-transform: uppercase;
      //font-size: 0.85rem;
      //font-weight: 600;
      //letter-spacing: 0.5px;
      //pointer-events: none;
    }

    input {
      margin: 0;
      border: 0;
      color: #fff;
      outline: none;
      width: 100%;
      display: block;
      font-size: 1.1rem;
      font-weight: 400;
      box-sizing: border-box;
      background-color: rgba(#abc, .05);
      border-radius: 5px;
      padding: 14px;
      overflow: hidden;
      margin-bottom: 32px;

      &.has-error {
        margin-bottom: 0;
      }

      &::placeholder {
        color: lighten(#68737e, 5%);
      }

      &:focus {

        &::placeholder {
          color: lighten(#68737e, 15%);
        }
      }
    }

    span {
      color: #ff9999;
      font-weight: 500;
      font-size: .9rem;
      min-height: 20px;
      margin-top: 12px;
      margin-bottom: 32px;

      &.show {
        display: block;
      }
    }
  }
</style>
