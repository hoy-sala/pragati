<script lang="ts">
  let { code }: { code: "" | "LC" | "NT" | "VU" | "EN" | "CR" } = $props();
  const labels: Record<string, string> = {
    LC: "Least Concern",
    NT: "Near Threatened",
    VU: "Vulnerable",
    EN: "Endangered",
    CR: "Critically Endangered"
  };
  const styles: Record<string, string> = {
    LC: "bg-emerald-50 text-emerald-700",
    NT: "bg-yellow-50 text-yellow-800",
    VU: "bg-orange-50 text-orange-700",
    EN: "bg-red-50 text-red-700",
    CR: "bg-rose-100 text-rose-800"
  };
</script>

{#if code}
  <span
    class="text-[11px] font-semibold px-2 py-0.5 rounded-full {styles[code]}"
    title="IUCN Red List status: {labels[code]}"
  >
    {labels[code]}
  </span>
{/if}
