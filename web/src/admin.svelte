<script lang="ts">
  import { onMount } from "svelte";

  let admin_key: string = "";

  let connected: boolean = false;
  let fetch_error: boolean = false;

  let show_session_keys: string[] = [];
  let show_streamer_keys: string[] = [];

  let sessions = [];
  let streamers = [];

  const default_streamer = {
    name: "",
    streamKey: "",
  };

  let new_streamer: {} = default_streamer;

  function showSessionKey(e, i) {
    show_session_keys[i] = e.target.checked ? "text" : "password";
  }

  function showStreamerKey(e, i) {
    show_streamer_keys[i] = e.target.checked ? "text" : "password";
  }

  $: getStreamerKeyVisibility = (i) => {
    return show_streamer_keys[i] || "password";
  };

  $: getSessionKeyVisibility = (i) => {
    return show_session_keys[i] || "password";
  };

  async function connect() {
    await getSessions();
    await getStreamers();
  }

  async function getSessions() {
    try {
      const res = await fetch(
        `/api/v1/sessions`,
        {
          headers: {
            Authorization: `Bearer ${admin_key}`,
          },
        }
      );

      console.log(admin_key);

      sessions = await res.json();
      connected = true;

      fetch_error = false;
    } catch (err) {
      fetch_error = true;
      console.error(err);
    }
  }

  async function getStreamers() {
    try {
      const res = await fetch(
        `/api/v1/streamers`,
        {
          headers: {
            Authorization: `Bearer ${admin_key}`,
          },
        }
      );

      console.log(admin_key);

      streamers = await res.json();

      fetch_error = false;
    } catch (err) {
      fetch_error = true;
      console.error(err);
    }
  }

  async function createStreamer(key: string) {
    const url = `/api/v1/streamers`;
    try {
      const promise = await fetch(url, {
        method: "POST",
        mode: "cors",
        cache: "no-cache",
        headers: {
          Authorization: `Bearer ${admin_key}`,
          "Content-Type": "application/json",
        },
        body: JSON.stringify(new_streamer),
      });

      await getStreamers();
      fetch_error = false;
    } catch (err) {
      fetch_error = true;
      console.error(err);
    }
  }

  async function removeStreamer(key: string) {
    const url = `/api/v1/streamers/${key}`;
    try {
      const promise = await fetch(url, {
        method: "DELETE",
        mode: "cors",
        cache: "no-cache",
        headers: {
          Authorization: `Bearer ${admin_key}`,
          "Content-Type": "application/json",
        },
      });

      await getStreamers();
      fetch_error = false;
    } catch (err) {
      fetch_error = true;
      console.error(err);
    }
  }

  async function removeRemoteDestination(key: string, dest_id: string) {
    const url = `/api/v1/sessions/${key}/destinations/${dest_id}`;
    try {
      const promise = await fetch(url, {
        method: "DELETE",
        mode: "cors",
        cache: "no-cache",
        headers: {
          "Content-Type": "application/json",
        },
      });

      await getSessions();
      fetch_error = false;
    } catch (err) {
      fetch_error = true;
      console.error(err);
    }
  }

  async function endSession(key: string) {
    const url = `/api/v1/sessions/${key}`;
    try {
      const promise = await fetch(url, {
        method: "DELETE",
        mode: "cors",
        cache: "no-cache",
        headers: {
          Authorization: `Bearer ${admin_key}`,
          "Content-Type": "application/json",
        },
      });

      await getSessions();
      fetch_error = false;
    } catch (err) {
      fetch_error = true;
      console.error(err);
    }
  }
</script>

<div id="main">
  <h1>Prism+</h1>

  <form>
    {#if !connected}
    <fieldset>
      <legend>Login to Admin</legend>

      <label for="server-port">Admin Key</label>
      <input
        id="server-port"
        name="server-port"
        type="password"
        bind:value={admin_key}
        placeholder="Admin Key"
      />

      <button type="button" on:click={connect}> Connect </button>
    </fieldset>

    {:else}
    <fieldset>
      <legend>Create new Streamer</legend>

      <label for="new-streamer-name">Streamer Name</label>
      <input
        id="new-streamer-name"
        name="new-streamer-name"
        type="text"
        bind:value={new_streamer.name}
        placeholder="Name of Streamer"
      />

      <label for="new-streamer-key">Streamer Key</label>
      <input
        id="new-streamer-key"
        name="new-streamer-key"
        type="text"
        bind:value={new_streamer.streamKey}
        placeholder="Key for Streamer (blank will cause server to generate)"
      />
      <br />

      <button type="button" on:click={createStreamer}> Create Streamer </button>
    </fieldset>
    {/if}
  </form>

  <div id="current-streamers">
    <h2>Streamers</h2>

    <form>
      {#if !fetch_error}
        {#each streamers as { id, name, streamKey }, i}
          <fieldset>
            <legend>Streamer {i}</legend>

            <label for="streamer-key-{i}">Streamer Name</label>
            <input
              id="streamer-key-{i}"
              name="streamer-key-{i}"
              type="text"
              value={name}
              readonly
            />

            <label for="streamer-key-{i}">Streamer Key</label>
            <input
              id="streamer-key-{i}"
              name="streamer-key-{i}"
              type={getStreamerKeyVisibility(i)}
              value={streamKey}
              readonly
            />
            <label for="show-remote-session-key-{i}">Show key</label>
            <input
              id="show-remote-session-key-{i}"
              name="show-remote-session-key-{i}"
              type="checkbox"
              on:change={(e) => showStreamerKey(e, i)}
            />

            <br />

            <button
              class="button-negative"
              type="button"
              on:click={() => removeStreamer(id)}>Remove Streamer</button
            >
          </fieldset>
        {/each}
      {:else}
        <p style="color: red">There was an error fetching resources.</p>
      {/if}
    </form>
  </div>

  <div id="current-sessions">
    <h2>Current Sessions</h2>

    <form>
      {#if !fetch_error}
        {#each sessions as { key, destinations, nextDestinationId, active, end, streamHeaders }, i}
          <fieldset>
            <legend>Session {i}</legend>

            <label for="session-key-{i}">Session Key</label>
            <input
              id="session-key-{i}"
              name="session-key-{i}"
              type={getSessionKeyVisibility(i)}
              value={key}
              readonly
            />
            <label for="show-remote-session-key-{i}">Show key</label>
            <input
              id="show-remote-session-key-{i}"
              name="show-remote-session-key-{i}"
              type="checkbox"
              on:change={(e) => showSessionKey(e, i)}
            />

            <h3>Destinations</h3>
            <ul>
              {#each Object.entries(destinations) as [_, { name, server, id }]}
                <li>
                  <button
                    class="button-negative"
                    type="button"
                    on:click={() => removeRemoteDestination(key, id)}
                    >Remove</button
                  >
                  <span>{name} - {server}</span>
                </li>
              {/each}
            </ul>

            <br />
            <br />

            <button
              class="button-negative"
              type="button"
              on:click={() => endSession(key)}>End Session</button
            >
          </fieldset>
        {/each}
      {:else}
        <p style="color: red">There was an error fetching resources.</p>
      {/if}
    </form>
  </div>
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

  .button-negative {
    background-color: #ff7777;
  }

  #main {
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
    #main {
      max-width: 640px;
    }
  }
</style>
