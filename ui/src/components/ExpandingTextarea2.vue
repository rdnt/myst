<template>
	<textarea
		ref="textarea"
		:style="{
				height: height,
			}"
		:value="value"
		@input="onInput($event.target.value)"
	/>
</template>

<script lang="ts">
import {defineComponent} from 'vue'

export default defineComponent({
	name: 'ExpandingTextarea',
	props: {
		disabled: {
			type: Boolean,
		},
		modelValue: {
			type: String,
			default: ''
		},
		// paddingTop: {
		// 	type: Number,
		// 	default: 0,
		// },
		// paddingBottom: {
		// 	type: Number,
		// 	default: 0,
		// },
	},
	data: (): {
		height: string,
		value: string,
	} => ({
		height: 'auto',
		value: '',
	}),
	watch: {
		modelValue(value: string) {
			console.log('watch from expanding textarea 2');
			this.value = value;
			this.$nextTick(this.update)
			// this.$emit('update:modelValue', value);
		}
	},
	mounted() {
		this.value = this.modelValue
		this.$nextTick(this.update)

		// window.addEventListener('resize', this.update)
	},
	unmounted() {
		// window.removeEventListener('resize', this.update);
	},
	methods: {
		update() {
			// console.log('update')

			// const style = this.$el.currentStyle || window.getComputedStyle(this.$el);
			//
			// const paddingTop = parseInt(style.paddingTop);
			// const paddingBottom = parseInt(style.paddingBottom);

			// console.log(paddingTop, paddingBottom)

			this.$el.style.height = '';
			this.$el.style.paddingTop = 0;
			this.$el.style.paddingBottom = 0;
			console.log(this.$el.scrollHeight - paddingTop - paddingBottom)
			this.$el.style.height = (this.$el.scrollHeight - paddingTop - paddingBottom) + "px";
			this.$el.style.paddingTop = paddingTop;
			this.$el.style.paddingBottom = paddingBottom;

			//
			// this.height = '';
			//
			// let scrollHeight = 0;
			// this.$nextTick(() => {
			// 	console.log(this.$el)
			// 	scrollHeight = this.$el.scrollHeight;
			// 	console.log(scrollHeight);
			// }).then(() => {
			// 	this.height = scrollHeight - 60 + 'px';
			// });


		},
		onInput(value: string) {
			this.value = value;
			this.$emit('update:modelValue', value);
		}
	}
})
</script>

<style lang="scss" scoped>
textarea {
	display: block;
	width: 100%;
	margin: 0;
	border: none;
	outline: none;
	resize: none;
}
</style>
