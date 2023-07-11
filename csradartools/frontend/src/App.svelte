<script>
  import {GetDemo} from "../wailsjs/go/main/App";
  import {GetTicksOfInterest} from "../wailsjs/go/main/App";
  import {GetTickCount} from "../wailsjs/go/main/App";
  import Radar from './Radar.svelte'
  import ProgressBar from './ProgressBar.svelte'
  import {currentTick} from './stores.js'
  import {radarHeight} from './stores.js'
  import {radarWidth} from './stores.js'
  
  let parseReady = false
  let tickcount;
  let ticksOfInterest = []
  let currenttickofinterest;
  let radarheight;
  let radarwidth;
  

  let currenttick

  radarHeight.subscribe((value) => {
    radarheight = value
  })

  radarWidth.subscribe((value) => {
    radarheight = value
  })

  currentTick.subscribe((value) => {
        currenttick = value;
  })

  function getdemo() {
    GetDemo().then(() => {
      GetTickCount().then((result) => {
        tickcount = result;
        parseReady = true;
      });
    })
  }

  function getTicksOfInterest() {
    GetTicksOfInterest().then(result => {
      console.log(result)
      let currTick = result[0]
      if(result.length != 0) {
        for(let i = 1; i < result.length; i++) {
          if(result[i] - currTick > 256) {
            ticksOfInterest.push(currTick)
            ticksOfInterest.push(result[i])
            currTick = result[i]
          }
        }
        //informs the svelte DOM compiler that ticksOfInterest has changed
        ticksOfInterest = ticksOfInterest
      }
    })
  }

  function nextTickOfInterest() {
    for(let i = 0; i < ticksOfInterest.length; i++) {
      if(ticksOfInterest[i] > currenttick) {
        currentTick.update((n) => n = ticksOfInterest[i])
        currenttickofinterest = i;
        break;
      }
    }
  }

  function lastTickOfInterest() {
    for(let i = 1; i < ticksOfInterest.length; i++) {
      if(ticksOfInterest[i] > currenttick) {
        currentTick.update((n) => n = ticksOfInterest[i - 1])
        currenttickofinterest = i - 1;
        break;
      }
    }
  }

  let radar;

</script>

<main>

{#if parseReady}

<div class="fullui">
  <div class="playback">
    <div class="tickbar">
      {#if ticksOfInterest.length == 0}
        <button class="btn" on:click={getTicksOfInterest}>Parse</button>
      {/if}
      {#if ticksOfInterest.length > 0}
      <button class="btn" on:click={lastTickOfInterest}>Last</button>
        <button class="btn" on:click={nextTickOfInterest}>Next</button>
      {/if}
    </div>
    <Radar tickcount={tickcount} bind:this={radar}/>
    <div class="demoinfo">
      <div id="roundcount">Round 15</div>
      <div id="roundtime">1:52</div>
    </div>
    <ProgressBar radar={radar} totaltime={tickcount}/>
  </div>
</div>
{/if}

{#if !parseReady}
  <button class="btn" on:click={getdemo}>Get Demo</button>
{/if}


</main>

<style>

  .tickbar {
    display:flex;
    width:100%;
    margin-top:2%;
    margin-bottom:2%;
    max-width:100%;
  }

  .fullui {
    display:flex;
  }

  .playback {
    margin:auto;
  }

  .demoinfo {
    display:grid;
  }

  .demoinfo #roundtime {
    font-size:3em;
  }

  .demoinfo #roundcount {
    margin-top:1%;
    font-size:1em;
    color:goldenrod;
  }

.btn {
    padding: 20px 30px;
    font-family:'Century Gothic';
    font-size:15px;
    border-style:none;
    border-radius:30px;
    background-color:#b86168;
    line-height: 30px;
    margin:auto;
  }

  .btn:hover {
    filter: brightness(1.1);
    animation-name:bigbtn;
    animation-duration:1s;
    font-size:20px;
  }
  .btn:not(:hover) {
    filter: brightness(1.0);
    animation-name:bigbtnout;
    animation-duration:1s;
    font-size:15px;
  }

  @keyframes bigbtn {
    from {font-size:15px}
    to {font-size:20px}
  }

  @keyframes bigbtnout {
    from {font-size:20px}
    to {font-size:15px}
  }


</style>


