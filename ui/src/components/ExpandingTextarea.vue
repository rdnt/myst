<template>
	<div :class="{focus: hasFocus, disabled}" class="field" tabindex="1" @focus="focus">
		<label>{{ label }}</label>
		<textarea
			:spellcheck="false"
			ref="textarea"
			:value="modelValue"
			:style="{
				height: height + 'px',
			}"
			:disabled="disabled"
			@blur="hasFocus = false"
			@focus="hasFocus = true"
			@input.passive="(event: InputEvent) => {
				$emit('update:modelValue', event.target.value);
				update();
			}"
		/>
	</div>
</template>

<script lang="ts">
import {defineComponent} from 'vue'

export default defineComponent({
	name: 'ExpandingTextarea',
	props: {
		placeholder: {
			type: String
		},
		disabled: {
			type: Boolean,
		},
		label: {
			type: String
		},
		modelValue: {
			type: String,
			default: ''
		}
	},
	data: (): {
		height: number,
		hasFocus: boolean
	} => ({
		height: 0,
		hasFocus: false,
	}),
	watch: {
		modelValue() {
			this.update()
		},
	},
	mounted() {
		console.log('mounted')

		this.$nextTick(() => {
			if (this.modelValue == "") {
				this.$emit('update:modelValue', this.$props.placeholder || '')
			}

			// this.value = this.$props.placeholder || ''

			this.update().then(() => {
				// this.currentValue = ''

				if (this.modelValue == "") {
					this.$emit('update:modelValue', this.modelValue)
				}

			})

			window.addEventListener('resize', this.update)
		})
	},
	unmounted() {
		console.log('unmounted')

		window.removeEventListener('resize', this.update);
	},
	methods: {
		update() {
			this.height = 0

			return this.$nextTick(() => {
				const textarea = this.$refs.textarea as HTMLTextAreaElement
				this.height = textarea.scrollHeight
			})
		},
		focus() {
			const textarea = this.$refs.textarea as HTMLTextAreaElement
			textarea.focus()
		}
	}
})
</script>

<style lang="scss" scoped>
.field {
	position: relative;
	overflow: hidden;
	display: flex;
	flex-direction: column;
	//margin-bottom: 30px;
	border-radius: 5px;
	padding-bottom: 16px;
	flex-shrink: 0;
	//background-color: rgba(#abc, .05);

	&.disabled {
		background-color: transparent;

		textarea {
			pointer-events: none;
		}

		&.focus {
			background-color: transparent;

			label {
				background-color: #111519;
				color: rgba(#8a8f9f, 1);
			}
		}
	}

	&.focus {
		background-color: rgba(#abc, .1);

		label {
			background-color: #20252b;
			color: rgba(#8a8f9f, 1);
		}
	}

	label {
		width: 100%;
		box-sizing: border-box;
		display: block;
		padding: 16px 16px 6px;
		background-color: #111519;

		color: rgba(#8a8f9f, .75);
		text-transform: uppercase;
		font-size: 0.85rem;
		font-weight: 600;
		letter-spacing: 0.5px;
		pointer-events: none;
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
		background-color: transparent;
		//background-color: rgba(#abc, .05);
		padding: 0 16px;
		color: #fff;
		//overflow: hidden;
		max-height: 200px;
		flex-shrink: 0;

		//&::placeholder {
		//	color: darken(#68737e, 10%);
		//}
		//
		//&:focus {
		//	&::placeholder {
		//		color: lighten(#68737e, 0%);
		//	}
		//}

		&:disabled {
			//background-color: transparent;
			//padding-left: 0;
			//padding-right: 0;

		}
	}
}
</style>
