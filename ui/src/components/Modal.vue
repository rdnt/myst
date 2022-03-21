<template>
	<div class="modal" :class="{show}">
		<div class="overlay" @click="setShow(false)"></div>
		<transition duration="180" name="show">
			<div class="body" v-if="show">
				<div class="header">
					<slot name="header"></slot>
				</div>
				<div class="content">
					<slot></slot>
				</div>
				<div class="footer">
					<slot name="footer"></slot>
				</div>
			</div>
		</transition>
	</div>
</template>

<script>
import {defineComponent} from 'vue'

export default defineComponent({
	name: 'modal',
	props: [
		'link',
		'show',
		'setShow'
	]
});
</script>

<style lang="scss" scoped>

.show-enter-active, .show-leave-active {
	transition: .18s;
}

.show-enter-from {
	opacity: 0;
	transform: scale(.95);
}

.show-enter-to {
	opacity: 1;
	transform: scale(1);
}

.show-leave-from {
	opacity: 1;
	transform: scale(1);
}

.show-leave-to {
	opacity: 0;
	transform: scale(.95);
}

//.show-enter-active,
//.show-leave-active {
//	opacity: 1;
//	.content {
//		transform: scale(1);
//	}
//}
//.show-enter,
//.show-leave-to {
//	opacity: 0;
//	.content {
//		transform: scale(0.95);
//	}
//}
.modal {
	//position: fixed;
	//top: 0;
	//left: 0;
	//width: 100vw;
	//height: 100vh;
	//display: flex;
	//align-items: center;
	//justify-content: center;
	//background: transparent;
	pointer-events: none;
	//transition: opacity 0.1s ease;

	position: fixed;
	top: 0;
	left: 0;
	width: 100%;
	height: 100%;
	height: -webkit-fill-available;
	padding: 50px 0;
	display: flex;
	flex-direction: row;
	align-items: center;
	justify-content: center;

	&.show {
		pointer-events: all;

		.overlay {
			opacity: .85;
		}
	}

	//display: none;
	//display: flex;
	//align-items: center;

	.body {
		position: relative;
		max-height: 100%;
		max-width: 500px;
		flex-basis: 100%;
		background: #111519;
		display: flex;
		flex-direction: column;
		overflow: hidden;
		margin: 0 auto;
		padding: 20px;
		border-radius: 5px;
		box-sizing: border-box;
		box-shadow: 0 8px 24px 4px rgba(0, 0, 0, .5);
		//margin-bottom: 10%;
		will-change: transform;

		.header {
			//margin-bottom: 16px;
		}

		.content {
			overflow-y: auto;
		}

		.footer {
			//margin-top: 16px;
		}
	}


	> .content {
		//background: #292e39;
		//box-shadow: 0px 0px 47px (#000, 0.5);
		//position: absolute;
		//display: flex;
		//max-width: 100%;
		//max-height: 100%;
		//backface-visibility: hidden;
		//pointer-events: all;
		//transition: transform 0.1s ease;
		//border-radius: $radius;
		//overflow: hidden;


	}

	.overlay {
		position: fixed;
		top: 0;
		left: 0;
		width: 100%;
		height: 100%;
		height: -webkit-fill-available;
		background: #000;
		//opacity: 0.86;
		opacity: 0;
		transition: opacity .18s ease;
	}
}
</style>
