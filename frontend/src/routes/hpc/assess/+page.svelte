<script lang="ts">
  import { goto } from "$app/navigation";
  import { onMount } from "svelte";

  onMount(() => {
    goto("/hpc?tab=assess", { replaceState: true });
  });
</script>
