<script lang="ts">
  import { api } from "$lib/api/client.svelte";
  import type { QuizListItem } from "$lib/types";
  import { onMount } from "svelte";
  import { goto } from "$app/navigation";
  import { Plus, ClipboardList } from "lucide-svelte";
  import Button from "$lib/components/Button.svelte";
  import PageHeader from "$lib/components/PageHeader.svelte";
  import EmptyState from "$lib/components/EmptyState.svelte";
  import Select from "$lib/components/Select.svelte";
  import SearchFilter from "$lib/components/SearchFilter.svelte";
  import Pagination from "$lib/components/Pagination.svelte";
  import { toast } from "$lib/stores/toast.svelte";

  let allQuizzes = $state<QuizListItem[]>([]);
  let loading = $state(true);
  let deleting = $state<string | null>(null);
  let search = $state("");
  let filterStatus = $state("");
  let page = $state(1);
  const pageSize = 20;

  const statusOptions = [
    { id: "", name: "All Status" },
    { id: "published", name: "Published" },
    { id: "draft", name: "Draft" },
  ];

  let filteredQuizzes = $derived(
    allQuizzes.filter((q) => {
      if (filterStatus === "published" && !q.is_published) return false;
      if (filterStatus === "draft" && q.is_published) return false;
      if (!search.trim()) return true;
      const q2 = search.toLowerCase();
      return (
        q.title?.toLowerCase().includes(q2) ||
        q.description?.toLowerCase().includes(q2)
      );
    }),
  );

  let paginatedQuizzes = $derived(
    filteredQuizzes.slice((page - 1) * pageSize, page * pageSize),
  );
  let totalQuizzes = $derived(filteredQuizzes.length);

  function onPageChange(p: number) {
    page = p;
  }
  function resetPage() {
    page = 1;
  }

  onMount(loadQuizzes);

  async function loadQuizzes() {
    loading = true;
    const res = await api<QuizListItem[]>("GET", "/quizzes");
    if (res.data) allQuizzes = res.data;
    loading = false;
  }

  async function publish(id: string) {
    const res = await api("POST", `/quizzes/${id}/publish`);
    if (res.error) {
      toast(res.error.message, "error");
      return;
    }
    toast("Quiz published", "success");
    loadQuizzes();
  }

  async function remove(id: string) {
    if (!confirm("Delete this quiz?")) return;
    deleting = id;
    const res = await api("DELETE", `/quizzes/${id}`);
    if (res.error) {
      toast(res.error.message, "error");
      deleting = null;
      return;
    }
    toast("Quiz deleted", "success");
    deleting = null;
    loadQuizzes();
  }

  function targetLabel(t: string): string {
    return t === "student" ? "Students" : "Staff";
  }
</script>

<div class="space-y-6">
  <PageHeader title="Quizzes" subtitle="{totalQuizzes} quizzes">
    {#snippet actions()}
      <Button icon={Plus} onclick={() => goto("/quizzes/create")}>
        Create Quiz
      </Button>
    {/snippet}
  </PageHeader>

  <div class="bg-white rounded-xl border border-slate-200">
    <div class="p-4 border-b border-slate-200 flex flex-wrap gap-3">
      <div class="flex-1 min-w-48">
        <SearchFilter
          bind:value={search}
          placeholder="Search quizzes..."
          onInput={resetPage}
        />
      </div>
      <div class="w-36">
        <Select
          bind:value={filterStatus}
          options={statusOptions}
          placeholder="All Status"
          onselect={resetPage}
        />
      </div>
    </div>

    {#if loading}
      <div class="p-8 text-center text-sm text-slate-400">Loading...</div>
    {:else if totalQuizzes === 0}
      <EmptyState
        icon={ClipboardList}
        title="No quizzes found."
        hint="Try adjusting filters, or create your first quiz."
      >
        {#snippet action()}
          <Button size="sm" icon={Plus} onclick={() => goto("/quizzes/create")}>Create Quiz</Button>
        {/snippet}
      </EmptyState>
    {:else}
      <div class="divide-y divide-slate-100">
        {#each paginatedQuizzes as q (q.id)}
          <div class="p-4 hover:bg-slate-50 transition-colors">
            <div class="flex items-center justify-between gap-4">
              <div class="flex-1 min-w-0">
                <div class="flex items-center gap-2">
                  <h3 class="text-sm font-medium text-slate-900">{q.title}</h3>
                  <span
                    class="text-xs px-2 py-0.5 rounded-full {q.is_published
                      ? 'bg-green-100 text-green-700'
                      : 'bg-yellow-100 text-yellow-700'}"
                  >
                    {q.is_published ? "Published" : "Draft"}
                  </span>
                </div>
                <p class="text-xs text-slate-500 mt-1 line-clamp-1">
                  {q.description || "No description"}
                </p>
                <div class="flex flex-wrap gap-3 mt-1.5 text-xs text-slate-400">
                  <span>{targetLabel(q.target_type)}</span>
                  <span>{q.question_count} questions</span>
                  <span>{q.attempt_count} attempts</span>
                  <span>Pass: {q.pass_pct}%</span>
                  <span>by {q.created_by_name}</span>
                </div>
              </div>
              <div class="flex gap-1.5 shrink-0">
                {#if !q.is_published}
                  <Button
                    size="sm"
                    onclick={() => publish(q.id)}
                    disabled={deleting !== null}
                  >
                    Publish
                  </Button>
                  <Button
                    size="sm"
                    variant="secondary"
                    onclick={() => goto(`/quizzes/${q.id}/edit`)}
                    disabled={deleting !== null}
                  >
                    Edit
                  </Button>
                  <Button
                    size="sm"
                    variant="danger"
                    onclick={() => remove(q.id)}
                    disabled={deleting !== null}
                    loading={deleting === q.id}
                  >
                    {deleting === q.id ? "Deleting..." : "Delete"}
                  </Button>
                {/if}
              </div>
            </div>
          </div>
        {/each}
      </div>
      <Pagination
        total={totalQuizzes}
        {pageSize}
        {page}
        onChange={onPageChange}
      />
    {/if}
  </div>
</div>
