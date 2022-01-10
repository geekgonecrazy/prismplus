<script lang="ts">
    import { onMount } from "svelte";

    const default_error_message = "Something went wrong.";

    let streamer_key: string = "";

    let streamer = {
        name: "",
        destinations: [],
    };

    let connected = false;

    let fetch_error = false;
    let err_message = default_error_message;

    let show_streamer_key = "password";
    let show_stream_key = "password";
    let show_new_destination_key = "password";

    const default_destination = {
        name: "",
        server: "",
        key: "",
    };

    let new_destination: {} = { ...default_destination };

    function showStreamerKey(e) {
        show_streamer_key = e.target.checked ? "text" : "password";
    }

    function showStreamKey(e) {
        show_stream_key = e.target.checked ? "text" : "password";
    }

    function showNewDestinationKey(e) {
        show_new_destination_key = e.target.checked ? "text" : "password";
    }

    async function onStreamerKeyInput(e) {
        streamer_key = e.target.value;
    }

    async function onStreamerKeyKeyDown(e) {
        const enter_keycode = 13;
        if (e.which === enter_keycode) {
            e.preventDefault();
            await connect();
        }
    }

    async function onNewDestinationKeyInput(e) {
        new_destination.key = e.target.value;
    }

    async function connect() {
        await getStreamer();
    }

    /**
     * Wrapper to provide error handling for fetch requests
     */
    async function fetchHelper(req) {
        let res;
        try {
            res = await req();
        } catch (err) {
            err_message =
                "Something went wrong fetching resources from the server.";
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
                    Authorization: `Bearer ${streamer_key}`,
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
                    Authorization: `Bearer ${streamer_key}`,
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(data),
            })
        );

        return res;
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
                    Authorization: `Bearer ${streamer_key}`,
                    "Content-Type": "application/json",
                },
            })
        );

        return res;
    }

    async function addDestination() {
        await postResource(`/api/v1/streamer/destinations`, new_destination);
        new_destination = { ...default_destination };
        await getStreamer();
    }

    async function removeDestination(dest_id: string) {
        await deleteResource(`/api/v1/streamer/destinations/${dest_id}`);
        await getStreamer();
    }

    async function getStreamer() {
        const data = await getResource(`/api/v1/streamer`);
        streamer = data;
        connected = true;
    }

    onMount(async () => {
        // if we store the streamer_key some how in session.. then we could probably add this back
        //await getStreamer();
    });
</script>

<div id="streamer">
    {#if !connected}
        <form>
            <fieldset>
                <legend>Login as Streamer</legend>

                <label for="login-streamer-key">Streamer Key</label>
                <input
                    id="login-streamer-key"
                    name="login-streamer-key"
                    type={show_streamer_key}
                    value={streamer_key}
                    on:keydown={onStreamerKeyKeyDown}
                    on:input={onStreamerKeyInput}
                    placeholder="Streamer Key"
                />

                <div class="oneline">
                    <label for="show-admin-key">Show key</label>
                    <input
                        id="show-admin-key"
                        name="show-admin-key"
                        type="checkbox"
                        on:change={showStreamerKey}
                    />
                </div>

                <button type="button" on:click={connect}> Connect </button>
            </fieldset>
        </form>
    {:else}
        <div id="current-sessions">
            <header>
                <h1>Stream Config</h1>
                <h2 id="streamer-name">{streamer.name}</h2>
            </header>

            <form>
                {#if !fetch_error}
                    <fieldset>
                        <legend>Session</legend>

                        <label for="stream-key">Stream Key</label>
                        <input
                            id="stream-key"
                            name="stream-key"
                            type={show_stream_key}
                            value={streamer.streamKey}
                            readonly
                        />

                        <div class="oneline">
                            <label for="show-streamer-key">Show key</label>
                            <input
                                id="show-streamer-key"
                                name="show-streamer-key"
                                type="checkbox"
                                on:change={showStreamKey}
                            />
                        </div>

                        <h3>Destinations</h3>
                        <ul>
                            {#each Object.entries(streamer.destinations) as [_, { name, server, key, id }]}
                                <li>
                                    <button
                                        class="button-negative"
                                        type="button"
                                        on:click={() => removeDestination(id)}
                                        >Remove</button
                                    >
                                    <span class="destination-entry">{name} - {server}</span>
                                </li>
                            {/each}
                        </ul>

                        <label for="new-destination-name"
                            >New Remote Destination Name</label
                        >
                        <input
                            class="destination-name"
                            id="new-destination-name"
                            name="new-destination-name"
                            bind:value={new_destination.name}
                            placeholder="My Owncast Server"
                        />

                        <label for="new-destination-server"
                            >New Remote Destination Server</label
                        >
                        <input
                            class="destination-server"
                            id="new-destination-server"
                            name="new-destination-server"
                            bind:value={new_destination.server}
                            placeholder="rtmp://host/live"
                        />

                        <label for="new-destination-key"
                            >New Remote Destination Key</label
                        >
                        <input
                            id="new-destination-key"
                            name="new-destination-key"
                            type={show_new_destination_key}
                            value={new_destination.key}
                            on:input={onNewDestinationKeyInput}
                            placeholder="Destination Key"
                        />

                        <div class="oneline">
                            <label for="show-new-destination-key"
                                >Show key</label
                            >
                            <input
                                id="show-new-destination-key"
                                name="show-new-destination-key"
                                type="checkbox"
                                on:change={showNewDestinationKey}
                            />
                        </div>

                        <button type="button" on:click={() => addDestination()}
                            >Add Remote Destination</button
                        >
                    </fieldset>
                {:else}
                    <p style="color: red">
                        {err_message}
                    </p>
                {/if}
            </form>
        </div>
    {/if}
</div>

<style>
    header {
        text-align: center;
    }

    ul {
        padding-left: 0;
    }

    .destination-name {
        width: 100%;
        max-width: 320px;
    }
    .destination-server {
        width: 100%;
        max-width: 640px;
    }

    .destination-entry {
        color: var(--white);
    }

    #streamer-name {
        font-size: 1.25em;
        font-weight: bold;
        color: var(--alt-color);
    }
</style>
