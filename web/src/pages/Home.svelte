<script lang="ts">
  import { Button } from "$lib/components/ui/button/index";
  import { fetchAddressesForPostcode, type PostcodeAddress } from "$lib/api";
  import Select from "$lib/components/ui/select/select.svelte";
  import SelectTrigger from "$lib/components/ui/select/select-trigger.svelte";
  import SelectContent from "$lib/components/ui/select/select-content.svelte";
  import SelectItem from "$lib/components/ui/select/select-item.svelte";
  import Card from "$lib/components/ui/card/card.svelte";
  import CardHeader from "$lib/components/ui/card/card-header.svelte";
  import CardTitle from "$lib/components/ui/card/card-title.svelte";
  import { CalendarDays } from "lucide-svelte";
  import CardContent from "$lib/components/ui/card/card-content.svelte";
  import { Label } from "$lib/components/ui/label";
  import { Input } from "$lib/components/ui/input";

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

<Card class="w-full max-w-md shadow-lg">
  <CardHeader>
    <CardTitle class="flex items-center gap-2 text-2xl">
      <CalendarDays class="h-6 w-6 text-primary" />
      Bin Day Feed 
    </CardTitle>
  </CardHeader>
  <CardContent class="space-y-4">
    <form onsubmit={handleSearch} class="space-y-4">
      <div class="space-y-2">
        <Label for="postcode">Postcode</Label>
                <div class="flex gap-2">
          <Input
            type="text"
            id="postcode"
            placeholder="e.g. SW1A 1AA"
            bind:value={postcode}
            disabled={isLoading}
            required
            class="uppercase placeholder:normal-case"
          />
          <Button type="submit" disabled={isLoading}>
            {#if isLoading}
              <Loader2 class="mr-2 h-4 w-4 animate-spin" />
              Searching...
            {:else}
              Find
            {/if}
          </Button>
        </div>
      </div>
    </form>

    {#if addressResults.length > 0}
    <div class="mt-4 space-y-2 animate-in fade-in-50 duration-200">
      <Label for="address-select">
        Select Address
      </Label>
      
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
  </CardContent>
</Card>