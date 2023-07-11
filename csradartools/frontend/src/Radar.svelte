
<script>
    import Popup from './Popup.svelte'
    import logo from './assets/images/radar.png';
    import { GetRayCast } from "../wailsjs/go/main/App";
    import { DefineControlVector } from "../wailsjs/go/main/App";
    import { WriteBoundary } from "../wailsjs/go/main/App";
    import { onMount } from "svelte";
    import {GetPlayers} from "../wailsjs/go/main/App";
    import { tick } from 'svelte';
    import Player from './Player.svelte'
    import {GetTickCount} from "../wailsjs/go/main/App";
    import {CheckControlVectors} from "../wailsjs/go/main/App";
    import {currentTick} from './stores.js'
    import {radarHeight} from './stores.js'
    import {radarWidth} from './stores.js'

    

    export let tickcount;

    let controlState = 0;
    let newLine = {a: {x: 0, y:0}, b: {x:0, y:0}}
    let mouseX = 0;
    let mouseY = 0;
    let canvas;
    let ctx;
    let lineCreate = false;
    let boundary = [];
    let controlFocus = false;
    let currenttick = 0;
    let playdemo = true;
    let radarheight;
    let radarwidth;

  radarHeight.subscribe((value) => {
    radarheight = value
  })

  radarWidth.subscribe((value) => {
    radarwidth = value
  })

    currentTick.subscribe((value) => {
        currenttick = value;
    })

    onMount(async () => {
    canvas = document.getElementById("myCanvas");

    // @ts-ignore
    if(!canvas.getContext) {
      console.log("bad context");
      return;
    }
    // @ts-ignore
    ctx = canvas.getContext("2d", {willReadFrequently: true});

    ctx.lineWidth = 5;
    ctx.strokeStyle = 'gray';

    console.log("Good context");

    renderdemo();

    
  })

  export function stopDemo() {
    playdemo = !playdemo;
  }

  function checkControlVectors() {
    CheckControlVectors(currenttick).then(result => {
      console.log(result)
    })
  }

export async function renderdemo() {
    
    GetTickCount().then((result) => tickcount = result).then(async () => {
        for(; currenttick < tickcount;) {
          if(playdemo) {

            await delay(2) 

            rendertick(currenttick).then(() => {
              currenttick++;
              currentTick.update((n) => n = currenttick)
            });
            await tick();
          }
          await delay(2)
        }
    });
}

const delay = ms => new Promise(res => setTimeout(res, ms));

async function rendertick(tick) {
    getPlayers(tick)
}

export let players = []

function getPlayers(tick) {

  resetCanvas()

  GetPlayers(tick).then(result => {
      if(result != null) {
        players = []
        players = result
        for(let i = 0; i < players.length; i++) {
          players[i].position.x = (players[i].position.x / 1024) * radarwidth 
          players[i].position.y = (players[i].position.y / 1024) * radarheight 
        }
      }
  });

  GetRayCast(tick).then(result => {
    if(result != null) {
      result.forEach((vec) => {
        vec.a.x = (vec.a.x / 1024) * radarwidth
        vec.a.y = (vec.a.y / 1024) * radarheight
        vec.b.x = (vec.b.x / 1024) * radarwidth
        vec.b.y = (vec.b.y / 1024) * radarheight

        drawLine(vec.a.x, vec.a.y, vec.b.x, vec.b.y)
      })
    }
  })
}

    
  function handleMousemove(event) {

    var pos = getXY(canvas, event);
    mouseX = pos.x;
    mouseY = pos.y;

    if(lineCreate == true) {

        loadCanvas();

        ctx.lineWidth = 5;
        ctx.strokeStyle = 'gray';

        drawLine(newLine.a.x, newLine.a.y, pos.x, pos.y)

    }

    //console.log(lineCreate, wallmode);
    //console.log(lineX, lineY);
  }
  


  function drawLine(x1, y1, x2, y2) {
        ctx.beginPath()
        ctx.moveTo(x1,y1);
        ctx.lineTo(x2, y2);
        ctx.stroke();
        ctx.closePath();
  }

  var _showWindow;

  
  async function setControl(arg) {
    controlState = arg;
    controlFocus = false;
    newLine.a.x = (newLine.a.x / radarwidth) * 1024
    newLine.a.y = (newLine.a.y / radarheight) * 1024
    newLine.b.x = (newLine.b.x / radarwidth) * 1024
    newLine.b.y = (newLine.b.y / radarheight) * 1024
    // @ts-ignore
    DefineControlVector(newLine, Number(arg));

  }


  async function createLine(event) {
    
    var pos = getXY(canvas, event);

    if(lineCreate == true) {
        ctx.lineWidth = 5;
        ctx.strokeStyle = "gray";

        newLine.b.x = pos.x
        newLine.b.y = pos.y

        drawLine(newLine.a.x, newLine.a.y, pos.x, pos.y)

        lineCreate = false;

        console.log("lineCreate")

      //fix vector and point stuct imports

       controlFocus = true

       //need to get button press before definecontrolvector is called

        //console.log(boundary);

        return;
    }

    saveCanvas();

    newLine.a.x = pos.x;
    newLine.a.y = pos.y;

    lineCreate = true;
  }

  let imageData;

  export function saveCanvas() {
    imageData = ctx.getImageData(0,0,canvas.width,canvas.height);
  }

  export function loadCanvas() {
    ctx.putImageData(imageData, 0, 0);
  }

  export function resetCanvas() {
    ctx.clearRect(0, 0, canvas.width, canvas.height);
  }

  export function saveBoundary() {
    WriteBoundary(boundary);
  }

  function getXY(canvas, event) {
    console.log(event)
    var rect = canvas.getBoundingClientRect();  // absolute position of canvas
    return {
        x: event.clientX - rect.left,
        y: event.clientY - rect.top
    }
}
    
</script>
{mouseX} x {mouseY}
  <div class="replay" id="players" bind:clientWidth={radarwidth} bind:clientHeight={radarheight}>
    <img alt="Wails logo" class="radar" id="radar" src="{logo}">
    <canvas class="canvas" on:mousemove={handleMousemove} on:mousedown={createLine} width={radarwidth} height={radarheight} id="myCanvas"></canvas>
    {#each players as player (player.position)} 
      <Player {...player}/>
    {/each}
    {#if controlFocus}
    <Popup> 
      <button on:click={() => setControl(0)}>No Team</button>
      <button on:click={() => setControl(1)}>All Teams</button>
      <button on:click={() => setControl(2)}>Terrorists</button>
      <button on:click={() => setControl(3)}>Counter-Terrorists</button>
    </Popup>
    {/if}
  </div>
    
<style>

  .radar {
    position:absolute;
    width:100%;
    z-index: 1;
    margin:auto;
    border-radius:10%;
    box-shadow: 0 1px 50px 1px rgb(121, 92, 97);
}

.canvas {
    position:absolute;
    z-index: 3;
}

.replay {
    display:flex;
    height:60vh;
    width:60vh;
    position:relative;
    padding:0;
    z-index:3;
    margin:auto;
}
</style>

