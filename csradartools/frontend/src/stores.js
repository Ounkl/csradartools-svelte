import { writable } from 'svelte/store';

export const count = writable(0);
export const currentTick = writable(0);
export const radarWidth = writable(0);
export const radarHeight = writable(0);
