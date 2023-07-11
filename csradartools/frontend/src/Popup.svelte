<script>
    export let windowX = 50 
    export let windowY = 50
</script>
    <div class="background"></div>
    <div class="window" style="width:{windowX}px; height:{windowY}px;">
        <slot />
    </div>
<style>

    .window {
        flex:1;
        position:absolute;
        background-color:black;
        border: 1px green solid;
        z-index:6;
    }

    .background {
        flex:1;
        position: absolute;
		top :0;
		left: 0;
		right: 0;
		bottom: 0;
		background: rgba(0,0,0,0.5);
        z-index:5;
    }

</style>

