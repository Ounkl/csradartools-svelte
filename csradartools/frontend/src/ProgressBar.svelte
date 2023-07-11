<script>
     import {currentTick} from './stores.js'
     import {GetTickCount} from "../wailsjs/go/main/App";
     import { onMount } from "svelte";

    export let totaltime = 1;
    export let radar;

    let progress = 0;
    let preview = 0;
    let mouseHover = false;
    let currenttime = 0;
    let barWidth;

    currentTick.subscribe((value) => {
        currenttime = value;
        progress = (currenttime / totaltime) * 100
    })
    

    function updateTime(event) {
        console.log(event.clientX, event.offsetX, this.offsetLeft)
   
            progress = event.offsetX

            currentTick.update((n) => n = Math.trunc((progress / barWidth) * totaltime))
        }

    function advanceTime(event) {
        console.log(event.keypress)
        //37, 39

        switch(event.keyCode) {
            //leftarrow
            case 37:
                currentTick.update((n) => n = n - 64)
                break;
            //rightarrow
            case 39:
                currentTick.update((n) => n = n + 64)
                break;
        }
    }

    function updatePreview(event) {
        console.log(event.clientX, event.offsetX, this.offsetLeft)
   
        preview = event.offsetX

    }

</script>

<svelte:window on:keydown|preventDefault={advanceTime} />

<div class="fullbar">
    <button class="btn" on:click={radar.stopDemo}>⏵︎</button>
    <button class="btn">⏩︎</button>
    <div class="timetext">{currenttime} / {totaltime}</div>
    <div class="stylebar"  on:mouseleave={() => mouseHover = false} on:focus={() => mouseHover = true} on:mouseover={() => mouseHover = true} on:keydown|preventDefault={advanceTime} on:click={updateTime} on:mousemove={updatePreview}>
        <div class="outerprogbar" bind:clientWidth={barWidth}>
            <div class="hoverprogbar" style="width:{(preview / barWidth) * 100}%">
                <div class="progbar" style="width:{progress}%"></div>

                {#if mouseHover}
                    <div class="progcircle" style="left:{progress}%"></div>
                {/if}

            </div>
        </div>
    </div>
</div>

    
<style> 

    .timetext {

        font-family: 'Roboto';
        font-size:0.75em;
        padding:2em;
        padding-right:1rem;
        padding-left:1rem;
        white-space:nowrap;
    }

    .fullbar {
        display:flex;
        width:100%;
        margin:auto;
        margin-top:1%;
        overflow-x:hidden;
    }

    .btn {
        padding: 1% 2%;
        background:none;
        border:none;
        cursor:pointer;
        color:white;
        background-color:#464447;
        margin-left: 1%;
        margin-right:1%;
        font-size:2em;
    }

    .stylebar {
        width:100%;
        height:30px;
        margin:auto;
        line-height: 30px;
    }

    .hoverprogbar {
        width:100%;
        height:100%;
        background-color:#585559;
        border-style:none;
        line-height: 30px;
    }

    .outerprogbar {
        width:100%;
        height:10%;
        background-color:#464447;
        margin:auto;
        line-height: 30px;
        top:50%;
        position:relative;
    }

    .progbar {
        border-style:none;
        width:100%;
        height:100%;
        background-color:#b86168;
        line-height: 30px;
        position:absolute;
        display:flex;
    }

    .progcircle {
        height:12px;
        width:12px;
        background-color: #b86168;
        border-radius: 50%;
        top:-3px;
        position:absolute;
        display:flex;
        animation-name:animcircle;
        animation-duration:250ms;
    }

    .stylebar:hover {
        animation-name:animbar;
        animation-duration:250ms;
        height:50px;
        cursor:pointer;
    }

    .stylebar:not(:hover) {
        animation-name:animbarout;
        animation-duration:250ms;
        height:30px;
    }

    @keyframes animcircle {
        from {height:0px;width:0px;top:0;}
        to {height:12px;width:12px;top:-3px;}
    }

    @keyframes animbar {
        from {height:30px;}
        to {height:50px;}
    }

    @keyframes animbarout {
        from {height:50px;}
        to {height:30px;}
    }


</style>

