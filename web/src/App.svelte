<script lang="ts">
  import { Router, Route } from "svelte-routing";
  import { onMount } from "svelte";

  import Admin from "./admin.svelte";
  import Streamer from "./streamer.svelte";

  export let url = "";

  const logo = "/img/PrismPlus_PlatformColorsTrans.png";
  const logo_small = "/img/PrismPlusLogo_LinePrismWhite.png";
</script>

<div id="main">
  <header>
    <h1>
      <img
        src={logo}
        alt="Prism plus logo --a 20 sided die with many colors streaming in and out of it."
      />
    </h1>
  </header>

  <div id="content">
    <Router {url}>
      <Route path="/" component={Streamer} />
      <Route path="/admin" component={Admin} />
    </Router>
  </div>

  <footer>
    <span>
      <img
        src={logo_small}
        alt="Small Prism+ logo --a white D20 with a plus sign in the middle."
      />
    </span>
    <span>
      Prism+ is free software. <a
        href="https://github.com/geekgonecrazy/prismplus"
        >Check out the source code.</a
      >
    </span>
  </footer>
</div>

<style>
  :global(h1) {
    margin: 0.5rem;
    padding: 0.25rem;

    font-size: 2rem;
  }
  :global(h2) {
    margin: 0.25rem;
    margin: 0.25rem;

    font-size: 1.5rem;
  }
  :global(h3) {
    margin: 0rem;
    padding: 0.25rem;

    font-size: 1.25rem;
  }

  :global(input) {
    width: 100%;
  }
  :global(input:read-only),
  :global(input:read-only:focus-visible) {
    background-color: rgba(0, 0, 0, 0);
    border: none;
    outline: none;
    color: var(--white);
  }

  :global(form) {
    max-width: 100%;
  }
  :global(fieldset) {
    margin: 1em 0;
  }

  :global(.oneline *) {
    display: inline-block;
  }

  :global(.button-negative) {
    background-color: var(--neg-color);
  }
  :global(.button-negative:hover) {
    background-color: rgba(var(--neg-color-rgb), 0.5);
  }
  :global(.button-negative:not(:disabled):active) {
    background-color: var(--white);
  }

  footer {
    display: flex;
    justify-content: space-between;
    gap: 1rem;

    flex-shrink: 0;

    padding: 1em;

    background-color: rgba(var(--black-rgb), 0.5);
  }
  footer img {
    max-width: 19px;
  }

  #main {
    display: flex;
    flex-direction: column;

    height: 100%;
  }
  #main header {
    margin: 0 auto;

    text-align: center;
  }
  #main header img {
    width: 100%;
    max-width: 320px;
  }

  #content {
    flex: 1 0 auto;

    width: 90%;
    max-width: 1024px;
    margin: 0 auto;
    padding: 1em;
  }
</style>
