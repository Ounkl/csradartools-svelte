<script>
  import logo from './assets/images/radar.png';
  import {Greet} from '../wailsjs/go/main/App.js';
  import {GetDemo} from "../wailsjs/go/main/App";
  import {GetPlayerPositions} from "../wailsjs/go/main/App";
  import {GetTickCount} from "../wailsjs/go/main/App";
  import {WriteBoundary} from "../wailsjs/go/main/App";
  import {GetRayCast} from "../wailsjs/go/main/App";
  import {GetViewDirections} from "../wailsjs/go/main/App";
  import { tick } from 'svelte';
  import { onMount } from 'svelte';

  var iconheight = 10;
  var iconwidth = 10;

  let resultText = "Please enter your name below 👇"
  let name
  let parseReady = false
  let Logo = logo
  let resetImg = Date.now();
  let render = 0

  let playdemo = true;

  async function stopDemo() {
    playdemo = !playdemo;
  }

  function greet() {
    Greet(name).then(result => resultText = result)
  }

  function getdemo() {
    GetDemo().then(updateImage).then(() => parseReady = true);
  }

  const updateImage = () => { 
        console.log('./assets/images/radar.png?' + resetImg);

        import('./assets/images/radar.png?' + resetImg).then(res => Logo = res.default)

        render += 1;

        resetImg = Date.now();
  }

  const delay = ms => new Promise(res => setTimeout(res, ms));

  async function renderdemo() {

    saveCanvas();
    var tickcount;
        
        GetTickCount().then((result) => tickcount = result).then(async () => {
            for(let i = 0; i < tickcount;) {
              if(playdemo) {
                await delay(2) 
                rendertick(i)
                await tick();
                i++;
              }
              await delay(2)
            }
        });
  }

  let players = []

  async function rendertick(tick) {

        GetPlayerPositions(tick).then((result) => {

            players = []

            //console.log(tick)
            //console.log(result)
            for(let j = 0; j < result.length; j = j + 2) {

                let coords = [result[j], result[j + 1]]

                players.push(coords)

                players = players
            }
        });

        GetRayCast(tick).then((result) => {

          //console.log(tick)
          console.log(result)

          for(let j = 0; j < result.length; j = j + 4) {
            drawLine(result[j], result[j + 1], result[j + 2], result[j + 3])
          }
        });

        loadCanvas();
    }

  let mouseX = 0;
  let mouseY = 0;
  let wallmode = false;

  let lineX = 0;
  let lineY = 0;

  let lineCreate = false;

  let canvas;
  let ctx;

  let boundary = [];

  onMount(async () => {
    canvas = document.getElementById("myCanvas");

    // @ts-ignore
    if(!canvas.getContext) {
      console.log("bad context");
      return;
    }
    // @ts-ignore
    ctx = canvas.getContext("2d");

    ctx.lineWidth = 5;
    ctx.strokeStyle = 'red';

    console.log("Good context");
    
  })

  function wallMode() {
    wallmode = true;

    console.log("wallMode on");
  }

  function handleMousemove(event) {
    mouseX = event.clientX;
    mouseY = event.clientY;

    if(lineCreate == true) {

        loadCanvas();

        ctx.lineWidth = 5;
        ctx.strokeStyle = 'red';

        drawLine(lineX, lineY, event.clientX, event.clientY)

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

  function createLine(event) {

    if(wallmode == false) {
      return;
    }

    if(lineCreate == true) {
        ctx.lineWidth = 5;
        ctx.strokeStyle = 'red';

        drawLine(lineX, lineY, event.clientX, event.clientY)

        lineCreate = false;

        boundary.push(lineX, lineY, event.clientX, event.clientY);

        //console.log(boundary);

        return;
    }

    saveCanvas();

    lineX = event.clientX;
    lineY = event.clientY;

    lineCreate = true;
  }

  let imageData;

  function saveCanvas() {
    imageData = ctx.getImageData(0,0,canvas.width,canvas.height);
  }

  function loadCanvas() {
    ctx.putImageData(imageData, 0, 0);
  }

  function saveBoundary() {
    WriteBoundary(boundary);
  }

</script>

<main>
  <div class="replay" id="players">

    <img alt="Wails logo" class="radar" id="radar" src="{Logo}">

    <div class="picon" style="top:-5px; left:-5px"/>
    {#each players as player}
	    <div class="picon" style="top:{player[1] - (iconheight / 2)}px; left:{player[0] - (iconwidth / 2)}px; height:{iconheight}px; width:{iconwidth}px"/>
    {/each}
    <canvas on:mousemove={handleMousemove} on:mousedown={createLine} id="myCanvas" width="1024px" height="1024px" style="flex:1; border:1px red solid; position:absolute;"></canvas>
  </div>


  <div class="result" id="result">{resultText}</div>
  <div class="input-box" id="input">
    <button class="btn" on:click={getdemo}>Get Demo</button>
    <button class="btn" on:click={renderdemo}>Render</button>
    <button class="btn" on:click={wallMode}>Define Walls</button>
    <button class="btn" on:click={saveCanvas}>Save</button>
    <button class="btn" on:click={loadCanvas}>Load</button>
    <button class="btn" on:click={saveBoundary}>Save Boundaries</button>
    <button class="btn" on:click={stopDemo}>Resume/Pause</button>
    Mouse is {mouseX} x {mouseY}
  </div>


</main>

<style>

  .result {
    height: 20px;
    line-height: 20px;
    margin: 1.5rem auto;
  }

  .input-box .btn {
    width: 60px;
    height: 30px;
    line-height: 30px;
    border-radius: 3px;
    border: none;
    margin: 0 0 0 20px;
    padding: 0 8px;
    cursor: pointer;
  }

  .input-box .btn:hover {
    background-image: linear-gradient(to top, #cfd9df 0%, #e2ebf0 100%);
    color: #333333;
  }

  .input-box .input {
    border: none;
    border-radius: 3px;
    outline: none;
    height: 30px;
    line-height: 30px;
    padding: 0 10px;
    background-color: rgba(240, 240, 240, 1);
    -webkit-font-smoothing: antialiased;
  }

  .input-box .input:hover {
    border: none;
    background-color: rgba(255, 255, 255, 1);
  }

  .input-box .input:focus {
    border: none;
    background-color: rgba(255, 255, 255, 1);
  }

</style>


