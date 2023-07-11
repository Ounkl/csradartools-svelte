
<script>

    export let name;
    export let position;
    export let viewangle;
    export let team;
    export let alive;

    let canvasHeight = 1024;
    let canvasWidth = 1024;

    const iconheight = 10;
    const iconwidth = 10

    let playerColor = 'gray';

    if(team == 2) {
        playerColor = 'red';
    }

    if(team == 3) {
        playerColor = 'blue';
    }

    console.log(name, team, playerColor)
    //const iconheight;
    //const iconwidth;

    //can add more variables, such as kills, deaths, etc.

    //{#each players as player}
	    //<div class="picon" style="top:{player[1] - (iconheight / 2)}px; left:{player[0] - (iconwidth / 2)}px; height:{iconheight}px; width:{iconwidth}px"/>
    //{/each}
</script>

    {#if alive}
        <playername class="pname" style="left:{(position.x / 1024.0) * canvasWidth + 10}px; top:{(position.y / 1024.0) * canvasHeight - 10}px">{name}</playername>
        <playericon class="picon" style="background-color:{playerColor}; left:{(position.x / 1024.0) * canvasWidth - (iconwidth / 2)}px; top:{(position.y / 1024.0) * canvasHeight - (iconheight / 2)}px"/>
    {/if}
    
<style>
.picon {
    flex:1;
    position:absolute;
    height:10px;
    width:10px;
    border-radius: 50%;
    z-index:4;
}

.pname {
    flex:1;
    position:absolute;
    z-index:4;
    font-size:0.5rem;
}
</style>

