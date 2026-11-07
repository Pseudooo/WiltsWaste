<script lang="ts">
  import { Router, Route } from "svelte-routing";
  import Home from "./pages/Home.svelte";
  import Schedule from "./pages/Schedule.svelte";
  import { GithubIcon } from "@hugeicons/core-free-icons";

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

  <footer class="w-full max-w-4xl mx-auto p-6 border-t border-border text-center text-sm text-muted-foreground">
    <p>Made by Pseudooo</p>
    <div class="mt-2 flex justify-center gap-4">
      <a 
        href="https://github.com/Pseudooo/WiltsWaste"
        target="_blank"
        rel="noreferrer"
        class="flex items-center gap-1.5 hover:text-foreground transition-colors"
      >
        <img src="/github.svg" alt="Github logo" class="w-4 h-4 opacity-70 group-hover:opacity-100 transition-opacity" />
        <span>Github</span>
      </a>
      <a
        href="https://www.wiltshire.gov.uk/rubbish-and-recycling"
        target="_blank"
        rel="noreferrer"
        class="flex items-center gap-1.5 hover:text-foreground transition-colors"
      >
        <img src="/wiltshire-council-logo.svg" alt="Wiltshire Council Logo" class="w-4 h-4 opacity-70 group-hover:opacity-100 transition-opacity" />
        <span>Wiltshire Council</span>
      </a>
    </div>
  </footer>
</Router>
