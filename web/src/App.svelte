<script lang="ts">
  import { Router, Route } from "svelte-routing";
  import Home from "./pages/Home.svelte";
  import Schedule from "./pages/Schedule.svelte";

  // This prop keeps the router working correctly if you host on a sub-path
  export let url = "";
</script>

<Router {url}>
  <main id="center" class="w-full max-w-4xl mx-auto p-6">
    <Route path="/">
      <Home />
    </Route>
    <Route path="/schedule">
      <Schedule />
    </Route>
  </main>
</Router>
