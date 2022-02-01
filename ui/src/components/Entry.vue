<template>
	<div class="entry" :class="{empty: !entry}">
		<template v-if="entry">
			<div class="field">
<!--				<label>Domain</label>-->
				<ExpandingTextarea v-model="domain"></ExpandingTextarea>
			</div>
			<div class="field">
<!--				<label>Username</label>-->
				<ExpandingTextarea v-model="username" ></ExpandingTextarea>
			</div>
			<div class="field">
<!--				<label>Password</label>-->
				<ExpandingTextarea v-model="password"></ExpandingTextarea>
			</div>
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
		domain: 'sddsasadsad',
		username: 'someuseadad',
		password:'dsdsasdaasd'
	}),
	watch: {
		entry(entry: Entry) {
			console.log('entry changed', entry)
			this.domain = entry.label
			this.username = entry.username
			this.password = entry.password
		}
	},
	computed: {},
	methods: {}
})
</script>

<style lang="scss" scoped>
.entry {
	display: flex;
	flex-direction: column;
	background-color: #101519;
	border-left: 2px solid #1a2025;
	height: 100%;
	padding: 20px;
	box-sizing: border-box;
	flex-basis: 40%;
	padding-top: 100px;

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
}

.field {
	display: flex;
	flex-direction: column;
	margin-bottom: 30px;

	label {
		font-size: 1.1rem;
		height: 30px;
		display: block;
		padding: 0 15px;
	}

	textarea {
		display: block;
		margin: 0;
		border: none;
		outline: none;
		width: 100%;
		resize: none;
		font-size: 1.1rem;
		font-weight: 400;
		box-sizing: border-box;
		background-color: rgba(#abc, .05);
		border-radius: 5px;
		padding: 15px 16px;
		color: #fff;
		overflow: hidden;

		&::placeholder {
			color: lighten(#68737e, 5%);
		}

		&:focus {
			&::placeholder {
				color: lighten(#68737e, 15%);
			}
		}

		&:disabled {
			//background-color: transparent;
			//padding-left: 0;
			//padding-right: 0;

		}
	}
}
</style>
