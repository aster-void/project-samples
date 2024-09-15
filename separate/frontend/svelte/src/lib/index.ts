// place files you want to import through the `$lib` alias in this folder.
export const sleep = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));

import { readable } from "svelte/store";
export let increment = readable(0, (_, update) => {
  setInterval(() => {
    update(val => val + 1);
  }, 1000);
});

