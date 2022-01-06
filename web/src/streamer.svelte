<script lang="ts">
    import { onMount } from "svelte";

    let stream_key: string = "";

    let streamer = {
        destinations: [],
    };

    let connected = false;

    let fetch_error = false;

    let show_stream_key = "password";

    const default_destination = {
        name: "",
        server: "",
        key: "",
    };

    let new_destination: {} = default_destination;

    function showStreamKey(e) {
        show_stream_key = e.target.checked ? "text" : "password";
    }

    async function connect() {
        await getStreamer();
        connected = true;
    }

    async function addDestination() {
        const url = `/api/v1/streamer/destinations`;

        try {
            await fetch(`/api/v1/streamer/destinations`, {
                method: "POST",
                mode: "cors",
                cache: "no-cache",
                headers: {
                    Authorization: `Bearer ${stream_key}`,
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(new_destination),
            });

            new_destination = {};
            new_destination = default_destination;

            await getStreamer();
            fetch_error = false;
        } catch (err) {
            fetch_error = true;
            console.error(err);
        }
    }

    async function removeDestination(dest_id: string) {
        try {
            await fetch(`/api/v1/streamer/destinations/${dest_id}`, {
                method: "DELETE",
                mode: "cors",
                cache: "no-cache",
                headers: {
                    Authorization: `Bearer ${stream_key}`,
                    "Content-Type": "application/json",
                },
            });

            await getStreamer();
            fetch_error = false;
        } catch (err) {
            fetch_error = true;
            console.error(err);
        }
    }

    async function getStreamer() {
        try {
            const res = await fetch(`/api/v1/streamer`, {
                method: "GET",
                mode: "cors",
                cache: "no-cache",
                headers: {
                    Authorization: `Bearer ${stream_key}`,
                    "Content-Type": "application/json",
                },
            });

            streamer = await res.json();
            fetch_error = false;
        } catch (err) {
            fetch_error = true;
            console.error(err);
        }
    }

    onMount(async () => {
        // if we store the stream_key some how in session.. then we could probably add this back
        //await getStreamer();
    });
</script>

<div id="main">
    <h1>Prism+</h1>

    {#if !connected}
    <form>
        <fieldset>
            <legend>Login as Streamer</legend>
            <label for="login-streamer-key">Stream Key</label>
            <input
                id="login-streamer-key"
                name="login-streamer-key"
                type="password"
                bind:value={stream_key}
                placeholder="Stream Key"
            />

            <button type="button" on:click={connect}> Connect </button>
        </fieldset>
    </form>
    {:else}

    <div id="current-sessions">
        <h2>Stream Config</h2>

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
                    <label for="show-streamer-key">Show key</label>
                    <input
                        id="show-streamer-key"
                        name="show-streamer-key"
                        type="checkbox"
                        on:change={(e) => showStreamKey(e)}
                    />

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
                                <span>{name} - {server}</span>
                            </li>
                        {/each}
                    </ul>

                    <br />

                    <label for="new-destination-name"
                        >New Remote Destination Name</label
                    >
                    <input
                        id="new-destination-name"
                        name="new-destination-name"
                        bind:value={new_destination.name}
                        placeholder="My Owncast Server"
                    />

                    <label for="new-destination-server"
                        >New Remote Destination Server</label
                    >
                    <input
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
                        type="text"
                        bind:value={new_destination.key}
                        placeholder="Destination Key"
                    />

                    <button type="button" on:click={() => addDestination()}
                        >Add Remote Destination</button
                    >
                </fieldset>
            {:else}
                <p style="color: red">There was an error fetching resources.</p>
            {/if}
        </form>
    </div>
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
