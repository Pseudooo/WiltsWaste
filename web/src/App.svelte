<script lang="ts">
  import Check from "lucide-svelte/icons/check";
  import ChevronsUpDown from "lucide-svelte/icons/chevrons-up-down";
  import * as Command from "$lib/components/ui/command/index";
  import * as Popover from "$lib/components/ui/popover/index";
  import { Button } from "$lib/components/ui/button/index";
  import { cn } from "$lib/utils";
  import { tick } from "svelte";
    import { fetchAddressesForPostcode } from "$lib/api";

  // State Management
  let postcode = $state('');
  let addressResults = $state<any[]>([]);
  let isLoading = $state(false);
  let errorMessage = $state('');
  let selectedIcsUrl = $state('');

  // Combobox Popover UI State
  let open = $state(false);
  let selectedValue = $state("");
  let selectedLabel = $state("Select your address...");

  // Triggered when an address is picked from the Combobox
  function handleSelectAddress(currentValue: string, currentLabel: string) {
    selectedValue = currentValue;
    selectedLabel = currentLabel;
    open = false;

    if (!currentValue) {
      selectedIcsUrl = '';
      return;
    }
    selectedIcsUrl = `https://your-api.com{currentValue}.ics`;
  }

  // Search for addresses via Postcode API
  async function handleSearch(event: Event) {
    event.preventDefault();
    if (!postcode.trim()) return;

    isLoading = true;
    errorMessage = '';
    addressResults = [];
    selectedIcsUrl = '';
    selectedValue = '';
    selectedLabel = "Select your address...";

    try {
      addressResults = await fetchAddressesForPostcode(postcode);
    } catch (error: any) {
      errorMessage = error.message || 'Something went wrong.';
    } finally {
      isLoading = false;
    }
  }

  // Clipboard Utility
  function copyToClipboard() {
    navigator.clipboard.writeText(selectedIcsUrl);
    alert('Link copied to clipboard!'); 
  }
</script>

<div class="flex min-h-screen items-center justify-center bg-zinc-50 p-4 dark:bg-zinc-950">
  <div class="w-full max-w-md rounded-xl border border-zinc-200 bg-white p-6 shadow-sm dark:border-zinc-800 dark:bg-zinc-900">
    
    <!-- Header -->
    <div class="space-y-1.5 mb-6">
      <h1 class="text-2xl font-semibold tracking-tight text-zinc-900 dark:text-zinc-50 flex items-center gap-2">
        📅 Bin Day Feed
      </h1>
      <p class="text-sm text-zinc-500 dark:text-zinc-400">
        Enter your postcode to sync your collection schedule directly to your calendar.
      </p>
    </div>

    <!-- Search Form -->
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

    <!-- Error Alert Banner -->
    {#if errorMessage}
      <div class="mt-4 rounded-lg border border-red-200 bg-red-50 p-3 text-sm text-red-600 dark:border-red-900/50 dark:bg-red-950/50 dark:text-red-400">
        {errorMessage}
      </div>
    {/if}

    <!-- Shadcn Combobox / Filterable Address Picker -->
    {#if addressResults.length > 0}
      <div class="mt-6 space-y-2 flex flex-col">
        <label for="address-combobox" class="text-sm font-medium leading-none text-zinc-950 dark:text-zinc-50">
          Select Your Address
        </label>
        
        <Popover.Root bind:open>
          <Popover.Trigger asChild let:builder>
            <Button
              builders={[builder]}
              variant="outline"
              role="combobox"
              aria-expanded={open}
              class="w-full justify-between font-normal text-left"
            >
              <span class="truncate">{selectedLabel}</span>
              <ChevronsUpDown class="ml-2 h-4 w-4 shrink-0 opacity-50" />
            </Button>
          </Popover.Trigger>
          <Popover.Content class="w-[var(--bits-popover-anchor-width)] p-0">
            <Command.Root>
              <Command.Input placeholder="Type to filter addresses..." />
              <Command.Empty>No address found.</Command.Empty>
              <Command.Group class="max-h-60 overflow-y-auto">
                {#each addressResults as item}
                  <Command.Item
                    value={item.address}
                    onSelect={() => handleSelectAddress(item.id, item.address)}
                  >
                    <Check
                      class={cn(
                        "mr-2 h-4 w-4",
                        selectedValue !== item.id && "text-transparent"
                      )}
                    />
                    <span class="truncate">{item.address}</span>
                  </Command.Item>
                {/each}
              </Command.Group>
            </Command.Root>
          </Popover.Content>
        </Popover.Root>
      </div>
    {/if}

    <!-- Final Success Layout & Copy Actions -->
    {#if selectedIcsUrl}
      <div class="mt-6 space-y-4 rounded-lg border border-emerald-200 bg-emerald-50/50 p-4 dark:border-emerald-900/30 dark:bg-emerald-950/20">
        <div class="space-y-1">
          <h3 class="text-sm font-semibold text-emerald-900 dark:text-emerald-400">🎉 Feed Ready!</h3>
          <p class="text-xs text-emerald-700 dark:text-emerald-500">
            Copy this link or subscribe directly in your calendar settings.
          </p>
        </div>
        
        <div class="flex gap-2">
          <input
            type="text"
            value={selectedIcsUrl}
            readonly
            class="flex h-9 w-full rounded-md border border-zinc-200 bg-white px-3 py-1 text-xs text-zinc-600 focus-visible:outline-none dark:border-zinc-800 dark:bg-zinc-950 dark:text-zinc-400"
          />
          <Button
            variant="outline"
            size="sm"
            onclick={copyToClipboard}
            class="h-9 px-3 text-xs"
          >
            Copy
          </Button>
        </div>

        <a
          href={selectedIcsUrl}
          download="bindays.ics"
          class="inline-flex h-9 w-full items-center justify-center rounded-md bg-emerald-600 text-xs font-medium text-white hover:bg-emerald-600/90 transition-colors"
        >
          Download .ics File
        </a>
      </div>
    {/if}

  </div>
</div>
