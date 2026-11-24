<script lang="ts">
  import { fetchWasteCollectionEvents, getWebCalUrl } from "$lib/api";
  import type { WasteCollectionEvent } from "$lib/api";
  import Button from "$lib/components/ui/button/button.svelte";
  import CardContent from "$lib/components/ui/card/card-content.svelte";
  import CardFooter from "$lib/components/ui/card/card-footer.svelte";
  import CardHeader from "$lib/components/ui/card/card-header.svelte";
  import CardTitle from "$lib/components/ui/card/card-title.svelte";
  import Card from "$lib/components/ui/card/card.svelte";
  import { ArrowLeft, CalendarDays, CalendarPlus, Loader2, Trash2 } from "lucide-svelte";
  import { Link } from "svelte-routing";

  let searchParams = $derived(new URLSearchParams(window.location.search));
  let postcode = $derived(searchParams.get("postcode") || "");
  let uprn = $derived(searchParams.get("uprn") || "");

  let wasteCollectionEvents = $state<WasteCollectionEvent[]>([])
  let isLoading = $state(false);
  let error = $state<string | null>(null);

  const dateFormatter = new Intl.DateTimeFormat("en-GB", {
      weekday: 'short',
      day: 'numeric',
      month: 'short'
    });

  const googleCalendarSubscriptionUrl = $derived(`https://www.google.com/calendar/render?cid=${encodeURIComponent(getWebCalUrl(postcode, uprn))}`);

  function formatDate(isoString: string): string {
    return dateFormatter.format(new Date(isoString))
  }

  $effect(() => {
    if (!postcode || !uprn) {
      return;
    }

    async function fetchCollections() {
      isLoading = true;
      error = null;

      try {
        wasteCollectionEvents = await fetchWasteCollectionEvents(postcode, uprn);
      } catch(err) {
        error = "Failed to load upcoming collections"
      } finally {
        isLoading = false;
      }
    }

    fetchCollections();
  })

</script>

<Card class="w-full max-w-md shadow-lg">
  <CardHeader>
    <CardTitle class="flex items-center gap-2 text-2xl">
      <CalendarDays class="h-6 w-6 text-primary" />
      Collection Schedule
    </CardTitle>
  </CardHeader>
  <CardContent class="space-y-2 text-left">
    <div class="rounded-lg bg-muted p-3 text-sm space-y-1">
      <p><strong>Postcode:</strong> <span class="uppercase">{postcode}</span></p>
      <p><strong>UPRN:</strong> <span class="font-mono text-xs">{uprn}</span></p>
    </div>
    {#if !postcode || !uprn}
    <p class="text-sm text-destructive">Missing lookup information. Please go back and try again.</p>
    {:else}
    <div class="space-y-2">
      <h3 class="text-sm font-medium leading-none">Upcoming Collections</h3>
      {#if isLoading}
        <div class="flex items-center justify-center rounded-md border p-6 text-sm text-muted-foreground">
          <Loader2 class="mr-2 h-4 w-4 animate-spin" />
          Loading dates...
        </div>
      {:else if error}
        <div class="rounded-md border border-destructive/50 bg-destructive/10 p-3 text-sm text-destructive">
          {error}
        </div>
      {:else if wasteCollectionEvents.length === 0}
        <div class="rounded-md border p-4 text-center text-sm text-muted-foreground">
          No upcoming collections found.
        </div>
      {:else}
        <div class="max-h-48 sm:max-h-72 overflow-y-auto rounded-md border divide-y">
          {#each wasteCollectionEvents as event}
            <div class="flex items-center justify-between p-3 text-sm hover:bg-muted/50 transition-colors">
              <div class="flex items-center gap-2.5 font-medium">
                <Trash2 class="h-4 w-4 text-muted-foreground shrink-0" />
                <span>{event.title}</span>
              </div>
              <span class="text-muted-foreground whitespace-nowrap ml-2">
                {formatDate(event.date)}
              </span>
            </div>
          {/each}
        </div>
      {/if}
    </div>

    <div class="space-y-2 text-center flex flex-col gap-1 justify-center">
      <p class="text-sm text-muted-foreground">
        Subscribe to the calendar feed on your device to keep up to date on your upcoming collection dates
      </p>

      <div class="flex flex-col sm:flex-row flex-wrap justify-center gap-1">
        <Button
          href={googleCalendarSubscriptionUrl} 
          target="_blank" 
          rel="noopener noreferrer" 
          class="w-full sm:w-auto gap-2"
        >
          <CalendarPlus class="h-4 w-4" />
          Subscribe on Google Calendar
        </Button>
        <Button
          href={getWebCalUrl(postcode, uprn)}
          rel="noopener noreferrer" 
          class="w-full sm:w-auto gap-2"
        >
          <CalendarPlus class="h-4 w-4" />
          Subscribe on Device Calendar
        </Button>
      </div>
    </div>
    {/if}
  </CardContent>
  <CardFooter class="justify-center border-t pt-4">
    <Link to="/">
      <Button variant="outline">
        <ArrowLeft class="w-4 h-4"/>
        Back to search
      </Button>
    </Link>
  </CardFooter>
</Card>