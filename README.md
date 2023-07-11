# csradartools-svelte

## About
This tool is being developed personally to make watching demos easier.
NOT meant to replace traditional ingame demo studying, just to find specific things in a demo MUCH faster.

## Features
The tool contains experimental ray-casting to simulate what a player in-game would see.
Every tick can be parsed through to find if a player has view of a specific position.
Experimental because the boundaries are a 2d simulation and will not give 100% accurate feedback to what a player would really see.
Still good enough to filter out 90% of the fluff of a demo, though.

## TODO
See TODO.md.

## Live Development

To run in live development mode, run `wails dev` in the project directory. This will run a Vite development
server that will provide very fast hot reload of your frontend changes. If you want to develop in a browser
and have access to your Go methods, there is also a dev server that runs on http://localhost:34115. Connect
to this in your browser, and you can call your Go code from devtools.

## Building

To build a redistributable, production mode package, use `wails build`.





