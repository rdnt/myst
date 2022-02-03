<template>
	<div class="entry" :class="{empty: !entry}">
		<template v-if="entry">
			<div class="actions">
				<button class="button">Edit</button>
				<button class="button">Delete</button>
			</div>
			<div class="header">
				<div class="image">
					<img :alt="entry.website" src="https://www.nicepng.com/png/full/52-520535_free-files-github-github-icon-png-white.png">
				</div>
				<div class="title">
					<h5>{{ entry.website }}</h5>
					<a>Login</a>
				</div>
			</div>
			<div class="separator"/>
			<ExpandingTextarea v-model="website" label="Website" :disabled="true"></ExpandingTextarea>
			<ExpandingTextarea v-model="username" label="Email Address" :disabled="true"></ExpandingTextarea>
			<ExpandingTextarea v-model="password" label="Password" :disabled="true"></ExpandingTextarea>
			<ExpandingTextarea v-model="notes" label="Notes" :disabled="false"></ExpandingTextarea>
		</template>
	</div>
</template>

<script lang="ts">
import {defineComponent} from 'vue'
import {Entry} from "../api/generated";
import ExpandingTextarea from "../components/ExpandingTextarea.vue";

export default defineComponent({
	name: 'Entry',
	components: {ExpandingTextarea},
	props: {
		entry: {
			type: Object as () => Entry,
			required: false
		}
	},
	data: () => ({
		website: 'sddsasadsad',
		username: 'someuseadad',
		password:'dsdsasdaasd',
		notes: '',
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
	width: calc(100% - 32px);
	height: 2px;
	background-color: #1b2025;
	margin: 10px auto 20px;
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
		flex-direction: row;
		justify-content: flex-end;
		align-items: center;
		//padding: 0 16px;
	}

	.header {
		display: flex;
		flex-direction: row;
		align-items: center;
		padding: 16px 16px;

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

			a {
				padding: 5px 0;
			}
		}
	}

	.fields {
		padding: 20px 0;
	}
}

.field {
	margin-bottom: 2px;

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
</style>
