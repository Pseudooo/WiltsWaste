<script lang="ts">
  import { Button } from "$lib/components/ui/button/index";
  import { fetchAddressesForPostcode, type PostcodeAddress } from "$lib/api";
  import Select from "$lib/components/ui/select/select.svelte";
  import SelectTrigger from "$lib/components/ui/select/select-trigger.svelte";
  import SelectContent from "$lib/components/ui/select/select-content.svelte";
  import SelectItem from "$lib/components/ui/select/select-item.svelte";

  let postcode = $state('');
  let addressResults = $state<PostcodeAddress[]>([]);
  let selectedAddressId = $state('');
  let isLoading = $state(false);
  let errorMessage = $state('');

  let selectedAddressLabel = $derived(
    addressResults.find(a => a.uprn === selectedAddressId)?.address || "Select your address"
  );

  async function handleSearch(event: Event) {
    event.preventDefault();
    if (!postcode.trim()) return;

    isLoading = true;
    errorMessage = '';
    addressResults = [];
    selectedAddressId = ''

    try {
      addressResults = await fetchAddressesForPostcode(postcode);
    } catch (error: any) {
      errorMessage = error.message || 'Something went wrong.';
    } finally {
      isLoading = false;
    }
  }
</script>

<div class="flex min-h-screen items-center justify-center bg-zinc-50 p-4 dark:bg-zinc-950">
  <div class="w-full max-w-md rounded-xl border border-zinc-200 bg-white p-6 shadow-sm dark:border-zinc-800 dark:bg-zinc-900">
    
    <div class="space-y-1.5 mb-6">
      <h1 class="text-2xl font-semibold tracking-tight text-zinc-900 dark:text-zinc-50 flex items-center gap-2">
        📅 Bin Day Feed
      </h1>
      <p class="text-sm text-zinc-500 dark:text-zinc-400">
        Enter your postcode to lookup your waste collection dates
      </p>
    </div>

    <form onsubmit={handleSearch} class="space-y-4">
      <div class="space-y-2">
        <label for="postcode" class="text-sm font-medium leading-none text-zinc-950 dark:text-zinc-50">
          Postcode
        </label>
        <div class="flex gap-2">
          <input
            type="text"
            id="postcode"
            placeholder="e.g. SW1A 1AA"
            bind:value={postcode}
            disabled={isLoading}
            required
            class="flex h-10 w-full rounded-md border border-zinc-200 bg-transparent px-3 py-2 text-sm ring-offset-white placeholder:text-zinc-500 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-zinc-950 disabled:cursor-not-allowed disabled:opacity-50 dark:border-zinc-800 dark:ring-offset-zinc-950 dark:placeholder:text-zinc-400 dark:focus-visible:ring-zinc-300"
          />
          <Button type="submit" disabled={isLoading}>
            {isLoading ? 'Searching...' : 'Find'}
          </Button>
        </div>
      </div>
    </form>

    {#if addressResults.length > 0}
      <div class="mt-4 space-y-2 animate-in fade-in-50 duration-200">
        <label for="address-select">
          Select Address
        </label>
        
        <Select type="single" bind:value={selectedAddressId}>
          <SelectTrigger id="address-select" class="w-full text-left">
            {selectedAddressLabel}
          </SelectTrigger>
          <SelectContent>
            {#each addressResults as addr}
              <SelectItem value={addr.uprn} label={addr.address}>
                {addr.address}
              </SelectItem>
            {/each}
          </SelectContent>
        </Select>
      </div>
    {/if}

    {#if errorMessage}
      <div class="mt-4 rounded-lg border border-red-200 bg-red-50 p-3 text-sm text-red-600 dark:border-red-900/50 dark:bg-red-950/50 dark:text-red-400">
        {errorMessage}
      </div>
    {/if}
  </div>
</div>
