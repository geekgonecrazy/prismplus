<script lang="ts">
  import { onMount } from "svelte";

  const default_error_message = "Something went wrong.";

  let admin_key: string = "";

  let connected: boolean = false;
  let fetch_error: boolean = false;
  let err_message: string = default_error_message;

  let show_admin_key = "password";
  let show_new_streamer_key = "password";

  let show_session_keys: string[] = [];
  let show_streamer_keys: string[] = [];

  let sessions = [];
  let streamers = [];

  const default_streamer = {
    name: "",
    streamKey: "",
  };

  let new_streamer: {} = { ...default_streamer };

  function showAdminKey(e) {
    show_admin_key = e.target.checked ? "text" : "password";
  }

  function showNewStreamerKey(e) {
    show_new_streamer_key = e.target.checked ? "text" : "password";
  }

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

  async function onAdminKeyInput(e) {
    admin_key = e.target.value;
  }

  async function onAdminKeyKeyDown(e) {
    const enter_keycode = 13;
    if (e.which === enter_keycode) {
      e.preventDefault();
      await connect();
    }
  }

  async function onNewStreamerKeyInput(e) {
    new_streamer.streamKey = e.target.value;
  }

  async function connect() {
    await getSessions();
    await getStreamers();
  }

  /**
   * Wrapper to provide error handling for fetch requests
   */
  async function fetchHelper(req) {
    let res;
    try {
      res = await req();
    } catch (err) {
      err_message = "Something went wrong fetching resources from the server.";
      fetch_error = true;
      console.error(err.message);
    }

    if (!res.ok) {
      if (res.status === 401) {
        err_message = "Invalid Key.";
      } else {
        console.error(
          `Fetch failed: ${res.url} - ${res.status}: ${res.statusText}`
        );
      }

      fetch_error = true;
      return null;
    }

    fetch_error = false;
    return res;
  }

  /**
   * Helper function to handle HTTP GET requests.
   */
  async function getResource(uri) {
    const res = await fetchHelper(() =>
      fetch(uri, {
        headers: {
          Authorization: `Bearer ${admin_key}`,
        },
      })
    );

    return res?.json();
  }

  /**
   * Helper function to handle HTTP POST requests.
   */
  async function postResource(uri, data) {
    const res = await fetchHelper(() =>
      fetch(uri, {
        method: "POST",
        mode: "cors",
        cache: "no-cache",
        headers: {
          Authorization: `Bearer ${admin_key}`,
          "Content-Type": "application/json",
        },
        body: JSON.stringify(data),
      })
    );

    return res?.json();
  }

  /**
   * Helper function to handle HTTP POST requests.
   */
  async function deleteResource(uri) {
    const res = await fetchHelper(() =>
      fetch(uri, {
        method: "DELETE",
        mode: "cors",
        cache: "no-cache",
        headers: {
          Authorization: `Bearer ${admin_key}`,
          "Content-Type": "application/json",
        },
      })
    );

    return res;
  }

  async function getSessions() {
    const data = await getResource(`/api/v1/sessions`);

    sessions = data || [];
    connected = data ? true : false;
  }

  async function getStreamers() {
    const data = await getResource(`/api/v1/streamers`);

    streamers = data || [];
  }

  async function createStreamer(key: string) {
    await postResource(`/api/v1/streamers`, new_streamer);
    await getStreamers();

    new_streamer = { ...default_streamer };
  }

  async function removeStreamer(key: string) {
    await deleteResource(`/api/v1/streamers/${key}`);
    await getStreamers();
  }

  async function removeRemoteDestination(key: string, dest_id: string) {
    await deleteResource(`/api/v1/sessions/${key}/destinations/${dest_id}`);
    await getSessions();
  }

  async function endSession(key: string) {
    await deleteResource(`/api/v1/sessions/${key}`);
    await getSessions();
  }
</script>

<div id="admin">
  <form>
    {#if !connected}
      <fieldset>
        <legend>Login to Admin</legend>

        <label for="admin-key">Admin Key</label>
        <input
          id="admin-key"
          name="admin-key"
          type={show_admin_key}
          value={admin_key}
          on:input={onAdminKeyInput}
          on:keydown={onAdminKeyKeyDown}
          placeholder="Admin Key"
        />

        <div class="oneline">
          <label class="showkey-label" for="show-admin-key">Show key</label>
          <input
            id="show-admin-key"
            name="show-admin-key"
            type="checkbox"
            on:change={showAdminKey}
          />
        </div>

        <button type="button" on:click={connect}> Connect </button>
      </fieldset>
    {:else}
      <header>
        <h1>Admin Config</h1>
      </header>

      <fieldset>
        <legend>Create new Streamer</legend>

        <label for="new-streamer-name">Streamer Name</label>
        <input
          class="streamer-name"
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
          type={show_new_streamer_key}
          value={new_streamer.streamKey}
          on:input={onNewStreamerKeyInput}
          placeholder="Key for Streamer (blank will cause server to generate)"
        />

        <div class="oneline">
          <label for="show-new-streamer-key">Show key</label>
          <input
            id="show-new-streamer-key"
            name="show-new-streamer-key"
            type="checkbox"
            on:change={showNewStreamerKey}
          />
        </div>

        <button type="button" on:click={createStreamer}>
          Create Streamer
        </button>
      </fieldset>
    {/if}
  </form>

  {#if streamers.length > 0}
    <div id="current-streamers">
      <h2>Streamers</h2>

      <form>
        {#if !fetch_error}
          {#each streamers as { id, name, streamKey }, i}
            <fieldset>
              <legend>Streamer ID: {id}</legend>

              <label for="streamer-name-{i}">Streamer Name</label>
              <input
                class="streamer-name"
                id="streamer-name-{i}"
                name="streamer-name-{i}"
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

              <div class="oneline">
                <label for="show-remote-streamer-key-{i}">Show key</label>
                <input
                  id="show-remote-streamer-key-{i}"
                  name="show-remote-streamer-key-{i}"
                  type="checkbox"
                  on:change={(e) => showStreamerKey(e, i)}
                />
              </div>

              <button
                class="button-negative"
                type="button"
                on:click={() => removeStreamer(id)}>Remove Streamer</button
              >
            </fieldset>
          {/each}
        {/if}
      </form>
    </div>
  {/if}

  {#if sessions.length > 0}
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

              <button
                class="button-negative"
                type="button"
                on:click={() => endSession(key)}>End Session</button
              >
            </fieldset>
          {/each}
        {:else}
          <p style="color: red">{err_message}</p>
        {/if}
      </form>
    </div>
  {/if}
</div>

<style>
  .streamer-name {
    width: 100%;
    max-width: 320px;

    color: var(--alt-color);
    font-size: 1.25rem;
  }
</style>
