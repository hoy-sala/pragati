<script lang="ts">
  import "../app.css";
  import "katex/dist/katex.min.css";
  import { onMount } from "svelte";
  import { initAuth, getAuthState } from "$lib/stores/auth.svelte";
  import { effectiveRole } from "$lib/utils/roles";
  import { goto } from "$app/navigation";
  import { page } from "$app/stores";
  import Sidebar from "$lib/components/layout/Sidebar.svelte";
  import MobileTopBar from "$lib/components/layout/MobileTopBar.svelte";
  import BottomNav from "$lib/components/layout/BottomNav.svelte";
  import NavSheet from "$lib/components/layout/NavSheet.svelte";
  import Toast from "$lib/components/Toast.svelte";
  import { LoaderCircle } from "lucide-svelte";

  let { children } = $props();

  const auth = getAuthState();
  let navOpen = $state(false);

  const publicRoutes = ["/login", "/timetable", "/play"];

  function isFullscreenRoute(path: string): boolean {
    return (
      path.startsWith("/quizzes/take") ||
      path.startsWith("/quizzes/results") ||
      path.startsWith("/certificates/print") ||
      path.startsWith("/play")
    );
  }

  let fullscreen = $derived(isFullscreenRoute($page.url.pathname));

  onMount(() => {
    initAuth();
  });

  $effect(() => {
    if (!auth.isLoading && !auth.isAuthenticated) {
      const path = $page.url.pathname;
      if (!publicRoutes.includes(path) && path !== "/") {
        goto("/login");
      }
    }
    if (
      !auth.isLoading &&
      auth.isAuthenticated &&
      $page.url.pathname === "/login"
    ) {
      goto(
        effectiveRole(auth.currentUser) === "student" ? "/reports" : "/home",
      );
    }
  });
</script>

{#if auth.isLoading}
  <div class="flex h-dvh items-center justify-center">
    <div class="flex flex-col items-center gap-2 text-slate-400">
      <LoaderCircle class="animate-spin text-primary-500" size={24} />
      <span class="text-sm">Loading...</span>
    </div>
  </div>
{:else if auth.isAuthenticated}
  <div class="flex h-dvh overflow-hidden">
    {#if !fullscreen}
      <Sidebar />
    {/if}
    <div class="flex-1 flex flex-col min-w-0 min-h-0">
      {#if !fullscreen}
        <MobileTopBar onmenu={() => (navOpen = true)} />
      {/if}
      <main
        class="flex-1 overflow-y-auto p-4 md:p-6"
        class:p-0={fullscreen}
        class:md:p-0={fullscreen}
      >
        {@render children()}
      </main>
      {#if !fullscreen}
        <BottomNav onmore={() => (navOpen = true)} />
      {/if}
    </div>
    <NavSheet bind:open={navOpen} />
  </div>
{:else}
  <main class="min-h-dvh">
    {@render children()}
  </main>
{/if}
<Toast />
