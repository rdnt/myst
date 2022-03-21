<template>
	<div class="field">
		<label class="label">{{ label }}</label>
		<expanding-textarea :modelValue="modelValue"
												class="textarea"
												@update:modelValue="onUpdate($event)"
		></expanding-textarea>
	</div>
</template>


<script lang="ts">
import {defineComponent} from 'vue'
import ExpandingTextarea from "../components/ExpandingTextarea2.vue";

export default defineComponent({
	name: 'entry-text-field',
	components: {ExpandingTextarea},
	data: (): {} => ({}),
	props: {
		label: {
			type: String,
			required: true
		},
		modelValue: {
			type: String,
		},
	},
	watch: {
		modelValue() {
			// console.log('watch from text field');
		}
	},
	methods: {
		onUpdate(value: string) {
			// console.log('onUpdate from text field', event);
			this.$emit('update:modelValue', value);

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
	margin: 12px 0;
	//margin-bottom: 2px;
	cursor: text;
	z-index: 0;

	.label {
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
		z-index: 2;
	}

	.textarea {
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
		box-sizing: content-box;
		background-color: transparent;
		//background-color: rgba(#abc, .05);
		padding: 0 16px;
		color: #fff;
		//overflow: hidden;
		//max-height: 80px;
		flex-shrink: 0;

		padding-top: 40px;
		padding-bottom: 20px !important;
		top: -40px;
		background-color: #191e23;
		transition: .18s ease;

		//background-color: rgba(#0c1d19, .5);
		margin-bottom: -60px;
		z-index: 1;

		&:focus {
			background-color: rgba(#abc, .1);
			//background-color: transparent;
		}
	}
}
</style>
