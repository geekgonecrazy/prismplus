<script lang="ts">
  import { Router, Route } from "svelte-routing";
  import { onMount } from "svelte";

  import Admin from "./admin.svelte";
  import Streamer from "./streamer.svelte";

  export let url = "";
</script>

<div id="main">
  <div id="content">
    <Router {url}>
      <Route path="/" component={Streamer} />
      <Route path="/admin" component={Admin} />
    </Router>
  </div>

  <footer>
    Prism+ is free software. <a
      href="https://github.com/geekgonecrazy/prismplus"
      >Check out the source code.</a
    >
  </footer>
</div>

<style>
  h1 {
    color: var(--prim-color);

    text-transform: uppercase;
    font-size: 4em;
    font-weight: 100;
  }
  h2 {
    color: var(--prim-color);

    font-size: 1.5em;
  }
  h3 {
    color: var(--prim-color);

    font-size: 1.25em;
  }

  button {
    margin: 1em 0;
  }
  label {
    margin: 0.5em;
  }
  input {
    width: 100%;
  }
  input:read-only {
    background-color: rgba(0, 0, 0, 0);
    border: none;
    color: rgb(var(--fg-light));
  }
  fieldset {
    margin: 1em 0;
  }
  footer {
    flex-shrink: 0;
    padding: 1em;

    text-align: right;

    background-color: rgba(255 255 255 200);
  }
  footer a {
    color: #afa;
  }

  .button-negative {
    background-color: #ff7777;
  }

  #main {
    height: 100%;
    display: flex;
    flex-direction: column;
  }

  #content {
    flex: 1 0 auto;

    max-width: 320px;
    margin: 0 auto;
    padding: 1em;

    text-align: center;
  }

  #current-sessions ul {
    margin: auto;
    width: calc(100% / 2);
    text-align: left;
  }
  #current-sessions ul {
    list-style-type: none;
  }
  #current-sessions li > span,
  #current-sessions li > button {
    display: inline-block;
    margin: 0.25em;
  }

  @media (min-width: 640px) {
    #content {
      min-width: 400px;
      max-width: 640px;
    }
  }
</style>
