<template>
	<div :class="{empty: !entry}" class="entry">
		<template v-if="entry">
			<div class="actions">
				<button class="button" @click="showEditModal = true">Edit</button>
				<button class="button red" @click="showDeleteModal = true">Delete</button>
			</div>
			<div class="header">
				<div class="image">
					<img :alt="entry.website"
							 src="https://www.nicepng.com/png/full/52-520535_free-files-github-github-icon-png-white.png">
				</div>
				<div class="title">
					<h5>{{ entry.website }}</h5>
					<a>{{ entry.username }}</a>
				</div>
			</div>
			<div class="separator"/>
			<div class="fields">
				<ExpandingTextarea :value="website" :disabled="true" label="Website"></ExpandingTextarea>
				<ExpandingTextarea :value="username" :disabled="true" label="Email Address"></ExpandingTextarea>
				<ExpandingTextarea :value="password" :disabled="true" label="Password"></ExpandingTextarea>
				<ExpandingTextarea :value="notes" :disabled="true" :placeholder="!notes ? '—' : ''"
													 label="Notes"></ExpandingTextarea>
			</div>

			<modal :show="showEditModal" :setShow="() => {this.showEditModal = false}">
				<template v-slot:header>
<!--					<div class="header modal-header">-->
<!--						<div class="image">-->
<!--							<img :alt="entry.website"-->
<!--									 src="https://www.nicepng.com/png/full/52-520535_free-files-github-github-icon-png-white.png">-->
<!--						</div>-->
<!--						<div class="title">-->
<!--							<h5>{{ entry.website }}</h5>-->
<!--							<a>{{ entry.username }}</a>-->
<!--						</div>-->
<!--					</div>-->
				</template>
				<div class="fields">
<!--					<ExpandingTextarea2 v-model="password"></ExpandingTextarea2>-->
					<entry-text-field v-model="password" label="Password"></entry-text-field>
					<ExpandingTextarea :value="website" :disabled="true" label="Website"></ExpandingTextarea>
					<ExpandingTextarea :value="username" :disabled="true" label="Email Address"></ExpandingTextarea>
					<ExpandingTextarea v-model="password" label="Password"></ExpandingTextarea>
					<ExpandingTextarea v-model="notes" :placeholder="!edit && !notes ? '—' : ''"
														 label="Notes"></ExpandingTextarea>
				</div>
				<template v-slot:footer>
					<div class="modal-footer">
						<button class="button transparent" @click="showEditModal = false">Cancel</button>
						<button class="button green">Save</button>
					</div>
				</template>
			</modal>

			<modal :show="showDeleteModal" :setShow="() => {this.showDeleteModal = false}">
				<template v-slot:header>
					<div class="delete-title">Are you sure you want to delete this entry?</div>
				</template>
				<template v-slot:footer>
					<div class="modal-footer">
						<button class="button transparent" @click="showDeleteModal = false">Cancel</button>
						<button class="button red">Delete</button>
					</div>
				</template>
			</modal>

		</template>
	</div>
</template>

<script lang="ts">
import {defineComponent} from 'vue'
import {Entry} from "../api/generated";
import ExpandingTextarea from "../components/ExpandingTextarea.vue";
import Modal from "../components/Modal.vue";
import ExpandingTextarea2 from "../components/ExpandingTextarea2.vue";
import EntryTextField from "../components/EntryTextField.vue";

export default defineComponent({
	name: 'Entry',
	components: {EntryTextField, ExpandingTextarea, Modal, ExpandingTextarea2},
	props: {
		entry: {
			type: Object as () => Entry,
			required: false
		}
	},
	data: () => ({
		edit: false,
		website: 'sddsasadsad',
		username: 'someuseadad',
		password: 'dsdsasdaasd',
		notes: '',
		showEditModal: false,
		showDeleteModal: false,
	}),
	watch: {
		entry(entry: Entry) {
			this.website = entry.website
			this.username = entry.username
			this.password = entry.password
			this.notes = entry.notes || '—'
		}
	},
	computed: {},
	methods: {}
})
</script>

<style lang="scss" scoped>
$accent: #00edb1;

.separator {
	//width: calc(100% - 32px);
	//height: 1px;
	//background-color: #1b2025;
	//margin: 10px auto 20px;
}

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
	flex-basis: 40%;
	overflow-y: auto;
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
		justify-content: flex-end;
		align-items: center;
		align-self: flex-start;
		margin-left: auto;

		&.bottom {
			margin: 16px 6px;
			margin-top: auto;
			justify-content: flex-start;
		}
	}

	.header {
		display: flex;
		flex-direction: row;
		align-items: center;
		padding: 16px;
		//padding: 16px 0 16px 16px;

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

			a {
				padding: 5px 0;
			}

		}

		.button {
			align-self: flex-start;
		}
	}

	.fields {
		//padding: 20px 0;
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
	padding: 0 20px;
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
