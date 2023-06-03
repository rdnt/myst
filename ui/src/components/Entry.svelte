<script lang="ts">
  import type {Entry, Keystore} from "@/api";
  import api from "@/api";
  import DeleteEntryModal from "@/components/DeleteEntryModal.svelte";
  import EditEntryModal from "@/components/EditEntryModal.svelte";
  import Field from "@/components/Field.svelte";
  import {getKeystores} from "@/stores/keystores";

  let showDeleteModal = false;
  let showEditModal = false;

  const updateEntry = async (password, notes: string) => {
    api.updateEntry({
      keystoreId: keystore.id,
      entryId: entry.id,
      requestBody: {
        password: password != entry.password ? password : undefined,
        notes: notes != entry.notes ? notes : undefined,
      }
    }).then((res) => {
      if (res) {
        entry.password = password;
        entry.notes = notes;
        showEditModal = false
      }
    });
  }

  export let keystore: Keystore;
  export let entry: Entry;

  const deleteEntry = async () => {
    await api.deleteEntry({
      keystoreId: keystore.id,
      entryId: entry.id
    }).then((res) => {
      if (res) {
        showDeleteModal = false
      }
    });

    await getKeystores()
  }
</script>

<!--<div class="entry">-->
<!--  {JSON.stringify(entry)}-->
<!--</div>-->

<div class="entry" class:empty={!entry}>
  {#if entry}

    <div class="header">
      <div class="image">
        <img alt={entry.website}
             src="https://www.nicepng.com/png/full/52-520535_free-files-github-github-icon-png-white.png">
      </div>
      <div class="title">
        <h5>{entry.website}</h5>
        <a>{entry.username}</a>
      </div>
      {#if !keystore.readOnly}
        <div class="actions">
          <button on:click={() => showEditModal = true} class="button"><span class="icon"></span> Edit</button>
          <button on:click={() => showDeleteModal = true} class="button red"><span class="icon"></span> Delete</button>
        </div>
      {/if}
    </div>
    <div class="fields">
      <Field label="Website" value={entry.website}/>
      <Field label="Username" value={entry.username}/>
      <Field label="Password" value={entry.password}/>
      <Field label="Notes" value={entry.notes}/>
      <!--    <ExpandingTextarea :disabled="true" :value="website" label="Website"></ExpandingTextarea>-->
      <!--    <ExpandingTextarea :disabled="true" :value="username" label="Email Address"></ExpandingTextarea>-->
      <!--    <ExpandingTextarea :disabled="true" :value="password" label="Password"></ExpandingTextarea>-->
      <!--    <ExpandingTextarea :disabled="true" :placeholder="!notes ? '—' : ''" :value="notes"-->
      <!--                       label="Notes"></ExpandingTextarea>-->
    </div>
  {/if}
</div>

<EditEntryModal bind:show={showEditModal} {entry} on:submit={(e) => {updateEntry(e.detail.password, e.detail.notes)}}/>

<!--    <Modal show={showEditModal}>-->
<!--      <div slot="header" class="modal-footer">-->
<!--        <div class="image">-->
<!--          <img src="https://www.nicepng.com/png/full/52-520535_free-files-github-github-icon-png-white.png" alt={entry.website}/>-->
<!--        </div>-->
<!--        <div class="title">-->
<!--          <h5>{ entry.website }</h5>-->
<!--          <a>{ entry.username }</a>-->
<!--        </div>-->
<!--      </div>-->
<!--      <div class="modal-footer" slot="footer">-->
<!--        <button class="button transparent" on:click={() => showEditModal = false} >Cancel</button>-->
<!--        <button class="button red">Delete</button>-->
<!--      </div>-->
<!--    </Modal>-->

<!--    <modal :show="showEditModal" :setShow="() => {this.showEditModal = false}">-->
<!--      <template v-slot:header>-->
<!--        &lt;!&ndash;					<div class="header modal-header">&ndash;&gt;-->
<!--        &lt;!&ndash;						<div class="image">&ndash;&gt;-->
<!--        &lt;!&ndash;							<img :alt="entry.website"&ndash;&gt;-->
<!--        &lt;!&ndash;									 src="https://www.nicepng.com/png/full/52-520535_free-files-github-github-icon-png-white.png">&ndash;&gt;-->
<!--        &lt;!&ndash;						</div>&ndash;&gt;-->
<!--        &lt;!&ndash;						<div class="title">&ndash;&gt;-->
<!--        &lt;!&ndash;							<h5>{{ entry.website }}</h5>&ndash;&gt;-->
<!--        &lt;!&ndash;							<a>{{ entry.username }}</a>&ndash;&gt;-->
<!--        &lt;!&ndash;						</div>&ndash;&gt;-->
<!--        &lt;!&ndash;					</div>&ndash;&gt;-->
<!--      </template>-->
<!--      <div class="fields">-->
<!--        &lt;!&ndash;					<ExpandingTextarea2 v-model="password"></ExpandingTextarea2>&ndash;&gt;-->
<!--        <entry-text-field v-model="password" label="Password"></entry-text-field>-->
<!--        <ExpandingTextarea :value="website" :disabled="true" label="Website"></ExpandingTextarea>-->
<!--        <ExpandingTextarea :value="username" :disabled="true" label="Email Address"></ExpandingTextarea>-->
<!--        <ExpandingTextarea v-model="password" label="Password"></ExpandingTextarea>-->
<!--        <ExpandingTextarea v-model="notes" :placeholder="!edit && !notes ? '—' : ''"-->
<!--                           label="Notes"></ExpandingTextarea>-->
<!--      </div>-->
<!--      <template v-slot:footer>-->
<!--        <div class="modal-footer">-->
<!--          <button class="button transparent" @click="showEditModal = false">Cancel</button>-->
<!--          <button class="button green">Save</button>-->
<!--        </div>-->
<!--      </template>-->
<!--    </modal>-->

<!--    <Modal bind:show={showDeleteModal}>-->
<!--        <div slot="header" class="delete-title">Are you sure you want to delete this entry?</div>-->
<!--        <div class="modal-footer" slot="footer">-->
<!--          <button class="button transparent" on:click={() => showDeleteModal = false} >Cancel</button>-->
<!--          <button class="button red">Delete</button>-->
<!--        </div>-->
<!--    </Modal>-->

<DeleteEntryModal bind:show={showDeleteModal} {entry} on:submit={() => {deleteEntry()}}/>
<!--</div>-->

<style lang="scss">
  $accent: #00edb1;

  h5 {
    font-weight: 600;
    font-size: 1.8rem;
    margin: 0;
  }

  .entry {
    display: flex;
    flex-direction: column;
    background-color: #101519;
    border-left: 2px solid #1a2025;
    height: 100%;
    padding: 20px;
    box-sizing: border-box;
    flex-basis: 50%;
    overflow-y: auto;
    padding: 32px;
    //flex-grow: 1;
    //padding-top: 100px;

    &.empty {
      &:after {
        content: "";
        background-color: #1b2025;
        border-radius: 50%;
        font-size: 1.5em;
        font-weight: bold;
        text-align: center;
        display: block;
        width: 100px;
        height: 100px;
        line-height: 100%;
      }
    }

    .actions {
      display: flex;
      //width: 100%;
      flex-direction: row;
      align-items: center;
      margin-left: auto;
      height: 100%;

      &.bottom {
        margin: 16px 6px;
        margin-top: auto;
        justify-content: flex-start;
      }

      .button {
        display: flex;
        align-items: center;

        .icon {
          display: block;
          width: 20px;
          height: 20px;
          background-color: rgba(#c0c2c8, 1);
          border-radius: 8px;
          margin-right: 10px;
        }

        &.red .icon {
          background-color: #ff9999;
        }
      }
    }

    .header {
      display: flex;
      flex-direction: row;
      align-items: center;
      margin-bottom: 40px;
      //padding: 16px 0;
      //padding: 16px 0 16px 16px;

      .image {
        width: 64px;
        height: 64px;
        padding-right: 20px;
        display: none;

        img {
          width: 64px;
          height: 64px;
        }
      }

      .title {
        display: flex;
        flex-direction: column;
        flex-grow: 1;

        a {
          padding: 5px 0;
          overflow: hidden;
          word-break: break-all;
          text-overflow: ellipsis;
        }

      }

      .button {
        align-self: flex-start;
      }
    }

    .fields {
      //padding: 16px 0;
    }
  }

  .field {

    &.disabled {
      //margin-left: -16px;
    }

    //margin-bottom: 2px;

    label {
      //font-size: 1.1rem;
      //height: 30px;
      //display: block;
      //padding: 0 15px;
    }

    //textarea {
    //	display: block;
    //	margin: 0;
    //	border: none;
    //	outline: none;
    //	width: 100%;
    //	resize: none;
    //	font-size: 1.1rem;
    //	font-weight: 400;
    //	box-sizing: border-box;
    //	//background-color: rgba(#abc, .05);
    //	padding: 15px 16px;
    //	color: #fff;
    //	overflow: hidden;
    //
    //	&::placeholder {
    //		color: lighten(#68737e, 5%);
    //	}
    //
    //	&:focus {
    //		&::placeholder {
    //			color: lighten(#68737e, 15%);
    //		}
    //	}
    //
    //	&:disabled {
    //		//background-color: transparent;
    //		//padding-left: 0;
    //		//padding-right: 0;
    //
    //	}
    //}
  }

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

  .modal-header {
    //margin-bottom: 22px;
  }

  .modal-footer {
    display: flex;
    flex-direction: row;
    justify-content: flex-end;
    margin-top: 22px;
  }

  .delete-title {
    padding: 4px;
    box-sizing: border-box;
    font-size: 1.1rem;
  }

  .fields {
    margin-bottom: 22px;
  }
</style>
