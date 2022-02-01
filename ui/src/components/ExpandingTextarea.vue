<template>
	<div class="field" :class="{focus: focus}">

		<label>Some label</label>
		<textarea
			ref="textarea"
			:placeholder="placeholder"
			:value="modelValue"
			:style="{
		height: height + 'px',
	}"
			@input="
	 (event: InputEvent) => {
		 $emit('update:modelValue', event.target.value)
		update()
	 }
	"
			@focus="focus = true"
			@blur="focus = false"
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
		modelValue: {
			type: String,
			default: ''
		}
	},
	data: (): {
		height: number,
		focus: boolean
	} => ({
		height: 0,
		focus: false,
	}),
	mounted() {
		this.$nextTick(() => {
			if (this.value == "") {
				$emit('update:modelValue', this.$props.placeholder || '')
			}

			// this.value = this.$props.placeholder || ''

			this.update().then(() => {
				// this.currentValue = ''

				if (this.value == "") {
					$emit('update:modelValue', this.modelValue)
				}

			})

			window.addEventListener('resize', this.update)
		})
	},
	unmounted() {
		window.removeEventListener('resize', this.update);
	},
	methods: {
		update() {
			this.height = 0

			return this.$nextTick(() => {
				const textarea = this.$refs.textarea as HTMLTextAreaElement
				this.height = textarea.scrollHeight
			})
		}
	}
})
</script>

<style lang="scss" scoped>
.field {
	display: flex;
	flex-direction: column;
	margin-bottom: 30px;

	&.focus {
		background-color: rgba(#abc, .05);
	}

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
