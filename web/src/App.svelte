<script lang="ts">
  import { Router, Link, Route } from "svelte-routing";
  import Home from "./pages/Home.svelte";

  // This prop keeps the router working correctly if you host on a sub-path
  export let url = "";
</script>

<Router {url}>
  <nav class="flex gap-4 p-4 border-b border-border bg-card">
    <Link to="/" class="text-primary hover:underline font-medium">Home</Link>
  </nav>

  <main id="center" class="w-full max-w-4xl mx-auto p-6">
    <Route path="/">
      <Home />
    </Route>
  </main>
</Router>
