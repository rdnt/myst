<template>
	<form class="login-form" @submit.prevent="submit">
		<div class="login-form-content">
			<h4>Myst</h4>
			<br/><br/><br/><br/>
			<h6>Use your master password to access your secrets.</h6>
			<br/>
			<div class="input">
				<label>Master password</label>
				<expanding-textarea ref="password" v-model="password" class="password"
														placeholder="Master password"
														@keydown.enter.prevent="submit"/>
			</div>
			<br/>
			<div class="login-form-footer">
				<button :class="{disabled: password === ''}" class="button green" type="submit">Login</button>
			</div>
		</div>
	</form>
</template>

<script lang="ts">
import {defineComponent} from 'vue'
import ExpandingTextarea from './ExpandingTextarea.vue'
import api from "../api";

export default defineComponent({
	name: 'Login',
	components: {ExpandingTextarea},
	data: () => ({
		password: '',
	}),
	computed: {
		valid() {
			return this.passwordValid
		},
		passwordValid() {
			if (this.password.length == 0) {
				return false
			} else if (this.password.length < 8) {
				return false
			}

			return true
		}
	},
	methods: {
		submit() {
			if (this.passwordValid) {

				// TODO: unlock all keystores at once
				// api.createKeystore({
				// 	createKeystoreRequest: {
				// 		password: this.password
				// 	}
				// }).then((keystore) => {
				// 	this.$emit('created', keystore)
				// }).catch((err) => {
				// 	console.error(err)
				// })
			}
		}
	}
})
</script>

<style lang="scss" scoped>
$bg: #0a0e11;
$accent: #00edb1;
$text-color: #fff;

h4 {
	font-weight: 700;
	margin: 0;
	font-size: 2.4rem;
	text-align: center;
}

h6 {
	margin: 0;
	font-size: 1.1rem;
	font-weight: 300;
}

.input {
	display: flex;
	flex-direction: row;
	flex-wrap: wrap;

	label {
		flex-basis: 100%;
		margin-bottom: 10px;
		opacity: .8;
		padding: 4px 8px;
		display: none;
	}

	input {
		margin: 0;
		border: 0;
		color: #fff;
		outline: none;
		width: 100%;
		height: 56px;
		display: block;
		font-size: 1.1rem;
		font-weight: 400;
		box-sizing: border-box;
		background-color: rgba(#abc, .05);
		border-radius: 5px;
		padding: 0 16px;
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

		&.selected {
			background: transparent;
			font-weight: 600;
			color: $accent;
			padding: 0;
			height: auto;
			font-size: 1.4rem;
		}
	}

	span {
		color: #ff9999;
		font-weight: 500;
		font-size: .9rem;
		padding: 10px 16px;
		opacity: 0;
	}

	&.invalid {
		span {
			opacity: 1;
		}
	}
}

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

//TODO focus

//box-shadow: 0 0 0 0 rgba($accent, .5);
//
//&:focus {
//	transition: .1s ease;
//	box-shadow: 0 0 0 2px rgba($accent, .5);
//}

.login-form {
	width: 100%;
	height: 100vh;
	display: flex;
	flex-direction: row;
	align-items: center;
	justify-content: center;
	background: $bg;

	.login-form-content {
		position: fixed;
		color: #e1e4e8;
		width: 500px;

		.login-form-footer {
			display: flex;
			flex-direction: column;
			align-items: center;
			margin-top: 64px;
		}
	}
}

</style>
