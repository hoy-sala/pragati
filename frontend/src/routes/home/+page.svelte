<script lang="ts">
  import { api } from "$lib/api/client.svelte";
  import { getAuthState } from "$lib/stores/auth.svelte";
  import { effectiveRole } from "$lib/utils/roles";
  import { goto } from "$app/navigation";
  import { onMount } from "svelte";
  import {
    Users,
    ClipboardCheck,
    Table,
    HelpCircle,
    ClipboardList,
    FileSpreadsheet,
    Heart,
    FileText,
    Award,
    Settings,
  } from "lucide-svelte";
  import PageHeader from "$lib/components/PageHeader.svelte";
  import EmptyState from "$lib/components/EmptyState.svelte";
  import type { User } from "$lib/types";

  const auth = getAuthState();
  let role = $derived(effectiveRole(auth.currentUser));
  let displayName = $derived(
    auth.currentUser ? (auth.currentUser as User).name?.split(" ")[0] || "there" : "there",
  );

  type PendingAssessment = {
    id: string;
    name: string;
    class_name: string;
    subject_name: string;
    max_marks: number;
    marks_count: number;
    student_count: number;
    due_date: string;
  };

  let pending = $state<PendingAssessment[]>([]);
  let loading = $state(true);

  const quickLinks = $derived(
    [
      { href: "/students", label: "Students", desc: "Enrolment & import", icon: Users, roles: ["admin", "principal", "teacher"] },
      { href: "/assessments", label: "Assessments", desc: "Create & publish", icon: ClipboardCheck, roles: ["admin", "principal", "teacher"] },
      { href: "/marks", label: "Marks Entry", desc: "Enter marks", icon: Table, roles: ["admin", "principal", "teacher"] },
      { href: "/questions", label: "Question Bank", desc: "Browse & import", icon: HelpCircle, roles: ["admin", "principal", "teacher"] },
      { href: "/quizzes", label: "Quizzes", desc: "Create & publish", icon: ClipboardList, roles: ["admin", "principal", "teacher"] },
      { href: "/hpc", label: "HPC Cards", desc: "Holistic progress", icon: FileSpreadsheet, roles: ["admin", "principal", "teacher"] },
      { href: "/mentors", label: "Mentors", desc: "Roster & logs", icon: Heart, roles: ["admin", "principal", "teacher", "special_educator"] },
      { href: "/reports", label: "Reports", desc: "Mark sheets & cards", icon: FileText, roles: ["admin", "principal", "teacher", "special_educator"] },
      { href: "/certificates", label: "Certificates", desc: "Events & printing", icon: Award, roles: ["admin"] },
      { href: "/settings", label: "Settings", desc: "Users & setup", icon: Settings, roles: ["admin"] },
    ].filter((l) => l.roles.includes(role)),
  );

  onMount(async () => {
    if (role === "student") {
      goto("/reports", { replaceState: true });
      return;
    }
    const res = await api<{
      pending_assessments: PendingAssessment[];
    }>("GET", "/dashboard/staff");
    if (res.data) pending = res.data.pending_assessments ?? [];
    loading = false;
  });
</script>

<svelte:head><title>Home — Pragati</title></svelte:head>

<div class="max-w-7xl mx-auto space-y-6">
  <PageHeader
    title={`Welcome back, ${displayName}`}
    subtitle="Your pending work and quick actions"
  />

  <div class="bg-white rounded-xl border border-slate-200 overflow-hidden">
    <div class="px-5 py-3 border-b border-slate-200 flex items-center justify-between">
      <h2 class="text-sm font-semibold text-slate-800">Pending assessments</h2>
      <a href="/assessments" class="text-xs text-primary-600 hover:text-primary-700 font-medium">View all</a>
    </div>
    {#if loading}
      <div class="p-8 text-center text-sm text-slate-400">Loading...</div>
    {:else if pending.length === 0}
      <EmptyState
        icon={ClipboardCheck}
        title="Nothing pending."
        hint="Draft assessments needing marks will show up here."
      />
    {:else}
      <div class="divide-y divide-slate-100">
        {#each pending as p (p.id)}
          <a
            href="/marks?assessment={p.id}"
            class="flex items-center justify-between gap-4 px-5 py-3 hover:bg-slate-50 transition-colors"
          >
            <div class="min-w-0">
              <div class="text-sm font-medium text-slate-800 truncate">{p.name || "Untitled assessment"}</div>
              <div class="text-xs text-slate-500 mt-0.5">{p.class_name} · {p.subject_name}</div>
            </div>
            <div class="text-xs text-slate-500 shrink-0">
              {p.marks_count}/{p.student_count} marked
            </div>
          </a>
        {/each}
      </div>
    {/if}
  </div>

  <div>
    <h2 class="text-sm font-semibold text-slate-800 mb-3">Quick actions</h2>
    <div class="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-5 gap-3">
      {#each quickLinks as l (l.href)}
        <a
          href={l.href}
          class="bg-white rounded-xl border border-slate-200 p-4 hover:border-primary-300 hover:shadow-sm transition-all group"
        >
          <l.icon size={20} class="text-primary-600 mb-2" />
          <div class="text-sm font-medium text-slate-800 group-hover:text-primary-700">{l.label}</div>
          <div class="text-xs text-slate-400 mt-0.5">{l.desc}</div>
        </a>
      {/each}
    </div>
  </div>
</div>
