<script lang="ts">
  import messageStore from "$lib/store"
    import { onMount } from "svelte";

let ws: WebSocket
  onMount(() => {
ws = new WebSocket("ws://localhost:8080/ws")
    })



  let message: any;
  let messages: any[] = [];
</script>

<h1>Welcome to SvelteKit</h1>
<p>
  Visit <a href="https://kit.svelte.dev">kit.svelte.dev</a> to read the documentation
</p>
