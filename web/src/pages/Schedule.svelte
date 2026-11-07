<script lang="ts">
  import { getWebCalUrl } from "$lib/api";
  import Button from "$lib/components/ui/button/button.svelte";
  import CardContent from "$lib/components/ui/card/card-content.svelte";
  import CardHeader from "$lib/components/ui/card/card-header.svelte";
  import CardTitle from "$lib/components/ui/card/card-title.svelte";
  import Card from "$lib/components/ui/card/card.svelte";
  import { CalendarDays, CalendarPlus } from "lucide-svelte";

  let searchParams = $derived(new URLSearchParams(window.location.search));
  let postcode = $derived(searchParams.get("postcode") || "");
  let uprn = $derived(searchParams.get("uprn") || "");

  const googleCalendarSubscriptionUrl = $derived(`https://www.google.com/calendar/render?cid=${encodeURIComponent(getWebCalUrl(postcode, uprn))}`);

</script>

<Card class="w-full max-w-md shadow-lg">
  <CardHeader>
    <CardTitle class="flex items-center gap-2 text-2xl">
      <CalendarDays class="h-6 w-6 text-primary" />
      Collection Schedule
    </CardTitle>
  </CardHeader>
  <CardContent class="space-y-4 text-left">
    <div class="rounded-lg bg-muted p-3 text-sm space-y-1">
      <p><strong>Postcode:</strong> <span class="uppercase">{postcode}</span></p>
      <p><strong>UPRN:</strong> <span class="font-mono text-xs">{uprn}</span></p>
    </div>
    {#if !postcode || !uprn}
      <p class="text-sm text-destructive">Missing lookup information. Please go back and try again.</p>
    {/if}

    <div class="space-y-4 text-center p-4">
      <p class="text-sm text-muted-foreground">
        Subscribe to get automatic garbage collection updates synced straight to your device.
      </p>

      <div>
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
  </CardContent>
</Card>