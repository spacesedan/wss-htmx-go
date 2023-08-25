<script lang="ts">
  import messageStore from "$lib/store";
  import { onMount } from "svelte";
  let connectionMessage=""

  let ws: WebSocket;
  onMount(() => {
    ws = new WebSocket("ws://localhost:8080/ws");
    ws.onopen = (e) => console.log(e);
    ws.onmessage = (e) => {
      const data = JSON.parse(e.data);
      console.log(data);

      switch (data.action) {
        case "connected":
connectionMessage = data.message
      }
    };
  });

  let message: any;
  let messages: any[] = [];
</script>

<h1>Welcome to SvelteKit</h1>
<p>{connectionMessage}</p>
<p>

  Visit <a href="https://kit.svelte.dev">kit.svelte.dev</a> to read the documentation
</p>
