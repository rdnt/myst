<template>
	<div :class="{focus: hasFocus, disabled}" class="field" tabindex="1" @focus="focus">
		<label>{{ label }}</label>
		<textarea
			:spellcheck="false"
			ref="textarea"
			:value="modelValue || value"
			:style="{
				height: height + 'px',
			}"
			:placeholder="placeholder"
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
		},
		value: {
			type: String,
			default: ''
		},
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
			// console.log(this.modelValue)
			this.update()
		},
		value() {
			// console.log(this.value)
			this.update()
		}
	},
	mounted() {
		// console.log('mounted')

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
		// console.log('unmounted')

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
	background-color: #191e23;
	margin: 12px 0;
	//margin-bottom: 2px;
	transition: .18s ease;
	cursor: text;

	&:last-child {
		margin-bottom: 0;
	}

	&.disabled {
		margin: 0;
		background-color: transparent;
		cursor: default;

		label {
			//background-color: transparent;
		}
		textarea {
			pointer-events: none;

			&::placeholder {
				color: #fff;
			}
		}


		&.focus {
			background-color: transparent;

			label {
				//background-color: transparent !important;

				color: rgba(#8a8f9f, 1);
			}
		}
	}

	&.focus {
		background-color: rgba(#abc, .1);

		label {
			//background-color: #20252b;
			color: rgba(#8a8f9f, 1);
		}
	}

	label {
		width: 100%;
		box-sizing: border-box;
		display: block;
		padding: 16px 16px 6px;
		color: rgba(#8a8f9f, .75);
		text-transform: uppercase;
		font-size: 0.85rem;
		font-weight: 600;
		letter-spacing: 0.5px;
		pointer-events: none;
	}

	textarea {
		position: relative;
		display: flex;
		flex-direction: column;
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

		padding-top: 40px;
		padding-bottom: 20px !important;
		top: -40px;

		//background-color: rgba(#0c1d19, .5);
		margin-bottom: -60px;

		&::placeholder {
			color: #334;
		}
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

.field-button {
	display: block;
	outline: none;
	border: none;
	height: 40px;
	width: 40px;
	font-size: 1.1rem;
	font-weight: 500;
	border-radius: 5px;
	background-color: rgba(#202228, 1);
	color: #fff;
	margin-left: 10px;
}
</style>
