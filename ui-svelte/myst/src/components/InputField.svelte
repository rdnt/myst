<script lang="ts">
  export let ref = null;
  export let label: string = '';
  export let error: string = '';
  export let placeholder: string = '';
  export let value: string;

  import {createEventDispatcher} from 'svelte';

  const dispatch = createEventDispatcher();
</script>

<div class={$$props.class} class:field={true}>
  {#if label}
    <label>{label}</label>
  {/if}
  <input bind:this={ref} bind:value={value} on:input={dispatch('input', value)} {placeholder}/>
  <span class:show={error}>{error}</span>
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
      padding: 14px 16px;
      overflow: hidden;

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
      padding: 0 16px;
      opacity: 0;
      pointer-events: none;
      margin-top: 12px;
      min-height: 20px;

      &.show {
        opacity: 1;
        pointer-events: unset;
      }
    }
  }
</style>
