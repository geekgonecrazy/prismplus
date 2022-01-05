<script lang="ts">
  import { onMount } from "svelte";

  const API_SESSIONS = "api/v1/sessions";

  let protocol: string = "http";
  let server_address: string = "localhost";
  let server_port: number = 5383;

  let session_key_type = "password";

  let sessions: [];
  let sessions_promise: Promise<void> = getSessions();

  let new_urls: string[] = [""];
  let new_keys: string[] = [""];
  let show_keys: string[] = ["password"];

  let new_remote_urls: string[] = [];
  let new_remote_keys: string[] = [];
  let show_remote_session_keys: string[] = [];
  let show_remote_dest_keys: string[] = [];

  let show_sessions_debug = false;

  let default_session = {
    key: "",
    destinations: [],
  };
  let new_session: {} = default_session;

  function showSessionKey(e) {
    session_key_type = e.target.checked ? "text" : "password";
  }
  function showKey(e, i) {
    show_keys[i] = e.target.checked ? "text" : "password";
  }
  function showRemoteSessionKey(e, i) {
    show_remote_session_keys[i] = e.target.checked ? "text" : "password";
  }
  function showRemoteDestKey(e, i) {
    show_remote_dest_keys[i] = e.target.checked ? "text" : "password";
  }

  function onInputNewSessionKey(e) {
    new_session.key = e.target.value;
  }
  function onInputNewKey(e, i) {
    new_keys[i] = e.target.value;
  }
  function onInputRemoteDestKey(e, i) {
    new_remote_keys[i] = e.target.value;
  }

  function addDestination() {
    let cp_urls = [...new_urls];
    let cp_keys = [...new_keys];
    let cp_show_keys = [...show_keys];
    cp_urls.push("");
    cp_keys.push("");
    cp_show_keys.push("password");
    new_urls = cp_urls;
    new_keys = cp_keys;
    show_keys = cp_show_keys;
  }
  function removeDestination(i: number) {
    let cp_urls = [...new_urls];
    let cp_keys = [...new_keys];
    let cp_show_keys = [...show_keys];
    cp_urls.splice(i, 1);
    cp_keys.splice(i, 1);
    cp_show_keys.splice(i, 1);
    new_urls = cp_urls;
    new_keys = cp_keys;
    show_keys = cp_show_keys;
  }

  async function createSession() {
    const url = `${protocol}://${server_address}:${server_port}/${API_SESSIONS}`;
    for (let i = 0; i < new_urls.length; i++) {
      new_session.destinations.push({ url: `${new_urls[i]}/${new_keys[i]}` });
    }
    try {
      const promise = await fetch(url, {
        method: "POST",
        mode: "cors",
        cache: "no-cache",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(new_session),
      });

      await getSessions();
      new_session = default_session;
      new_urls = [""];
      new_keys = [""];
      show_keys = ["password"];
    } catch (err) {
      console.error(err);
    }
  }

  async function endSession(key: string) {
    const url = `${protocol}://${server_address}:${server_port}/${API_SESSIONS}/${key}`;
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
    } catch (err) {
      console.error(err);
    }
  }

  async function addRemoteDestination(
    key: string,
    dest_url: string,
    dest_key: string
  ) {
    const url = `${protocol}://${server_address}:${server_port}/${API_SESSIONS}/${key}/destinations`;
    const data = { url: `${dest_url}/${dest_key}` };
    try {
      const promise = await fetch(url, {
        method: "POST",
        mode: "cors",
        cache: "no-cache",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
      });

      await getSessions();
    } catch (err) {
      console.error(err);
    }
  }

  async function removeRemoteDestination(key: string, dest_id: string) {
    const url = `${protocol}://${server_address}:${server_port}/${API_SESSIONS}/${key}/destinations/${dest_id}`;
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
    } catch (err) {
      console.error(err);
    }
  }

  async function getSessions() {
    const res = await fetch(
      `${protocol}://${server_address}:${server_port}/${API_SESSIONS}`
    );

    if (res.ok) {
      sessions = await res.json();
      let tmp_urls = [];
      let tmp_keys = [];
      let tmp_shows = [];
      for (let i = 0; i < sessions.length; i++) {
        tmp_urls.push("");
        tmp_keys.push("");
        tmp_shows.push("password");
      }
      new_remote_urls = tmp_urls;
      new_remote_keys = tmp_keys;
      show_remote_session_keys = tmp_shows;
      show_remote_dest_keys = tmp_shows;
    }
  }
</script>

<div id="main">
  <h1>Prism+</h1>

  <form>
    <fieldset>
      <legend>Server information</legend>

      <label for="server-address">Server Address</label>
      <input
        id="server-address"
        name="server-address"
        bind:value={server_address}
        placeholder="Server Address"
      />

      <label for="server-port">Server Port</label>
      <input
        id="server-port"
        name="server-port"
        type="number"
        bind:value={server_port}
        placeholder="Server Port"
      />
    </fieldset>

    <fieldset>
      <legend>Create new session</legend>

      <label for="session-key">Session Key</label>
      <input
        id="session-key"
        name="session-key"
        type={session_key_type}
        value={new_session.key}
        placeholder="New Session Key"
        on:input={onInputNewSessionKey}
      />
      <label for="show-session-key">Show key</label>
      <input
        id="show-session-key"
        name="show-session-key"
        type="checkbox"
        on:change={showSessionKey}
      />
      <br />

      {#each { length: new_urls.length } as _, i}
        <label for="dest-{i}">Destination URL {i}</label>
        <input
          name="dest-{i}"
          bind:value={new_urls[i]}
          placeholder="Destination URL"
        />

        <label for="key-{i}">Destination Key {i}</label>
        <input
          id="key-{i}"
          name="key-{i}"
          type={show_keys[i]}
          value={new_keys[i]}
          placeholder="Destination Key"
          on:input={(e) => onInputNewKey(e, i)}
        />

        <label for="show-dest-key-{i}">Show key</label>
        <input
          id="show-dest-key-{i}"
          name="show-dest-key-{i}"
          type="checkbox"
          on:change={(e) => showKey(e, i)}
        />

        {#if new_urls.length > 1}
          <button
            class="button-negative"
            type="button"
            on:click={() => removeDestination(i)}
          >
            Remove #{i}
          </button>
        {/if}
      {/each}

      <br />
      <button type="button" on:click={addDestination}>
        Add Another Destination
      </button>
      <br />

      <button type="button" on:click={createSession}> Create Session </button>
    </fieldset>
  </form>

  <div id="current-sessions">
    <h2>Current Sessions</h2>

    <form>
      {#await sessions_promise}
        <p>Fetching session info...</p>
      {:then}
        {#each sessions as { key, destinations, nextDestinationId, active, end, streamHeaders }, i}
          <fieldset>
            <legend>Session {i}</legend>

            <label for="session-key-{i}">Session Key</label>
            <input
              id="session-key-{i}"
              name="session-key-{i}"
              type={show_remote_session_keys[i]}
              value={key}
              readonly
            />
            <label for="show-remote-session-key-{i}">Show key</label>
            <input
              id="show-remote-session-key-{i}"
              name="show-remote-session-key-{i}"
              type="checkbox"
              on:change={(e) => showRemoteSessionKey(e, i)}
            />

            <h3>Destinations</h3>
            <ul>
              {#each Object.entries(destinations) as [_, { url, id }]}
                <li>
                  <button
                    class="button-negative"
                    type="button"
                    on:click={() => removeRemoteDestination(key, id)}
                    >Remove</button
                  >
                  <span>{url.split("/").slice(0, 1)}</span>
                </li>
              {/each}
            </ul>

            <br />
            <label for="remote-destination-url-{i}"
              >New Remote Destination URL</label
            >
            <input
              id="remote-destination-url-{i}"
              name="remote-destination-url-{i}"
              bind:value={new_remote_urls[i]}
              placeholder="Destination URL"
            />

            <label for="remote-destination-key-{i}"
              >New Remote Destination Key</label
            >
            <input
              id="remote-destination-key-{i}"
              name="remote-destination-key-{i}"
              type={show_remote_dest_keys[i]}
              value={new_remote_keys[i]}
              on:input={(e) => onInputRemoteDestKey(e, i)}
              placeholder="Destination Key"
            />
            <label for="show-remote-dest-key-{i}">Show key</label>
            <input
              id="show-remote-dest-key-{i}"
              name="show-remote-dest-key-{i}"
              type="checkbox"
              on:change={(e) => showRemoteDestKey(e, i)}
            />

            <button
              type="button"
              on:click={() =>
                addRemoteDestination(
                  key,
                  new_remote_urls[i],
                  new_remote_keys[i]
                )}>Add Remote Destination</button
            >

            <br />

            <!-- Disabling this button until end session is useful -->
            <button
              class="button-negative"
              type="button"
              on:click={() => endSession(key)}
              disabled>End Session</button
            >
          </fieldset>
        {/each}
      {:catch err}
        <p style="color: red">{err.message}</p>
      {/await}
    </form>
  </div>

  <label style="color: red;" for="show-sessions-debug"
    >Show Sessions Debug (will expose keys)</label
  >
  <input
    id="show-sessions-debug"
    name="show-sessions-debug"
    type="checkbox"
    bind:checked={show_sessions_debug}
  />
  {#if show_sessions_debug}
    <pre style="text-align: left;">{JSON.stringify(sessions, null, 2)}</pre>
  {/if}
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
